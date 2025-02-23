package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type SasaranOpdRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx, KodeOpd string, tahunAwal string, tahunAkhir string) ([]domain.SasaranOpd, error)
}
