package controllers

import (
	"KlinikKu/dto"
	"KlinikKu/middleware"
	"KlinikKu/utils"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllResep(c *gin.Context) {
	rows, err := db.Query(`
		SELECT resep_id, rekam_medis_id, created_at, status
		FROM resep
		ORDER BY created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data resep"})
		return
	}
	defer rows.Close()

	var resepList []map[string]interface{}

	for rows.Next() {
		var resepID, rekamMedisID int
		var createdAt, status string

		if err := rows.Scan(&resepID, &rekamMedisID, &createdAt, &status); err != nil {
			continue
		}

		obatRows, _ := db.Query(`
			SELECT o.obat_id, o.nama_obat, rd.jumlah, rd.dosis, rd.keterangan
			FROM resep_detail rd
			JOIN obat o ON o.obat_id = rd.obat_id
			WHERE rd.resep_id = $1
		`, resepID)

		var obats []map[string]interface{}
		for obatRows.Next() {
			var oid, jumlah int
			var nama, dosis, keterangan string
			obatRows.Scan(&oid, &nama, &jumlah, &dosis, &keterangan)

			obats = append(obats, map[string]interface{}{
				"obat_id":    oid,
				"nama_obat":  nama,
				"jumlah":     jumlah,
				"dosis":      dosis,
				"keterangan": keterangan,
			})
		}
		obatRows.Close()

		resepList = append(resepList, map[string]interface{}{
			"resep_id":       resepID,
			"rekam_medis_id": rekamMedisID,
			"tanggal":        createdAt,
			"status":         status,
			"obats":          obats,
		})
	}

	c.JSON(http.StatusOK, resepList)
}

func CreateResep(c *gin.Context) {
	var input dto.CreateResepRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdBy := middleware.GetUsername(c)
	var resepID int

	err := utils.RunTx(db, func(tx *sql.Tx) error {
		var exists bool
		err := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM rekam_medis WHERE rekam_medis_id = $1)", input.RekamMedisID).Scan(&exists)
		if err != nil {
			return err
		}
		if !exists {
			return sql.ErrNoRows
		}

		// Insert resep
		err = tx.QueryRow(`
			INSERT INTO resep (rekam_medis_id, status, created_by)
			VALUES ($1, 'menunggu'::resep_status, $2)
			RETURNING resep_id
		`, input.RekamMedisID, createdBy).Scan(&resepID)
		if err != nil {
			return err
		}

		// Insert resep_detail
		for _, obat := range input.Obats {
			_, err := tx.Exec(`
				INSERT INTO resep_detail (resep_id, obat_id, jumlah, dosis, keterangan, created_by)
				VALUES ($1, $2, $3, $4, $5, $6)
			`, resepID, obat.ObatID, obat.Jumlah, obat.Dosis, obat.Keterangan, createdBy)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Rekam medis tidak ditemukan"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan resep"})
		}
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

func GetResepMenunggu(c *gin.Context) {
	rows, err := db.Query(`
		SELECT 
			r.resep_id,
			r.rekam_medis_id,
			rm.kunjungan_id,
			p.nama AS nama_pasien,
			r.created_at
		FROM resep r
		JOIN rekam_medis rm ON r.rekam_medis_id = rm.rekam_medis_id
		JOIN kunjungan k ON rm.kunjungan_id = k.kunjungan_id
		JOIN pasien p ON k.pasien_id = p.pasien_id
		WHERE r.status != 'selesai'::resep_status
		ORDER BY r.created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil resep menunggu"})
		return
	}
	defer rows.Close()

	var daftar []map[string]interface{}
	for rows.Next() {
		var resepID, rekamMedisID, kunjunganID int
		var namaPasien, createdAt string

		err := rows.Scan(&resepID, &rekamMedisID, &kunjunganID, &namaPasien, &createdAt)
		if err != nil {
			continue
		}

		daftar = append(daftar, map[string]interface{}{
			"resep_id":       resepID,
			"rekam_medis_id": rekamMedisID,
			"kunjungan_id":   kunjunganID,
			"nama_pasien":    namaPasien,
			"tanggal":        createdAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"resep_menunggu": daftar})
}
