package rencanakinerja

type ManualIKResponse struct {
	Id                  int           `json:"id"`
	IndikatorId         string        `json:"indikator_id"`
	DataIndikator       RekinResponse `json:"rencana_kinerja"`
	Perspektif          string        `json:"perspektif"`
	TujuanRekin         string        `json:"tujuan_rekin"`
	Definisi            string        `json:"definisi"`
	KeyActivities       string        `json:"key_activities"`
	Formula             string        `json:"formula"`
	JenisIndikator      string        `json:"jenis_indikator"`
	OutputData          OutputData    `json:"output_data"`
	UnitPenanggungJawab string        `json:"unit_penanggung_jawab"`
	UnitPenyediaData    string        `json:"unit_penyedia_data"`
	SumberData          string        `json:"sumber_data"`
	JangkaWaktuAwal     string        `json:"jangka_waktu_awal"`
	JangkaWaktuAkhir    string        `json:"jangka_waktu_akhir"`
	PeriodePelaporan    string        `json:"periode_pelaporan"`
}

type OutputData struct {
	Kinerja  bool `json:"kinerja"`
	Penduduk bool `json:"penduduk"`
	Spatial  bool `json:"spatial"`
}

type RekinResponse struct {
	RencanaKinerja string              `json:"Nama_rencana_kinerja,omitempty"`
	Indikator      []IndikatorResponse `json:"indikator,omitempty"`
}
