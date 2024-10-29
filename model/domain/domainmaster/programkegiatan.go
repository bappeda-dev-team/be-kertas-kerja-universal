package domainmaster

type ProgramKegiatan struct {
	Id          string
	KodeProgram string
	NamaProgram string
	KodeOPD     string
	IsActive    bool //default true
	Tahun       string
	// urusan->bidangurusan->program->kegiatan->subkegiatan
}
