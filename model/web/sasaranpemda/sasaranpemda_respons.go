package sasaranpemda

type SasaranPemdaResponse struct {
	Id               int                 `json:"id"`
	SubtemaId        int                 `json:"subtema_id,omitempty"`
	NamaSubtema      string              `json:"nama_subtema,omitempty"`
	SasaranPemda     string              `json:"sasaran_pemda,omitempty"`
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

type SasaranPemdaWithPokinResponse struct {
	PokinId      int                   `json:"pokin_id"`
	NamaPohon    string                `json:"nama_tematik"`
	JenisPohon   string                `json:"jenis_pohon"`
	LevelPohon   int                   `json:"level_pohon"`
	Keterangan   string                `json:"keterangan"`
	TahunPokin   string                `json:"tahun_pokin"`
	SasaranPemda *SasaranPemdaResponse `json:"sasaran_pemda,omitempty"`
}
