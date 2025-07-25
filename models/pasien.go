package models

type Pasien struct {
	PasienID      int     `json:"pasien_id"`
	Nama          string  `json:"nama"`
	TanggalLahir  string  `json:"tanggal_lahir"`
	JenisKelamin  string  `json:"jenis_kelamin"`
	Alamat        string  `json:"alamat"`
	NoHP          string  `json:"no_hp"`
	NoTelepon     *string `json:"no_telepon"`
	KTP           string  `json:"ktp"`
	Email         *string `json:"email"`
	GolonganDarah *string `json:"golongan_darah"`
}
