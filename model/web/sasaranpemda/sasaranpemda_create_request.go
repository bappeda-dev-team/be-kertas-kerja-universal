package sasaranpemda

type SasaranPemdaCreateRequest struct {
	SubtemaId     int    `json:"subtema_id"`
	TujuanPemdaId int    `json:"tujuan_pemda_id"`
	SasaranPemda  string `json:"sasaran_pemda"`
	// PeriodeId        int                      `json:"periode_id"`
	// Indikator []IndikatorCreateRequest `json:"indikator"`
}

// type IndikatorCreateRequest struct {
// 	Indikator string                `json:"indikator"`
// 	Target    []TargetCreateRequest `json:"target"`
// }

// type TargetCreateRequest struct {
// 	Target string `json:"target"`
// 	Satuan string `json:"satuan"`
// 	Tahun  string `json:"tahun"`
// }
