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

type TujuanPemdaRepositoryImpl struct {
}

func NewTujuanPemdaRepositoryImpl() *TujuanPemdaRepositoryImpl {
	return &TujuanPemdaRepositoryImpl{}
}

func (repository *TujuanPemdaRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, tujuanPemda domain.TujuanPemda) (domain.TujuanPemda, error) {
	query := "INSERT INTO tb_tujuan_pemda(id, tujuan_pemda_id, rumus_perhitungan, sumber_data, periode_id) VALUES (?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, query, tujuanPemda.Id, tujuanPemda.TujuanPemdaId, tujuanPemda.RumusPerhitungan, tujuanPemda.SumberData, tujuanPemda.PeriodeId)
	if err != nil {
		return tujuanPemda, err
	}
	return tujuanPemda, nil
}

func (repository *TujuanPemdaRepositoryImpl) CreateIndikator(ctx context.Context, tx *sql.Tx, indikator domain.Indikator) (domain.Indikator, error) {
	query := "INSERT INTO tb_indikator(id, tujuan_pemda_id, indikator) VALUES (?, ?, ?)"
	_, err := tx.ExecContext(ctx, query, indikator.Id, indikator.TujuanPemdaId, indikator.Indikator)
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
	query := "UPDATE tb_tujuan_pemda SET rumus_perhitungan = ?, sumber_data = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, tujuanPemda.RumusPerhitungan, tujuanPemda.SumberData, tujuanPemda.Id)
	if err != nil {
		return tujuanPemda, err
	}

	// Proses indikator
	for _, indikator := range tujuanPemda.Indikator {
		// Update atau insert indikator
		scriptUpdateIndikator := `
            INSERT INTO tb_indikator (id, tujuan_pemda_id, indikator) 
            VALUES (?, ?, ?)
            ON DUPLICATE KEY UPDATE indikator = VALUES(indikator)`

		_, err := tx.ExecContext(ctx, scriptUpdateIndikator,
			indikator.Id,
			tujuanPemda.Id,
			indikator.Indikator)
		if err != nil {
			return tujuanPemda, err
		}

		// Hapus target lama untuk indikator ini
		scriptDeleteTargets := "DELETE FROM tb_target WHERE indikator_id = ?"
		_, err = tx.ExecContext(ctx, scriptDeleteTargets, indikator.Id)
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
            tp.tujuan_pemda_id,
            tp.rumus_perhitungan,
            tp.sumber_data,
            tp.periode_id,
            p.tahun_awal,
            p.tahun_akhir,
            i.id as indikator_id,
            i.indikator as indikator_text,
            t.id as target_id,
            t.target as target_value,
            t.satuan as target_satuan,
            t.tahun as target_tahun
        FROM 
            tb_tujuan_pemda tp
            INNER JOIN tb_periode p ON tp.periode_id = p.id
            LEFT JOIN tb_indikator i ON tp.id = i.tujuan_pemda_id
            LEFT JOIN tb_target t ON t.indikator_id = i.id
        WHERE tp.id = ?
        ORDER BY 
            tp.id, i.id, CAST(t.tahun AS SIGNED)`

	rows, err := tx.QueryContext(ctx, query, tujuanPemdaId)
	if err != nil {
		return domain.TujuanPemda{}, err
	}
	defer rows.Close()

	var tujuanPemda domain.TujuanPemda
	var firstRow = true
	indikatorMap := make(map[string]*domain.Indikator)

	for rows.Next() {
		var id int
		var tujuanPemdaId int
		var rumusPerhitungan, sumberData string
		var periodeId int
		var tahunAwal, tahunAkhir string
		var indikatorId, indikatorText sql.NullString
		var targetId, targetValue, targetSatuan, targetTahun sql.NullString

		err := rows.Scan(
			&id,
			&tujuanPemdaId,
			&rumusPerhitungan,
			&sumberData,
			&periodeId,
			&tahunAwal,
			&tahunAkhir,
			&indikatorId,
			&indikatorText,
			&targetId,
			&targetValue,
			&targetSatuan,
			&targetTahun,
		)
		if err != nil {
			return domain.TujuanPemda{}, err
		}

		if firstRow {
			tujuanPemda = domain.TujuanPemda{
				Id:               id,
				TujuanPemdaId:    tujuanPemdaId,
				PeriodeId:        periodeId,
				RumusPerhitungan: rumusPerhitungan,
				SumberData:       sumberData,
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
				// Buat slice target untuk semua tahun dalam periode
				tahunAwalInt, err := strconv.Atoi(tahunAwal)
				if err != nil {
					return domain.TujuanPemda{}, fmt.Errorf("gagal konversi tahun awal: %v", err)
				}

				tahunAkhirInt, err := strconv.Atoi(tahunAkhir)
				if err != nil {
					return domain.TujuanPemda{}, fmt.Errorf("gagal konversi tahun akhir: %v", err)
				}

				var targets []domain.Target
				// Inisialisasi target kosong untuk setiap tahun
				for tahun := tahunAwalInt; tahun <= tahunAkhirInt; tahun++ {
					tahunStr := strconv.Itoa(tahun)
					targets = append(targets, domain.Target{
						Id:          "-",
						IndikatorId: indikatorId.String,
						Target:      "-",
						Satuan:      "-",
						Tahun:       tahunStr,
					})
				}

				indikator := domain.Indikator{
					Id:            indikatorId.String,
					TujuanPemdaId: id,
					Indikator:     indikatorText.String,
					Target:        targets,
				}
				tujuanPemda.Indikator = append(tujuanPemda.Indikator, indikator)
				indikatorMap[indikatorId.String] = &tujuanPemda.Indikator[len(tujuanPemda.Indikator)-1]
				currentIndikator = &tujuanPemda.Indikator[len(tujuanPemda.Indikator)-1]
			}

			if targetId.Valid && targetValue.Valid && targetTahun.Valid {
				// Update target yang sudah ada dengan data sebenarnya
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

	if tujuanPemda.Id == 0 {
		return tujuanPemda, errors.New("tujuan pemda not found")
	}

	return tujuanPemda, nil
}

func (repository *TujuanPemdaRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, tahun string) ([]domain.TujuanPemda, error) {
	query := `
        SELECT DISTINCT
            tp.id,
			tp.tujuan_pemda_id,
            tp.rumus_perhitungan,
            tp.sumber_data,
            tp.periode_id,
            p.tahun_awal,
            p.tahun_akhir,
            i.id as indikator_id,
            i.indikator as indikator_text,
            t.id as target_id,
            t.target as target_value,
            t.satuan as target_satuan,
            t.tahun as target_tahun
        FROM 
            tb_tujuan_pemda tp
            INNER JOIN tb_periode p ON tp.periode_id = p.id
            LEFT JOIN tb_indikator i ON tp.id = i.tujuan_pemda_id
            LEFT JOIN tb_target t ON t.indikator_id = i.id
        WHERE 
            CASE 
                WHEN ? = '' THEN 1
                ELSE EXISTS (
                    SELECT 1 
                    FROM tb_target t2 
                    WHERE t2.indikator_id = i.id 
                    AND t2.tahun = ?
                )
            END
        ORDER BY 
            tp.id, i.id, CAST(t.tahun AS SIGNED)`

	rows, err := tx.QueryContext(ctx, query, tahun, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tujuanPemdaMap := make(map[int]*domain.TujuanPemda)

	for rows.Next() {
		var id int
		var tujuanPemdaId int
		var rumusPerhitungan, sumberData string
		var periodeId int
		var tahunAwal, tahunAkhir string
		var indikatorId, indikatorText sql.NullString
		var targetId, targetValue, targetSatuan, targetTahun sql.NullString

		err := rows.Scan(
			&id,
			&tujuanPemdaId,
			&rumusPerhitungan,
			&sumberData,
			&periodeId,
			&tahunAwal,
			&tahunAkhir,
			&indikatorId,
			&indikatorText,
			&targetId,
			&targetValue,
			&targetSatuan,
			&targetTahun,
		)
		if err != nil {
			return nil, err
		}

		tujuanPemda, exists := tujuanPemdaMap[id]
		if !exists {
			tahunAwalInt, _ := strconv.Atoi(tahunAwal)
			tahunAkhirInt, _ := strconv.Atoi(tahunAkhir)

			// Inisialisasi target untuk semua tahun dalam periode
			allTargets := make(map[string]bool)
			for tahun := tahunAwalInt; tahun <= tahunAkhirInt; tahun++ {
				allTargets[strconv.Itoa(tahun)] = false
			}

			tujuanPemda = &domain.TujuanPemda{
				Id:               id,
				TujuanPemdaId:    tujuanPemdaId,
				PeriodeId:        periodeId,
				RumusPerhitungan: rumusPerhitungan,
				SumberData:       sumberData,
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
					Target:        make([]domain.Target, 0),
				}
				tujuanPemda.Indikator = append(tujuanPemda.Indikator, newIndikator)
				currentIndikator = &tujuanPemda.Indikator[len(tujuanPemda.Indikator)-1]
			}

			// Update target jika ada
			if targetId.Valid && targetValue.Valid && targetTahun.Valid {
				target := domain.Target{
					Id:          targetId.String,
					IndikatorId: indikatorId.String,
					Target:      targetValue.String,
					Satuan:      targetSatuan.String,
					Tahun:       targetTahun.String,
				}

				// Cek apakah target untuk tahun ini sudah ada
				targetExists := false
				for i, existingTarget := range currentIndikator.Target {
					if existingTarget.Tahun == targetTahun.String {
						currentIndikator.Target[i] = target
						targetExists = true
						break
					}
				}

				// Jika target belum ada, tambahkan baru
				if !targetExists {
					currentIndikator.Target = append(currentIndikator.Target, target)
				}
			}
		}
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
