package controllers

import (
	"KlinikKu/dto"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllPasien(c *gin.Context) {
	rows, err := db.Query(`SELECT pasien_id, nama, tanggal_lahir, jenis_kelamin, alamat, no_hp, no_telepon, ktp, email, golongan_darah FROM pasien ORDER BY pasien_id`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pasien"})
		return
	}
	defer rows.Close()

	var result []dto.PasienResponse
	for rows.Next() {
		var p dto.PasienResponse
		err := rows.Scan(&p.PasienID, &p.Nama, &p.TanggalLahir, &p.JenisKelamin, &p.Alamat,
			&p.NoHP, &p.NoTelepon, &p.KTP, &p.Email, &p.GolonganDarah)
		if err == nil {
			result = append(result, p)
		}
	}

	c.JSON(http.StatusOK, result)
}

func SearchPasien(c *gin.Context) {
	ktp := c.Query("ktp")
	noHP := c.Query("no_hp")

	if ktp == "" && noHP == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "KTP atau No HP harus diisi untuk pencarian"})
		return
	}

	query := `SELECT pasien_id, nama, tanggal_lahir, jenis_kelamin, alamat, no_hp, no_telepon, ktp, email, golongan_darah 
			  FROM pasien WHERE `
	var rows *sql.Rows
	var err error

	if ktp != "" {
		query += `ktp = $1`
		rows, err = db.Query(query, ktp)
	} else {
		query += `no_hp = $1`
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
		err := rows.Scan(&p.PasienID, &p.Nama, &p.TanggalLahir, &p.JenisKelamin, &p.Alamat,
			&p.NoHP, &p.NoTelepon, &p.KTP, &p.Email, &p.GolonganDarah)
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

func CreatePasien(c *gin.Context) {
	var input dto.CreatePasienRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec(`
		INSERT INTO pasien 
		(nama, tanggal_lahir, jenis_kelamin, alamat, no_hp, no_telepon, ktp, email, golongan_darah)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`, input.Nama, input.TanggalLahir, input.JenisKelamin, input.Alamat,
		input.NoHP, input.NoTelepon, input.KTP, input.Email, input.GolonganDarah)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat pasien"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pasien berhasil ditambahkan"})
}

func UpdatePasien(c *gin.Context) {
	pasienID := c.Param("id")
	var input dto.CreatePasienRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := db.Exec(`
		UPDATE pasien SET
			nama = $1,
			tanggal_lahir = $2,
			jenis_kelamin = $3,
			alamat = $4,
			no_hp = $5,
			no_telepon = $6,
			ktp = $7,
			email = $8,
			golongan_darah = $9
		WHERE pasien_id = $10
	`, input.Nama, input.TanggalLahir, input.JenisKelamin, input.Alamat, input.NoHP,
		input.NoTelepon, input.KTP, input.Email, input.GolonganDarah, pasienID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update data pasien"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pasien tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data pasien diperbarui"})
}

func DeletePasien(c *gin.Context) {
	pasienID := c.Param("id")
	res, err := db.Exec(`DELETE FROM pasien WHERE pasien_id = $1`, pasienID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hapus pasien"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pasien tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pasien berhasil dihapus"})
}
