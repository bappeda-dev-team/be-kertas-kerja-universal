package rencanaaksi

import "ekak_kabupaten_madiun/model/web"

type PelaksanaanRencanaAksiResponse struct {
	Id            string             `json:"id"`
	RencanaAksiId string             `json:"rencana_aksi_id"`
	Bobot         int                `json:"bobot"`
	Bulan         int                `json:"bulan"`
	BobotAvail    int                `json:"bobot_tersedia,omitempty"`
	Action        []web.ActionButton `json:"action,omitempty"`
}
