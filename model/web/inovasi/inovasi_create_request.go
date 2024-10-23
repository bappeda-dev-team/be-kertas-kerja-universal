package inovasi

type InovasiCreateRequest struct {
	RekinId               string `json:"rencana_kinerja_id"`
	PegawaiId             string `json:"pegawai_id"`
	JudulInovasi          string `json:"judul_inovasi"`
	JenisInovasi          string `json:"jenis_inovasi"`
	GambaranNilaiKebaruan string `json:"gambaran_nilai_kebaruan"`
}
