package domain

type Indikator struct {
	Id               string
	PokinId          string
	RencanaKinerjaId string
	Indikator        string
	Tahun            string
	Target           []Target
}
