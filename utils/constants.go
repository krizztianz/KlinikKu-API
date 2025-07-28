package utils

// Gender Enum
const (
	GenderLakiLaki  = "L"
	GenderPerempuan = "P"
)

// BloodType Enum
const (
	BloodTypeA  = "A"
	BloodTypeB  = "B"
	BloodTypeAB = "AB"
	BloodTypeO  = "O"
)

// VisitType Enum
const (
	VisitTypeUmum       = "Umum"
	VisitTypeRawatJalan = "Rawat Jalan"
	VisitTypeKontrol    = "Kontrol"
)

// VisitStatus Enum
const (
	VisitStatusRegistrasi    = "registrasi"
	VisitStatusPemeriksaan   = "pemeriksaan"
	VisitStatusMenungguResep = "menunggu_resep"
	VisitStatusSelesai       = "selesai"
	VisitStatusBatal         = "batal"
)

// VisitPriority Enum
const (
	VisitPriorityNormal    = "normal"
	VisitPriorityPrioritas = "prioritas"
)

// ResepStatus Enum
const (
	ResepStatusMenunggu = "menunggu"
	ResepStatusSelesai  = "selesai"
)
