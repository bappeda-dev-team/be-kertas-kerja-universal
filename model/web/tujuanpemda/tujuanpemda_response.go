package tujuanpemda

type TujuanPemdaResponse struct {
	Id               int                 `json:"id"`
	TujuanPemdaId    int                 `json:"tujuan_pemda_id"`
	NamaTujuanPemda  string              `json:"nama_tujuan_pemda"`
	Periode          PeriodeResponse     `json:"periode"`
	RumusPerhitungan string              `json:"rumus_perhitungan"`
	SumberData       string              `json:"sumber_data"`
	Indikator        []IndikatorResponse `json:"indikator"`
}

type IndikatorResponse struct {
	Id        string           `json:"id"`
	Indikator string           `json:"indikator"`
	Target    []TargetResponse `json:"target"`
}

type TargetResponse struct {
	Id     string `json:"id"`
	Target string `json:"target"`
	Satuan string `json:"satuan"`
	Tahun  string `json:"tahun"`
}

type PeriodeResponse struct {
	TahunAwal  string `json:"tahun_awal"`
	TahunAkhir string `json:"tahun_akhir"`
}
