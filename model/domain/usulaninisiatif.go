package domain

import "time"

type UsulanInisiatif struct {
	Id        string
	Usulan    string
	Manfaat   string
	Uraian    string
	Tahun     string
	RekinId   string
	PegawaiId string
	KodeOpd   string
	IsActive  bool
	Status    string
	CreatedAt time.Time
}
