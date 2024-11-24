package pegawai

type PegawaiResponse struct {
	Id          string `json:"id"`
	NamaPegawai string `json:"nama_pegawai"`
	Nip         string `json:"nip,omitempty"`
	KodeOpd     string `json:"kode_opd,omitempty"`
}
