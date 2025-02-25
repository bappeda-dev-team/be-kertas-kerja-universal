package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"sort"
	"strconv"
)

type IkuRepositoryImpl struct {
}

func NewIkuRepositoryImpl() *IkuRepositoryImpl {
	return &IkuRepositoryImpl{}
}

func (repository *IkuRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]domain.Indikator, error) {
	query := `
    WITH indikator_tujuan AS (
        SELECT 
            i.id as indikator_id,
            i.indikator,
            i.rumus_perhitungan,
            i.sumber_data,
            i.created_at as indikator_created_at,
            t.id as target_id,
            t.target,
            t.satuan,
            t.tahun as target_tahun,
            'tujuan_pemda' as sumber,
            tp.id as parent_id,
            tp.tujuan_pemda as parent_name,
            tp.tahun_awal_periode,
            tp.tahun_akhir_periode,
            tp.jenis_periode
        FROM tb_indikator i
        INNER JOIN tb_tujuan_pemda tp ON i.tujuan_pemda_id = tp.id
        LEFT JOIN tb_target t ON t.indikator_id = i.id
        WHERE tp.tahun_awal_periode = ? 
        AND tp.tahun_akhir_periode = ?
        AND tp.jenis_periode = ?
    ),
    indikator_sasaran AS (
        SELECT 
            i.id as indikator_id,
            i.indikator,
            i.rumus_perhitungan,
            i.sumber_data,
            i.created_at as indikator_created_at,
            t.id as target_id,
            t.target,
            t.satuan,
            t.tahun as target_tahun,
            'sasaran_pemda' as sumber,
            sp.id as parent_id,
            sp.sasaran_pemda as parent_name,
            sp.tahun_awal,
            sp.tahun_akhir,
            sp.jenis_periode
        FROM tb_indikator i
        INNER JOIN tb_sasaran_pemda sp ON i.sasaran_pemda_id = sp.id
        LEFT JOIN tb_target t ON t.indikator_id = i.id
        WHERE sp.tahun_awal = ? 
        AND sp.tahun_akhir = ?
        AND sp.jenis_periode = ?
    )
    SELECT * FROM (
        SELECT * FROM indikator_tujuan
        UNION ALL
        SELECT * FROM indikator_sasaran
    ) combined
    WHERE indikator IS NOT NULL
    ORDER BY indikator_created_at ASC`

	rows, err := tx.QueryContext(ctx, query,
		tahunAwal, tahunAkhir, jenisPeriode,
		tahunAwal, tahunAkhir, jenisPeriode,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	indikatorMap := make(map[string]*domain.Indikator)

	for rows.Next() {
		var (
			indikatorId        sql.NullString
			indikator          sql.NullString
			rumusPerhitungan   sql.NullString
			sumberData         sql.NullString
			indikatorCreatedAt sql.NullTime
			targetId           sql.NullString
			target             sql.NullString
			satuan             sql.NullString
			targetTahun        sql.NullString
			sumber             string
			parentId           sql.NullInt64
			parentName         sql.NullString
			tahunAwal          string
			tahunAkhir         string
			// periodeId          int
			jenisPeriodeData string
		)

		err := rows.Scan(
			&indikatorId,
			&indikator,
			&rumusPerhitungan,
			&sumberData,
			&indikatorCreatedAt,
			&targetId,
			&target,
			&satuan,
			&targetTahun,
			&sumber,
			&parentId,
			&parentName,
			&tahunAwal,
			&tahunAkhir,
			// &periodeId,
			&jenisPeriodeData,
		)
		if err != nil {
			return nil, err
		}
		if !indikator.Valid || !indikatorId.Valid {
			continue
		}

		item, exists := indikatorMap[indikatorId.String]
		if !exists {
			// Inisialisasi target kosong untuk semua tahun dalam periode
			tahunAwalInt, _ := strconv.Atoi(tahunAwal)
			tahunAkhirInt, _ := strconv.Atoi(tahunAkhir)

			var targets []domain.Target
			for tahun := tahunAwalInt; tahun <= tahunAkhirInt; tahun++ {
				tahunStr := strconv.Itoa(tahun)
				targets = append(targets, domain.Target{
					Id:          "-",
					IndikatorId: indikatorId.String,
					Target:      "",
					Satuan:      "",
					Tahun:       tahunStr,
				})
			}

			item = &domain.Indikator{
				Id:               indikatorId.String,
				Indikator:        indikator.String,
				RumusPerhitungan: rumusPerhitungan,
				SumberData:       sumberData,
				CreatedAt:        indikatorCreatedAt.Time,
				Sumber:           sumber,
				ParentId:         int(parentId.Int64),
				ParentName:       parentName.String,
				Target:           targets,
			}
			indikatorMap[indikatorId.String] = item
		}

		// Update target yang memiliki data
		if targetId.Valid && targetTahun.Valid {
			tahunInt, _ := strconv.Atoi(targetTahun.String)
			tahunAwalInt, _ := strconv.Atoi(tahunAwal)
			if tahunInt >= tahunAwalInt {
				idx := tahunInt - tahunAwalInt
				if idx >= 0 && idx < len(item.Target) {
					item.Target[idx] = domain.Target{
						Id:          targetId.String,
						IndikatorId: indikatorId.String,
						Target:      target.String,
						Satuan:      satuan.String,
						Tahun:       targetTahun.String,
					}
				}
			}
		}
	}

	result := make([]domain.Indikator, 0, len(indikatorMap))
	for _, item := range indikatorMap {
		result = append(result, *item)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].CreatedAt.Equal(result[j].CreatedAt) {
			return result[i].Indikator < result[j].Indikator
		}
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})

	return result, nil
}
