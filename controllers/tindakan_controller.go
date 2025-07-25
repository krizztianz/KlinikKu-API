package controllers

import (
	"KlinikKu/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateTindakan(c *gin.Context) {
	var input dto.CreateTindakanRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec(`
		INSERT INTO tindakan (kode_icd9, nama_tindakan, deskripsi, biaya_dasar)
		VALUES ($1, $2, $3, $4)
	`, input.KodeICD9, input.NamaTindakan, input.Deskripsi, input.BiayaDasar)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan tindakan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tindakan berhasil ditambahkan"})
}

func GetAllTindakan(c *gin.Context) {
	rows, err := db.Query(`
		SELECT tindakan_id, kode_icd9, nama_tindakan, deskripsi, biaya_dasar
		FROM tindakan ORDER BY tindakan_id`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data tindakan"})
		return
	}
	defer rows.Close()

	var result []dto.TindakanResponse
	for rows.Next() {
		var t dto.TindakanResponse
		err := rows.Scan(&t.TindakanID, &t.KodeICD9, &t.NamaTindakan, &t.Deskripsi, &t.BiayaDasar)
		if err == nil {
			result = append(result, t)
		}
	}

	c.JSON(http.StatusOK, result)
}

func UpdateTindakan(c *gin.Context) {
	id := c.Param("id")
	var input dto.CreateTindakanRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := db.Exec(`
		UPDATE tindakan SET 
			kode_icd9 = $1,
			nama_tindakan = $2,
			deskripsi = $3,
			biaya_dasar = $4
		WHERE tindakan_id = $5
	`, input.KodeICD9, input.NamaTindakan, input.Deskripsi, input.BiayaDasar, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update tindakan"})
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tindakan tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tindakan diperbarui"})
}

func DeleteTindakan(c *gin.Context) {
	id := c.Param("id")
	res, err := db.Exec(`DELETE FROM tindakan WHERE tindakan_id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus tindakan"})
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tindakan tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tindakan berhasil dihapus"})
}
