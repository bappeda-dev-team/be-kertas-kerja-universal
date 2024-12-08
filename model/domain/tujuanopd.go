package domain

type TujuanOpd struct {
	Id               int
	KodeOpd          string
	NamaOpd          string
	Tujuan           string
	RumusPerhitungan string
	SumberData       string
	TahunAwal        string
	TahunAkhir       string
	Indikator        []Indikator
}
