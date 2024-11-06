package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"log"
)

type SubKegiatanRepositoryImpl struct {
}

func NewSubKegiatanRepositoryImpl() *SubKegiatanRepositoryImpl {
	return &SubKegiatanRepositoryImpl{}
}

func (repository *SubKegiatanRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, subKegiatan domain.SubKegiatan) (domain.SubKegiatan, error) {
	script := `INSERT INTO tb_subkegiatan (id, kode_subkegiatan, nama_subkegiatan, kode_opd, tahun) VALUES (?, ?, ?, ?, ?)`

	_, err := tx.ExecContext(ctx, script, subKegiatan.Id, subKegiatan.KodeSubKegiatan, subKegiatan.NamaSubKegiatan, subKegiatan.KodeOpd, subKegiatan.Tahun)
	if err != nil {
		return domain.SubKegiatan{}, err
	}

	return subKegiatan, nil
}

func (repository *SubKegiatanRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, subKegiatan domain.SubKegiatan) (domain.SubKegiatan, error) {
	script := `UPDATE tb_subkegiatan SET nama_subkegiatan = ?, kode_opd = ?, tahun = ? WHERE id = ?`

	_, err := tx.ExecContext(ctx, script, subKegiatan.NamaSubKegiatan, subKegiatan.KodeOpd, subKegiatan.Tahun, subKegiatan.Id)
	if err != nil {
		log.Printf("Error updating subkegiatan: %v", err)
		return domain.SubKegiatan{}, err
	}

	return subKegiatan, nil
}

func (repository *SubKegiatanRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string, pegawaiId string) ([]domain.SubKegiatan, error) {
	script := `SELECT id, kode_subkegiatan, nama_subkegiatan, kode_opd, tahun, created_at FROM tb_subkegiatan WHERE 1=1`
	var params []interface{}

	if kodeOpd != "" {
		script += ` AND kode_opd = ?`
		params = append(params, kodeOpd)
	}
	if pegawaiId != "" {
		script += ` AND pegawai_id = ?`
		params = append(params, pegawaiId)
	}
	script += ` ORDER BY created_at ASC`

	rows, err := tx.QueryContext(ctx, script, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subKegiatans []domain.SubKegiatan
	for rows.Next() {
		subKegiatan := domain.SubKegiatan{}
		err := rows.Scan(&subKegiatan.Id, &subKegiatan.KodeSubKegiatan, &subKegiatan.NamaSubKegiatan, &subKegiatan.KodeOpd, &subKegiatan.Tahun, &subKegiatan.CreatedAt)
		if err != nil {
			return nil, err
		}
		subKegiatans = append(subKegiatans, subKegiatan)
	}

	return subKegiatans, nil
}

func (repository *SubKegiatanRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, subKegiatanId string) (domain.SubKegiatan, error) {
	script := `SELECT id, kode_subkegiatan, pegawai_id, nama_subkegiatan, kode_opd, tahun FROM tb_subkegiatan WHERE id = ?`

	rows, err := tx.QueryContext(ctx, script, subKegiatanId)
	if err != nil {
		return domain.SubKegiatan{}, err
	}

	subKegiatan := domain.SubKegiatan{}
	if rows.Next() {
		rows.Scan(&subKegiatan.Id, &subKegiatan.KodeSubKegiatan, &subKegiatan.PegawaiId, &subKegiatan.NamaSubKegiatan, &subKegiatan.KodeOpd, &subKegiatan.Tahun)
	}

	return subKegiatan, nil
}

func (repository *SubKegiatanRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, subKegiatanId string) error {
	script := `DELETE FROM tb_subkegiatan WHERE id = ?`
	_, err := tx.ExecContext(ctx, script, subKegiatanId)
	return err
}
