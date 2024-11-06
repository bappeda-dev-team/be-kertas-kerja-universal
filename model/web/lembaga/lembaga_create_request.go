package lembaga

type LembagaCreateRequest struct {
	KodeLembaga string `json:"kode_lembaga" validate:"required"`
	NamaLembaga string `json:"nama_lembaga" validate:"required"`
}
