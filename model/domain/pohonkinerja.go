package domain

type PohonKinerja struct {
	Id         string
	Title      string
	Parent     string
	KodeOpd    string
	Keterangan string
	Tahun      string
	JenisPohon string
	LevelPohon int
	Indikator  []Indikator
}
