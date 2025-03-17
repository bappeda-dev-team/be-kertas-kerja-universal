package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
	"sort"
)

type LaporanRepositoryImpl struct{}

func NewLaporanRepositoryImpl() *LaporanRepositoryImpl {
	return &LaporanRepositoryImpl{}
}

// THIS IS FOR SQL
func (repo *LaporanRepositoryImpl) OpdSupportingPokin(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.OpdSupportingPokin, error) {
	script := `
		WITH RECURSIVE pohon_hierarchy AS (
			SELECT
				pk.id,
				pk.parent,
				pk.nama_pohon,
				pk.jenis_pohon,
				pk.level_pohon,
				pk.kode_opd,
				pk.keterangan,
				pk.tahun,
				pk.status
			FROM tb_pohon_kinerja pk
			WHERE
				pk.kode_opd = ?
				AND pk.tahun = ?
				AND pk.jenis_pohon = 'Strategic Pemda'
				AND pk.parent != 0

			UNION ALL

			SELECT
				pk.id,
				pk.parent,
				pk.nama_pohon,
				pk.jenis_pohon,
				pk.level_pohon,
				pk.kode_opd,
				pk.keterangan,
				pk.tahun,
				pk.status
			FROM tb_pohon_kinerja pk
				INNER JOIN pohon_hierarchy ph ON pk.id = ph.parent
		)
		SELECT
			ph.id,
            ph.nama_pohon,
            ph.parent,
            ph.jenis_pohon,
            ph.level_pohon,
            ph.kode_opd,
            ph.keterangan,
            ph.tahun,
            i.id as indikator_id,
            i.indikator as nama_indikator,
            t.id as target_id,
            t.target as target_value,
            t.satuan as target_satuan
		FROM
            pohon_hierarchy ph
		LEFT JOIN
            tb_indikator i ON ph.id = i.pokin_id
		LEFT JOIN
            tb_target t ON i.id = t.indikator_id
        WHERE
            i.tahun = ? AND t.tahun = ?
		ORDER BY ph.id ASC;`

	rows, err := tx.QueryContext(ctx, script, kodeOpd, tahun, tahun, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map untuk menyimpan pohon kinerja yang sudah diproses
	pokinMap := make(map[int]domain.OpdSupportingPokin)
	indikatorMap := make(map[string]domain.Indikator)
	for rows.Next() {
		var (
			pokinId, parent, levelPohon                            int
			namaPohon, jenisPohon, kodeOpd, keterangan, tahunPokin string
			indikatorId, namaIndikator                             sql.NullString
			targetId, targetValue, targetSatuan                    sql.NullString
		)
		err := rows.Scan(
			&pokinId, &namaPohon, &parent, &jenisPohon, &levelPohon,
			&kodeOpd, &keterangan, &tahunPokin,
			&indikatorId, &namaIndikator,
			&targetId, &targetValue, &targetSatuan,
		)
		if err != nil {
			return nil, err
		}
		// Proses Pohon Kinerja
		pokin, exists := pokinMap[pokinId]
		if !exists {
			pokin = domain.OpdSupportingPokin{
				Id:         pokinId,
				NamaPohon:  namaPohon,
				Parent:     parent,
				JenisPohon: jenisPohon,
				LevelPohon: levelPohon,
				KodeOpd:    kodeOpd,
				Keterangan: keterangan,
				Tahun:      tahunPokin,
			}
			pokinMap[pokinId] = pokin
		}

		// Proses Indikator jika ada
		if indikatorId.Valid && namaIndikator.Valid {
			indikator, exists := indikatorMap[indikatorId.String]
			if !exists {
				indikator = domain.Indikator{
					Id:        indikatorId.String,
					PokinId:   fmt.Sprint(pokinId),
					Indikator: namaIndikator.String,
					Tahun:     tahunPokin,
				}
			}

			// Proses Target jika ada
			if targetId.Valid && targetValue.Valid && targetSatuan.Valid {
				target := domain.Target{
					Id:          targetId.String,
					IndikatorId: indikatorId.String,
					Target:      targetValue.String,
					Satuan:      targetSatuan.String,
					Tahun:       tahunPokin,
				}
				indikator.Target = append(indikator.Target, target)
			}

			indikatorMap[indikatorId.String] = indikator

			// Update indikator di pokin
			pokin.Indikators = append(pokin.Indikators, indikator)
			pokinMap[pokinId] = pokin
		}
	}
	// Konversi map ke slice
	var result []domain.OpdSupportingPokin
	for _, pokin := range pokinMap {
		result = append(result, pokin)
	}

	// Urutkan berdasarkan level dan ID
	sort.Slice(result, func(i, j int) bool {
		if result[i].LevelPohon == result[j].LevelPohon {
			return result[i].Id < result[j].Id
		}
		return result[i].LevelPohon < result[j].LevelPohon
	})

	return result, nil

}
