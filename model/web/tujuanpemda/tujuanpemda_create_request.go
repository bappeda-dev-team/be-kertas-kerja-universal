package tujuanpemda

type TujuanPemdaCreateRequest struct {
	TujuanPemda string                   `json:"tujuan_pemda"`
	TematikId   int                      `json:"tema_id"`
	PeriodeId   int                      `json:"periode_id"`
	Indikator   []IndikatorCreateRequest `json:"indikator"`
}

type IndikatorCreateRequest struct {
	Indikator        string                `json:"indikator"`
	RumusPerhitungan string                `json:"rumus_perhitungan"`
	SumberData       string                `json:"sumber_data"`
	Target           []TargetCreateRequest `json:"target"`
}

type TargetCreateRequest struct {
	Target string `json:"target"`
	Satuan string `json:"satuan"`
	Tahun  string `json:"tahun"`
}
