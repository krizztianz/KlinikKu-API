package models

type Kunjungan struct {
	KunjunganID      int      `json:"kunjungan_id"`
	PasienID         int      `json:"pasien_id"`
	DokterID         int      `json:"dokter_id"`
	TanggalKunjungan string   `json:"tanggal_kunjungan"`
	Keluhan          string   `json:"keluhan"`
	TinggiBadan      *float64 `json:"tinggi_badan"`
	BeratBadan       *float64 `json:"berat_badan"`
	TekananDarah     *string  `json:"tekanan_darah"`
	SuhuTubuh        *float64 `json:"suhu_tubuh"`
	JenisKunjungan   string   `json:"jenis_kunjungan"`
	Status           string   `json:"status"`
	Prioritas        *string  `json:"prioritas"`
}
