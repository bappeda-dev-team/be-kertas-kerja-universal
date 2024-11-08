package subkegiatan

import "ekak_kabupaten_madiun/model/web"

type SubKegiatanResponse struct {
	Id                   string                         `json:"id"`
	PegawaiId            string                         `json:"pegawai_id"`
	NamaSubKegiatan      string                         `json:"nama_sub_kegiatan"`
	KodeOpd              string                         `json:"kode_opd"`
	Tahun                string                         `json:"tahun"`
	Indikator            []IndikatorResponse            `json:"indikator,omitempty"`
	IndikatorSubkegiatan []IndikatorSubKegiatanResponse `json:"indikator_subkegiatan,omitempty"`
	PaguSubKegiatan      []PaguSubKegiatanResponse      `json:"pagu,omitempty"`
	Action               []web.ActionButton             `json:"action,omitempty"`
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

type IndikatorSubKegiatanResponse struct {
	Id            string `json:"id"`
	SubKegiatanId string `json:"sub_kegiatan_id"`
	NamaIndikator string `json:"indikator"`
}

type PaguSubKegiatanResponse struct {
	Id            string `json:"id"`
	SubKegiatanId string `json:"sub_kegiatan_id"`
	JenisPagu     string `json:"jenis"`
	PaguAnggaran  int    `json:"pagu_anggaran"`
	Tahun         string `json:"tahun"`
}
