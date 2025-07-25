package dto

type CreateDokterRequest struct {
	NoIzinPraktek *string `json:"no_izin_praktek,omitempty"`
	Nama          string  `json:"nama" binding:"required"`
	TanggalLahir  string  `json:"tanggal_lahir" binding:"required"`
	JenisKelamin  *string `json:"jenis_kelamin,omitempty"`
	Alamat        string  `json:"alamat" binding:"required"`
	NoHP          string  `json:"no_hp" binding:"required"`
	NoTelepon     *string `json:"no_telepon,omitempty"`
	KTP           string  `json:"ktp" binding:"required"`
	Email         *string `json:"email,omitempty"`
}

type DokterResponse struct {
	DokterID      int     `json:"dokter_id"`
	NoIzinPraktek *string `json:"no_izin_praktek"`
	Nama          string  `json:"nama"`
	TanggalLahir  string  `json:"tanggal_lahir"`
	JenisKelamin  *string `json:"jenis_kelamin"`
	Alamat        string  `json:"alamat"`
	NoHP          string  `json:"no_hp"`
	NoTelepon     *string `json:"no_telepon"`
	KTP           string  `json:"ktp"`
	Email         *string `json:"email"`
}

type DokterWithSpesialisasi struct {
	DokterID     int                     `json:"dokter_id"`
	Nama         string                  `json:"nama"`
	NoHP         string                  `json:"no_hp"`
	Spesialisasi []SpesialisasiSimpleDTO `json:"spesialisasi"`
}

type SpesialisasiSimpleDTO struct {
	SpesialisasiID   int    `json:"spesialisasi_id"`
	KodeSpesialisasi string `json:"kode_spesialisasi"`
	NamaSpesialisasi string `json:"nama_spesialisasi"`
}
