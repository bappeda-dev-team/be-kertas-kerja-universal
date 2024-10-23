package subkegiatan

type SubKegiatanUpdateRequest struct {
	Id              string `json:"id"`
	NamaSubKegiatan string `json:"nama_subkegiatan"`
	KodeOpd         string `json:"kode_opd"`
	Tahun           string `json:"tahun"`
}
