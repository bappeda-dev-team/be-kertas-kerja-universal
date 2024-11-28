package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
	"sort"
	"strings"
)

type PohonKinerjaRepositoryImpl struct {
}

func NewPohonKinerjaRepositoryImpl() *PohonKinerjaRepositoryImpl {
	return &PohonKinerjaRepositoryImpl{}
}

func (repository *PohonKinerjaRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, pohonKinerja domain.PohonKinerja) (domain.PohonKinerja, error) {
	scriptPokin := "INSERT INTO tb_pohon_kinerja (nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, scriptPokin,
		pohonKinerja.NamaPohon,
		pohonKinerja.Parent,
		pohonKinerja.JenisPohon,
		pohonKinerja.LevelPohon,
		pohonKinerja.KodeOpd,
		pohonKinerja.Keterangan,
		pohonKinerja.Tahun,
		pohonKinerja.Status)
	if err != nil {
		return pohonKinerja, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return pohonKinerja, err
	}
	pohonKinerja.Id = int(lastInsertId)

	scriptPelaksana := "INSERT INTO tb_pelaksana_pokin (id, pohon_kinerja_id, pegawai_id) VALUES (?, ?, ?)"
	for _, pelaksana := range pohonKinerja.Pelaksana {
		_, err = tx.ExecContext(ctx, scriptPelaksana,
			pelaksana.Id,                // id pelaksana_pokin
			fmt.Sprint(pohonKinerja.Id), // pohon_kinerja_id dalam string
			pelaksana.PegawaiId)         // pegawai_id
		if err != nil {
			return pohonKinerja, err
		}
	}

	// Insert indikator
	for _, indikator := range pohonKinerja.Indikator {
		scriptIndikator := "INSERT INTO tb_indikator (id, pokin_id, indikator, tahun) VALUES (?, ?, ?, ?)"
		_, err := tx.ExecContext(ctx, scriptIndikator,
			indikator.Id,
			pohonKinerja.Id,
			indikator.Indikator,
			indikator.Tahun)
		if err != nil {
			return pohonKinerja, err
		}

		// Insert target untuk setiap indikator
		for _, target := range indikator.Target {
			scriptTarget := "INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)"
			_, err := tx.ExecContext(ctx, scriptTarget,
				target.Id,
				indikator.Id,
				target.Target,
				target.Satuan,
				target.Tahun)
			if err != nil {
				return pohonKinerja, err
			}
		}
	}

	return pohonKinerja, nil
}

func (repository *PohonKinerjaRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, pohonKinerja domain.PohonKinerja) (domain.PohonKinerja, error) {
	// Update tb_pohon_kinerja
	scriptPokin := "UPDATE tb_pohon_kinerja SET nama_pohon = ?, parent = ?, jenis_pohon = ?, level_pohon = ?, kode_opd = ?, keterangan = ?, tahun = ?, status = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, scriptPokin,
		pohonKinerja.NamaPohon,
		pohonKinerja.Parent,
		pohonKinerja.JenisPohon,
		pohonKinerja.LevelPohon,
		pohonKinerja.KodeOpd,
		pohonKinerja.Keterangan,
		pohonKinerja.Tahun,
		pohonKinerja.Status,
		pohonKinerja.Id)
	if err != nil {
		return pohonKinerja, err
	}

	// Hapus data pelaksana yang lama
	scriptDeletePelaksana := "DELETE FROM tb_pelaksana_pokin WHERE pohon_kinerja_id = ?"
	_, err = tx.ExecContext(ctx, scriptDeletePelaksana, fmt.Sprint(pohonKinerja.Id))
	if err != nil {
		return pohonKinerja, err
	}

	// Insert data pelaksana baru
	for _, pelaksana := range pohonKinerja.Pelaksana {
		scriptPelaksana := "INSERT INTO tb_pelaksana_pokin (id, pohon_kinerja_id, pegawai_id) VALUES (?, ?, ?)"
		_, err := tx.ExecContext(ctx, scriptPelaksana,
			pelaksana.Id,
			fmt.Sprint(pohonKinerja.Id),
			pelaksana.PegawaiId)
		if err != nil {
			return pohonKinerja, err
		}
	}

	// Hapus indikator dan target yang lama
	scriptDeleteTarget := "DELETE FROM tb_target WHERE indikator_id IN (SELECT id FROM tb_indikator WHERE pokin_id = ?)"
	_, err = tx.ExecContext(ctx, scriptDeleteTarget, pohonKinerja.Id)
	if err != nil {
		return pohonKinerja, err
	}

	scriptDeleteIndikator := "DELETE FROM tb_indikator WHERE pokin_id = ?"
	_, err = tx.ExecContext(ctx, scriptDeleteIndikator, pohonKinerja.Id)
	if err != nil {
		return pohonKinerja, err
	}

	// Insert indikator baru
	for _, indikator := range pohonKinerja.Indikator {
		scriptIndikator := "INSERT INTO tb_indikator (id, pokin_id, indikator, tahun) VALUES (?, ?, ?, ?)"
		_, err := tx.ExecContext(ctx, scriptIndikator,
			indikator.Id,
			pohonKinerja.Id,
			indikator.Indikator,
			indikator.Tahun)
		if err != nil {
			return pohonKinerja, err
		}

		// Insert target baru untuk setiap indikator
		for _, target := range indikator.Target {
			scriptTarget := "INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)"
			_, err := tx.ExecContext(ctx, scriptTarget,
				target.Id,
				indikator.Id,
				target.Target,
				target.Satuan,
				target.Tahun)
			if err != nil {
				return pohonKinerja, err
			}
		}
	}

	return pohonKinerja, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.PohonKinerja, error) {
	scriptPokin := `
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
        FROM 
            tb_pohon_kinerja pk 
        WHERE 
            pk.id = ?`

	rows, err := tx.QueryContext(ctx, scriptPokin, id)
	if err != nil {
		return domain.PohonKinerja{}, err
	}
	defer rows.Close()

	pohonKinerja := domain.PohonKinerja{}
	if rows.Next() {
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
		)
		if err != nil {
			return domain.PohonKinerja{}, err
		}
	}

	return pohonKinerja, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string) ([]domain.PohonKinerja, error) {
	script := `
        SELECT 
            pk.id,
            pk.nama_pohon,
            pk.parent,
            pk.jenis_pohon,
            pk.level_pohon,
            pk.kode_opd,
            pk.keterangan,
            pk.tahun,
            pk.created_at,
            pk.status
        FROM 
            tb_pohon_kinerja pk
        WHERE 
            pk.kode_opd = ? 
            AND pk.tahun = ?
            AND pk.status = 'disetujui'
        ORDER BY 
            pk.level_pohon, pk.id, pk.created_at asc
    `

	rows, err := tx.QueryContext(ctx, script, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.PohonKinerja
	for rows.Next() {
		var pokin domain.PohonKinerja
		err := rows.Scan(
			&pokin.Id,
			&pokin.NamaPohon,
			&pokin.Parent,
			&pokin.JenisPohon,
			&pokin.LevelPohon,
			&pokin.KodeOpd,
			&pokin.Keterangan,
			&pokin.Tahun,
			&pokin.CreatedAt,
			&pokin.Status,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, pokin)
	}

	// Query untuk mendapatkan data pelaksana
	for i, pokin := range result {
		// Query pelaksana
		scriptPelaksana := `
            SELECT 
                id, pegawai_id
            FROM 
                tb_pelaksana_pokin
            WHERE 
                pohon_kinerja_id = ?
        `

		pelaksanaRows, err := tx.QueryContext(ctx, scriptPelaksana, fmt.Sprint(pokin.Id))
		if err != nil {
			return nil, err
		}
		defer pelaksanaRows.Close()

		var pelaksanaList []domain.PelaksanaPokin
		for pelaksanaRows.Next() {
			var pelaksana domain.PelaksanaPokin
			err := pelaksanaRows.Scan(
				&pelaksana.Id,
				&pelaksana.PegawaiId,
			)
			if err != nil {
				return nil, err
			}
			pelaksanaList = append(pelaksanaList, pelaksana)
		}
		result[i].Pelaksana = pelaksanaList

		// Query indikator dan target
		scriptIndikator := `
            SELECT 
                i.id,
                i.indikator,
                i.tahun,
                t.id as target_id,
                t.target,
                t.satuan,
                t.tahun as target_tahun
            FROM 
                tb_indikator i
            LEFT JOIN 
                tb_target t ON i.id = t.indikator_id
            WHERE 
                i.pokin_id = ?
        `

		indikatorRows, err := tx.QueryContext(ctx, scriptIndikator, pokin.Id)
		if err != nil {
			return nil, err
		}
		defer indikatorRows.Close()

		indikatorMap := make(map[string]domain.Indikator)
		for indikatorRows.Next() {
			var (
				indikatorId, indikatorNama, indikatorTahun string
				targetId, target, satuan, targetTahun      sql.NullString
			)

			err := indikatorRows.Scan(
				&indikatorId,
				&indikatorNama,
				&indikatorTahun,
				&targetId,
				&target,
				&satuan,
				&targetTahun,
			)
			if err != nil {
				return nil, err
			}

			indikator, exists := indikatorMap[indikatorId]
			if !exists {
				indikator = domain.Indikator{
					Id:        indikatorId,
					PokinId:   fmt.Sprint(pokin.Id),
					Indikator: indikatorNama,
					Tahun:     indikatorTahun,
				}
			}

			if targetId.Valid && target.Valid && satuan.Valid && targetTahun.Valid {
				targetObj := domain.Target{
					Id:          targetId.String,
					IndikatorId: indikatorId,
					Target:      target.String,
					Satuan:      satuan.String,
					Tahun:       targetTahun.String,
				}
				indikator.Target = append(indikator.Target, targetObj)
			}

			indikatorMap[indikatorId] = indikator
		}

		// Convert indikator map to slice
		var indikatorList []domain.Indikator
		for _, indikator := range indikatorMap {
			indikatorList = append(indikatorList, indikator)
		}
		result[i].Indikator = indikatorList
	}

	// Modifikasi bagian pemrosesan data
	pohonMap := make(map[int]map[int][]domain.PohonKinerja)
	maxLevel := 0

	for _, pokin := range result {
		if pokin.LevelPohon > maxLevel {
			maxLevel = pokin.LevelPohon
		}
	}

	for i := 4; i <= maxLevel; i++ {
		pohonMap[i] = make(map[int][]domain.PohonKinerja)
	}

	for _, p := range result {
		if p.LevelPohon >= 4 {
			pohonMap[p.LevelPohon][p.Parent] = append(
				pohonMap[p.LevelPohon][p.Parent],
				p,
			)
		}
	}

	return result, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindStrategicNoParent(ctx context.Context, tx *sql.Tx, levelPohon, parent int, kodeOpd, tahun string) ([]domain.PohonKinerja, error) {
	script := "SELECT id, nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun FROM tb_pohon_kinerja WHERE level_pohon = ? AND parent = ? AND kode_opd = ? AND tahun = ?"
	rows, err := tx.QueryContext(ctx, script, levelPohon, parent, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.PohonKinerja
	for rows.Next() {
		var pokin domain.PohonKinerja
		err := rows.Scan(&pokin.Id, &pokin.NamaPohon, &pokin.Parent, &pokin.JenisPohon, &pokin.LevelPohon, &pokin.KodeOpd, &pokin.Keterangan, &pokin.Tahun)
		if err != nil {
			return nil, err
		}
		result = append(result, pokin)
	}
	return result, nil
}

func (repository *PohonKinerjaRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	scriptDeleteTarget := `
        DELETE FROM tb_target 
        WHERE indikator_id IN (
            SELECT id FROM tb_indikator WHERE pokin_id = ?
        )`
	_, err := tx.ExecContext(ctx, scriptDeleteTarget, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus data target: %v", err)
	}

	scriptDeleteIndikator := "DELETE FROM tb_indikator WHERE pokin_id = ?"
	_, err = tx.ExecContext(ctx, scriptDeleteIndikator, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus data indikator: %v", err)
	}

	scriptDeletePelaksana := "DELETE FROM tb_pelaksana_pokin WHERE pohon_kinerja_id = ?"
	_, err = tx.ExecContext(ctx, scriptDeletePelaksana, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus data pelaksana: %v", err)
	}

	scriptDeletePokin := "DELETE FROM tb_pohon_kinerja WHERE id = ?"
	_, err = tx.ExecContext(ctx, scriptDeletePokin, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus data pohon kinerja: %v", err)
	}

	return nil
}

func (repository *PohonKinerjaRepositoryImpl) FindPelaksanaPokin(ctx context.Context, tx *sql.Tx, pohonKinerjaId string) ([]domain.PelaksanaPokin, error) {
	script := "SELECT id, pohon_kinerja_id, pegawai_id FROM tb_pelaksana_pokin WHERE pohon_kinerja_id = ?"
	rows, err := tx.QueryContext(ctx, script, pohonKinerjaId)
	helper.PanicIfError(err)
	defer rows.Close()

	var result []domain.PelaksanaPokin
	for rows.Next() {
		var pelaksana domain.PelaksanaPokin
		err := rows.Scan(&pelaksana.Id, &pelaksana.PohonKinerjaId, &pelaksana.PegawaiId)
		helper.PanicIfError(err)
		result = append(result, pelaksana)
	}
	return result, nil
}

func (repository *PohonKinerjaRepositoryImpl) DeletePelaksanaPokin(ctx context.Context, tx *sql.Tx, pelaksanaId string) error {
	script := "DELETE FROM tb_pelaksana_pokin WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, pelaksanaId)
	return err
}

// admin pokin
func (repository *PohonKinerjaRepositoryImpl) CreatePokinAdmin(ctx context.Context, tx *sql.Tx, pokinAdmin domain.PohonKinerja) (domain.PohonKinerja, error) {
	// Insert pohon kinerja tanpa ID
	scriptPokin := "INSERT INTO tb_pohon_kinerja (nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, scriptPokin,
		pokinAdmin.NamaPohon, pokinAdmin.Parent, pokinAdmin.JenisPohon, pokinAdmin.LevelPohon, pokinAdmin.KodeOpd, pokinAdmin.Keterangan, pokinAdmin.Tahun, pokinAdmin.Status)
	if err != nil {
		return pokinAdmin, err
	}

	// Dapatkan ID yang baru dibuat
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return pokinAdmin, err
	}
	pokinAdmin.Id = int(lastInsertId)

	// Tambahkan insert pelaksana
	scriptPelaksana := "INSERT INTO tb_pelaksana_pokin (id, pohon_kinerja_id, pegawai_id) VALUES (?, ?, ?)"
	for _, pelaksana := range pokinAdmin.Pelaksana {
		_, err = tx.ExecContext(ctx, scriptPelaksana,
			pelaksana.Id,
			fmt.Sprint(pokinAdmin.Id),
			pelaksana.PegawaiId)
		if err != nil {
			return pokinAdmin, err
		}
	}

	// Insert indikator
	for _, indikator := range pokinAdmin.Indikator {
		scriptIndikator := "INSERT INTO tb_indikator (id, pokin_id, indikator, tahun) VALUES (?, ?, ?, ?)"
		_, err := tx.ExecContext(ctx, scriptIndikator,
			indikator.Id, pokinAdmin.Id, indikator.Indikator,
			indikator.Tahun)
		if err != nil {
			return pokinAdmin, err
		}

		// Insert target untuk setiap indikator
		for _, target := range indikator.Target {
			scriptTarget := "INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)"
			_, err := tx.ExecContext(ctx, scriptTarget, target.Id, indikator.Id, target.Target, target.Satuan, target.Tahun)
			if err != nil {
				return pokinAdmin, err
			}
		}
	}

	return pokinAdmin, nil
}

func (repository *PohonKinerjaRepositoryImpl) UpdatePokinAdmin(ctx context.Context, tx *sql.Tx, pokinAdmin domain.PohonKinerja) (domain.PohonKinerja, error) {
	// Update data pohon kinerja
	scriptPokin := "UPDATE tb_pohon_kinerja SET nama_pohon = ?, parent = ?, jenis_pohon = ?, level_pohon = ?, kode_opd = ?, keterangan = ?, tahun = ?, status = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, scriptPokin,
		pokinAdmin.NamaPohon,
		pokinAdmin.Parent,
		pokinAdmin.JenisPohon,
		pokinAdmin.LevelPohon,
		pokinAdmin.KodeOpd,
		pokinAdmin.Keterangan,
		pokinAdmin.Tahun,
		pokinAdmin.Status,
		pokinAdmin.Id)
	if err != nil {
		return pokinAdmin, err
	}

	// Hapus data pelaksana yang lama
	scriptDeletePelaksana := "DELETE FROM tb_pelaksana_pokin WHERE pohon_kinerja_id = ?"
	_, err = tx.ExecContext(ctx, scriptDeletePelaksana, fmt.Sprint(pokinAdmin.Id))
	if err != nil {
		return pokinAdmin, err
	}

	// Insert data pelaksana baru
	for _, pelaksana := range pokinAdmin.Pelaksana {
		scriptPelaksana := "INSERT INTO tb_pelaksana_pokin (id, pohon_kinerja_id, pegawai_id) VALUES (?, ?, ?)"
		_, err := tx.ExecContext(ctx, scriptPelaksana,
			pelaksana.Id,
			fmt.Sprint(pokinAdmin.Id),
			pelaksana.PegawaiId)
		if err != nil {
			return pokinAdmin, err
		}
	}

	// Hapus indikator dan target yang lama
	scriptDeleteTarget := "DELETE FROM tb_target WHERE indikator_id IN (SELECT id FROM tb_indikator WHERE pokin_id = ?)"
	_, err = tx.ExecContext(ctx, scriptDeleteTarget, pokinAdmin.Id)
	if err != nil {
		return pokinAdmin, err
	}

	scriptDeleteIndikator := "DELETE FROM tb_indikator WHERE pokin_id = ?"
	_, err = tx.ExecContext(ctx, scriptDeleteIndikator, pokinAdmin.Id)
	if err != nil {
		return pokinAdmin, err
	}

	// Insert indikator baru
	for _, indikator := range pokinAdmin.Indikator {
		scriptIndikator := "INSERT INTO tb_indikator (id, pokin_id, indikator, tahun) VALUES (?, ?, ?, ?)"
		_, err := tx.ExecContext(ctx, scriptIndikator,
			indikator.Id,
			pokinAdmin.Id,
			indikator.Indikator,
			indikator.Tahun)
		if err != nil {
			return pokinAdmin, err
		}

		// Insert target baru untuk setiap indikator
		for _, target := range indikator.Target {
			scriptTarget := "INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)"
			_, err := tx.ExecContext(ctx, scriptTarget,
				target.Id,
				indikator.Id,
				target.Target,
				target.Satuan,
				target.Tahun)
			if err != nil {
				return pokinAdmin, err
			}
		}
	}

	return pokinAdmin, nil
}

func (repository *PohonKinerjaRepositoryImpl) DeletePokinAdmin(ctx context.Context, tx *sql.Tx, id int) error {
	findIdsScript := `
        WITH RECURSIVE pohon_hierarki AS (
            -- Base case: node yang akan dihapus
            SELECT id 
            FROM tb_pohon_kinerja 
            WHERE id = ?
            
            UNION ALL
            
            -- Recursive case: semua child nodes
            SELECT pk.id
            FROM tb_pohon_kinerja pk
            INNER JOIN pohon_hierarki ph ON pk.parent = ph.id
        )
        SELECT id FROM pohon_hierarki;
    `

	rows, err := tx.QueryContext(ctx, findIdsScript, id)
	if err != nil {
		return fmt.Errorf("gagal mengambil hierarki ID: %v", err)
	}
	defer rows.Close()

	// Kumpulkan semua ID yang akan dihapus
	var idsToDelete []interface{}
	for rows.Next() {
		var idToDelete int
		if err := rows.Scan(&idToDelete); err != nil {
			return fmt.Errorf("gagal scan ID: %v", err)
		}
		idsToDelete = append(idsToDelete, idToDelete)
	}

	if len(idsToDelete) == 0 {
		return fmt.Errorf("tidak ada data yang akan dihapus")
	}

	// Buat placeholder untuk query IN clause
	placeholders := make([]string, len(idsToDelete))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	inClause := strings.Join(placeholders, ",")

	// Hapus target terlebih dahulu
	scriptDeleteTarget := fmt.Sprintf("DELETE FROM tb_target WHERE indikator_id IN (SELECT id FROM tb_indikator WHERE pokin_id IN (%s))", inClause)
	_, err = tx.ExecContext(ctx, scriptDeleteTarget, idsToDelete...)
	if err != nil {
		return fmt.Errorf("gagal menghapus target: %v", err)
	}

	// Hapus indikator
	scriptDeleteIndikator := fmt.Sprintf("DELETE FROM tb_indikator WHERE pokin_id IN (%s)", inClause)
	_, err = tx.ExecContext(ctx, scriptDeleteIndikator, idsToDelete...)
	if err != nil {
		return fmt.Errorf("gagal menghapus indikator: %v", err)
	}

	// Hapus pelaksana
	scriptDeletePelaksana := fmt.Sprintf("DELETE FROM tb_pelaksana_pokin WHERE pohon_kinerja_id IN (%s)", inClause)
	_, err = tx.ExecContext(ctx, scriptDeletePelaksana, idsToDelete...)
	if err != nil {
		return fmt.Errorf("gagal menghapus pelaksana: %v", err)
	}

	// Hapus pohon kinerja
	scriptDeletePokin := fmt.Sprintf("DELETE FROM tb_pohon_kinerja WHERE id IN (%s)", inClause)
	_, err = tx.ExecContext(ctx, scriptDeletePokin, idsToDelete...)
	if err != nil {
		return fmt.Errorf("gagal menghapus pohon kinerja: %v", err)
	}

	return nil
}

func (repository *PohonKinerjaRepositoryImpl) FindPokinAdminById(ctx context.Context, tx *sql.Tx, id int) (domain.PohonKinerja, error) {
	script := "SELECT id, nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun FROM tb_pohon_kinerja WHERE id = ?"
	rows, err := tx.QueryContext(ctx, script, id)
	if err != nil {
		return domain.PohonKinerja{}, err
	}
	defer rows.Close()

	pokinAdmin := domain.PohonKinerja{}
	if rows.Next() {
		rows.Scan(&pokinAdmin.Id, &pokinAdmin.NamaPohon, &pokinAdmin.Parent, &pokinAdmin.JenisPohon, &pokinAdmin.LevelPohon, &pokinAdmin.KodeOpd, &pokinAdmin.Keterangan, &pokinAdmin.Tahun)
	}
	return pokinAdmin, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindPokinAdminAll(ctx context.Context, tx *sql.Tx, tahun string) ([]domain.PohonKinerja, error) {
	script := `
        SELECT 
            pk.id,
            pk.nama_pohon,
            pk.parent,
            pk.jenis_pohon,
            pk.level_pohon,
            pk.kode_opd,
            pk.keterangan,
            pk.tahun,
            i.id as indikator_id,
            i.indikator as nama_indikator,
            t.id as target_id,
            t.target as target_value,
            t.satuan as target_satuan
        FROM 
            tb_pohon_kinerja pk
        LEFT JOIN 
            tb_indikator i ON pk.id = i.pokin_id
        LEFT JOIN 
            tb_target t ON i.id = t.indikator_id
        WHERE 
            pk.tahun = ?
        ORDER BY 
            pk.level_pohon, pk.id, i.id, t.id
    `

	rows, err := tx.QueryContext(ctx, script, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map untuk menyimpan pohon kinerja yang sudah diproses
	pokinMap := make(map[int]domain.PohonKinerja)
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
			pokin = domain.PohonKinerja{
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
			pokin.Indikator = append(pokin.Indikator, indikator)
			pokinMap[pokinId] = pokin
		}
	}

	// Konversi map ke slice
	var result []domain.PohonKinerja
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

func (repository *PohonKinerjaRepositoryImpl) FindPokinAdminByIdHierarki(ctx context.Context, tx *sql.Tx, idPokin int) ([]domain.PohonKinerja, error) {
	script := `
        WITH RECURSIVE pohon_hierarki AS (
            -- Base case: pilih node yang diminta
            SELECT id, nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun, status
            FROM tb_pohon_kinerja 
            WHERE id = ?
            
            UNION ALL
            
            -- Recursive case: ambil semua child nodes
            SELECT pk.id, pk.nama_pohon, pk.parent, pk.jenis_pohon, pk.level_pohon, pk.kode_opd, pk.keterangan, pk.tahun, pk.status
            FROM tb_pohon_kinerja pk
            INNER JOIN pohon_hierarki ph ON pk.parent = ph.id
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
            ph.status,
            i.id as indikator_id,
            i.indikator as nama_indikator,
            t.id as target_id,
            t.target as target_value,
            t.satuan as target_satuan,
            pp.id as pelaksana_id,
            pp.pegawai_id
        FROM 
            pohon_hierarki ph
        LEFT JOIN 
            tb_indikator i ON ph.id = i.pokin_id
        LEFT JOIN 
            tb_target t ON i.id = t.indikator_id
        LEFT JOIN 
            tb_pelaksana_pokin pp ON ph.id = pp.pohon_kinerja_id
        ORDER BY 
            ph.level_pohon, ph.id, i.id, t.id, pp.id
    `

	rows, err := tx.QueryContext(ctx, script, idPokin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map untuk menyimpan pohon kinerja yang sudah diproses
	pokinMap := make(map[int]domain.PohonKinerja)
	indikatorMap := make(map[string]domain.Indikator)

	for rows.Next() {
		var (
			pokinId, parent, levelPohon                                    int
			namaPohon, jenisPohon, kodeOpd, keterangan, tahunPokin, status string
			indikatorId, namaIndikator                                     sql.NullString
			targetId, targetValue, targetSatuan                            sql.NullString
			pelaksanaId, pegawaiId                                         sql.NullString
		)

		err := rows.Scan(
			&pokinId, &namaPohon, &parent, &jenisPohon, &levelPohon,
			&kodeOpd, &keterangan, &tahunPokin, &status,
			&indikatorId, &namaIndikator,
			&targetId, &targetValue, &targetSatuan,
			&pelaksanaId, &pegawaiId,
		)
		if err != nil {
			return nil, err
		}

		// Proses Pohon Kinerja
		pokin, exists := pokinMap[pokinId]
		if !exists {
			pokin = domain.PohonKinerja{
				Id:         pokinId,
				NamaPohon:  namaPohon,
				Parent:     parent,
				JenisPohon: jenisPohon,
				LevelPohon: levelPohon,
				KodeOpd:    kodeOpd,
				Keterangan: keterangan,
				Tahun:      tahunPokin,
				Status:     status,
			}
		}

		// Proses Pelaksana jika ada
		if pelaksanaId.Valid && pegawaiId.Valid {
			pelaksana := domain.PelaksanaPokin{
				Id:        pelaksanaId.String,
				PegawaiId: pegawaiId.String,
			}
			// Cek duplikasi pelaksana
			isDuplicate := false
			for _, p := range pokin.Pelaksana {
				if p.Id == pelaksana.Id {
					isDuplicate = true
					break
				}
			}
			if !isDuplicate {
				pokin.Pelaksana = append(pokin.Pelaksana, pelaksana)
			}
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
			pokin.Indikator = append(pokin.Indikator, indikator)
		}

		pokinMap[pokinId] = pokin
	}

	// Konversi map ke slice
	var result []domain.PohonKinerja
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

func (repository *PohonKinerjaRepositoryImpl) FindIndikatorByPokinId(ctx context.Context, tx *sql.Tx, pokinId string) ([]domain.Indikator, error) {
	script := "SELECT id, pokin_id, indikator, tahun FROM tb_indikator WHERE pokin_id = ?"
	rows, err := tx.QueryContext(ctx, script, pokinId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indikators []domain.Indikator
	for rows.Next() {
		var indikator domain.Indikator
		err := rows.Scan(&indikator.Id, &indikator.PokinId, &indikator.Indikator, &indikator.Tahun)
		if err != nil {
			return nil, err
		}
		indikators = append(indikators, indikator)
	}
	return indikators, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindTargetByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string) ([]domain.Target, error) {
	script := "SELECT id, indikator_id, target, satuan, tahun FROM tb_target WHERE indikator_id = ?"
	rows, err := tx.QueryContext(ctx, script, indikatorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var targets []domain.Target
	for rows.Next() {
		var target domain.Target
		err := rows.Scan(&target.Id, &target.IndikatorId, &target.Target, &target.Satuan, &target.Tahun)
		if err != nil {
			return nil, err
		}
		targets = append(targets, target)
	}
	return targets, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindPokinToClone(ctx context.Context, tx *sql.Tx, id int) (domain.PohonKinerja, error) {
	script := "SELECT id, nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun FROM tb_pohon_kinerja WHERE id = ?"
	rows, err := tx.QueryContext(ctx, script, id)
	if err != nil {
		return domain.PohonKinerja{}, fmt.Errorf("gagal memeriksa data yang akan di-clone: %v", err)
	}
	defer rows.Close()

	var existingPokin domain.PohonKinerja
	if rows.Next() {
		err := rows.Scan(
			&existingPokin.Id,
			&existingPokin.NamaPohon,
			&existingPokin.Parent,
			&existingPokin.JenisPohon,
			&existingPokin.LevelPohon,
			&existingPokin.KodeOpd,
			&existingPokin.Keterangan,
			&existingPokin.Tahun,
		)
		if err != nil {
			return domain.PohonKinerja{}, fmt.Errorf("gagal membaca data yang akan di-clone: %v", err)
		}
		return existingPokin, nil
	}
	return domain.PohonKinerja{}, fmt.Errorf("data dengan ID %d tidak ditemukan", id)
}

func (repository *PohonKinerjaRepositoryImpl) ValidateParentLevel(ctx context.Context, tx *sql.Tx, parentId int, levelPohon int) error {
	script := "SELECT level_pohon FROM tb_pohon_kinerja WHERE id = ?"
	var parentLevel int
	err := tx.QueryRowContext(ctx, script, parentId).Scan(&parentLevel)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("parent dengan ID %d tidak ditemukan", parentId)
		}
		return fmt.Errorf("gagal memeriksa level parent: %v", err)
	}

	if levelPohon == 5 && parentLevel != 4 {
		return fmt.Errorf("untuk level pohon 5, parent harus memiliki level 4")
	} else if levelPohon == 6 && parentLevel != 5 {
		return fmt.Errorf("untuk level pohon 6, parent harus memiliki level 5")
	}

	return nil
}

func (repository *PohonKinerjaRepositoryImpl) FindIndikatorToClone(ctx context.Context, tx *sql.Tx, pokinId int) ([]domain.Indikator, error) {
	script := "SELECT id, indikator, tahun FROM tb_indikator WHERE pokin_id = ?"
	rows, err := tx.QueryContext(ctx, script, pokinId)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data indikator: %v", err)
	}
	defer rows.Close()

	var indikators []domain.Indikator
	for rows.Next() {
		var indikator domain.Indikator
		err := rows.Scan(&indikator.Id, &indikator.Indikator, &indikator.Tahun)
		if err != nil {
			return nil, fmt.Errorf("gagal membaca data indikator: %v", err)
		}
		indikators = append(indikators, indikator)
	}
	return indikators, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindTargetToClone(ctx context.Context, tx *sql.Tx, indikatorId string) ([]domain.Target, error) {
	script := "SELECT target, satuan, tahun FROM tb_target WHERE indikator_id = ?"
	rows, err := tx.QueryContext(ctx, script, indikatorId)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data target: %v", err)
	}
	defer rows.Close()

	var targets []domain.Target
	for rows.Next() {
		var target domain.Target
		err := rows.Scan(&target.Target, &target.Satuan, &target.Tahun)
		if err != nil {
			return nil, fmt.Errorf("gagal membaca data target: %v", err)
		}
		targets = append(targets, target)
	}
	return targets, nil
}

func (repository *PohonKinerjaRepositoryImpl) InsertClonedPokin(ctx context.Context, tx *sql.Tx, pokin domain.PohonKinerja) (int64, error) {
	script := "INSERT INTO tb_pohon_kinerja (nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, script,
		pokin.NamaPohon,
		pokin.Parent,
		pokin.JenisPohon,
		pokin.LevelPohon,
		pokin.KodeOpd,
		pokin.Keterangan,
		pokin.Tahun,
	)
	if err != nil {
		return 0, fmt.Errorf("gagal menyimpan data pohon kinerja yang di-clone: %v", err)
	}
	return result.LastInsertId()
}

func (repository *PohonKinerjaRepositoryImpl) InsertClonedIndikator(ctx context.Context, tx *sql.Tx, indikatorId string, pokinId int64, indikator domain.Indikator) error {
	script := "INSERT INTO tb_indikator (id, pokin_id, indikator, tahun) VALUES (?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, script, indikatorId, pokinId, indikator.Indikator, indikator.Tahun)
	if err != nil {
		return fmt.Errorf("gagal menyimpan indikator baru: %v", err)
	}
	return nil
}

func (repository *PohonKinerjaRepositoryImpl) InsertClonedTarget(ctx context.Context, tx *sql.Tx, targetId string, indikatorId string, target domain.Target) error {
	script := "INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, script, targetId, indikatorId, target.Target, target.Satuan, target.Tahun)
	if err != nil {
		return fmt.Errorf("gagal menyimpan target baru: %v", err)
	}
	return nil
}

func (repository *PohonKinerjaRepositoryImpl) FindPokinByJenisPohon(ctx context.Context, tx *sql.Tx, jenisPohon string, levelPohon int, tahun string, kodeOpd string) ([]domain.PohonKinerja, error) {
	script := "SELECT id, nama_pohon, jenis_pohon, level_pohon, kode_opd, tahun FROM tb_pohon_kinerja WHERE 1=1"
	parameters := []interface{}{}
	if jenisPohon != "" {
		script += " AND jenis_pohon = ?"
		parameters = append(parameters, jenisPohon)
	}
	if levelPohon != 0 {
		script += " AND level_pohon = ?"
		parameters = append(parameters, levelPohon)
	}
	if kodeOpd != "" {
		script += " AND kode_opd = ?"
		parameters = append(parameters, kodeOpd)
	}
	if tahun != "" {
		script += " AND tahun = ?"
		parameters = append(parameters, tahun)
	}
	script += " ORDER BY nama_pohon asc"

	rows, err := tx.QueryContext(ctx, script, parameters...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pokins []domain.PohonKinerja
	for rows.Next() {
		var pokin domain.PohonKinerja
		err := rows.Scan(&pokin.Id, &pokin.NamaPohon, &pokin.JenisPohon, &pokin.LevelPohon, &pokin.KodeOpd, &pokin.Tahun)
		if err != nil {
			return nil, err
		}
		pokins = append(pokins, pokin)
	}
	return pokins, nil
}

// Tambahkan fungsi helper untuk mendapatkan child nodes
func (repository *PohonKinerjaRepositoryImpl) GetChildNodes(ctx context.Context, tx *sql.Tx, parentId int) ([]domain.PohonKinerja, error) {
	script := `
        SELECT 
            id, nama_pohon, parent, jenis_pohon, level_pohon, 
            kode_opd, keterangan, tahun
        FROM 
            tb_pohon_kinerja 
        WHERE 
            parent = ?
        ORDER BY 
            id ASC
    `
	rows, err := tx.QueryContext(ctx, script, parentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var children []domain.PohonKinerja
	for rows.Next() {
		var child domain.PohonKinerja
		err := rows.Scan(
			&child.Id, &child.NamaPohon, &child.Parent,
			&child.JenisPohon, &child.LevelPohon,
			&child.KodeOpd, &child.Keterangan, &child.Tahun,
		)
		if err != nil {
			return nil, err
		}
		children = append(children, child)
	}
	return children, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindPokinByPelaksana(ctx context.Context, tx *sql.Tx, pegawaiId string, tahun string) ([]domain.PohonKinerja, error) {
	script := `
        SELECT DISTINCT
            pk.id,
            pk.nama_pohon,
            pk.parent,
            pk.jenis_pohon,
            pk.level_pohon,
            pk.kode_opd,
            pk.keterangan,
            pk.tahun,
            pk.created_at,
            pp.id as pelaksana_id,
            pp.pegawai_id
        FROM 
            tb_pohon_kinerja pk
        INNER JOIN 
            tb_pelaksana_pokin pp ON pk.id = pp.pohon_kinerja_id
        WHERE 
            pp.pegawai_id = ?
            AND pk.tahun = ?
        ORDER BY 
            pk.level_pohon, pk.id, pk.created_at ASC
    `

	rows, err := tx.QueryContext(ctx, script, pegawaiId, tahun)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data pohon kinerja: %v", err)
	}
	defer rows.Close()

	pokinMap := make(map[int]domain.PohonKinerja)

	for rows.Next() {
		var pokin domain.PohonKinerja
		var pelaksana domain.PelaksanaPokin

		err := rows.Scan(
			&pokin.Id,
			&pokin.NamaPohon,
			&pokin.Parent,
			&pokin.JenisPohon,
			&pokin.LevelPohon,
			&pokin.KodeOpd,
			&pokin.Keterangan,
			&pokin.Tahun,
			&pokin.CreatedAt,
			&pelaksana.Id,
			&pelaksana.PegawaiId,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal scan data pohon kinerja: %v", err)
		}

		// Cek apakah pohon kinerja sudah ada di map
		existingPokin, exists := pokinMap[pokin.Id]
		if exists {
			// Jika sudah ada, tambahkan pelaksana baru ke slice pelaksana yang ada
			existingPokin.Pelaksana = append(existingPokin.Pelaksana, pelaksana)
			pokinMap[pokin.Id] = existingPokin
		} else {
			// Jika belum ada, buat entry baru dengan pelaksana pertama
			pokin.Pelaksana = []domain.PelaksanaPokin{pelaksana}
			pokinMap[pokin.Id] = pokin
		}
	}

	// Konversi map ke slice untuk hasil akhir
	var result []domain.PohonKinerja
	for _, pokin := range pokinMap {
		result = append(result, pokin)
	}

	if len(result) == 0 {
		return nil, sql.ErrNoRows
	}

	return result, nil
}
