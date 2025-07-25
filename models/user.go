package models

type User struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Role     string `json:"role"`
	DokterID *int   `json:"dokter_id"`
}
