package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type SubKegiatanRepository interface {
	Create(ctx context.Context, tx *sql.Tx, subKegiatan domain.SubKegiatan) (domain.SubKegiatan, error)
	FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string, pegawaiId string) ([]domain.SubKegiatan, error)
	Update(ctx context.Context, tx *sql.Tx, subKegiatan domain.SubKegiatan) (domain.SubKegiatan, error)
	FindById(ctx context.Context, tx *sql.Tx, subKegiatanId string) (domain.SubKegiatan, error)
	Delete(ctx context.Context, tx *sql.Tx, subKegiatanId string) error
}
