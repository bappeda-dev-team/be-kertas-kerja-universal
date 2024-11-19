package pegawai

type PegawaiResponse struct {
	Id          string `json:"id"`
	NamaPegawai string `json:"nama_pegawai"`
<<<<<<< HEAD
	Nip         string `json:"nip"`
	KodeOpd     string `json:"kode_opd"`
	NamaOpd     string `json:"nama_opd"`
=======
	Nip         string `json:"nip,omitempty"`
>>>>>>> develop_pokin
}
