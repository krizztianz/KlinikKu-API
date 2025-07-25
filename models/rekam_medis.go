package models

type RekamMedis struct {
	RekamMedisID     int    `json:"rekam_medis_id"`
	KunjunganID      int    `json:"kunjungan_id"`
	DokterID         int    `json:"dokter_id"`
	Anamnesa         string `json:"anamnesa"`
	PemeriksaanFisik string `json:"pemeriksaan_fisik"`
	Catatan          string `json:"catatan"`
}
