package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type SubKegiatanRepository interface {
	Create(ctx context.Context, tx *sql.Tx, subKegiatan domain.SubKegiatan) (domain.SubKegiatan, error)
	FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string, rekinId string, status string) ([]domain.SubKegiatan, error)
	Update(ctx context.Context, tx *sql.Tx, subKegiatan domain.SubKegiatan) (domain.SubKegiatan, error)
	FindById(ctx context.Context, tx *sql.Tx, subKegiatanId string) (domain.SubKegiatan, error)
	Delete(ctx context.Context, tx *sql.Tx, subKegiatanId string) error
	FindIndikatorBySubKegiatanId(ctx context.Context, tx *sql.Tx, subKegiatanId string) ([]domain.Indikator, error)
	FindTargetByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string) ([]domain.Target, error)
	FindByKodeSubKegiatan(ctx context.Context, tx *sql.Tx, kodeSubKegiatan string) (domain.SubKegiatan, error)
	CreateRekin(ctx context.Context, tx *sql.Tx, idUsulan string, rekinId string) error
	DeleteSubKegiatanTerpilih(ctx context.Context, tx *sql.Tx, idSubKegiatan string) error
}
