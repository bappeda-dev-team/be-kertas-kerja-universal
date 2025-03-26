package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
	"log"
)

type OpdRepositoryImpl struct {
}

func NewOpdRepositoryImpl() *OpdRepositoryImpl {
	return &OpdRepositoryImpl{}
}

func (repository *OpdRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, opd domainmaster.Opd) (domainmaster.Opd, error) {
	script := `INSERT INTO tb_operasional_daerah (
		id, kode_opd, nama_opd, singkatan, alamat, telepon, fax, 
		email, website, nama_kepala_opd, nip_kepala_opd, pangkat_kepala,
nama_admin_opd, no_wa_admin_opd,
id_lembaga
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := tx.ExecContext(ctx, script,
		opd.Id, opd.KodeOpd, opd.NamaOpd, opd.Singkatan, opd.Alamat,
		opd.Telepon, opd.Fax, opd.Email, opd.Website, opd.NamaKepalaOpd,
		opd.NIPKepalaOpd, opd.PangkatKepala,
		opd.NamaAdmin, opd.NoWaAdmin,
		opd.IdLembaga)
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
nama_admin_opd = ?, no_wa_admin_opd = ?,
		id_lembaga = ? 
		WHERE id = ?`

	_, err := tx.ExecContext(ctx, script,
		opd.KodeOpd, opd.NamaOpd, opd.Singkatan, opd.Alamat,
		opd.Telepon, opd.Fax, opd.Email, opd.Website,
		opd.NamaKepalaOpd, opd.NIPKepalaOpd, opd.PangkatKepala,
		opd.NamaAdmin, opd.NoWaAdmin,
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

func nullToEmpty(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// TODO: add kode lembaga filter
func (repository *OpdRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domainmaster.OpdWithBidangUrusan, error) {
	query := `SELECT
    opd.id,
	opd.kode_opd,
	opd.nama_opd,
	u1.kode_urusan AS kode_urusan_1,
	u1.nama_urusan AS nama_urusan_1,
	u2.kode_urusan AS kode_urusan_2,
	u2.nama_urusan AS nama_urusan_2,
	u3.kode_urusan AS kode_urusan_3,
	u3.nama_urusan AS nama_urusan_3,
	bu1.kode_bidang_urusan AS kode_bidang_urusan_1,
	bu1.nama_bidang_urusan AS nama_bidang_urusan_1,
	bu2.kode_bidang_urusan AS kode_bidang_urusan_2,
	bu2.nama_bidang_urusan AS nama_bidang_urusan_2,
	bu3.kode_bidang_urusan AS kode_bidang_urusan_3,
	bu3.nama_bidang_urusan AS nama_bidang_urusan_3,
    opd.nama_kepala_opd,
    opd.nip_kepala_opd,
    opd.pangkat_kepala,
    opd.nama_admin_opd,
    opd.no_wa_admin_opd
	FROM
		tb_operasional_daerah opd
    LEFT JOIN tb_urusan u1
    ON u1.kode_urusan = LEFT(opd.kode_opd, LOCATE('.', opd.kode_opd) - 1)
	LEFT JOIN tb_urusan u2
    ON u2.kode_urusan = REGEXP_SUBSTR(opd.kode_opd, '(?<=\\.)(\\d+)', 1, 2)
	LEFT JOIN tb_urusan u3
    ON u3.kode_urusan = REGEXP_SUBSTR(opd.kode_opd, '(?<=\\.)(\\d+)', 1, 4)
	LEFT JOIN tb_bidang_urusan bu1
	ON bu1.kode_bidang_urusan = REGEXP_SUBSTR(opd.kode_opd, '(\\d+)\\.\\d{2}', 1, 1)
	LEFT JOIN tb_bidang_urusan bu2
	ON bu2.kode_bidang_urusan = REGEXP_SUBSTR(opd.kode_opd, '(\\d+)\\.\\d{2}', 1, 2)
	LEFT JOIN tb_bidang_urusan bu3
	ON bu3.kode_bidang_urusan = REGEXP_SUBSTR(opd.kode_opd, '(\\d+)\\.\\d{2}', 1, 3);`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var opds []domainmaster.OpdWithBidangUrusan
	for rows.Next() {
		opd := domainmaster.OpdWithBidangUrusan{}

		// Use sql.NullString for nullable string fields
		var kodeUrusan1, namaUrusan1 sql.NullString
		var kodeUrusan2, namaUrusan2 sql.NullString
		var kodeUrusan3, namaUrusan3 sql.NullString
		var kodeBidangUrusan1, namaBidangUrusan1 sql.NullString
		var kodeBidangUrusan2, namaBidangUrusan2 sql.NullString
		var kodeBidangUrusan3, namaBidangUrusan3 sql.NullString
		var namaKepalaOpd, nipKepalaOpd sql.NullString
		var pangkatKepala, namaAdmin, noWaAdmin sql.NullString

		err := rows.Scan(
			&opd.Id,
			&opd.KodeOpd,
			&opd.NamaOpd,
			&kodeUrusan1,
			&namaUrusan1,
			&kodeUrusan2,
			&namaUrusan2,
			&kodeUrusan3,
			&namaUrusan3,
			&kodeBidangUrusan1,
			&namaBidangUrusan1,
			&kodeBidangUrusan2,
			&namaBidangUrusan2,
			&kodeBidangUrusan3,
			&namaBidangUrusan3,
			&namaKepalaOpd,
			&nipKepalaOpd,
			&pangkatKepala,
			&namaAdmin,
			&noWaAdmin,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		// Convert NULL values to empty strings
		opd.KodeUrusan1 = nullToEmpty(kodeUrusan1)
		opd.NamaUrusan1 = nullToEmpty(namaUrusan1)
		opd.KodeUrusan2 = nullToEmpty(kodeUrusan2)
		opd.NamaUrusan2 = nullToEmpty(namaUrusan2)
		opd.KodeUrusan3 = nullToEmpty(kodeUrusan3)
		opd.NamaUrusan3 = nullToEmpty(namaUrusan3)
		opd.KodeBidangUrusan1 = nullToEmpty(kodeBidangUrusan1)
		opd.NamaBidangUrusan1 = nullToEmpty(namaBidangUrusan1)
		opd.KodeBidangUrusan2 = nullToEmpty(kodeBidangUrusan2)
		opd.NamaBidangUrusan2 = nullToEmpty(namaBidangUrusan2)
		opd.KodeBidangUrusan3 = nullToEmpty(kodeBidangUrusan3)
		opd.NamaBidangUrusan3 = nullToEmpty(namaBidangUrusan3)
		opd.NamaKepalaOpd = nullToEmpty(namaKepalaOpd)
		opd.NIPKepalaOpd = nullToEmpty(nipKepalaOpd)
		opd.PangkatKepala = nullToEmpty(pangkatKepala)
		opd.NamaAdmin = nullToEmpty(namaAdmin)
		opd.NoWaAdmin = nullToEmpty(noWaAdmin)

		opds = append(opds, opd)
	}
	return opds, nil
}

func (repository *OpdRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, opdId string) (domainmaster.Opd, error) {
	script := `SELECT id, kode_opd, nama_opd, singkatan, alamat,
telepon, fax, email, website, nama_kepala_opd,
nip_kepala_opd, pangkat_kepala,
nama_admin_opd, no_wa_admin_opd,
id_lembaga FROM tb_operasional_daerah WHERE id = ?`
	rows, err := tx.QueryContext(ctx, script, opdId)
	if err != nil {
		return domainmaster.Opd{}, err
	}
	defer rows.Close()

	var opd domainmaster.Opd
	if rows.Next() {
		err := rows.Scan(&opd.Id,
			&opd.KodeOpd, &opd.NamaOpd, &opd.Singkatan, &opd.Alamat,
			&opd.Telepon, &opd.Fax, &opd.Email, &opd.Website,
			&opd.NamaKepalaOpd,
			&opd.NIPKepalaOpd,
			&opd.PangkatKepala,
			&opd.NamaAdmin,
			&opd.NoWaAdmin,
			&opd.IdLembaga)
		helper.PanicIfError(err)
	}
	return opd, nil
}

func (repository *OpdRepositoryImpl) InfoOpd(ctx context.Context, tx *sql.Tx, kodeOpd string, kodeLembaga string) (domainmaster.OpdWithBidangUrusan, error) {
	query := `
    SELECT
    opd.id,
	opd.kode_opd,
	opd.nama_opd,
	u1.kode_urusan AS kode_urusan_1,
	u1.nama_urusan AS nama_urusan_1,
	u2.kode_urusan AS kode_urusan_2,
	u2.nama_urusan AS nama_urusan_2,
	u3.kode_urusan AS kode_urusan_3,
	u3.nama_urusan AS nama_urusan_3,
	bu1.kode_bidang_urusan AS kode_bidang_urusan_1,
	bu1.nama_bidang_urusan AS nama_bidang_urusan_1,
	bu2.kode_bidang_urusan AS kode_bidang_urusan_2,
	bu2.nama_bidang_urusan AS nama_bidang_urusan_2,
	bu3.kode_bidang_urusan AS kode_bidang_urusan_3,
	bu3.nama_bidang_urusan AS nama_bidang_urusan_3,
    opd.nama_kepala_opd,
    opd.nip_kepala_opd,
    opd.pangkat_kepala,
    opd.nama_admin_opd,
    opd.no_wa_admin_opd
	FROM
		tb_operasional_daerah opd
    LEFT JOIN tb_urusan u1
    ON u1.kode_urusan = LEFT(opd.kode_opd, LOCATE('.', opd.kode_opd) - 1)
	LEFT JOIN tb_urusan u2
    ON u2.kode_urusan = REGEXP_SUBSTR(opd.kode_opd, '(?<=\\.)(\\d+)', 1, 2)
	LEFT JOIN tb_urusan u3
    ON u3.kode_urusan = REGEXP_SUBSTR(opd.kode_opd, '(?<=\\.)(\\d+)', 1, 4)
	LEFT JOIN tb_bidang_urusan bu1
	ON bu1.kode_bidang_urusan = REGEXP_SUBSTR(opd.kode_opd, '(\\d+)\\.\\d{2}', 1, 1)
	LEFT JOIN tb_bidang_urusan bu2
	ON bu2.kode_bidang_urusan = REGEXP_SUBSTR(opd.kode_opd, '(\\d+)\\.\\d{2}', 1, 2)
	LEFT JOIN tb_bidang_urusan bu3
	ON bu3.kode_bidang_urusan = REGEXP_SUBSTR(opd.kode_opd, '(\\d+)\\.\\d{2}', 1, 3)
    LEFT JOIN tb_lembaga lem
    ON lem.id = opd.id_lembaga
    WHERE opd.kode_opd = ? AND lem.kode_lembaga = ?;`

	rows, err := tx.QueryContext(ctx, query, kodeOpd, kodeLembaga)
	if err != nil {
		return domainmaster.OpdWithBidangUrusan{}, err
	}
	defer rows.Close()

	opd := domainmaster.OpdWithBidangUrusan{}
	for rows.Next() {

		// Use sql.NullString for nullable string fields
		var kodeUrusan1, namaUrusan1 sql.NullString
		var kodeUrusan2, namaUrusan2 sql.NullString
		var kodeUrusan3, namaUrusan3 sql.NullString
		var kodeBidangUrusan1, namaBidangUrusan1 sql.NullString
		var kodeBidangUrusan2, namaBidangUrusan2 sql.NullString
		var kodeBidangUrusan3, namaBidangUrusan3 sql.NullString
		var namaKepalaOpd, nipKepalaOpd sql.NullString
		var pangkatKepala, namaAdmin, noWaAdmin sql.NullString

		err := rows.Scan(
			&opd.Id,
			&opd.KodeOpd,
			&opd.NamaOpd,
			&kodeUrusan1,
			&namaUrusan1,
			&kodeUrusan2,
			&namaUrusan2,
			&kodeUrusan3,
			&namaUrusan3,
			&kodeBidangUrusan1,
			&namaBidangUrusan1,
			&kodeBidangUrusan2,
			&namaBidangUrusan2,
			&kodeBidangUrusan3,
			&namaBidangUrusan3,
			&namaKepalaOpd,
			&nipKepalaOpd,
			&pangkatKepala,
			&namaAdmin,
			&noWaAdmin,
		)
		if err != nil {
			log.Println(err)
			return domainmaster.OpdWithBidangUrusan{}, err
		}
		// Convert NULL values to empty strings
		opd.KodeUrusan1 = nullToEmpty(kodeUrusan1)
		opd.NamaUrusan1 = nullToEmpty(namaUrusan1)
		opd.KodeUrusan2 = nullToEmpty(kodeUrusan2)
		opd.NamaUrusan2 = nullToEmpty(namaUrusan2)
		opd.KodeUrusan3 = nullToEmpty(kodeUrusan3)
		opd.NamaUrusan3 = nullToEmpty(namaUrusan3)
		opd.KodeBidangUrusan1 = nullToEmpty(kodeBidangUrusan1)
		opd.NamaBidangUrusan1 = nullToEmpty(namaBidangUrusan1)
		opd.KodeBidangUrusan2 = nullToEmpty(kodeBidangUrusan2)
		opd.NamaBidangUrusan2 = nullToEmpty(namaBidangUrusan2)
		opd.KodeBidangUrusan3 = nullToEmpty(kodeBidangUrusan3)
		opd.NamaBidangUrusan3 = nullToEmpty(namaBidangUrusan3)
		opd.NamaKepalaOpd = nullToEmpty(namaKepalaOpd)
		opd.NIPKepalaOpd = nullToEmpty(nipKepalaOpd)
		opd.PangkatKepala = nullToEmpty(pangkatKepala)
		opd.NamaAdmin = nullToEmpty(namaAdmin)
		opd.NoWaAdmin = nullToEmpty(noWaAdmin)
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
		err := rows.Scan(&opd.Id,
			&opd.KodeOpd, &opd.NamaOpd, &opd.Singkatan, &opd.Alamat,
			&opd.Telepon, &opd.Fax, &opd.Email, &opd.Website,
			&opd.NamaKepalaOpd,
			&opd.NIPKepalaOpd,
			&opd.PangkatKepala,
			&opd.IdLembaga)
		helper.PanicIfError(err)
	}
	return opd, nil
}
