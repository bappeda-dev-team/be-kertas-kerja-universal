package rencanaaksi

import "ekak_kabupaten_madiun/model/web"

type RencanaAksiResponse struct {
	Id                     string                           `json:"id"`
	RencanaKinerjaId       string                           `json:"rekin_id"`
	KodeOpd                string                           `json:"kode_opd,omitempty"`
	Urutan                 int                              `json:"urutan"`
	NamaRencanaAksi        string                           `json:"nama_rencana_aksi"`
	PelaksanaanRencanaAksi []PelaksanaanRencanaAksiResponse `json:"pelaksanaan"`
	JumlahBobot            int                              `json:"jumlah_bobot,omitempty"`
	Action                 []web.ActionButton               `json:"action,omitempty"`
}
