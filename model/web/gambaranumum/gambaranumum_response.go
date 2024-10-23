package gambaranumum

import "ekak_kabupaten_madiun/model/web"

type GambaranUmumResponse struct {
	Id           string             `json:"id"`
	RekinId      string             `json:"rencana_kinerja_id"`
	PegawaiId    string             `json:"pegawai_id"`
	Urutan       int                `json:"urutan"`
	GambaranUmum string             `json:"gambaran_umum"`
	Action       []web.ActionButton `json:"action,omitempty"`
}
