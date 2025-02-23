package domain

type SasaranOpd struct {
	Id                  int
	IdPohon             int
	NamaPohon           string
	JenisPohon          string
	LevelPohon          int
	IdRencanaKinerja    string
	NamaRencanaKinerja  string
	PegawaiId           string
	Pelaksana           []PelaksanaPokin
	IndikatorSasaranOpd []Indikator
	TahunPohon          string
	TahunAwalRencana    string
	TahunAkhirRencana   string
	TahunAwalPeriode    string
	TahunAkhirPeriode   string
	NipRencanaKinerja   string
}

// type PelaksanaPokin struct {
//     Id        string
//     PegawaiId string
// }

// type Indikator struct {
//     Id        string
//     Indikator string
//     ManualIK  *ManualIK
//     Target    []Target
// }

type ManualIKSasaranOpd struct {
	IndikatorId string
	Formula     string
	SumberData  string
	ManualIK    *ManualIKSasaranOpd
	Target      []Target
}

// type Target struct {
//     Id          string
//     IndikatorId string
//     Target      string
//     Satuan      string
//     Tahun       string
// }
