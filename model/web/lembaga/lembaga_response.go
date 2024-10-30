package lembaga

type LembagaResponse struct {
	Id          string `json:"id"`
	NamaLembaga string `json:"nama_lembaga"`
	IsActive    bool   `json:"is_active"`
}
