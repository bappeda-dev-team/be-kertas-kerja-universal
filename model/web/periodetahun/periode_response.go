package periodetahun

type PeriodeResponse struct {
	Id         int      `json:"id"`
	TahunAwal  string   `json:"tahun_awal"`
	TahunAkhir string   `json:"tahun_akhir"`
	TahunList  []string `json:"tahun_list"`
}

type TahunPeriodeResponse struct {
	Id    int    `json:"id"`
	Tahun string `json:"tahun"`
}
