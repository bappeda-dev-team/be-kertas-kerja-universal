package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type UsulanPokokPikiranRepository interface {
	Create(ctx context.Context, tx *sql.Tx, usulan domain.UsulanPokokPikiran) (domain.UsulanPokokPikiran, error)
	Update(ctx context.Context, tx *sql.Tx, usulan domain.UsulanPokokPikiran) (domain.UsulanPokokPikiran, error)
	FindAll(ctx context.Context, tx *sql.Tx, pegawaiId *string, isActive *bool, rekinId *string) ([]domain.UsulanPokokPikiran, error)
	FindById(ctx context.Context, tx *sql.Tx, idUsulan string) (domain.UsulanPokokPikiran, error)
	Delete(ctx context.Context, tx *sql.Tx, idUsulan string) error
}
