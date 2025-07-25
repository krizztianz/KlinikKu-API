package controllers

import (
	"KlinikKu/dto"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateResep(c *gin.Context) {
	var input dto.CreateResepRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuka transaksi DB"})
		return
	}
	defer tx.Rollback()

	// Validasi rekam_medis_id
	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM rekam_medis WHERE rekam_medis_id = $1)", input.RekamMedisID).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rekam medis tidak ditemukan"})
		return
	}

	// Insert ke tabel resep
	var resepID int
	err = tx.QueryRow(`
		INSERT INTO resep (rekam_medis_id)
		VALUES ($1)
		RETURNING resep_id
	`, input.RekamMedisID).Scan(&resepID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan resep"})
		return
	}

	// Insert ke resep_detail
	for _, obat := range input.Obats {
		_, err := tx.Exec(`
			INSERT INTO resep_detail (resep_id, obat_id, jumlah, dosis, keterangan)
			VALUES ($1, $2, $3, $4, $5)
		`, resepID, obat.ObatID, obat.Jumlah, obat.Dosis, obat.Keterangan)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan detail resep"})
			return
		}
	}

	// Commit
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal commit resep"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Resep berhasil disimpan", "resep_id": resepID})
}

func GetResepByID(c *gin.Context) {
	id := c.Param("id")

	var resep struct {
		ResepID      int    `json:"resep_id"`
		RekamMedisID int    `json:"rekam_medis_id"`
		Tanggal      string `json:"tanggal_dibuat"`
	}

	err := db.QueryRow(`
		SELECT resep_id, rekam_medis_id, created_at FROM resep WHERE resep_id = $1
	`, id).Scan(&resep.ResepID, &resep.RekamMedisID, &resep.Tanggal)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	obatRows, _ := db.Query(`
		SELECT o.obat_id, o.nama_obat, rd.jumlah, rd.dosis, rd.keterangan
		FROM resep_detail rd
		JOIN obat o ON o.obat_id = rd.obat_id
		WHERE rd.resep_id = $1
	`, id)
	defer obatRows.Close()

	var obats []map[string]interface{}
	for obatRows.Next() {
		var oid int
		var nama, dosis, keterangan string
		var jumlah int
		obatRows.Scan(&oid, &nama, &jumlah, &dosis, &keterangan)
		obats = append(obats, map[string]interface{}{
			"obat_id":    oid,
			"nama_obat":  nama,
			"jumlah":     jumlah,
			"dosis":      dosis,
			"keterangan": keterangan,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"resep": resep,
		"obats": obats,
	})
}

func GetResepByRekamMedisID(c *gin.Context) {
	rekamMedisID := c.Param("id")

	var resep struct {
		ResepID      int    `json:"resep_id"`
		RekamMedisID int    `json:"rekam_medis_id"`
		Tanggal      string `json:"tanggal_dibuat"`
	}

	err := db.QueryRow(`
		SELECT resep_id, rekam_medis_id, created_at
		FROM resep
		WHERE rekam_medis_id = $1
	`, rekamMedisID).Scan(&resep.ResepID, &resep.RekamMedisID, &resep.Tanggal)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada resep untuk rekam medis ini"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data resep"})
		return
	}

	// Ambil detail obat
	rows, _ := db.Query(`
		SELECT o.obat_id, o.nama_obat, rd.jumlah, rd.dosis, rd.keterangan
		FROM resep_detail rd
		JOIN obat o ON o.obat_id = rd.obat_id
		WHERE rd.resep_id = $1
	`, resep.ResepID)
	defer rows.Close()

	var obats []map[string]interface{}
	for rows.Next() {
		var oid int
		var nama, dosis, keterangan string
		var jumlah int
		rows.Scan(&oid, &nama, &jumlah, &dosis, &keterangan)
		obats = append(obats, map[string]interface{}{
			"obat_id":    oid,
			"nama_obat":  nama,
			"jumlah":     jumlah,
			"dosis":      dosis,
			"keterangan": keterangan,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"resep": resep,
		"obats": obats,
	})
}
