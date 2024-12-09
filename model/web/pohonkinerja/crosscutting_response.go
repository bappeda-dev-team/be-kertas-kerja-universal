package pohonkinerja

import "time"

type CrosscuttingOpdResponse struct {
	Id         int                 `json:"id"`
	NamaPohon  string              `json:"nama_pohon"`
	JenisPohon string              `json:"jenis_pohon"`
	LevelPohon int                 `json:"level_pohon"`
	KodeOpd    string              `json:"kode_opd"`
	Keterangan string              `json:"keterangan"`
	Tahun      string              `json:"tahun"`
	Status     string              `json:"status"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
	Indikator  []IndikatorResponse `json:"indikator"`
}
