package iku

import "time"

type IkuResponse struct {
	IndikatorId  string           `json:"indikator_id"`
	Sumber       string           `json:"sumber"`
	Indikator    string           `json:"indikator"`
	CreatedAt    time.Time        `json:"created_at"`
	TahunAwal    string           `json:"tahun_awal"`
	TahunAkhir   string           `json:"tahun_akhir"`
	JenisPeriode string           `json:"jenis_periode"`
	Target       []TargetResponse `json:"target"`
}

type TargetResponse struct {
	Target string `json:"target"`
	Satuan string `json:"satuan"`
	Tahun  string `json:"tahun"`
}
