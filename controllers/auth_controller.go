package controllers

import (
	"KlinikKu/models"
	"KlinikKu/utils"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB(database *sql.DB) {
	db = database
}

// @Summary LoginHandler
// @Tags Auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /login [post]

func LoginHandler(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	err := db.QueryRow(`SELECT user_id, username, password, role FROM users WHERE username=$1`, req.Username).
		Scan(&user.UserID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
		return
	}

	// Generate tokens
	accessToken, _ := utils.GenerateJWT(user.UserID, user.Username, user.Role)
	refreshToken, _ := utils.GenerateRefreshToken()

	// Simpan refresh token ke DB
	_, err = db.Exec(`UPDATE users SET refresh_token=$1 WHERE user_id=$2`, refreshToken, user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// @Summary RefreshHandler
// @Tags Auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /refresh [get]

func RefreshHandler(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid refresh token"})
		return
	}

	var user models.User
	err := db.QueryRow(`SELECT user_id, username, role FROM users WHERE refresh_token=$1`, req.RefreshToken).
		Scan(&user.UserID, &user.Username, &user.Role)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	newAccessToken, _ := utils.GenerateJWT(user.UserID, user.Username, user.Role)
	newRefreshToken, _ := utils.GenerateRefreshToken()

	_, err = db.Exec(`UPDATE users SET refresh_token=$1 WHERE user_id=$2`, newRefreshToken, user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

// @Summary GeneratePassword
// @Tags Auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /generate-password [get]

func GeneratePassword(c *gin.Context) {
	var req struct {
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Password"})
		return
	}

	pass, err := utils.HashPassword(req.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Geneate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Plain Text": req.Password,
		"Password":   pass,
	})
}
