package sasaranpemda

type SasaranPemdaCreateRequest struct {
	SasaranPemdaId   int                      `json:"sasaran_pemda_id"`
	PeriodeId        int                      `json:"periode_id"`
	RumusPerhitungan string                   `json:"rumus_perhitungan"`
	SumberData       string                   `json:"sumber_data"`
	Indikator        []IndikatorCreateRequest `json:"indikator"`
}

type IndikatorCreateRequest struct {
	Indikator string                `json:"indikator"`
	Target    []TargetCreateRequest `json:"target"`
}

type TargetCreateRequest struct {
	Target string `json:"target"`
	Satuan string `json:"satuan"`
	Tahun  string `json:"tahun"`
}
