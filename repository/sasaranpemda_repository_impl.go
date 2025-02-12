package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
	"sort"
	"strconv"
)

type SasaranPemdaRepositoryImpl struct {
}

func NewSasaranPemdaRepositoryImpl() *SasaranPemdaRepositoryImpl {
	return &SasaranPemdaRepositoryImpl{}
}

func (repository *SasaranPemdaRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, sasaranPemda domain.SasaranPemda) (domain.SasaranPemda, error) {
	query := "INSERT INTO tb_sasaran_pemda(id, tujuan_pemda_id, subtema_id, sasaran_pemda) VALUES (?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, query, sasaranPemda.Id, sasaranPemda.TujuanPemdaId, sasaranPemda.SubtemaId, sasaranPemda.SasaranPemda)
	if err != nil {
		return sasaranPemda, err
	}
	return sasaranPemda, nil
}

func (repository *SasaranPemdaRepositoryImpl) CreateIndikator(ctx context.Context, tx *sql.Tx, indikator domain.Indikator) (domain.Indikator, error) {
	query := "INSERT INTO tb_indikator(id, sasaran_pemda_id, indikator) VALUES (?, ?, ?)"
	_, err := tx.ExecContext(ctx, query, indikator.Id, indikator.SasaranPemdaId, indikator.Indikator)
	if err != nil {
		return indikator, err
	}
	return indikator, nil
}

func (repository *SasaranPemdaRepositoryImpl) CreateTarget(ctx context.Context, tx *sql.Tx, target domain.Target) error {
	query := "INSERT INTO tb_target(id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, query, target.Id, target.IndikatorId, target.Target, target.Satuan, target.Tahun)
	return err
}

func (repository *SasaranPemdaRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, sasaranPemda domain.SasaranPemda) (domain.SasaranPemda, error) {
	// Update tujuan pemda
	query := "UPDATE tb_sasaran_pemda SET tujuan_pemda_id = ?, sasaran_pemda = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, sasaranPemda.TujuanPemdaId, sasaranPemda.SasaranPemda, sasaranPemda.Id)
	if err != nil {
		return sasaranPemda, err
	}

	return sasaranPemda, nil
}

func (repository *SasaranPemdaRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, sasaranPemdaId int) error {
	query := "DELETE FROM tb_sasaran_pemda WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, sasaranPemdaId)
	if err != nil {
		return err
	}

	return nil
}

func (repository *SasaranPemdaRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, sasaranPemdaId int) (domain.SasaranPemda, error) {

	query := `
			SELECT 
				sp.id,
				sp.subtema_id,
				sp.tujuan_pemda_id,
				sp.sasaran_pemda
			FROM 
				tb_sasaran_pemda sp
			WHERE 
				sp.id = ?`

	var sasaranPemda domain.SasaranPemda
	err := tx.QueryRowContext(ctx, query, sasaranPemdaId).Scan(
		&sasaranPemda.Id,
		&sasaranPemda.SubtemaId,
		&sasaranPemda.TujuanPemdaId,
		&sasaranPemda.SasaranPemda,
	)
	if err != nil {
		return domain.SasaranPemda{}, err
	}

	return sasaranPemda, nil

}

func (repository *SasaranPemdaRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, tahun string) ([]domain.SasaranPemda, error) {
	query := `
        SELECT DISTINCT
            tp.id,
            tp.subtema_id,
            tp.sasaran_pemda,
            tp.periode_id,
             COALESCE(p.tahun_awal, '') as tahun_awal,
            COALESCE(p.tahun_akhir, '') as tahun_akhir,
            i.id as indikator_id,
            i.indikator as indikator_text,
            t.id as target_id,
            t.target as target_value,
            t.satuan as target_satuan,
            t.tahun as target_tahun
        FROM 
            tb_sasaran_pemda tp
            LEFT JOIN tb_periode p ON tp.periode_id = p.id
            LEFT JOIN tb_indikator i ON tp.id = i.sasaran_pemda_id
            LEFT JOIN tb_target t ON t.indikator_id = i.id
        ORDER BY 
            tp.id, i.id, CAST(t.tahun AS SIGNED)`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sasaranPemdaMap := make(map[int]*domain.SasaranPemda)

	for rows.Next() {
		var id int
		var subtemaId int
		var namaSubtema string
		var periodeId int
		var tahunAwal, tahunAkhir string
		var indikatorId, indikatorText sql.NullString
		var targetId, targetValue, targetSatuan, targetTahun sql.NullString

		err := rows.Scan(
			&id,
			&subtemaId,
			&namaSubtema,
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

		sasaranPemda, exists := sasaranPemdaMap[id]
		if !exists {
			sasaranPemda = &domain.SasaranPemda{
				Id:          id,
				SubtemaId:   subtemaId,
				NamaSubtema: namaSubtema,
				Periode: domain.Periode{
					TahunAwal:  tahunAwal,
					TahunAkhir: tahunAkhir,
				},
				Indikator: []domain.Indikator{},
			}

			// Update periode hanya jika ada dan valid
			if periodeId != 0 && tahunAwal != "" && tahunAkhir != "" {
				sasaranPemda.Periode.TahunAwal = tahunAwal
				sasaranPemda.Periode.TahunAkhir = tahunAkhir
			}

			sasaranPemdaMap[id] = sasaranPemda
		}

		if indikatorId.Valid && indikatorText.Valid {
			var currentIndikator *domain.Indikator

			// Cari indikator yang sudah ada
			for i := range sasaranPemda.Indikator {
				if sasaranPemda.Indikator[i].Id == indikatorId.String {
					currentIndikator = &sasaranPemda.Indikator[i]
					break
				}
			}

			// Buat indikator baru jika belum ada
			if currentIndikator == nil {
				newIndikator := domain.Indikator{
					Id:             indikatorId.String,
					SasaranPemdaId: id,
					Indikator:      indikatorText.String,
					Target:         make([]domain.Target, 0),
				}

				// Tambahkan target kosong untuk semua tahun dalam periode
				if periodeId != 0 && tahunAwal != "" && tahunAkhir != "" {
					tahunAwalInt, _ := strconv.Atoi(tahunAwal)
					tahunAkhirInt, _ := strconv.Atoi(tahunAkhir)

					for tahun := tahunAwalInt; tahun <= tahunAkhirInt; tahun++ {
						tahunStr := strconv.Itoa(tahun)
						newIndikator.Target = append(newIndikator.Target, domain.Target{
							Id:          "",
							IndikatorId: indikatorId.String,
							Target:      "-",
							Satuan:      "-",
							Tahun:       tahunStr,
						})
					}
				}

				sasaranPemda.Indikator = append(sasaranPemda.Indikator, newIndikator)
				currentIndikator = &sasaranPemda.Indikator[len(sasaranPemda.Indikator)-1]
			}

			// Update target dengan data sebenarnya jika ada
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

	// Convert map to slice dan urutkan hasilnya
	result := make([]domain.SasaranPemda, 0, len(sasaranPemdaMap))
	for _, sasaranPemda := range sasaranPemdaMap {
		result = append(result, *sasaranPemda)
	}

	// Urutkan berdasarkan ID
	sort.Slice(result, func(i, j int) bool {
		return result[i].Id < result[j].Id
	})

	return result, nil
}

func (repository *SasaranPemdaRepositoryImpl) DeleteIndikator(ctx context.Context, tx *sql.Tx, sasaranPemdaId int) error {
	query := "DELETE FROM tb_indikator WHERE sasaran_pemda_id = ?"
	_, err := tx.ExecContext(ctx, query, sasaranPemdaId)
	return err
}

func (repository *SasaranPemdaRepositoryImpl) IsIdExists(ctx context.Context, tx *sql.Tx, id int) bool {
	query := "SELECT COUNT(*) FROM tb_sasaran_pemda WHERE id = ?"
	var count int
	err := tx.QueryRowContext(ctx, query, id).Scan(&count)
	if err != nil {
		return true
	}
	return count > 0
}

func (repository *SasaranPemdaRepositoryImpl) UpdatePeriode(ctx context.Context, tx *sql.Tx, sasaranPemda domain.SasaranPemda) (domain.SasaranPemda, error) {
	// Update hanya periode_id
	query := "UPDATE tb_sasaran_pemda SET periode_id = ? WHERE id = ?"
	result, err := tx.ExecContext(ctx, query, sasaranPemda.PeriodeId, sasaranPemda.Id)
	if err != nil {
		return domain.SasaranPemda{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.SasaranPemda{}, err
	}

	if rowsAffected == 0 {
		return domain.SasaranPemda{}, fmt.Errorf("periode dengan id %d sudah digunakan", sasaranPemda.PeriodeId)
	}

	// Ambil data terbaru setelah update
	query = `
        SELECT 
            tp.id,
            tp.sasaran_pemda_id,
            tp.periode_id,
            COALESCE(p.tahun_awal, 'Pilih periode') as tahun_awal,
            COALESCE(p.tahun_akhir, 'Pilih periode') as tahun_akhir
        FROM 
            tb_sasaran_pemda tp
            LEFT JOIN tb_periode p ON tp.periode_id = p.id
        WHERE tp.id = ?`

	var updatedSasaranPemda domain.SasaranPemda
	err = tx.QueryRowContext(ctx, query, sasaranPemda.Id).Scan(
		&updatedSasaranPemda.Id,
		&updatedSasaranPemda.SubtemaId,
		&updatedSasaranPemda.SasaranPemda,
		&updatedSasaranPemda.PeriodeId,
		&updatedSasaranPemda.Periode.TahunAwal,
		&updatedSasaranPemda.Periode.TahunAkhir,
	)
	if err != nil {
		return domain.SasaranPemda{}, fmt.Errorf("gagal mengambil data setelah update: %v", err)
	}

	return updatedSasaranPemda, nil
}

func (repository *SasaranPemdaRepositoryImpl) FindAllWithPokin(ctx context.Context, tx *sql.Tx, tahun string) ([]domain.SasaranPemdaWithPokin, error) {
	checkQuery := `
	SELECT COUNT(*) 
	FROM tb_pohon_kinerja 
	WHERE level_pohon = 0 AND tahun = ?
`
	var count int
	err := tx.QueryRowContext(ctx, checkQuery, tahun).Scan(&count)
	if err != nil {
		return nil, err
	}

	// Jika tidak ada tematik level 0, return array kosong
	if count == 0 {
		return []domain.SasaranPemdaWithPokin{}, nil
	}
	query := `
    WITH tematik_tanpa_turunan AS (
        SELECT 
            t.id as subtematik_id,
            t.nama_pohon as nama_subtematik,
            t.jenis_pohon,
            t.level_pohon,
            t.keterangan,
            t.id as tematik_id,
            t.nama_pohon as nama_tematik,
            0 as id_sasaran_pemda,
            '' as sasaran_pemda,
            NULL as indikator,
            NULL as target,
            NULL as satuan,
            t.created_at as created_at
        FROM tb_pohon_kinerja t
        LEFT JOIN tb_pohon_kinerja child ON child.parent = t.id AND child.tahun = ?
        WHERE t.level_pohon = 0 
        AND t.tahun = ?
        AND child.id IS NULL
    )
    SELECT * FROM (
        SELECT 
            pk.id as subtematik_id,
            pk.nama_pohon as nama_subtematik,
            pk.jenis_pohon,
            pk.level_pohon,
            pk.keterangan,
            COALESCE(
                CASE 
                    WHEN pk.level_pohon = 1 THEN parent.id
                    WHEN pk.level_pohon = 2 THEN grandparent.id
                    WHEN pk.level_pohon = 3 THEN great_grandparent.id
                END, 
                0
            ) as tematik_id,
            COALESCE(
                CASE 
                    WHEN pk.level_pohon = 1 THEN parent.nama_pohon
                    WHEN pk.level_pohon = 2 THEN grandparent.nama_pohon
                    WHEN pk.level_pohon = 3 THEN great_grandparent.nama_pohon
                END, 
                ''
            ) as nama_tematik,
            COALESCE(sp.id, 0) as id_sasaran_pemda,
            COALESCE(sp.sasaran_pemda, '') as sasaran_pemda,
            COALESCE(i_sp.indikator, i_pk.indikator) as indikator,
            COALESCE(t_sp.target, t_pk.target) as target,
            COALESCE(t_sp.satuan, t_pk.satuan) as satuan,
            pk.created_at as created_at
        FROM 
            tb_pohon_kinerja pk
        LEFT JOIN 
            tb_pohon_kinerja parent ON pk.parent = parent.id AND parent.tahun = ?
        LEFT JOIN 
            tb_pohon_kinerja grandparent ON parent.parent = grandparent.id AND grandparent.tahun = ?
        LEFT JOIN 
            tb_pohon_kinerja great_grandparent ON grandparent.parent = great_grandparent.id AND great_grandparent.tahun = ?
        LEFT JOIN 
            tb_sasaran_pemda sp ON pk.id = sp.subtema_id
        LEFT JOIN 
            tb_indikator i_sp ON sp.id = i_sp.sasaran_pemda_id
        LEFT JOIN 
            tb_target t_sp ON i_sp.id = t_sp.indikator_id
        LEFT JOIN 
            tb_indikator i_pk ON pk.id = i_pk.pokin_id AND i_sp.id IS NULL
        LEFT JOIN 
            tb_target t_pk ON i_pk.id = t_pk.indikator_id AND t_sp.id IS NULL
        WHERE 
            pk.level_pohon BETWEEN 1 AND 3
            AND pk.tahun = ?

        UNION ALL

        SELECT * FROM tematik_tanpa_turunan
    ) combined
    ORDER BY 
        nama_tematik ASC,
        CASE 
            WHEN level_pohon = 0 THEN 2  -- Tematik tanpa turunan
            WHEN level_pohon = 1 THEN 0  -- Subtematik prioritas pertama
            WHEN level_pohon = 2 THEN 1  -- Sub-subtematik prioritas kedua
            WHEN level_pohon = 3 THEN 1  -- Sub-sub-subtematik prioritas ketiga
        END ASC,
        level_pohon ASC,
        created_at ASC,
        nama_subtematik ASC`

	rows, err := tx.QueryContext(ctx, query, tahun, tahun, tahun, tahun, tahun, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subtematikMap := make(map[int]*domain.SasaranPemdaWithPokin)

	for rows.Next() {
		var (
			subtematikId   int
			namaSubtematik string
			jenisPohon     string
			levelPohon     int
			keterangan     sql.NullString
			tematikId      int
			namaTematik    string
			idSasaranPemda int
			sasaranPemda   string
			indikator      sql.NullString
			target         sql.NullString
			satuan         sql.NullString
			createdAt      sql.NullTime
		)

		err := rows.Scan(
			&subtematikId,
			&namaSubtematik,
			&jenisPohon,
			&levelPohon,
			&keterangan,
			&tematikId,
			&namaTematik,
			&idSasaranPemda,
			&sasaranPemda,
			&indikator,
			&target,
			&satuan,
			&createdAt,
		)
		if err != nil {
			return nil, err
		}

		item, exists := subtematikMap[subtematikId]
		if !exists {
			item = &domain.SasaranPemdaWithPokin{
				SubtematikId:        subtematikId,
				NamaSubtematik:      namaSubtematik,
				JenisPohon:          jenisPohon,
				LevelPohon:          levelPohon,
				Keterangan:          keterangan.String,
				TematikId:           tematikId,
				NamaTematik:         namaTematik,
				IdsasaranPemda:      idSasaranPemda,
				SasaranPemda:        sasaranPemda,
				IndikatorSubtematik: []domain.Indikator{},
			}
			subtematikMap[subtematikId] = item
		}

		if indikator.Valid {
			var indikatorFound bool
			for i := range item.IndikatorSubtematik {
				if item.IndikatorSubtematik[i].Indikator == indikator.String {
					if target.Valid {
						item.IndikatorSubtematik[i].Target = append(
							item.IndikatorSubtematik[i].Target,
							domain.Target{
								Target: target.String,
								Satuan: satuan.String,
							},
						)
					}
					indikatorFound = true
					break
				}
			}

			if !indikatorFound {
				newIndikator := domain.Indikator{
					Indikator: indikator.String,
					Target:    []domain.Target{},
				}
				if target.Valid {
					newIndikator.Target = append(newIndikator.Target, domain.Target{
						Target: target.String,
						Satuan: satuan.String,
					})
				}
				item.IndikatorSubtematik = append(item.IndikatorSubtematik, newIndikator)
			}
		}
	}

	result := make([]domain.SasaranPemdaWithPokin, 0, len(subtematikMap))
	for _, item := range subtematikMap {
		if item.LevelPohon == 0 {
			result = append(result, domain.SasaranPemdaWithPokin{
				TematikId:   item.TematikId,
				NamaTematik: item.NamaTematik,
			})
		} else {
			result = append(result, *item)
		}
	}

	return result, nil

}

func (repository *SasaranPemdaRepositoryImpl) IsSubtemaIdExists(ctx context.Context, tx *sql.Tx, subtemaId int) bool {
	query := "SELECT COUNT(*) FROM tb_sasaran_pemda WHERE subtema_id = ?"
	var count int
	err := tx.QueryRowContext(ctx, query, subtemaId).Scan(&count)
	if err != nil {
		return false // Ubah return value jika error menjadi false
	}
	return count > 0
}
