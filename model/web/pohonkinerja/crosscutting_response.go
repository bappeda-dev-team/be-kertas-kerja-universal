package pohonkinerja

import "time"

type CrosscuttingOpdResponse struct {
	Id            int                 `json:"id"`
	NamaPohon     string              `json:"nama_pohon"`
	JenisPohon    string              `json:"jenis_pohon"`
	LevelPohon    int                 `json:"level_pohon"`
	KodeOpd       string              `json:"kode_opd"`
	Keterangan    string              `json:"keterangan"`
	Tahun         string              `json:"tahun"`
	Status        string              `json:"status"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	PegawaiAction interface{}         `json:"pegawai_action,omitempty"`
	Indikator     []IndikatorResponse `json:"indikator"`
}

type CrosscuttingApproveRequest struct {
	Approve    bool   `json:"approve"`
	ParentId   int    `json:"parent_id"`
	NipPegawai string `json:"nip_pegawai"`
	LevelPohon int    `json:"level_pohon"`
	JenisPohon string `json:"jenis_pohon"`
}

type CrosscuttingApproveResponse struct {
	Id         int        `json:"id"`
	Status     string     `json:"status"`
	ApprovedBy *string    `json:"approved_by,omitempty"`
	ApprovedAt *time.Time `json:"approved_at,omitempty"`
	RejectedBy *string    `json:"rejected_by,omitempty"`
	RejectedAt *time.Time `json:"rejected_at,omitempty"`
	Message    string     `json:"message"`
}
