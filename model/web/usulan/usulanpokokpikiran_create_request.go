package usulan

type UsulanPokokPikiranCreateRequest struct {
	Usulan    string `json:"usulan"`
	Alamat    string `json:"alamat"`
	Uraian    string `json:"uraian"`
	Tahun     string `json:"tahun"`
	RekinId   string `json:"rencana_kinerja_id"`
	PegawaiId string `json:"pegawai_id"`
	KodeOpd   string `json:"kode_opd"`
	Status    string `json:"status"`
}
