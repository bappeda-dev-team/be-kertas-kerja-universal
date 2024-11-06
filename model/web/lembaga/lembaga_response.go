package lembaga

type LembagaResponse struct {
	Id          string `json:"id"`
	KodeLembaga string `json:"kode_lembaga"`
	NamaLembaga string `json:"nama_lembaga"`
	IsActive    bool   `json:"is_active"`
}
