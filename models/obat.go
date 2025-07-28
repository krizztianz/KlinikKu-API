package models

type Obat struct {
	ObatID     int     `json:"obat_id"`
	NamaObat   string  `json:"nama_obat"`
	Deskripsi  string  `json:"deskripsi"`
	Stok       int     `json:"stok"`
	Harga      float64 `json:"harga"`
	Satuan     string  `json:"satuan"`
	CreatedAt  string  `json:"created_at"`
	CreatedBy  string  `json:"created_by"`
	ModifiedAt string  `json:"modified_at"`
	ModifiedBy string  `json:"modified_by"`
}
