package domain

type Indikator struct {
	Id               string
	PokinId          string
	ProgramId        string
	RencanaKinerjaId string
	KegiatanId       string
	SubKegiatanId    string
	TujuanOpdId      int
	TujuanPemdaId    int
	Indikator        string
	Tahun            string
	CloneFrom        string
	Target           []Target
	RencanaKinerja   RencanaKinerja
}
