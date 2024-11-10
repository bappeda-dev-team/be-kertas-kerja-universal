package subkegiatan

type SubKegiatanTerpilihResponse struct {
	Id              string              `json:"id,omitempty"`
	KodeSubKegiatan SubKegiatanResponse `json:"kode_subkegiatan"`
}
