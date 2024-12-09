package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type CrosscuttingOpdRepository interface {
	CreateCrosscutting(ctx context.Context, tx *sql.Tx, pokin domain.PohonKinerja, parentId int) (domain.PohonKinerja, error)
	UpdateCrosscutting(ctx context.Context, tx *sql.Tx, pokin domain.PohonKinerja, parentId int) (domain.PohonKinerja, error)
	DeleteCrosscutting(ctx context.Context, tx *sql.Tx, id int) error
	FindAllCrosscutting(ctx context.Context, tx *sql.Tx, crosscuttingId int) ([]domain.PohonKinerja, error)
	ValidateKodeOpdChange(ctx context.Context, tx *sql.Tx, id int) error
	FindTargetByIndikatorIds(ctx context.Context, tx *sql.Tx, indikatorIds []string) ([]domain.Target, error)
	FindIndikatorByPokinId(ctx context.Context, tx *sql.Tx, pokinIds []int) ([]domain.Indikator, error)
	ApproveOrRejectCrosscutting(ctx context.Context, tx *sql.Tx, crosscuttingId int, approve bool) error
}
