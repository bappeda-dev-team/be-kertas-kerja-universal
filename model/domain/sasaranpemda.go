package domain

type SasaranPemda struct {
	Id               int
	SubtemaId        int
	NamaSubtema      string
	SasaranPemda     string
	PeriodeId        int
	Periode          Periode
	RumusPerhitungan string
	SumberData       string
	Indikator        []Indikator
}

type SasaranPemdaWithPokin struct {
	PokinId      int           `json:"pokin_id"`
	NamaPohon    string        `json:"nama_pohon"`
	JenisPohon   string        `json:"jenis_pohon"`
	LevelPohon   int           `json:"level_pohon"`
	KodeOpd      string        `json:"kode_opd"`
	Keterangan   string        `json:"keterangan"`
	TahunPokin   string        `json:"tahun_pokin"`
	SasaranPemda *SasaranPemda `json:"sasaran_pemda,omitempty"`
}
