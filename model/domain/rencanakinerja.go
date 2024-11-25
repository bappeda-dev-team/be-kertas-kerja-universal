package domain

import (
	"time"
)

type RencanaKinerja struct {
	Id                   string
	IdPohon              int
	NamaPohon            string
	NamaRencanaKinerja   string
	Tahun                string
	StatusRencanaKinerja string
	Catatan              string
	KodeOpd              string
	NamaOpd              string
	PegawaiId            string
	NamaPegawai          string
	CreatedAt            time.Time
	Indikator            []Indikator
	KodeSubKegiatan      string
}
