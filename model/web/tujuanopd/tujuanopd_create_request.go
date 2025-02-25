package tujuanopd

type TujuanOpdCreateRequest struct {
	KodeOpd          string                   `json:"kode_opd"`
	KodeBidangUrusan string                   `json:"kode_bidang_urusan"`
	Tujuan           string                   `json:"tujuan"`
	PeriodeId        int                      `json:"periode_id"`
	TahunAwal        string                   `json:"tahun_awal"`
	TahunAkhir       string                   `json:"tahun_akhir"`
	JenisPeriode     string                   `json:"jenis_periode"`
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
