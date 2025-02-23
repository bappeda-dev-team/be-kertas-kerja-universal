package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"errors"
	"math/rand"
	"time"
)

type MisiPemdaRepositoryImpl struct{}

func NewMisiPemdaRepositoryImpl() *MisiPemdaRepositoryImpl {
	return &MisiPemdaRepositoryImpl{}
}

func (repository *MisiPemdaRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, visiPemda domain.MisiPemda) (domain.MisiPemda, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	visiPemda.Id = r.Intn(9000) + 1000

	script := "INSERT INTO tb_misi_pemda (id, id_visi, misi, urutan, tahun_awal_periode, tahun_akhir_periode, jenis_periode, keterangan) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, script, visiPemda.Id, visiPemda.IdVisi, visiPemda.Misi, visiPemda.Urutan, visiPemda.TahunAwalPeriode, visiPemda.TahunAkhirPeriode, visiPemda.JenisPeriode, visiPemda.Keterangan)
	if err != nil {
		return domain.MisiPemda{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return domain.MisiPemda{}, err
	}
	visiPemda.Id = int(id)

	return visiPemda, nil
}

func (repository *MisiPemdaRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, visiPemda domain.MisiPemda) (domain.MisiPemda, error) {
	script := "UPDATE tb_misi_pemda SET id_visi = ?, misi = ?, urutan = ?, tahun_awal_periode = ?, tahun_akhir_periode = ?, jenis_periode = ?, keterangan = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, visiPemda.IdVisi, visiPemda.Misi, visiPemda.Urutan, visiPemda.TahunAwalPeriode, visiPemda.TahunAkhirPeriode, visiPemda.JenisPeriode, visiPemda.Keterangan, visiPemda.Id)
	if err != nil {
		return domain.MisiPemda{}, err
	}

	return visiPemda, nil
}

func (repository *MisiPemdaRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, visiPemdaId int) error {
	script := "DELETE FROM tb_misi_pemda WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, visiPemdaId)
	if err != nil {
		return err
	}

	return nil
}

func (repository *MisiPemdaRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, visiPemdaId int) (domain.MisiPemda, error) {
	script := "SELECT id, id_visi, misi, urutan, tahun_awal_periode, tahun_akhir_periode, jenis_periode, keterangan FROM tb_misi_pemda WHERE id = ?"
	rows, err := tx.QueryContext(ctx, script, visiPemdaId)
	if err != nil {
		return domain.MisiPemda{}, err
	}
	defer rows.Close()

	var visiPemda domain.MisiPemda
	if rows.Next() {
		err := rows.Scan(
			&visiPemda.Id,
			&visiPemda.IdVisi,
			&visiPemda.Misi,
			&visiPemda.Urutan,
			&visiPemda.TahunAwalPeriode,
			&visiPemda.TahunAkhirPeriode,
			&visiPemda.JenisPeriode,
			&visiPemda.Keterangan,
		)
		if err != nil {
			return domain.MisiPemda{}, err
		}

		return visiPemda, nil
	}

	return domain.MisiPemda{}, errors.New("visi pemda not found")
}

func (repository *MisiPemdaRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]domain.MisiPemda, error) {
	script := "SELECT id, id_visi, misi, urutan, tahun_awal_periode, tahun_akhir_periode, jenis_periode, keterangan FROM tb_misi_pemda WHERE tahun_awal_periode = ? AND tahun_akhir_periode = ? AND jenis_periode = ?"
	rows, err := tx.QueryContext(ctx, script, tahunAwal, tahunAkhir, jenisPeriode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var visiPemdaList []domain.MisiPemda
	for rows.Next() {
		var visiPemda domain.MisiPemda
		err := rows.Scan(
			&visiPemda.Id,
			&visiPemda.IdVisi,
			&visiPemda.Misi,
			&visiPemda.Urutan,
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

func (repository *MisiPemdaRepositoryImpl) FindByIdWithDefault(ctx context.Context, tx *sql.Tx, visiPemdaId int) (domain.MisiPemda, error) {
	if visiPemdaId == 0 {
		return domain.MisiPemda{
			Id:                0,
			IdVisi:            0,
			Misi:              "",
			Urutan:            0,
			TahunAwalPeriode:  "",
			TahunAkhirPeriode: "",
			JenisPeriode:      "",
			Keterangan:        "",
		}, nil
	}

	return repository.FindById(ctx, tx, visiPemdaId)
}

func (repository *MisiPemdaRepositoryImpl) CheckUrutanExists(ctx context.Context, tx *sql.Tx, idVisi int, urutan int) (bool, error) {
	script := "SELECT EXISTS(SELECT 1 FROM tb_misi_pemda WHERE id_visi = ? AND urutan = ?)"
	var exists bool
	err := tx.QueryRowContext(ctx, script, idVisi, urutan).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (repository *MisiPemdaRepositoryImpl) CheckUrutanExistsExcept(ctx context.Context, tx *sql.Tx, idVisi int, urutan int, id int) (bool, error) {
	script := "SELECT EXISTS(SELECT 1 FROM tb_misi_pemda WHERE id_visi = ? AND urutan = ? AND id != ?)"
	var exists bool
	err := tx.QueryRowContext(ctx, script, idVisi, urutan, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
