package app

import (
	"ekak_kabupaten_madiun/controller"

	"github.com/julienschmidt/httprouter"
)

type RouteController struct {
}

func NewRouter(
	rencanaKinerjaController controller.RencanaKinerjaController,
	rencanaAksiController controller.RencanaAksiController,
	pelaksanaanRencanaAksiController controller.PelaksanaanRencanaAksiController,
	usulanMusrebangController controller.UsulanMusrebangController,
	usulanMandatoriController controller.UsulanMandatoriController,
	usulanPokokPikiranController controller.UsulanPokokPikiranController,
	usulanInisiatifController controller.UsulanInisiatifController,
	usulanTerpilihController controller.UsulanTerpilihController,
	gambaranUmumController controller.GambaranUmumController,
	dasarHukumController controller.DasarHukumController,
	inovasiController controller.InovasiController,
	subKegiatanController controller.SubKegiatanController,
	subKegiatanTerpilihController controller.SubKegiatanTerpilihController,
	pohonKinerjaOpdController controller.PohonKinerjaOpdController,
	pegawaiController controller.PegawaiController,
	lembagaController controller.LembagaController,
	jabatanController controller.JabatanController,
	pohonKinerjaAdminController controller.PohonKinerjaAdminController,
	opdController controller.OpdController,
	programController controller.ProgramController,
	urusanController controller.UrusanController,
	bidangUrusanController controller.BidangUrusanController,
	kegiatanController controller.KegiatanController,
) *httprouter.Router {
	router := httprouter.New()

	//rencana_kinerja
	router.POST("/rencana_kinerja/create", rencanaKinerjaController.Create)
	router.GET("/get_rencana_kinerja/pegawai/:pegawai_id", rencanaKinerjaController.FindAllRencanaKinerja)
	router.GET("/detail-rencana_kinerja/:rencana_kinerja_id", rencanaKinerjaController.FindById)
	router.PUT("/rencana_kinerja/update/:id", rencanaKinerjaController.Update)
	router.DELETE("/rencana_kinerja/delete/:id", rencanaKinerjaController.Delete)

	//rencana_aksi
	router.GET("/rencana_aksi/findall/:rencana_kinerja_id", rencanaAksiController.FindAll)
	// router.GET("/rencana_kinerja/:rekin_id/rincian_kak", rencanaAksiController.FindAll)
	router.GET("/detail-rencana_aksi/:rencanaaksiId", rencanaAksiController.FindById)
	router.POST("/rencana_aksi/create/rencanaaksi/:rekin_id", rencanaAksiController.Create)
	router.PUT("/rencana_aksi/update/rencanaaksi/:rencanaaksiId", rencanaAksiController.Update)
	router.DELETE("/rencana_aksi/delete/rencanaaksi/:rencanaaksiId", rencanaAksiController.Delete)

	//pelaksanaan_rencana_aksi
	router.POST("/pelaksanaan_rencana_aksi/create/:rencanaAksiId", pelaksanaanRencanaAksiController.Create)
	router.PUT("/pelaksanaan_rencana_aksi/update/:pelaksanaanRencanaAksiId", pelaksanaanRencanaAksiController.Update)
	router.GET("/pelaksanaan_rencana_aksi/detail/:id", pelaksanaanRencanaAksiController.FindById)
	router.DELETE("/pelaksanaan_rencanxa_aksi/delete/:id", pelaksanaanRencanaAksiController.Delete)

	//usulan musrebang
	router.POST("/usulan_musrebang/create", usulanMusrebangController.Create)
	router.POST("/usulan_musrebang/create/:pegawai_id", usulanMusrebangController.Create)
	router.PUT("/usulan_musrebang/update/:id", usulanMusrebangController.Update)
	router.PUT("/usulan_musrebang/update/:id/:pegawai_id", usulanMusrebangController.Update)
	router.GET("/usulan_musrebang/detail/:id", usulanMusrebangController.FindById)
	router.DELETE("/usulan_musrebang/delete/:id", usulanMusrebangController.Delete)
	router.GET("/usulan_musrebang/pilihan", usulanMusrebangController.FindAll)
	router.GET("/usulan_musrebang/findall", usulanMusrebangController.FindAll)

	router.GET("/usulan_musrebang/pegawai/:pegawai_id", usulanMusrebangController.FindAll)

	//usulan mandatori
	router.POST("/usulan_mandatori/create", usulanMandatoriController.Create)
	router.POST("/usulan_mandatori/create/:pegawai_id", usulanMandatoriController.Create)
	router.PUT("/usulan_mandatori/update/:id", usulanMandatoriController.Update)
	router.PUT("/usulan_mandatori/update/:id/:pegawai_id", usulanMandatoriController.Update)
	router.GET("/usulan_mandatori/detail/:id", usulanMandatoriController.FindById)
	router.DELETE("/usulan_mandatori/delete/:id", usulanMandatoriController.Delete)
	router.GET("/usulan_mandatori/findall", usulanMandatoriController.FindAll)
	router.GET("/usulan_mandatori/pilihan", usulanMandatoriController.FindAll)
	router.GET("/usulan_mandatori/pegawai/:pegawai_id", usulanMandatoriController.FindAll)

	//usulan pokok pikiran
	router.POST("/usulan_pokok_pikiran/create", usulanPokokPikiranController.Create)
	router.POST("/usulan_pokok_pikiran/create/:pegawai_id", usulanPokokPikiranController.Create)
	router.PUT("/usulan_pokok_pikiran/update/:id", usulanPokokPikiranController.Update)
	router.PUT("/usulan_pokok_pikiran/update/:id/:pegawai_id", usulanPokokPikiranController.Update)
	router.GET("/usulan_pokok_pikiran/detail/:id", usulanPokokPikiranController.FindById)
	router.DELETE("/usulan_pokok_pikiran/delete/:id", usulanPokokPikiranController.Delete)
	router.GET("/usulan_pokok_pikiran/findall", usulanPokokPikiranController.FindAll)
	router.GET("/usulan_pokok_pikiran/pilihan", usulanPokokPikiranController.FindAll)
	router.GET("/usulan_pokok_pikiran/pegawai/:pegawai_id", usulanPokokPikiranController.FindAll)

	//usulan inisiatif
	router.POST("/usulan_inisiatif/create", usulanInisiatifController.Create)
	router.POST("/usulan_inisiatif/create/:pegawai_id", usulanInisiatifController.Create)
	router.PUT("/usulan_inisiatif/update/:id", usulanInisiatifController.Update)
	router.PUT("/usulan_inisiatif/update/:id/:pegawai_id", usulanInisiatifController.Update)
	router.GET("/usulan_inisiatif/detail/:id", usulanInisiatifController.FindById)
	router.DELETE("/usulan_inisiatif/delete/:id", usulanInisiatifController.Delete)
	router.GET("/usulan_inisiatif/findall", usulanInisiatifController.FindAll)
	router.GET("/usulan_inisiatif/pilihan", usulanInisiatifController.FindAll)
	router.GET("/usulan_inisiatif/pegawai/:pegawai_id", usulanInisiatifController.FindAll)

	//usulan terpilih
	router.POST("/usulan_terpilih/create", usulanTerpilihController.Create)
	router.DELETE("/usulan_terpilih/delete/:id_usulan", usulanTerpilihController.Delete)

	//gambaran umum
	router.POST("/gambaran_umum/create/:rencana_kinerja_id", gambaranUmumController.Create)
	router.GET("/gambaran_umum/findall/:rencana_kinerja_id", gambaranUmumController.FindAll)
	router.GET("/gambaran_umum/detail/:id", gambaranUmumController.FindById)
	router.PUT("/gambaran_umum/update/:id", gambaranUmumController.Update)
	router.DELETE("/gambaran_umum/delete/:id", gambaranUmumController.Delete)

	//dasar hukum
	router.POST("/dasar_hukum/create/:rencana_kinerja_id", dasarHukumController.Create)
	router.GET("/dasar_hukum/findall/:rencana_kinerja_id", dasarHukumController.FindAll)
	router.GET("/dasar_hukum/detail/:id", dasarHukumController.FindById)
	router.PUT("/dasar_hukum/update/:id", dasarHukumController.Update)
	router.DELETE("/dasar_hukum/delete/:id", dasarHukumController.Delete)

	//inovasi
	router.POST("/inovasi/create/:rencana_kinerja_id", inovasiController.Create)
	router.GET("/inovasi/findall/:rencana_kinerja_id", inovasiController.FindAll)
	router.GET("/inovasi/detail/:id", inovasiController.FindById)
	router.PUT("/inovasi/update/:id", inovasiController.Update)
	router.DELETE("/inovasi/delete/:id", inovasiController.Delete)

	//sub kegiatan
	router.POST("/sub_kegiatan/create", subKegiatanController.Create)
	router.PUT("/sub_kegiatan/update/:id", subKegiatanController.Update)
	router.GET("/sub_kegiatan/detail/:id", subKegiatanController.FindById)
	router.GET("/sub_kegiatan/findall", subKegiatanController.FindAll)
	router.GET("/sub_kegiatan/pilihan/:kode_opd", subKegiatanController.FindAll)
	router.GET("/sub_kegiatan/byrekinid/:rencana_kinerja_id", subKegiatanController.FindAll)
	router.DELETE("/sub_kegiatan/delete/:id", subKegiatanController.Delete)

	//sub kegiatan terpilih
	router.PUT("/subkegiatanterpilih/create/:rencana_kinerja_id", subKegiatanTerpilihController.Update)
	router.DELETE("/subkegiatanterpilih/delete/:rencana_kinerja_id/:kode_subkegiatan", subKegiatanTerpilihController.Delete)
	router.GET("/subkegiatanterpilih/findbykodesubkegiatan/:kode_subkegiatan", subKegiatanTerpilihController.FindByKodeSubKegiatan)

	//pohon kinerja opd
	router.POST("/pohon_kinerja_opd/create", pohonKinerjaOpdController.Create)
	router.PUT("/pohon_kinerja_opd/update/:id", pohonKinerjaOpdController.Update)
	router.GET("/pohon_kinerja_opd/detail/:id", pohonKinerjaOpdController.FindById)
	router.DELETE("/pohon_kinerja_opd/delete/:id", pohonKinerjaOpdController.Delete)
	router.GET("/pohon_kinerja_opd/findall/:kode_opd/:tahun", pohonKinerjaOpdController.FindAll)
	router.GET("/pohon_kinerja_opd/strategic_no_parent/:kode_opd/:tahun", pohonKinerjaOpdController.FindStrategicNoParent)

	//pohon kinerja admin
	router.POST("/pohon_kinerja_admin/create", pohonKinerjaAdminController.Create)
	router.PUT("/pohon_kinerja_admin/update/:pohonKinerjaId", pohonKinerjaAdminController.Update)
	router.GET("/pohon_kinerja_admin/detail/:id", pohonKinerjaAdminController.FindById)
	router.DELETE("/pohon_kinerja_admin/delete/:pohonKinerjaId", pohonKinerjaAdminController.Delete)
	router.GET("/pohon_kinerja_admin/findall/:tahun", pohonKinerjaAdminController.FindAll)
	router.GET("/pohon_kinerja_admin/subtematik/:tahun", pohonKinerjaAdminController.FindSubTematik)
	router.GET("/pohon_kinerja_admin/tematik/:idPokin", pohonKinerjaAdminController.FindPokinAdminByIdHierarki)

	//DATA MASTER
	//pegawai
	router.POST("/pegawai/create", pegawaiController.Create)
	router.PUT("/pegawai/update/:id", pegawaiController.Update)
	router.GET("/pegawai/detail/:id", pegawaiController.FindById)
	router.DELETE("/pegawai/delete/:id", pegawaiController.Delete)
	router.GET("/pegawai/findall", pegawaiController.FindAll)

	//lembaga
	router.POST("/lembaga/create", lembagaController.Create)
	router.PUT("/lembaga/update/:id", lembagaController.Update)
	router.GET("/lembaga/detail/:id", lembagaController.FindById)
	router.DELETE("/lembaga/delete/:id", lembagaController.Delete)
	router.GET("/lembaga/findall", lembagaController.FindAll)

	//jabatan
	router.POST("/jabatan/create", jabatanController.Create)
	router.PUT("/jabatan/update/:id", jabatanController.Update)
	router.GET("/jabatan/detail/:id", jabatanController.FindById)
	router.DELETE("/jabatan/delete/:id", jabatanController.Delete)
	router.GET("/jabatan/findall/:kode_opd", jabatanController.FindAll)
	router.GET("/jabatan/findall/:kode_opd/:tahun", jabatanController.FindAll)

	//opd
	router.POST("/opd/create", opdController.Create)
	router.PUT("/opd/update/:opdId", opdController.Update)
	router.GET("/opd/detail/:opdId", opdController.FindById)
	router.DELETE("/opd/delete/:opdId", opdController.Delete)
	router.GET("/opd/findall", opdController.FindAll)

	//program
	router.POST("/program_kegiatan/create", programController.Create)
	router.PUT("/program_kegiatan/update/:programId", programController.Update)
	router.GET("/program_kegiatan/detail/:id", programController.FindById)
	router.DELETE("/program_kegiatan/delete/:id", programController.Delete)
	router.GET("/program_kegiatan/findall", programController.FindAll)

	//urusan
	router.POST("/urusan/create", urusanController.Create)
	router.PUT("/urusan/update/:id", urusanController.Update)
	router.GET("/urusan/detail/:id", urusanController.FindById)
	router.DELETE("/urusan/delete/:id", urusanController.Delete)
	router.GET("/urusan/findall", urusanController.FindAll)

	//bidang urusan
	router.POST("/bidang_urusan/create", bidangUrusanController.Create)
	router.PUT("/bidang_urusan/update/:id", bidangUrusanController.Update)
	router.GET("/bidang_urusan/detail/:id", bidangUrusanController.FindById)
	router.DELETE("/bidang_urusan/delete/:id", bidangUrusanController.Delete)
	router.GET("/bidang_urusan/findall", bidangUrusanController.FindAll)

	//kegiatan
	router.POST("/kegiatan/create", kegiatanController.Create)
	router.PUT("/kegiatan/update/:id", kegiatanController.Update)
	router.GET("/kegiatan/detail/:id", kegiatanController.FindById)
	router.DELETE("/kegiatan/delete/:id", kegiatanController.Delete)
	router.GET("/kegiatan/findall", kegiatanController.FindAll)

	//rincian kak
	router.GET("/rencana_kinerja/:rencana_kinerja_id/pegawai/:pegawai_id/input_rincian_kak", rencanaKinerjaController.FindAllRincianKak)

	// router.GET("/rencana_kinerja/:rencana_kinerja_id/pegawai/:pegawai_id/input_rincian_kak", combineHandlers(
	// 	rencanaKinerjaController.FindAll,
	// 	rencanaAksiController.FindAllByRekin,
	// 	usulanMusrebangController.FindAllRekin,
	// 	usulanMandatoriController.FindAllByRekin,
	// 	usulanPokokPikiranController.FindAllByRekin,
	// 	usulanInisiatifController.FindAllByRekin,
	// 	subKegiatanController.FindAllByRekin,
	// 	dasarHukumController.FindAllByRekinId,
	// 	gambaranUmumController.FindAllByRekinId,
	// 	inovasiController.FindAllByRekinId,
	// ))

	return router
}

// Buat fungsi wrapper
// func combineHandlers(handlers ...httprouter.Handle) httprouter.Handle {
// 	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 		for _, handler := range handlers {
// 			handler(w, r, ps)
// 		}
// 	}
// }
