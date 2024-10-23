package domain

import "time"

type UsulanMusrebang struct {
	Id        string
	Usulan    string
	Alamat    string
	Uraian    string
	Tahun     string
	RekinId   string
	PegawaiId string
	KodeOpd   string
	IsActive  bool
	Status    string
	CreatedAt time.Time
}
