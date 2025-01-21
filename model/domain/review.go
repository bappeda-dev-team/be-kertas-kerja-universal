package domain

import "time"

type Review struct {
	Id             int
	IdPohonKinerja int
	Review         string
	Keterangan     string
	CreatedBy      string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
