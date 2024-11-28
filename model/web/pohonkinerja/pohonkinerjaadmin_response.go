package pohonkinerja

import "ekak_kabupaten_madiun/model/web/opdmaster"

type PohonKinerjaAdminResponse struct {
	Tahun   string            `json:"tahun,omitempty"`
	Tematik []TematikResponse `json:"tematiks"`
}

type PohonKinerjaAdminResponseData struct {
	Id         int                    `json:"id"`
	Parent     int                    `json:"parent,omitempty"`
	NamaPohon  string                 `json:"nama_pohon"`
	KodeOpd    string                 `json:"kode_opd,omitempty"`
	NamaOpd    string                 `json:"nama_opd,omitempty"`
	Keterangan string                 `json:"keterangan,omitempty"`
	Tahun      string                 `json:"tahun"`
	JenisPohon string                 `json:"jenis_pohon"`
	LevelPohon int                    `json:"level_pohon"`
	Status     string                 `json:"status"`
	Pelaksana  []PelaksanaOpdResponse `json:"pelaksana,omitempty"`
	Indikators []IndikatorResponse    `json:"indikators,omitempty"`
	// SubTematiks []SubtematikResponse `json:"sub_tematiks,omitempty"`
}

type TematikResponse struct {
	Id         int                 `json:"id"`
	Parent     *int                `json:"parent"`
	Tema       string              `json:"tema"`
	JenisPohon string              `json:"jenis_pohon"`
	LevelPohon int                 `json:"level_pohon"`
	Keterangan string              `json:"keterangan"`
	Indikators []IndikatorResponse `json:"indikators"`
	// SubTematiks []SubtematikResponse `json:"childs,omitempty"`
	// Strategics  []StrategicResponse  `json:"strategics,omitempty"`
	Child []interface{} `json:"childs,omitempty"`
}

type SubtematikResponse struct {
	Id         int                 `json:"id"`
	Parent     int                 `json:"parent"`
	Tema       string              `json:"tema"`
	JenisPohon string              `json:"jenis_pohon"`
	LevelPohon int                 `json:"level_pohon"`
	Keterangan string              `json:"keterangan"`
	Indikators []IndikatorResponse `json:"indikators"`
	// SubSubTematiks []SubSubTematikResponse `json:"childs,omitempty"`
	// Strategics     []StrategicResponse     `json:"strategics,omitempty"`
	Child []interface{} `json:"childs,omitempty"`
}

type SubSubTematikResponse struct {
	Id         int                 `json:"id"`
	Parent     int                 `json:"parent"`
	Tema       string              `json:"tema"`
	JenisPohon string              `json:"jenis_pohon"`
	LevelPohon int                 `json:"level_pohon"`
	Keterangan string              `json:"keterangan"`
	Indikators []IndikatorResponse `json:"indikators"`
	// SuperSubTematiks []SuperSubTematikResponse `json:"childs,omitempty"`
	// Strategics       []StrategicResponse       `json:"strategics,omitempty"`
	Child []interface{} `json:"childs,omitempty"`
}

type SuperSubTematikResponse struct {
	Id         int                 `json:"id"`
	Parent     int                 `json:"parent"`
	Tema       string              `json:"tema"`
	JenisPohon string              `json:"jenis_pohon"`
	LevelPohon int                 `json:"level_pohon"`
	Keterangan string              `json:"keterangan"`
	Indikators []IndikatorResponse `json:"indikators"`
	Childs     []interface{}       `json:"childs,omitempty"`
}

type StrategicResponse struct {
	Id         int                          `json:"id"`
	Parent     int                          `json:"parent"`
	Strategi   string                       `json:"tema"`
	JenisPohon string                       `json:"jenis_pohon"`
	LevelPohon int                          `json:"level_pohon"`
	Keterangan string                       `json:"keterangan"`
	KodeOpd    *opdmaster.OpdResponseForAll `json:"perangkat_daerah,omitempty"`
	Pelaksana  []PelaksanaOpdResponse       `json:"pelaksana,omitempty"`
	Indikators []IndikatorResponse          `json:"indikators"`
	Childs     []interface{}                `json:"childs,omitempty"`
}

type TacticalResponse struct {
	Id         int                          `json:"id"`
	Parent     int                          `json:"parent"`
	Strategi   string                       `json:"tema"`
	JenisPohon string                       `json:"jenis_pohon"`
	LevelPohon int                          `json:"level_pohon"`
	Keterangan *string                      `json:"keterangan"`
	KodeOpd    *opdmaster.OpdResponseForAll `json:"perangkat_daerah,omitempty"`
	Pelaksana  []PelaksanaOpdResponse       `json:"pelaksana,omitempty"`
	Indikators []IndikatorResponse          `json:"indikators"`
	Childs     []interface{}                `json:"childs,omitempty"`
}

type OperationalResponse struct {
	Id         int                          `json:"id"`
	Parent     int                          `json:"parent"`
	Strategi   string                       `json:"tema"`
	JenisPohon string                       `json:"jenis_pohon"`
	LevelPohon int                          `json:"level_pohon"`
	Keterangan *string                      `json:"keterangan"`
	KodeOpd    *opdmaster.OpdResponseForAll `json:"perangkat_daerah,omitempty"`
	Pelaksana  []PelaksanaOpdResponse       `json:"pelaksana,omitempty"`
	Indikators []IndikatorResponse          `json:"indikators"`
	Childs     []interface{}                `json:"childs,omitempty"`
}

type OperationalNResponse struct {
	Id         int                          `json:"id"`
	Parent     int                          `json:"parent"`
	Strategi   string                       `json:"tema"`
	JenisPohon string                       `json:"jenis_pohon"`
	LevelPohon int                          `json:"level_pohon"`
	Keterangan *string                      `json:"keterangan"`
	KodeOpd    *opdmaster.OpdResponseForAll `json:"perangkat_daerah,omitempty"`
	Pelaksana  []PelaksanaOpdResponse       `json:"pelaksana,omitempty"`
	Indikators []IndikatorResponse          `json:"indikators"`
	Childs     []OperationalNResponse       `json:"childs,omitempty"`
}
