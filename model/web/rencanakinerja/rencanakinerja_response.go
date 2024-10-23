package rencanakinerja

import "ekak_kabupaten_madiun/model/web"

type RencanaKinerjaResponse struct {
	Id                   string              `json:"id_rencana_kinerja"`
	NamaRencanaKinerja   string              `json:"nama_rencana_kinerja"`
	Tahun                string              `json:"tahun"`
	StatusRencanaKinerja string              `json:"status_rencana_kinerja"`
	Catatan              string              `json:"catatan"`
	KodeOpd              string              `json:"kode_opd"`
	PegawaiId            string              `json:"pegawai_id"`
	Indikator            []IndikatorResponse `json:"indikator"`
	Action               []web.ActionButton  `json:"action"`
}

type IndikatorResponse struct {
	Id               string           `json:"id_indikator"`
	RencanaKinerjaId string           `json:"rencana_kinerja_id"`
	NamaIndikator    string           `json:"nama_indikator"`
	Target           []TargetResponse `json:"targets"`
}

type TargetResponse struct {
	Id              string `json:"id_target"`
	IndikatorId     string `json:"indikator_id"`
	TargetIndikator int    `json:"target"`
	SatuanIndikator string `json:"satuan"`
}
