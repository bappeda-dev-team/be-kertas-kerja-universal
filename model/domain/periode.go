package domain

type Periode struct {
	Id         int
	TahunAwal  string
	TahunAkhir string
}

type TahunPeriode struct {
	Id        int
	IdPeriode int
	Tahun     string
}
