package helper

import (
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/dasarhukum"
	"ekak_kabupaten_madiun/model/web/gambaranumum"
	"ekak_kabupaten_madiun/model/web/inovasi"
	"ekak_kabupaten_madiun/model/web/jabatan"
	"ekak_kabupaten_madiun/model/web/opdmaster"
	"ekak_kabupaten_madiun/model/web/pegawai"
	"ekak_kabupaten_madiun/model/web/rencanakinerja"
	"ekak_kabupaten_madiun/model/web/subkegiatan"
	"ekak_kabupaten_madiun/model/web/usulan"
	"fmt"
	"os"
)

func ToRencanaKinerjaResponse(rencanaKinerja domain.RencanaKinerja) rencanakinerja.RencanaKinerjaResponse {
	indikatorResponses := make([]rencanakinerja.IndikatorResponse, 0)
	for _, indikator := range rencanaKinerja.Indikator {
		targetResponses := make([]rencanakinerja.TargetResponse, 0)
		for _, target := range indikator.Target {
			targetResponses = append(targetResponses, rencanakinerja.TargetResponse{
				Id:              target.Id,
				IndikatorId:     target.IndikatorId,
				TargetIndikator: target.Target,
				SatuanIndikator: target.Satuan,
			})
		}
		indikatorResponses = append(indikatorResponses, rencanakinerja.IndikatorResponse{
			Id:            indikator.Id,
			NamaIndikator: indikator.Indikator,
			Target:        targetResponses,
		})
	}

	opdResponse := opdmaster.OpdResponseForAll{
		KodeOpd: rencanaKinerja.KodeOpd,
		NamaOpd: rencanaKinerja.NamaOpd,
	}

	return rencanakinerja.RencanaKinerjaResponse{
		Id:                   rencanaKinerja.Id,
		NamaRencanaKinerja:   rencanaKinerja.NamaRencanaKinerja,
		Tahun:                rencanaKinerja.Tahun,
		StatusRencanaKinerja: rencanaKinerja.StatusRencanaKinerja,
		Catatan:              rencanaKinerja.Catatan,
		KodeOpd:              opdResponse,
		PegawaiId:            rencanaKinerja.PegawaiId,
		Indikator:            indikatorResponses,
	}
}

func ToRencanaKinerjaResponses(rencanaKinerjas []domain.RencanaKinerja) []rencanakinerja.RencanaKinerjaResponse {
	var rencanaKinerjaResponses []rencanakinerja.RencanaKinerjaResponse
	for _, rencanaKinerja := range rencanaKinerjas {
		rencanaKinerjaResponses = append(rencanaKinerjaResponses, ToRencanaKinerjaResponse(rencanaKinerja))
	}
	return rencanaKinerjaResponses
}
func ToUsulanMusrebangResponse(usulanMusrebang domain.UsulanMusrebang) usulan.UsulanMusrebangResponse {
	host := os.Getenv("host")
	port := os.Getenv("port")
	buttonActions := []web.ActionButton{
		{
			NameAction: "Find Id Usulan Musrebang",
			Method:     "GET",
			Url:        fmt.Sprintf("%s:%s/usulan_musrebang/detail/:id", host, port),
		},
		{
			NameAction: "Delete Usulan Musrebang Terpilih",
			Method:     "DELETE",
			Url:        fmt.Sprintf("%s:%s/usulan_terpilih/delete/:id_usulan", host, port),
		},
	}
	return usulan.UsulanMusrebangResponse{
		Id:        usulanMusrebang.Id,
		Usulan:    usulanMusrebang.Usulan,
		Alamat:    usulanMusrebang.Alamat,
		Uraian:    usulanMusrebang.Uraian,
		Tahun:     usulanMusrebang.Tahun,
		RekinId:   usulanMusrebang.RekinId,
		PegawaiId: usulanMusrebang.PegawaiId,
		KodeOpd:   usulanMusrebang.KodeOpd,
		IsActive:  usulanMusrebang.IsActive,
		Status:    usulanMusrebang.Status,
		CreatedAt: usulanMusrebang.CreatedAt.Format("2006-01-02"),
		Action:    buttonActions,
	}
}

func ToUsulanMusrebangResponses(usulanMusrebangs []domain.UsulanMusrebang) []usulan.UsulanMusrebangResponse {
	var usulanMusrebangResponses []usulan.UsulanMusrebangResponse
	for _, usulanMusrebang := range usulanMusrebangs {
		usulanMusrebangResponses = append(usulanMusrebangResponses, ToUsulanMusrebangResponse(usulanMusrebang))
	}
	return usulanMusrebangResponses
}

func ToUsulanMandatoriResponse(usulanMandatori domain.UsulanMandatori) usulan.UsulanMandatoriResponse {
	host := os.Getenv("host")
	port := os.Getenv("port")
	buttonActions := []web.ActionButton{
		{
			NameAction: "Find Id Usulan Mandatori",
			Method:     "GET",
			Url:        fmt.Sprintf("%s:%s/usulan_mandatori/detail/:id", host, port),
		},
		{
			NameAction: "Delete Usulan Mandatori Terpilih",
			Method:     "DELETE",
			Url:        fmt.Sprintf("%s:%s/usulan_terpilih/delete/:id_usulan", host, port),
		},
	}
	return usulan.UsulanMandatoriResponse{
		Id:               usulanMandatori.Id,
		Usulan:           usulanMandatori.Usulan,
		PeraturanTerkait: usulanMandatori.PeraturanTerkait,
		Uraian:           usulanMandatori.Uraian,
		Tahun:            usulanMandatori.Tahun,
		RekinId:          usulanMandatori.RekinId,
		PegawaiId:        usulanMandatori.PegawaiId,
		KodeOpd:          usulanMandatori.KodeOpd,
		Status:           usulanMandatori.Status,
		IsActive:         usulanMandatori.IsActive,
		CreatedAt:        usulanMandatori.CreatedAt.Format("2006-01-02"),
		Action:           buttonActions,
	}
}

func ToUsulanMandatoriResponses(usulanMandatoris []domain.UsulanMandatori) []usulan.UsulanMandatoriResponse {
	var usulanMandatoriResponses []usulan.UsulanMandatoriResponse
	for _, usulanMandatori := range usulanMandatoris {
		usulanMandatoriResponses = append(usulanMandatoriResponses, ToUsulanMandatoriResponse(usulanMandatori))
	}
	return usulanMandatoriResponses
}

func ToUsulanPokokPikiranResponse(usulanPokokPikiran domain.UsulanPokokPikiran) usulan.UsulanPokokPikiranResponse {
	host := os.Getenv("host")
	port := os.Getenv("port")
	buttonActions := []web.ActionButton{
		{
			NameAction: "Find Id Usulan Pokok Pikiran",
			Method:     "GET",
			Url:        fmt.Sprintf("%s:%s/usulan_pokok_pikiran/detail/:id", host, port),
		},
		{
			NameAction: "Delete Usulan Pokok Pikiran Terpilih",
			Method:     "DELETE",
			Url:        fmt.Sprintf("%s:%s/usulan_terpilih/delete/:id_usulan", host, port),
		},
	}
	return usulan.UsulanPokokPikiranResponse{
		Id:        usulanPokokPikiran.Id,
		Usulan:    usulanPokokPikiran.Usulan,
		Alamat:    usulanPokokPikiran.Alamat,
		Uraian:    usulanPokokPikiran.Uraian,
		Tahun:     usulanPokokPikiran.Tahun,
		RekinId:   usulanPokokPikiran.RekinId,
		PegawaiId: usulanPokokPikiran.PegawaiId,
		KodeOpd:   usulanPokokPikiran.KodeOpd,
		Status:    usulanPokokPikiran.Status,
		IsActive:  usulanPokokPikiran.IsActive,
		CreatedAt: usulanPokokPikiran.CreatedAt.Format("2006-01-02"),
		Action:    buttonActions,
	}
}

func ToUsulanPokokPikiranResponses(usulanPokokPikirans []domain.UsulanPokokPikiran) []usulan.UsulanPokokPikiranResponse {
	var usulanPokokPikiranResponses []usulan.UsulanPokokPikiranResponse
	for _, usulanPokokPikiran := range usulanPokokPikirans {
		usulanPokokPikiranResponses = append(usulanPokokPikiranResponses, ToUsulanPokokPikiranResponse(usulanPokokPikiran))
	}
	return usulanPokokPikiranResponses
}

func ToUsulanInisiatifResponse(usulanInisiatif domain.UsulanInisiatif) usulan.UsulanInisiatifResponse {
	host := os.Getenv("host")
	port := os.Getenv("port")
	buttonActions := []web.ActionButton{
		{
			NameAction: "Find Id Usulan Inisiatif",
			Method:     "GET",
			Url:        fmt.Sprintf("%s:%s/usulan_inisiatif/detail/:id", host, port),
		},
		{
			NameAction: "Delete Usulan Inisiatif Terpilih",
			Method:     "DELETE",
			Url:        fmt.Sprintf("%s:%s/usulan_terpilih/delete/:id_usulan", host, port),
		},
	}
	return usulan.UsulanInisiatifResponse{
		Id:        usulanInisiatif.Id,
		Usulan:    usulanInisiatif.Usulan,
		Manfaat:   usulanInisiatif.Manfaat,
		Uraian:    usulanInisiatif.Uraian,
		Tahun:     usulanInisiatif.Tahun,
		RekinId:   usulanInisiatif.RekinId,
		PegawaiId: usulanInisiatif.PegawaiId,
		KodeOpd:   usulanInisiatif.KodeOpd,
		Status:    usulanInisiatif.Status,
		CreatedAt: usulanInisiatif.CreatedAt.Format("2006-01-02"),
		Action:    buttonActions,
	}
}

func ToUsulanInisiatifResponses(usulanInovasis []domain.UsulanInisiatif) []usulan.UsulanInisiatifResponse {
	var usulanInovasiResponses []usulan.UsulanInisiatifResponse
	for _, usulanInovasi := range usulanInovasis {
		usulanInovasiResponses = append(usulanInovasiResponses, ToUsulanInisiatifResponse(usulanInovasi))
	}
	return usulanInovasiResponses
}

func ToUsulanTerpilihResponse(usulanTerpilih domain.UsulanTerpilih) usulan.UsulanTerpilihResponse {
	return usulan.UsulanTerpilihResponse{
		Id:          usulanTerpilih.Id,
		Keterangan:  usulanTerpilih.Keterangan,
		JenisUsulan: usulanTerpilih.JenisUsulan,
		UsulanId:    usulanTerpilih.UsulanId,
		RekinId:     usulanTerpilih.RekinId,
		Tahun:       usulanTerpilih.Tahun,
		KodeOpd:     usulanTerpilih.KodeOpd,
	}
}

func ToUsulanTerpilihResponses(usulanTerpilihs []domain.UsulanTerpilih) []usulan.UsulanTerpilihResponse {
	var usulanTerpilihResponses []usulan.UsulanTerpilihResponse
	for _, usulanTerpilih := range usulanTerpilihs {
		usulanTerpilihResponses = append(usulanTerpilihResponses, ToUsulanTerpilihResponse(usulanTerpilih))
	}
	return usulanTerpilihResponses
}

func ToGambaranUmumResponse(gambaranUmum domain.GambaranUmum) gambaranumum.GambaranUmumResponse {
	host := os.Getenv("host")
	port := os.Getenv("port")
	buttonActions := []web.ActionButton{
		{
			NameAction: "Find By Id Gambaran Umum",
			Method:     "GET",
			Url:        fmt.Sprintf("%s:%s/gambaran_umum/detail/:id", host, port),
		},
		{
			NameAction: "Update Gambaran Umum",
			Method:     "PUT",
			Url:        fmt.Sprintf("%s:%s/gambaran_umum/update/:id", host, port),
		},
		{
			NameAction: "Delete Gambaran Umum",
			Method:     "DELETE",
			Url:        fmt.Sprintf("%s:%s/gambaran_umum/delete/:id", host, port),
		},
	}
	return gambaranumum.GambaranUmumResponse{
		Id:           gambaranUmum.Id,
		RekinId:      gambaranUmum.RekinId,
		PegawaiId:    gambaranUmum.PegawaiId,
		KodeOpd:      gambaranUmum.KodeOpd,
		Urutan:       gambaranUmum.Urutan,
		GambaranUmum: gambaranUmum.GambaranUmum,
		Action:       buttonActions,
	}
}

func ToGambaranUmumResponses(gambaranUmums []domain.GambaranUmum) []gambaranumum.GambaranUmumResponse {
	var gambaranUmumResponses []gambaranumum.GambaranUmumResponse
	for _, gambaranUmum := range gambaranUmums {
		gambaranUmumResponses = append(gambaranUmumResponses, ToGambaranUmumResponse(gambaranUmum))
	}
	return gambaranUmumResponses
}

func ToDasarHukumResponse(dasarHukum domain.DasarHukum) dasarhukum.DasarHukumResponse {
	host := os.Getenv("host")
	port := os.Getenv("port")
	buttonActions := []web.ActionButton{
		{
			NameAction: "Find By Id Dasar Hukum",
			Method:     "GET",
			Url:        fmt.Sprintf("%s:%s/dasar_hukum/detail/:id", host, port),
		},
		{
			NameAction: "Update Dasar Hukum",
			Method:     "PUT",
			Url:        fmt.Sprintf("%s:%s/dasar_hukum/update/:id", host, port),
		},
		{
			NameAction: "Delete Dasar Hukum",
			Method:     "DELETE",
			Url:        fmt.Sprintf("%s:%s/dasar_hukum/delete/:id", host, port),
		},
	}
	return dasarhukum.DasarHukumResponse{
		Id:               dasarHukum.Id,
		RekinId:          dasarHukum.RekinId,
		PegawaiId:        dasarHukum.PegawaiId,
		KodeOpd:          dasarHukum.KodeOpd,
		Urutan:           dasarHukum.Urutan,
		PeraturanTerkait: dasarHukum.PeraturanTerkait,
		Uraian:           dasarHukum.Uraian,
		Action:           buttonActions,
	}
}

func ToDasarHukumResponses(dasarHukums []domain.DasarHukum) []dasarhukum.DasarHukumResponse {
	var dasarHukumResponses []dasarhukum.DasarHukumResponse
	for _, dasarHukum := range dasarHukums {
		dasarHukumResponses = append(dasarHukumResponses, ToDasarHukumResponse(dasarHukum))
	}
	return dasarHukumResponses
}

func ToInovasiResponse(Inovasi domain.Inovasi) inovasi.InovasiResponse {
	host := os.Getenv("host")
	port := os.Getenv("port")
	buttonActions := []web.ActionButton{
		{
			NameAction: "Find By Id Inovasi",
			Method:     "GET",
			Url:        fmt.Sprintf("%s:%s/inovasi/detail/:id", host, port),
		},
		{
			NameAction: "Update Inovasi",
			Method:     "PUT",
			Url:        fmt.Sprintf("%s:%s/inovasi/update/:id", host, port),
		},
		{
			NameAction: "Delete Inovasi",
			Method:     "DELETE",
			Url:        fmt.Sprintf("%s:%s/inovasi/delete/:id", host, port),
		},
	}
	return inovasi.InovasiResponse{
		Id:                    Inovasi.Id,
		RekinId:               Inovasi.RekinId,
		PegawaiId:             Inovasi.PegawaiId,
		KodeOpd:               Inovasi.KodeOpd,
		JudulInovasi:          Inovasi.JudulInovasi,
		JenisInovasi:          Inovasi.JenisInovasi,
		GambaranNilaiKebaruan: Inovasi.GambaranNilaiKebaruan,
		Action:                buttonActions,
	}
}

func ToInovasiResponses(Inovasis []domain.Inovasi) []inovasi.InovasiResponse {
	var inovasiResponses []inovasi.InovasiResponse
	for _, inovasi := range Inovasis {
		inovasiResponses = append(inovasiResponses, ToInovasiResponse(inovasi))
	}
	return inovasiResponses
}

func ToSubKegiatanResponse(subKegiatan domain.SubKegiatan) subkegiatan.SubKegiatanResponse {
	indikatorResponses := make([]subkegiatan.IndikatorSubKegiatanResponse, 0)
	for _, indikator := range subKegiatan.IndikatorSubKegiatan {
		indikatorResponses = append(indikatorResponses, subkegiatan.IndikatorSubKegiatanResponse{
			Id:            indikator.Id,
			SubKegiatanId: indikator.SubKegiatanId,
			NamaIndikator: indikator.NamaIndikator,
		})
	}

	paguResponses := make([]subkegiatan.PaguSubKegiatanResponse, 0)
	for _, pagu := range subKegiatan.PaguSubKegiatan {
		paguResponses = append(paguResponses, subkegiatan.PaguSubKegiatanResponse{
			Id:            pagu.Id,
			SubKegiatanId: pagu.SubKegiatanId,
			JenisPagu:     pagu.JenisPagu,
			PaguAnggaran:  pagu.PaguAnggaran,
			Tahun:         pagu.Tahun,
		})
	}

	host := os.Getenv("host")
	port := os.Getenv("port")
	buttonActions := []web.ActionButton{
		{
			NameAction: "Find By IdSubkegiatan",
			Method:     "GET",
			Url:        fmt.Sprintf("%s:%s/sub_kegiatan/detail/:idsubkegiatan", host, port),
		},
		{
			NameAction: "Delete Subkegiatan",
			Method:     "DELETE",
			Url:        fmt.Sprintf("%s:%s/subkegiatanterpilih/delete/:subkegiatan_id", host, port),
		},
	}

	return subkegiatan.SubKegiatanResponse{
		Id:                   subKegiatan.Id,
		PegawaiId:            subKegiatan.PegawaiId,
		NamaSubKegiatan:      subKegiatan.NamaSubKegiatan,
		KodeOpd:              subKegiatan.KodeOpd,
		Tahun:                subKegiatan.Tahun,
		IndikatorSubkegiatan: indikatorResponses,
		PaguSubKegiatan:      paguResponses,
		Action:               buttonActions,
	}
}

func ToSubKegiatanResponses(subKegiatans []domain.SubKegiatan) []subkegiatan.SubKegiatanResponse {
	var subKegiatanResponses []subkegiatan.SubKegiatanResponse
	for _, subKegiatan := range subKegiatans {
		subKegiatanResponses = append(subKegiatanResponses, ToSubKegiatanResponse(subKegiatan))
	}
	return subKegiatanResponses
}

func ToSubKegiatanTerpilihResponse(subKegiatanTerpilih domain.SubKegiatanTerpilih) subkegiatan.SubKegiatanTerpilihResponse {
	return subkegiatan.SubKegiatanTerpilihResponse{
		Id:               subKegiatanTerpilih.Id,
		RencanaKinerjaId: subKegiatanTerpilih.RencanaKinerjaId,
		SubKegiatanId:    subKegiatanTerpilih.SubKegiatanId,
	}
}

func ToSubKegiatanTerpilihResponses(subKegiatanTerpilihs []domain.SubKegiatanTerpilih) []subkegiatan.SubKegiatanTerpilihResponse {
	var subKegiatanTerpilihResponses []subkegiatan.SubKegiatanTerpilihResponse
	for _, subKegiatanTerpilih := range subKegiatanTerpilihs {
		subKegiatanTerpilihResponses = append(subKegiatanTerpilihResponses, ToSubKegiatanTerpilihResponse(subKegiatanTerpilih))
	}
	return subKegiatanTerpilihResponses
}

func ToPegawaiResponse(pegawais domainmaster.Pegawai) pegawai.PegawaiResponse {
	return pegawai.PegawaiResponse{ // Langsung return struct PegawaiResponse
		Id:          pegawais.Id,
		NamaPegawai: pegawais.NamaPegawai,
		Nip:         pegawais.Nip,
	}
}
func ToPegawaiResponses(pegawais []domainmaster.Pegawai) []pegawai.PegawaiResponse {
	var pegawaiResponses []pegawai.PegawaiResponse
	for _, pegawai := range pegawais {
		pegawaiResponses = append(pegawaiResponses, ToPegawaiResponse(pegawai))
	}
	return pegawaiResponses
}

func ToJabatanResponse(jabatans domainmaster.Jabatan) jabatan.JabatanResponse {
	opd := opdmaster.OpdResponseForAll{
		KodeOpd: jabatans.KodeOpd,
		NamaOpd: jabatans.NamaOpd,
	}
	return jabatan.JabatanResponse{
		Id:           jabatans.Id,
		KodeJabatan:  jabatans.KodeJabatan,
		NamaJabatan:  jabatans.NamaJabatan,
		KelasJabatan: jabatans.KelasJabatan,
		JenisJabatan: jabatans.JenisJabatan,
		NilaiJabatan: jabatans.NilaiJabatan,
		KodeOpd:      opd,
		IndexJabatan: jabatans.IndexJabatan,
		Tahun:        jabatans.Tahun,
		Esselon:      jabatans.Esselon,
	}
}

func ToJabatanResponses(jabatans []domainmaster.Jabatan) []jabatan.JabatanResponse {
	var jabatanResponses []jabatan.JabatanResponse
	for _, jabatan := range jabatans {
		jabatanResponses = append(jabatanResponses, ToJabatanResponse(jabatan))
	}
	return jabatanResponses
}
