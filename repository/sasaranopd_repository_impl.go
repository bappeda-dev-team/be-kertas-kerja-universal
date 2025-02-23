package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
	"sort"
)

type SasaranOpdRepositoryImpl struct {
}

func NewSasaranOpdRepositoryImpl() *SasaranOpdRepositoryImpl {
	return &SasaranOpdRepositoryImpl{}
}

// func (repository *SasaranOpdRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, KodeOpd string, tahunAwal string, tahunAkhir string) ([]domain.SasaranOpd, error) {
// 	script := `
//         WITH RECURSIVE periode_tahun AS (
//             SELECT CAST(? AS CHAR(4)) as tahun
//             UNION ALL
//             SELECT CAST(CAST(tahun AS UNSIGNED) + 1 AS CHAR(4))
//             FROM periode_tahun
//             WHERE CAST(tahun AS UNSIGNED) < CAST(? AS UNSIGNED)
//         )
//         SELECT DISTINCT
//             pk.id as pokin_id,
//             pk.nama_pohon,
//             pk.jenis_pohon,
//             pk.level_pohon,
//             pk.tahun as tahun_pohon,
//             pp.id as pelaksana_id,
//             pp.pegawai_id,
//             u.id as user_id,
//             u.nip as pelaksana_nip,
//             rk.id as rekin_id,
//             rk.nama_rencana_kinerja,
//             rk.pegawai_id as rekin_pegawai_id,
//             p.id as periode_id,
//             p.tahun_awal,
//             p.tahun_akhir,
//             i.id as indikator_id,
//             i.indikator,
//             mik.formula,
//             mik.sumber_data,
//             t.id as target_id,
//             t.tahun as target_tahun,
//             t.target,
//             t.satuan,
//             pt.tahun as tahun_periode
//         FROM tb_pohon_kinerja pk
//         INNER JOIN tb_pelaksana_pokin pp ON pk.id = pp.pohon_kinerja_id
//         INNER JOIN tb_users u ON CAST(pp.pegawai_id AS CHAR) = CAST(u.id AS CHAR)
//         LEFT JOIN tb_rencana_kinerja rk ON pk.id = rk.id_pohon
//         LEFT JOIN tb_periode p ON rk.periode_id = p.id
//         LEFT JOIN tb_indikator i ON rk.id = i.rencana_kinerja_id
//         LEFT JOIN tb_manual_ik mik ON i.id = mik.indikator_id
//         LEFT JOIN tb_target t ON i.id = t.indikator_id
//         CROSS JOIN periode_tahun pt
//         WHERE pk.level_pohon = 4
//         AND pk.kode_opd = ?
//         AND pk.tahun BETWEEN ? AND ?
//         ORDER BY pk.id, i.id, pt.tahun`

// 	rows, err := tx.QueryContext(ctx, script, tahunAwal, tahunAkhir, KodeOpd, tahunAwal, tahunAkhir)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	pokinMap := make(map[int]*domain.SasaranOpd)
// 	indikatorMap := make(map[string]*domain.Indikator)
// 	pelaksanaMap := make(map[string]bool)

// 	for rows.Next() {
// 		var (
// 			pokinId, levelPohon                  int
// 			namaPohon, jenisPohon, tahunPohon    string
// 			pelaksanaId, pegawaiId, pelaksanaNip sql.NullString
// 			rekinId, namaRencanaKinerja          sql.NullString
// 			rekinPegawaiId                       sql.NullString
// 			periodeId                            sql.NullInt64
// 			periodeTahunAwal, periodeTahunAkhir  sql.NullString
// 			indikatorId, indikator               sql.NullString
// 			formula, sumberData                  sql.NullString
// 			targetId, targetTahun                sql.NullString
// 			targetValue, targetSatuan            sql.NullString
// 			tahunPeriode                         string
// 			userId                               sql.NullInt64
// 		)

// 		err := rows.Scan(
// 			&pokinId, &namaPohon, &jenisPohon, &levelPohon, &tahunPohon,
// 			&pelaksanaId, &pegawaiId, &userId, &pelaksanaNip,
// 			&rekinId, &namaRencanaKinerja, &rekinPegawaiId,
// 			&periodeId, &periodeTahunAwal, &periodeTahunAkhir,
// 			&indikatorId, &indikator,
// 			&formula, &sumberData,
// 			&targetId, &targetTahun, &targetValue, &targetSatuan,
// 			&tahunPeriode,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		sasaranOpd, exists := pokinMap[pokinId]
// 		if !exists {
// 			sasaranOpd = &domain.SasaranOpd{
// 				Id:                  pokinId,
// 				IdPohon:             pokinId,
// 				NamaPohon:           namaPohon,
// 				JenisPohon:          jenisPohon,
// 				LevelPohon:          levelPohon,
// 				TahunPohon:          tahunPohon,
// 				TahunAwalPeriode:    tahunAwal,
// 				TahunAkhirPeriode:   tahunAkhir,
// 				Pelaksana:           []domain.PelaksanaPokin{},
// 				IndikatorSasaranOpd: []domain.Indikator{},
// 			}
// 			pokinMap[pokinId] = sasaranOpd
// 		}

// 		// Proses Pelaksana
// 		if pelaksanaId.Valid && pegawaiId.Valid && pelaksanaNip.Valid {
// 			key := fmt.Sprintf("%s-%s", pelaksanaId.String, pegawaiId.String)
// 			if !pelaksanaMap[key] {
// 				pelaksanaMap[key] = true
// 				sasaranOpd.Pelaksana = append(sasaranOpd.Pelaksana, domain.PelaksanaPokin{
// 					Id:        pelaksanaId.String,
// 					PegawaiId: pegawaiId.String,
// 					Nip:       pelaksanaNip.String,
// 				})
// 			}
// 		}

// 		// Proses RencanaKinerja
// 		if rekinId.Valid && namaRencanaKinerja.Valid {
// 			sasaranOpd.IdRencanaKinerja = rekinId.String
// 			sasaranOpd.NamaRencanaKinerja = namaRencanaKinerja.String
// 			if rekinPegawaiId.Valid {
// 				sasaranOpd.PegawaiId = rekinPegawaiId.String
// 			}
// 			if periodeTahunAwal.Valid && periodeTahunAkhir.Valid {
// 				sasaranOpd.TahunAwalRencana = periodeTahunAwal.String
// 				sasaranOpd.TahunAkhirRencana = periodeTahunAkhir.String
// 			}

// 			// Proses Indikator
// 			if indikatorId.Valid && indikator.Valid {
// 				indKey := fmt.Sprintf("%s-%s", rekinId.String, indikatorId.String)
// 				ind, exists := indikatorMap[indKey]
// 				if !exists {
// 					ind = &domain.Indikator{
// 						Id:        indikatorId.String,
// 						Indikator: indikator.String,
// 						Target:    []domain.Target{},
// 					}

// 					if formula.Valid && sumberData.Valid {
// 						ind.ManualIK = &domain.ManualIKSasaranOpd{
// 							IndikatorId: indikatorId.String,
// 							Formula:     formula.String,
// 							SumberData:  sumberData.String,
// 						}
// 					}

// 					indikatorMap[indKey] = ind
// 					sasaranOpd.IndikatorSasaranOpd = append(sasaranOpd.IndikatorSasaranOpd, *ind)
// 				}

// 				// Proses Target
// 				target := domain.Target{
// 					Id:          targetId.String,
// 					IndikatorId: indikatorId.String,
// 					Tahun:       tahunPeriode,
// 					Target:      "",
// 					Satuan:      "",
// 				}

// 				if targetTahun.Valid && targetValue.Valid &&
// 					targetTahun.String == tahunPeriode {
// 					target.Target = targetValue.String
// 					if targetSatuan.Valid {
// 						target.Satuan = targetSatuan.String
// 					}
// 				}

// 				for i := range sasaranOpd.IndikatorSasaranOpd {
// 					if sasaranOpd.IndikatorSasaranOpd[i].Id == indikatorId.String {
// 						sasaranOpd.IndikatorSasaranOpd[i].Target = append(
// 							sasaranOpd.IndikatorSasaranOpd[i].Target,
// 							target,
// 						)
// 						break
// 					}
// 				}
// 			}
// 		}
// 	}

// 	var result []domain.SasaranOpd
// 	for _, sasaranOpd := range pokinMap {
// 		for i := range sasaranOpd.IndikatorSasaranOpd {
// 			sort.Slice(sasaranOpd.IndikatorSasaranOpd[i].Target, func(a, b int) bool {
// 				return sasaranOpd.IndikatorSasaranOpd[i].Target[a].Tahun <
// 					sasaranOpd.IndikatorSasaranOpd[i].Target[b].Tahun
// 			})
// 		}
// 		result = append(result, *sasaranOpd)
// 	}

// 	return result, nil
// }

// editet last and true
func (repository *SasaranOpdRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, KodeOpd string, tahunAwal string, tahunAkhir string) ([]domain.SasaranOpd, error) {
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
            u.nip as pelaksana_nip,
            rk.id as rekin_id,
            rk.nama_rencana_kinerja,
            rk.pegawai_id as rekin_pegawai_id,
            p.id as periode_id,
            p.tahun_awal,
            p.tahun_akhir,
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
        LEFT JOIN tb_users u ON pp.pegawai_id = u.id
        LEFT JOIN tb_rencana_kinerja rk ON pk.id = rk.id_pohon
        LEFT JOIN tb_periode p ON rk.periode_id = p.id
        LEFT JOIN tb_indikator i ON rk.id = i.rencana_kinerja_id
        LEFT JOIN tb_manual_ik mik ON i.id = mik.indikator_id
        LEFT JOIN tb_target t ON i.id = t.indikator_id
        CROSS JOIN periode_tahun pt
        WHERE pk.level_pohon = 4
        AND pk.kode_opd = ?
        AND pk.tahun BETWEEN ? AND ?
        ORDER BY pk.id, i.id, pt.tahun`

	rows, err := tx.QueryContext(ctx, script, tahunAwal, tahunAkhir, KodeOpd, tahunAwal, tahunAkhir)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pokinMap := make(map[int]*domain.SasaranOpd)
	indikatorMap := make(map[string]*domain.Indikator)
	pelaksanaMap := make(map[string]bool)

	for rows.Next() {
		var (
			pokinId, levelPohon                  int
			namaPohon, jenisPohon, tahunPohon    string
			pelaksanaId, pegawaiId, pelaksanaNip sql.NullString
			rekinId, namaRencanaKinerja          sql.NullString
			rekinPegawaiId                       sql.NullString
			periodeId                            sql.NullInt64
			periodeTahunAwal, periodeTahunAkhir  sql.NullString
			indikatorId, indikator               sql.NullString
			formula, sumberData                  sql.NullString
			targetId, targetTahun                sql.NullString
			targetValue, targetSatuan            sql.NullString
			tahunPeriode                         string
		)

		err := rows.Scan(
			&pokinId, &namaPohon, &jenisPohon, &levelPohon, &tahunPohon,
			&pelaksanaId, &pegawaiId, &pelaksanaNip,
			&rekinId, &namaRencanaKinerja, &rekinPegawaiId,
			&periodeId, &periodeTahunAwal, &periodeTahunAkhir,
			&indikatorId, &indikator,
			&formula, &sumberData,
			&targetId, &targetTahun, &targetValue, &targetSatuan,
			&tahunPeriode,
		)
		if err != nil {
			return nil, err
		}

		sasaranOpd, exists := pokinMap[pokinId]
		if !exists {
			sasaranOpd = &domain.SasaranOpd{
				Id:                  pokinId,
				IdPohon:             pokinId,
				NamaPohon:           namaPohon,
				JenisPohon:          jenisPohon,
				LevelPohon:          levelPohon,
				TahunPohon:          tahunPohon,
				TahunAwalPeriode:    tahunAwal,
				TahunAkhirPeriode:   tahunAkhir,
				Pelaksana:           []domain.PelaksanaPokin{},
				IndikatorSasaranOpd: []domain.Indikator{},
			}
			pokinMap[pokinId] = sasaranOpd
		}

		// Proses Pelaksana
		if pelaksanaId.Valid && pegawaiId.Valid && pelaksanaNip.Valid {
			key := fmt.Sprintf("%s-%s", pelaksanaId.String, pegawaiId.String)
			if !pelaksanaMap[key] {
				pelaksanaMap[key] = true
				sasaranOpd.Pelaksana = append(sasaranOpd.Pelaksana, domain.PelaksanaPokin{
					Id:        pelaksanaId.String,
					PegawaiId: pegawaiId.String,
					Nip:       pelaksanaNip.String,
				})
			}
		}

		// Proses RencanaKinerja
		if rekinId.Valid && namaRencanaKinerja.Valid {
			sasaranOpd.IdRencanaKinerja = rekinId.String
			sasaranOpd.NamaRencanaKinerja = namaRencanaKinerja.String
			if rekinPegawaiId.Valid {
				sasaranOpd.PegawaiId = rekinPegawaiId.String
			}
			if periodeTahunAwal.Valid && periodeTahunAkhir.Valid {
				sasaranOpd.TahunAwalRencana = periodeTahunAwal.String
				sasaranOpd.TahunAkhirRencana = periodeTahunAkhir.String
			}

			// Proses Indikator
			if indikatorId.Valid && indikator.Valid {
				indKey := fmt.Sprintf("%s-%s", rekinId.String, indikatorId.String)
				ind, exists := indikatorMap[indKey]
				if !exists {
					ind = &domain.Indikator{
						Id:        indikatorId.String,
						Indikator: indikator.String,
						Target:    []domain.Target{},
					}

					if formula.Valid && sumberData.Valid {
						ind.ManualIK = &domain.ManualIKSasaranOpd{
							IndikatorId: indikatorId.String,
							Formula:     formula.String,
							SumberData:  sumberData.String,
						}
					}

					indikatorMap[indKey] = ind
					sasaranOpd.IndikatorSasaranOpd = append(sasaranOpd.IndikatorSasaranOpd, *ind)
				}

				// Proses Target
				target := domain.Target{
					Id:          targetId.String,
					IndikatorId: indikatorId.String,
					Tahun:       tahunPeriode,
					Target:      "",
					Satuan:      "",
				}

				if targetTahun.Valid && targetValue.Valid &&
					targetTahun.String == tahunPeriode {
					target.Target = targetValue.String
					if targetSatuan.Valid {
						target.Satuan = targetSatuan.String
					}
				}

				for i := range sasaranOpd.IndikatorSasaranOpd {
					if sasaranOpd.IndikatorSasaranOpd[i].Id == indikatorId.String {
						sasaranOpd.IndikatorSasaranOpd[i].Target = append(
							sasaranOpd.IndikatorSasaranOpd[i].Target,
							target,
						)
						break
					}
				}
			}
		}
	}

	var result []domain.SasaranOpd
	for _, sasaranOpd := range pokinMap {
		for i := range sasaranOpd.IndikatorSasaranOpd {
			sort.Slice(sasaranOpd.IndikatorSasaranOpd[i].Target, func(a, b int) bool {
				return sasaranOpd.IndikatorSasaranOpd[i].Target[a].Tahun <
					sasaranOpd.IndikatorSasaranOpd[i].Target[b].Tahun
			})
		}
		result = append(result, *sasaranOpd)
	}

	return result, nil
}

// func (repository *SasaranOpdRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, KodeOpd string, tahun string) ([]domain.SasaranOpd, error) {
// 	tahunInt, err := strconv.Atoi(tahun)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid tahun format: %v", err)
// 	}
// 	tahunAkhir := strconv.Itoa(tahunInt + 4)

// 	script := `
//         WITH RECURSIVE periode AS (
//             SELECT CAST(? AS CHAR(4)) as tahun
//             UNION ALL
//             SELECT CAST(CAST(tahun AS UNSIGNED) + 1 AS CHAR(4))
//             FROM periode
//             WHERE CAST(tahun AS UNSIGNED) < CAST(? AS UNSIGNED)
//         )
//         SELECT
//             pk.id as pokin_id,
//             pk.nama_pohon,
//             pk.jenis_pohon,
//             pk.level_pohon,
//             pk.tahun as tahun_pohon,
//             pp.id as pelaksana_id,
//             pp.pegawai_id,
//             rk.id as rekin_id,
//             rk.nama_rencana_kinerja,
//             rk.pegawai_id as rekin_pegawai_id,
//             rk.tahun as tahun_rencana,
//             i.id as indikator_id,
//             i.indikator,
//             mik.formula,
//             mik.sumber_data,
//             t.id as target_id,
//             t.tahun as target_tahun,
//             t.target,
//             t.satuan,
//             p.tahun as periode_tahun
//         FROM tb_pohon_kinerja pk
//         LEFT JOIN tb_pelaksana_pokin pp ON pk.id = pp.pohon_kinerja_id
//         LEFT JOIN tb_rencana_kinerja rk ON pk.id = rk.id_pohon
//         LEFT JOIN tb_indikator i ON rk.id = i.rencana_kinerja_id
//         LEFT JOIN tb_manual_ik mik ON i.id = mik.indikator_id
//         LEFT JOIN tb_target t ON i.id = t.indikator_id
//         CROSS JOIN periode p
//         WHERE pk.level_pohon = 4
//         AND pk.kode_opd = ?
//         AND pk.tahun = ?
//         ORDER BY pk.id, i.id, p.tahun`

// 	rows, err := tx.QueryContext(ctx, script, tahun, tahunAkhir, KodeOpd, tahun)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	pokinMap := make(map[int]*domain.SasaranOpd)
// 	indikatorMap := make(map[string]*domain.Indikator)
// 	pelaksanaMap := make(map[string]bool)

// 	for rows.Next() {
// 		var (
// 			pokinId, levelPohon        int
// 			namaPohon, jenisPohon      string
// 			tahunPohon                 string
// 			pelaksanaId, pegawaiId     sql.NullString
// 			rekinId, rekinPegawaiId    sql.NullString
// 			namaRencanaKinerja         sql.NullString
// 			tahunRencana               sql.NullString
// 			indikatorId, indikator     sql.NullString
// 			formula, sumberData        sql.NullString
// 			targetId                   sql.NullString
// 			targetTahun, targetSasaran sql.NullString
// 			satuan, periodeTahun       sql.NullString
// 		)

// 		err := rows.Scan(
// 			&pokinId, &namaPohon, &jenisPohon, &levelPohon, &tahunPohon,
// 			&pelaksanaId, &pegawaiId,
// 			&rekinId, &namaRencanaKinerja, &rekinPegawaiId, &tahunRencana,
// 			&indikatorId, &indikator,
// 			&formula, &sumberData,
// 			&targetId, &targetTahun, &targetSasaran, &satuan,
// 			&periodeTahun,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// Inisialisasi atau ambil SasaranOpd
// 		sasaranOpd, exists := pokinMap[pokinId]
// 		if !exists {
// 			sasaranOpd = &domain.SasaranOpd{
// 				Id:         pokinId,
// 				IdPohon:    pokinId,
// 				NamaPohon:  namaPohon,
// 				JenisPohon: jenisPohon,
// 				LevelPohon: levelPohon,
// 				TahunPohon: tahunPohon,
// 			}
// 			pokinMap[pokinId] = sasaranOpd
// 		}

// 		// Proses Pelaksana
// 		if pelaksanaId.Valid && pegawaiId.Valid {
// 			key := fmt.Sprintf("%s-%s", pelaksanaId.String, pegawaiId.String)
// 			if !pelaksanaMap[key] {
// 				pelaksanaMap[key] = true
// 				sasaranOpd.Pelaksana = append(sasaranOpd.Pelaksana, domain.PelaksanaPokin{
// 					Id:        pelaksanaId.String,
// 					PegawaiId: pegawaiId.String,
// 				})
// 			}
// 		}

// 		// Proses RencanaKinerja
// 		if rekinId.Valid && namaRencanaKinerja.Valid {
// 			sasaranOpd.IdRencanaKinerja = rekinId.String
// 			sasaranOpd.NamaRencanaKinerja = namaRencanaKinerja.String
// 			if rekinPegawaiId.Valid {
// 				sasaranOpd.PegawaiId = rekinPegawaiId.String
// 			}
// 			if tahunRencana.Valid {
// 				sasaranOpd.TahunAwalRencana = tahun
// 				tahunAkhirInt, _ := strconv.Atoi(tahun)
// 				sasaranOpd.TahunAkhirRencana = strconv.Itoa(tahunAkhirInt + 4)
// 			}

// 			// Proses Indikator
// 			if indikatorId.Valid && indikator.Valid {
// 				indKey := fmt.Sprintf("%s-%s", rekinId.String, indikatorId.String)
// 				ind, exists := indikatorMap[indKey]
// 				if !exists {
// 					ind = &domain.Indikator{
// 						Id:        indikatorId.String,
// 						Indikator: indikator.String,
// 					}

// 					// Proses Manual IK jika ada
// 					if formula.Valid && sumberData.Valid {
// 						ind.ManualIK = &domain.ManualIKSasaranOpd{
// 							IndikatorId: indikatorId.String,
// 							Formula:     formula.String,
// 							SumberData:  sumberData.String,
// 						}
// 					}

// 					indikatorMap[indKey] = ind
// 					sasaranOpd.IndikatorSasaranOpd = append(sasaranOpd.IndikatorSasaranOpd, *ind)
// 				}

// 				// Proses Target untuk setiap periode
// 				target := domain.Target{
// 					Id:          targetId.String,
// 					IndikatorId: indikatorId.String,
// 					Tahun:       periodeTahun.String,
// 					Target:      "",
// 					Satuan:      "",
// 				}

// 				if targetTahun.Valid && targetSasaran.Valid &&
// 					targetTahun.String == periodeTahun.String {
// 					target.Target = targetSasaran.String
// 					if satuan.Valid {
// 						target.Satuan = satuan.String
// 					}
// 				}

// 				for i := range sasaranOpd.IndikatorSasaranOpd {
// 					if sasaranOpd.IndikatorSasaranOpd[i].Id == indikatorId.String {
// 						sasaranOpd.IndikatorSasaranOpd[i].Target = append(
// 							sasaranOpd.IndikatorSasaranOpd[i].Target,
// 							target,
// 						)
// 						break
// 					}
// 				}
// 			}
// 		}
// 	}

// 	var result []domain.SasaranOpd
// 	for _, sasaranOpd := range pokinMap {
// 		for i := range sasaranOpd.IndikatorSasaranOpd {
// 			sort.Slice(sasaranOpd.IndikatorSasaranOpd[i].Target, func(a, b int) bool {
// 				return sasaranOpd.IndikatorSasaranOpd[i].Target[a].Tahun <
// 					sasaranOpd.IndikatorSasaranOpd[i].Target[b].Tahun
// 			})
// 		}
// 		result = append(result, *sasaranOpd)
// 	}

// 	return result, nil
// }
