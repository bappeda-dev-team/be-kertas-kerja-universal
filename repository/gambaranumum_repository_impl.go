package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type GambaranUmumRepositoryImpl struct {
}

func NewGambaranUmumRepositoryImpl() *GambaranUmumRepositoryImpl {
	return &GambaranUmumRepositoryImpl{}
}

func (repository *GambaranUmumRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, gambaranUmum domain.GambaranUmum) (domain.GambaranUmum, error) {
	query := "INSERT INTO tb_gambaran_umum (id, rekin_id, kode_opd, urutan, gambaran_umum) VALUES (?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, query, gambaranUmum.Id, gambaranUmum.RekinId, gambaranUmum.KodeOpd, gambaranUmum.Urutan, gambaranUmum.GambaranUmum)
	if err != nil {
		return domain.GambaranUmum{}, err
	}
	return gambaranUmum, nil
}

func (repository *GambaranUmumRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, gambaranUmum domain.GambaranUmum) (domain.GambaranUmum, error) {
	query := "UPDATE tb_gambaran_umum SET gambaran_umum = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, gambaranUmum.GambaranUmum, gambaranUmum.Id)
	if err != nil {
		return domain.GambaranUmum{}, err
	}
	return gambaranUmum, nil
}

func (repository *GambaranUmumRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	query := "DELETE FROM tb_gambaran_umum WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *GambaranUmumRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (domain.GambaranUmum, error) {
	query := "SELECT id, rekin_id, kode_opd, urutan, gambaran_umum FROM tb_gambaran_umum WHERE id = ? ORDER BY urutan ASC"
	row := tx.QueryRowContext(ctx, query, id)
	var gambaranUmum domain.GambaranUmum
	err := row.Scan(&gambaranUmum.Id, &gambaranUmum.RekinId, &gambaranUmum.KodeOpd, &gambaranUmum.Urutan, &gambaranUmum.GambaranUmum)
	if err != nil {
		return domain.GambaranUmum{}, err
	}
	return gambaranUmum, nil
}

func (repository *GambaranUmumRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, rekinId string) ([]domain.GambaranUmum, error) {
	query := "SELECT id, rekin_id, kode_opd, urutan, gambaran_umum FROM tb_gambaran_umum WHERE rekin_id = ? ORDER BY urutan ASC"
	rows, err := tx.QueryContext(ctx, query, rekinId)
	if err != nil {
		return []domain.GambaranUmum{}, err
	}
	defer rows.Close()

	var gambaranUmumList []domain.GambaranUmum
	for rows.Next() {
		var gambaranUmum domain.GambaranUmum
		err := rows.Scan(&gambaranUmum.Id, &gambaranUmum.RekinId, &gambaranUmum.KodeOpd, &gambaranUmum.Urutan, &gambaranUmum.GambaranUmum)
		if err != nil {
			return []domain.GambaranUmum{}, err
		}

		gambaranUmumList = append(gambaranUmumList, gambaranUmum)
	}

	err = rows.Err()
	if err != nil {
		return []domain.GambaranUmum{}, err
	}

	return gambaranUmumList, nil
}

func (repository *GambaranUmumRepositoryImpl) GetLastUrutanByRekinId(ctx context.Context, tx *sql.Tx, rekinId string) (int, error) {
	SQL := "SELECT COALESCE(MAX(urutan), 0) FROM tb_gambaran_umum WHERE rekin_id = ?"
	var lastUrutan int
	err := tx.QueryRowContext(ctx, SQL, rekinId).Scan(&lastUrutan)
	if err != nil {
		return 0, err
	}
	return lastUrutan, nil
}
