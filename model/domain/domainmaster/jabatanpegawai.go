package domainmaster

type JabatanPegawai struct {
	Id        string
	IdJabatan string
	IdPegawai string
	Status    string
	IsActive  bool //default true
	Bulan     int
	Tahun     int
	Pangkat   string
	Golongan  string
}
