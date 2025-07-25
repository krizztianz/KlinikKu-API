package dto

type CreateRekamMedisRequest struct {
	KunjunganID      int    `json:"kunjungan_id" binding:"required"`
	Anamnesa         string `json:"anamnesa" binding:"required"`
	PemeriksaanFisik string `json:"pemeriksaan_fisik" binding:"required"`
	Catatan          string `json:"catatan,omitempty"`
}

type AddDiagnosaRequest struct {
	DiagnosaIDs []int `json:"diagnosa_ids" binding:"required"` // bisa satu atau banyak
}

type AddTindakan struct {
	TindakanID  int     `json:"tindakan_id" binding:"required"`
	Jumlah      int     `json:"jumlah" binding:"required,min=1"`
	BiayaAktual float64 `json:"biaya_aktual" binding:"required"`
	Catatan     string  `json:"catatan,omitempty"`
}

type AddTindakanRequest struct {
	Tindakans []AddTindakan `json:"tindakans" binding:"required"`
}
