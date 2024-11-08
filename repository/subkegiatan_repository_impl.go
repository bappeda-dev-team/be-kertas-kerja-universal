package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
	"log"
)

type SubKegiatanRepositoryImpl struct {
}

func NewSubKegiatanRepositoryImpl() *SubKegiatanRepositoryImpl {
	return &SubKegiatanRepositoryImpl{}
}

func (repository *SubKegiatanRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, subKegiatan domain.SubKegiatan) (domain.SubKegiatan, error) {
	scriptSubKegiatan := `INSERT INTO tb_subkegiatan (id, kode_subkegiatan, nama_subkegiatan, kode_opd, tahun, pegawai_id) 
                         VALUES (?, ?, ?, ?, ?, ?)`

	_, err := tx.ExecContext(ctx, scriptSubKegiatan,
		subKegiatan.Id,
		subKegiatan.KodeSubKegiatan,
		subKegiatan.NamaSubKegiatan,
		subKegiatan.KodeOpd,
		subKegiatan.Tahun,
		subKegiatan.PegawaiId)
	if err != nil {
		return domain.SubKegiatan{}, err
	}

	for _, indikator := range subKegiatan.Indikator {
		scriptIndikator := `INSERT INTO tb_indikator (id, subkegiatan_id, indikator, tahun) 
						   VALUES (?, ?, ?, ?)`

		_, err = tx.ExecContext(ctx, scriptIndikator,
			indikator.Id,
			subKegiatan.Id,
			indikator.Indikator,
			indikator.Tahun)
		if err != nil {
			return domain.SubKegiatan{}, err
		}

		for _, target := range indikator.Target {
			scriptTarget := `INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) 
						   VALUES (?, ?, ?, ?, ?)`

			_, err = tx.ExecContext(ctx, scriptTarget,
				target.Id,
				indikator.Id,
				target.Target,
				target.Satuan,
				target.Tahun)
			if err != nil {
				return domain.SubKegiatan{}, err
			}
		}
	}

	return subKegiatan, nil
}

func (repository *SubKegiatanRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, subKegiatan domain.SubKegiatan) (domain.SubKegiatan, error) {
	// Update SubKegiatan
	scriptSubKegiatan := `UPDATE tb_subkegiatan 
                         SET nama_subkegiatan = ?, 
                             kode_opd = ?, 
                             tahun = ? 
                         WHERE id = ?`

	_, err := tx.ExecContext(ctx, scriptSubKegiatan,
		subKegiatan.NamaSubKegiatan,
		subKegiatan.KodeOpd,
		subKegiatan.Tahun,
		subKegiatan.Id)
	if err != nil {
		log.Printf("Error updating subkegiatan: %v", err)
		return domain.SubKegiatan{}, err
	}

	// Hapus indikator dan target yang lama
	scriptDeleteTarget := `DELETE FROM tb_target 
                          WHERE indikator_id IN (
                              SELECT id FROM tb_indikator 
                              WHERE subkegiatan_id = ?
                          )`
	_, err = tx.ExecContext(ctx, scriptDeleteTarget, subKegiatan.Id)
	if err != nil {
		log.Printf("Error deleting old targets: %v", err)
		return domain.SubKegiatan{}, err
	}

	scriptDeleteIndikator := `DELETE FROM tb_indikator WHERE subkegiatan_id = ?`
	_, err = tx.ExecContext(ctx, scriptDeleteIndikator, subKegiatan.Id)
	if err != nil {
		log.Printf("Error deleting old indicators: %v", err)
		return domain.SubKegiatan{}, err
	}

	// Insert indikator baru
	for _, indikator := range subKegiatan.Indikator {
		scriptIndikator := `INSERT INTO tb_indikator (id, subkegiatan_id, indikator, tahun) 
						   VALUES (?, ?, ?, ?)`

		_, err = tx.ExecContext(ctx, scriptIndikator,
			indikator.Id,
			subKegiatan.Id,
			indikator.Indikator,
			indikator.Tahun)
		if err != nil {
			log.Printf("Error inserting new indicator: %v", err)
			return domain.SubKegiatan{}, err
		}

		// Insert target baru untuk setiap indikator
		for _, target := range indikator.Target {
			scriptTarget := `INSERT INTO tb_target (id, indikator_id, target, satuan, tahun) 
						   VALUES (?, ?, ?, ?, ?)`

			_, err = tx.ExecContext(ctx, scriptTarget,
				target.Id,
				indikator.Id,
				target.Target,
				target.Satuan,
				target.Tahun)
			if err != nil {
				log.Printf("Error inserting new target: %v", err)
				return domain.SubKegiatan{}, err
			}
		}
	}

	return subKegiatan, nil
}

func (repository *SubKegiatanRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string, pegawaiId string) ([]domain.SubKegiatan, error) {
	script := `SELECT id, kode_subkegiatan, nama_subkegiatan, kode_opd, tahun, created_at FROM tb_subkegiatan WHERE 1=1`
	var params []interface{}

	if kodeOpd != "" {
		script += ` AND kode_opd = ?`
		params = append(params, kodeOpd)
	}
	if pegawaiId != "" {
		script += ` AND pegawai_id = ?`
		params = append(params, pegawaiId)
	}
	script += ` ORDER BY created_at ASC`

	rows, err := tx.QueryContext(ctx, script, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subKegiatans []domain.SubKegiatan
	for rows.Next() {
		subKegiatan := domain.SubKegiatan{}
		err := rows.Scan(&subKegiatan.Id, &subKegiatan.KodeSubKegiatan, &subKegiatan.NamaSubKegiatan, &subKegiatan.KodeOpd, &subKegiatan.Tahun, &subKegiatan.CreatedAt)
		if err != nil {
			return nil, err
		}
		subKegiatans = append(subKegiatans, subKegiatan)
	}

	return subKegiatans, nil
}

func (repository *SubKegiatanRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, subKegiatanId string) (domain.SubKegiatan, error) {
	script := `SELECT id, kode_subkegiatan, pegawai_id, nama_subkegiatan, kode_opd, tahun FROM tb_subkegiatan WHERE id = ?`

	rows, err := tx.QueryContext(ctx, script, subKegiatanId)
	if err != nil {
		return domain.SubKegiatan{}, err
	}

	subKegiatan := domain.SubKegiatan{}
	if rows.Next() {
		rows.Scan(&subKegiatan.Id, &subKegiatan.KodeSubKegiatan, &subKegiatan.PegawaiId, &subKegiatan.NamaSubKegiatan, &subKegiatan.KodeOpd, &subKegiatan.Tahun)
	}

	return subKegiatan, nil
}

func (repository *SubKegiatanRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, subKegiatanId string) error {
	script := `DELETE FROM tb_subkegiatan WHERE id = ?`
	_, err := tx.ExecContext(ctx, script, subKegiatanId)
	return err
}

func (repository *SubKegiatanRepositoryImpl) FindIndikatorBySubKegiatanId(ctx context.Context, tx *sql.Tx, subKegiatanId string) ([]domain.Indikator, error) {
	script := "SELECT id, subkegiatan_id, indikator, tahun FROM tb_indikator WHERE subkegiatan_id = ?"
	params := []interface{}{subKegiatanId}

	rows, err := tx.QueryContext(ctx, script, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indikators []domain.Indikator
	for rows.Next() {
		var indikator domain.Indikator
		err := rows.Scan(&indikator.Id, &indikator.ProgramId, &indikator.Indikator, &indikator.Tahun)
		if err != nil {
			return nil, err
		}
		indikators = append(indikators, indikator)
	}

	return indikators, nil
}

func (repository *SubKegiatanRepositoryImpl) FindTargetByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string) ([]domain.Target, error) {
	script := "SELECT id, indikator_id, target, satuan, tahun FROM tb_target WHERE indikator_id = ?"
	params := []interface{}{indikatorId}

	rows, err := tx.QueryContext(ctx, script, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var targets []domain.Target
	for rows.Next() {
		var target domain.Target
		err := rows.Scan(&target.Id, &target.IndikatorId, &target.Target, &target.Satuan, &target.Tahun)
		if err != nil {
			return nil, err
		}
		targets = append(targets, target)
	}

	return targets, nil
}
