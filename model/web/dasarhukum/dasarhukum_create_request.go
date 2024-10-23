package dasarhukum

type DasarHukumCreateRequest struct {
	RekinId          string `json:"rencana_kinerja_id"`
	PegawaiId        string `json:"pegawai_id"`
	Urutan           int    `json:"urutan"`
	PeraturanTerkait string `json:"peraturan_terkait"`
	Uraian           string `json:"uraian"`
}
