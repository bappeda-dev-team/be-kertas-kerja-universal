package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type PohonKinerjaOpdRepositoryImpl struct {
}

func NewPohonKinerjaOpdRepositoryImpl() *PohonKinerjaOpdRepositoryImpl {
	return &PohonKinerjaOpdRepositoryImpl{}
}

func (repository *PohonKinerjaOpdRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, pohonKinerja domain.PohonKinerja) (domain.PohonKinerja, error) {
	script := "INSERT INTO tb_pohon_kinerja (nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, script, pohonKinerja.NamaPohon, pohonKinerja.Parent, pohonKinerja.JenisPohon, pohonKinerja.LevelPohon, pohonKinerja.KodeOpd, pohonKinerja.Keterangan, pohonKinerja.Tahun)
	if err != nil {
		return pohonKinerja, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return pohonKinerja, err
	}

	pohonKinerja.Id = int(lastInsertId)
	return pohonKinerja, nil
}

func (repository *PohonKinerjaOpdRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, pohonKinerja domain.PohonKinerja) (domain.PohonKinerja, error) {
	script := "UPDATE tb_pohon_kinerja SET nama_pohon = ?, parent = ?, jenis_pohon = ?, level_pohon = ?, kode_opd = ?, keterangan = ?, tahun = ? WHERE id = ?"
	result, err := tx.ExecContext(ctx, script, pohonKinerja.NamaPohon, pohonKinerja.Parent, pohonKinerja.JenisPohon, pohonKinerja.LevelPohon, pohonKinerja.KodeOpd, pohonKinerja.Keterangan, pohonKinerja.Tahun, pohonKinerja.Id)
	if err != nil {
		return pohonKinerja, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return pohonKinerja, err
	}

	pohonKinerja.Id = int(lastInsertId)
	return pohonKinerja, nil
}

func (repository *PohonKinerjaOpdRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.PohonKinerja, error) {
	script := "SELECT id, parent, nama_pohon, jenis_pohon, level_pohon, kode_opd, keterangan, tahun FROM tb_pohon_kinerja WHERE id = ?"
	rows, err := tx.QueryContext(ctx, script, id)
	if err != nil {
		return domain.PohonKinerja{}, err
	}
	defer rows.Close()

	pohonKinerja := domain.PohonKinerja{}
	if rows.Next() {
		rows.Scan(&pohonKinerja.Id, &pohonKinerja.Parent, &pohonKinerja.NamaPohon, &pohonKinerja.JenisPohon, &pohonKinerja.LevelPohon, &pohonKinerja.KodeOpd, &pohonKinerja.Keterangan, &pohonKinerja.Tahun)
	}
	return pohonKinerja, nil
}

func (repository *PohonKinerjaOpdRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string) ([]domain.PohonKinerja, error) {
	script := "SELECT id, parent, nama_pohon, jenis_pohon, level_pohon, kode_opd, keterangan, tahun FROM tb_pohon_kinerja WHERE 1=1"
	params := []interface{}{}
	if kodeOpd != "" {
		script += " AND kode_opd = ?"
		params = append(params, kodeOpd)
	}
	if tahun != "" {
		script += " AND tahun = ?"
		params = append(params, tahun)
	}

	rows, err := tx.QueryContext(ctx, script, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pohonKinerjas []domain.PohonKinerja
	for rows.Next() {
		pohonKinerja := domain.PohonKinerja{}
		err := rows.Scan(&pohonKinerja.Id, &pohonKinerja.Parent, &pohonKinerja.NamaPohon, &pohonKinerja.JenisPohon, &pohonKinerja.LevelPohon, &pohonKinerja.KodeOpd, &pohonKinerja.Keterangan, &pohonKinerja.Tahun)
		if err != nil {
			return nil, err
		}
		pohonKinerjas = append(pohonKinerjas, pohonKinerja)
	}
	return pohonKinerjas, nil
}

func (repository *PohonKinerjaOpdRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	script := "DELETE FROM tb_pohon_kinerja WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}
	return nil
}
