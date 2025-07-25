package dto

type CreateKunjunganRequest struct {
	PasienID       int     `json:"pasien_id" binding:"required"`
	DokterID       int     `json:"dokter_id" binding:"required"`
	Keluhan        string  `json:"keluhan" binding:"required"`
	JenisKunjungan string  `json:"jenis_kunjungan" binding:"required"`
	Prioritas      *string `json:"prioritas,omitempty"`
}

type KunjunganResponse struct {
	KunjunganID      int     `json:"kunjungan_id"`
	PasienID         int     `json:"pasien_id"`
	DokterID         int     `json:"dokter_id"`
	TanggalKunjungan string  `json:"tanggal_kunjungan"`
	Keluhan          string  `json:"keluhan"`
	JenisKunjungan   string  `json:"jenis_kunjungan"`
	Status           string  `json:"status"`
	Prioritas        *string `json:"prioritas,omitempty"`
}

type UpdateStatusKunjunganRequest struct {
	Status string `json:"status" binding:"required"`
}
