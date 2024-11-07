package pohonkinerja

import "ekak_kabupaten_madiun/model/web/opdmaster"

type PohonKinerjaAdminResponse struct {
	Tahun   string            `json:"tahun"`
	Tematik []TematikResponse `json:"tematiks"`
}

type PohonKinerjaAdminResponseData struct {
	Id          int                  `json:"id"`
	Parent      int                  `json:"parent"`
	NamaPohon   string               `json:"nama_pohon"`
	KodeOpd     string               `json:"kode_opd,omitempty"`
	NamaOpd     string               `json:"nama_opd,omitempty"`
	Keterangan  string               `json:"keterangan"`
	Tahun       string               `json:"tahun"`
	JenisPohon  string               `json:"jenis_pohon"`
	LevelPohon  int                  `json:"level_pohon"`
	Indikators  []IndikatorResponse  `json:"indikators"`
	SubTematiks []SubtematikResponse `json:"sub_tematiks,omitempty"`
}

type TematikResponse struct {
	Id          int                  `json:"id"`
	Parent      *int                 `json:"parent"`
	Tema        string               `json:"tema"`
	Keterangan  string               `json:"keterangan"`
	Indikators  []IndikatorResponse  `json:"indikators"`
	SubTematiks []SubtematikResponse `json:"sub_tematiks,omitempty"`
	Strategics  []StrategicResponse  `json:"strategics,omitempty"`
}

type SubtematikResponse struct {
	Id             int                     `json:"id"`
	Parent         int                     `json:"parent"`
	Tema           string                  `json:"tema_sub_tematik"`
	Keterangan     string                  `json:"keterangan"`
	Indikators     []IndikatorResponse     `json:"indikators"`
	SubSubTematiks []SubSubTematikResponse `json:"sub_sub_tematiks,omitempty"`
	Strategics     []StrategicResponse     `json:"strategics,omitempty"`
}

type SubSubTematikResponse struct {
	Id               int                       `json:"id"`
	Parent           int                       `json:"parent"`
	Tema             string                    `json:"tema_sub_sub_tematik"`
	Keterangan       string                    `json:"keterangan"`
	Indikators       []IndikatorResponse       `json:"indikators"`
	SuperSubTematiks []SuperSubTematikResponse `json:"super_sub_tematiks,omitempty"`
	Strategics       []StrategicResponse       `json:"strategics,omitempty"`
}

type SuperSubTematikResponse struct {
	Id         int                 `json:"id"`
	Parent     int                 `json:"parent"`
	Tema       string              `json:"tema_super_sub_tematik"`
	Keterangan string              `json:"keterangan"`
	Indikators []IndikatorResponse `json:"indikators"`
	Strategics []StrategicResponse `json:"strategics,omitempty"`
}

type IndikatorResponse struct {
	Id            string           `json:"id_indikator"`
	IdPokin       string           `json:"id_pokin"`
	NamaIndikator string           `json:"nama_indikator"`
	Target        []TargetResponse `json:"targets"`
}

type TargetResponse struct {
	Id              string `json:"id_target"`
	IndikatorId     string `json:"indikator_id"`
	TargetIndikator string `json:"target"`
	SatuanIndikator string `json:"satuan"`
}

type StrategicResponse struct {
	Id         int                         `json:"id"`
	Parent     int                         `json:"parent"`
	Strategi   string                      `json:"strategi"`
	Keterangan string                      `json:"keterangan"`
	KodeOpd    opdmaster.OpdResponseForAll `json:"perangkat_daerah"`
	Indikators []IndikatorResponse         `json:"indikators"`
	Tacticals  []TacticalResponse          `json:"tacticals,omitempty"`
}

type TacticalResponse struct {
	Id           int                         `json:"id"`
	Parent       int                         `json:"parent"`
	Strategi     string                      `json:"strategi"`
	Keterangan   *string                     `json:"keterangan"`
	KodeOpd      opdmaster.OpdResponseForAll `json:"kode_perangkat_daerah"`
	Indikators   []IndikatorResponse         `json:"indikators"`
	Operationals []OperationalResponse       `json:"operationals"`
}

type OperationalResponse struct {
	Id         int                         `json:"id"`
	Parent     int                         `json:"parent"`
	Strategi   string                      `json:"strategi"`
	Keterangan *string                     `json:"keterangan"`
	KodeOpd    opdmaster.OpdResponseForAll `json:"kode_perangkat_daerah"`
	Indikators []IndikatorResponse         `json:"indikators"`
}
