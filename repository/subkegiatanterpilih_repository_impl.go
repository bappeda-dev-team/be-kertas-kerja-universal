package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"fmt"

	"github.com/google/uuid"
)

type SubKegiatanTerpilihRepositoryImpl struct {
}

func NewSubKegiatanTerpilihRepositoryImpl() *SubKegiatanTerpilihRepositoryImpl {
	return &SubKegiatanTerpilihRepositoryImpl{}
}

func (repository *SubKegiatanTerpilihRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, subKegiatanTerpilih domain.SubKegiatanTerpilih) (domain.SubKegiatanTerpilih, error) {
	script := "UPDATE tb_rencana_kinerja SET kode_subkegiatan = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, subKegiatanTerpilih.KodeSubKegiatan, subKegiatanTerpilih.Id)
	if err != nil {
		return subKegiatanTerpilih, err
	}

	return subKegiatanTerpilih, nil
}

func (repository *SubKegiatanTerpilihRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id string, kodeSubKegiatan string) error {
	scriptDelete := "UPDATE tb_rencana_kinerja SET kode_subkegiatan = '' WHERE id = ? AND kode_subkegiatan = ?"
	_, err := tx.ExecContext(ctx, scriptDelete, id, kodeSubKegiatan)
	if err != nil {
		return err
	}

	return nil
}

func (repository *SubKegiatanTerpilihRepositoryImpl) FindByIdAndKodeSubKegiatan(ctx context.Context, tx *sql.Tx, id string, kodeSubKegiatan string) (domain.SubKegiatanTerpilih, error) {
	script := "SELECT id, kode_subkegiatan FROM tb_rencana_kinerja WHERE id = ? AND kode_subkegiatan = ?"
	var subKegiatanTerpilih domain.SubKegiatanTerpilih
	err := tx.QueryRowContext(ctx, script, id, kodeSubKegiatan).Scan(&subKegiatanTerpilih.Id, &subKegiatanTerpilih.KodeSubKegiatan)
	return subKegiatanTerpilih, err
}

func (repository *SubKegiatanTerpilihRepositoryImpl) CreateRekin(ctx context.Context, tx *sql.Tx, idSubKegiatan string, rekinId string) error {
	// Validasi keberadaan subkegiatan di tb_subkegiatan
	checkSubkegiatanScript := "SELECT COUNT(*) FROM tb_subkegiatan WHERE id = ?"
	var subkegiatanCount int
	err := tx.QueryRowContext(ctx, checkSubkegiatanScript, idSubKegiatan).Scan(&subkegiatanCount)
	if err != nil {
		return fmt.Errorf("error saat memeriksa data subkegiatan: %v", err)
	}
	if subkegiatanCount == 0 {
		return fmt.Errorf("subkegiatan dengan id %s tidak ditemukan di tb_subkegiatan", idSubKegiatan)
	}

	// Hapus data subkegiatan terpilih yang lama untuk rekin_id yang sama
	deleteScript := "DELETE FROM tb_subkegiatan_terpilih WHERE rekin_id = ?"
	_, err = tx.ExecContext(ctx, deleteScript, rekinId)
	if err != nil {
		return fmt.Errorf("error saat menghapus data subkegiatan terpilih yang lama: %v", err)
	}

	// Generate UUID baru untuk primary key
	newId := uuid.New().String()

	// Insert data baru ke tb_subkegiatan_terpilih
	script := "INSERT INTO tb_subkegiatan_terpilih (id, subkegiatan_id, rekin_id) VALUES (?, ?, ?)"
	result, err := tx.ExecContext(ctx, script, newId, idSubKegiatan, rekinId)
	if err != nil {
		return fmt.Errorf("error saat menyimpan subkegiatan terpilih: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error saat memeriksa rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("gagal menyimpan subkegiatan terpilih")
	}

	return nil
}

func (repository *SubKegiatanTerpilihRepositoryImpl) DeleteSubKegiatanTerpilih(ctx context.Context, tx *sql.Tx, idSubKegiatan string) error {
	script := "DELETE FROM tb_subkegiatan_terpilih WHERE id = ?"
	result, err := tx.ExecContext(ctx, script, idSubKegiatan)
	if err != nil {
		return fmt.Errorf("error saat menghapus subkegiatan terpilih: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error saat memeriksa rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subkegiatan dengan id %s tidak ditemukan", idSubKegiatan)
	}

	return nil
}

func (repository *SubKegiatanTerpilihRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, rekinId string) ([]domain.SubKegiatanTerpilih, error) {
	script := "SELECT id, subkegiatan_id, rekin_id FROM tb_subkegiatan_terpilih WHERE rekin_id = ?"
	rows, err := tx.QueryContext(ctx, script, rekinId)
	if err != nil {
		return nil, fmt.Errorf("error saat mengambil data subkegiatan terpilih: %v", err)
	}
	defer rows.Close()

	var result []domain.SubKegiatanTerpilih
	for rows.Next() {
		var subKegiatanTerpilih domain.SubKegiatanTerpilih
		err := rows.Scan(&subKegiatanTerpilih.Id, &subKegiatanTerpilih.SubkegiatanId, &subKegiatanTerpilih.RekinId)
		if err != nil {
			return nil, fmt.Errorf("error saat scanning data subkegiatan terpilih: %v", err)
		}
		result = append(result, subKegiatanTerpilih)
	}

	return result, nil
}
