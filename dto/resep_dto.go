package dto

type ResepObatDetail struct {
	ObatID     int    `json:"obat_id" binding:"required"`
	Jumlah     int    `json:"jumlah" binding:"required,min=1"`
	Dosis      string `json:"dosis" binding:"required"` // contoh: "3x1"
	Keterangan string `json:"keterangan"`
}

type CreateResepRequest struct {
	RekamMedisID int               `json:"rekam_medis_id" binding:"required"`
	Obats        []ResepObatDetail `json:"obats" binding:"required"`
}
