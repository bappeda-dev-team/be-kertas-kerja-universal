package dataseeder

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/repository"
	"log"
)

type RoleSeederImpl struct {
	RoleRepository repository.RoleRepository
}

func NewRoleSeederImpl(roleRepository repository.RoleRepository) *RoleSeederImpl {
	return &RoleSeederImpl{
		RoleRepository: roleRepository,
	}
}

func (seeder *RoleSeederImpl) Seed(ctx context.Context, tx *sql.Tx) error {
	roles, err := seeder.RoleRepository.FindAll(ctx, tx)
	if err != nil {
		return err
	}

	if len(roles) > 0 {
		log.Println("Roles sudah ada, skip seeding roles")
		return nil
	}

	defaultRoles := []domain.Roles{
		{
			Role: "super_admin",
		},
		{
			Role: "admin_dinas_pendidikan_kebudayaan",
		},
		{
			Role: "admin_dinas_kesehatan",
		},
		{
			Role: "admin_puskesmas_kebonsari",
		},
		{
			Role: "admin_puskesmas_gantrung",
		},
		{
			Role: "admin_puskesmas_geger",
		},
		{
			Role: "admin_puskesmas_kaibon",
		},
		{
			Role: "admin_puskesmas_mlilir",
		},
		{
			Role: "admin_puskesmas_bangunsari",
		},
		{
			Role: "admin_puskesmas_dagangan",
		},
		{
			Role: "admin_puskesmas_jetis",
		},
		{
			Role: "admin_puskesmas_wungu",
		},
		{
			Role: "admin_puskesmas_mojopurno",
		},
		{
			Role: "admin_puskesmas_kare",
		},
		{
			Role: "admin_puskesmas_gemarang",
		},
		{
			Role: "admin_puskesmas_saradan",
		},
		{
			Role: "admin_puskesmas_sumbersari",
		},
		{
			Role: "admin_puskesmas_pilangkenceng",
		},
		{
			Role: "admin_puskesmas_krebet",
		},
		{
			Role: "admin_puskesmas_mejayan",
		},
		{
			Role: "admin_puskesmas_klecorejo",
		},
		{
			Role: "admin_puskesmas_wonoasri",
		},
		{
			Role: "admin_puskesmas_balerejo",
		},
		{
			Role: "admin_puskesmas_simo",
		},
		{
			Role: "admin_puskesmas_madiun",
		},
		{
			Role: "admin_puskesmas_dimong",
		},
		{
			Role: "admin_puskesmas_sawahan",
		},
		{
			Role: "admin_puskesmas_klagenserut",
		},
		{
			Role: "admin_puskesmas_jiwan",
		},
		{
			Role: "admin_rsud_caruban",
		},
		{
			Role: "admin_rsud_dolopo",
		},
		{
			Role: "admin_dinas_pekerjaan_umum_dan_penataan_ruang",
		},
		{
			Role: "admin_dinas_perumahan_dan_kawasan_permukiman",
		},
		{
			Role: "admin_satuan_polisi_pamong_praja_dan_pemadam_kebakaran",
		},
		{
			Role: "admin_badan_penanggulangan_bencana_daerah",
		},
		{
			Role: "admin_dinas_sosial",
		},
		{
			Role: "admin_dinas_tenaga_kerja_dan_perindustrian",
		},
		{
			Role: "admin_dinas_ketahanan_pangan_dan_peternakan",
		},
		{
			Role: "admin_dinas_lingkungan_hidup",
		},
		{
			Role: "admin_dinas_kependudukan_dan_pencatatan_sipil",
		},
		{
			Role: "admin_dinas_pemberdayaan_masyarakat_dan_desa",
		},
		{
			Role: "admin_dinas_pengendalian_penduduk_dan_keluarga_berencana_pemberdayaan_perempuan_dan_perlindungan_anak",
		},
		{
			Role: "admin_dinas_perhubungan",
		},
		{
			Role: "admin_dinas_komunikasi_dan_informatika",
		},
		{
			Role: "admin_dinas_perdagangan_koperasi_dan_usaha_mikro",
		},
		{
			Role: "admin_dinas_penanaman_modal_dan_pelayanan_terpadu_satu_pintu",
		},
		{
			Role: "admin_dinas_pariwisata_pemuda_dan_olahraga",
		},
		{
			Role: "admin_dinas_perpustakaan_dan_kearsipan",
		},
		{
			Role: "admin_dinas_pertanian_dan_perikanan",
		},
		{
			Role: "admin_sekretariat_daerah",
		},
		{
			Role: "admin_bagian_administrasi_pemerintahan",
		},
		{
			Role: "admin_bagian_hukum",
		},
		{
			Role: "admin_bagian_pengadaan_barang_dan_jasa",
		},
		{
			Role: "admin_bagian_kesejahteraan_rakyat",
		},
		{
			Role: "admin_bagian_perekonomian",
		},
		{
			Role: "admin_bagian_umum",
		},
		{
			Role: "admin_bagian_organisasi",
		},
		{
			Role: "admin_bagian_protokol_dan_komunikasi_pimpinan",
		},
		{
			Role: "admin_bagian_administrasi_pembangunan",
		},
		{
			Role: "admin_sekretariat_dewan_perwakilan_rakyat_daerah",
		},
		{
			Role: "admin_badan_perencanaan_pembangunan_riset_dan_inovasi_daerah",
		},
		{
			Role: "admin_badan_pengelolaan_keuangan_dan_aset_daerah",
		},
		{
			Role: "admin_badan_pendapatan_daerah",
		},
		{
			Role: "admin_badan_kepegawaian_dan_pengembangan_sumber_daya_manusia",
		},
		{
			Role: "admin_inspektorat",
		},
		{
			Role: "admin_kecamatan_balerejo",
		},
		{
			Role: "admin_kecamatan_dagangan",
		},
		{
			Role: "admin_kecamatan_dolopo",
		},
		{
			Role: "admin_kelurahan_bangunsari_dolopo",
		},
		{
			Role: "admin_kelurahan_mlilir",
		},
		{
			Role: "admin_kecamatan_geger",
		},
		{
			Role: "admin_kecamatan_gemarang",
		},
		{
			Role: "admin_kecamatan_jiwan",
		},
		{
			Role: "admin_kecamatan_kebonsari",
		},
		{
			Role: "admin_kecamatan_kare",
		},
		{
			Role: "admin_kecamatan_madiun",
		},
		{
			Role: "admin_kelurahan_nglames",
		},
		{
			Role: "admin_kecamatan_mejayan",
		},
		{
			Role: "admin_kelurahan_bangunsari_mejayan",
		},
		{
			Role: "admin_kelurahan_krajan",
		},
		{
			Role: "admin_kelurahan_pandean",
		},
		{
			Role: "admin_kecamatan_pilangkenceng",
		},
		{
			Role: "admin_kecamatan_sawahan",
		},
		{
			Role: "admin_kecamatan_saradan",
		},
		{
			Role: "admin_kecamatan_wungu",
		},
		{
			Role: "admin_kelurahan_wungu",
		},
		{
			Role: "admin_kelurahan_munggut",
		},
		{
			Role: "admin_kecamatan_wonoasri",
		},
		{
			Role: "admin_badan_kesatuan_bangsa_dan_politik",
		},
	}

	for _, role := range defaultRoles {
		_, err := seeder.RoleRepository.Create(ctx, tx, role)
		if err != nil {
			return err
		}
	}
	log.Println("Roles berhasil di-seed")
	return nil
}
