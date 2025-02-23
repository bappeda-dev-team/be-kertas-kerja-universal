package visimisipemda

type VisiPemdaUpdateRequest struct {
	Id                int    `json:"id" validate:"required"`
	Visi              string `json:"visi" validate:"required"`
	TahunAwalPeriode  string `json:"tahun_awal_periode" validate:"required"`
	TahunAkhirPeriode string `json:"tahun_akhir_periode" validate:"required"`
	JenisPeriode      string `json:"jenis_periode" validate:"required"`
	Keterangan        string `json:"keterangan"`
}
