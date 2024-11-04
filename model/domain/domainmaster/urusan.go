package domainmaster

import "time"

type Urusan struct {
	Id         string    `json:"id"`
	KodeUrusan string    `json:"kode_urusan"`
	NamaUrusan string    `json:"nama_urusan"`
	CreatedAt  time.Time `json:"created_at"`
}
