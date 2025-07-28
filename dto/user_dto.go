package dto

type CreateUserRequest struct {
	Username    string `json:"username" binding:"required"`
	NamaLengkap string `json:"nama_lengkap" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Role        string `json:"role" binding:"required,oneof=admin frontliner dokter farmasi"`
	DokterID    *int   `json:"dokter_id,omitempty` // Optional
}

type UserResponse struct {
	UserID      int    `json:"user_id"`
	Username    string `json:"username"`
	NamaLengkap string `json:"nama_lengkap"`
	Role        string `json:"role"`
	DokterID    *int   `json:"dokter_id"`
}
