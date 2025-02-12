package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type IkuRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx, tahun string) ([]domain.Indikator, error)
}
