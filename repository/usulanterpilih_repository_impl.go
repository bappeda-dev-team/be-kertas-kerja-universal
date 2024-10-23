package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"errors"
	"fmt"
)

type UsulanTerpilihRepositoryImpl struct {
}

func NewUsulanTerpilihRepositoryImpl() *UsulanTerpilihRepositoryImpl {
	return &UsulanTerpilihRepositoryImpl{}
}

func (repository *UsulanTerpilihRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, usulan domain.UsulanTerpilih) (domain.UsulanTerpilih, error) {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO tb_usulan_terpilih (id, keterangan, jenis_usulan, usulan_id, rekin_id, tahun, kode_opd)
		VALUES (?, ?, ?, ?, ?, ?, ?)
`, usulan.Id, usulan.Keterangan, usulan.JenisUsulan, usulan.UsulanId, usulan.RekinId, usulan.Tahun, usulan.KodeOpd)
	if err != nil {
		return domain.UsulanTerpilih{}, err
	}

	var updateQuery string
	switch usulan.JenisUsulan {
	case "mandatori":
		updateQuery = "UPDATE tb_usulan_mandatori SET is_active = true, status = 'usulan telah diambil', rekin_id = ? WHERE id = ?"
	case "musrebang":
		updateQuery = "UPDATE tb_usulan_musrebang SET is_active = true, status = 'usulan telah diambil', rekin_id = ? WHERE id = ?"
	case "inisiatif":
		updateQuery = "UPDATE tb_usulan_inisiatif SET is_active = true, status = 'usulan telah diambil', rekin_id = ? WHERE id = ?"
	case "pokok_pikiran":
		updateQuery = "UPDATE tb_usulan_pokok_pikiran SET is_active = true, status = 'usulan telah diambil', rekin_id = ? WHERE id = ?"
	default:
		return domain.UsulanTerpilih{}, errors.New("jenis usulan tidak valid")
	}

	_, err = tx.ExecContext(ctx, updateQuery, usulan.RekinId, usulan.UsulanId)
	if err != nil {
		return domain.UsulanTerpilih{}, err
	}

	return usulan, nil
}

func (repository *UsulanTerpilihRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, idUsulan string) error {
	// Ambil informasi usulan terpilih sebelum dihapus
	var jenisUsulan, usulanId string
	err := tx.QueryRowContext(ctx, "SELECT jenis_usulan, usulan_id FROM tb_usulan_terpilih WHERE usulan_id = ?", idUsulan).Scan(&jenisUsulan, &usulanId)
	if err != nil {
		return fmt.Errorf("error saat mengambil informasi usulan terpilih: %v", err)
	}

	// Hapus dari tb_usulan_terpilih
	script := "DELETE FROM tb_usulan_terpilih WHERE usulan_id = ?"
	_, err = tx.ExecContext(ctx, script, idUsulan)
	if err != nil {
		return fmt.Errorf("error saat menghapus usulan terpilih: %v", err)
	}

	// Update tabel sesuai dengan jenis usulan
	var updateQuery string
	switch jenisUsulan {
	case "mandatori":
		updateQuery = "UPDATE tb_usulan_mandatori SET is_active = false, status = 'usulan dibatalkan', rekin_id = '' WHERE id = ?"
	case "musrebang":
		updateQuery = "UPDATE tb_usulan_musrebang SET is_active = false, status = 'usulan dibatalkan', rekin_id = '' WHERE id = ?"
	case "inisiatif":
		updateQuery = "UPDATE tb_usulan_inisiatif SET is_active = false, status = 'usulan dibatalkan', rekin_id = '' WHERE id = ?"
	case "pokok_pikiran":
		updateQuery = "UPDATE tb_usulan_pokok_pikiran SET is_active = false, status = 'usulan dibatalkan', rekin_id = '' WHERE id = ?"
	default:
		return fmt.Errorf("jenis usulan tidak valid: %s", jenisUsulan)
	}

	_, err = tx.ExecContext(ctx, updateQuery, usulanId)
	if err != nil {
		return fmt.Errorf("error saat memperbarui status usulan: %v", err)
	}

	return nil
}

func (repository *UsulanTerpilihRepositoryImpl) ExistsByJenisAndUsulanId(ctx context.Context, tx *sql.Tx, jenisUsulan string, usulanId string) (bool, error) {
	SQL := "SELECT COUNT(*) FROM tb_usulan_terpilih WHERE jenis_usulan = ? AND usulan_id = ?"
	var count int
	err := tx.QueryRowContext(ctx, SQL, jenisUsulan, usulanId).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repository *UsulanTerpilihRepositoryImpl) ValidateJenisAndUsulanId(ctx context.Context, tx *sql.Tx, jenisUsulan string, usulanId string) (bool, error) {
	var tabelUsulan string
	switch jenisUsulan {
	case "mandatori":
		tabelUsulan = "tb_usulan_mandatori"
	case "musrebang":
		tabelUsulan = "tb_usulan_musrebang"
	case "inisiatif":
		tabelUsulan = "tb_usulan_inisiatif"
	case "pokok_pikiran":
		tabelUsulan = "tb_usulan_pokok_pikiran"
	default:
		return false, errors.New("jenis usulan tidak valid")
	}

	SQL := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id = ?", tabelUsulan)
	var count int
	err := tx.QueryRowContext(ctx, SQL, usulanId).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
