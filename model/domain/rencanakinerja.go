package domain

import (
	"time"
)

type RencanaKinerja struct {
	Id                   string
	NamaRencanaKinerja   string
	Tahun                string
	StatusRencanaKinerja string
	Catatan              string
	KodeOpd              string
	NamaOpd              string
	PegawaiId            string
	CreatedAt            time.Time
	Indikator            []Indikator
	KodeSubKegiatan      string
}
