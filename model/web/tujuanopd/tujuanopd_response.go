package tujuanopd

type TujuanOpdResponse struct {
	Id               int                 `json:"id_tujuan_opd"`
	KodeBidangUrusan string              `json:"kode_bidang_urusan"`
	KodeOpd          string              `json:"kode_opd"`
	NamaOpd          string              `json:"nama_opd"`
	Tujuan           string              `json:"tujuan"`
	RumusPerhitungan string              `json:"rumus_perhitungan"`
	SumberData       string              `json:"sumber_data"`
	Periode          PeriodeResponse     `json:"periode"`
	Indikator        []IndikatorResponse `json:"indikator"`
}

type IndikatorResponse struct {
	Id            string           `json:"id"`
	IdTujuanOpd   int              `json:"id_tujuan_opd"`
	NamaIndikator string           `json:"nama_indikator"`
	Target        []TargetResponse `json:"target"`
}

type TargetResponse struct {
	Id              string `json:"id"`
	IndikatorId     string `json:"indikator_id"`
	Tahun           string `json:"tahun"`
	TargetIndikator string `json:"target"`
	SatuanIndikator string `json:"satuan"`
}

type PeriodeResponse struct {
	Id         int    `json:"id"`
	TahunAwal  string `json:"tahun_awal"`
	TahunAkhir string `json:"tahun_akhir"`
}
