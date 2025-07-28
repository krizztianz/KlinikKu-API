package controllers

import (
	"KlinikKu/dto"
	"KlinikKu/middleware"
	"KlinikKu/utils"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary CreateDiagnosa
// @Tags Diagnosa
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /create-diagnosa [post]

func CreateDiagnosa(c *gin.Context) {
	var input dto.CreateDiagnosaRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdBy := middleware.GetUsername(c)

	err := utils.RunTx(db, func(tx *sql.Tx) error {
		_, err := tx.Exec(`
			INSERT INTO diagnosa (kode_icd10, nama_diagnosa, deskripsi, created_by)
			VALUES ($1, $2, $3, $4)
		`, input.KodeICD10, input.NamaDiagnosa, input.Deskripsi, createdBy)
		return err
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan diagnosa"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Diagnosa berhasil ditambahkan"})
}

// @Summary GetAllDiagnosa
// @Tags Diagnosa
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /get-all-diagnosa [get]

func GetAllDiagnosa(c *gin.Context) {
	rows, err := db.Query(`
		SELECT diagnosa_id, kode_icd10, nama_diagnosa, deskripsi
		FROM diagnosa ORDER BY diagnosa_id`)
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

// @Summary UpdateDiagnosa
// @Tags Diagnosa
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /update-diagnosa [put]

func UpdateDiagnosa(c *gin.Context) {
	id := c.Param("id")
	var input dto.CreateDiagnosaRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	modifiedBy := middleware.GetUsername(c)

	err := utils.RunTx(db, func(tx *sql.Tx) error {
		res, err := tx.Exec(`
			UPDATE diagnosa SET 
				kode_icd10 = $1,
				nama_diagnosa = $2,
				deskripsi = $3,
				modified_by = $4,
				modified_at = CURRENT_TIMESTAMP
			WHERE diagnosa_id = $5
		`, input.KodeICD10, input.NamaDiagnosa, input.Deskripsi, modifiedBy, id)
		if err != nil {
			return err
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return sql.ErrNoRows
		}
		return nil
	})

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diagnosa tidak ditemukan"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update diagnosa"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Diagnosa diperbarui"})
}

// @Summary DeleteDiagnosa
// @Tags Diagnosa
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /delete-diagnosa [delete]

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
