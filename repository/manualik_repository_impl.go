package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type ManualIKRepositoryImpl struct {
}

func NewManualIKRepositoryImpl() *ManualIKRepositoryImpl {
	return &ManualIKRepositoryImpl{}
}

func (repository *ManualIKRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, manualik domain.ManualIK) (domain.ManualIK, error) {
	script := `INSERT INTO tb_manual_ik (
        id, indikator_id, perspektif, tujuan_rekin, definisi, key_activities, 
        formula, jenis_indikator, kinerja, penduduk, spasial,
        unit_penanggung_jawab, unit_penyedia_jasa, sumber_data,
        jangka_waktu_awal, jangka_waktu_akhir, periode_pelaporan
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?)`

	_, err := tx.ExecContext(ctx, script,
		manualik.Id,
		manualik.IndikatorId,
		manualik.Perspektif,
		manualik.TujuanRekin,
		manualik.Definisi,
		manualik.KeyActivities,
		manualik.Formula,
		manualik.JenisIndikator,
		manualik.Kinerja,
		manualik.Penduduk,
		manualik.Spatial,
		manualik.UnitPenanggungJawab,
		manualik.UnitPenyediaJasa,
		manualik.SumberData,
		manualik.JangkaWaktuAwal,
		manualik.JangkaWaktuAkhir,
		manualik.PeriodePelaporan,
	)
	if err != nil {
		return manualik, err
	}

	return manualik, nil
}

func (repository *ManualIKRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, manualik domain.ManualIK) (domain.ManualIK, error) {
	script := `UPDATE tb_manual_ik SET 
        perspektif = ?, 
        tujuan_rekin = ?,
        definisi = ?,
        key_activities = ?,
        formula = ?,
        jenis_indikator = ?,
        kinerja = ?,
        penduduk = ?,
        spasial = ?,
        unit_penanggung_jawab = ?,
        unit_penyedia_jasa = ?,
        sumber_data = ?,
        jangka_waktu_awal = ?,
        jangka_waktu_akhir = ?,
        periode_pelaporan = ?
    WHERE indikator_id = ?`

	_, err := tx.ExecContext(ctx, script,
		manualik.Perspektif,
		manualik.TujuanRekin,
		manualik.Definisi,
		manualik.KeyActivities,
		manualik.Formula,
		manualik.JenisIndikator,
		manualik.Kinerja,
		manualik.Penduduk,
		manualik.Spatial,
		manualik.UnitPenanggungJawab,
		manualik.UnitPenyediaJasa,
		manualik.SumberData,
		manualik.JangkaWaktuAwal,
		manualik.JangkaWaktuAkhir,
		manualik.PeriodePelaporan,
		manualik.IndikatorId,
	)
	if err != nil {
		return manualik, err
	}

	// Ambil data yang baru diupdate menggunakan SELECT
	script = `SELECT id, indikator_id, perspektif, tujuan_rekin, definisi, key_activities, 
        formula, jenis_indikator, kinerja, penduduk, spasial,
        unit_penanggung_jawab, unit_penyedia_jasa, sumber_data,
        jangka_waktu_awal, jangka_waktu_akhir, periode_pelaporan 
        FROM tb_manual_ik WHERE indikator_id = ?`

	var result domain.ManualIK
	err = tx.QueryRowContext(ctx, script, manualik.IndikatorId).Scan(
		&result.Id,
		&result.IndikatorId,
		&result.Perspektif,
		&result.TujuanRekin,
		&result.Definisi,
		&result.KeyActivities,
		&result.Formula,
		&result.JenisIndikator,
		&result.Kinerja,
		&result.Penduduk,
		&result.Spatial,
		&result.UnitPenanggungJawab,
		&result.UnitPenyediaJasa,
		&result.SumberData,
		&result.JangkaWaktuAwal,
		&result.JangkaWaktuAkhir,
		&result.PeriodePelaporan,
	)
	if err != nil {
		return manualik, err
	}

	return result, nil
}

func (repository *ManualIKRepositoryImpl) FindManualIKByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string) ([]domain.ManualIK, error) {
	script := `SELECT id, indikator_id, perspektif, tujuan_rekin, definisi, key_activities, 
        formula, jenis_indikator, kinerja, penduduk, spasial,
        unit_penanggung_jawab, unit_penyedia_jasa, sumber_data,
        jangka_waktu_awal, jangka_waktu_akhir, periode_pelaporan 
        FROM tb_manual_ik WHERE indikator_id = ?`
	rows, err := tx.QueryContext(ctx, script, indikatorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var manualiks []domain.ManualIK
	for rows.Next() {
		manualik := domain.ManualIK{}
		err := rows.Scan(&manualik.Id, &manualik.IndikatorId, &manualik.Perspektif, &manualik.TujuanRekin, &manualik.Definisi, &manualik.KeyActivities, &manualik.Formula, &manualik.JenisIndikator, &manualik.Kinerja, &manualik.Penduduk, &manualik.Spatial, &manualik.UnitPenanggungJawab, &manualik.UnitPenyediaJasa, &manualik.SumberData, &manualik.JangkaWaktuAwal, &manualik.JangkaWaktuAkhir, &manualik.PeriodePelaporan)
		if err != nil {
			return nil, err
		}
		manualiks = append(manualiks, manualik)
	}

	return manualiks, nil
}
