package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"errors"
)

type PeriodeRepositoryImpl struct {
}

func NewPeriodeRepositoryImpl() *PeriodeRepositoryImpl {
	return &PeriodeRepositoryImpl{}
}

func (repository *PeriodeRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, periode domain.Periode) (domain.Periode, error) {
	query := "INSERT INTO tb_periode(id, tahun_awal, tahun_akhir) VALUES (?, ?, ?)"
	_, err := tx.ExecContext(ctx, query, periode.Id, periode.TahunAwal, periode.TahunAkhir)
	if err != nil {
		return periode, err
	}
	return periode, nil
}

func (repository *PeriodeRepositoryImpl) IsIdExists(ctx context.Context, tx *sql.Tx, id int) bool {
	query := "SELECT COUNT(*) FROM tb_periode WHERE id = ?"
	var count int
	err := tx.QueryRowContext(ctx, query, id).Scan(&count)
	if err != nil {
		return true // Assume exists on error to be safe
	}
	return count > 0
}

func (repository *PeriodeRepositoryImpl) SaveTahunPeriode(ctx context.Context, tx *sql.Tx, tahunPeriode domain.TahunPeriode) error {
	query := "INSERT INTO tb_tahun_periode(id_periode, tahun) VALUES (?, ?)"
	_, err := tx.ExecContext(ctx, query, tahunPeriode.IdPeriode, tahunPeriode.Tahun)
	return err
}

func (repository *PeriodeRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, periodeId int) (domain.Periode, error) {
	query := "SELECT id, tahun_awal, tahun_akhir FROM tb_periode WHERE id = ?"
	rows, err := tx.QueryContext(ctx, query, periodeId)
	if err != nil {
		return domain.Periode{}, err
	}
	defer rows.Close()

	periode := domain.Periode{}
	if rows.Next() {
		err := rows.Scan(&periode.Id, &periode.TahunAwal, &periode.TahunAkhir)
		if err != nil {
			return periode, err
		}
		return periode, nil
	}

	return periode, errors.New("periode not found")
}

func (repository *PeriodeRepositoryImpl) FindOverlappingPeriodes(ctx context.Context, tx *sql.Tx, tahunAwal, tahunAkhir string) ([]domain.Periode, error) {
	query := `
		SELECT id, tahun_awal, tahun_akhir 
		FROM tb_periode 
		WHERE (tahun_awal <= ? AND tahun_akhir >= ?) 
		   OR (tahun_awal <= ? AND tahun_akhir >= ?)
		   OR (tahun_awal >= ? AND tahun_akhir <= ?)`

	rows, err := tx.QueryContext(ctx, query, tahunAkhir, tahunAwal, tahunAkhir, tahunAwal, tahunAwal, tahunAkhir)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var periodes []domain.Periode
	for rows.Next() {
		periode := domain.Periode{}
		err := rows.Scan(&periode.Id, &periode.TahunAwal, &periode.TahunAkhir)
		if err != nil {
			return nil, err
		}
		periodes = append(periodes, periode)
	}
	return periodes, nil
}

func (repository *PeriodeRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, periode domain.Periode) (domain.Periode, error) {
	query := "UPDATE tb_periode SET tahun_awal = ?, tahun_akhir = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, periode.TahunAwal, periode.TahunAkhir, periode.Id)
	if err != nil {
		return periode, err
	}
	return periode, nil
}

func (repository *PeriodeRepositoryImpl) DeleteTahunPeriode(ctx context.Context, tx *sql.Tx, periodeId int) error {
	query := "DELETE FROM tb_tahun_periode WHERE id_periode = ?"
	_, err := tx.ExecContext(ctx, query, periodeId)
	return err
}

func (repository *PeriodeRepositoryImpl) FindByTahun(ctx context.Context, tx *sql.Tx, tahun string) (domain.Periode, error) {
	query := `
		SELECT p.id, p.tahun_awal, p.tahun_akhir 
		FROM tb_periode p
		JOIN tb_tahun_periode tp ON p.id = tp.id_periode
		WHERE tp.tahun = ?
		LIMIT 1`

	rows, err := tx.QueryContext(ctx, query, tahun)
	if err != nil {
		return domain.Periode{}, err
	}
	defer rows.Close()

	if rows.Next() {
		periode := domain.Periode{}
		err := rows.Scan(&periode.Id, &periode.TahunAwal, &periode.TahunAkhir)
		if err != nil {
			return periode, err
		}
		return periode, nil
	}

	return domain.Periode{}, errors.New("periode not found")
}

func (repository *PeriodeRepositoryImpl) FindOverlappingPeriodesExcludeCurrent(ctx context.Context, tx *sql.Tx, currentId int, tahunAwal, tahunAkhir string) ([]domain.Periode, error) {
	query := `
		SELECT id, tahun_awal, tahun_akhir 
		FROM tb_periode 
		WHERE id != ? AND (
			(tahun_awal <= ? AND tahun_akhir >= ?) 
			OR (tahun_awal <= ? AND tahun_akhir >= ?)
			OR (tahun_awal >= ? AND tahun_akhir <= ?)
		)`

	rows, err := tx.QueryContext(ctx, query, currentId, tahunAkhir, tahunAwal, tahunAkhir, tahunAwal, tahunAwal, tahunAkhir)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var periodes []domain.Periode
	for rows.Next() {
		periode := domain.Periode{}
		err := rows.Scan(&periode.Id, &periode.TahunAwal, &periode.TahunAkhir)
		if err != nil {
			return nil, err
		}
		periodes = append(periodes, periode)
	}
	return periodes, nil
}
