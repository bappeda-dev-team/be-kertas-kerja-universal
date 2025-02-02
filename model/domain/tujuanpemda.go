package domain

type TujuanPemda struct {
	Id               int
	TujuanPemdaId    int
	NamaTujuanPemda  string
	PeriodeId        int
	Periode          Periode
	RumusPerhitungan string
	SumberData       string
	Indikator        []Indikator
}
