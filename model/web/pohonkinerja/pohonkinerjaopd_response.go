package pohonkinerja

type PohonKinerjaOpdResponse struct {
	Id         int    `json:"id"`
	Parent     string `json:"parent"`
	NamaPohon  string `json:"nama_pohon"`
	JenisPohon string `json:"jenis_pohon"`
	LevelPohon int    `json:"level_pohon"`
	KodeOpd    string `json:"kode_opd,omitempty"`
	Keterangan string `json:"keterangan,omitempty"`
	Tahun      string `json:"tahun,omitempty"`
}

type PohonKinerjaOpdAllResponse struct {
	KodeOpd       string                    `json:"kode_opd"`
	Tahun         string                    `json:"tahun"`
	NamaOpd       string                    `json:"nama_opd"`
	PohonKinerjas []PohonKinerjaOpdResponse `json:"pohon_kinerjas"`
}
