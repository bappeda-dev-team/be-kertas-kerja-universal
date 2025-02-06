package domain

type TujuanPemda struct {
	Id          int
	TematikId   int
	NamaTematik string
	TujuanPemda string
	// TujuanPemdaId    int
	// NamaTujuanPemda  string
	PeriodeId        int
	Periode          Periode
	RumusPerhitungan string
	SumberData       string
	Indikator        []Indikator
}
