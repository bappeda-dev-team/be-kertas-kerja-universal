package rencanakinerja

type RencanaKinerjaUpdateRequest struct {
	Id                   string                   `json:"id"`
	NamaRencanaKinerja   string                   `json:"nama_rencana_kinerja"`
	Tahun                string                   `json:"tahun"`
	StatusRencanaKinerja string                   `json:"status_rencana_kinerja"`
	Catatan              string                   `json:"catatan"`
	KodeOpd              string                   `json:"kode_opd"`
	PegawaiId            string                   `json:"pegawai_id"`
	Indikator            []IndikatorUpdateRequest `json:"indikator"`
}

type IndikatorUpdateRequest struct {
	Id               string                `json:"id_indikator"`
	RencanaKinerjaId string                `json:"rencana_kinerja_id"`
	Indikator        string                `json:"nama_indikator"`
	Tahun            string                `json:"tahun"`
	Target           []TargetUpdateRequest `json:"target"`
}

type TargetUpdateRequest struct {
	Id              string `json:"id_target"`
	IndikatorId     string `json:"indikator_id"`
	Tahun           string `json:"tahun"`
	Target          int    `json:"target"`
	SatuanIndikator string `json:"satuan"`
}
