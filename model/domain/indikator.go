package domain

import (
	"database/sql"
	"time"
)

type Indikator struct {
	Id               string
	PokinId          string
	ProgramId        string
	RencanaKinerjaId string
	KegiatanId       string
	SubKegiatanId    string
	TujuanOpdId      int
	TujuanPemdaId    int
	SasaranPemdaId   int
	Indikator        string
	Tahun            string
	CloneFrom        string
	Sumber           string
	ParentId         int
	ParentName       string
	CreatedAt        time.Time
	Target           []Target
	RencanaKinerja   RencanaKinerja
	RumusPerhitungan sql.NullString
	SumberData       sql.NullString
}
