package rencanaaksi

type RencanaAksiCreateRequest struct {
	RencanaKinerjaId string `json:"rekin_id"`
	PegawaiId        string `json:"pegawai_id"`
	Urutan           int    `json:"urutan"`
	NamaRencanaAksi  string `json:"nama_rencana_aksi"`
}
