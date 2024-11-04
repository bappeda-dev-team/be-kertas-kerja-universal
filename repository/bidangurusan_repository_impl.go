package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
)

type BidangUrusanRepositoryImpl struct {
}

func NewBidangUrusanRepositoryImpl() *BidangUrusanRepositoryImpl {
	return &BidangUrusanRepositoryImpl{}
}

func (repository *BidangUrusanRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, bidangurusan domainmaster.BidangUrusan) domainmaster.BidangUrusan {
	script := "INSERT INTO tb_bidang_urusan (id, kode_bidang_urusan, nama_bidang_urusan) VALUES (?, ?, ?)"
	_, err := tx.ExecContext(ctx, script, bidangurusan.Id, bidangurusan.KodeBidangUrusan, bidangurusan.NamaBidangUrusan)
	if err != nil {
		return bidangurusan
	}
	return bidangurusan
}

func (repository *BidangUrusanRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, bidangurusan domainmaster.BidangUrusan) domainmaster.BidangUrusan {
	script := "UPDATE tb_bidang_urusan SET kode_bidang_urusan = ?, nama_bidang_urusan = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, bidangurusan.KodeBidangUrusan, bidangurusan.NamaBidangUrusan, bidangurusan.Id)
	if err != nil {
		return bidangurusan
	}
	return bidangurusan
}

func (repository *BidangUrusanRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	script := "DELETE FROM tb_bidang_urusan WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *BidangUrusanRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (domainmaster.BidangUrusan, error) {
	script := "SELECT id, kode_bidang_urusan, nama_bidang_urusan FROM tb_bidang_urusan WHERE id = ?"
	rows, err := tx.QueryContext(ctx, script, id)
	if err != nil {
		return domainmaster.BidangUrusan{}, err
	}
	defer rows.Close()

	bidangurusan := domainmaster.BidangUrusan{}
	if rows.Next() {
		rows.Scan(&bidangurusan.Id, &bidangurusan.KodeBidangUrusan, &bidangurusan.NamaBidangUrusan)
	}
	return bidangurusan, nil
}

func (repository *BidangUrusanRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domainmaster.BidangUrusan, error) {
	script := "SELECT id, kode_bidang_urusan, nama_bidang_urusan FROM tb_bidang_urusan"
	rows, err := tx.QueryContext(ctx, script)
	if err != nil {
		return []domainmaster.BidangUrusan{}, err
	}
	defer rows.Close()

	var bidangurusans []domainmaster.BidangUrusan
	for rows.Next() {
		bidangurusan := domainmaster.BidangUrusan{}
		rows.Scan(&bidangurusan.Id, &bidangurusan.KodeBidangUrusan, &bidangurusan.NamaBidangUrusan)
		bidangurusans = append(bidangurusans, bidangurusan)
	}
	return bidangurusans, nil
}
