package dto

type CreateObatRequest struct {
	NamaObat  string  `json:"nama_obat" binding:"required"`
	Deskripsi string  `json:"deskripsi"`
	Stok      int     `json:"stok" binding:"min=0"`
	Harga     float64 `json:"harga" binding:"min=0"`
	Satuan    string  `json:"satuan"`
}

type UpdateObatRequest struct {
	NamaObat  string  `json:"nama_obat"`
	Deskripsi string  `json:"deskripsi"`
	Stok      int     `json:"stok" binding:"min=0"`
	Harga     float64 `json:"harga" binding:"min=0"`
	Satuan    string  `json:"satuan"`
}
