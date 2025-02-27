package bidangurusanresponse

type BidangUrusanResponse struct {
	Id               string `json:"id,omitempty"`
	KodeBidangUrusan string `json:"kode_bidang_urusan"`
	NamaBidangUrusan string `json:"nama_bidang_urusan"`
	Tahun            string `json:"tahun"`
}
