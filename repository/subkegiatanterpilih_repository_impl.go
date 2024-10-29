package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type SubKegiatanTerpilihRepositoryImpl struct {
}

func NewSubKegiatanTerpilihRepositoryImpl() *SubKegiatanTerpilihRepositoryImpl {
	return &SubKegiatanTerpilihRepositoryImpl{}
}

func (repository *SubKegiatanTerpilihRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, subKegiatanTerpilih domain.SubKegiatanTerpilih) (domain.SubKegiatanTerpilih, error) {
	// Cek dan hapus data lama jika ada
	scriptCheck := "SELECT subkegiatan_id FROM tb_subkegiatan_terpilih WHERE rekin_id = ?"
	var oldSubKegiatanId string
	err := tx.QueryRowContext(ctx, scriptCheck, subKegiatanTerpilih.RencanaKinerjaId).Scan(&oldSubKegiatanId)
	if err != nil && err != sql.ErrNoRows {
		return domain.SubKegiatanTerpilih{}, err
	}

	if oldSubKegiatanId != "" {
		// Update tb_subkegiatan lama, set rekin_id menjadi kosong
		scriptUpdateOld := "UPDATE tb_subkegiatan SET rekin_id = '' WHERE id = ?"
		_, err = tx.ExecContext(ctx, scriptUpdateOld, oldSubKegiatanId)
		if err != nil {
			return domain.SubKegiatanTerpilih{}, err
		}

		// Hapus data lama dari tb_subkegiatan_terpilih
		scriptDelete := "DELETE FROM tb_subkegiatan_terpilih WHERE rekin_id = ?"
		_, err = tx.ExecContext(ctx, scriptDelete, subKegiatanTerpilih.RencanaKinerjaId)
		if err != nil {
			return domain.SubKegiatanTerpilih{}, err
		}
	}

	// Insert data baru ke tb_subkegiatan_terpilih
	script1 := "INSERT INTO tb_subkegiatan_terpilih (id, rekin_id, subkegiatan_id) VALUES (?, ?, ?)"
	_, err = tx.ExecContext(ctx, script1, subKegiatanTerpilih.Id, subKegiatanTerpilih.RencanaKinerjaId, subKegiatanTerpilih.SubKegiatanId)
	if err != nil {
		return domain.SubKegiatanTerpilih{}, err
	}

	// Update tb_subkegiatan baru dengan rekin_id
	script2 := "UPDATE tb_subkegiatan SET rekin_id = ? WHERE id = ?"
	_, err = tx.ExecContext(ctx, script2, subKegiatanTerpilih.RencanaKinerjaId, subKegiatanTerpilih.SubKegiatanId)
	if err != nil {
		return domain.SubKegiatanTerpilih{}, err
	}

	return subKegiatanTerpilih, nil
}

// func (repository *SubKegiatanTerpilihRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, subKegiatanTerpilih domain.SubKegiatanTerpilih) (domain.SubKegiatanTerpilih, error) {
// 	// Insert ke tb_subkegiatan_terpilih
// 	script1 := "INSERT INTO tb_subkegiatan_terpilih (id, rekin_id, subkegiatan_id) VALUES (?, ?, ?)"
// 	_, err := tx.ExecContext(ctx, script1, subKegiatanTerpilih.Id, subKegiatanTerpilih.RencanaKinerjaId, subKegiatanTerpilih.SubKegiatanId)
// 	if err != nil {
// 		return domain.SubKegiatanTerpilih{}, err
// 	}

// 	// Update tb_subkegiatan dengan rekin_id
// 	script2 := "UPDATE tb_subkegiatan SET rekin_id = ? WHERE id = ?"
// 	_, err = tx.ExecContext(ctx, script2, subKegiatanTerpilih.RencanaKinerjaId, subKegiatanTerpilih.SubKegiatanId)
// 	if err != nil {
// 		return domain.SubKegiatanTerpilih{}, err
// 	}

// 	return subKegiatanTerpilih, nil
// }

func (repository *SubKegiatanTerpilihRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, subKegiatanId string) error {
	// Hapus dari tb_subkegiatan_terpilih
	scriptDelete := "DELETE FROM tb_subkegiatan_terpilih WHERE subkegiatan_id = ?"
	_, err := tx.ExecContext(ctx, scriptDelete, subKegiatanId)
	if err != nil {
		return err
	}

	// Update tb_subkegiatan, set rekin_id menjadi NULL
	scriptUpdate := "UPDATE tb_subkegiatan SET rekin_id = '' WHERE id = ?"
	_, err = tx.ExecContext(ctx, scriptUpdate, subKegiatanId)
	return err
}

func (repository *SubKegiatanTerpilihRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, subKegiatanTerpilihId string) (domain.SubKegiatanTerpilih, error) {
	script := "SELECT * FROM tb_subkegiatan_terpilih WHERE id = ?"
	var subKegiatanTerpilih domain.SubKegiatanTerpilih
	err := tx.QueryRowContext(ctx, script, subKegiatanTerpilihId).Scan(&subKegiatanTerpilih.Id, &subKegiatanTerpilih.RencanaKinerjaId, &subKegiatanTerpilih.SubKegiatanId)
	return subKegiatanTerpilih, err
}

func (repository *SubKegiatanTerpilihRepositoryImpl) ExistsByRencanaKinerjaIdAndSubKegiatanId(ctx context.Context, tx *sql.Tx, rencanaKinerjaId string, subKegiatanId string) (bool, error) {
	SQL := "SELECT COUNT(*) FROM tb_subkegiatan_terpilih WHERE rekin_id = ? AND subkegiatan_id = ?"

	var count int
	err := tx.QueryRowContext(ctx, SQL, rencanaKinerjaId, subKegiatanId).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repository *SubKegiatanTerpilihRepositoryImpl) ExistsByRekinAndSubKegiatan(ctx context.Context, tx *sql.Tx, rekinId string, subKegiatanId string) (bool, error) {
	script := "SELECT COUNT(*) FROM tb_subkegiatan_terpilih WHERE rekin_id = ? AND subkegiatan_id = ?"
	var count int
	err := tx.QueryRowContext(ctx, script, rekinId, subKegiatanId).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repository *SubKegiatanTerpilihRepositoryImpl) ExistsInSubKegiatan(ctx context.Context, tx *sql.Tx, subKegiatanId string) (bool, error) {
	script := "SELECT COUNT(*) FROM tb_subkegiatan WHERE id = ?"
	var count int
	err := tx.QueryRowContext(ctx, script, subKegiatanId).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
