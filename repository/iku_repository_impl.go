package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"sort"
)

type IkuRepositoryImpl struct {
}

func NewIkuRepositoryImpl() *IkuRepositoryImpl {
	return &IkuRepositoryImpl{}
}
func (repository *IkuRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, tahun string) ([]domain.Indikator, error) {
	query := `
    WITH indikator_tujuan AS (
        -- Indikator dan target dari tujuan pemda
        SELECT 
            i.id as indikator_id,
            i.indikator,
            i.created_at as indikator_created_at,
            t.id as target_id,
            t.target,
            t.satuan,
            'tujuan_pemda' as sumber,
            tp.id as parent_id,
            tp.tujuan_pemda as parent_name
        FROM tb_indikator i
        INNER JOIN tb_tujuan_pemda tp ON i.tujuan_pemda_id = tp.id
        INNER JOIN tb_target t ON t.indikator_id = i.id
        INNER JOIN tb_periode p ON tp.periode_id = p.id
        WHERE t.tahun = ?  
        AND ? BETWEEN p.tahun_awal AND p.tahun_akhir  
    ),
    indikator_sasaran AS (
        -- Indikator dan target dari sasaran pemda dan pohon kinerja
        SELECT 
            COALESCE(i_sp.id, i_pk.id) as indikator_id,
            COALESCE(i_sp.indikator, i_pk.indikator) as indikator,
            COALESCE(i_sp.created_at, i_pk.created_at) as indikator_created_at,
            COALESCE(t_sp.id, t_pk.id) as target_id,
            COALESCE(t_sp.target, t_pk.target) as target,
            COALESCE(t_sp.satuan, t_pk.satuan) as satuan,
            'sasaran_pemda' as sumber,
            sp.id as parent_id,
            sp.sasaran_pemda as parent_name
        FROM tb_sasaran_pemda sp
        LEFT JOIN tb_indikator i_sp ON sp.id = i_sp.sasaran_pemda_id
        LEFT JOIN tb_target t_sp ON i_sp.id = t_sp.indikator_id
        LEFT JOIN tb_pohon_kinerja pk ON sp.subtema_id = pk.id
        LEFT JOIN tb_indikator i_pk ON pk.id = i_pk.pokin_id
        LEFT JOIN tb_target t_pk ON i_pk.id = t_pk.indikator_id
        WHERE pk.tahun = ?
    )
    -- Gabungkan kedua sumber data
    SELECT 
        indikator_id, 
        indikator,
        indikator_created_at, 
        target_id, 
        target, 
        satuan, 
        sumber,
        parent_id,
        parent_name
    FROM (
        SELECT * FROM indikator_tujuan
        UNION ALL
        SELECT * FROM indikator_sasaran
    ) combined
    WHERE indikator IS NOT NULL
    ORDER BY indikator_created_at ASC`

	rows, err := tx.QueryContext(ctx, query, tahun, tahun, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	indikatorMap := make(map[string]*domain.Indikator)

	for rows.Next() {
		var (
			indikatorId        sql.NullString
			indikator          sql.NullString
			indikatorCreatedAt sql.NullTime
			targetId           sql.NullString
			target             sql.NullString
			satuan             sql.NullString
			sumber             string
			parentId           sql.NullInt64
			parentName         sql.NullString
		)

		err := rows.Scan(
			&indikatorId,
			&indikator,
			&indikatorCreatedAt,
			&targetId,
			&target,
			&satuan,
			&sumber,
			&parentId,
			&parentName,
		)
		if err != nil {
			return nil, err
		}

		if !indikator.Valid || !indikatorId.Valid {
			continue
		}

		item, exists := indikatorMap[indikatorId.String]
		if !exists {
			item = &domain.Indikator{
				Id:         indikatorId.String,
				Indikator:  indikator.String,
				CreatedAt:  indikatorCreatedAt.Time,
				Sumber:     sumber,
				ParentId:   int(parentId.Int64),
				ParentName: parentName.String,
				Target:     []domain.Target{},
			}
			indikatorMap[indikatorId.String] = item
		}

		if target.Valid && targetId.Valid {
			item.Target = append(item.Target, domain.Target{
				Id:     targetId.String,
				Target: target.String,
				Satuan: satuan.String,
			})
		}
	}

	// Konversi map ke slice dan urutkan berdasarkan created_at dan indikator
	result := make([]domain.Indikator, 0, len(indikatorMap))
	for _, item := range indikatorMap {
		result = append(result, *item)
	}

	// Urutkan slice berdasarkan created_at dan indikator
	sort.Slice(result, func(i, j int) bool {
		if result[i].CreatedAt.Equal(result[j].CreatedAt) {
			return result[i].Indikator < result[j].Indikator
		}
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})

	return result, nil
}
