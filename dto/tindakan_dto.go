package dto

type CreateTindakanRequest struct {
	KodeICD9     string  `json:"kode_icd9" binding:"required"`
	NamaTindakan string  `json:"nama_tindakan" binding:"required"`
	Deskripsi    string  `json:"deskripsi"`
	BiayaDasar   float64 `json:"biaya_dasar" binding:"required"`
}

type TindakanResponse struct {
	TindakanID   int     `json:"tindakan_id"`
	KodeICD9     string  `json:"kode_icd9"`
	NamaTindakan string  `json:"nama_tindakan"`
	Deskripsi    string  `json:"deskripsi"`
	BiayaDasar   float64 `json:"biaya_dasar"`
}
