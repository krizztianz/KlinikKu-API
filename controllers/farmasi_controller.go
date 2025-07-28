package controllers

import (
	"KlinikKu/middleware"
	"KlinikKu/utils"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPendingResep(c *gin.Context) {
	rows, err := db.Query(`
		SELECT r.resep_id, p.nama AS nama_pasien, d.nama AS nama_dokter, r.created_at
		FROM resep r
		JOIN rekam_medis rm ON rm.rekam_medis_id = r.rekam_medis_id
		JOIN kunjungan k ON k.kunjungan_id = rm.kunjungan_id
		JOIN pasien p ON p.pasien_id = k.pasien_id
		JOIN dokter d ON d.dokter_id = k.dokter_id
		WHERE r.status != 'selesai'
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar resep"})
		return
	}
	defer rows.Close()

	var resepList []map[string]interface{}
	for rows.Next() {
		var resepID int
		var namaPasien, namaDokter string
		var createdAt string
		rows.Scan(&resepID, &namaPasien, &namaDokter, &createdAt)

		resepList = append(resepList, gin.H{
			"resep_id":    resepID,
			"nama_pasien": namaPasien,
			"nama_dokter": namaDokter,
			"created_at":  createdAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"resep": resepList})
}

func GetResepDetail(c *gin.Context) {
	id := c.Param("id")

	var exists bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM resep WHERE resep_id = $1)`, id).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan"})
		return
	}

	// Reuse GetResepByID
	GetResepByID(c)
}

func MarkResepAsCompleted(c *gin.Context) {
	id := c.Param("id")
	resepID, _ := strconv.Atoi(id)
	modifiedBy := middleware.GetUsername(c)

	err := utils.RunTx(db, func(tx *sql.Tx) error {
		// Ambil data resep_detail
		rows, err := tx.Query(`
			SELECT obat_id, jumlah FROM resep_detail WHERE resep_id = $1
		`, resepID)
		if err != nil {
			return err
		}
		defer rows.Close()

		// Kurangi stok obat
		for rows.Next() {
			var obatID, jumlah int
			if err := rows.Scan(&obatID, &jumlah); err != nil {
				return err
			}

			res, err := tx.Exec(`
				UPDATE obat SET stok = stok - $1, modified_by = $3, modified_at = CURRENT_TIMESTAMP
				WHERE obat_id = $2 AND stok >= $1
			`, jumlah, obatID, modifiedBy)
			if err != nil {
				return err
			}
			affected, _ := res.RowsAffected()
			if affected == 0 {
				return sql.ErrNoRows // stok tidak cukup
			}
		}

		// Update status resep
		_, err = tx.Exec(`
			UPDATE resep SET status = 'selesai', modified_by = $2, modified_at = CURRENT_TIMESTAMP
			WHERE resep_id = $1
		`, resepID, modifiedBy)
		if err != nil {
			return err
		}

		// Update status kunjungan
		_, err = tx.Exec(`
			UPDATE kunjungan SET status = 'selesai', modified_by = $2, modified_at = CURRENT_TIMESTAMP
			WHERE kunjungan_id = (
				SELECT rm.kunjungan_id FROM rekam_medis rm
				JOIN resep r ON r.rekam_medis_id = rm.rekam_medis_id
				WHERE r.resep_id = $1
			)
		`, resepID, modifiedBy)
		if err != nil {
			return err
		}

		return nil
	})

	if err == sql.ErrNoRows {
		c.JSON(http.StatusConflict, gin.H{"error": "Stok obat tidak cukup"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan saat menebus resep"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Resep berhasil ditebus & stok obat dikurangi"})
}
