package repository

import (
	"context"
	"database/sql"
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
	script := "INSERT INTO tb_pohon_kinerja (nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, script, pohonKinerja.NamaPohon, pohonKinerja.Parent, pohonKinerja.JenisPohon, pohonKinerja.LevelPohon, pohonKinerja.KodeOpd, pohonKinerja.Keterangan, pohonKinerja.Tahun)
	if err != nil {
		return pohonKinerja, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return pohonKinerja, err
	}

	pohonKinerja.Id = int(lastInsertId)
	return pohonKinerja, nil
}

func (repository *PohonKinerjaRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, pohonKinerja domain.PohonKinerja) (domain.PohonKinerja, error) {
	script := "UPDATE tb_pohon_kinerja SET nama_pohon = ?, parent = ?, jenis_pohon = ?, level_pohon = ?, kode_opd = ?, keterangan = ?, tahun = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, pohonKinerja.NamaPohon, pohonKinerja.Parent, pohonKinerja.JenisPohon, pohonKinerja.LevelPohon, pohonKinerja.KodeOpd, pohonKinerja.Keterangan, pohonKinerja.Tahun, pohonKinerja.Id)
	if err != nil {
		return pohonKinerja, err
	}
	return pohonKinerja, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.PohonKinerja, error) {
	script := "SELECT id, parent, nama_pohon, jenis_pohon, level_pohon, kode_opd, keterangan, tahun FROM tb_pohon_kinerja WHERE id = ?"
	rows, err := tx.QueryContext(ctx, script, id)
	if err != nil {
		return domain.PohonKinerja{}, err
	}
	defer rows.Close()

	pohonKinerja := domain.PohonKinerja{}
	if rows.Next() {
		rows.Scan(&pohonKinerja.Id, &pohonKinerja.Parent, &pohonKinerja.NamaPohon, &pohonKinerja.JenisPohon, &pohonKinerja.LevelPohon, &pohonKinerja.KodeOpd, &pohonKinerja.Keterangan, &pohonKinerja.Tahun)
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
            pk.created_at
        FROM 
            tb_pohon_kinerja pk
        WHERE 
            pk.kode_opd = ? 
            AND pk.tahun = ?
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
		)
		if err != nil {
			return nil, err
		}
		result = append(result, pokin)
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
	script := "DELETE FROM tb_pohon_kinerja WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}
	return nil
}

// admin pokin
func (repository *PohonKinerjaRepositoryImpl) CreatePokinAdmin(ctx context.Context, tx *sql.Tx, pokinAdmin domain.PohonKinerja) (domain.PohonKinerja, error) {
	// Insert pohon kinerja tanpa ID
	scriptPokin := "INSERT INTO tb_pohon_kinerja (nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, scriptPokin,
		pokinAdmin.NamaPohon, pokinAdmin.Parent, pokinAdmin.JenisPohon, pokinAdmin.LevelPohon, pokinAdmin.KodeOpd, pokinAdmin.Keterangan, pokinAdmin.Tahun)
	if err != nil {
		return pokinAdmin, err
	}

	// Dapatkan ID yang baru dibuat
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return pokinAdmin, err
	}
	pokinAdmin.Id = int(lastInsertId)

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
	scriptPokin := "UPDATE tb_pohon_kinerja SET nama_pohon = ?, parent = ?, jenis_pohon = ?, level_pohon = ?, kode_opd = ?, keterangan = ?, tahun = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, scriptPokin,
		pokinAdmin.NamaPohon,
		pokinAdmin.Parent,
		pokinAdmin.JenisPohon,
		pokinAdmin.LevelPohon,
		pokinAdmin.KodeOpd,
		pokinAdmin.Keterangan,
		pokinAdmin.Tahun,
		pokinAdmin.Id)
	if err != nil {
		return pokinAdmin, err
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
            SELECT id, nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun
            FROM tb_pohon_kinerja 
            WHERE id = ?
            
            UNION ALL
            
            -- Recursive case: ambil semua child nodes
            SELECT pk.id, pk.nama_pohon, pk.parent, pk.jenis_pohon, pk.level_pohon, pk.kode_opd, pk.keterangan, pk.tahun
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
            i.id as indikator_id,
            i.indikator as nama_indikator,
            t.id as target_id,
            t.target as target_value,
            t.satuan as target_satuan
        FROM 
            pohon_hierarki ph
        LEFT JOIN 
            tb_indikator i ON ph.id = i.pokin_id
        LEFT JOIN 
            tb_target t ON i.id = t.indikator_id
        ORDER BY 
            ph.level_pohon, ph.id, i.id, t.id
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
