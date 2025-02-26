package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type SasaranOpdRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx, KodeOpd string, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]domain.SasaranOpd, error)
	FindByIdRencanaKinerja(ctx context.Context, tx *sql.Tx, idRencanaKinerja string) (*domain.SasaranOpd, error)
}
