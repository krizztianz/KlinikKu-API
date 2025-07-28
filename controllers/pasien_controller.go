package controllers

import (
	"KlinikKu/dto"
	"KlinikKu/utils"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary GetAllPasien
// @Tags Pasien
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /get-all-pasien [get]

func GetAllPasien(c *gin.Context) {
	rows, err := db.Query(`
		SELECT pasien_id, nama, tanggal_lahir, jenis_kelamin, alamat, 
			   no_hp, no_telepon, ktp, email, golongan_darah 
		FROM pasien 
		ORDER BY pasien_id
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pasien"})
		return
	}
	defer rows.Close()

	var result []dto.PasienResponse
	for rows.Next() {
		var p dto.PasienResponse
		err := rows.Scan(
			&p.PasienID, &p.Nama, &p.TanggalLahir, &p.JenisKelamin, &p.Alamat,
			&p.NoHP, &p.NoTelepon, &p.KTP, &p.Email, &p.GolonganDarah,
		)
		if err == nil {
			result = append(result, p)
		}
	}

	c.JSON(http.StatusOK, result)
}

// @Summary SearchPasien
// @Tags Pasien
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /search-pasien [get]

func SearchPasien(c *gin.Context) {
	ktp := c.Query("ktp")
	noHP := c.Query("no_hp")

	if ktp == "" && noHP == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "KTP atau No HP harus diisi untuk pencarian"})
		return
	}

	var (
		query string
		rows  *sql.Rows
		err   error
	)

	if ktp != "" {
		query = `
			SELECT pasien_id, nama, tanggal_lahir, jenis_kelamin, alamat, 
				   no_hp, no_telepon, ktp, email, golongan_darah 
			FROM pasien 
			WHERE ktp = $1
		`
		rows, err = db.Query(query, ktp)
	} else {
		query = `
			SELECT pasien_id, nama, tanggal_lahir, jenis_kelamin, alamat, 
				   no_hp, no_telepon, ktp, email, golongan_darah 
			FROM pasien 
			WHERE no_hp = $1
		`
		rows, err = db.Query(query, noHP)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencari pasien"})
		return
	}
	defer rows.Close()

	var results []dto.PasienResponse
	for rows.Next() {
		var p dto.PasienResponse
		err := rows.Scan(
			&p.PasienID, &p.Nama, &p.TanggalLahir, &p.JenisKelamin, &p.Alamat,
			&p.NoHP, &p.NoTelepon, &p.KTP, &p.Email, &p.GolonganDarah,
		)
		if err == nil {
			results = append(results, p)
		}
	}

	if len(results) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Pasien tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// @Summary CreatePasien
// @Tags Pasien
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /create-pasien [post]

func CreatePasien(c *gin.Context) {
	var input dto.CreatePasienRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := c.GetString("username")

	err := utils.RunTx(db, func(tx *sql.Tx) error {
		_, err := tx.Exec(`
			INSERT INTO pasien (
				nama, tanggal_lahir, jenis_kelamin, alamat, no_hp, no_telepon,
				ktp, email, golongan_darah, created_by
			)
			VALUES ($1,$2,$3::gender,$4,$5,$6,$7,$8,$9::blood_type,$10)
		`,
			input.Nama, input.TanggalLahir, input.JenisKelamin, input.Alamat,
			input.NoHP, input.NoTelepon, input.KTP, input.Email, input.GolonganDarah, username,
		)
		return err
	})

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan pasien"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Pasien berhasil ditambahkan"})
}

// @Summary UpdatePasien
// @Tags Pasien
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /update-pasien [put]

func UpdatePasien(c *gin.Context) {
	pasienID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var input dto.UpdatePasienRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := c.GetString("username")

	err = utils.RunTx(db, func(tx *sql.Tx) error {
		_, err := tx.Exec(`
			UPDATE pasien SET
				nama = $1,
				tanggal_lahir = $2,
				jenis_kelamin = $3::gender,
				alamat = $4,
				no_hp = $5,
				no_telepon = $6,
				ktp = $7,
				email = $8,
				golongan_darah = $9::blood_type,
				modified_at = CURRENT_TIMESTAMP,
				modified_by = $10
			WHERE pasien_id = $11
		`,
			input.Nama, input.TanggalLahir, input.JenisKelamin, input.Alamat,
			input.NoHP, input.NoTelepon, input.KTP, input.Email, input.GolonganDarah, username, pasienID,
		)
		return err
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui pasien"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pasien berhasil diperbarui"})
}

// @Summary DeletePasien
// @Tags Pasien
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /delete-pasien [delete]

func DeletePasien(c *gin.Context) {
	pasienID := c.Param("id")

	err := utils.RunTx(db, func(tx *sql.Tx) error {
		res, err := tx.Exec(`DELETE FROM pasien WHERE pasien_id = $1`, pasienID)
		if err != nil {
			return err
		}
		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			return sql.ErrNoRows
		}
		return nil
	})

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pasien tidak ditemukan"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hapus pasien"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pasien berhasil dihapus"})
}
