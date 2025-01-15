package rencanakinerja

import (
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/opdmaster"
)

type RencanaKinerjaResponse struct {
	Id                   string                      `json:"id_rencana_kinerja,omitempty"`
	IdPohon              int                         `json:"id_pohon,omitempty"`
	NamaPohon            string                      `json:"nama_pohon,omitempty"`
	NamaRencanaKinerja   string                      `json:"nama_rencana_kinerja,omitempty"`
	Tahun                string                      `json:"tahun,omitempty"`
	StatusRencanaKinerja string                      `json:"status_rencana_kinerja,omitempty"`
	Catatan              string                      `json:"catatan,omitempty"`
	KodeOpd              opdmaster.OpdResponseForAll `json:"operasioanl_daerah,omitempty"`
	PegawaiId            string                      `json:"pegawai_id,omitempty"`
	NamaPegawai          string                      `json:"nama_pegawai,omitempty"`
	Indikator            []IndikatorResponse         `json:"indikator,omitempty"`
	// SubKegiatan          subkegiatan.SubKegiatanResponse `json:"sub_kegiatan,omitempty"`
	Action []web.ActionButton `json:"action,omitempty"`
}

type IndikatorResponse struct {
	Id               string           `json:"id_indikator,omitempty"`
	RencanaKinerjaId string           `json:"rencana_kinerja_id,omitempty"`
	NamaIndikator    string           `json:"nama_indikator,omitempty"`
	Target           []TargetResponse `json:"targets,omitempty"`
	ManualIK         *DataOutput      `json:"data_output,omitempty"`
}

type TargetResponse struct {
	Id              string `json:"id_target,omitempty"`
	IndikatorId     string `json:"indikator_id,omitempty"`
	TargetIndikator string `json:"target,omitempty"`
	SatuanIndikator string `json:"satuan,omitempty"`
}
