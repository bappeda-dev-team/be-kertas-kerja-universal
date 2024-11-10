package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type SubKegiatanTerpilihRepository interface {
	Update(ctx context.Context, tx *sql.Tx, subKegiatanTerpilih domain.SubKegiatanTerpilih) (domain.SubKegiatanTerpilih, error)
	Delete(ctx context.Context, tx *sql.Tx, id string, kodeSubKegiatan string) error
	FindByIdAndKodeSubKegiatan(ctx context.Context, tx *sql.Tx, id string, kodeSubKegiatan string) (domain.SubKegiatanTerpilih, error)
}
