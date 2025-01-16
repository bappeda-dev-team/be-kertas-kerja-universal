package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
)

type RencanaAksiRepositoryImpl struct {
}

func NewRencanaAksiRepositoryImpl() *RencanaAksiRepositoryImpl {
	return &RencanaAksiRepositoryImpl{}
}

func (repository *RencanaAksiRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, rencanaAksi domain.RencanaAksi) (domain.RencanaAksi, error) {
	script := "INSERT INTO tb_rencana_aksi (id, rencana_kinerja_id, kode_opd, urutan, nama_rencana_aksi) VALUES (?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, script, rencanaAksi.Id, rencanaAksi.RencanaKinerjaId, rencanaAksi.KodeOpd, rencanaAksi.Urutan, rencanaAksi.NamaRencanaAksi)
	if err != nil {
		return domain.RencanaAksi{}, err
	}
	return rencanaAksi, nil
}

func (repository *RencanaAksiRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, rencanaAksi domain.RencanaAksi) (domain.RencanaAksi, error) {
	script := "UPDATE tb_rencana_aksi SET urutan = ?, nama_rencana_aksi = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, rencanaAksi.Urutan, rencanaAksi.NamaRencanaAksi, rencanaAksi.Id)
	if err != nil {
		return domain.RencanaAksi{}, err
	}
	return rencanaAksi, nil
}

func (repository *RencanaAksiRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	// Hapus terlebih dahulu dari tb_pelaksanaan_rencana_aksi
	scriptPelaksanaan := "DELETE FROM tb_pelaksanaan_rencana_aksi WHERE rencana_aksi_id = ?"
	_, err := tx.ExecContext(ctx, scriptPelaksanaan, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus data dari tb_pelaksanaan_rencana_aksi: %v", err)
	}

	// Kemudian hapus dari rencana_aksi
	scriptRencanaAksi := "DELETE FROM tb_rencana_aksi WHERE id = ?"
	_, err = tx.ExecContext(ctx, scriptRencanaAksi, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus data dari rencana_aksi: %v", err)
	}

	return nil
}

func (repository *RencanaAksiRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (domain.RencanaAksi, error) {
	script := "SELECT id, rencana_kinerja_id, kode_opd, urutan, nama_rencana_aksi FROM tb_rencana_aksi WHERE id = ?"
	row := tx.QueryRowContext(ctx, script, id)
	var rencanaAksi domain.RencanaAksi
	err := row.Scan(&rencanaAksi.Id, &rencanaAksi.RencanaKinerjaId, &rencanaAksi.KodeOpd, &rencanaAksi.Urutan, &rencanaAksi.NamaRencanaAksi)
	if err != nil {
		return domain.RencanaAksi{}, err
	}
	return rencanaAksi, nil
}

func (repository *RencanaAksiRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, rencanaKinerjaId string) ([]domain.RencanaAksi, error) {
	script := "SELECT id, rencana_kinerja_id, kode_opd, urutan, nama_rencana_aksi FROM tb_rencana_aksi"
	var args []interface{}

	if rencanaKinerjaId != "" {
		script += " WHERE rencana_kinerja_id = ?"
		args = append(args, rencanaKinerjaId)
	}

	script += " ORDER BY urutan ASC"

	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return []domain.RencanaAksi{}, err
	}
	defer rows.Close()

	var rencanaAksis []domain.RencanaAksi
	for rows.Next() {
		var rencanaAksi domain.RencanaAksi
		err := rows.Scan(&rencanaAksi.Id, &rencanaAksi.RencanaKinerjaId, &rencanaAksi.KodeOpd, &rencanaAksi.Urutan, &rencanaAksi.NamaRencanaAksi)
		if err != nil {
			return []domain.RencanaAksi{}, err
		}
		rencanaAksis = append(rencanaAksis, rencanaAksi)
	}

	return rencanaAksis, nil
}

func (repository *RencanaAksiRepositoryImpl) IsUrutanExistsForRencanaKinerja(ctx context.Context, tx *sql.Tx, rencanaKinerjaId string, urutan int) (bool, error) {
	script := "SELECT COUNT(*) FROM tb_rencana_aksi WHERE rencana_kinerja_id = ? AND urutan = ?"
	var count int
	err := tx.QueryRowContext(ctx, script, rencanaKinerjaId, urutan).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repository *RencanaAksiRepositoryImpl) IsUrutanExistsForRencanaKinerjaExcludingId(ctx context.Context, tx *sql.Tx, rencanaKinerjaId string, urutan int, excludeId string) (bool, error) {
	script := "SELECT COUNT(*) FROM tb_rencana_aksi WHERE rencana_kinerja_id = ? AND urutan = ? AND id != ?"
	var count int
	err := tx.QueryRowContext(ctx, script, rencanaKinerjaId, urutan, excludeId).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repository *RencanaAksiRepositoryImpl) GetTotalBobotForRencanaKinerja(ctx context.Context, tx *sql.Tx, rencanaKinerjaId string) (int, error) {
	script := `
		SELECT COALESCE(SUM(pra.bobot), 0)
		FROM tb_rencana_aksi ra
		JOIN tb_pelaksanaan_rencana_aksi pra ON ra.id = pra.rencana_aksi_id
		WHERE ra.rencana_kinerja_id = ?
	`
	var totalBobot int
	err := tx.QueryRowContext(ctx, script, rencanaKinerjaId).Scan(&totalBobot)
	return totalBobot, err
}
