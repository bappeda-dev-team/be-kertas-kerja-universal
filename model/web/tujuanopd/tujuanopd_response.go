package tujuanopd

type TujuanOpdResponse struct {
	Id               int                 `json:"id_tujuan_opd"`
	KodeOpd          string              `json:"kode_opd"`
	NamaOpd          string              `json:"nama_opd"`
	Tujuan           string              `json:"tujuan"`
	RumusPerhitungan string              `json:"rumus_perhitungan"`
	SumberData       string              `json:"sumber_data"`
	TahunAwal        string              `json:"tahun_awal"`
	TahunAkhir       string              `json:"tahun_akhir"`
	Indikator        []IndikatorResponse `json:"indikator"`
}

type IndikatorResponse struct {
	Id            string           `json:"id_indikator"`
	IdTujuanOpd   string           `json:"id_tujuan_opd"`
	NamaIndikator string           `json:"nama_indikator"`
	Target        []TargetResponse `json:"targets"`
}

type TargetResponse struct {
	Id              string `json:"id_target"`
	IndikatorId     string `json:"indikator_id"`
	Tahun           string `json:"tahun"`
	TargetIndikator string `json:"target"`
	SatuanIndikator string `json:"satuan"`
}
