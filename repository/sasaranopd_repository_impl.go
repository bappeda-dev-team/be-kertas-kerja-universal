package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"errors"
	"fmt"
	"sort"
	"strconv"
)

type SasaranOpdRepositoryImpl struct {
}

func NewSasaranOpdRepositoryImpl() *SasaranOpdRepositoryImpl {
	return &SasaranOpdRepositoryImpl{}
}

func (repository *SasaranOpdRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, KodeOpd string, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]domain.SasaranOpd, error) {
	script := `
     WITH RECURSIVE periode_tahun AS (
        SELECT CAST(? AS CHAR(4)) as tahun
        UNION ALL
        SELECT CAST(CAST(tahun AS UNSIGNED) + 1 AS CHAR(4))
        FROM periode_tahun
        WHERE CAST(tahun AS UNSIGNED) < CAST(? AS UNSIGNED)
    )
    SELECT DISTINCT
        pk.id as pokin_id,
        pk.nama_pohon,
        pk.jenis_pohon,
        pk.level_pohon,
        pk.tahun as tahun_pohon,
        pp.id as pelaksana_id,
        pp.pegawai_id,
        p.nip as pelaksana_nip,
        p.nama as nama_pegawai,
        rk.id as rekin_id,
        rk.nama_rencana_kinerja,
        rk.pegawai_id as rekin_pegawai_id,
        rk.tahun_awal as rekin_tahun_awal,
        rk.tahun_akhir as rekin_tahun_akhir,
        rk.jenis_periode as rekin_jenis_periode,
        i.id as indikator_id,
        i.indikator,
        mik.formula,
        mik.sumber_data,
        t.id as target_id,
        t.tahun as target_tahun,
        t.target,
        t.satuan,
        pt.tahun as tahun_periode
    FROM tb_pohon_kinerja pk
    LEFT JOIN tb_pelaksana_pokin pp ON pk.id = pp.pohon_kinerja_id
    LEFT JOIN tb_pegawai p ON pp.pegawai_id = p.id
    LEFT JOIN tb_rencana_kinerja rk ON pk.id = rk.id_pohon 
        AND rk.tahun_awal = ?
        AND rk.tahun_akhir = ?
        AND rk.jenis_periode = ?
        AND rk.pegawai_id = p.nip  -- Tambahkan kondisi ini untuk memastikan pegawai_id di rencana_kinerja sesuai dengan nip pelaksana
    LEFT JOIN tb_indikator i ON rk.id = i.rencana_kinerja_id
    LEFT JOIN tb_manual_ik mik ON i.id = mik.indikator_id
    LEFT JOIN tb_target t ON i.id = t.indikator_id
    CROSS JOIN periode_tahun pt
    WHERE pk.level_pohon = 4
    AND pk.kode_opd = ?
    AND CAST(pk.tahun AS UNSIGNED) BETWEEN CAST(? AS UNSIGNED) AND CAST(? AS UNSIGNED)
    ORDER BY pk.id, i.id, pt.tahun`

	rows, err := tx.QueryContext(ctx, script,
		tahunAwal, tahunAkhir,
		tahunAwal, tahunAkhir, jenisPeriode,
		KodeOpd,
		tahunAwal, tahunAkhir,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pokinMap := make(map[int]*domain.SasaranOpd)
	rekinMap := make(map[string]bool)     // Ubah menjadi map untuk tracking saja
	indikatorMap := make(map[string]bool) // Ubah menjadi map untuk tracking saja
	pelaksanaMap := make(map[string]bool)

	for rows.Next() {
		var (
			pokinId, levelPohon                  int
			namaPohon, jenisPohon, tahunPohon    string
			pelaksanaId, pegawaiId, pelaksanaNip sql.NullString
			namaPegawai                          sql.NullString
			rekinId, namaRencanaKinerja          sql.NullString
			rekinPegawaiId                       sql.NullString
			rekinTahunAwal, rekinTahunAkhir      sql.NullString
			rekinJenisPeriode                    sql.NullString
			indikatorId, indikator               sql.NullString
			formula, sumberData                  sql.NullString
			targetId, targetTahun                sql.NullString
			targetValue, targetSatuan            sql.NullString
			tahunPeriode                         string
		)

		err := rows.Scan(
			&pokinId, &namaPohon, &jenisPohon, &levelPohon, &tahunPohon,
			&pelaksanaId, &pegawaiId, &pelaksanaNip, &namaPegawai,
			&rekinId, &namaRencanaKinerja, &rekinPegawaiId,
			&rekinTahunAwal, &rekinTahunAkhir, &rekinJenisPeriode,
			&indikatorId, &indikator,
			&formula, &sumberData,
			&targetId, &targetTahun, &targetValue, &targetSatuan,
			&tahunPeriode,
		)
		if err != nil {
			return nil, err
		}

		// Proses SasaranOpd
		sasaranOpd, exists := pokinMap[pokinId]
		if !exists {
			sasaranOpd = &domain.SasaranOpd{
				Id:             pokinId,
				IdPohon:        pokinId,
				NamaPohon:      namaPohon,
				JenisPohon:     jenisPohon,
				LevelPohon:     levelPohon,
				TahunPohon:     tahunPohon,
				RencanaKinerja: make([]domain.RencanaKinerja, 0),
				Pelaksana:      make([]domain.PelaksanaPokin, 0),
			}
			pokinMap[pokinId] = sasaranOpd
		}
		// Proses Pelaksana
		if pelaksanaId.Valid && pegawaiId.Valid && pelaksanaNip.Valid && namaPegawai.Valid {
			pelaksanaKey := fmt.Sprintf("%d-%s", pokinId, pelaksanaId.String)
			if !pelaksanaMap[pelaksanaKey] {
				pelaksanaMap[pelaksanaKey] = true
				sasaranOpd.Pelaksana = append(sasaranOpd.Pelaksana, domain.PelaksanaPokin{
					Id:          pelaksanaId.String,
					PegawaiId:   pegawaiId.String,
					Nip:         pelaksanaNip.String,
					NamaPegawai: namaPegawai.String,
				})
			}
		}

		// Proses RencanaKinerja
		if rekinId.Valid && namaRencanaKinerja.Valid {
			rekinKey := fmt.Sprintf("%d-%s", pokinId, rekinId.String)
			if !rekinMap[rekinKey] {
				rekinMap[rekinKey] = true
				newRekin := domain.RencanaKinerja{
					Id:                 rekinId.String,
					NamaRencanaKinerja: namaRencanaKinerja.String,
					PegawaiId:          rekinPegawaiId.String,
					TahunAwal:          rekinTahunAwal.String,
					TahunAkhir:         rekinTahunAkhir.String,
					Indikator:          make([]domain.Indikator, 0),
				}
				sasaranOpd.RencanaKinerja = append(sasaranOpd.RencanaKinerja, newRekin)
			}

			// Proses Indikator
			if indikatorId.Valid && indikator.Valid {
				indKey := fmt.Sprintf("%s-%s", rekinId.String, indikatorId.String)
				if !indikatorMap[indKey] {
					indikatorMap[indKey] = true

					newInd := domain.Indikator{
						Id:        indikatorId.String,
						Indikator: indikator.String,
						Target:    make([]domain.Target, 0),
					}

					if formula.Valid && sumberData.Valid {
						newInd.ManualIK = &domain.ManualIKSasaranOpd{
							IndikatorId: indikatorId.String,
							Formula:     formula.String,
							SumberData:  sumberData.String,
						}
					}

					// Inisialisasi target kosong untuk semua tahun dalam periode
					tahunAwalInt, _ := strconv.Atoi(rekinTahunAwal.String)
					tahunAkhirInt, _ := strconv.Atoi(rekinTahunAkhir.String)

					for tahun := tahunAwalInt; tahun <= tahunAkhirInt; tahun++ {
						tahunStr := strconv.Itoa(tahun)
						newInd.Target = append(newInd.Target, domain.Target{
							Id:          "",
							IndikatorId: indikatorId.String,
							Tahun:       tahunStr,
							Target:      "",
							Satuan:      "",
						})
					}

					for i := range sasaranOpd.RencanaKinerja {
						if sasaranOpd.RencanaKinerja[i].Id == rekinId.String {
							sasaranOpd.RencanaKinerja[i].Indikator = append(
								sasaranOpd.RencanaKinerja[i].Indikator,
								newInd,
							)
							break
						}
					}
				}

				// Update target jika ada data
				if targetId.Valid && targetTahun.Valid && targetValue.Valid {
					for i := range sasaranOpd.RencanaKinerja[len(sasaranOpd.RencanaKinerja)-1].Indikator {
						if sasaranOpd.RencanaKinerja[len(sasaranOpd.RencanaKinerja)-1].Indikator[i].Id == indikatorId.String {
							for j := range sasaranOpd.RencanaKinerja[len(sasaranOpd.RencanaKinerja)-1].Indikator[i].Target {
								if sasaranOpd.RencanaKinerja[len(sasaranOpd.RencanaKinerja)-1].Indikator[i].Target[j].Tahun == targetTahun.String {
									sasaranOpd.RencanaKinerja[len(sasaranOpd.RencanaKinerja)-1].Indikator[i].Target[j] = domain.Target{
										Id:          targetId.String,
										IndikatorId: indikatorId.String,
										Tahun:       targetTahun.String,
										Target:      targetValue.String,
										Satuan:      targetSatuan.String,
									}
									break
								}
							}
							break
						}
					}
				}
			}
		}
	}

	// Konversi map ke slice dan sort
	var result []domain.SasaranOpd
	for _, sasaranOpd := range pokinMap {
		// Sort indikator
		for i := range sasaranOpd.RencanaKinerja {
			// Pastikan target terurut berdasarkan tahun
			sort.Slice(sasaranOpd.RencanaKinerja[i].Indikator, func(a, b int) bool {
				tahunA, _ := strconv.Atoi(sasaranOpd.RencanaKinerja[i].Indikator[a].Target[0].Tahun)
				tahunB, _ := strconv.Atoi(sasaranOpd.RencanaKinerja[i].Indikator[b].Target[0].Tahun)
				return tahunA < tahunB
			})
		}
		result = append(result, *sasaranOpd)
	}

	return result, nil
}

func (repository *SasaranOpdRepositoryImpl) FindByIdRencanaKinerja(ctx context.Context, tx *sql.Tx, idRencanaKinerja string) (*domain.SasaranOpd, error) {
	script := `
    SELECT DISTINCT
        pk.id as pokin_id,
        pk.nama_pohon,
        pk.jenis_pohon,
        pk.level_pohon,
        pk.tahun as tahun_pohon,
        pp.id as pelaksana_id,
        pp.pegawai_id,
        p.nip as pelaksana_nip,
        p.nama as nama_pegawai,
        rk.id as rekin_id,
        rk.nama_rencana_kinerja,
        rk.pegawai_id as rekin_pegawai_id,
        rk.tahun_awal as rekin_tahun_awal,
        rk.tahun_akhir as rekin_tahun_akhir,
        rk.jenis_periode as rekin_jenis_periode,
        i.id as indikator_id,
        i.indikator,
        mik.formula,
        mik.sumber_data,
        t.id as target_id,
        t.tahun as target_tahun,
        t.target,
        t.satuan
    FROM tb_pohon_kinerja pk
    LEFT JOIN tb_pelaksana_pokin pp ON pk.id = pp.pohon_kinerja_id
    LEFT JOIN tb_pegawai p ON pp.pegawai_id = p.id
    LEFT JOIN tb_rencana_kinerja rk ON pk.id = rk.id_pohon 
    LEFT JOIN tb_indikator i ON rk.id = i.rencana_kinerja_id
    LEFT JOIN tb_manual_ik mik ON i.id = mik.indikator_id
    LEFT JOIN tb_target t ON i.id = t.indikator_id
    WHERE rk.id = ?
    ORDER BY pk.id, i.id, t.tahun`

	rows, err := tx.QueryContext(ctx, script, idRencanaKinerja)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sasaranOpd *domain.SasaranOpd
	rekinMap := make(map[string]bool)
	indikatorMap := make(map[string]bool)
	pelaksanaMap := make(map[string]bool)

	for rows.Next() {
		var (
			pokinId, levelPohon                  int
			namaPohon, jenisPohon, tahunPohon    string
			pelaksanaId, pegawaiId, pelaksanaNip sql.NullString
			namaPegawai                          sql.NullString
			rekinId, namaRencanaKinerja          sql.NullString
			rekinPegawaiId                       sql.NullString
			rekinTahunAwal, rekinTahunAkhir      sql.NullString
			rekinJenisPeriode                    sql.NullString
			indikatorId, indikator               sql.NullString
			formula, sumberData                  sql.NullString
			targetId, targetTahun                sql.NullString
			targetValue, targetSatuan            sql.NullString
		)

		err := rows.Scan(
			&pokinId, &namaPohon, &jenisPohon, &levelPohon, &tahunPohon,
			&pelaksanaId, &pegawaiId, &pelaksanaNip, &namaPegawai,
			&rekinId, &namaRencanaKinerja, &rekinPegawaiId,
			&rekinTahunAwal, &rekinTahunAkhir, &rekinJenisPeriode,
			&indikatorId, &indikator,
			&formula, &sumberData,
			&targetId, &targetTahun, &targetValue, &targetSatuan,
		)
		if err != nil {
			return nil, err
		}

		// Inisialisasi SasaranOpd jika belum ada
		if sasaranOpd == nil {
			sasaranOpd = &domain.SasaranOpd{
				Id:             pokinId,
				IdPohon:        pokinId,
				NamaPohon:      namaPohon,
				JenisPohon:     jenisPohon,
				LevelPohon:     levelPohon,
				TahunPohon:     tahunPohon,
				RencanaKinerja: make([]domain.RencanaKinerja, 0),
				Pelaksana:      make([]domain.PelaksanaPokin, 0),
			}
		}

		// Proses Pelaksana
		if pelaksanaId.Valid && pegawaiId.Valid && pelaksanaNip.Valid && namaPegawai.Valid {
			pelaksanaKey := fmt.Sprintf("%d-%s", pokinId, pelaksanaId.String)
			if !pelaksanaMap[pelaksanaKey] {
				pelaksanaMap[pelaksanaKey] = true
				sasaranOpd.Pelaksana = append(sasaranOpd.Pelaksana, domain.PelaksanaPokin{
					Id:          pelaksanaId.String,
					PegawaiId:   pegawaiId.String,
					Nip:         pelaksanaNip.String,
					NamaPegawai: namaPegawai.String,
				})
			}
		}

		// Proses RencanaKinerja
		if rekinId.Valid && namaRencanaKinerja.Valid {
			rekinKey := fmt.Sprintf("%d-%s", pokinId, rekinId.String)
			if !rekinMap[rekinKey] {
				rekinMap[rekinKey] = true
				newRekin := domain.RencanaKinerja{
					Id:                 rekinId.String,
					NamaRencanaKinerja: namaRencanaKinerja.String,
					PegawaiId:          rekinPegawaiId.String,
					TahunAwal:          rekinTahunAwal.String,
					TahunAkhir:         rekinTahunAkhir.String,
					JenisPeriode:       rekinJenisPeriode.String,
					Indikator:          make([]domain.Indikator, 0),
				}
				sasaranOpd.RencanaKinerja = append(sasaranOpd.RencanaKinerja, newRekin)
			}

			// Proses Indikator
			if indikatorId.Valid && indikator.Valid {
				indKey := fmt.Sprintf("%s-%s", rekinId.String, indikatorId.String)
				if !indikatorMap[indKey] {
					indikatorMap[indKey] = true

					newInd := domain.Indikator{
						Id:        indikatorId.String,
						Indikator: indikator.String,
						Target:    make([]domain.Target, 0),
					}

					if formula.Valid && sumberData.Valid {
						newInd.ManualIK = &domain.ManualIKSasaranOpd{
							IndikatorId: indikatorId.String,
							Formula:     formula.String,
							SumberData:  sumberData.String,
						}
					}

					// Generate target untuk semua tahun dalam periode
					tahunAwalInt, _ := strconv.Atoi(rekinTahunAwal.String)
					tahunAkhirInt, _ := strconv.Atoi(rekinTahunAkhir.String)

					for tahun := tahunAwalInt; tahun <= tahunAkhirInt; tahun++ {
						tahunStr := strconv.Itoa(tahun)
						newInd.Target = append(newInd.Target, domain.Target{
							Id:          "",
							IndikatorId: indikatorId.String,
							Tahun:       tahunStr,
							Target:      "",
							Satuan:      "",
						})
					}

					// Tambahkan indikator ke rencana kinerja yang sesuai
					for i := range sasaranOpd.RencanaKinerja {
						if sasaranOpd.RencanaKinerja[i].Id == rekinId.String {
							sasaranOpd.RencanaKinerja[i].Indikator = append(
								sasaranOpd.RencanaKinerja[i].Indikator,
								newInd,
							)
							break
						}
					}
				}

				// Update target jika ada data
				if targetId.Valid && targetTahun.Valid && targetValue.Valid {
					for i := range sasaranOpd.RencanaKinerja {
						if sasaranOpd.RencanaKinerja[i].Id == rekinId.String {
							for j := range sasaranOpd.RencanaKinerja[i].Indikator {
								if sasaranOpd.RencanaKinerja[i].Indikator[j].Id == indikatorId.String {
									for k := range sasaranOpd.RencanaKinerja[i].Indikator[j].Target {
										if sasaranOpd.RencanaKinerja[i].Indikator[j].Target[k].Tahun == targetTahun.String {
											sasaranOpd.RencanaKinerja[i].Indikator[j].Target[k] = domain.Target{
												Id:          targetId.String,
												IndikatorId: indikatorId.String,
												Tahun:       targetTahun.String,
												Target:      targetValue.String,
												Satuan:      targetSatuan.String,
											}
											break
										}
									}
									break
								}
							}
							break
						}
					}
				}
			}
		}
	}

	if sasaranOpd == nil {
		return nil, errors.New("rencana kinerja tidak ditemukan")
	}

	return sasaranOpd, nil
}

func (repository *SasaranOpdRepositoryImpl) FindIdPokinSasaran(ctx context.Context, tx *sql.Tx, id int) (domain.PohonKinerja, error) {
	scriptPokin := `
    SELECT DISTINCT
        pk.id, 
        pk.parent, 
        pk.nama_pohon, 
        pk.jenis_pohon, 
        pk.level_pohon, 
        pk.kode_opd, 
        pk.keterangan, 
        pk.tahun,
        pk.status,
        i.id as indikator_id,
        i.indikator as nama_indikator,
        t.id as target_id,
        t.target,
        t.satuan,
        t.tahun as tahun_target
    FROM 
        tb_pohon_kinerja pk 
        LEFT JOIN tb_indikator i ON pk.id = i.pokin_id
        LEFT JOIN tb_target t ON i.id = t.indikator_id
    WHERE 
        pk.id = ?
    ORDER BY t.id DESC
    LIMIT 1`

	rows, err := tx.QueryContext(ctx, scriptPokin, id)
	if err != nil {
		return domain.PohonKinerja{}, fmt.Errorf("error querying pohon kinerja: %v", err)
	}
	defer rows.Close()

	var pohonKinerja domain.PohonKinerja
	indikatorMap := make(map[string]*domain.Indikator)
	dataFound := false

	for rows.Next() {
		var (
			indikatorId, namaIndikator            sql.NullString
			targetId, target, satuan, tahunTarget sql.NullString
		)

		err := rows.Scan(
			&pohonKinerja.Id,
			&pohonKinerja.Parent,
			&pohonKinerja.NamaPohon,
			&pohonKinerja.JenisPohon,
			&pohonKinerja.LevelPohon,
			&pohonKinerja.KodeOpd,
			&pohonKinerja.Keterangan,
			&pohonKinerja.Tahun,
			&pohonKinerja.Status,
			&indikatorId,
			&namaIndikator,
			&targetId,
			&target,
			&satuan,
			&tahunTarget,
		)
		if err != nil {
			return domain.PohonKinerja{}, fmt.Errorf("error scanning row: %v", err)
		}

		dataFound = true

		if indikatorId.Valid && namaIndikator.Valid {
			ind := &domain.Indikator{
				Id:        indikatorId.String,
				Indikator: namaIndikator.String,
				PokinId:   fmt.Sprint(pohonKinerja.Id),
				Target:    []domain.Target{},
			}

			if targetId.Valid && target.Valid && satuan.Valid {
				targetObj := domain.Target{
					Id:          targetId.String,
					IndikatorId: indikatorId.String,
					Target:      target.String,
					Satuan:      satuan.String,
					Tahun:       tahunTarget.String,
				}
				ind.Target = append(ind.Target, targetObj)
			}

			indikatorMap[indikatorId.String] = ind
			pohonKinerja.Indikator = append(pohonKinerja.Indikator, *ind)
		}
	}

	if !dataFound {
		return domain.PohonKinerja{}, fmt.Errorf("pohon kinerja with id %d not found", id)
	}

	return pohonKinerja, nil
}
