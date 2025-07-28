package controllers

import (
	"KlinikKu/dto"
	"KlinikKu/models"
	"KlinikKu/utils"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateObat(c *gin.Context) {
	var input dto.CreateObatRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := c.GetString("username")

	var id int
	err := utils.RunTx(db, func(tx *sql.Tx) error {
		return tx.QueryRow(`
			INSERT INTO obat (nama_obat, deskripsi, stok, harga, satuan, created_by)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING obat_id
		`, input.NamaObat, input.Deskripsi, input.Stok, input.Harga, input.Satuan, username).Scan(&id)
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat obat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Obat berhasil dibuat", "obat_id": id})
}

func GetAllObat(c *gin.Context) {
	rows, err := db.Query(`
		SELECT obat_id, nama_obat, deskripsi, stok, harga, satuan
		FROM obat ORDER BY nama_obat
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data obat"})
		return
	}
	defer rows.Close()

	var obats []models.Obat
	for rows.Next() {
		var o models.Obat
		err := rows.Scan(&o.ObatID, &o.NamaObat, &o.Deskripsi, &o.Stok, &o.Harga, &o.Satuan)
		if err == nil {
			obats = append(obats, o)
		}
	}

	c.JSON(http.StatusOK, obats)
}

func GetObatByID(c *gin.Context) {
	id := c.Param("id")
	var o models.Obat
	err := db.QueryRow(`
		SELECT obat_id, nama_obat, deskripsi, stok, harga, satuan
		FROM obat WHERE obat_id = $1
	`, id).Scan(&o.ObatID, &o.NamaObat, &o.Deskripsi, &o.Stok, &o.Harga, &o.Satuan)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Obat tidak ditemukan"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, o)
}

func UpdateObat(c *gin.Context) {
	obatID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID obat tidak valid"})
		return
	}

	var input dto.UpdateObatRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := c.GetString("username")

	err = utils.RunTx(db, func(tx *sql.Tx) error {
		_, err := tx.Exec(`
			UPDATE obat
			SET nama_obat = $1, deskripsi = $2, stok = $3, harga = $4, satuan = $5,
				modified_at = CURRENT_TIMESTAMP, modified_by = $6
			WHERE obat_id = $7
		`, input.NamaObat, input.Deskripsi, input.Stok, input.Harga, input.Satuan, username, obatID)
		return err
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui obat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Obat berhasil diperbarui"})
}

func DeleteObat(c *gin.Context) {
	id := c.Param("id")

	err := utils.RunTx(db, func(tx *sql.Tx) error {
		_, err := tx.Exec(`DELETE FROM obat WHERE obat_id = $1`, id)
		return err
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus obat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Obat berhasil dihapus"})
}
