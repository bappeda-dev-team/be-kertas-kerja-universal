package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type DasarHukumRepositoryImpl struct {
}

func NewDasarHukumRepositoryImpl() *DasarHukumRepositoryImpl {
	return &DasarHukumRepositoryImpl{}
}

func (repository *DasarHukumRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, dasarHukum domain.DasarHukum) (domain.DasarHukum, error) {
	script := "INSERT INTO tb_dasar_hukum (id, rekin_id, pegawai_id, urutan, peraturan_terkait, uraian) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, script, dasarHukum.Id, dasarHukum.RekinId, dasarHukum.PegawaiId, dasarHukum.Urutan, dasarHukum.PeraturanTerkait, dasarHukum.Uraian)
	if err != nil {
		return domain.DasarHukum{}, err
	}

	return dasarHukum, nil
}

func (repository *DasarHukumRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, dasarHukum domain.DasarHukum) (domain.DasarHukum, error) {
	script := "UPDATE tb_dasar_hukum SET urutan = ?, peraturan_terkait = ?, uraian = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, dasarHukum.Urutan, dasarHukum.PeraturanTerkait, dasarHukum.Uraian, dasarHukum.Id)
	if err != nil {
		return domain.DasarHukum{}, err
	}

	return dasarHukum, nil
}
func (repository *DasarHukumRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, rekinId string, pegawaiId string) ([]domain.DasarHukum, error) {
	script := "SELECT id, rekin_id, pegawai_id, urutan, peraturan_terkait, uraian FROM tb_dasar_hukum WHERE rekin_id = ? AND pegawai_id = ? ORDER BY urutan ASC"
	rows, err := tx.QueryContext(ctx, script, rekinId, pegawaiId)
	if err != nil {
		return []domain.DasarHukum{}, err
	}
	defer rows.Close()

	var dasarHukum []domain.DasarHukum
	for rows.Next() {
		var dh domain.DasarHukum
		err = rows.Scan(&dh.Id, &dh.RekinId, &dh.PegawaiId, &dh.Urutan, &dh.PeraturanTerkait, &dh.Uraian)
		if err != nil {
			return []domain.DasarHukum{}, err
		}
		dasarHukum = append(dasarHukum, dh)
	}

	return dasarHukum, nil
}

func (repository *DasarHukumRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	script := "DELETE FROM tb_dasar_hukum WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *DasarHukumRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (domain.DasarHukum, error) {
	script := "SELECT id, rekin_id, pegawai_id, urutan, peraturan_terkait, uraian FROM tb_dasar_hukum WHERE id = ? ORDER BY urutan ASC"
	var dh domain.DasarHukum
	err := tx.QueryRowContext(ctx, script, id).Scan(&dh.Id, &dh.RekinId, &dh.PegawaiId, &dh.Urutan, &dh.PeraturanTerkait, &dh.Uraian)
	if err != nil {
		return domain.DasarHukum{}, err
	}
	return dh, nil
}

func (repository *DasarHukumRepositoryImpl) GetLastUrutan(ctx context.Context, tx *sql.Tx) (int, error) {
	script := "SELECT COALESCE(MAX(urutan), 0) FROM tb_dasar_hukum"
	var lastUrutan int
	err := tx.QueryRowContext(ctx, script).Scan(&lastUrutan)
	if err != nil {
		return 0, err
	}
	return lastUrutan, nil
}
