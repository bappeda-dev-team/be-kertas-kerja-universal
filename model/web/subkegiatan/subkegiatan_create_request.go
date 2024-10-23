package subkegiatan

type SubKegiatanCreateRequest struct {
	RekinId         string `json:"rencana_kinerja_id"`
	PegawaiId       string `json:"pegawai_id"`
	NamaSubKegiatan string `json:"nama_subkegiatan" validate:"required"`
	KodeOpd         string `json:"kode_opd" validate:"required"`
	Tahun           string `json:"tahun" validate:"required"`
}
