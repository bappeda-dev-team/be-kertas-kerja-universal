package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type PohonKinerjaRepository interface {
	Create(ctx context.Context, tx *sql.Tx, pohonKinerja domain.PohonKinerja) (domain.PohonKinerja, error)
	Update(ctx context.Context, tx *sql.Tx, pohonKinerja domain.PohonKinerja) (domain.PohonKinerja, error)
	Delete(ctx context.Context, tx *sql.Tx, id string) error
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.PohonKinerja, error)
	FindAll(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string) ([]domain.PohonKinerja, error)
}
