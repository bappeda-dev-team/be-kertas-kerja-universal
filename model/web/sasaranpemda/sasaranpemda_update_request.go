package sasaranpemda

type SasaranPemdaUpdateRequest struct {
	Id            int    `json:"id"`
	TujuanPemdaId int    `json:"tujuan_pemda_id"`
	SubtemaId     int    `json:"subtema_id"`
	SasaranPemda  string `json:"sasaran_pemda"`
	// PeriodeId        int                      `json:"periode_id"`
	// Indikator        []IndikatorUpdateRequest `json:"indikator"`
}

// type IndikatorUpdateRequest struct {
// 	Id             string                `json:"id"`
// 	SasaranPemdaId string                `json:"sasaran_id"`
// 	Indikator      string                `json:"indikator"`
// 	Target         []TargetUpdateRequest `json:"target"`
// }

// type TargetUpdateRequest struct {
// 	Id     string `json:"id"`
// 	Target string `json:"target"`
// 	Satuan string `json:"satuan"`
// 	Tahun  string `json:"tahun"`
// }
