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
