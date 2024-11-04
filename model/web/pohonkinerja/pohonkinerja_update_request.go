package pohonkinerja

type PohonKinerjaUpdateRequest struct {
	Id         int    `json:"id"`
	Parent     int    `json:"parent"`
	NamaPohon  string `json:"nama_pohon"`
	JenisPohon string `json:"jenis_pohon"`
	LevelPohon int    `json:"level_pohon"`
	KodeOpd    string `json:"kode_opd"`
	Keterangan string `json:"keterangan"`
	Tahun      string `json:"tahun"`
}

type PohonKinerjaAdminUpdateRequest struct {
	Id         int                      `json:"id"`
	Parent     int                      `json:"parent"`
	NamaPohon  string                   `json:"nama_pohon"`
	JenisPohon string                   `json:"jenis_pohon"`
	LevelPohon int                      `json:"level_pohon"`
	KodeOpd    string                   `json:"kode_opd"`
	Keterangan string                   `json:"keterangan"`
	Tahun      string                   `json:"tahun"`
	Indikator  []IndikatorUpdateRequest `json:"indikator"`
}

type IndikatorUpdateRequest struct {
	Id             string                `json:"id"`
	PohonKinerjaId int                   `json:"pohon_id"`
	NamaIndikator  string                `json:"indikator"`
	Target         []TargetUpdateRequest `json:"target"`
}

type TargetUpdateRequest struct {
	Id          string `json:"id"`
	IndikatorId int    `json:"indikator_id"`
	Target      string `json:"target"`
	Satuan      string `json:"satuan"`
}
