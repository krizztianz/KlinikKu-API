package controllers

import (
	"KlinikKu/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateSpesialisasi(c *gin.Context) {
	var input dto.CreateSpesialisasiRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec(`
		INSERT INTO spesialisasi (kode_spesialisasi, nama_spesialisasi, deskripsi)
		VALUES ($1, $2, $3)
	`, input.KodeSpesialisasi, input.NamaSpesialisasi, input.Deskripsi)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan spesialisasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Spesialisasi berhasil ditambahkan"})
}

func GetAllSpesialisasi(c *gin.Context) {
	rows, err := db.Query(`
		SELECT spesialisasi_id, kode_spesialisasi, nama_spesialisasi, deskripsi
		FROM spesialisasi ORDER BY spesialisasi_id`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data spesialisasi"})
		return
	}
	defer rows.Close()

	var result []dto.SpesialisasiResponse
	for rows.Next() {
		var s dto.SpesialisasiResponse
		err := rows.Scan(&s.SpesialisasiID, &s.KodeSpesialisasi, &s.NamaSpesialisasi, &s.Deskripsi)
		if err == nil {
			result = append(result, s)
		}
	}

	c.JSON(http.StatusOK, result)
}

func UpdateSpesialisasi(c *gin.Context) {
	id := c.Param("id")
	var input dto.CreateSpesialisasiRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := db.Exec(`
		UPDATE spesialisasi SET 
			kode_spesialisasi = $1,
			nama_spesialisasi = $2,
			deskripsi = $3
		WHERE spesialisasi_id = $4
	`, input.KodeSpesialisasi, input.NamaSpesialisasi, input.Deskripsi, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update spesialisasi"})
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Spesialisasi tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Spesialisasi diperbarui"})
}

func DeleteSpesialisasi(c *gin.Context) {
	id := c.Param("id")
	res, err := db.Exec(`DELETE FROM spesialisasi WHERE spesialisasi_id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus spesialisasi"})
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Spesialisasi tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Spesialisasi berhasil dihapus"})
}
