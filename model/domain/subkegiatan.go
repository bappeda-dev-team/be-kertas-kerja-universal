package domain

import "time"

type SubKegiatan struct {
	Id                   string
	RekinId              string
	NamaSubKegiatan      string
	KodeOpd              string
	Tahun                string
	PegawaiId            string
	CreatedAt            time.Time
	IndikatorSubKegiatan []IndikatorSubKegiatan
	PaguSubKegiatan      []PaguSubKegiatan
}

type IndikatorSubKegiatan struct {
	Id            string
	SubKegiatanId string
	NamaIndikator string
}

type PaguSubKegiatan struct {
	Id            string
	SubKegiatanId string
	JenisPagu     string
	PaguAnggaran  int
	Tahun         string
}
