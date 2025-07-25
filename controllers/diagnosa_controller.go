package controllers

import (
	"KlinikKu/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateDiagnosa(c *gin.Context) {
	var input dto.CreateDiagnosaRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec(`
		INSERT INTO diagnosa (kode_icd10, nama_diagnosa, deskripsi)
		VALUES ($1, $2, $3)
	`, input.KodeICD10, input.NamaDiagnosa, input.Deskripsi)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan diagnosa"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Diagnosa berhasil ditambahkan"})
}

func GetAllDiagnosa(c *gin.Context) {
	rows, err := db.Query(`SELECT diagnosa_id, kode_icd10, nama_diagnosa, deskripsi FROM diagnosa ORDER BY diagnosa_id`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data diagnosa"})
		return
	}
	defer rows.Close()

	var result []dto.DiagnosaResponse
	for rows.Next() {
		var d dto.DiagnosaResponse
		err := rows.Scan(&d.DiagnosaID, &d.KodeICD10, &d.NamaDiagnosa, &d.Deskripsi)
		if err == nil {
			result = append(result, d)
		}
	}

	c.JSON(http.StatusOK, result)
}

func UpdateDiagnosa(c *gin.Context) {
	id := c.Param("id")
	var input dto.CreateDiagnosaRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := db.Exec(`
		UPDATE diagnosa SET 
			kode_icd10 = $1,
			nama_diagnosa = $2,
			deskripsi = $3
		WHERE diagnosa_id = $4
	`, input.KodeICD10, input.NamaDiagnosa, input.Deskripsi, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update diagnosa"})
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diagnosa tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Diagnosa diperbarui"})
}

func DeleteDiagnosa(c *gin.Context) {
	id := c.Param("id")
	res, err := db.Exec(`DELETE FROM diagnosa WHERE diagnosa_id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus diagnosa"})
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diagnosa tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Diagnosa berhasil dihapus"})
}
