package laporan

type OpdSupportingPokinResponseData struct {
	Tahun         string            `json:"tahun"`
	PohonKinerjas []PokinSupporting `json:"pohon_kinerjas"`
}

type PokinSupporting struct {
	Id         int                 `json:"id"`
	Parent     int                 `json:"parent,omitempty"`
	Tema       string              `json:"tema"`
	JenisPohon string              `json:"jenis_pohon"`
	Keterangan string              `json:"keterangan,omitempty"`
	Tahun      string              `json:"tahun"`
	LevelPohon int                 `json:"level_pohon"`
	Indikators []IndikatorResponse `json:"indikator,omitempty"`
	Childs     []PokinSupporting       `json:"childs,omitempty"`
}

type IndikatorResponse struct {
	Id            string           `json:"id_indikator"`
	IdPokin       string           `json:"id_pokin,omitempty"`
	NamaIndikator string           `json:"nama_indikator"`
	Targets       []TargetResponse `json:"targets"`
}

type TargetResponse struct {
	Id              string `json:"id_target"`
	IndikatorId     string `json:"indikator_id"`
	TargetIndikator string `json:"target"`
	SatuanIndikator string `json:"satuan"`
	TahunSasaran    string `json:"tahun,omitempty"`
}
