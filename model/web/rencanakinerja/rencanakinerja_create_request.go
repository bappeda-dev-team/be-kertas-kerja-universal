package rencanakinerja

type RencanaKinerjaCreateRequest struct {
	IdPohon              int                      `json:"id_pohon"`
	NamaRencanaKinerja   string                   `json:"nama_rencana_kinerja"`
	Tahun                string                   `json:"tahun"`
	StatusRencanaKinerja string                   `json:"status_rencana_kinerja"`
	Catatan              string                   `json:"catatan"`
	KodeOpd              string                   `json:"kode_opd"`
	PegawaiId            string                   `json:"pegawai_id"`
	Indikator            []IndikatorCreateRequest `json:"indikator"`
}

type IndikatorCreateRequest struct {
	NamaIndikator string                `json:"nama_indikator"`
	Tahun         string                `json:"tahun"`
	Target        []TargetCreateRequest `json:"target"`
}

type TargetCreateRequest struct {
	Tahun           string `json:"tahun"`
	Target          string `json:"target"`
	SatuanIndikator string `json:"satuan"`
}
