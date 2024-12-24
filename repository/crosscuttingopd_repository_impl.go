package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

type CrosscuttingOpdRepositoryImpl struct {
}

func NewCrosscuttingOpdRepositoryImpl() *CrosscuttingOpdRepositoryImpl {
	return &CrosscuttingOpdRepositoryImpl{}
}

func (repository *CrosscuttingOpdRepositoryImpl) CreateCrosscutting(ctx context.Context, tx *sql.Tx, pokin domain.PohonKinerja, parentId int) (domain.PohonKinerja, error) {
	scriptPokin := "INSERT INTO tb_pohon_kinerja (nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun, status) VALUES (?, 0, ?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, scriptPokin,
		pokin.NamaPohon,
		pokin.JenisPohon,
		pokin.LevelPohon,
		pokin.KodeOpd,
		pokin.Keterangan,
		pokin.Tahun,
		pokin.Status)
	if err != nil {
		return pokin, err
	}

	newPokinId, err := result.LastInsertId()
	if err != nil {
		return pokin, err
	}
	pokin.Id = int(newPokinId)

	scriptCrosscutting := "INSERT INTO tb_crosscutting (crosscutting_from, crosscutting_to, status) VALUES (?, ?, ?)"
	_, err = tx.ExecContext(ctx, scriptCrosscutting,
		parentId,
		newPokinId,
		pokin.Status)
	if err != nil {
		return pokin, err
	}

	for _, indikator := range pokin.Indikator {
		scriptIndikator := "INSERT INTO tb_indikator (id, pokin_id, indikator, tahun) VALUES (?, ?, ?, ?)"
		_, err := tx.ExecContext(ctx, scriptIndikator,
			indikator.Id,
			pokin.Id,
			indikator.Indikator,
			indikator.Tahun)
		if err != nil {
			return pokin, err
		}

		for _, target := range indikator.Target {
			scriptTarget := "INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)"
			_, err := tx.ExecContext(ctx, scriptTarget,
				target.Id,
				indikator.Id,
				target.Target,
				target.Satuan,
				target.Tahun)
			if err != nil {
				return pokin, err
			}
		}
	}

	return pokin, nil
}
func (repository *CrosscuttingOpdRepositoryImpl) FindAllCrosscutting(ctx context.Context, tx *sql.Tx, parentId int) ([]domain.PohonKinerja, error) {
	script := `
        SELECT p.id, p.nama_pohon, p.parent, p.jenis_pohon, p.level_pohon, 
               p.kode_opd, p.keterangan, p.tahun, p.status, p.pegawai_action,
               p.created_at
        FROM tb_pohon_kinerja p 
        JOIN tb_crosscutting c ON p.id = c.crosscutting_to 
        WHERE c.crosscutting_from = ?
    `
	rows, err := tx.QueryContext(ctx, script, parentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.PohonKinerja
	for rows.Next() {
		pokin := domain.PohonKinerja{}
		var pegawaiActionJSON sql.NullString
		err := rows.Scan(
			&pokin.Id,
			&pokin.NamaPohon,
			&pokin.Parent,
			&pokin.JenisPohon,
			&pokin.LevelPohon,
			&pokin.KodeOpd,
			&pokin.Keterangan,
			&pokin.Tahun,
			&pokin.Status,
			&pegawaiActionJSON,
			&pokin.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if pegawaiActionJSON.Valid {
			var pegawaiAction interface{}
			err = json.Unmarshal([]byte(pegawaiActionJSON.String), &pegawaiAction)
			if err != nil {
				return nil, err
			}
			pokin.PegawaiAction = pegawaiAction
		}

		result = append(result, pokin)
	}
	return result, nil
}

func (repository *CrosscuttingOpdRepositoryImpl) FindIndikatorByPokinId(ctx context.Context, tx *sql.Tx, pokinIds []int) ([]domain.Indikator, error) {
	// Cek jika array kosong
	if len(pokinIds) == 0 {
		return []domain.Indikator{}, nil
	}

	// Buat placeholder untuk IN clause
	placeholders := make([]string, len(pokinIds))
	for i := range pokinIds {
		placeholders[i] = "?"
	}

	query := "SELECT id, pokin_id, indikator FROM tb_indikator WHERE pokin_id IN (" + strings.Join(placeholders, ",") + ")"

	args := make([]interface{}, len(pokinIds))
	for i, id := range pokinIds {
		args[i] = id
	}

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indikators []domain.Indikator
	for rows.Next() {
		var indikator domain.Indikator
		var pokinId int
		err := rows.Scan(
			&indikator.Id,
			&pokinId,
			&indikator.Indikator,
		)
		if err != nil {
			return nil, err
		}
		indikator.PokinId = strconv.Itoa(pokinId)
		indikators = append(indikators, indikator)
	}
	return indikators, nil
}

func (repository *CrosscuttingOpdRepositoryImpl) FindTargetByIndikatorIds(ctx context.Context, tx *sql.Tx, indikatorIds []string) ([]domain.Target, error) {
	// Cek jika array kosong
	if len(indikatorIds) == 0 {
		return []domain.Target{}, nil
	}

	// Buat placeholder untuk IN clause
	placeholders := make([]string, len(indikatorIds))
	for i := range indikatorIds {
		placeholders[i] = "?"
	}

	query := "SELECT id, indikator_id, target, satuan FROM tb_target WHERE indikator_id IN (" + strings.Join(placeholders, ",") + ")"

	args := make([]interface{}, len(indikatorIds))
	for i, id := range indikatorIds {
		args[i] = id
	}

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var targets []domain.Target
	for rows.Next() {
		var target domain.Target
		err := rows.Scan(&target.Id, &target.IndikatorId, &target.Target, &target.Satuan)
		if err != nil {
			return nil, err
		}
		targets = append(targets, target)
	}
	return targets, nil
}

func (repository *CrosscuttingOpdRepositoryImpl) UpdateCrosscutting(ctx context.Context, tx *sql.Tx, pokin domain.PohonKinerja, parentId int) (domain.PohonKinerja, error) {
	// Query update pohon kinerja tanpa status
	scriptPokin := `
        UPDATE tb_pohon_kinerja 
        SET nama_pohon = ?,
            jenis_pohon = ?,
            level_pohon = ?,
            kode_opd = ?,
            keterangan = ?,
            tahun = ?
        WHERE id = ?
    `
	args := []interface{}{
		pokin.NamaPohon,
		pokin.JenisPohon,
		pokin.LevelPohon,
		pokin.KodeOpd,
		pokin.Keterangan,
		pokin.Tahun,
		pokin.Id,
	}

	_, err := tx.ExecContext(ctx, scriptPokin, args...)
	if err != nil {
		return pokin, err
	}

	// Hapus dan insert ulang indikator dan target
	scriptDeleteTarget := `
        DELETE t FROM tb_target t
        INNER JOIN tb_indikator i ON t.indikator_id = i.id
        WHERE i.pokin_id = ?
    `
	_, err = tx.ExecContext(ctx, scriptDeleteTarget, pokin.Id)
	if err != nil {
		return pokin, err
	}

	scriptDeleteIndikator := "DELETE FROM tb_indikator WHERE pokin_id = ?"
	_, err = tx.ExecContext(ctx, scriptDeleteIndikator, pokin.Id)
	if err != nil {
		return pokin, err
	}

	// Insert indikator dan target baru
	for _, indikator := range pokin.Indikator {
		scriptIndikator := "INSERT INTO tb_indikator (id, pokin_id, indikator, tahun) VALUES (?, ?, ?, ?)"
		_, err := tx.ExecContext(ctx, scriptIndikator,
			indikator.Id,
			pokin.Id,
			indikator.Indikator,
			indikator.Tahun)
		if err != nil {
			return pokin, err
		}

		for _, target := range indikator.Target {
			scriptTarget := "INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)"
			_, err := tx.ExecContext(ctx, scriptTarget,
				target.Id,
				indikator.Id,
				target.Target,
				target.Satuan,
				target.Tahun)
			if err != nil {
				return pokin, err
			}
		}
	}

	return pokin, nil
}
func (repository *CrosscuttingOpdRepositoryImpl) ValidateKodeOpdChange(ctx context.Context, tx *sql.Tx, id int) error {
	var status string
	err := tx.QueryRowContext(ctx, "SELECT status FROM tb_crosscutting WHERE crosscutting_to = ?", id).Scan(&status)
	if err != nil {
		return err
	}

	if status != "crosscutting_menunggu" {
		return errors.New("kode OPD hanya dapat diubah saat status crosscutting_menunggu")
	}

	return nil
}

func (repository *CrosscuttingOpdRepositoryImpl) DeleteCrosscutting(ctx context.Context, tx *sql.Tx, pokinId int, nipPegawai string) error {
	// Validasi status
	var currentStatus string
	query := `
        SELECT status FROM tb_pohon_kinerja 
        WHERE id = ?
    `
	err := tx.QueryRowContext(ctx, query, pokinId).Scan(&currentStatus)
	if err != nil {
		return err
	}

	if currentStatus != "crosscutting_disetujui" {
		return errors.New("crosscutting hanya dapat dihapus saat status crosscutting_disetujui")
	}

	// Buat pegawai_action
	currentTime := time.Now()
	pegawaiAction := map[string]interface{}{
		"reject_by": nipPegawai,
		"reject_at": currentTime,
	}

	pegawaiActionJSON, err := json.Marshal(pegawaiAction)
	if err != nil {
		return err
	}

	scriptUpdatePokin := `
        UPDATE tb_pohon_kinerja 
        SET parent = 0,
            status = 'crosscutting_ditolak',
            pegawai_action = ?
        WHERE id = ?
    `
	_, err = tx.ExecContext(ctx, scriptUpdatePokin, pegawaiActionJSON, pokinId)
	if err != nil {
		return err
	}

	// Update status di tb_crosscutting
	scriptUpdateCross := `
        UPDATE tb_crosscutting 
        SET status = 'crosscutting_ditolak'
        WHERE crosscutting_to = ?
    `
	result, err := tx.ExecContext(ctx, scriptUpdateCross, pokinId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows updated in tb_crosscutting, check crosscutting_to value")
	}

	return nil
}

func (repository *CrosscuttingOpdRepositoryImpl) ApproveOrRejectCrosscutting(ctx context.Context, tx *sql.Tx, crosscuttingId int, approve bool, pegawaiAction map[string]interface{}, levelPohon int, jenisPohon string, parentId int) error {
	var currentStatus string
	query := `
        SELECT status FROM tb_pohon_kinerja 
        WHERE id = ?
    `
	err := tx.QueryRowContext(ctx, query, crosscuttingId).Scan(&currentStatus)
	if err != nil {
		return err
	}

	if currentStatus != "crosscutting_menunggu" && currentStatus != "crosscutting_ditolak" {
		return errors.New("crosscutting sudah disetujui")
	}

	var newStatus string
	if approve {
		newStatus = "crosscutting_disetujui"
	} else {
		newStatus = "crosscutting_ditolak"
	}

	pegawaiActionJSON, err := json.Marshal(pegawaiAction)
	if err != nil {
		return err
	}

	if approve {
		scriptUpdateParent := `
            UPDATE tb_pohon_kinerja 
            SET parent = ?,
                status = ?,
                pegawai_action = ?,
                level_pohon = ?,
                jenis_pohon = ?
            WHERE id = ?
        `
		_, err = tx.ExecContext(ctx, scriptUpdateParent, parentId, newStatus, pegawaiActionJSON, levelPohon, jenisPohon, crosscuttingId)
		if err != nil {
			return err
		}
	} else {
		scriptPokin := `
            UPDATE tb_pohon_kinerja 
            SET status = ?,
                pegawai_action = ?
            WHERE id = ?
        `
		_, err = tx.ExecContext(ctx, scriptPokin, newStatus, pegawaiActionJSON, crosscuttingId)
		if err != nil {
			return err
		}
	}

	scriptCrosscutting := `
        UPDATE tb_crosscutting 
        SET status = ?
        WHERE crosscutting_to = ?
    `
	result, err := tx.ExecContext(ctx, scriptCrosscutting, newStatus, crosscuttingId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows updated in tb_crosscutting, check crosscutting_to value")
	}

	return nil
}

func (repository *CrosscuttingOpdRepositoryImpl) DeleteUnused(ctx context.Context, tx *sql.Tx, crosscuttingId int) error {
	// Cek apakah ada data yang bisa dihapus
	checkQuery := `
        SELECT COUNT(p.id) 
        FROM tb_pohon_kinerja p
        JOIN tb_crosscutting c ON p.id = c.crosscutting_to
        WHERE p.status IN ('crosscutting_menunggu', 'crosscutting_ditolak')
        AND c.crosscutting_to = ?
    `
	var count int
	err := tx.QueryRowContext(ctx, checkQuery, crosscuttingId).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("crosscutting tidak dapat dihapus")
	}

	query := `
        SELECT p.id 
        FROM tb_pohon_kinerja p
        JOIN tb_crosscutting c ON p.id = c.crosscutting_to
        WHERE p.status IN ('crosscutting_menunggu', 'crosscutting_ditolak')
        AND c.crosscutting_to = ?
    `

	rows, err := tx.QueryContext(ctx, query, crosscuttingId)
	if err != nil {
		return err
	}
	defer rows.Close()

	var pokinIds []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return err
		}
		pokinIds = append(pokinIds, id)
	}

	placeholders := make([]string, len(pokinIds))
	args := make([]interface{}, len(pokinIds))
	for i, id := range pokinIds {
		placeholders[i] = "?"
		args[i] = id
	}
	inClause := strings.Join(placeholders, ",")

	// Hapus pelaksana (perbaikan nama kolom dari pokin_id menjadi id_pokin)
	pelaksanaQuery := "DELETE FROM tb_pelaksana_pokin WHERE pohon_kinerja_id IN (" + inClause + ")"
	_, err = tx.ExecContext(ctx, pelaksanaQuery, args...)
	if err != nil {
		return err
	}

	// Hapus target
	targetQuery := `
        DELETE t FROM tb_target t
        INNER JOIN tb_indikator i ON t.indikator_id = i.id
        WHERE i.pokin_id IN (` + inClause + ")"
	_, err = tx.ExecContext(ctx, targetQuery, args...)
	if err != nil {
		return err
	}

	// Hapus indikator
	indikatorQuery := "DELETE FROM tb_indikator WHERE pokin_id IN (" + inClause + ")"
	_, err = tx.ExecContext(ctx, indikatorQuery, args...)
	if err != nil {
		return err
	}

	// Hapus crosscutting
	crosscuttingQuery := "DELETE FROM tb_crosscutting WHERE crosscutting_to IN (" + inClause + ")"
	_, err = tx.ExecContext(ctx, crosscuttingQuery, args...)
	if err != nil {
		return err
	}

	// Hapus pohon kinerja
	pokinQuery := "DELETE FROM tb_pohon_kinerja WHERE id IN (" + inClause + ")"
	_, err = tx.ExecContext(ctx, pokinQuery, args...)
	if err != nil {
		return err
	}

	return nil
}
