package subkegiatan

import "ekak_kabupaten_madiun/model/web"

type SubKegiatanResponse struct {
	Id                   string                         `json:"id"`
	PegawaiId            string                         `json:"pegawai_id"`
	NamaSubKegiatan      string                         `json:"nama_sub_kegiatan"`
	KodeOpd              string                         `json:"kode_opd"`
	Tahun                string                         `json:"tahun"`
	IndikatorSubkegiatan []IndikatorSubKegiatanResponse `json:"indikator_subkegiatan"`
	PaguSubKegiatan      []PaguSubKegiatanResponse      `json:"pagu"`
	Action               []web.ActionButton             `json:"action,omitempty"`
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
