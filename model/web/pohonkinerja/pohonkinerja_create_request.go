package pohonkinerja

type PohonKinerjaCreateRequest struct {
	Parent     int    `json:"parent"`
	NamaPohon  string `json:"nama_pohon"`
	JenisPohon string `json:"jenis_pohon"`
	LevelPohon int    `json:"level_pohon"`
	KodeOpd    string `json:"kode_opd"`
	Keterangan string `json:"keterangan"`
	Tahun      string `json:"tahun"`
}

type PohonKinerjaAdminCreateRequest struct {
	Parent     int                      `json:"parent"`
	NamaPohon  string                   `json:"nama_pohon"`
	JenisPohon string                   `json:"jenis_pohon"`
	KodeOpd    string                   `json:"kode_opd,omitempty"`
	LevelPohon int                      `json:"level_pohon"`
	Keterangan string                   `json:"keterangan"`
	Tahun      string                   `json:"tahun"`
	Indikator  []IndikatorCreateRequest `json:"indikator"`
}

type PohonKinerjaAdminStrategicCreateRequest struct {
	IdToClone  int `json:"id"`
	Parent     int `json:"parent"`
	LevelPohon int `json:"level_pohon"`
}

type IndikatorCreateRequest struct {
	PohonKinerjaId int                   `json:"pohon_id"`
	NamaIndikator  string                `json:"indikator"`
	Target         []TargetCreateRequest `json:"target"`
}

type TargetCreateRequest struct {
	IndikatorId int    `json:"indikator_id"`
	Target      string `json:"target"`
	Satuan      string `json:"satuan"`
}
