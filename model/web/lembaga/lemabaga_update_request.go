package lembaga

type LembagaUpdateRequest struct {
	Id          string `json:"id"`
	KodeLembaga string `json:"kode_lembaga" validate:"required"`
	NamaLembaga string `json:"nama_lembaga" validate:"required"`
	IsActive    bool   `json:"is_active"`
}
