package bidangurusanresponse

type BidangUrusanUpdateRequest struct {
	Id               string `json:"id"`
	KodeBidangUrusan string `json:"kode_bidang_urusan"`
	NamaBidangUrusan string `json:"nama_bidang_urusan"`
}
