package pohonkinerja

import "time"

type CrosscuttingOpdResponse struct {
	Id            int                 `json:"id"`
	NamaPohon     string              `json:"nama_pohon,omitempty"`
	JenisPohon    string              `json:"jenis_pohon,omitempty"`
	LevelPohon    int                 `json:"level_pohon,omitempty"`
	KodeOpd       string              `json:"kode_opd,omitempty"`
	NamaOpd       string              `json:"nama_opd"`
	Keterangan    string              `json:"keterangan"`
	Tahun         string              `json:"tahun"`
	Status        string              `json:"status"`
	CreatedAt     time.Time           `json:"created_at,omitempty"`
	UpdatedAt     time.Time           `json:"updated_at,omitempty"`
	PegawaiAction interface{}         `json:"pegawai_action,omitempty"`
	Indikator     []IndikatorResponse `json:"indikator,omitempty"`
}

type CrosscuttingApproveRequest struct {
	Approve     bool   `json:"approve"`
	CreateNew   bool   `json:"create_new"`
	UseExisting bool   `json:"use_existing"`
	NamaPohon   string `json:"nama_pohon"`
	ParentId    int    `json:"parent_id"`
	NipPegawai  string `json:"nip_pegawai"`
	LevelPohon  int    `json:"level_pohon"`
	JenisPohon  string `json:"jenis_pohon"`
	ExistingId  int    `json:"existing_id,omitempty"`
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

type CrosscuttingFromResponse struct {
	KodeOpd string `json:"kode_opd"`
	NamaOpd string `json:"nama_opd"`
}
