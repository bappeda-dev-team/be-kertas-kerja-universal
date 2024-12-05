package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain"
)

type PohonKinerjaRepository interface {
	//pokin opd
	Create(ctx context.Context, tx *sql.Tx, pohonKinerja domain.PohonKinerja) (domain.PohonKinerja, error)
	Update(ctx context.Context, tx *sql.Tx, pohonKinerja domain.PohonKinerja) (domain.PohonKinerja, error)
	Delete(ctx context.Context, tx *sql.Tx, id string) error
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.PohonKinerja, error)
	FindAll(ctx context.Context, tx *sql.Tx, kodeOpd, tahun string) ([]domain.PohonKinerja, error)
	FindStrategicNoParent(ctx context.Context, tx *sql.Tx, levelPohon, parent int, kodeOpd, tahun string) ([]domain.PohonKinerja, error)
	FindPelaksanaPokin(ctx context.Context, tx *sql.Tx, pohonKinerjaId string) ([]domain.PelaksanaPokin, error)
	DeletePelaksanaPokin(ctx context.Context, tx *sql.Tx, pelaksanaId string) error
	//admin pokin
	CreatePokinAdmin(ctx context.Context, tx *sql.Tx, pokinAdmin domain.PohonKinerja) (domain.PohonKinerja, error)
	UpdatePokinAdmin(ctx context.Context, tx *sql.Tx, pokinAdmin domain.PohonKinerja) (domain.PohonKinerja, error)
	DeletePokinAdmin(ctx context.Context, tx *sql.Tx, id int) error
	FindPokinAdminById(ctx context.Context, tx *sql.Tx, id int) (domain.PohonKinerja, error)
	FindPokinAdminAll(ctx context.Context, tx *sql.Tx, tahun string) ([]domain.PohonKinerja, error)
	FindPokinAdminByIdHierarki(ctx context.Context, tx *sql.Tx, idPokin int) ([]domain.PohonKinerja, error)
	FindIndikatorByPokinId(ctx context.Context, tx *sql.Tx, pokinId string) ([]domain.Indikator, error)
	FindTargetByIndikatorId(ctx context.Context, tx *sql.Tx, indikatorId string) ([]domain.Target, error)
	FindPokinToClone(ctx context.Context, tx *sql.Tx, id int) (domain.PohonKinerja, error)
	ValidateParentLevel(ctx context.Context, tx *sql.Tx, parentId int, levelPohon int) error
	FindIndikatorToClone(ctx context.Context, tx *sql.Tx, pokinId int) ([]domain.Indikator, error)
	FindTargetToClone(ctx context.Context, tx *sql.Tx, indikatorId string) ([]domain.Target, error)
	InsertClonedPokin(ctx context.Context, tx *sql.Tx, pokin domain.PohonKinerja) (int64, error)
	InsertClonedIndikator(ctx context.Context, tx *sql.Tx, indikatorId string, pokinId int64, indikator domain.Indikator) error
	InsertClonedTarget(ctx context.Context, tx *sql.Tx, targetId string, indikatorId string, target domain.Target) error
	GetChildNodes(ctx context.Context, tx *sql.Tx, parentId int) ([]domain.PohonKinerja, error)
	UpdatePokinStatus(ctx context.Context, tx *sql.Tx, id int, status string) error
	CheckPokinStatus(ctx context.Context, tx *sql.Tx, id int) (string, error)
	InsertClonedPokinWithStatus(ctx context.Context, tx *sql.Tx, pokin domain.PohonKinerja) (int64, error)
	UpdatePokinStatusTolak(ctx context.Context, tx *sql.Tx, id int, status string) error
	CheckCloneFrom(ctx context.Context, tx *sql.Tx, id int) (int, error)
	FindPokinByCloneFrom(ctx context.Context, tx *sql.Tx, cloneFromId int) ([]domain.PohonKinerja, error)
	FindIndikatorByCloneFrom(ctx context.Context, tx *sql.Tx, pokinId int, cloneFromId string) (domain.Indikator, error)
	FindTargetByCloneFrom(ctx context.Context, tx *sql.Tx, indikatorId string, cloneFromId string) (domain.Target, error)

	//find pokin for dropdown
	FindPokinByJenisPohon(ctx context.Context, tx *sql.Tx, jenisPohon string, levelPohon int, tahun string, kodeOpd string) ([]domain.PohonKinerja, error)
	FindPokinByPelaksana(ctx context.Context, tx *sql.Tx, pelaksanaId string, tahun string) ([]domain.PohonKinerja, error)
	FindPokinByStatus(ctx context.Context, tx *sql.Tx, kodeOpd string, tahun string, status string) ([]domain.PohonKinerja, error)
}
