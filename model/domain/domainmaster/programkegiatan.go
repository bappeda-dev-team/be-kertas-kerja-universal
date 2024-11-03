package domainmaster

import "ekak_kabupaten_madiun/model/domain"

type ProgramKegiatan struct {
	Id          string
	KodeProgram string
	NamaProgram string
	KodeOPD     string
	IsActive    bool //default true
	Tahun       string
	Indikator   []domain.Indikator
}
