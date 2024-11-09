package subkegiatan

type SubKegiatanCreateRequest struct {
	PegawaiId       string                   `json:"pegawai_id"`
	NamaSubKegiatan string                   `json:"nama_subkegiatan" validate:"required"`
	KodeOpd         string                   `json:"kode_opd" validate:"required"`
	Tahun           string                   `json:"tahun" validate:"required"`
	Indikator       []IndikatorCreateRequest `json:"indikator"`
}

type IndikatorCreateRequest struct {
	Id            string                `json:"id_indikator"`
	NamaIndikator string                `json:"indikator"`
	Target        []TargetCreateRequest `json:"targets"`
}

type TargetCreateRequest struct {
	Id              string `json:"id_target"`
	IndikatorId     string `json:"indikator_id"`
	TargetIndikator string `json:"target"`
	SatuanIndikator string `json:"satuan"`
}
