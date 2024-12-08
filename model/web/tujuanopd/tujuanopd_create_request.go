package tujuanopd

type TujuanOpdCreateRequest struct {
	KodeOpd          string                   `json:"kode_opd"`
	Tujuan           string                   `json:"tujuan"`
	RumusPerhitungan string                   `json:"rumus_perhitungan"`
	SumberData       string                   `json:"sumber_data"`
	TahunAwal        string                   `json:"tahun_awal"`
	TahunAkhir       string                   `json:"tahun_akhir"`
	Indikator        []IndikatorCreateRequest `json:"indikator"`
}

type IndikatorCreateRequest struct {
	IdTujuanOpd string                `json:"id_tujuan_opd"`
	Indikator   string                `json:"indikator"`
	Target      []TargetCreateRequest `json:"target"`
}

type TargetCreateRequest struct {
	IndikatorId string `json:"indikator_id"`
	Target      string `json:"target"`
	Tahun       string `json:"tahun"`
	Satuan      string `json:"satuan"`
}
