package domain

type Indikator struct {
	Id               string
	PokinId          string
	ProgramId        string
	RencanaKinerjaId string
	KegiatanId       string
	SubKegiatanId    string
	Indikator        string
	Tahun            string
	Target           []Target
}
