package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"errors"
)

type InovasiRepositoryImpl struct {
}

func NewInovasiRepositoryImpl() *InovasiRepositoryImpl {
	return &InovasiRepositoryImpl{}
}

func (repository *InovasiRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, inovasi domain.Inovasi) (domain.Inovasi, error) {
	script := "INSERT INTO tb_inovasi (id, rekin_id, kode_opd, pegawai_id, judul_inovasi, jenis_inovasi, gambaran_nilai_kebaruan) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, script, inovasi.Id, inovasi.RekinId, inovasi.KodeOpd, inovasi.PegawaiId, inovasi.JudulInovasi, inovasi.JenisInovasi, inovasi.GambaranNilaiKebaruan)
	if err != nil {
		return inovasi, err
	}
	return inovasi, nil
}

func (repository *InovasiRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, inovasiId string) (domain.Inovasi, error) {
	script := "SELECT id, rekin_id, kode_opd, pegawai_id, judul_inovasi, jenis_inovasi, gambaran_nilai_kebaruan FROM tb_inovasi WHERE id = ?"
	rows, err := tx.QueryContext(ctx, script, inovasiId)
	if err != nil {
		return domain.Inovasi{}, err
	}
	defer rows.Close()

	inovasi := domain.Inovasi{}
	if rows.Next() {
		err := rows.Scan(&inovasi.Id, &inovasi.RekinId, &inovasi.KodeOpd, &inovasi.PegawaiId, &inovasi.JudulInovasi, &inovasi.JenisInovasi, &inovasi.GambaranNilaiKebaruan)
		if err != nil {
			return domain.Inovasi{}, err
		}
		return inovasi, nil
	}
	return domain.Inovasi{}, sql.ErrNoRows
}

func (repository *InovasiRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, rekinId string, pegawaiId string) ([]domain.Inovasi, error) {
	script := "SELECT id, rekin_id, kode_opd, pegawai_id, judul_inovasi, jenis_inovasi, gambaran_nilai_kebaruan, created_at FROM tb_inovasi WHERE rekin_id = ? AND pegawai_id = ? ORDER BY created_at ASC"
	rows, err := tx.QueryContext(ctx, script, rekinId, pegawaiId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inovasis []domain.Inovasi
	for rows.Next() {
		inovasi := domain.Inovasi{}
		err := rows.Scan(&inovasi.Id, &inovasi.RekinId, &inovasi.KodeOpd, &inovasi.PegawaiId, &inovasi.JudulInovasi, &inovasi.JenisInovasi, &inovasi.GambaranNilaiKebaruan, &inovasi.CreatedAt)
		if err != nil {
			return nil, err
		}
		inovasis = append(inovasis, inovasi)
	}
	return inovasis, nil
}

func (repository *InovasiRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, inovasi domain.Inovasi) (domain.Inovasi, error) {
	script := "UPDATE tb_inovasi SET rekin_id = ?, judul_inovasi = ?, jenis_inovasi = ?, gambaran_nilai_kebaruan = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, inovasi.RekinId, inovasi.JudulInovasi, inovasi.JenisInovasi, inovasi.GambaranNilaiKebaruan, inovasi.Id)
	if err != nil {
		return inovasi, err
	}
	return inovasi, nil
}

func (repository *InovasiRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, inovasiId string) error {
	script := "DELETE FROM tb_inovasi WHERE id = ?"
	result, err := tx.ExecContext(ctx, script, inovasiId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("inovasi tidak ditemukan")
	}

	return nil
}
