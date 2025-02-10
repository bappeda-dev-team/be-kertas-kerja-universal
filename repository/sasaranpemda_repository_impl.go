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

type SasaranPemdaRepositoryImpl struct {
}

func NewSasaranPemdaRepositoryImpl() *SasaranPemdaRepositoryImpl {
	return &SasaranPemdaRepositoryImpl{}
}

func (repository *SasaranPemdaRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, sasaranPemda domain.SasaranPemda) (domain.SasaranPemda, error) {
	query := "INSERT INTO tb_sasaran_pemda(id,subtema_id, sasaran_pemda) VALUES (?, ?, ?)"
	_, err := tx.ExecContext(ctx, query, sasaranPemda.Id, sasaranPemda.SubtemaId, sasaranPemda.SasaranPemda)
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
	query := "UPDATE tb_sasaran_pemda SET sasaran_pemda = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, sasaranPemda.SasaranPemda, sasaranPemda.Id)
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
        WITH periode_aktif AS (
            SELECT id, tahun_awal, tahun_akhir
            FROM tb_periode
            WHERE id = (
                SELECT periode_id 
                FROM tb_sasaran_pemda 
                WHERE id = ?
            )
        )
        SELECT DISTINCT
            sp.id,
            sp.subtema_id,
            pk.nama_pohon as nama_subtema,
            sp.sasaran_pemda,
            sp.periode_id,
            p.tahun_awal,
            p.tahun_akhir,
            i.id as indikator_id,
            i.indikator as indikator_text,
            t.id as target_id,
            t.target as target_value,
            t.satuan as target_satuan,
            t.tahun as target_tahun
        FROM 
            tb_sasaran_pemda sp
            INNER JOIN tb_pohon_kinerja pk ON sp.subtema_id = pk.id
            LEFT JOIN tb_periode p ON sp.periode_id = p.id
            LEFT JOIN tb_indikator i ON sp.id = i.sasaran_pemda_id
            LEFT JOIN tb_target t ON i.id = t.indikator_id
        WHERE sp.id = ?
        ORDER BY 
            i.id, CAST(t.tahun AS SIGNED)`

	rows, err := tx.QueryContext(ctx, query, sasaranPemdaId, sasaranPemdaId)
	if err != nil {
		return domain.SasaranPemda{}, err
	}
	defer rows.Close()

	var sasaranPemda domain.SasaranPemda
	var firstRow = true
	indikatorMap := make(map[string]*domain.Indikator)
	targetMap := make(map[string]map[string]domain.Target) // map[indikatorId]map[tahun]Target

	for rows.Next() {
		var (
			id, subtemaId, periodeId                         int
			namaSubtema, sasaranPemdaText                    string
			tahunAwal, tahunAkhir                            string
			indikatorId, indikatorText                       sql.NullString
			targetId, targetValue, targetSatuan, targetTahun sql.NullString
		)

		err := rows.Scan(
			&id,
			&subtemaId,
			&namaSubtema,
			&sasaranPemdaText,

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
			return domain.SasaranPemda{}, err
		}

		if firstRow {
			sasaranPemda = domain.SasaranPemda{
				Id:           id,
				SubtemaId:    subtemaId,
				NamaSubtema:  namaSubtema,
				SasaranPemda: sasaranPemdaText,
				PeriodeId:    periodeId,
				Periode: domain.Periode{
					Id:         periodeId,
					TahunAwal:  tahunAwal,
					TahunAkhir: tahunAkhir,
				},
				Indikator: []domain.Indikator{},
			}
			firstRow = false
		}

		// Proses indikator jika ada
		if indikatorId.Valid && indikatorText.Valid {
			if _, exists := indikatorMap[indikatorId.String]; !exists {
				indikatorMap[indikatorId.String] = &domain.Indikator{
					Id:             indikatorId.String,
					SasaranPemdaId: id,
					Indikator:      indikatorText.String,
					Target:         []domain.Target{},
				}
				targetMap[indikatorId.String] = make(map[string]domain.Target)
			}

			// Proses target jika ada
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

	// Proses final untuk menambahkan target sesuai periode
	if sasaranPemda.Periode.TahunAwal != "" && sasaranPemda.Periode.TahunAkhir != "" {
		tahunAwal, _ := strconv.Atoi(sasaranPemda.Periode.TahunAwal)
		tahunAkhir, _ := strconv.Atoi(sasaranPemda.Periode.TahunAkhir)

		for indikatorId, indikator := range indikatorMap {
			var targets []domain.Target
			for tahun := tahunAwal; tahun <= tahunAkhir; tahun++ {
				tahunStr := strconv.Itoa(tahun)
				if target, exists := targetMap[indikatorId][tahunStr]; exists {
					targets = append(targets, target)
				} else {
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
			sasaranPemda.Indikator = append(sasaranPemda.Indikator, *indikator)
		}
	}

	// Sort indikator berdasarkan ID
	sort.Slice(sasaranPemda.Indikator, func(i, j int) bool {
		return sasaranPemda.Indikator[i].Id < sasaranPemda.Indikator[j].Id
	})

	if sasaranPemda.Id == 0 {
		return sasaranPemda, errors.New("sasaran pemda not found")
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
	query := `
 WITH RECURSIVE parent_hierarchy AS (
        -- Level 1-3 nodes
        SELECT 
            id,
            nama_pohon,
            jenis_pohon,
            level_pohon,
            keterangan,
            parent,
            id as original_id,
            CASE 
                WHEN level_pohon = 1 THEN parent
                ELSE NULL 
            END as tematik_id
        FROM 
            tb_pohon_kinerja
        WHERE 
            level_pohon BETWEEN 1 AND 3
            AND tahun = ?

        UNION ALL

        -- Recursive untuk mencari parent sampai level 0
        SELECT 
            pk.id,
            pk.nama_pohon,
            pk.jenis_pohon,
            pk.level_pohon,
            pk.keterangan,
            pk.parent,
            ph.original_id,
            CASE 
                WHEN pk.level_pohon = 0 THEN pk.id
                WHEN pk.level_pohon = 1 THEN pk.parent
                ELSE NULL
            END as tematik_id
        FROM 
            tb_pohon_kinerja pk
        INNER JOIN 
            parent_hierarchy ph ON pk.id = ph.parent
        WHERE 
            pk.tahun = ?
    )
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
        COALESCE(sp.sasaran_pemda, '') as sasaran_pemda,
        COALESCE(i_sp.indikator, i_pk.indikator) as indikator,
        COALESCE(t_sp.target, t_pk.target) as target,
        COALESCE(t_sp.satuan, t_pk.satuan) as satuan
    FROM 
        tb_pohon_kinerja pk
    -- Parent untuk level 1
    LEFT JOIN 
        tb_pohon_kinerja parent ON pk.parent = parent.id AND parent.tahun = ?
    -- Parent untuk level 2
    LEFT JOIN 
        tb_pohon_kinerja grandparent ON parent.parent = grandparent.id AND grandparent.tahun = ?
    -- Parent untuk level 3
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
    ORDER BY 
        pk.nama_pohon ASC,
        pk.id, 
        COALESCE(i_sp.id, i_pk.id)`

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
			sasaranPemda   string
			indikator      sql.NullString
			target         sql.NullString
			satuan         sql.NullString
		)

		err := rows.Scan(
			&subtematikId,
			&namaSubtematik,
			&jenisPohon,
			&levelPohon,
			&keterangan,
			&tematikId,
			&namaTematik,
			&sasaranPemda,
			&indikator,
			&target,
			&satuan,
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
		result = append(result, *item)
	}

	return result, nil
}

func (repository *SasaranPemdaRepositoryImpl) IsSubtemaIdExists(ctx context.Context, tx *sql.Tx, subtemaId int) bool {
	query := "SELECT COUNT(*) FROM tb_sasaran_pemda WHERE subtema_id = ?"
	var count int
	err := tx.QueryRowContext(ctx, query, subtemaId).Scan(&count)
	if err != nil {
		return true // Mengembalikan true untuk berjaga-jaga jika terjadi error
	}
	return count > 0
}
