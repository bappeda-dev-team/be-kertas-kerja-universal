package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type SubKegiatanTerpilihRepository interface {
	Create(ctx context.Context, tx *sql.Tx, subKegiatanTerpilih domain.SubKegiatanTerpilih) (domain.SubKegiatanTerpilih, error)
	Delete(ctx context.Context, tx *sql.Tx, subKegiatanTerpilihId string) error
	FindById(ctx context.Context, tx *sql.Tx, subKegiatanTerpilihId string) (domain.SubKegiatanTerpilih, error)
	ExistsByRekinAndSubKegiatan(ctx context.Context, tx *sql.Tx, rekinId string, subKegiatanId string) (bool, error)
	ExistsInSubKegiatan(ctx context.Context, tx *sql.Tx, subKegiatanId string) (bool, error)
}
