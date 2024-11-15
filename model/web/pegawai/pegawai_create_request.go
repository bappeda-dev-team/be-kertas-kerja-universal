package pegawai

type PegawaiCreateRequest struct {
	NamaPegawai string `json:"nama_pegawai"`
	Nip         string `json:"nip"`
	KodeOpd     string `json:"kode_opd"`
}
