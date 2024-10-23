package dasarhukum

import "ekak_kabupaten_madiun/model/web"

type DasarHukumResponse struct {
	Id               string             `json:"id"`
	RekinId          string             `json:"rencana_kinerja_id"`
	PegawaiId        string             `json:"pegawai_id"`
	Urutan           int                `json:"urutan"`
	PeraturanTerkait string             `json:"peraturan_terkait"`
	Uraian           string             `json:"uraian"`
	Action           []web.ActionButton `json:"action,omitempty"`
}
