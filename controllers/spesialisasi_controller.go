package controllers

import (
	"KlinikKu/dto"
	"KlinikKu/middleware"
	"KlinikKu/utils"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary CreateSpesialisasi
// @Tags Spesialisasi
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /create-spesialisasi [post]

func CreateSpesialisasi(c *gin.Context) {
	var input dto.CreateSpesialisasiRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdBy := middleware.GetUsername(c)

	err := utils.RunTx(db, func(tx *sql.Tx) error {
		_, err := tx.Exec(`
			INSERT INTO spesialisasi (kode_spesialisasi, nama_spesialisasi, deskripsi, created_by)
			VALUES ($1, $2, $3, $4)
		`, input.KodeSpesialisasi, input.NamaSpesialisasi, input.Deskripsi, createdBy)
		return err
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan spesialisasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Spesialisasi berhasil ditambahkan"})
}

// @Summary GetAllSpesialisasi
// @Tags Spesialisasi
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /get-all-spesialisasi [get]

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

// @Summary UpdateSpesialisasi
// @Tags Spesialisasi
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /update-spesialisasi [put]

func UpdateSpesialisasi(c *gin.Context) {
	id := c.Param("id")
	var input dto.CreateSpesialisasiRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	modifiedBy := middleware.GetUsername(c)

	err := utils.RunTx(db, func(tx *sql.Tx) error {
		res, err := tx.Exec(`
			UPDATE spesialisasi SET 
				kode_spesialisasi = $1,
				nama_spesialisasi = $2,
				deskripsi = $3,
				modified_by = $4,
				modified_at = CURRENT_TIMESTAMP
			WHERE spesialisasi_id = $5
		`, input.KodeSpesialisasi, input.NamaSpesialisasi, input.Deskripsi, modifiedBy, id)
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Spesialisasi tidak ditemukan"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update spesialisasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Spesialisasi diperbarui"})
}

// @Summary DeleteSpesialisasi
// @Tags Spesialisasi
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /delete-spesialisasi [delete]

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
