package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
	"sort"
	"strconv"
)

type TujuanPemdaRepositoryImpl struct {
}

func NewTujuanPemdaRepositoryImpl() *TujuanPemdaRepositoryImpl {
	return &TujuanPemdaRepositoryImpl{}
}

func (repository *TujuanPemdaRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, tujuanPemda domain.TujuanPemda) (domain.TujuanPemda, error) {
	query := "INSERT INTO tb_tujuan_pemda(id, tujuan_pemda, tematik_id, periode_id) VALUES (?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, query, tujuanPemda.Id, tujuanPemda.TujuanPemda, tujuanPemda.TematikId, tujuanPemda.PeriodeId)
	if err != nil {
		return tujuanPemda, err
	}
	return tujuanPemda, nil
}

func (repository *TujuanPemdaRepositoryImpl) CreateIndikator(ctx context.Context, tx *sql.Tx, indikator domain.Indikator) (domain.Indikator, error) {
	query := "INSERT INTO tb_indikator(id, tujuan_pemda_id, indikator, rumus_perhitungan, sumber_data) VALUES (?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, query, indikator.Id, indikator.TujuanPemdaId, indikator.Indikator, indikator.RumusPerhitungan, indikator.SumberData)
	if err != nil {
		return indikator, err
	}
	return indikator, nil
}

func (repository *TujuanPemdaRepositoryImpl) CreateTarget(ctx context.Context, tx *sql.Tx, target domain.Target) error {
	query := "INSERT INTO tb_target(id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, query, target.Id, target.IndikatorId, target.Target, target.Satuan, target.Tahun)
	return err
}

func (repository *TujuanPemdaRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, tujuanPemda domain.TujuanPemda) (domain.TujuanPemda, error) {
	// Update tujuan pemda
	query := "UPDATE tb_tujuan_pemda SET tujuan_pemda = ?, tematik_id = ?, periode_id = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, tujuanPemda.TujuanPemda, tujuanPemda.TematikId, tujuanPemda.PeriodeId, tujuanPemda.Id)
	if err != nil {
		return tujuanPemda, err
	}

	// Hapus semua indikator lama beserta targetnya
	scriptDeleteOldIndicators := "DELETE FROM tb_indikator WHERE tujuan_pemda_id = ?"
	_, err = tx.ExecContext(ctx, scriptDeleteOldIndicators, tujuanPemda.Id)
	if err != nil {
		return tujuanPemda, err
	}

	// Insert indikator baru
	for _, indikator := range tujuanPemda.Indikator {
		scriptInsertIndikator := `
            INSERT INTO tb_indikator 
                (id, tujuan_pemda_id, indikator, rumus_perhitungan, sumber_data) 
            VALUES 
                (?, ?, ?, ?, ?)`

		_, err := tx.ExecContext(ctx, scriptInsertIndikator,
			indikator.Id,
			tujuanPemda.Id,
			indikator.Indikator,
			indikator.RumusPerhitungan,
			indikator.SumberData)
		if err != nil {
			return tujuanPemda, err
		}

		// Hapus semua target lama untuk indikator ini
		scriptDeleteOldTargets := "DELETE FROM tb_target WHERE indikator_id = ?"
		_, err = tx.ExecContext(ctx, scriptDeleteOldTargets, indikator.Id)
		if err != nil {
			return tujuanPemda, err
		}

		// Insert target baru
		for _, target := range indikator.Target {
			scriptInsertTarget := `
                INSERT INTO tb_target 
                    (id, indikator_id, target, satuan, tahun)
                VALUES 
                    (?, ?, ?, ?, ?)`

			_, err := tx.ExecContext(ctx, scriptInsertTarget,
				target.Id,
				indikator.Id,
				target.Target,
				target.Satuan,
				target.Tahun)
			if err != nil {
				return tujuanPemda, err
			}
		}
	}

	return tujuanPemda, nil
}

func (repository *TujuanPemdaRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, tujuanPemdaId int) error {
	// 1. Hapus target berdasarkan indikator yang terkait dengan tujuan pemda
	queryDeleteTarget := `
        DELETE t FROM tb_target t
        INNER JOIN tb_indikator i ON t.indikator_id = i.id
        WHERE i.tujuan_pemda_id = ?`
	_, err := tx.ExecContext(ctx, queryDeleteTarget, tujuanPemdaId)
	if err != nil {
		return err
	}

	// 2. Hapus indikator yang terkait dengan tujuan pemda
	queryDeleteIndikator := "DELETE FROM tb_indikator WHERE tujuan_pemda_id = ?"
	_, err = tx.ExecContext(ctx, queryDeleteIndikator, tujuanPemdaId)
	if err != nil {
		return err
	}

	// 3. Hapus tujuan pemda
	queryDeleteTujuanPemda := "DELETE FROM tb_tujuan_pemda WHERE id = ?"
	_, err = tx.ExecContext(ctx, queryDeleteTujuanPemda, tujuanPemdaId)
	if err != nil {
		return err
	}

	return nil
}

func (repository *TujuanPemdaRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, tujuanPemdaId int) (domain.TujuanPemda, error) {
	query := `
        SELECT DISTINCT
            tp.id,
            tp.tujuan_pemda,
            tp.tematik_id,
            tp.periode_id,
            COALESCE(p.tahun_awal, '') as tahun_awal,
            COALESCE(p.tahun_akhir, '') as tahun_akhir,
            i.id as indikator_id,
            i.indikator as indikator_text,
            i.rumus_perhitungan,
            i.sumber_data,
            t.id as target_id,
            t.target as target_value,
            t.satuan as target_satuan,
            t.tahun as target_tahun
        FROM 
            tb_tujuan_pemda tp
            LEFT JOIN tb_periode p ON tp.periode_id = p.id
            LEFT JOIN tb_indikator i ON tp.id = i.tujuan_pemda_id
            LEFT JOIN tb_target t ON t.indikator_id = i.id
        WHERE tp.id = ?
        ORDER BY 
            tp.id, i.id, CAST(t.tahun AS SIGNED)`

	rows, err := tx.QueryContext(ctx, query, tujuanPemdaId)
	if err != nil {
		return domain.TujuanPemda{}, fmt.Errorf("error querying tujuan pemda: %v", err)
	}
	defer rows.Close()

	var result domain.TujuanPemda
	var firstRow = true
	indikatorMap := make(map[string]*domain.Indikator)

	for rows.Next() {
		var (
			id, tematikId, periodeId                         int
			tujuanPemdaText, tahunAwal, tahunAkhir           string
			indikatorId, indikatorText                       sql.NullString
			rumusPerhitunganNull, sumberDataNull             sql.NullString // Ubah ke sql.NullString
			targetId, targetValue, targetSatuan, targetTahun sql.NullString
		)

		err := rows.Scan(
			&id,
			&tujuanPemdaText,
			&tematikId,
			&periodeId,
			&tahunAwal,
			&tahunAkhir,
			&indikatorId,
			&indikatorText,
			&rumusPerhitunganNull, // Scan ke NullString
			&sumberDataNull,       // Scan ke NullString
			&targetId,
			&targetValue,
			&targetSatuan,
			&targetTahun,
		)
		if err != nil {
			return domain.TujuanPemda{}, fmt.Errorf("error scanning row: %v", err)
		}

		if firstRow {
			result = domain.TujuanPemda{
				Id:          id,
				TujuanPemda: tujuanPemdaText,
				TematikId:   tematikId,
				PeriodeId:   periodeId,
				Periode: domain.Periode{
					TahunAwal:  tahunAwal,
					TahunAkhir: tahunAkhir,
				},
				Indikator: []domain.Indikator{},
			}
			firstRow = false
		}

		if indikatorId.Valid && indikatorText.Valid {
			currentIndikator, exists := indikatorMap[indikatorId.String]

			if !exists {
				// Konversi NullString ke string biasa, jika null maka gunakan string kosong
				rumusPerhitungan := ""
				if rumusPerhitunganNull.Valid {
					rumusPerhitungan = rumusPerhitunganNull.String
				}

				sumberData := ""
				if sumberDataNull.Valid {
					sumberData = sumberDataNull.String
				}

				indikator := domain.Indikator{
					Id:               indikatorId.String,
					TujuanPemdaId:    id,
					Indikator:        indikatorText.String,
					RumusPerhitungan: sql.NullString{String: rumusPerhitungan, Valid: rumusPerhitungan != ""},
					SumberData:       sql.NullString{String: sumberData, Valid: sumberData != ""},
					Target:           []domain.Target{},
				}

				// Generate target untuk setiap tahun dalam periode
				if periodeId != 0 && tahunAwal != "" && tahunAkhir != "" {
					tahunAwalInt, errAwal := strconv.Atoi(tahunAwal)
					tahunAkhirInt, errAkhir := strconv.Atoi(tahunAkhir)

					if errAwal == nil && errAkhir == nil {
						for tahun := tahunAwalInt; tahun <= tahunAkhirInt; tahun++ {
							indikator.Target = append(indikator.Target, domain.Target{
								Id:          "-",
								IndikatorId: indikatorId.String,
								Target:      "-",
								Satuan:      "-",
								Tahun:       strconv.Itoa(tahun),
							})
						}
					}
				}

				result.Indikator = append(result.Indikator, indikator)
				indikatorMap[indikatorId.String] = &result.Indikator[len(result.Indikator)-1]
				currentIndikator = &result.Indikator[len(result.Indikator)-1]
			}

			// Update target yang ada dengan data sebenarnya
			if targetId.Valid && targetValue.Valid && targetTahun.Valid {
				for i := range currentIndikator.Target {
					if currentIndikator.Target[i].Tahun == targetTahun.String {
						currentIndikator.Target[i] = domain.Target{
							Id:          targetId.String,
							IndikatorId: indikatorId.String,
							Target:      targetValue.String,
							Satuan:      targetSatuan.String,
							Tahun:       targetTahun.String,
						}
						break
					}
				}
			}
		}
	}

	if err = rows.Err(); err != nil {
		return domain.TujuanPemda{}, fmt.Errorf("error iterating rows: %v", err)
	}

	if result.Id == 0 {
		return domain.TujuanPemda{}, fmt.Errorf("tujuan pemda dengan id %d tidak ditemukan", tujuanPemdaId)
	}

	return result, nil
}

func (repository *TujuanPemdaRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, tahun string) ([]domain.TujuanPemda, error) {
	query := `
        SELECT DISTINCT
            tp.id,
            tp.tujuan_pemda,
            tp.tematik_id,
            tp.periode_id,
            COALESCE(p.tahun_awal, '') as tahun_awal,
            COALESCE(p.tahun_akhir, '') as tahun_akhir,
            i.id as indikator_id,
            i.indikator as indikator_text,
			i.rumus_perhitungan,
			i.sumber_data,
            t.id as target_id,
            t.target as target_value,
            t.satuan as target_satuan,
            t.tahun as target_tahun
        FROM 
            tb_tujuan_pemda tp
            LEFT JOIN tb_periode p ON tp.periode_id = p.id
            LEFT JOIN tb_indikator i ON tp.id = i.tujuan_pemda_id
            LEFT JOIN tb_target t ON t.indikator_id = i.id
        ORDER BY 
            tp.id, i.id, CAST(t.tahun AS SIGNED)`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying tujuan pemda: %v", err)
	}
	defer rows.Close()

	tujuanPemdaMap := make(map[int]*domain.TujuanPemda)

	for rows.Next() {
		var (
			id, tematikId, periodeId                                             int
			tujuanPemdaText, rumusPerhitungan, sumberData, tahunAwal, tahunAkhir string
			indikatorId, indikatorText                                           sql.NullString
			targetId, targetValue, targetSatuan, targetTahun                     sql.NullString
		)

		err := rows.Scan(
			&id,
			&tujuanPemdaText,
			&tematikId,
			&periodeId,
			&tahunAwal,
			&tahunAkhir,
			&indikatorId,
			&indikatorText,
			&rumusPerhitungan,
			&sumberData,
			&targetId,
			&targetValue,
			&targetSatuan,
			&targetTahun,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		tujuanPemda, exists := tujuanPemdaMap[id]
		if !exists {
			tujuanPemda = &domain.TujuanPemda{
				Id:          id,
				TujuanPemda: tujuanPemdaText,
				TematikId:   tematikId,
				PeriodeId:   periodeId,
				Periode: domain.Periode{
					TahunAwal:  tahunAwal,
					TahunAkhir: tahunAkhir,
				},
				Indikator: []domain.Indikator{},
			}
			tujuanPemdaMap[id] = tujuanPemda
		}

		if indikatorId.Valid && indikatorText.Valid {
			var currentIndikator *domain.Indikator

			// Cari indikator yang sudah ada
			for i := range tujuanPemda.Indikator {
				if tujuanPemda.Indikator[i].Id == indikatorId.String {
					currentIndikator = &tujuanPemda.Indikator[i]
					break
				}
			}

			// Buat indikator baru jika belum ada
			if currentIndikator == nil {
				newIndikator := domain.Indikator{
					Id:            indikatorId.String,
					TujuanPemdaId: id,
					Indikator:     indikatorText.String,
					Target:        []domain.Target{},
				}

				// Generate target untuk setiap tahun dalam periode
				if periodeId != 0 && tahunAwal != "" && tahunAkhir != "" {
					tahunAwalInt, errAwal := strconv.Atoi(tahunAwal)
					tahunAkhirInt, errAkhir := strconv.Atoi(tahunAkhir)

					if errAwal == nil && errAkhir == nil {
						for tahun := tahunAwalInt; tahun <= tahunAkhirInt; tahun++ {
							newIndikator.Target = append(newIndikator.Target, domain.Target{
								Id:          "-",
								IndikatorId: indikatorId.String,
								Target:      "-",
								Satuan:      "-",
								Tahun:       strconv.Itoa(tahun),
							})
						}
					}
				}

				tujuanPemda.Indikator = append(tujuanPemda.Indikator, newIndikator)
				currentIndikator = &tujuanPemda.Indikator[len(tujuanPemda.Indikator)-1]
			}

			// Update target yang ada dengan data sebenarnya
			if targetId.Valid && targetValue.Valid && targetTahun.Valid {
				for i := range currentIndikator.Target {
					if currentIndikator.Target[i].Tahun == targetTahun.String {
						currentIndikator.Target[i] = domain.Target{
							Id:          targetId.String,
							IndikatorId: indikatorId.String,
							Target:      targetValue.String,
							Satuan:      targetSatuan.String,
							Tahun:       targetTahun.String,
						}
						break
					}
				}
			}
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	// Convert map to slice dan urutkan hasilnya
	result := make([]domain.TujuanPemda, 0, len(tujuanPemdaMap))
	for _, tujuanPemda := range tujuanPemdaMap {
		result = append(result, *tujuanPemda)
	}

	// Urutkan berdasarkan ID
	sort.Slice(result, func(i, j int) bool {
		return result[i].Id < result[j].Id
	})

	return result, nil
}

func (repository *TujuanPemdaRepositoryImpl) DeleteIndikator(ctx context.Context, tx *sql.Tx, tujuanPemdaId int) error {
	query := "DELETE FROM tb_indikator WHERE tujuan_pemda_id = ?"
	_, err := tx.ExecContext(ctx, query, tujuanPemdaId)
	return err
}

func (repository *TujuanPemdaRepositoryImpl) IsIdExists(ctx context.Context, tx *sql.Tx, id int) bool {
	query := "SELECT COUNT(*) FROM tb_tujuan_pemda WHERE id = ?"
	var count int
	err := tx.QueryRowContext(ctx, query, id).Scan(&count)
	if err != nil {
		return true
	}
	return count > 0
}

func (repository *TujuanPemdaRepositoryImpl) UpdatePeriode(ctx context.Context, tx *sql.Tx, tujuanPemda domain.TujuanPemda) (domain.TujuanPemda, error) {
	// Update hanya periode_id
	query := "UPDATE tb_tujuan_pemda SET periode_id = ? WHERE id = ?"
	result, err := tx.ExecContext(ctx, query, tujuanPemda.PeriodeId, tujuanPemda.Id)
	if err != nil {
		return domain.TujuanPemda{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.TujuanPemda{}, err
	}

	if rowsAffected == 0 {
		return domain.TujuanPemda{}, fmt.Errorf("periode dengan id %d sudah digunakan", tujuanPemda.PeriodeId)
	}

	// Ambil data terbaru setelah update
	query = `
        SELECT 
            tp.id,
            tp.tujuan_pemda,
            tp.tematik_id,
            tp.periode_id,
            COALESCE(p.tahun_awal, 'Pilih periode') as tahun_awal,
            COALESCE(p.tahun_akhir, 'Pilih periode') as tahun_akhir
        FROM 
            tb_tujuan_pemda tp
            LEFT JOIN tb_periode p ON tp.periode_id = p.id
        WHERE tp.id = ?`

	var updatedTujuanPemda domain.TujuanPemda
	err = tx.QueryRowContext(ctx, query, tujuanPemda.Id).Scan(
		&updatedTujuanPemda.Id,
		&updatedTujuanPemda.TujuanPemda,
		&updatedTujuanPemda.TematikId,
		&updatedTujuanPemda.PeriodeId,
		&updatedTujuanPemda.Periode.TahunAwal,
		&updatedTujuanPemda.Periode.TahunAkhir,
	)
	if err != nil {
		return domain.TujuanPemda{}, fmt.Errorf("gagal mengambil data setelah update: %v", err)
	}

	return updatedTujuanPemda, nil
}

func (repository *TujuanPemdaRepositoryImpl) FindAllWithPokin(ctx context.Context, tx *sql.Tx, tahun string) ([]domain.TujuanPemdaWithPokin, error) {
	query := `
    WITH periode_aktif AS (
        SELECT id, tahun_awal, tahun_akhir
        FROM tb_periode
        WHERE CAST(? AS SIGNED) BETWEEN CAST(tahun_awal AS SIGNED) AND CAST(tahun_akhir AS SIGNED)
    )
    SELECT 
        pk.id as pokin_id,
        pk.nama_pohon,
        pk.jenis_pohon,
        pk.level_pohon,
        pk.kode_opd,
        pk.keterangan,
        pk.tahun as tahun_pokin,
        tp.id as tujuan_id,
        tp.tujuan_pemda,
        COALESCE(pa.id, p.id) as periode_id,
        COALESCE(pa.tahun_awal, p.tahun_awal) as tahun_awal,
        COALESCE(pa.tahun_akhir, p.tahun_akhir) as tahun_akhir,
        i.id as indikator_id,
        i.indikator,
        i.rumus_perhitungan,
        i.sumber_data,
        t.id as target_id,
        t.target,
        t.satuan,
        t.tahun
    FROM 
        tb_pohon_kinerja pk
    LEFT JOIN 
        tb_tujuan_pemda tp ON pk.id = tp.tematik_id
    LEFT JOIN 
        tb_periode pa ON tp.periode_id = pa.id
    LEFT JOIN
        periode_aktif p ON CAST(pk.tahun AS SIGNED) BETWEEN CAST(p.tahun_awal AS SIGNED) AND CAST(p.tahun_akhir AS SIGNED)
    LEFT JOIN 
        tb_indikator i ON tp.id = i.tujuan_pemda_id
    LEFT JOIN 
        tb_target t ON i.id = t.indikator_id
    WHERE 
        pk.level_pohon = 0
        AND (
            (tp.id IS NOT NULL AND CAST(? AS SIGNED) BETWEEN CAST(pa.tahun_awal AS SIGNED) AND CAST(pa.tahun_akhir AS SIGNED))
            OR (tp.id IS NULL AND CAST(pk.tahun AS SIGNED) BETWEEN CAST(p.tahun_awal AS SIGNED) AND CAST(p.tahun_akhir AS SIGNED))
        )
    ORDER BY 
        pk.id, tp.id, i.id, t.tahun`

	rows, err := tx.QueryContext(ctx, query, tahun, tahun)
	if err != nil {
		return nil, fmt.Errorf("error querying data: %v", err)
	}
	defer rows.Close()

	pokinMap := make(map[int]*domain.TujuanPemdaWithPokin)
	indikatorMap := make(map[string]*domain.Indikator)
	targetMap := make(map[string]map[string]domain.Target)

	for rows.Next() {
		var (
			pokinId                                                int
			namaPohon, jenisPohon, kodeOpd, keterangan, tahunPokin string
			levelPohon                                             int
			tujuanId                                               sql.NullInt64
			tujuanPemda                                            sql.NullString
			periodeId                                              sql.NullInt64
			tahunAwal, tahunAkhir                                  sql.NullString
			indikatorId, indikatorText                             sql.NullString
			rumusPerhitungan, sumberData                           sql.NullString
			targetId, targetValue, targetSatuan, targetTahun       sql.NullString
		)

		err := rows.Scan(
			&pokinId,
			&namaPohon,
			&jenisPohon,
			&levelPohon,
			&kodeOpd,
			&keterangan,
			&tahunPokin,
			&tujuanId,
			&tujuanPemda,
			&periodeId,
			&tahunAwal,
			&tahunAkhir,
			&indikatorId,
			&indikatorText,
			&rumusPerhitungan,
			&sumberData,
			&targetId,
			&targetValue,
			&targetSatuan,
			&targetTahun,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Inisialisasi atau ambil data pokin dari map
		pokin, exists := pokinMap[pokinId]
		if !exists {
			pokin = &domain.TujuanPemdaWithPokin{
				PokinId:     pokinId,
				NamaPohon:   namaPohon,
				JenisPohon:  jenisPohon,
				LevelPohon:  levelPohon,
				KodeOpd:     kodeOpd,
				Keterangan:  keterangan,
				TahunPokin:  tahunPokin,
				TujuanPemda: nil,
			}
			pokinMap[pokinId] = pokin
		}

		// Jika ada data tujuan pemda
		if tujuanId.Valid && tujuanPemda.Valid {
			if pokin.TujuanPemda == nil {
				pokin.TujuanPemda = &domain.TujuanPemda{
					Id:          int(tujuanId.Int64),
					TujuanPemda: tujuanPemda.String,
					TematikId:   pokinId,
					PeriodeId:   int(periodeId.Int64),
					Periode: domain.Periode{
						Id:         int(periodeId.Int64),
						TahunAwal:  tahunAwal.String,
						TahunAkhir: tahunAkhir.String,
					},
					Indikator: []domain.Indikator{},
				}
			}

			// Proses indikator
			if indikatorId.Valid && indikatorText.Valid {
				if _, exists := indikatorMap[indikatorId.String]; !exists {
					// Konversi NullString ke string biasa
					rumusPerhitungan := ""
					if rumusPerhitungan != "" {
						_ = rumusPerhitungan
					}

					sumberData := ""
					if sumberData != "" {
						_ = sumberData
					}

					indikatorMap[indikatorId.String] = &domain.Indikator{
						Id:               indikatorId.String,
						TujuanPemdaId:    int(tujuanId.Int64),
						Indikator:        indikatorText.String,
						RumusPerhitungan: sql.NullString{String: rumusPerhitungan, Valid: rumusPerhitungan != ""},
						SumberData:       sql.NullString{String: sumberData, Valid: sumberData != ""},
						Target:           []domain.Target{},
					}
					targetMap[indikatorId.String] = make(map[string]domain.Target)
				}

				// Proses target
				if targetId.Valid && targetTahun.Valid {
					targetMap[indikatorId.String][targetTahun.String] = domain.Target{
						Id:          targetId.String,
						IndikatorId: indikatorId.String,
						Target:      targetValue.String,
						Satuan:      targetSatuan.String,
						Tahun:       targetTahun.String,
					}
				}
			}
		}
	}

	// Proses final untuk menambahkan target sesuai periode
	for pokinId, pokin := range pokinMap {
		if pokin.TujuanPemda != nil && pokin.TujuanPemda.Periode.TahunAwal != "" {
			tahunAwal, _ := strconv.Atoi(pokin.TujuanPemda.Periode.TahunAwal)
			tahunAkhir, _ := strconv.Atoi(pokin.TujuanPemda.Periode.TahunAkhir)

			for indikatorId, indikator := range indikatorMap {
				if indikator.TujuanPemdaId == pokin.TujuanPemda.Id {
					var targets []domain.Target
					for tahun := tahunAwal; tahun <= tahunAkhir; tahun++ {
						tahunStr := strconv.Itoa(tahun)
						if target, exists := targetMap[indikatorId][tahunStr]; exists {
							targets = append(targets, target)
						} else {
							// Tambahkan target kosong jika tidak ada data
							targets = append(targets, domain.Target{
								Id:          "-",
								IndikatorId: indikatorId,
								Target:      "-",
								Satuan:      "-",
								Tahun:       tahunStr,
							})
						}
					}

					// Sort targets berdasarkan tahun
					sort.Slice(targets, func(i, j int) bool {
						return targets[i].Tahun < targets[j].Tahun
					})

					indikator.Target = targets
					pokin.TujuanPemda.Indikator = append(pokin.TujuanPemda.Indikator, *indikator)
				}
			}
		}
		pokinMap[pokinId] = pokin
	}

	// Konversi map ke slice
	result := make([]domain.TujuanPemdaWithPokin, 0, len(pokinMap))
	for _, pokin := range pokinMap {
		// Sort indikator jika ada
		if pokin.TujuanPemda != nil {
			sort.Slice(pokin.TujuanPemda.Indikator, func(i, j int) bool {
				return pokin.TujuanPemda.Indikator[i].Id < pokin.TujuanPemda.Indikator[j].Id
			})
		}
		result = append(result, *pokin)
	}

	// Sort hasil akhir berdasarkan pokin_id
	sort.Slice(result, func(i, j int) bool {
		return result[i].PokinId < result[j].PokinId
	})

	return result, nil
}

func (repository *TujuanPemdaRepositoryImpl) IsPokinIdExists(ctx context.Context, tx *sql.Tx, pokinId int) (bool, error) {
	query := "SELECT COUNT(*) FROM tb_tujuan_pemda WHERE tematik_id = ?"
	var count int
	err := tx.QueryRowContext(ctx, query, pokinId).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
