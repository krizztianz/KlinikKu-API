package dto

type CreateDiagnosaRequest struct {
	KodeICD10    string `json:"kode_icd10" binding:"required"`
	NamaDiagnosa string `json:"nama_diagnosa" binding:"required"`
	Deskripsi    string `json:"deskripsi"`
}

type DiagnosaResponse struct {
	DiagnosaID   int    `json:"diagnosa_id"`
	KodeICD10    string `json:"kode_icd10"`
	NamaDiagnosa string `json:"nama_diagnosa"`
	Deskripsi    string `json:"deskripsi"`
}
