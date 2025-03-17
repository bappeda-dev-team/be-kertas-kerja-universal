package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type LaporanRepository interface {
	OpdSupportingPokin(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domain.OpdSupportingPokin, error)
}
