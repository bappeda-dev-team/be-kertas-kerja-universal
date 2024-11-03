package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
)

type OpdRepositoryImpl struct {
}

func NewOpdRepositoryImpl() *OpdRepositoryImpl {
	return &OpdRepositoryImpl{}
}

func (repository *OpdRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, opd domainmaster.Opd) (domainmaster.Opd, error) {
	script := `INSERT INTO tb_operasional_daerah (
		id, kode_opd, nama_opd, singkatan, alamat, telepon, fax, 
		email, website, nama_kepala_opd, nip_kepala_opd, pangkat_kepala, id_lembaga
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := tx.ExecContext(ctx, script,
		opd.Id, opd.KodeOpd, opd.NamaOpd, opd.Singkatan, opd.Alamat,
		opd.Telepon, opd.Fax, opd.Email, opd.Website, opd.NamaKepalaOpd,
		opd.NIPKepalaOpd, opd.PangkatKepala, opd.IdLembaga)
	if err != nil {
		return opd, err
	}
	return opd, nil
}

func (repository *OpdRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, opd domainmaster.Opd) (domainmaster.Opd, error) {
	script := `UPDATE tb_operasional_daerah SET 
		kode_opd = ?, nama_opd = ?, singkatan = ?, alamat = ?, 
		telepon = ?, fax = ?, email = ?, website = ?, 
		nama_kepala_opd = ?, nip_kepala_opd = ?, pangkat_kepala = ?, 
		id_lembaga = ? 
		WHERE id = ?`

	_, err := tx.ExecContext(ctx, script,
		opd.KodeOpd, opd.NamaOpd, opd.Singkatan, opd.Alamat,
		opd.Telepon, opd.Fax, opd.Email, opd.Website,
		opd.NamaKepalaOpd, opd.NIPKepalaOpd, opd.PangkatKepala,
		opd.IdLembaga, opd.Id)
	if err != nil {
		return opd, err
	}
	return opd, nil
}

func (repository *OpdRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, opdId string) error {
	script := "DELETE FROM tb_operasional_daerah WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, opdId)
	if err != nil {
		return err
	}
	return nil
}

func (repository *OpdRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domainmaster.Opd, error) {
	script := `SELECT 
		id, kode_opd, nama_opd, singkatan, alamat, telepon, fax,
		email, website, nama_kepala_opd, nip_kepala_opd, pangkat_kepala,
		id_lembaga 
		FROM tb_operasional_daerah`
	rows, err := tx.QueryContext(ctx, script)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var opds []domainmaster.Opd
	for rows.Next() {
		opd := domainmaster.Opd{}
		err := rows.Scan(
			&opd.Id, &opd.KodeOpd, &opd.NamaOpd, &opd.Singkatan,
			&opd.Alamat, &opd.Telepon, &opd.Fax, &opd.Email,
			&opd.Website, &opd.NamaKepalaOpd, &opd.NIPKepalaOpd,
			&opd.PangkatKepala, &opd.IdLembaga)
		if err != nil {
			return nil, err
		}
		opds = append(opds, opd)
	}
	return opds, nil
}

func (repository *OpdRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, opdId string) (domainmaster.Opd, error) {
	script := "SELECT id, kode_opd, nama_opd, singkatan, alamat, telepon, fax, email, website, nama_kepala_opd, nip_kepala_opd, pangkat_kepala, id_lembaga FROM tb_operasional_daerah WHERE id = ?"
	rows, err := tx.QueryContext(ctx, script, opdId)
	if err != nil {
		return domainmaster.Opd{}, err
	}
	defer rows.Close()

	var opd domainmaster.Opd
	if rows.Next() {
		err := rows.Scan(&opd.Id, &opd.KodeOpd, &opd.NamaOpd, &opd.Singkatan, &opd.Alamat, &opd.Telepon, &opd.Fax, &opd.Email, &opd.Website, &opd.NamaKepalaOpd, &opd.NIPKepalaOpd, &opd.PangkatKepala, &opd.IdLembaga)
		helper.PanicIfError(err)
	}
	return opd, nil
}

func (repository *OpdRepositoryImpl) FindByKodeOpd(ctx context.Context, tx *sql.Tx, kodeOpd string) (domainmaster.Opd, error) {
	script := "SELECT id, kode_opd, nama_opd, singkatan, alamat, telepon, fax, email, website, nama_kepala_opd, nip_kepala_opd, pangkat_kepala, id_lembaga FROM tb_operasional_daerah WHERE kode_opd = ?"
	rows, err := tx.QueryContext(ctx, script, kodeOpd)
	if err != nil {
		return domainmaster.Opd{}, err
	}
	defer rows.Close()

	var opd domainmaster.Opd
	if rows.Next() {
		err := rows.Scan(&opd.Id, &opd.KodeOpd, &opd.NamaOpd, &opd.Singkatan, &opd.Alamat, &opd.Telepon, &opd.Fax, &opd.Email, &opd.Website, &opd.NamaKepalaOpd, &opd.NIPKepalaOpd, &opd.PangkatKepala, &opd.IdLembaga)
		helper.PanicIfError(err)
	}
	return opd, nil
}
