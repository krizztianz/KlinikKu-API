package controllers

import (
	"KlinikKu/dto"
	"KlinikKu/utils"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateKunjungan(c *gin.Context) {
	var input dto.CreateKunjunganRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tgl := time.Now()
	createdBy := c.GetString("username")

	var kunjunganID int
	err := utils.RunTx(db, func(tx *sql.Tx) error {
		return tx.QueryRow(`
			INSERT INTO kunjungan (
				pasien_id, dokter_id, tanggal_kunjungan,
				keluhan, jenis_kunjungan, status, prioritas,
				created_by
			)
			VALUES (
				$1, $2, $3, $4,
				$5::visit_type,
				'registrasi'::visit_status,
				COALESCE($6, 'normal')::visit_priority,
				$7
			)
			RETURNING kunjungan_id
		`, input.PasienID, input.DokterID, tgl, input.Keluhan, input.JenisKunjungan, input.Prioritas, createdBy).
			Scan(&kunjunganID)
	})

	if err != nil {
		fmt.Println("Input:", input)
		fmt.Println("Error tx:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan kunjungan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Kunjungan berhasil dibuat", "kunjungan_id": kunjunganID})
}

func GetKunjunganList(c *gin.Context) {
	status := c.Query("status")
	rows, err := db.Query(`
		SELECT kunjungan_id, pasien_id, dokter_id, tanggal_kunjungan, keluhan, jenis_kunjungan, status, prioritas
		FROM kunjungan
		WHERE ($1 = '' OR status = $1::visit_status)
		ORDER BY tanggal_kunjungan DESC`, status)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data kunjungan"})
		return
	}
	defer rows.Close()

	var results []dto.KunjunganResponse
	for rows.Next() {
		var k dto.KunjunganResponse
		err := rows.Scan(&k.KunjunganID, &k.PasienID, &k.DokterID, &k.TanggalKunjungan, &k.Keluhan, &k.JenisKunjungan, &k.Status, &k.Prioritas)
		if err != nil {
			continue
		}
		results = append(results, k)
	}

	c.JSON(http.StatusOK, results)
}

func UpdateKunjunganStatus(c *gin.Context) {
	kunjunganID := c.Param("id")
	var input dto.UpdateStatusKunjunganRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status wajib diisi"})
		return
	}

	newStatus := input.Status
	modifiedBy := c.GetString("username")

	err := utils.RunTx(db, func(tx *sql.Tx) error {
		var currentStatus string
		err := tx.QueryRow(`SELECT status FROM kunjungan WHERE kunjungan_id = $1`, kunjunganID).
			Scan(&currentStatus)
		if err != nil {
			return err
		}

		allowedTransitions := map[string][]string{
			"registrasi":     {"pemeriksaan", "batal"},
			"pemeriksaan":    {"menunggu_resep", "batal"},
			"menunggu_resep": {"selesai", "batal"},
		}

		allowedNext, exists := allowedTransitions[currentStatus]
		if !exists {
			return fmt.Errorf("Transisi dari status ini tidak diperbolehkan")
		}

		valid := false
		for _, s := range allowedNext {
			if newStatus == s {
				valid = true
				break
			}
		}

		if !valid {
			return fmt.Errorf("Transisi status tidak valid")
		}

		_, err = tx.Exec(`
			UPDATE kunjungan 
			SET status = $1::visit_status, modified_at = CURRENT_TIMESTAMP, modified_by = $2
			WHERE kunjungan_id = $3
		`, newStatus, modifiedBy, kunjunganID)

		return err
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status kunjungan berhasil diupdate"})
}
