package domain

type TujuanOpd struct {
	Id               int
	KodeOpd          string
	NamaOpd          string
	KodeBidangUrusan string
	Tujuan           string
	RumusPerhitungan string
	SumberData       string
	PeriodeId        Periode
	TahunAwal        string
	TahunAkhir       string
	Indikator        []Indikator
}
