package domain

import "time"

type PohonKinerja struct {
	Id         int
	Parent     int
	NamaPohon  string
	KodeOpd    string
	NamaOpd    string
	Keterangan string
	Tahun      string
	JenisPohon string
	LevelPohon int
	CreatedAt  time.Time
	Indikator  []Indikator
}
