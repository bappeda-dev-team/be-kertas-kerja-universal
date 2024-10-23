package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type UsulanMusrebangRepository interface {
	Create(ctx context.Context, tx *sql.Tx, usulan domain.UsulanMusrebang) (domain.UsulanMusrebang, error)
	Update(ctx context.Context, tx *sql.Tx, usulan domain.UsulanMusrebang) (domain.UsulanMusrebang, error)
	FindById(ctx context.Context, tx *sql.Tx, idUsulan string) (domain.UsulanMusrebang, error)
	FindAll(ctx context.Context, tx *sql.Tx, pegawaiId *string, is_active *bool, rekinId *string) ([]domain.UsulanMusrebang, error)
	Delete(ctx context.Context, tx *sql.Tx, idUsulan string) error
}
