package domainmaster

import (
	"ekak_kabupaten_madiun/model/domain"
	"time"
)

type Kegiatan struct {
	Id           string
	NamaKegiatan string
	KodeKegiatan string
	KodeOPD      string
	CreatedAt    time.Time
	Indikator    []domain.Indikator
}
