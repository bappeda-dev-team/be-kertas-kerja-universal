package jabatan

type JabatanCreateRequest struct {
	NamaJabatan  string
	KelasJabatan string
	JenisJabatan string
	NilaiJabatan int
	KodeOpd      string
	IndexJabatan int
	Tahun        string
	Esselon      string
}
