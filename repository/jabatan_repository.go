package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
)

type JabatanRepository interface {
	Create(ctx context.Context, tx *sql.Tx, jabatan domainmaster.Jabatan) domainmaster.Jabatan
	Update(ctx context.Context, tx *sql.Tx, jabatan domainmaster.Jabatan) domainmaster.Jabatan
	Delete(ctx context.Context, tx *sql.Tx, id string) error
	FindById(ctx context.Context, tx *sql.Tx, id string) (domainmaster.Jabatan, error)
	FindAll(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string) ([]domainmaster.Jabatan, error)
}
