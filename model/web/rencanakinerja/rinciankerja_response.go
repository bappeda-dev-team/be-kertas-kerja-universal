package rencanakinerja

import (
	"ekak_kabupaten_madiun/model/web/dasarhukum"
	"ekak_kabupaten_madiun/model/web/gambaranumum"
	"ekak_kabupaten_madiun/model/web/inovasi"
	"ekak_kabupaten_madiun/model/web/rencanaaksi"
	"ekak_kabupaten_madiun/model/web/subkegiatan"
	"ekak_kabupaten_madiun/model/web/usulan"
)

type DataRincianKerja struct {
	RencanaKinerja     RencanaKinerjaResponse              `json:"rencana_kinerja"`
	RencanaAksi        []rencanaaksi.RencanaAksiResponse   `json:"rencana_aksis"`
	UsulanMusrebang    []usulan.UsulanMusrebangResponse    `json:"usulan_musrebang"`
	UsulanMandatori    []usulan.UsulanMandatoriResponse    `json:"usulan_mandatori"`
	UsulanPokokPikiran []usulan.UsulanPokokPikiranResponse `json:"usulan_pokok_pikiran"`
	UsulanInisiatif    []usulan.UsulanInisiatifResponse    `json:"usulan_inisiatif"`
	SubKegiatan        []subkegiatan.SubKegiatanResponse   `json:"subkegiatan"`
	DasarHukum         []dasarhukum.DasarHukumResponse     `json:"dasar_hukum"`
	GambaranUmum       []gambaranumum.GambaranUmumResponse `json:"gambaran_umum"`
	Inovasi            []inovasi.InovasiResponse           `json:"inovasi"`
}
