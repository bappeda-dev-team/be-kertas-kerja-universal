package tujuanopd

type TujuanOpdCreateRequest struct {
	KodeOpd          string                   `json:"kode_opd"`
	KodeBidangUrusan string                   `json:"kode_bidang_urusan"`
	Tujuan           string                   `json:"tujuan"`
	PeriodeId        int                      `json:"periode_id"`
	Indikator        []IndikatorCreateRequest `json:"indikator"`
}

type IndikatorCreateRequest struct {
	IdTujuanOpd      string                `json:"id_tujuan_opd"`
	Indikator        string                `json:"indikator"`
	RumusPerhitungan string                `json:"rumus_perhitungan"`
	SumberData       string                `json:"sumber_data"`
	Target           []TargetCreateRequest `json:"target"`
}

type TargetCreateRequest struct {
	IndikatorId string `json:"indikator_id"`
	Target      string `json:"target"`
	Tahun       string `json:"tahun"`
	Satuan      string `json:"satuan"`
}
