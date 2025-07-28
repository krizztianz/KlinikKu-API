package controllers

import (
	"KlinikKu/dto"
	"KlinikKu/middleware"
	"KlinikKu/utils"
	"database/sql"
	"fmt"
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

	createdBy := middleware.GetUsername(c)
	var rekamMedisID int

	err := utils.RunTx(db, func(tx *sql.Tx) error {
		var status string
		err := tx.QueryRow("SELECT status FROM kunjungan WHERE kunjungan_id = $1", input.KunjunganID).Scan(&status)
		if err != nil {
			return err
		}
		if status != "pemeriksaan" {
			return fmt.Errorf("Kunjungan belum dalam status 'pemeriksaan'")
		}

		err = tx.QueryRow(`
			INSERT INTO rekam_medis (kunjungan_id, anamnesa, pemeriksaan_fisik, catatan, created_by)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING rekam_medis_id
		`, input.KunjunganID, input.Anamnesa, input.PemeriksaanFisik, input.Catatan, createdBy).Scan(&rekamMedisID)
		if err != nil {
			return err
		}

		for _, diagnosaID := range input.DiagnosaIDs {
			_, err := tx.Exec(`
				INSERT INTO rekam_medis_diagnosa (rekam_medis_id, diagnosa_id, created_by)
				VALUES ($1, $2, $3)
			`, rekamMedisID, diagnosaID, createdBy)
			if err != nil {
				return err
			}
		}

		for _, tindakan := range input.Tindakans {
			_, err := tx.Exec(`
				INSERT INTO rekam_medis_tindakan (rekam_medis_id, tindakan_id, jumlah, biaya_aktual, catatan, created_by)
				VALUES ($1, $2, $3, $4, $5, $6)
			`, rekamMedisID, tindakan.TindakanID, tindakan.Jumlah, tindakan.BiayaAktual, tindakan.Catatan, createdBy)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Rekam medis berhasil disimpan",
		"rekam_medis_id": rekamMedisID,
	})
}

func GetRekamMedisByID(c *gin.Context) {
	rekamMedisID := c.Param("id")

	row := db.QueryRow(`
		SELECT rm.rekam_medis_id, rm.kunjungan_id, rm.anamnesa, rm.pemeriksaan_fisik, rm.catatan, rm.created_at, rm.created_by
		FROM rekam_medis rm
		WHERE rm.rekam_medis_id = $1
	`, rekamMedisID)

	var result map[string]interface{}
	var rmID, kunjunganID int
	var anamnesa, fisik, catatan, createdAt, createdBy string

	err := row.Scan(&rmID, &kunjunganID, &anamnesa, &fisik, &catatan, &createdAt, &createdBy)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rekam medis tidak ditemukan"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	// Query diagnosa
	diagnosaRows, err := db.Query(`
		SELECT d.kode, d.nama
		FROM rekam_medis_diagnosa rmd
		JOIN diagnosa d ON rmd.diagnosa_id = d.diagnosa_id
		WHERE rmd.rekam_medis_id = $1
	`, rekamMedisID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil diagnosa"})
		return
	}
	defer diagnosaRows.Close()

	var diagnosaList []map[string]string
	for diagnosaRows.Next() {
		var kode, nama string
		if err := diagnosaRows.Scan(&kode, &nama); err == nil {
			diagnosaList = append(diagnosaList, map[string]string{
				"kode": kode, "nama": nama,
			})
		}
	}

	// Query tindakan
	tindakanRows, err := db.Query(`
		SELECT t.nama_tindakan, rmt.jumlah, rmt.biaya_aktual, rmt.catatan
		FROM rekam_medis_tindakan rmt
		JOIN tindakan t ON rmt.tindakan_id = t.tindakan_id
		WHERE rmt.rekam_medis_id = $1
	`, rekamMedisID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil tindakan"})
		return
	}
	defer tindakanRows.Close()

	var tindakanList []map[string]interface{}
	for tindakanRows.Next() {
		var nama string
		var jumlah int
		var biaya int
		var catatan string
		if err := tindakanRows.Scan(&nama, &jumlah, &biaya, &catatan); err == nil {
			tindakanList = append(tindakanList, map[string]interface{}{
				"nama_tindakan": nama,
				"jumlah":        jumlah,
				"biaya_aktual":  biaya,
				"catatan":       catatan,
			})
		}
	}

	result = map[string]interface{}{
		"rekam_medis_id":    rmID,
		"kunjungan_id":      kunjunganID,
		"anamnesa":          anamnesa,
		"pemeriksaan_fisik": fisik,
		"catatan":           catatan,
		"created_at":        createdAt,
		"created_by":        createdBy,
		"diagnosa":          diagnosaList,
		"tindakan":          tindakanList,
	}

	c.JSON(http.StatusOK, result)
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
