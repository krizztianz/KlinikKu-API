package controllers

import (
	"KlinikKu/dto"
	"KlinikKu/middleware"
	"KlinikKu/utils"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary CreateDokter
// @Tags Dokter
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /create-dokter [post]

func CreateDokter(c *gin.Context) {
	var input dto.CreateDokterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdBy := middleware.GetUsername(c)

	err := utils.RunTx(db, func(tx *sql.Tx) error {
		_, err := tx.Exec(`
			INSERT INTO dokter (
				no_izin_praktek, nama, tanggal_lahir, jenis_kelamin, alamat,
				no_hp, no_telepon, ktp, email, created_by
			) VALUES ($1,$2,$3,$4::gender,$5,$6,$7,$8,$9,$10)
		`, input.NoIzinPraktek, input.Nama, input.TanggalLahir, input.JenisKelamin,
			input.Alamat, input.NoHP, input.NoTelepon, input.KTP, input.Email, createdBy)
		return err
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan dokter"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dokter berhasil ditambahkan"})
}

// @Summary GetAllDokter
// @Tags Dokter
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /get-all-dokter [get]

func GetAllDokter(c *gin.Context) {
	rows, err := db.Query(`SELECT dokter_id, nama, no_hp FROM dokter`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data dokter"})
		return
	}
	defer rows.Close()

	var allDokter []dto.DokterWithSpesialisasi

	for rows.Next() {
		var d dto.DokterWithSpesialisasi
		if err := rows.Scan(&d.DokterID, &d.Nama, &d.NoHP); err != nil {
			continue
		}

		// Ambil spesialisasi untuk tiap dokter
		spesialisasiRows, err := db.Query(`
			SELECT s.spesialisasi_id, s.kode_spesialisasi, s.nama_spesialisasi
			FROM dokter_spesialisasi ds
			JOIN spesialisasi s ON s.spesialisasi_id = ds.spesialisasi_id
			WHERE ds.dokter_id = $1
		`, d.DokterID)
		if err == nil {
			for spesialisasiRows.Next() {
				var s dto.SpesialisasiSimpleDTO
				spesialisasiRows.Scan(&s.SpesialisasiID, &s.KodeSpesialisasi, &s.NamaSpesialisasi)
				d.Spesialisasi = append(d.Spesialisasi, s)
			}
			spesialisasiRows.Close()
		}

		allDokter = append(allDokter, d)
	}

	c.JSON(http.StatusOK, allDokter)
}

// @Summary GetDokterByID
// @Tags Dokter
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /get-dokter-by-i-d [get]

func GetDokterByID(c *gin.Context) {
	id := c.Param("id")

	var dokter dto.DokterWithSpesialisasi
	err := db.QueryRow(`SELECT dokter_id, nama, no_hp FROM dokter WHERE dokter_id = $1`, id).
		Scan(&dokter.DokterID, &dokter.Nama, &dokter.NoHP)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dokter tidak ditemukan"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil dokter"})
		return
	}

	// Ambil spesialisasi dokter
	rows, err := db.Query(`
		SELECT s.spesialisasi_id, s.kode_spesialisasi, s.nama_spesialisasi
		FROM dokter_spesialisasi ds
		JOIN spesialisasi s ON s.spesialisasi_id = ds.spesialisasi_id
		WHERE ds.dokter_id = $1
	`, id)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var s dto.SpesialisasiSimpleDTO
			rows.Scan(&s.SpesialisasiID, &s.KodeSpesialisasi, &s.NamaSpesialisasi)
			dokter.Spesialisasi = append(dokter.Spesialisasi, s)
		}
	}

	c.JSON(http.StatusOK, dokter)
}

// @Summary UpdateDokter
// @Tags Dokter
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /update-dokter [put]

func UpdateDokter(c *gin.Context) {
	dokterID := c.Param("id")
	var input dto.CreateDokterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	modifiedBy := middleware.GetUsername(c)

	err := utils.RunTx(db, func(tx *sql.Tx) error {
		res, err := tx.Exec(`
			UPDATE dokter SET
				no_izin_praktek=$1,
				nama=$2,
				tanggal_lahir=$3,
				jenis_kelamin=$4::gender,
				alamat=$5,
				no_hp=$6,
				no_telepon=$7,
				ktp=$8,
				email=$9,
				modified_by=$10,
				modified_at=CURRENT_TIMESTAMP
			WHERE dokter_id=$11
		`, input.NoIzinPraktek, input.Nama, input.TanggalLahir, input.JenisKelamin,
			input.Alamat, input.NoHP, input.NoTelepon, input.KTP, input.Email, modifiedBy, dokterID)

		if err != nil {
			return err
		}

		rows, _ := res.RowsAffected()
		if rows == 0 {
			return sql.ErrNoRows
		}
		return nil
	})

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dokter tidak ditemukan"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update data dokter"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data dokter diperbarui"})
}

// @Summary DeleteDokter
// @Tags Dokter
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /delete-dokter [delete]

func DeleteDokter(c *gin.Context) {
	dokterID := c.Param("id")

	_, err := db.Exec(`DELETE FROM dokter WHERE dokter_id = $1`, dokterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hapus dokter"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dokter berhasil dihapus"})
}

// @Summary AssignSpesialisasiToDokter
// @Tags Dokter
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /assign-spesialisasi-to-dokter [get]

func AssignSpesialisasiToDokter(c *gin.Context) {
	dokterID := c.Param("id")
	var body struct {
		SpesialisasiID int `json:"spesialisasi_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	createdBy := middleware.GetUsername(c)

	err := utils.RunTx(db, func(tx *sql.Tx) error {
		_, err := tx.Exec(`
			INSERT INTO dokter_spesialisasi (dokter_id, spesialisasi_id, created_by)
			VALUES ($1, $2, $3) ON CONFLICT DO NOTHING
		`, dokterID, body.SpesialisasiID, createdBy)
		return err
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan spesialisasi ke dokter"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Spesialisasi berhasil ditambahkan ke dokter"})
}

// @Summary RemoveSpesialisasiFromDokter
// @Tags Dokter
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /remove-spesialisasi-from-dokter [get]

func RemoveSpesialisasiFromDokter(c *gin.Context) {
	dokterID := c.Param("id")
	spesialisasiID := c.Param("spesialisasi_id")

	err := utils.RunTx(db, func(tx *sql.Tx) error {
		_, err := tx.Exec(`
			DELETE FROM dokter_spesialisasi
			WHERE dokter_id = $1 AND spesialisasi_id = $2
		`, dokterID, spesialisasiID)
		return err
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus spesialisasi dari dokter"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Spesialisasi berhasil dihapus dari dokter"})
}
