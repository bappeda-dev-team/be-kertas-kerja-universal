package domainmaster

type Jabatan struct {
	Id           string
	NamaJabatan  string
	KelasJabatan string
	JenisJabatan string
	NilaiJabatan int //default 0
	KodeOpd      string
	IndexJabatan int //default 0
	Tahun        int
	Esselon      string
}
