package domainmaster

import "time"

type BidangUrusan struct {
	Id               string
	KodeBidangUrusan string
	NamaBidangUrusan string
	CreatedAt        time.Time
}
