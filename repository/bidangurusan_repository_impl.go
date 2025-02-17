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

func (repository *BidangUrusanRepositoryImpl) FindByKodeOpd(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domainmaster.BidangUrusan, error) {
	// Memisahkan kode OPD untuk mendapatkan kode bidang urusan
	kodeBidangUrusans := make([]string, 0)

	// Format kode OPD: 1.01.2.22.0.00.01.0000
	// Kode bidang urusan terdiri dari 3 bagian: 1.01 | 2.22 | 0.00

	// Mengambil kode bidang urusan pertama (1.01)
	if len(kodeOpd) >= 4 {
		kode1 := kodeOpd[:4]
		if kode1 != "0.00" {
			kodeBidangUrusans = append(kodeBidangUrusans, kode1)
		}
	}

	// Mengambil kode bidang urusan kedua (2.22)
	if len(kodeOpd) >= 9 {
		kode2 := kodeOpd[5:9]
		if kode2 != "0.00" {
			kodeBidangUrusans = append(kodeBidangUrusans, kode2)
		}
	}

	// Mengambil kode bidang urusan ketiga (0.00)
	if len(kodeOpd) >= 14 {
		kode3 := kodeOpd[10:14]
		if kode3 != "0.00" {
			kodeBidangUrusans = append(kodeBidangUrusans, kode3)
		}
	}

	// Jika tidak ada kode bidang urusan yang valid
	if len(kodeBidangUrusans) == 0 {
		return []domainmaster.BidangUrusan{}, nil
	}

	// Membuat query dengan IN clause
	query := "SELECT id, kode_bidang_urusan, nama_bidang_urusan FROM tb_bidang_urusan WHERE kode_bidang_urusan IN ("
	params := make([]interface{}, len(kodeBidangUrusans))
	for i := range kodeBidangUrusans {
		if i > 0 {
			query += ","
		}
		query += "?"
		params[i] = kodeBidangUrusans[i]
	}
	query += ")"

	rows, err := tx.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bidangUrusans []domainmaster.BidangUrusan
	for rows.Next() {
		bidangUrusan := domainmaster.BidangUrusan{}
		err := rows.Scan(&bidangUrusan.Id, &bidangUrusan.KodeBidangUrusan, &bidangUrusan.NamaBidangUrusan)
		if err != nil {
			return nil, err
		}
		bidangUrusans = append(bidangUrusans, bidangUrusan)
	}

	return bidangUrusans, nil
}
