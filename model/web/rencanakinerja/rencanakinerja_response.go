package rencanakinerja

import (
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/opdmaster"
)

type RencanaKinerjaResponse struct {
	Id                   string                      `json:"id_rencana_kinerja"`
	IdPohon              int                         `json:"id_pohon"`
	NamaPohon            string                      `json:"nama_pohon"`
	NamaRencanaKinerja   string                      `json:"nama_rencana_kinerja"`
	Tahun                string                      `json:"tahun"`
	StatusRencanaKinerja string                      `json:"status_rencana_kinerja"`
	Catatan              string                      `json:"catatan"`
	KodeOpd              opdmaster.OpdResponseForAll `json:"operasioanl_daerah"`
	PegawaiId            string                      `json:"pegawai_id"`
	NamaPegawai          string                      `json:"nama_pegawai"`
	Indikator            []IndikatorResponse         `json:"indikator"`
	// SubKegiatan          subkegiatan.SubKegiatanResponse `json:"sub_kegiatan,omitempty"`
	Action []web.ActionButton `json:"action"`
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
	TargetIndikator string `json:"target"`
	SatuanIndikator string `json:"satuan"`
}
