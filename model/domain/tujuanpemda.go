package domain

type TujuanPemda struct {
	Id          int
	TematikId   int
	NamaTematik string
	TujuanPemda string
	JenisPohon  string
	// TujuanPemdaId    int
	// NamaTujuanPemda  string
	PeriodeId int
	Periode   Periode
	Indikator []Indikator
}

type TujuanPemdaWithPokin struct {
	PokinId     int          `json:"pokin_id"`
	NamaPohon   string       `json:"nama_pohon"`
	JenisPohon  string       `json:"jenis_pohon"`
	LevelPohon  int          `json:"level_pohon"`
	KodeOpd     string       `json:"kode_opd"`
	Keterangan  string       `json:"keterangan"`
	TahunPokin  string       `json:"tahun_pokin"`
	TujuanPemda *TujuanPemda `json:"tujuan_pemda,omitempty"`
}
