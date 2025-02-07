package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type PermasalahanRekinRepositoryImpl struct {
}

func NewPermasalahanRekinRepositoryImpl() *PermasalahanRekinRepositoryImpl {
	return &PermasalahanRekinRepositoryImpl{}
}

func (repository *PermasalahanRekinRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, permasalahan domain.PermasalahanRekin) (domain.PermasalahanRekin, error) {
	script := "INSERT INTO tb_permasalahan (id, rekin_id, permasalahan, penyebab_internal, penyebab_eksternal, jenis_permasalahan) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, script, permasalahan.Id, permasalahan.RekinId, permasalahan.Permasalahan, permasalahan.PenyebabInternal, permasalahan.PenyebabEksternal, permasalahan.JenisPermasalahan)
	if err != nil {
		return domain.PermasalahanRekin{}, err
	}
	return permasalahan, nil
}

func (repository *PermasalahanRekinRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, permasalahan domain.PermasalahanRekin) (domain.PermasalahanRekin, error) {
	script := "UPDATE tb_permasalahan SET permasalahan = ?, penyebab_internal = ?, penyebab_eksternal = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, permasalahan.Permasalahan, permasalahan.PenyebabInternal, permasalahan.PenyebabEksternal, permasalahan.Id)
	if err != nil {
		return domain.PermasalahanRekin{}, err
	}
	return permasalahan, nil
}

func (repository *PermasalahanRekinRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	script := "DELETE FROM tb_permasalahan WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *PermasalahanRekinRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, rekinId *string) ([]domain.PermasalahanRekin, error) {
	script := "SELECT id, rekin_id, permasalahan, penyebab_internal, penyebab_eksternal, jenis_permasalahan, created_at FROM tb_permasalahan WHERE 1=1"

	var args []interface{}

	if rekinId != nil {
		script += " AND rekin_id = ?"
		args = append(args, *rekinId)
	}

	script += " order by created_at ascs"

	rows, err := tx.QueryContext(ctx, script, args...)
	if err != nil {
		return []domain.PermasalahanRekin{}, err
	}
	defer rows.Close()

	var permasalahanList []domain.PermasalahanRekin
	for rows.Next() {
		var permasalahan domain.PermasalahanRekin
		err := rows.Scan(&permasalahan.Id, &permasalahan.RekinId, &permasalahan.Permasalahan, &permasalahan.PenyebabInternal, &permasalahan.PenyebabEksternal, &permasalahan.JenisPermasalahan)
		if err != nil {
			return []domain.PermasalahanRekin{}, err
		}
		permasalahanList = append(permasalahanList, permasalahan)
	}
	return permasalahanList, nil
}

func (repository *PermasalahanRekinRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.PermasalahanRekin, error) {
	script := "SELECT id, rekin_id, permasalahan, penyebab_internal, penyebab_eksternal, jenis_permasalahan FROM tb_permasalahan WHERE id = ?"
	var permasalahan domain.PermasalahanRekin
	err := tx.QueryRowContext(ctx, script, id).Scan(&permasalahan.Id, &permasalahan.Permasalahan, &permasalahan.PenyebabInternal, &permasalahan.PenyebabEksternal, &permasalahan.JenisPermasalahan)
	if err != nil {
		return domain.PermasalahanRekin{}, err
	}
	return permasalahan, nil
}
