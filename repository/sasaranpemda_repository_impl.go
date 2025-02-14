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
	// Insert sasaran pemda
	query := "INSERT INTO tb_sasaran_pemda(id, tujuan_pemda_id, subtema_id, sasaran_pemda, periode_id) VALUES (?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, query, sasaranPemda.Id, sasaranPemda.TujuanPemdaId, sasaranPemda.SubtemaId, sasaranPemda.SasaranPemda, sasaranPemda.PeriodeId)
	if err != nil {
		return sasaranPemda, err
	}

	// Insert indikator
	for _, indikator := range sasaranPemda.Indikator {
		scriptInsertIndikator := `
            INSERT INTO tb_indikator 
                (id, sasaran_pemda_id, indikator, rumus_perhitungan, sumber_data) 
            VALUES 
                (?, ?, ?, ?, ?)`

		_, err := tx.ExecContext(ctx, scriptInsertIndikator,
			indikator.Id,
			sasaranPemda.Id,
			indikator.Indikator,
			indikator.RumusPerhitungan,
			indikator.SumberData)
		if err != nil {
			return sasaranPemda, err
		}

		// Insert target untuk setiap indikator
		for _, target := range indikator.Target {
			// Skip jika target kosong
			if target.Target == "" && target.Satuan == "" {
				continue
			}

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
				return sasaranPemda, err
			}
		}
	}

	return sasaranPemda, nil
}

func (repository *SasaranPemdaRepositoryImpl) CreateIndikator(ctx context.Context, tx *sql.Tx, indikator domain.Indikator) (domain.Indikator, error) {
	query := "INSERT INTO tb_indikator(id, sasaran_pemda_id, indikator, rumus_perhitungan, sumber_data) VALUES (?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, query, indikator.Id, indikator.SasaranPemdaId, indikator.Indikator, indikator.RumusPerhitungan, indikator.SumberData)
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
	// Update sasaran pemda
	query := "UPDATE tb_sasaran_pemda SET sasaran_pemda = ?, tujuan_pemda_id = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, sasaranPemda.SasaranPemda, sasaranPemda.TujuanPemdaId, sasaranPemda.Id)
	if err != nil {
		return sasaranPemda, err
	}

	// Hapus semua indikator lama beserta targetnya
	scriptDeleteOldIndicators := "DELETE FROM tb_indikator WHERE sasaran_pemda_id = ?"
	_, err = tx.ExecContext(ctx, scriptDeleteOldIndicators, sasaranPemda.Id)
	if err != nil {
		return sasaranPemda, err
	}

	// Insert indikator baru
	for _, indikator := range sasaranPemda.Indikator {
		scriptInsertIndikator := `
            INSERT INTO tb_indikator 
                (id, sasaran_pemda_id, indikator, rumus_perhitungan, sumber_data) 
            VALUES 
                (?, ?, ?, ?, ?)`

		_, err := tx.ExecContext(ctx, scriptInsertIndikator,
			indikator.Id,
			sasaranPemda.Id,
			indikator.Indikator,
			indikator.RumusPerhitungan,
			indikator.SumberData)
		if err != nil {
			return sasaranPemda, err
		}

		// Hapus semua target lama untuk indikator ini
		scriptDeleteOldTargets := "DELETE FROM tb_target WHERE indikator_id = ?"
		_, err = tx.ExecContext(ctx, scriptDeleteOldTargets, indikator.Id)
		if err != nil {
			return sasaranPemda, err
		}

		// Insert target baru
		for _, target := range indikator.Target {
			// Skip jika target kosong
			if target.Target == "" && target.Satuan == "" {
				continue
			}

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
				return sasaranPemda, err
			}
		}
	}

	return sasaranPemda, nil
}

func (repository *SasaranPemdaRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, sasaranPemdaId int) error {
	queryDeleteTarget := `
	DELETE t FROM tb_target t
	INNER JOIN tb_indikator i ON t.indikator_id = i.id
	WHERE i.sasaran_pemda_id = ?`
	_, err := tx.ExecContext(ctx, queryDeleteTarget, sasaranPemdaId)
	if err != nil {
		return err
	}

	// 2. Hapus indikator yang terkait dengan tujuan pemda
	queryDeleteIndikator := "DELETE FROM tb_indikator WHERE sasaran_pemda_id = ?"
	_, err = tx.ExecContext(ctx, queryDeleteIndikator, sasaranPemdaId)
	if err != nil {
		return err
	}
	query := "DELETE FROM tb_sasaran_pemda WHERE id = ?"
	_, err = tx.ExecContext(ctx, query, sasaranPemdaId)
	if err != nil {
		return err
	}

	return nil
}

func (repository *SasaranPemdaRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, sasaranPemdaId int) (domain.SasaranPemda, error) {
	query := `
        SELECT DISTINCT
            tp.id,
            tp.tujuan_pemda_id,
            tp.subtema_id,
            tp.sasaran_pemda,
            tp.periode_id,
            COALESCE(p.tahun_awal, '') as tahun_awal,
            COALESCE(p.tahun_akhir, '') as tahun_akhir,
            pk.jenis_pohon,
            i.id as indikator_id,
            i.indikator as indikator_text,
            i.rumus_perhitungan,
            i.sumber_data,
            t.id as target_id,
            t.target as target_value,
            t.satuan as target_satuan,
            t.tahun as target_tahun
        FROM 
            tb_sasaran_pemda tp
            LEFT JOIN tb_periode p ON tp.periode_id = p.id
            LEFT JOIN tb_pohon_kinerja pk ON tp.subtema_id = pk.id
            LEFT JOIN tb_indikator i ON tp.id = i.sasaran_pemda_id
            LEFT JOIN tb_target t ON t.indikator_id = i.id
        WHERE tp.id = ?
        ORDER BY 
            tp.id, i.id, CAST(t.tahun AS SIGNED)`

	rows, err := tx.QueryContext(ctx, query, sasaranPemdaId)
	if err != nil {
		return domain.SasaranPemda{}, fmt.Errorf("error querying sasaran pemda: %v", err)
	}
	defer rows.Close()

	var result domain.SasaranPemda
	var firstRow = true
	indikatorMap := make(map[string]*domain.Indikator)

	for rows.Next() {
		var (
			id, subtemaId, periodeId, tujuanPemdaId             int
			sasaranPemdaText, tahunAwal, tahunAkhir, jenisPohon string // Tambahkan jenisPohon
			indikatorId, indikatorText                          sql.NullString
			rumusPerhitunganNull, sumberDataNull                sql.NullString
			targetId, targetValue, targetSatuan, targetTahun    sql.NullString
		)

		err := rows.Scan(
			&id,
			&tujuanPemdaId,
			&subtemaId,
			&sasaranPemdaText,
			&periodeId,
			&tahunAwal,
			&tahunAkhir,
			&jenisPohon,
			&indikatorId,
			&indikatorText,
			&rumusPerhitunganNull,
			&sumberDataNull,
			&targetId,
			&targetValue,
			&targetSatuan,
			&targetTahun,
		)
		if err != nil {
			return domain.SasaranPemda{}, fmt.Errorf("error scanning row: %v", err)
		}

		if firstRow {
			result = domain.SasaranPemda{
				Id:            id,
				TujuanPemdaId: tujuanPemdaId,
				SubtemaId:     subtemaId,
				SasaranPemda:  sasaranPemdaText,
				PeriodeId:     periodeId,
				JenisPohon:    jenisPohon,
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
				// Konversi NullString ke string biasa
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
					SasaranPemdaId:   id,
					Indikator:        indikatorText.String,
					RumusPerhitungan: sql.NullString{String: rumusPerhitungan, Valid: rumusPerhitungan != ""},
					SumberData:       sql.NullString{String: sumberData, Valid: sumberData != ""},
					Target:           []domain.Target{},
				}

				// Generate target untuk setiap tahun dalam periode
				if periodeId != 0 && tahunAwal != "" && tahunAkhir != "" {
					tahunAwalInt, _ := strconv.Atoi(tahunAwal)
					tahunAkhirInt, _ := strconv.Atoi(tahunAkhir)

					for tahun := tahunAwalInt; tahun <= tahunAkhirInt; tahun++ {
						tahunStr := strconv.Itoa(tahun)
						indikator.Target = append(indikator.Target, domain.Target{
							Id:          "-",
							IndikatorId: indikatorId.String,
							Target:      "",
							Satuan:      "",
							Tahun:       tahunStr,
						})
					}
				}

				result.Indikator = append(result.Indikator, indikator)
				indikatorMap[indikatorId.String] = &result.Indikator[len(result.Indikator)-1]
				currentIndikator = &result.Indikator[len(result.Indikator)-1]
			}

			// Update target yang ada dengan data sebenarnya
			if targetId.Valid && targetValue.Valid && targetTahun.Valid {
				tahunInt, _ := strconv.Atoi(targetTahun.String)
				tahunAwalInt, _ := strconv.Atoi(tahunAwal)
				if tahunInt >= tahunAwalInt {
					idx := tahunInt - tahunAwalInt
					if idx >= 0 && idx < len(currentIndikator.Target) {
						currentIndikator.Target[idx] = domain.Target{
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
	}

	if err = rows.Err(); err != nil {
		return domain.SasaranPemda{}, fmt.Errorf("error iterating rows: %v", err)
	}

	if result.Id == 0 {
		return domain.SasaranPemda{}, fmt.Errorf("sasaran pemda dengan id %d tidak ditemukan", sasaranPemdaId)
	}

	// Sort indikator berdasarkan ID
	sort.Slice(result.Indikator, func(i, j int) bool {
		return result.Indikator[i].Id < result.Indikator[j].Id
	})

	return result, nil
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
    WITH RECURSIVE tahun_periode AS (
        SELECT DISTINCT p.id as periode_id, p.tahun_awal, p.tahun_akhir
        FROM tb_periode p
        JOIN tb_tahun_periode tp ON p.id = tp.id_periode
        WHERE CAST(? AS SIGNED) BETWEEN CAST(p.tahun_awal AS SIGNED) AND CAST(p.tahun_akhir AS SIGNED)
    ),
    pohon_kinerja_aktif AS (
        SELECT DISTINCT pk.*
        FROM tb_pohon_kinerja pk
        JOIN tahun_periode tp 
        WHERE CAST(pk.tahun AS SIGNED) BETWEEN CAST(tp.tahun_awal AS SIGNED) AND CAST(tp.tahun_akhir AS SIGNED)
    ),
    pohon_hierarchy AS (
        -- Ambil semua level
        SELECT 
            pk.id,
            pk.nama_pohon,
            pk.parent,
            pk.level_pohon,
            pk.jenis_pohon,
            pk.keterangan,
            CAST(pk.id AS CHAR(50)) as path,
            pk.id as root_id,
            pk.nama_pohon as root_nama
        FROM pohon_kinerja_aktif pk
        WHERE pk.level_pohon = 0

        UNION ALL

        SELECT 
            c.id,
            c.nama_pohon,
            c.parent,
            c.level_pohon,
            c.jenis_pohon,
            c.keterangan,
            CONCAT(ph.path, ',', c.id),
            ph.root_id,
            ph.root_nama
        FROM pohon_kinerja_aktif c
        JOIN pohon_hierarchy ph ON c.parent = ph.id
        WHERE c.level_pohon BETWEEN 1 AND 3
    )
    SELECT DISTINCT
        pk.id as subtematik_id,
        pk.nama_pohon as nama_subtematik,
        pk.jenis_pohon,
        pk.level_pohon,
        pk.keterangan,
        pk.root_id as tematik_id,
        pk.root_nama as nama_tematik,
        COALESCE(sp.id, 0) as id_sasaran_pemda,
        COALESCE(sp.sasaran_pemda, '') as sasaran_pemda,
        sp.periode_id,
        tp.tahun_awal,
        tp.tahun_akhir,
        i.id as indikator_id,
        i.indikator,
        t.id as target_id,
        t.target,
        t.satuan,
        t.tahun
    FROM pohon_hierarchy pk
    LEFT JOIN tb_sasaran_pemda sp ON pk.id = sp.subtema_id
    CROSS JOIN tahun_periode tp
    LEFT JOIN tb_indikator i ON sp.id = i.sasaran_pemda_id
    LEFT JOIN tb_target t ON i.id = t.indikator_id
    WHERE pk.level_pohon BETWEEN 0 AND 3
        AND (sp.periode_id IS NULL OR sp.periode_id = tp.periode_id)
    ORDER BY 
        pk.root_nama,
        pk.level_pohon,
        pk.id`

	rows, err := tx.QueryContext(ctx, query, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pokinMap := make(map[int]*domain.SasaranPemdaWithPokin)
	indikatorMap := make(map[int]map[string]*domain.Indikator)

	for rows.Next() {
		var (
			subtematikId                           int
			namaSubtematik, jenisPohon, keterangan string
			levelPohon                             int
			tematikId                              int
			namaTematik                            string
			idSasaranPemda                         int
			sasaranPemda                           string
			periodeId                              sql.NullInt64
			tahunAwal, tahunAkhir                  string
			indikatorId, indikator                 sql.NullString
			targetId, target, satuan, targetTahun  sql.NullString
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
			&periodeId,
			&tahunAwal,
			&tahunAkhir,
			&indikatorId,
			&indikator,
			&targetId,
			&target,
			&satuan,
			&targetTahun,
		)
		if err != nil {
			return nil, err
		}

		if _, exists := pokinMap[subtematikId]; !exists {
			pokinMap[subtematikId] = &domain.SasaranPemdaWithPokin{
				SubtematikId:   subtematikId,
				NamaSubtematik: namaSubtematik,
				JenisPohon:     jenisPohon,
				LevelPohon:     levelPohon,
				Keterangan:     keterangan,
				TematikId:      tematikId,
				NamaTematik:    namaTematik,
				IdsasaranPemda: idSasaranPemda,
				SasaranPemda:   sasaranPemda,
			}

			if idSasaranPemda > 0 {
				indikatorMap[idSasaranPemda] = make(map[string]*domain.Indikator)
			}
		}

		// Proses indikator dan target
		if idSasaranPemda > 0 && indikatorId.Valid {
			if _, exists := indikatorMap[idSasaranPemda][indikatorId.String]; !exists {
				indikatorMap[idSasaranPemda][indikatorId.String] = &domain.Indikator{
					Id:        indikatorId.String,
					Indikator: indikator.String,
					Target:    []domain.Target{},
				}
			}

			currentIndikator := indikatorMap[idSasaranPemda][indikatorId.String]

			// Buat target untuk semua tahun dalam periode
			tahunAwalInt, _ := strconv.Atoi(tahunAwal)
			tahunAkhirInt, _ := strconv.Atoi(tahunAkhir)

			// Hanya inisialisasi target jika belum ada
			if len(currentIndikator.Target) == 0 {
				for tahun := tahunAwalInt; tahun <= tahunAkhirInt; tahun++ {
					tahunStr := strconv.Itoa(tahun)
					currentIndikator.Target = append(currentIndikator.Target, domain.Target{
						Id:          "-",
						IndikatorId: indikatorId.String,
						Target:      "",
						Satuan:      "",
						Tahun:       tahunStr,
					})
				}
			}

			// Update target yang memiliki data
			if targetId.Valid && targetTahun.Valid {
				tahunInt, _ := strconv.Atoi(targetTahun.String)
				if tahunInt >= tahunAwalInt && tahunInt <= tahunAkhirInt {
					idx := tahunInt - tahunAwalInt
					if idx >= 0 && idx < len(currentIndikator.Target) {
						currentIndikator.Target[idx] = domain.Target{
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
	}

	var result []domain.SasaranPemdaWithPokin
	for _, pokin := range pokinMap {
		if pokin.IdsasaranPemda > 0 {
			var indikators []domain.Indikator
			for _, indikator := range indikatorMap[pokin.IdsasaranPemda] {
				indikators = append(indikators, *indikator)
			}
			// Sort indikator berdasarkan ID
			sort.Slice(indikators, func(i, j int) bool {
				return indikators[i].Id < indikators[j].Id
			})
			pokin.IndikatorSubtematik = indikators
		}
		result = append(result, *pokin)
	}

	// Sort hasil akhir
	sort.Slice(result, func(i, j int) bool {
		if result[i].TematikId != result[j].TematikId {
			return result[i].TematikId < result[j].TematikId
		}
		if result[i].LevelPohon != result[j].LevelPohon {
			return result[i].LevelPohon < result[j].LevelPohon
		}
		return result[i].SubtematikId < result[j].SubtematikId
	})

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
