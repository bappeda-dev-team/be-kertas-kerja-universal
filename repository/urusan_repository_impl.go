package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
)

type UrusanRepositoryImpl struct {
}

func NewUrusanRepositoryImpl() *UrusanRepositoryImpl {
	return &UrusanRepositoryImpl{}
}

func (repository *UrusanRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, urusan domainmaster.Urusan) (domainmaster.Urusan, error) {
	script := "INSERT INTO urusan(kode_urusan, nama_urusan) VALUES (?, ?)"
	_, err := tx.ExecContext(ctx, script, urusan.KodeUrusan, urusan.NamaUrusan)
	if err != nil {
		return urusan, err
	}

	return urusan, nil
}

func (repository *UrusanRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, urusan domainmaster.Urusan) (domainmaster.Urusan, error) {
	script := "UPDATE urusan SET kode_urusan = ?, nama_urusan = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, urusan.KodeUrusan, urusan.NamaUrusan, urusan.Id)
	if err != nil {
		return urusan, err
	}

	return urusan, nil
}

func (repository *UrusanRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domainmaster.Urusan, error) {
	script := "SELECT id, kode_urusan, nama_urusan, created_at FROM urusan"
	rows, err := tx.QueryContext(ctx, script)
	if err != nil {
		return []domainmaster.Urusan{}, err
	}

	defer rows.Close()

	var urusans []domainmaster.Urusan
	for rows.Next() {
		urusan := domainmaster.Urusan{}
		rows.Scan(&urusan.Id, &urusan.KodeUrusan, &urusan.NamaUrusan, &urusan.CreatedAt)
		urusans = append(urusans, urusan)
	}

	return urusans, nil
}

func (repository *UrusanRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (domainmaster.Urusan, error) {
	script := "SELECT id, kode_urusan, nama_urusan, created_at FROM urusan WHERE id = ?"
	rows, err := tx.QueryContext(ctx, script, id)
	if err != nil {
		return domainmaster.Urusan{}, err
	}

	defer rows.Close()

	urusan := domainmaster.Urusan{}
	rows.Scan(&urusan.Id, &urusan.KodeUrusan, &urusan.NamaUrusan, &urusan.CreatedAt)

	return urusan, nil
}

func (repository *UrusanRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	script := "DELETE FROM urusan WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}

	return nil
}
