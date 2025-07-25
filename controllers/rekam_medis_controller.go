package controllers

import (
	"KlinikKu/dto"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateRekamMedis(c *gin.Context) {
	var input struct {
		dto.CreateRekamMedisRequest
		dto.AddDiagnosaRequest
		dto.AddTindakanRequest
	}

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

	var status string
	err = tx.QueryRow("SELECT status FROM kunjungan WHERE kunjungan_id = $1", input.KunjunganID).Scan(&status)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kunjungan tidak ditemukan"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil status kunjungan"})
		return
	}

	if status != "diperiksa" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kunjungan belum dalam status 'diperiksa'"})
		return
	}

	var rekamMedisID int
	err = tx.QueryRow(`
		INSERT INTO rekam_medis (kunjungan_id, anamnesa, pemeriksaan_fisik, catatan)
		VALUES ($1, $2, $3, $4)
		RETURNING rekam_medis_id
	`, input.KunjunganID, input.Anamnesa, input.PemeriksaanFisik, input.Catatan).Scan(&rekamMedisID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan rekam medis"})
		return
	}

	for _, diagnosaID := range input.DiagnosaIDs {
		_, err := tx.Exec(`
			INSERT INTO rekam_medis_diagnosa (rekam_medis_id, diagnosa_id)
			VALUES ($1, $2)
		`, rekamMedisID, diagnosaID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan diagnosa"})
			return
		}
	}

	for _, tindakan := range input.Tindakans {
		_, err := tx.Exec(`
			INSERT INTO rekam_medis_tindakan (rekam_medis_id, tindakan_id, jumlah, biaya_aktual, catatan)
			VALUES ($1, $2, $3, $4, $5)
		`, rekamMedisID, tindakan.TindakanID, tindakan.Jumlah, tindakan.BiayaAktual, tindakan.Catatan)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan tindakan"})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal commit transaksi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Rekam medis berhasil disimpan",
		"rekam_medis_id": rekamMedisID,
	})
}

func GetRekamMedisByID(c *gin.Context) {
	id := c.Param("id")

	// Ambil data utama rekam medis
	var rekamMedis struct {
		RekamMedisID     int    `json:"rekam_medis_id"`
		KunjunganID      int    `json:"kunjungan_id"`
		Anamnesa         string `json:"anamnesa"`
		PemeriksaanFisik string `json:"pemeriksaan_fisik"`
		Catatan          string `json:"catatan"`
		TanggalDibuat    string `json:"tanggal_dibuat"`
	}

	err := db.QueryRow(`
		SELECT rekam_medis_id, kunjungan_id, anamnesa, pemeriksaan_fisik, catatan, created_at
		FROM rekam_medis WHERE rekam_medis_id = $1
	`, id).Scan(&rekamMedis.RekamMedisID, &rekamMedis.KunjunganID,
		&rekamMedis.Anamnesa, &rekamMedis.PemeriksaanFisik,
		&rekamMedis.Catatan, &rekamMedis.TanggalDibuat)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rekam medis tidak ditemukan"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	// Diagnosa
	diagnosaRows, _ := db.Query(`
		SELECT d.diagnosa_id, d.kode_icd10, d.nama_diagnosa
		FROM rekam_medis_diagnosa rmd
		JOIN diagnosa d ON d.diagnosa_id = rmd.diagnosa_id
		WHERE rmd.rekam_medis_id = $1
	`, id)
	defer diagnosaRows.Close()

	var diagnosaList []map[string]interface{}
	for diagnosaRows.Next() {
		var dID int
		var kode, nama string
		diagnosaRows.Scan(&dID, &kode, &nama)
		diagnosaList = append(diagnosaList, map[string]interface{}{
			"diagnosa_id":   dID,
			"kode_icd10":    kode,
			"nama_diagnosa": nama,
		})
	}

	// Tindakan
	tindakanRows, _ := db.Query(`
		SELECT t.tindakan_id, t.nama_tindakan, rt.jumlah, rt.biaya_aktual, rt.catatan
		FROM rekam_medis_tindakan rt
		JOIN tindakan t ON t.tindakan_id = rt.tindakan_id
		WHERE rt.rekam_medis_id = $1
	`, id)
	defer tindakanRows.Close()

	var tindakanList []map[string]interface{}
	for tindakanRows.Next() {
		var tid int
		var nama string
		var jumlah int
		var biaya float64
		var catatan string
		tindakanRows.Scan(&tid, &nama, &jumlah, &biaya, &catatan)
		tindakanList = append(tindakanList, map[string]interface{}{
			"tindakan_id":   tid,
			"nama_tindakan": nama,
			"jumlah":        jumlah,
			"biaya_aktual":  biaya,
			"catatan":       catatan,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"rekam_medis": rekamMedis,
		"diagnosa":    diagnosaList,
		"tindakan":    tindakanList,
	})
}

func GetRekamMedisByKunjungan(c *gin.Context) {
	kunjunganID := c.Param("id")

	rows, err := db.Query(`
		SELECT rm.rekam_medis_id, rm.anamnesa, rm.pemeriksaan_fisik, rm.catatan, rm.created_at
		FROM rekam_medis rm
		WHERE rm.kunjungan_id = $1
	`, kunjunganID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	var rekamMedisList []map[string]interface{}

	for rows.Next() {
		var rmID int
		var anamnesa, fisik, catatan, createdAt string
		err := rows.Scan(&rmID, &anamnesa, &fisik, &catatan, &createdAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data"})
			return
		}
		rekamMedisList = append(rekamMedisList, map[string]interface{}{
			"rekam_medis_id":    rmID,
			"anamnesa":          anamnesa,
			"pemeriksaan_fisik": fisik,
			"catatan":           catatan,
			"tanggal_dibuat":    createdAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"kunjungan_id": kunjunganID,
		"rekam_medis":  rekamMedisList,
	})
}
