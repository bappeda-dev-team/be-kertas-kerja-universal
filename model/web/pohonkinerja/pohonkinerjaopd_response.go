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
	Pelaksana  []PelaksanaOpdResponse `json:"pelaksana"`
}

type PohonKinerjaOpdAllResponse struct {
	KodeOpd    string                 `json:"kode_opd"`
	NamaOpd    string                 `json:"nama_opd"`
	Tahun      string                 `json:"tahun"`
	Strategics []StrategicOpdResponse `json:"childs"`
}

type StrategicOpdResponse struct {
	Id         int                         `json:"id"`
	Parent     *int                        `json:"parent"`
	Strategi   string                      `json:"nama_pohon"`
	JenisPohon string                      `json:"jenis_pohon"`
	LevelPohon int                         `json:"level_pohon"`
	Keterangan string                      `json:"keterangan"`
	KodeOpd    opdmaster.OpdResponseForAll `json:"perangkat_daerah"`
	Pelaksana  []PelaksanaOpdResponse      `json:"pelaksana"`
	Tacticals  []TacticalOpdResponse       `json:"childs,omitempty"`
}

type TacticalOpdResponse struct {
	Id           int                         `json:"id"`
	Parent       int                         `json:"parent"`
	Strategi     string                      `json:"nama_pohon"`
	JenisPohon   string                      `json:"jenis_pohon"`
	LevelPohon   int                         `json:"level_pohon"`
	Keterangan   string                      `json:"keterangan"`
	KodeOpd      opdmaster.OpdResponseForAll `json:"perangkat_daerah"`
	Pelaksana    []PelaksanaOpdResponse      `json:"pelaksana"`
	Operationals []OperationalOpdResponse    `json:"childs,omitempty"`
}

type OperationalOpdResponse struct {
	Id         int                         `json:"id"`
	Parent     int                         `json:"parent"`
	Strategi   string                      `json:"nama_pohon"`
	JenisPohon string                      `json:"jenis_pohon"`
	LevelPohon int                         `json:"level_pohon"`
	Keterangan string                      `json:"keterangan"`
	KodeOpd    opdmaster.OpdResponseForAll `json:"perangkat_daerah"`
	Pelaksana  []PelaksanaOpdResponse      `json:"pelaksana"`
}

type PelaksanaOpdResponse struct {
	Id             string `json:"id_pelaksana"`
	PohonKinerjaId string `json:"pohon_kinerja_id,omitempty"`
	PegawaiId      string `json:"pegawai_id"`
	NamaPegawai    string `json:"nama_pegawai"`
}
