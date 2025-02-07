package sasaranpemda

type SasaranPemdaResponse struct {
	Id               int                 `json:"id"`
	SasaranPemdaId   int                 `json:"sasaran_pemda_id,omitempty"`
	NamaSasaranPemda string              `json:"nama_sasaran_pemda,omitempty"`
	Periode          PeriodeResponse     `json:"periode,omitempty"`
	RumusPerhitungan string              `json:"rumus_perhitungan,omitempty"`
	SumberData       string              `json:"sumber_data,omitempty"`
	Indikator        []IndikatorResponse `json:"indikator,omitempty"`
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
