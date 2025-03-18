package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
)

type LembagaRepositoryImpl struct {
}

func NewLembagaRepositoryImpl() *LembagaRepositoryImpl {
	return &LembagaRepositoryImpl{}
}

func (repository *LembagaRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, lembaga domainmaster.Lembaga) domainmaster.Lembaga {
	script := "INSERT INTO tb_lembaga (id, kode_lembaga, nama_lembaga) VALUES (?, ?, ?)"
	_, err := tx.ExecContext(ctx, script, lembaga.Id, lembaga.KodeLembaga, lembaga.NamaLembaga)
	if err != nil {
		return lembaga
	}
	return lembaga
}

func (repository *LembagaRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, lembaga domainmaster.Lembaga) domainmaster.Lembaga {
	script := "UPDATE tb_lembaga SET kode_lembaga = ?, nama_lembaga = ?, is_active = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, lembaga.KodeLembaga, lembaga.NamaLembaga, lembaga.IsActive, lembaga.Id)
	if err != nil {
		return lembaga
	}
	return lembaga
}

func (repository *LembagaRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	script := "DELETE FROM tb_lembaga WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *LembagaRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (domainmaster.Lembaga, error) {
	script := "SELECT id, kode_lembaga, nama_lembaga, is_active FROM tb_lembaga WHERE id = ?"
	var lembaga domainmaster.Lembaga
	err := tx.QueryRowContext(ctx, script, id).Scan(&lembaga.Id, &lembaga.KodeLembaga, &lembaga.NamaLembaga, &lembaga.IsActive)
	if err != nil {
		return domainmaster.Lembaga{}, err
	}
	return lembaga, nil
}

func (r *LembagaRepositoryImpl) FindByKode(ctx context.Context, tx *sql.Tx, kodeLembaga string) (domainmaster.Lembaga, error) {
	script := "SELECT id, kode_lembaga, nama_lembaga, is_active FROM tb_lembaga WHERE kode_lembaga = ?"
	var lembaga domainmaster.Lembaga
	err := tx.QueryRowContext(ctx, script, kodeLembaga).Scan(&lembaga.Id, &lembaga.KodeLembaga, &lembaga.NamaLembaga, &lembaga.IsActive)
	if err != nil {
		return domainmaster.Lembaga{}, err
	}
	return lembaga, nil
}

func (repository *LembagaRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domainmaster.Lembaga, error) {
	script := "SELECT id, kode_lembaga, nama_lembaga, is_active FROM tb_lembaga"
	rows, err := tx.QueryContext(ctx, script)
	if err != nil {
		return []domainmaster.Lembaga{}, err
	}
	defer rows.Close()
	var lembagas []domainmaster.Lembaga
	for rows.Next() {
		lembaga := domainmaster.Lembaga{}
		err := rows.Scan(&lembaga.Id, &lembaga.KodeLembaga, &lembaga.NamaLembaga, &lembaga.IsActive)
		if err != nil {
			return []domainmaster.Lembaga{}, err
		}
		lembagas = append(lembagas, lembaga)
	}
	return lembagas, nil
}
