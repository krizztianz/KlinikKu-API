package controllers

import (
	"KlinikKu/dto"
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

	var kunjunganID int
	err := db.QueryRow(`
		INSERT INTO kunjungan (pasien_id, dokter_id, tanggal_kunjungan, keluhan, jenis_kunjungan, status, prioritas)
		VALUES ($1, $2, $3, $4, $5, 'registrasi', COALESCE($6, 'normal'))
		RETURNING kunjungan_id
	`, input.PasienID, input.DokterID, tgl, input.Keluhan, input.JenisKunjungan, input.Prioritas).
		Scan(&kunjunganID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan kunjungan"})
		return
	}

	resp := dto.KunjunganResponse{
		KunjunganID:      kunjunganID,
		PasienID:         input.PasienID,
		DokterID:         input.DokterID,
		TanggalKunjungan: tgl.Format("2006-01-02 15:04:05"),
		Keluhan:          input.Keluhan,
		JenisKunjungan:   input.JenisKunjungan,
		Status:           "registrasi",
		Prioritas:        input.Prioritas,
	}

	c.JSON(http.StatusOK, resp)
}

func GetKunjunganList(c *gin.Context) {
	status := c.Query("status")
	rows, err := db.Query(`
		SELECT kunjungan_id, pasien_id, dokter_id, tanggal_kunjungan, keluhan, jenis_kunjungan, status, prioritas
		FROM kunjungan
		WHERE ($1 = '' OR status = $1)
		ORDER BY tanggal_kunjungan DESC`, status)
	if err != nil {
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

	// Cek status sebelumnya
	var currentStatus string
	err := db.QueryRow(`SELECT status FROM kunjungan WHERE kunjungan_id = $1`, kunjunganID).
		Scan(&currentStatus)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kunjungan tidak ditemukan"})
		return
	}

	// Validasi transisi status yang sah
	allowedTransitions := map[string][]string{
		"registrasi":     {"pemeriksaan", "batal"},
		"pemeriksaan":    {"menunggu_resep", "batal"},
		"menunggu_resep": {"selesai", "batal"},
	}

	allowedNext, exists := allowedTransitions[currentStatus]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transisi dari status ini tidak diperbolehkan"})
		return
	}

	valid := false
	for _, s := range allowedNext {
		if newStatus == s {
			valid = true
			break
		}
	}

	if !valid {
		msg := fmt.Sprintf("Transisi status dari '%s' ke '%s' tidak diperbolehkan", currentStatus, newStatus)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	_, err = db.Exec(`UPDATE kunjungan SET status = $1 WHERE kunjungan_id = $2`, newStatus, kunjunganID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update status kunjungan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status kunjungan diperbarui"})
}
