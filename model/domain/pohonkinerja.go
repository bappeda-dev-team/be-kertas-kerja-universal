package domain

type PohonKinerja struct {
	Id         int
	Parent     int
	NamaPohon  string
	KodeOpd    string
	Keterangan string
	Tahun      string
	JenisPohon string
	LevelPohon int
	Indikator  []Indikator
}
