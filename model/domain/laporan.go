package domain

type OpdSupportingPokin struct {
	Id         int
	Parent     int
	LevelPohon int
	JenisPohon string
	NamaPohon  string
	KodeOpd    string
	NamaOpd    string
	Keterangan string
	Tahun      string
	Indikators []Indikator
	Status     string
}
