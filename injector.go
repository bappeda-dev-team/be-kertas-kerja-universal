//go:build wireinject
// +build wireinject

package main

import (
	"ekak_kabupaten_madiun/app"
	"ekak_kabupaten_madiun/controller"
	"ekak_kabupaten_madiun/dataseeder"
	"ekak_kabupaten_madiun/middleware"
	"ekak_kabupaten_madiun/repository"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/julienschmidt/httprouter"
)

var rencanaKinerjaSet = wire.NewSet(
	repository.NewRencanaKinerjaRepositoryImpl,
	wire.Bind(new(repository.RencanaKinerjaRepository), new(*repository.RencanaKinerjaRepositoryImpl)),
	service.NewRencanaKinerjaServiceImpl,
	wire.Bind(new(service.RencanaKinerjaService), new(*service.RencanaKinerjaServiceImpl)),
	controller.NewRencanaKinerjaControllerImpl,
	wire.Bind(new(controller.RencanaKinerjaController), new(*controller.RencanaKinerjaControllerImpl)),
)

var rencanaAksiSet = wire.NewSet(
	repository.NewRencanaAksiRepositoryImpl,
	wire.Bind(new(repository.RencanaAksiRepository), new(*repository.RencanaAksiRepositoryImpl)),
	service.NewRencanaAksiServiceImpl,
	wire.Bind(new(service.RencanaAksiService), new(*service.RencanaAksiServiceImpl)),
	controller.NewRencanaAksiControllerImpl,
	wire.Bind(new(controller.RencanaAksiController), new(*controller.RencanaAksiControllerImpl)),
)

var pelaksanaanRencanaAksiSet = wire.NewSet(
	repository.NewPelaksanaanRencanaAksiRepositoryImpl,
	wire.Bind(new(repository.PelaksanaanRencanaAksiRepository), new(*repository.PelaksanaanRencanaAksiRepositoryImpl)),
	service.NewPelaksanaanRencanaAksiServiceImpl,
	wire.Bind(new(service.PelaksanaanRencanaAksiService), new(*service.PelaksanaanRencanaAksiServiceImpl)),
	controller.NewPelaksanaanRencanaAksiControllerImpl,
	wire.Bind(new(controller.PelaksanaanRencanaAksiController), new(*controller.PelaksanaanRencanaAksiControllerImpl)),
)

var usulanMusrebangSet = wire.NewSet(
	repository.NewUsulanMusrebangRepositoryImpl,
	wire.Bind(new(repository.UsulanMusrebangRepository), new(*repository.UsulanMusrebangRepositoryImpl)),
	service.NewUsulanMusrebangServiceImpl,
	wire.Bind(new(service.UsulanMusrebangService), new(*service.UsulanMusrebangServiceImpl)),
	controller.NewUsulanMusrebangControllerImpl,
	wire.Bind(new(controller.UsulanMusrebangController), new(*controller.UsulanMusrebangControllerImpl)),
)

var usulanMandatoriSet = wire.NewSet(
	repository.NewUsulanMandatoriRepositoryImpl,
	wire.Bind(new(repository.UsulanMandatoriRepository), new(*repository.UsulanMandatoriRepositoryImpl)),
	service.NewUsulanMandatoriServiceImpl,
	wire.Bind(new(service.UsulanMandatoriService), new(*service.UsulanMandatoriServiceImpl)),
	controller.NewUsulanMandatoriControllerImpl,
	wire.Bind(new(controller.UsulanMandatoriController), new(*controller.UsulanMandatoriControllerImpl)),
)

var usulanPokokPikiranSet = wire.NewSet(
	repository.NewUsulanPokokPikiranRepositoryImpl,
	wire.Bind(new(repository.UsulanPokokPikiranRepository), new(*repository.UsulanPokokPikiranRepositoryImpl)),
	service.NewUsulanPokokPikiranServiceImpl,
	wire.Bind(new(service.UsulanPokokPikiranService), new(*service.UsulanPokokPikiranServiceImpl)),
	controller.NewUsulanPokokPikiranControllerImpl,
	wire.Bind(new(controller.UsulanPokokPikiranController), new(*controller.UsulanPokokPikiranControllerImpl)),
)

var usulanInisiatifSet = wire.NewSet(
	repository.NewUsulanInisiatifRepositoryImpl,
	wire.Bind(new(repository.UsulanInisiatifRepository), new(*repository.UsulanInisiatifRepositoryImpl)),
	service.NewUsulanInisiatifServiceImpl,
	wire.Bind(new(service.UsulanInisiatifService), new(*service.UsulanInisiatifServiceImpl)),
	controller.NewUsulanInisiatifControllerImpl,
	wire.Bind(new(controller.UsulanInisiatifController), new(*controller.UsulanInisiatifControllerImpl)),
)

var usulanTerpilihSet = wire.NewSet(
	repository.NewUsulanTerpilihRepositoryImpl,
	wire.Bind(new(repository.UsulanTerpilihRepository), new(*repository.UsulanTerpilihRepositoryImpl)),
	service.NewUsulanTerpilihServiceImpl,
	wire.Bind(new(service.UsulanTerpilihService), new(*service.UsulanTerpilihServiceImpl)),
	controller.NewUsulanTerpilihControllerImpl,
	wire.Bind(new(controller.UsulanTerpilihController), new(*controller.UsulanTerpilihControllerImpl)),
)

var gambaranUmumSet = wire.NewSet(
	repository.NewGambaranUmumRepositoryImpl,
	wire.Bind(new(repository.GambaranUmumRepository), new(*repository.GambaranUmumRepositoryImpl)),
	service.NewGambaranUmumServiceImpl,
	wire.Bind(new(service.GambaranUmumService), new(*service.GambaranUmumServiceImpl)),
	controller.NewGambaranUmumControllerImpl,
	wire.Bind(new(controller.GambaranUmumController), new(*controller.GambaranUmumControllerImpl)),
)

var dasarHukumSet = wire.NewSet(
	repository.NewDasarHukumRepositoryImpl,
	wire.Bind(new(repository.DasarHukumRepository), new(*repository.DasarHukumRepositoryImpl)),
	service.NewDasarHukumServiceImpl,
	wire.Bind(new(service.DasarHukumService), new(*service.DasarHukumServiceImpl)),
	controller.NewDasarHukumControllerImpl,
	wire.Bind(new(controller.DasarHukumController), new(*controller.DasarHukumControllerImpl)),
)

var inovasiSet = wire.NewSet(
	repository.NewInovasiRepositoryImpl,
	wire.Bind(new(repository.InovasiRepository), new(*repository.InovasiRepositoryImpl)),
	service.NewInovasiServiceImpl,
	wire.Bind(new(service.InovasiService), new(*service.InovasiServiceImpl)),
	controller.NewInovasiControllerImpl,
	wire.Bind(new(controller.InovasiController), new(*controller.InovasiControllerImpl)),
)

var subKegiatanSet = wire.NewSet(
	repository.NewSubKegiatanRepositoryImpl,
	wire.Bind(new(repository.SubKegiatanRepository), new(*repository.SubKegiatanRepositoryImpl)),
	service.NewSubKegiatanServiceImpl,
	wire.Bind(new(service.SubKegiatanService), new(*service.SubKegiatanServiceImpl)),
	controller.NewSubKegiatanControllerImpl,
	wire.Bind(new(controller.SubKegiatanController), new(*controller.SubKegiatanControllerImpl)),
)

var subKegiatanTerpilihSet = wire.NewSet(
	repository.NewSubKegiatanTerpilihRepositoryImpl,
	wire.Bind(new(repository.SubKegiatanTerpilihRepository), new(*repository.SubKegiatanTerpilihRepositoryImpl)),
	service.NewSubKegiatanTerpilihServiceImpl,
	wire.Bind(new(service.SubKegiatanTerpilihService), new(*service.SubKegiatanTerpilihServiceImpl)),
	controller.NewSubKegiatanTerpilihControllerImpl,
	wire.Bind(new(controller.SubKegiatanTerpilihController), new(*controller.SubKegiatanTerpilihControllerImpl)),
)

var pohonKinerjaOpdSet = wire.NewSet(
	repository.NewPohonKinerjaRepositoryImpl,
	wire.Bind(new(repository.PohonKinerjaRepository), new(*repository.PohonKinerjaRepositoryImpl)),
	service.NewPohonKinerjaOpdServiceImpl,
	wire.Bind(new(service.PohonKinerjaOpdService), new(*service.PohonKinerjaOpdServiceImpl)),
	controller.NewPohonKinerjaOpdControllerImpl,
	wire.Bind(new(controller.PohonKinerjaOpdController), new(*controller.PohonKinerjaOpdControllerImpl)),
)

var pegawaiSet = wire.NewSet(
	repository.NewPegawaiRepositoryImpl,
	wire.Bind(new(repository.PegawaiRepository), new(*repository.PegawaiRepositoryImpl)),
	service.NewPegawaiServiceImpl,
	wire.Bind(new(service.PegawaiService), new(*service.PegawaiServiceImpl)),
	controller.NewPegawaiControllerImpl,
	wire.Bind(new(controller.PegawaiController), new(*controller.PegawaiControllerImpl)),
)

var lembagaSet = wire.NewSet(
	repository.NewLembagaRepositoryImpl,
	wire.Bind(new(repository.LembagaRepository), new(*repository.LembagaRepositoryImpl)),
	service.NewLembagaServiceImpl,
	wire.Bind(new(service.LembagaService), new(*service.LembagaServiceImpl)),
	controller.NewLembagaControllerImpl,
	wire.Bind(new(controller.LembagaController), new(*controller.LembagaControllerImpl)),
)

var jabatanSet = wire.NewSet(
	repository.NewJabatanRepositoryImpl,
	wire.Bind(new(repository.JabatanRepository), new(*repository.JabatanRepositoryImpl)),
	service.NewJabatanServiceImpl,
	wire.Bind(new(service.JabatanService), new(*service.JabatanServiceImpl)),
	controller.NewJabatanControllerImpl,
	wire.Bind(new(controller.JabatanController), new(*controller.JabatanControllerImpl)),
)

var pohonKinerjaAdminSet = wire.NewSet(

	service.NewPohonKinerjaAdminServiceImpl,
	wire.Bind(new(service.PohonKinerjaAdminService), new(*service.PohonKinerjaAdminServiceImpl)),
	controller.NewPohonKinerjaAdminControllerImpl,
	wire.Bind(new(controller.PohonKinerjaAdminController), new(*controller.PohonKinerjaAdminControllerImpl)),
)

var opdSet = wire.NewSet(
	repository.NewOpdRepositoryImpl,
	wire.Bind(new(repository.OpdRepository), new(*repository.OpdRepositoryImpl)),
	service.NewOpdServiceImpl,
	wire.Bind(new(service.OpdService), new(*service.OpdServiceImpl)),
	controller.NewOpdControllerImpl,
	wire.Bind(new(controller.OpdController), new(*controller.OpdControllerImpl)),
)

var programSet = wire.NewSet(
	repository.NewProgramRepositoryImpl,
	wire.Bind(new(repository.ProgramRepository), new(*repository.ProgramRepositoryImpl)),
	service.NewProgramServiceImpl,
	wire.Bind(new(service.ProgramService), new(*service.ProgramServiceImpl)),
	controller.NewProgramControllerImpl,
	wire.Bind(new(controller.ProgramController), new(*controller.ProgramControllerImpl)),
)

var urusanSet = wire.NewSet(
	repository.NewUrusanRepositoryImpl,
	wire.Bind(new(repository.UrusanRepository), new(*repository.UrusanRepositoryImpl)),
	service.NewUrusanServiceImpl,
	wire.Bind(new(service.UrusanService), new(*service.UrusanServiceImpl)),
	controller.NewUrusanControllerImpl,
	wire.Bind(new(controller.UrusanController), new(*controller.UrusanControllerImpl)),
)

var bidangUrusanSet = wire.NewSet(
	repository.NewBidangUrusanRepositoryImpl,
	wire.Bind(new(repository.BidangUrusanRepository), new(*repository.BidangUrusanRepositoryImpl)),
	service.NewBidangUrusanServiceImpl,
	wire.Bind(new(service.BidangUrusanService), new(*service.BidangUrusanServiceImpl)),
	controller.NewBidangUrusanControllerImpl,
	wire.Bind(new(controller.BidangUrusanController), new(*controller.BidangUrusanControllerImpl)),
)

var kegiatanSet = wire.NewSet(
	repository.NewKegiatanRepositoryImpl,
	wire.Bind(new(repository.KegiatanRepository), new(*repository.KegiatanRepositoryImpl)),
	service.NewKegiatanServiceImpl,
	wire.Bind(new(service.KegiatanService), new(*service.KegiatanServiceImpl)),
	controller.NewKegiatanControllerImpl,
	wire.Bind(new(controller.KegiatanController), new(*controller.KegiatanControllerImpl)),
)

var roleSet = wire.NewSet(
	repository.NewRoleRepositoryImpl,
	wire.Bind(new(repository.RoleRepository), new(*repository.RoleRepositoryImpl)),
	service.NewRoleServiceImpl,
	wire.Bind(new(service.RoleService), new(*service.RoleServiceImpl)),
	controller.NewRoleControllerImpl,
	wire.Bind(new(controller.RoleController), new(*controller.RoleControllerImpl)),
)

var userSet = wire.NewSet(
	repository.NewUserRepositoryImpl,
	wire.Bind(new(repository.UserRepository), new(*repository.UserRepositoryImpl)),
	service.NewUserServiceImpl,
	wire.Bind(new(service.UserService), new(*service.UserServiceImpl)),
	controller.NewUserControllerImpl,
	wire.Bind(new(controller.UserController), new(*controller.UserControllerImpl)),
)

var seederProviderSet = wire.NewSet(
	dataseeder.NewSeederImpl,
	wire.Bind(new(dataseeder.Seeder), new(*dataseeder.SeederImpl)),
	dataseeder.NewRoleSeederImpl,
	wire.Bind(new(dataseeder.RoleSeeder), new(*dataseeder.RoleSeederImpl)),
	dataseeder.NewUserSeederImpl,
	wire.Bind(new(dataseeder.UserSeeder), new(*dataseeder.UserSeederImpl)),
	dataseeder.NewPegawaiSeederImpl,
	wire.Bind(new(dataseeder.PegawaiSeeder), new(*dataseeder.PegawaiSeederImpl)),
)

var tujuanOpdSet = wire.NewSet(
	repository.NewTujuanOpdRepositoryImpl,
	wire.Bind(new(repository.TujuanOpdRepository), new(*repository.TujuanOpdRepositoryImpl)),
	service.NewTujuanOpdServiceImpl,
	wire.Bind(new(service.TujuanOpdService), new(*service.TujuanOpdServiceImpl)),
	controller.NewTujuanOpdControllerImpl,
	wire.Bind(new(controller.TujuanOpdController), new(*controller.TujuanOpdControllerImpl)),
)

var crosscuttingOpdSet = wire.NewSet(
	repository.NewCrosscuttingOpdRepositoryImpl,
	wire.Bind(new(repository.CrosscuttingOpdRepository), new(*repository.CrosscuttingOpdRepositoryImpl)),
	service.NewCrosscuttingOpdServiceImpl,
	wire.Bind(new(service.CrosscuttingOpdService), new(*service.CrosscuttingOpdServiceImpl)),
	controller.NewCrosscuttingOpdControllerImpl,
	wire.Bind(new(controller.CrosscuttingOpdController), new(*controller.CrosscuttingOpdControllerImpl)),
)

var manualIKSet = wire.NewSet(
	repository.NewManualIKRepositoryImpl,
	wire.Bind(new(repository.ManualIKRepository), new(*repository.ManualIKRepositoryImpl)),
	service.NewManualIKServiceImpl,
	wire.Bind(new(service.ManualIKService), new(*service.ManualIKServiceImpl)),
	controller.NewManualIKControllerImpl,
	wire.Bind(new(controller.ManualIKController), new(*controller.ManualIKControllerImpl)),
)

var reviewSet = wire.NewSet(
	repository.NewReviewRepositoryImpl,
	wire.Bind(new(repository.ReviewRepository), new(*repository.ReviewRepositoryImpl)),
	service.NewReviewServiceImpl,
	wire.Bind(new(service.ReviewService), new(*service.ReviewServiceImpl)),
	controller.NewReviewControllerImpl,
	wire.Bind(new(controller.ReviewController), new(*controller.ReviewControllerImpl)),
)

var periodeSet = wire.NewSet(
	repository.NewPeriodeRepositoryImpl,
	wire.Bind(new(repository.PeriodeRepository), new(*repository.PeriodeRepositoryImpl)),
	service.NewPeriodeServiceImpl,
	wire.Bind(new(service.PeriodeService), new(*service.PeriodeServiceImpl)),
	controller.NewPeriodeControllerImpl,
	wire.Bind(new(controller.PeriodeController), new(*controller.PeriodeControllerImpl)),
)

var tujuanPemdaSet = wire.NewSet(
	repository.NewTujuanPemdaRepositoryImpl,
	wire.Bind(new(repository.TujuanPemdaRepository), new(*repository.TujuanPemdaRepositoryImpl)),
	service.NewTujuanPemdaServiceImpl,
	wire.Bind(new(service.TujuanPemdaService), new(*service.TujuanPemdaServiceImpl)),
	controller.NewTujuanPemdaControllerImpl,
	wire.Bind(new(controller.TujuanPemdaController), new(*controller.TujuanPemdaControllerImpl)),
)

var sasaranPemdaSet = wire.NewSet(
	repository.NewSasaranPemdaRepositoryImpl,
	wire.Bind(new(repository.SasaranPemdaRepository), new(*repository.SasaranPemdaRepositoryImpl)),
	service.NewSasaranPemdaServiceImpl,
	wire.Bind(new(service.SasaranPemdaService), new(*service.SasaranPemdaServiceImpl)),
	controller.NewSasaranPemdaControllerImpl,
	wire.Bind(new(controller.SasaranPemdaController), new(*controller.SasaranPemdaControllerImpl)),
)

var permasalahanRekinSet = wire.NewSet(
	repository.NewPermasalahanRekinRepositoryImpl,
	wire.Bind(new(repository.PermasalahanRekinRepository), new(*repository.PermasalahanRekinRepositoryImpl)),
	service.NewPermasalahanRekinServiceImpl,
	wire.Bind(new(service.PermasalahanRekinService), new(*service.PermasalahanRekinServiceImpl)),
	controller.NewPermasalahanRekinControllerImpl,
	wire.Bind(new(controller.PermasalahanRekinController), new(*controller.PermasalahanRekinControllerImpl)),
)

var ikuSet = wire.NewSet(
	repository.NewIkuRepositoryImpl,
	wire.Bind(new(repository.IkuRepository), new(*repository.IkuRepositoryImpl)),
	service.NewIkuServiceImpl,
	wire.Bind(new(service.IkuService), new(*service.IkuServiceImpl)),
	controller.NewIkuControllerImpl,
	wire.Bind(new(controller.IkuController), new(*controller.IkuControllerImpl)),
)

var sasaranOpdSet = wire.NewSet(
	repository.NewSasaranOpdRepositoryImpl,
	wire.Bind(new(repository.SasaranOpdRepository), new(*repository.SasaranOpdRepositoryImpl)),
	service.NewSasaranOpdServiceImpl,
	wire.Bind(new(service.SasaranOpdService), new(*service.SasaranOpdServiceImpl)),
	controller.NewSasaranOpdControllerImpl,
	wire.Bind(new(controller.SasaranOpdController), new(*controller.SasaranOpdControllerImpl)),
)

var visiPemdaSet = wire.NewSet(
	repository.NewVisiPemdaRepositoryImpl,
	wire.Bind(new(repository.VisiPemdaRepository), new(*repository.VisiPemdaRepositoryImpl)),
	service.NewVisiPemdaServiceImpl,
	wire.Bind(new(service.VisiPemdaService), new(*service.VisiPemdaServiceImpl)),
	controller.NewVisiPemdaControllerImpl,
	wire.Bind(new(controller.VisiPemdaController), new(*controller.VisiPemdaControllerImpl)),
)

var misiPemdaSet = wire.NewSet(
	repository.NewMisiPemdaRepositoryImpl,
	wire.Bind(new(repository.MisiPemdaRepository), new(*repository.MisiPemdaRepositoryImpl)),
	service.NewMisiPemdaServiceImpl,
	wire.Bind(new(service.MisiPemdaService), new(*service.MisiPemdaServiceImpl)),
	controller.NewMisiPemdaControllerImpl,
	wire.Bind(new(controller.MisiPemdaController), new(*controller.MisiPemdaControllerImpl)),
)

var laporanSet = wire.NewSet(
	repository.NewLaporanRepositoryImpl,
	wire.Bind(new(repository.LaporanRepository), new(*repository.LaporanRepositoryImpl)),
	service.NewLaporanServiceImpl,
	wire.Bind(new(service.LaporanService), new(*service.LaporanServiceImpl)),
	controller.NewLaporanControllerImpl,
	wire.Bind(new(controller.LaporanController), new(*controller.LaporanControllerImpl)),
)

func InitializeServer() *http.Server {

	wire.Build(
		app.GetConnection,
		wire.Value([]validator.Option{}),
		validator.New,
		rencanaKinerjaSet,
		rencanaAksiSet,
		pelaksanaanRencanaAksiSet,
		usulanMusrebangSet,
		usulanMandatoriSet,
		usulanPokokPikiranSet,
		usulanInisiatifSet,
		usulanTerpilihSet,
		gambaranUmumSet,
		dasarHukumSet,
		inovasiSet,
		subKegiatanSet,
		subKegiatanTerpilihSet,
		pohonKinerjaOpdSet,
		pegawaiSet,
		lembagaSet,
		jabatanSet,
		pohonKinerjaAdminSet,
		opdSet,
		programSet,
		urusanSet,
		bidangUrusanSet,
		kegiatanSet,
		roleSet,
		userSet,
		tujuanOpdSet,
		crosscuttingOpdSet,
		manualIKSet,
		reviewSet,
		periodeSet,
		tujuanPemdaSet,
		sasaranPemdaSet,
		permasalahanRekinSet,
		ikuSet,
		sasaranOpdSet,
		visiPemdaSet,
		misiPemdaSet,
		laporanSet,
		app.NewRouter,
		wire.Bind(new(http.Handler), new(*httprouter.Router)),
		middleware.NewAuthMiddleware,
		NewServer,
	)

	return nil
}

func InitializeSeeder() dataseeder.Seeder {
	wire.Build(
		app.GetConnection,
		roleSet,
		userSet,
		pegawaiSet,
		seederProviderSet,
	)
	return nil
}
