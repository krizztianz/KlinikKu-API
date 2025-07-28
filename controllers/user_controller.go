package controllers

import (
	"KlinikKu/dto"
	"KlinikKu/middleware"
	"KlinikKu/utils"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateUser - Tambah user baru
// @Summary CreateUser
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /create-user [post]
func CreateUser(c *gin.Context) {
	var input dto.CreateUserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdBy := middleware.GetUsername(c)

	err := utils.RunTx(db, func(tx *sql.Tx) error {
		_, err := tx.Exec(`
			INSERT INTO users (username, nama_lengkap, password, role, dokter_id, created_at, created_by)
			VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, $6)
		`, input.Username, input.NamaLengkap, input.Password, input.Role, input.DokterID, createdBy)
		return err
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User berhasil dibuat"})
}

// GetAllUsers - Ambil semua user
// @Summary GetAllUsers
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /get-all-users [get]
func GetAllUsers(c *gin.Context) {
	rows, err := db.Query(`
		SELECT user_id, username, nama_lengkap, role, dokter_id
		FROM users ORDER BY user_id
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data user"})
		return
	}
	defer rows.Close()

	var result []dto.UserResponse
	for rows.Next() {
		var u dto.UserResponse
		err := rows.Scan(&u.UserID, &u.Username, &u.NamaLengkap, &u.Role, &u.DokterID)
		if err == nil {
			result = append(result, u)
		}
	}

	c.JSON(http.StatusOK, result)
}

// GetUserByID - Ambil user berdasarkan ID
// @Summary GetUserByID
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /get-user-by-i-d [get]
func GetUserByID(c *gin.Context) {
	id := c.Param("id")

	var u dto.UserResponse
	err := db.QueryRow(`
		SELECT user_id, username, nama_lengkap, role, dokter_id
		FROM users WHERE user_id = $1
	`, id).Scan(&u.UserID, &u.Username, &u.NamaLengkap, &u.Role, &u.DokterID)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data user"})
		return
	}

	c.JSON(http.StatusOK, u)
}

// UpdateUser - Ubah data user
// @Summary UpdateUser
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /update-user [put]
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var input dto.CreateUserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	modifiedBy := middleware.GetUsername(c)

	err := utils.RunTx(db, func(tx *sql.Tx) error {
		res, err := tx.Exec(`
			UPDATE users SET
				username = $1,
				nama_lengkap = $2,
				password = $3,
				role = $4,
				dokter_id = $5,
				modified_by = $6,
				modified_at = CURRENT_TIMESTAMP
			WHERE user_id = $7
		`, input.Username, input.NamaLengkap, input.Password, input.Role, input.DokterID, modifiedBy, id)
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
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update user"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "User berhasil diperbarui"})
	}
}

// DeleteUser - Hapus user berdasarkan ID
// @Summary DeleteUser
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /delete-user [delete]

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	res, err := db.Exec(`DELETE FROM users WHERE user_id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus user"})
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User berhasil dihapus"})
}
