package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"errors"
	"math/rand"
	"time"
)

type VisiPemdaRepositoryImpl struct {
}

func NewVisiPemdaRepositoryImpl() *VisiPemdaRepositoryImpl {
	return &VisiPemdaRepositoryImpl{}
}

func (repository *VisiPemdaRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, visiPemda domain.VisiPemda) (domain.VisiPemda, error) {
	// Generate random ID 4 digit
	rand.Seed(time.Now().UnixNano())
	visiPemda.Id = rand.Intn(9000) + 1000 // Generate angka antara 1000-9999

	script := "INSERT INTO tb_visi_pemda (id, visi, tahun_awal_periode, tahun_akhir_periode, jenis_periode, keterangan) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, script,
		visiPemda.Id,
		visiPemda.Visi,
		visiPemda.TahunAwalPeriode,
		visiPemda.TahunAkhirPeriode,
		visiPemda.JenisPeriode,
		visiPemda.Keterangan)
	if err != nil {
		return visiPemda, err
	}

	return visiPemda, nil
}

func (repository *VisiPemdaRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, visiPemda domain.VisiPemda) (domain.VisiPemda, error) {
	script := `UPDATE tb_visi_pemda SET 
        visi = ?, 
        tahun_awal_periode = ?, 
        tahun_akhir_periode = ?, 
        jenis_periode = ?, 
        keterangan = ? 
    WHERE id = ?`

	_, err := tx.ExecContext(ctx, script,
		visiPemda.Visi,
		visiPemda.TahunAwalPeriode,
		visiPemda.TahunAkhirPeriode,
		visiPemda.JenisPeriode,
		visiPemda.Keterangan,
		visiPemda.Id)
	if err != nil {
		return visiPemda, err
	}

	return visiPemda, nil
}

func (repository *VisiPemdaRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, visiPemdaId int) error {
	script := "DELETE FROM tb_visi_pemda WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, visiPemdaId)
	return err
}

func (repository *VisiPemdaRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, visiPemdaId int) (domain.VisiPemda, error) {
	script := "SELECT id, visi, tahun_awal_periode, tahun_akhir_periode, jenis_periode, keterangan FROM tb_visi_pemda WHERE id = ?"
	rows, err := tx.QueryContext(ctx, script, visiPemdaId)
	if err != nil {
		return domain.VisiPemda{}, err
	}
	defer rows.Close()

	var visiPemda domain.VisiPemda
	if rows.Next() {
		err := rows.Scan(
			&visiPemda.Id,
			&visiPemda.Visi,
			&visiPemda.TahunAwalPeriode,
			&visiPemda.TahunAkhirPeriode,
			&visiPemda.JenisPeriode,
			&visiPemda.Keterangan,
		)
		if err != nil {
			return domain.VisiPemda{}, err
		}
		return visiPemda, nil
	}

	return domain.VisiPemda{}, errors.New("visi pemda not found")
}

func (repository *VisiPemdaRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]domain.VisiPemda, error) {
	script := "SELECT id, visi, tahun_awal_periode, tahun_akhir_periode, jenis_periode, keterangan FROM tb_visi_pemda WHERE tahun_awal_periode = ? AND tahun_akhir_periode = ? AND jenis_periode = ?"
	rows, err := tx.QueryContext(ctx, script, tahunAwal, tahunAkhir, jenisPeriode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var visiPemdaList []domain.VisiPemda
	for rows.Next() {
		var visiPemda domain.VisiPemda
		err := rows.Scan(
			&visiPemda.Id,
			&visiPemda.Visi,
			&visiPemda.TahunAwalPeriode,
			&visiPemda.TahunAkhirPeriode,
			&visiPemda.JenisPeriode,
			&visiPemda.Keterangan,
		)
		if err != nil {
			return nil, err
		}
		visiPemdaList = append(visiPemdaList, visiPemda)
	}

	return visiPemdaList, nil
}
