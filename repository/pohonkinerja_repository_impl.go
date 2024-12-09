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
	scriptPokin := `
        UPDATE tb_pohon_kinerja 
        SET nama_pohon = ?, 
            parent = CASE 
                WHEN clone_from = 0 THEN ? 
                ELSE parent 
            END,
            jenis_pohon = ?, 
            level_pohon = ?, 
            kode_opd = ?, 
            keterangan = ?, 
            tahun = ?
        WHERE id = ?`

	_, err := tx.ExecContext(ctx, scriptPokin,
		pohonKinerja.NamaPohon,
		pohonKinerja.Parent,
		pohonKinerja.JenisPohon,
		pohonKinerja.LevelPohon,
		pohonKinerja.KodeOpd,
		pohonKinerja.Keterangan,
		pohonKinerja.Tahun,
		pohonKinerja.Id)
	if err != nil {
		return pohonKinerja, err
	}

	// Update pelaksana
	scriptDeletePelaksana := "DELETE FROM tb_pelaksana_pokin WHERE pohon_kinerja_id = ?"
	_, err = tx.ExecContext(ctx, scriptDeletePelaksana, fmt.Sprint(pohonKinerja.Id))
	if err != nil {
		return pohonKinerja, err
	}

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

	// Proses indikator
	for _, indikator := range pohonKinerja.Indikator {
		// Update atau insert indikator dengan clone_from
		scriptUpdateIndikator := `
			INSERT INTO tb_indikator (id, pokin_id, indikator, tahun, clone_from) 
			VALUES (?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE 
				indikator = VALUES(indikator),
				tahun = VALUES(tahun),
				clone_from = VALUES(clone_from)`

		_, err := tx.ExecContext(ctx, scriptUpdateIndikator,
			indikator.Id,
			pohonKinerja.Id,
			indikator.Indikator,
			indikator.Tahun,
			indikator.CloneFrom)
		if err != nil {
			return pohonKinerja, err
		}

		// Hapus target lama untuk indikator ini
		scriptDeleteTargets := "DELETE FROM tb_target WHERE indikator_id = ?"
		_, err = tx.ExecContext(ctx, scriptDeleteTargets, indikator.Id)
		if err != nil {
			return pohonKinerja, err
		}

		// Insert target baru dengan clone_from
		for _, target := range indikator.Target {
			// Log untuk debugging
			fmt.Printf("Inserting target: ID=%s, IndikatorID=%s, CloneFrom=%s\n",
				target.Id, target.IndikatorId, target.CloneFrom)

			scriptInsertTarget := `
				INSERT INTO tb_target 
					(id, indikator_id, target, satuan, tahun, clone_from)
				VALUES 
					(?, ?, ?, ?, ?, ?)`

			_, err := tx.ExecContext(ctx, scriptInsertTarget,
				target.Id,
				target.IndikatorId,
				target.Target,
				target.Satuan,
				target.Tahun,
				target.CloneFrom) // Pastikan clone_from dimasukkan
			if err != nil {
				return pohonKinerja, fmt.Errorf("error inserting target: %v", err)
			}
		}
	}

	// Hapus indikator yang tidak ada dalam request
	var existingIndikatorIds []string
	scriptGetExisting := "SELECT id FROM tb_indikator WHERE pokin_id = ?"
	rows, err := tx.QueryContext(ctx, scriptGetExisting, pohonKinerja.Id)
	if err != nil {
		return pohonKinerja, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return pohonKinerja, err
		}
		existingIndikatorIds = append(existingIndikatorIds, id)
	}

	// Buat map untuk indikator baru
	newIndikatorIds := make(map[string]bool)
	for _, ind := range pohonKinerja.Indikator {
		newIndikatorIds[ind.Id] = true
	}

	// Hapus indikator yang tidak ada dalam request
	for _, existingId := range existingIndikatorIds {
		if !newIndikatorIds[existingId] {
			// Hapus target terlebih dahulu
			scriptDeleteTargets := "DELETE FROM tb_target WHERE indikator_id = ?"
			_, err = tx.ExecContext(ctx, scriptDeleteTargets, existingId)
			if err != nil {
				return pohonKinerja, err
			}

			// Kemudian hapus indikator
			scriptDeleteIndikator := "DELETE FROM tb_indikator WHERE id = ?"
			_, err = tx.ExecContext(ctx, scriptDeleteIndikator, existingId)
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
        WITH RECURSIVE pohon_hierarki AS (
            -- Base case: ambil semua node dari OPD yang diminta dan OPD lain yang memiliki parent dari OPD yang diminta
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
                (pk.kode_opd = ? 
                OR pk.parent IN (
                    SELECT id FROM tb_pohon_kinerja 
                    WHERE kode_opd = ? AND tahun = ?
                ))
                AND pk.tahun = ?
				AND pk.status NOT IN ('menunggu_disetujui', 'tarik pokin opd', 'disetujui', 'ditolak')
            
            UNION 
            
            -- Recursive case: ambil parent nodes
            SELECT 
                p.id,
                p.nama_pohon,
                p.parent,
                p.jenis_pohon,
                p.level_pohon,
                p.kode_opd,
                p.keterangan,
                p.tahun,
                p.created_at,
                p.status
            FROM 
                tb_pohon_kinerja p
            INNER JOIN 
                pohon_hierarki ph ON p.id = ph.parent
            WHERE 
                p.tahun = ?
				AND p.status NOT IN ('menunggu_disetujui', 'tarik pokin opd', 'disetujui', 'ditolak','crosscutting_menunggu')
        )
        SELECT DISTINCT 
            h.id,
            h.nama_pohon,
            h.parent,
            h.jenis_pohon,
            h.level_pohon,
            h.kode_opd,
            h.keterangan,
            h.tahun,
            h.created_at,
            h.status,
            p.id as pelaksana_id,
            p.pegawai_id
        FROM 
            pohon_hierarki h
        LEFT JOIN 
            tb_pelaksana_pokin p ON h.id = p.pohon_kinerja_id
        ORDER BY 
            h.level_pohon, h.id, h.created_at ASC`

	rows, err := tx.QueryContext(ctx, script, kodeOpd, kodeOpd, tahun, tahun, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pokinMap := make(map[int]domain.PohonKinerja)

	for rows.Next() {
		var pokin domain.PohonKinerja
		var pelaksana domain.PelaksanaPokin
		var pelaksanaID, pegawaiID sql.NullString

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
			&pelaksanaID,
			&pegawaiID,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal scan data pohon kinerja: %v", err)
		}

		existingPokin, exists := pokinMap[pokin.Id]
		if exists {
			if pelaksanaID.Valid && pegawaiID.Valid {
				pelaksana.Id = pelaksanaID.String
				pelaksana.PegawaiId = pegawaiID.String
				existingPokin.Pelaksana = append(existingPokin.Pelaksana, pelaksana)
				pokinMap[pokin.Id] = existingPokin
			}
		} else {
			if pelaksanaID.Valid && pegawaiID.Valid {
				pelaksana.Id = pelaksanaID.String
				pelaksana.PegawaiId = pegawaiID.String
				pokin.Pelaksana = []domain.PelaksanaPokin{pelaksana}
			} else {
				pokin.Pelaksana = []domain.PelaksanaPokin{}
			}
			pokinMap[pokin.Id] = pokin
		}
	}

	var result []domain.PohonKinerja
	for _, pokin := range pokinMap {
		result = append(result, pokin)
	}

	if len(result) == 0 {
		return nil, sql.ErrNoRows
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
	// Update tb_pohon_kinerja dengan mempertahankan status
	scriptPokin := `
        UPDATE tb_pohon_kinerja 
        SET nama_pohon = ?, 
            parent = CASE 
                WHEN clone_from = 0 THEN ? 
                ELSE parent 
            END,
            jenis_pohon = ?, 
            level_pohon = ?, 
            kode_opd = ?, 
            keterangan = ?, 
            tahun = ?
        WHERE id = ?`

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

	// Update pelaksana
	scriptDeletePelaksana := "DELETE FROM tb_pelaksana_pokin WHERE pohon_kinerja_id = ?"
	_, err = tx.ExecContext(ctx, scriptDeletePelaksana, fmt.Sprint(pokinAdmin.Id))
	if err != nil {
		return pokinAdmin, err
	}

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

	// Proses indikator
	for _, indikator := range pokinAdmin.Indikator {
		// Update atau insert indikator
		scriptUpdateIndikator := `
			INSERT INTO tb_indikator (id, pokin_id, indikator, tahun, clone_from) 
			VALUES (?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE 
				indikator = VALUES(indikator),
				tahun = VALUES(tahun)`

		_, err := tx.ExecContext(ctx, scriptUpdateIndikator,
			indikator.Id,
			pokinAdmin.Id,
			indikator.Indikator,
			indikator.Tahun,
			indikator.CloneFrom)
		if err != nil {
			return pokinAdmin, err
		}

		// Hapus target lama untuk indikator ini
		scriptDeleteTargets := "DELETE FROM tb_target WHERE indikator_id = ?"
		_, err = tx.ExecContext(ctx, scriptDeleteTargets, indikator.Id)
		if err != nil {
			return pokinAdmin, err
		}

		// Insert target baru
		for _, target := range indikator.Target {
			scriptInsertTarget := `
				INSERT INTO tb_target (id, indikator_id, target, satuan, tahun, clone_from)
				VALUES (?, ?, ?, ?, ?, ?)`

			_, err := tx.ExecContext(ctx, scriptInsertTarget,
				target.Id,
				indikator.Id,
				target.Target,
				target.Satuan,
				target.Tahun,
				target.CloneFrom)
			if err != nil {
				return pokinAdmin, err
			}
		}
	}

	// Hapus indikator yang tidak ada dalam request
	var existingIndikatorIds []string
	scriptGetExisting := "SELECT id FROM tb_indikator WHERE pokin_id = ?"
	rows, err := tx.QueryContext(ctx, scriptGetExisting, pokinAdmin.Id)
	if err != nil {
		return pokinAdmin, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return pokinAdmin, err
		}
		existingIndikatorIds = append(existingIndikatorIds, id)
	}

	// Buat map untuk indikator baru
	newIndikatorIds := make(map[string]bool)
	for _, ind := range pokinAdmin.Indikator {
		newIndikatorIds[ind.Id] = true
	}

	// Hapus indikator yang tidak ada dalam request
	for _, existingId := range existingIndikatorIds {
		if !newIndikatorIds[existingId] {
			// Hapus target terlebih dahulu
			scriptDeleteTargets := "DELETE FROM tb_target WHERE indikator_id = ?"
			_, err = tx.ExecContext(ctx, scriptDeleteTargets, existingId)
			if err != nil {
				return pokinAdmin, err
			}

			// Kemudian hapus indikator
			scriptDeleteIndikator := "DELETE FROM tb_indikator WHERE id = ?"
			_, err = tx.ExecContext(ctx, scriptDeleteIndikator, existingId)
			if err != nil {
				return pokinAdmin, err
			}
		}
	}

	return pokinAdmin, nil
}

func (repository *PohonKinerjaRepositoryImpl) DeletePokinAdmin(ctx context.Context, tx *sql.Tx, id int) error {
	// Query untuk mendapatkan semua ID yang akan dihapus
	findIdsScript := `
        WITH RECURSIVE pohon_hierarki AS (
            -- Base case: node yang akan dihapus
            SELECT id, parent, level_pohon, clone_from 
            FROM tb_pohon_kinerja 
            WHERE id = ?
            
            UNION ALL
            
            -- Recursive case: child nodes dan data clone
            SELECT pk.id, pk.parent, pk.level_pohon, pk.clone_from
            FROM tb_pohon_kinerja pk
            INNER JOIN pohon_hierarki ph ON 
                -- Ambil child nodes langsung
                pk.parent = ph.id OR 
                -- Jika data asli, ambil yang mengclone-nya
                (ph.clone_from = 0 AND pk.clone_from = ph.id)
        ),
        clone_hierarki AS (
            -- Base case: data yang mengclone dan data yang parent-nya terhubung dengan id yang dihapus
            SELECT id, parent, level_pohon, clone_from
            FROM tb_pohon_kinerja
            WHERE clone_from IN (SELECT id FROM pohon_hierarki)
            OR parent IN (SELECT id FROM pohon_hierarki)
            
            UNION ALL
            
            -- Recursive case: child nodes dari data clone
            SELECT pk.id, pk.parent, pk.level_pohon, pk.clone_from
            FROM tb_pohon_kinerja pk
            INNER JOIN clone_hierarki ch ON 
                pk.parent = ch.id
        ),
        parent_hierarki AS (
            -- Ambil data yang parent-nya adalah id yang akan dihapus
            SELECT id, parent, level_pohon, clone_from
            FROM tb_pohon_kinerja
            WHERE parent = ?
        )
        SELECT id FROM pohon_hierarki
        UNION
        SELECT id FROM clone_hierarki
        UNION
        SELECT id FROM parent_hierarki;
    `

	rows, err := tx.QueryContext(ctx, findIdsScript, id, id)
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
	scriptDeleteTarget := fmt.Sprintf(`
        DELETE FROM tb_target 
        WHERE indikator_id IN (
            SELECT id FROM tb_indikator 
            WHERE pokin_id IN (%s)
        )`, inClause)
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
	script := "SELECT id, nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun, status FROM tb_pohon_kinerja WHERE id = ?"
	rows, err := tx.QueryContext(ctx, script, id)
	if err != nil {
		return domain.PohonKinerja{}, err
	}
	defer rows.Close()

	pokinAdmin := domain.PohonKinerja{}
	if rows.Next() {
		rows.Scan(&pokinAdmin.Id, &pokinAdmin.NamaPohon, &pokinAdmin.Parent, &pokinAdmin.JenisPohon, &pokinAdmin.LevelPohon, &pokinAdmin.KodeOpd, &pokinAdmin.Keterangan, &pokinAdmin.Tahun, &pokinAdmin.Status)
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

	// Gunakan map untuk melacak indikator yang sudah diproses
	processedIndikators := make(map[string]bool)

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
			// Cek apakah indikator sudah diproses
			if !processedIndikators[indikatorId.String] {
				processedIndikators[indikatorId.String] = true

				indikator := domain.Indikator{
					Id:        indikatorId.String,
					PokinId:   fmt.Sprint(pokinId),
					Indikator: namaIndikator.String,
					Tahun:     tahunPokin,
				}

				// Gunakan map untuk melacak target yang unik
				processedTargets := make(map[string]bool)

				// Proses Target jika ada
				if targetId.Valid && targetValue.Valid && targetSatuan.Valid {
					if !processedTargets[targetId.String] {
						processedTargets[targetId.String] = true
						target := domain.Target{
							Id:          targetId.String,
							IndikatorId: indikatorId.String,
							Target:      targetValue.String,
							Satuan:      targetSatuan.String,
							Tahun:       tahunPokin,
						}
						indikator.Target = append(indikator.Target, target)
					}
				}

				pokin.Indikator = append(pokin.Indikator, indikator)
			}
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
	script := `
        SELECT i.id, i.pokin_id, i.indikator, i.tahun, i.clone_from,
               t.id, t.indikator_id, t.target, t.satuan, t.tahun, t.clone_from
        FROM tb_indikator i
        LEFT JOIN tb_target t ON i.id = t.indikator_id
        WHERE i.pokin_id = ?`

	rows, err := tx.QueryContext(ctx, script, pokinId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	indikatorMap := make(map[string]*domain.Indikator)

	for rows.Next() {
		var indId, pokinId, indikator, indTahun, indCloneFrom string
		var targetId, indikatorId, target, satuan, targetTahun sql.NullString
		var targetCloneFrom sql.NullString

		err := rows.Scan(
			&indId, &pokinId, &indikator, &indTahun, &indCloneFrom,
			&targetId, &indikatorId, &target, &satuan, &targetTahun, &targetCloneFrom)
		if err != nil {
			return nil, err
		}

		// Proses Indikator
		ind, exists := indikatorMap[indId]
		if !exists {
			ind = &domain.Indikator{
				Id:        indId,
				Indikator: indikator,
				Tahun:     indTahun,
				CloneFrom: indCloneFrom,
				Target:    []domain.Target{},
			}
			indikatorMap[indId] = ind
		}

		// Proses Target jika ada
		if targetId.Valid && indikatorId.Valid {
			target := domain.Target{
				Id:          targetId.String,
				IndikatorId: indikatorId.String,
				Target:      target.String,
				Satuan:      satuan.String,
				Tahun:       targetTahun.String,
			}
			if targetCloneFrom.Valid {
				target.CloneFrom = targetCloneFrom.String
			}
			ind.Target = append(ind.Target, target)
		}
	}

	// Convert map to slice
	var result []domain.Indikator
	for _, ind := range indikatorMap {
		result = append(result, *ind)
	}

	return result, nil
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
	script := "SELECT id, nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun, status FROM tb_pohon_kinerja WHERE id = ?"
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
			&existingPokin.Status,
		)
		if err != nil {
			return domain.PohonKinerja{}, fmt.Errorf("gagal membaca data yang akan di-clone: %v", err)
		}
		return existingPokin, nil
	}
	return domain.PohonKinerja{}, fmt.Errorf("data dengan ID %d tidak ditemukan", id)
}

func (repository *PohonKinerjaRepositoryImpl) ValidateParentLevel(ctx context.Context, tx *sql.Tx, parentId int, levelPohon int) error {
	// Validasi dasar: level tidak boleh kurang dari 4
	if levelPohon < 4 {
		return fmt.Errorf("level pohon tidak boleh kurang dari 4")
	}

	// Untuk level 4, parent bisa memiliki level 0 hingga 3
	if levelPohon == 4 {
		if parentId == 0 {
			return nil
		}
		// Cek level parentnya
		script := "SELECT level_pohon FROM tb_pohon_kinerja WHERE id = ?"
		var parentLevel int
		err := tx.QueryRowContext(ctx, script, parentId).Scan(&parentLevel)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("parent dengan ID %d tidak ditemukan", parentId)
			}
			return fmt.Errorf("gagal memeriksa level parent: %v", err)
		}

		// Validasi level parent untuk level 4
		if parentLevel < 0 || parentLevel > 3 {
			return fmt.Errorf("level pohon 4 harus memiliki parent dengan level 0 hingga 3, bukan level %d", parentLevel)
		}
		return nil
	}

	// Untuk level > 4, parent tidak boleh 0
	if parentId == 0 {
		return fmt.Errorf("level pohon %d harus memiliki parent", levelPohon)
	}

	// Cek level parent untuk level > 4
	script := "SELECT level_pohon FROM tb_pohon_kinerja WHERE id = ?"
	var parentLevel int
	err := tx.QueryRowContext(ctx, script, parentId).Scan(&parentLevel)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("parent dengan ID %d tidak ditemukan", parentId)
		}
		return fmt.Errorf("gagal memeriksa level parent: %v", err)
	}

	// Validasi: level parent harus tepat 1 tingkat di atas level saat ini
	expectedParentLevel := levelPohon - 1
	if parentLevel != expectedParentLevel {
		return fmt.Errorf("level pohon %d harus memiliki parent dengan level %d, bukan level %d",
			levelPohon, expectedParentLevel, parentLevel)
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
	script := "SELECT id, target, satuan, tahun FROM tb_target WHERE indikator_id = ?"
	rows, err := tx.QueryContext(ctx, script, indikatorId)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data target: %v", err)
	}
	defer rows.Close()

	var targets []domain.Target
	for rows.Next() {
		var target domain.Target
		err := rows.Scan(&target.Id, &target.Target, &target.Satuan, &target.Tahun)
		if err != nil {
			return nil, fmt.Errorf("gagal membaca data target: %v", err)
		}
		targets = append(targets, target)
	}
	return targets, nil
}

func (repository *PohonKinerjaRepositoryImpl) InsertClonedPokin(ctx context.Context, tx *sql.Tx, pokin domain.PohonKinerja) (int64, error) {
	script := `INSERT INTO tb_pohon_kinerja 
        (nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun, status, clone_from) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := tx.ExecContext(ctx, script,
		pokin.NamaPohon,
		pokin.Parent,
		pokin.JenisPohon,
		pokin.LevelPohon,
		pokin.KodeOpd,
		pokin.Keterangan,
		pokin.Tahun,
		pokin.Status,
		pokin.CloneFrom,
	)
	if err != nil {
		return 0, fmt.Errorf("gagal menyimpan data pohon kinerja yang di-clone: %v", err)
	}
	return result.LastInsertId()
}

func (repository *PohonKinerjaRepositoryImpl) InsertClonedIndikator(ctx context.Context, tx *sql.Tx, indikatorId string, pokinId int64, indikator domain.Indikator) error {
	script := "INSERT INTO tb_indikator (id, pokin_id, indikator, tahun, clone_from) VALUES (?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, script,
		indikatorId,
		pokinId,
		indikator.Indikator,
		indikator.Tahun,
		indikator.Id, // Id indikator asli sebagai clone_from
	)
	if err != nil {
		return fmt.Errorf("gagal menyimpan indikator baru: %v", err)
	}
	return nil
}

func (repository *PohonKinerjaRepositoryImpl) InsertClonedTarget(ctx context.Context, tx *sql.Tx, targetId string, indikatorId string, target domain.Target) error {
	fmt.Printf("Inserting target with clone_from: %s\n", target.Id) // Log sementara
	script := "INSERT INTO tb_target (id, indikator_id, target, satuan, tahun, clone_from) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, script,
		targetId,
		indikatorId,
		target.Target,
		target.Satuan,
		target.Tahun,
		target.Id, // Id target asli sebagai clone_from
	)
	if err != nil {
		return fmt.Errorf("gagal menyimpan target baru: %v", err)
	}
	return nil
}

func (repository *PohonKinerjaRepositoryImpl) FindPokinByJenisPohon(ctx context.Context, tx *sql.Tx, jenisPohon string, levelPohon int, tahun string, kodeOpd string, status string) ([]domain.PohonKinerja, error) {
	script := "SELECT id, nama_pohon, jenis_pohon, level_pohon, kode_opd, tahun, status FROM tb_pohon_kinerja WHERE 1=1"
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
	if status != "" {
		script += " AND status = ?"
		parameters = append(parameters, status)
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
		err := rows.Scan(&pokin.Id, &pokin.NamaPohon, &pokin.JenisPohon, &pokin.LevelPohon, &pokin.KodeOpd, &pokin.Tahun, &pokin.Status)
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

func (repository *PohonKinerjaRepositoryImpl) FindPokinByStatus(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string, status string) ([]domain.PohonKinerja, error) {
	SQL := `SELECT id, nama_pohon, kode_opd, tahun, jenis_pohon, level_pohon, parent, status 
            FROM tb_pohon_kinerja 
            WHERE kode_opd = ? AND tahun = ? AND status = ?`

	rows, err := tx.QueryContext(ctx, SQL, kodeOpd, tahun, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pokins []domain.PohonKinerja
	for rows.Next() {
		pokin := domain.PohonKinerja{}
		err := rows.Scan(
			&pokin.Id,
			&pokin.NamaPohon,
			&pokin.KodeOpd,
			&pokin.Tahun,
			&pokin.JenisPohon,
			&pokin.LevelPohon,
			&pokin.Parent,
			&pokin.Status,
		)
		if err != nil {
			return nil, err
		}
		pokins = append(pokins, pokin)
	}
	return pokins, nil
}

func (repository *PohonKinerjaRepositoryImpl) UpdatePokinStatus(ctx context.Context, tx *sql.Tx, id int, status string) error {
	script := "UPDATE tb_pohon_kinerja SET status = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, status, id)
	if err != nil {
		return fmt.Errorf("gagal mengupdate status: %v", err)
	}
	return nil
}

func (repository *PohonKinerjaRepositoryImpl) CheckPokinStatus(ctx context.Context, tx *sql.Tx, id int) (string, error) {
	script := "SELECT status FROM tb_pohon_kinerja WHERE id = ?"
	var status string
	err := tx.QueryRowContext(ctx, script, id).Scan(&status)
	if err != nil {
		return "", fmt.Errorf("gagal mengecek status: %v", err)
	}
	return status, nil
}

func (repository *PohonKinerjaRepositoryImpl) InsertClonedPokinWithStatus(ctx context.Context, tx *sql.Tx, pokin domain.PohonKinerja) (int64, error) {
	script := `INSERT INTO tb_pohon_kinerja 
        (nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun, status, clone_from) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := tx.ExecContext(ctx, script,
		pokin.NamaPohon,
		pokin.Parent,
		pokin.JenisPohon,
		pokin.LevelPohon,
		pokin.KodeOpd,
		pokin.Keterangan,
		pokin.Tahun,
		pokin.Status,
		pokin.CloneFrom,
	)
	if err != nil {
		return 0, fmt.Errorf("gagal menyimpan data pohon kinerja yang di-clone: %v", err)
	}
	return result.LastInsertId()
}

func (repository *PohonKinerjaRepositoryImpl) UpdatePokinStatusTolak(ctx context.Context, tx *sql.Tx, id int, status string) error {
	script := "UPDATE tb_pohon_kinerja SET status = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, status, id)
	if err != nil {
		return fmt.Errorf("gagal mengupdate status dan alasan: %v", err)
	}
	return nil
}

func (repository *PohonKinerjaRepositoryImpl) CheckCloneFrom(ctx context.Context, tx *sql.Tx, id int) (int, error) {
	script := "SELECT COALESCE(clone_from, 0) FROM tb_pohon_kinerja WHERE id = ?"
	var cloneFrom int
	err := tx.QueryRowContext(ctx, script, id).Scan(&cloneFrom)
	if err != nil {
		return 0, fmt.Errorf("gagal mengecek clone_from: %v", err)
	}
	return cloneFrom, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindPokinByCloneFrom(ctx context.Context, tx *sql.Tx, cloneFromId int) ([]domain.PohonKinerja, error) {
	script := "SELECT id, parent, nama_pohon, jenis_pohon, level_pohon, kode_opd, keterangan, tahun, status, clone_from FROM tb_pohon_kinerja WHERE clone_from = ?"
	rows, err := tx.QueryContext(ctx, script, cloneFromId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pokins []domain.PohonKinerja
	for rows.Next() {
		var pokin domain.PohonKinerja
		err := rows.Scan(
			&pokin.Id,
			&pokin.Parent,
			&pokin.NamaPohon,
			&pokin.JenisPohon,
			&pokin.LevelPohon,
			&pokin.KodeOpd,
			&pokin.Keterangan,
			&pokin.Tahun,
			&pokin.Status,
			&pokin.CloneFrom,
		)
		if err != nil {
			return nil, err
		}
		pokins = append(pokins, pokin)
	}
	return pokins, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindIndikatorByCloneFrom(ctx context.Context, tx *sql.Tx, pokinId int, cloneFromId string) (domain.Indikator, error) {
	script := "SELECT id, indikator, tahun FROM tb_indikator WHERE pokin_id = ? AND clone_from = ?"
	var indikator domain.Indikator
	err := tx.QueryRowContext(ctx, script, pokinId, cloneFromId).Scan(
		&indikator.Id,
		&indikator.Indikator,
		&indikator.Tahun,
	)
	if err != nil {
		return domain.Indikator{}, err
	}
	return indikator, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindTargetByCloneFrom(ctx context.Context, tx *sql.Tx, indikatorId string, cloneFromId string) (domain.Target, error) {
	script := "SELECT id, target, satuan, tahun FROM tb_target WHERE indikator_id = ? AND clone_from = ?"
	var target domain.Target
	err := tx.QueryRowContext(ctx, script, indikatorId, cloneFromId).Scan(
		&target.Id,
		&target.Target,
		&target.Satuan,
		&target.Tahun,
	)
	if err != nil {
		return domain.Target{}, err
	}
	return target, nil
}

// Tambahkan method baru untuk FindPokinByCrosscuttingStatus
func (repository *PohonKinerjaRepositoryImpl) FindPokinByCrosscuttingStatus(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.PohonKinerja, error) {
	script := `SELECT 
        id, nama_pohon, parent, jenis_pohon, level_pohon, 
        kode_opd, keterangan, tahun, status 
        FROM tb_pohon_kinerja 
        WHERE kode_opd = ? 
        AND tahun = ? 
        AND status = 'crosscutting_menunggu'
        ORDER BY level_pohon, id ASC`

	rows, err := tx.QueryContext(ctx, script, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pokins []domain.PohonKinerja
	for rows.Next() {
		var pokin domain.PohonKinerja
		err := rows.Scan(
			&pokin.Id, &pokin.NamaPohon, &pokin.Parent, &pokin.JenisPohon, &pokin.LevelPohon, &pokin.KodeOpd, &pokin.Keterangan, &pokin.Tahun, &pokin.Status,
		)
		if err != nil {
			return nil, err
		}
		pokins = append(pokins, pokin)
	}
	return pokins, nil
}

func (repository *PohonKinerjaRepositoryImpl) DeleteClonedPokinHierarchy(ctx context.Context, tx *sql.Tx, id int) error {
	// Query untuk mendapatkan hierarki dari data clone
	findIdsScript := `
        WITH RECURSIVE clone_hierarki AS (
            -- Base case: node clone yang akan dihapus
            SELECT id, parent, level_pohon, clone_from 
            FROM tb_pohon_kinerja 
            WHERE id = ?
            
            UNION ALL
            
            -- Recursive case: child nodes dari data clone
            SELECT pk.id, pk.parent, pk.level_pohon, pk.clone_from
            FROM tb_pohon_kinerja pk
            INNER JOIN clone_hierarki ch ON 
                pk.parent = ch.id
        )
        SELECT id FROM clone_hierarki;
    `

	rows, err := tx.QueryContext(ctx, findIdsScript, id)
	if err != nil {
		return fmt.Errorf("gagal mengambil hierarki clone: %v", err)
	}
	defer rows.Close()

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
	scriptDeleteTarget := fmt.Sprintf(`
        DELETE FROM tb_target 
        WHERE indikator_id IN (
            SELECT id FROM tb_indikator 
            WHERE pokin_id IN (%s)
        )`, inClause)
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
