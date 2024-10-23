package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type UsulanMandatoriRepository interface {
	Create(ctx context.Context, tx *sql.Tx, usulan domain.UsulanMandatori) (domain.UsulanMandatori, error)
	FindAll(ctx context.Context, tx *sql.Tx, pegawaiId *string, isActive *bool, rekinId *string) ([]domain.UsulanMandatori, error)
	FindById(ctx context.Context, tx *sql.Tx, idUsulan string) (domain.UsulanMandatori, error)
	Update(ctx context.Context, tx *sql.Tx, usulan domain.UsulanMandatori) (domain.UsulanMandatori, error)
	Delete(ctx context.Context, tx *sql.Tx, idUsulan string) error
}
