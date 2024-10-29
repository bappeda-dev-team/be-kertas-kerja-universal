package inovasi

import "ekak_kabupaten_madiun/model/web"

type InovasiResponse struct {
	Id                    string             `json:"id"`
	RekinId               string             `json:"rencana_kinerja_id"`
	PegawaiId             string             `json:"pegawai_id"`
	KodeOpd               string             `json:"kode_opd"`
	JudulInovasi          string             `json:"judul_inovasi"`
	JenisInovasi          string             `json:"jenis_inovasi"`
	GambaranNilaiKebaruan string             `json:"gambaran_nilai_kebaruan"`
	Action                []web.ActionButton `json:"action,omitempty"`
}
