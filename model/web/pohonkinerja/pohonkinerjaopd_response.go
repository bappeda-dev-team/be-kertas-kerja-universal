package pohonkinerja

import (
	"ekak_kabupaten_madiun/model/web/opdmaster"
)

type PohonKinerjaOpdResponse struct {
	Id         int                    `json:"id"`
	Parent     string                 `json:"parent"`
	NamaPohon  string                 `json:"nama_pohon"`
	JenisPohon string                 `json:"jenis_pohon"`
	LevelPohon int                    `json:"level_pohon"`
	KodeOpd    string                 `json:"kode_opd,omitempty"`
	NamaOpd    string                 `json:"nama_opd,omitempty"`
	Keterangan string                 `json:"keterangan,omitempty"`
	Tahun      string                 `json:"tahun,omitempty"`
	Status     string                 `json:"status"`
	Pelaksana  []PelaksanaOpdResponse `json:"pelaksana"`
	Indikator  []IndikatorResponse    `json:"indikator"`
}

type PohonKinerjaOpdAllResponse struct {
	KodeOpd    string                 `json:"kode_opd"`
	NamaOpd    string                 `json:"nama_opd"`
	Tahun      string                 `json:"tahun"`
	TujuanOpd  []TujuanOpdResponse    `json:"tujuan_opd"`
	Strategics []StrategicOpdResponse `json:"childs"`
}

type StrategicOpdResponse struct {
	Id           int                         `json:"id"`
	Parent       *int                        `json:"parent"`
	Strategi     string                      `json:"nama_pohon"`
	JenisPohon   string                      `json:"jenis_pohon"`
	LevelPohon   int                         `json:"level_pohon"`
	Keterangan   string                      `json:"keterangan"`
	Status       string                      `json:"status"`
	KodeOpd      opdmaster.OpdResponseForAll `json:"perangkat_daerah"`
	Pelaksana    []PelaksanaOpdResponse      `json:"pelaksana"`
	Indikator    []IndikatorResponse         `json:"indikator"`
	Tacticals    []TacticalOpdResponse       `json:"childs,omitempty"`
	Crosscutting []CrosscuttingOpdResponse   `json:"ccrosscutting_strategic,omitempty"`
}

type TacticalOpdResponse struct {
	Id           int                         `json:"id"`
	Parent       int                         `json:"parent"`
	Strategi     string                      `json:"nama_pohon"`
	JenisPohon   string                      `json:"jenis_pohon"`
	LevelPohon   int                         `json:"level_pohon"`
	Keterangan   string                      `json:"keterangan"`
	Status       string                      `json:"status"`
	KodeOpd      opdmaster.OpdResponseForAll `json:"perangkat_daerah"`
	Pelaksana    []PelaksanaOpdResponse      `json:"pelaksana"`
	Indikator    []IndikatorResponse         `json:"indikator"`
	Operationals []OperationalOpdResponse    `json:"childs,omitempty"`
	Crosscutting []CrosscuttingOpdResponse   `json:"ccrosscutting_tactical,omitempty"`
}

type OperationalOpdResponse struct {
	Id           int                         `json:"id"`
	Parent       int                         `json:"parent"`
	Strategi     string                      `json:"nama_pohon"`
	JenisPohon   string                      `json:"jenis_pohon"`
	LevelPohon   int                         `json:"level_pohon"`
	Keterangan   string                      `json:"keterangan"`
	Status       string                      `json:"status"`
	KodeOpd      opdmaster.OpdResponseForAll `json:"perangkat_daerah"`
	Pelaksana    []PelaksanaOpdResponse      `json:"pelaksana"`
	Indikator    []IndikatorResponse         `json:"indikator"`
	Childs       []OperationalNOpdResponse   `json:"childs,omitempty"`
	Crosscutting []CrosscuttingOpdResponse   `json:"ccrosscutting_operational,omitempty"`
}

type OperationalNOpdResponse struct {
	Id         int                         `json:"id"`
	Parent     int                         `json:"parent"`
	Strategi   string                      `json:"nama_pohon"`
	JenisPohon string                      `json:"jenis_pohon"`
	LevelPohon int                         `json:"level_pohon"`
	Keterangan string                      `json:"keterangan"`
	Status     string                      `json:"status"`
	KodeOpd    opdmaster.OpdResponseForAll `json:"perangkat_daerah"`
	Pelaksana  []PelaksanaOpdResponse      `json:"pelaksana"`
	Indikator  []IndikatorResponse         `json:"indikator"`
	Childs     []OperationalNOpdResponse   `json:"childs,omitempty"`
}

type PelaksanaOpdResponse struct {
	Id             string `json:"id_pelaksana"`
	PohonKinerjaId string `json:"pohon_kinerja_id,omitempty"`
	PegawaiId      string `json:"pegawai_id"`
	NamaPegawai    string `json:"nama_pegawai"`
}

type TujuanOpdResponse struct {
	Id        int                       `json:"id"`
	KodeOpd   string                    `json:"kode_opd"`
	Tujuan    string                    `json:"tujuan"`
	Indikator []IndikatorTujuanResponse `json:"indikator"`
}

type IndikatorTujuanResponse struct {
	Indikator string `json:"indikator"`
}
