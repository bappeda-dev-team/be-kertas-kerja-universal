package opdmaster

import "ekak_kabupaten_madiun/model/web/lembaga"

type OpdResponse struct {
	Id            string                  `json:"id"`
	KodeOpd       string                  `json:"kode_opd"`
	NamaOpd       string                  `json:"nama_opd"`
	Singkatan     string                  `json:"singkatan"`
	NamaAdmin     string                  `json:"nama_admin_opd"`
	NoWaAdmin     string                  `json:"no_wa_admin_opd"`
	Alamat        string                  `json:"alamat"`
	Telepon       string                  `json:"telepon"`
	Fax           string                  `json:"fax"`
	Email         string                  `json:"email"`
	Website       string                  `json:"website"`
	NamaKepalaOpd string                  `json:"nama_kepala_opd"`
	NIPKepalaOpd  string                  `json:"nip_kepala_opd"`
	PangkatKepala string                  `json:"pangkat_kepala"`
	IdLembaga     lembaga.LembagaResponse `json:"id_lembaga"`
}

type OpdWithBidangUrusan struct {
	Id                string `json:"id"`
	KodeOpd           string `json:"kode_opd"`
	NamaOpd           string `json:"nama_opd"`
	KodeBidangUrusan1 string `json:"kode_bidang_urusan1"`
	NamaBidangUrusan1 string `json:"nama_bidang_urusan1"`
	KodeBidangUrusan2 string `json:"kode_bidang_urusan2"`
	NamaBidangUrusan2 string `json:"nama_bidang_urusan2"`
	KodeBidangUrusan3 string `json:"kode_bidang_urusan3"`
	NamaBidangUrusan3 string `json:"nama_bidang_urusan3"`
	NamaAdmin         string `json:"nama_admin_opd"`
	NoWaAdmin         string `json:"no_wa_admin_opd"`
	NamaKepalaOpd     string `json:"nama_kepala_opd"`
	NIPKepalaOpd      string `json:"nip_kepala_opd"`
	PangkatKepala     string `json:"pangkat_kepala"`
}

type OpdResponseForAll struct {
	KodeOpd string `json:"kode_opd,omitempty"`
	NamaOpd string `json:"nama_opd,omitempty"`
}
