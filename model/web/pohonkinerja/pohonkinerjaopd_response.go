package pohonkinerja

import "ekak_kabupaten_madiun/model/web/opdmaster"

type PohonKinerjaOpdResponse struct {
	Id         int    `json:"id"`
	Parent     string `json:"parent"`
	NamaPohon  string `json:"nama_pohon"`
	JenisPohon string `json:"jenis_pohon"`
	LevelPohon int    `json:"level_pohon"`
	KodeOpd    string `json:"kode_opd,omitempty"`
	NamaOpd    string `json:"nama_opd,omitempty"`
	Keterangan string `json:"keterangan,omitempty"`
	Tahun      string `json:"tahun,omitempty"`
}

type PohonKinerjaOpdAllResponse struct {
	KodeOpd    string                 `json:"kode_opd"`
	NamaOpd    string                 `json:"nama_opd"`
	Tahun      string                 `json:"tahun"`
	Strategics []StrategicOpdResponse `json:"strategic"`
}

type StrategicOpdResponse struct {
	Id         int                         `json:"id"`
	Parent     *int                        `json:"parent"`
	Strategi   string                      `json:"strategi"`
	Keterangan string                      `json:"keterangan"`
	KodeOpd    opdmaster.OpdResponseForAll `json:"perangkat_daerah"`
	Tacticals  []TacticalOpdResponse       `json:"tacticals,omitempty"`
}

type TacticalOpdResponse struct {
	Id           int                         `json:"id"`
	Parent       int                         `json:"parent"`
	Strategi     string                      `json:"strategi"`
	Keterangan   string                      `json:"keterangan"`
	KodeOpd      opdmaster.OpdResponseForAll `json:"perangkat_daerah"`
	Operationals []OperationalOpdResponse    `json:"operational"`
}

type OperationalOpdResponse struct {
	Id         int                         `json:"id"`
	Parent     int                         `json:"parent"`
	Strategi   string                      `json:"strategi"`
	Keterangan string                      `json:"keterangan"`
	KodeOpd    opdmaster.OpdResponseForAll `json:"perangkat_daerah"`
}
