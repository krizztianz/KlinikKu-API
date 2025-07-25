package dto

type CreateSpesialisasiRequest struct {
	KodeSpesialisasi string `json:"kode_spesialisasi" binding:"required"`
	NamaSpesialisasi string `json:"nama_spesialisasi" binding:"required"`
	Deskripsi        string `json:"deskripsi"`
}

type SpesialisasiResponse struct {
	SpesialisasiID   int    `json:"spesialisasi_id"`
	KodeSpesialisasi string `json:"kode_spesialisasi"`
	NamaSpesialisasi string `json:"nama_spesialisasi"`
	Deskripsi        string `json:"deskripsi"`
}
