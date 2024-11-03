package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type PohonKinerjaRepositoryImpl struct {
}

func NewPohonKinerjaRepositoryImpl() *PohonKinerjaRepositoryImpl {
	return &PohonKinerjaRepositoryImpl{}
}

func (repository *PohonKinerjaRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, pohonKinerja domain.PohonKinerja) (domain.PohonKinerja, error) {
	script := "INSERT INTO tb_pohon_kinerja (nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, script, pohonKinerja.NamaPohon, pohonKinerja.Parent, pohonKinerja.JenisPohon, pohonKinerja.LevelPohon, pohonKinerja.KodeOpd, pohonKinerja.Keterangan, pohonKinerja.Tahun)
	if err != nil {
		return pohonKinerja, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return pohonKinerja, err
	}

	pohonKinerja.Id = int(lastInsertId)
	return pohonKinerja, nil
}

func (repository *PohonKinerjaRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, pohonKinerja domain.PohonKinerja) (domain.PohonKinerja, error) {
	script := "UPDATE tb_pohon_kinerja SET nama_pohon = ?, parent = ?, jenis_pohon = ?, level_pohon = ?, kode_opd = ?, keterangan = ?, tahun = ? WHERE id = ?"
	result, err := tx.ExecContext(ctx, script, pohonKinerja.NamaPohon, pohonKinerja.Parent, pohonKinerja.JenisPohon, pohonKinerja.LevelPohon, pohonKinerja.KodeOpd, pohonKinerja.Keterangan, pohonKinerja.Tahun, pohonKinerja.Id)
	if err != nil {
		return pohonKinerja, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return pohonKinerja, err
	}

	pohonKinerja.Id = int(lastInsertId)
	return pohonKinerja, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.PohonKinerja, error) {
	script := "SELECT id, parent, nama_pohon, jenis_pohon, level_pohon, kode_opd, keterangan, tahun FROM tb_pohon_kinerja WHERE id = ?"
	rows, err := tx.QueryContext(ctx, script, id)
	if err != nil {
		return domain.PohonKinerja{}, err
	}
	defer rows.Close()

	pohonKinerja := domain.PohonKinerja{}
	if rows.Next() {
		rows.Scan(&pohonKinerja.Id, &pohonKinerja.Parent, &pohonKinerja.NamaPohon, &pohonKinerja.JenisPohon, &pohonKinerja.LevelPohon, &pohonKinerja.KodeOpd, &pohonKinerja.Keterangan, &pohonKinerja.Tahun)
	}
	return pohonKinerja, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string) ([]domain.PohonKinerja, error) {
	script := "SELECT id, parent, nama_pohon, jenis_pohon, level_pohon, kode_opd, keterangan, tahun FROM tb_pohon_kinerja WHERE 1=1"
	params := []interface{}{}
	if kodeOpd != "" {
		script += " AND kode_opd = ?"
		params = append(params, kodeOpd)
	}
	if tahun != "" {
		script += " AND tahun = ?"
		params = append(params, tahun)
	}

	rows, err := tx.QueryContext(ctx, script, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pohonKinerjas []domain.PohonKinerja
	for rows.Next() {
		pohonKinerja := domain.PohonKinerja{}
		err := rows.Scan(&pohonKinerja.Id, &pohonKinerja.Parent, &pohonKinerja.NamaPohon, &pohonKinerja.JenisPohon, &pohonKinerja.LevelPohon, &pohonKinerja.KodeOpd, &pohonKinerja.Keterangan, &pohonKinerja.Tahun)
		if err != nil {
			return nil, err
		}
		pohonKinerjas = append(pohonKinerjas, pohonKinerja)
	}
	return pohonKinerjas, nil
}

func (repository *PohonKinerjaRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	script := "DELETE FROM tb_pohon_kinerja WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}
	return nil
}

// admin pokin
func (repository *PohonKinerjaRepositoryImpl) CreatePokinAdmin(ctx context.Context, tx *sql.Tx, pokinAdmin domain.PohonKinerja) (domain.PohonKinerja, error) {
	// Insert pohon kinerja tanpa ID
	scriptPokin := "INSERT INTO tb_pohon_kinerja (nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, scriptPokin,
		pokinAdmin.NamaPohon, pokinAdmin.Parent, pokinAdmin.JenisPohon, pokinAdmin.LevelPohon, pokinAdmin.KodeOpd, pokinAdmin.Keterangan, pokinAdmin.Tahun)
	if err != nil {
		return pokinAdmin, err
	}

	// Dapatkan ID yang baru dibuat
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return pokinAdmin, err
	}
	pokinAdmin.Id = int(lastInsertId)

	// Insert indikator
	for _, indikator := range pokinAdmin.Indikator {
		scriptIndikator := "INSERT INTO tb_indikator (id, pokin_id, indikator, tahun) VALUES (?, ?, ?, ?)"
		_, err := tx.ExecContext(ctx, scriptIndikator,
			indikator.Id, pokinAdmin.Id, indikator.Indikator,
			indikator.Tahun)
		if err != nil {
			return pokinAdmin, err
		}

		// Insert target untuk setiap indikator
		for _, target := range indikator.Target {
			scriptTarget := "INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)"
			_, err := tx.ExecContext(ctx, scriptTarget, target.Id, indikator.Id, target.Target, target.Satuan, target.Tahun)
			if err != nil {
				return pokinAdmin, err
			}
		}
	}

	return pokinAdmin, nil
}

func (repository *PohonKinerjaRepositoryImpl) UpdatePokinAdmin(ctx context.Context, tx *sql.Tx, pokinAdmin domain.PohonKinerja) (domain.PohonKinerja, error) {
	// Update data pohon kinerja
	scriptPokin := "UPDATE tb_pohon_kinerja SET nama_pohon = ?, parent = ?, jenis_pohon = ?, level_pohon = ?, kode_opd = ?, keterangan = ?, tahun = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, scriptPokin,
		pokinAdmin.NamaPohon,
		pokinAdmin.Parent,
		pokinAdmin.JenisPohon,
		pokinAdmin.LevelPohon,
		pokinAdmin.KodeOpd,
		pokinAdmin.Keterangan,
		pokinAdmin.Tahun,
		pokinAdmin.Id)
	if err != nil {
		return pokinAdmin, err
	}

	// Hapus indikator dan target yang lama
	scriptDeleteTarget := "DELETE FROM tb_target WHERE indikator_id IN (SELECT id FROM tb_indikator WHERE pokin_id = ?)"
	_, err = tx.ExecContext(ctx, scriptDeleteTarget, pokinAdmin.Id)
	if err != nil {
		return pokinAdmin, err
	}

	scriptDeleteIndikator := "DELETE FROM tb_indikator WHERE pokin_id = ?"
	_, err = tx.ExecContext(ctx, scriptDeleteIndikator, pokinAdmin.Id)
	if err != nil {
		return pokinAdmin, err
	}

	// Insert indikator baru
	for _, indikator := range pokinAdmin.Indikator {
		scriptIndikator := "INSERT INTO tb_indikator (id, pokin_id, indikator, tahun) VALUES (?, ?, ?, ?)"
		_, err := tx.ExecContext(ctx, scriptIndikator,
			indikator.Id,
			pokinAdmin.Id,
			indikator.Indikator,
			indikator.Tahun)
		if err != nil {
			return pokinAdmin, err
		}

		// Insert target baru untuk setiap indikator
		for _, target := range indikator.Target {
			scriptTarget := "INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) VALUES (?, ?, ?, ?, ?)"
			_, err := tx.ExecContext(ctx, scriptTarget,
				target.Id,
				indikator.Id,
				target.Target,
				target.Satuan,
				target.Tahun)
			if err != nil {
				return pokinAdmin, err
			}
		}
	}

	return pokinAdmin, nil
}

func (repository *PohonKinerjaRepositoryImpl) DeletePokinAdmin(ctx context.Context, tx *sql.Tx, id int) error {
	// Hapus target terlebih dahulu
	scriptDeleteTarget := "DELETE FROM tb_target WHERE indikator_id IN (SELECT id FROM tb_indikator WHERE pokin_id = ?)"
	_, err := tx.ExecContext(ctx, scriptDeleteTarget, id)
	if err != nil {
		return err
	}

	// Hapus indikator
	scriptDeleteIndikator := "DELETE FROM tb_indikator WHERE pokin_id = ?"
	_, err = tx.ExecContext(ctx, scriptDeleteIndikator, id)
	if err != nil {
		return err
	}

	// Hapus pohon kinerja
	scriptDeletePokin := "DELETE FROM tb_pohon_kinerja WHERE id = ?"
	_, err = tx.ExecContext(ctx, scriptDeletePokin, id)
	if err != nil {
		return err
	}

	return nil
}

func (repository *PohonKinerjaRepositoryImpl) FindPokinAdminById(ctx context.Context, tx *sql.Tx, id int) (domain.PohonKinerja, error) {
	script := "SELECT id, nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun FROM tb_pohon_kinerja WHERE id = ?"
	rows, err := tx.QueryContext(ctx, script, id)
	if err != nil {
		return domain.PohonKinerja{}, err
	}
	defer rows.Close()

	pokinAdmin := domain.PohonKinerja{}
	if rows.Next() {
		rows.Scan(&pokinAdmin.Id, &pokinAdmin.NamaPohon, &pokinAdmin.Parent, &pokinAdmin.JenisPohon, &pokinAdmin.LevelPohon, &pokinAdmin.KodeOpd, &pokinAdmin.Keterangan, &pokinAdmin.Tahun)
	}
	return pokinAdmin, nil
}

func (repository *PohonKinerjaRepositoryImpl) FindPokinAdminAll(ctx context.Context, tx *sql.Tx, tahun string) ([]domain.PohonKinerja, error) {
	// Query utama untuk pohon kinerja
	script := "SELECT id, nama_pohon, parent, jenis_pohon, level_pohon, kode_opd, keterangan, tahun FROM tb_pohon_kinerja WHERE tahun = ?"
	rows, err := tx.QueryContext(ctx, script, tahun)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pohonKinerjas []domain.PohonKinerja
	for rows.Next() {
		var pohonKinerja domain.PohonKinerja
		err := rows.Scan(&pohonKinerja.Id, &pohonKinerja.NamaPohon, &pohonKinerja.Parent, &pohonKinerja.JenisPohon,
			&pohonKinerja.LevelPohon, &pohonKinerja.KodeOpd, &pohonKinerja.Keterangan, &pohonKinerja.Tahun)
		if err != nil {
			return nil, err
		}

		// Ambil indikator untuk pohon kinerja ini
		scriptIndikator := "SELECT id, pokin_id, indikator, tahun FROM tb_indikator WHERE pokin_id = ?"
		indikatorRows, err := tx.QueryContext(ctx, scriptIndikator, pohonKinerja.Id)
		if err != nil {
			return nil, err
		}

		// Proses indikator
		for indikatorRows.Next() {
			var indikator domain.Indikator
			err := indikatorRows.Scan(&indikator.Id, &indikator.PokinId, &indikator.Indikator, &indikator.Tahun)
			if err != nil {
				indikatorRows.Close()
				return nil, err
			}

			// Ambil target untuk indikator ini
			scriptTarget := "SELECT id, indikator_id, target, satuan, tahun FROM tb_target WHERE indikator_id = ?"
			targetRows, err := tx.QueryContext(ctx, scriptTarget, indikator.Id)
			if err != nil {
				indikatorRows.Close()
				return nil, err
			}

			// Proses target
			for targetRows.Next() {
				var target domain.Target
				err := targetRows.Scan(&target.Id, &target.IndikatorId, &target.Target, &target.Satuan, &target.Tahun)
				if err != nil {
					targetRows.Close()
					indikatorRows.Close()
					return nil, err
				}
				indikator.Target = append(indikator.Target, target)
			}
			targetRows.Close()

			pohonKinerja.Indikator = append(pohonKinerja.Indikator, indikator)
		}
		indikatorRows.Close()

		pohonKinerjas = append(pohonKinerjas, pohonKinerja)
	}

	return pohonKinerjas, nil
}
