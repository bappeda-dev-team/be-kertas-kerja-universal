package domain

type SasaranPemda struct {
	Id               int
	SasaranPemdaId   int
	NamaSasaranPemda string
	PeriodeId        int
	Periode          Periode
	RumusPerhitungan string
	SumberData       string
	Indikator        []Indikator
}
