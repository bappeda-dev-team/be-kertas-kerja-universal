package domain

type SasaranOpd struct {
	Id                int
	IdPohon           int
	NamaPohon         string
	JenisPohon        string
	LevelPohon        int
	TahunPohon        string
	TahunAwalPeriode  string
	TahunAkhirPeriode string
	JenisPeriode      string
	Pelaksana         []PelaksanaPokin
	RencanaKinerja    []RencanaKinerja
}

type ManualIKSasaranOpd struct {
	IndikatorId string
	Formula     string
	SumberData  string
	ManualIK    *ManualIKSasaranOpd
	Target      []Target
}
