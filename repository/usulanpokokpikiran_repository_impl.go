package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"
)

type UsulanPokokPikiranRepositoryImpl struct {
}

func NewUsulanPokokPikiranRepositoryImpl() *UsulanPokokPikiranRepositoryImpl {
	return &UsulanPokokPikiranRepositoryImpl{}
}

func (repository *UsulanPokokPikiranRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, usulan domain.UsulanPokokPikiran) (domain.UsulanPokokPikiran, error) {
	script := "INSERT INTO tb_usulan_pokok_pikiran (id, usulan, alamat, uraian, tahun, rekin_id, pegawai_id, kode_opd, status) VALUES (?,?,?,?,?,?,?,?,?)"
	_, err := tx.ExecContext(ctx, script, usulan.Id, usulan.Usulan, usulan.Alamat, usulan.Uraian, usulan.Tahun, usulan.RekinId, usulan.PegawaiId, usulan.KodeOpd, usulan.Status)
	if err != nil {
		return domain.UsulanPokokPikiran{}, fmt.Errorf("error saat menyimpan usulan pokok pikiran: %v", err)
	}
	return usulan, nil
}

func (repository *UsulanPokokPikiranRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, usulan domain.UsulanPokokPikiran) (domain.UsulanPokokPikiran, error) {
	script := "UPDATE tb_usulan_pokok_pikiran SET usulan = ?, alamat = ?, uraian = ?, tahun = ?, pegawai_id = ?, kode_opd = ?, status = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, usulan.Usulan, usulan.Alamat, usulan.Uraian, usulan.Tahun, usulan.PegawaiId, usulan.KodeOpd, usulan.Status, usulan.Id)
	if err != nil {
		return domain.UsulanPokokPikiran{}, fmt.Errorf("error saat mengupdate usulan pokok pikiran: %v", err)
	}
	return usulan, nil
}

func (repository *UsulanPokokPikiranRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, idUsulan string) (domain.UsulanPokokPikiran, error) {
	script := "SELECT id, usulan, alamat, uraian, tahun, rekin_id, pegawai_id, kode_opd, is_active, status, created_at FROM tb_usulan_pokok_pikiran WHERE id = ?"
	row := tx.QueryRowContext(ctx, script, idUsulan)

	var usulan domain.UsulanPokokPikiran
	err := row.Scan(&usulan.Id, &usulan.Usulan, &usulan.Alamat, &usulan.Uraian, &usulan.Tahun, &usulan.RekinId, &usulan.PegawaiId, &usulan.KodeOpd, &usulan.IsActive, &usulan.Status, &usulan.CreatedAt)
	if err != nil {
		return domain.UsulanPokokPikiran{}, fmt.Errorf("error saat mencari usulan pokok pikiran: %v", err)
	}

	return usulan, nil
}

func (repository *UsulanPokokPikiranRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, pegawaiId *string, isActive *bool, rekinId *string) ([]domain.UsulanPokokPikiran, error) {
	script := "SELECT id, usulan, alamat, uraian, tahun, rekin_id, pegawai_id, kode_opd, is_active, status, created_at FROM tb_usulan_pokok_pikiran WHERE 1=1"
	var params []interface{}

	if pegawaiId != nil {
		script += " AND pegawai_id = ?"
		params = append(params, *pegawaiId)
	}

	if isActive != nil {
		script += " AND is_active = ?"
		params = append(params, *isActive)
	}

	if rekinId != nil {
		script += " AND rekin_id = ?"
		params = append(params, *rekinId)
	}

	script += " ORDER BY created_at DESC"

	rows, err := tx.QueryContext(ctx, script, params...)
	if err != nil {
		return nil, fmt.Errorf("error saat mencari usulan pokok pikiran: %v", err)
	}

	defer rows.Close()

	var usulans []domain.UsulanPokokPikiran
	for rows.Next() {
		var usulan domain.UsulanPokokPikiran
		err := rows.Scan(&usulan.Id, &usulan.Usulan, &usulan.Alamat, &usulan.Uraian, &usulan.Tahun, &usulan.RekinId, &usulan.PegawaiId, &usulan.KodeOpd, &usulan.IsActive, &usulan.Status, &usulan.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error saat membaca usulan pokok pikiran: %v", err)
		}
		usulans = append(usulans, usulan)
	}

	return usulans, nil
}

func (repository *UsulanPokokPikiranRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, idUsulan string) error {
	script := "DELETE FROM tb_usulan_pokok_pikiran WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, idUsulan)
	if err != nil {
		return fmt.Errorf("error saat menghapus usulan pokok pikiran: %v", err)
	}
	return nil
}
