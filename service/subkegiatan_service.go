package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/subkegiatan"
)

type SubKegiatanService interface {
	Create(ctx context.Context, request subkegiatan.SubKegiatanCreateRequest) (subkegiatan.SubKegiatanResponse, error)
	Update(ctx context.Context, request subkegiatan.SubKegiatanUpdateRequest) (subkegiatan.SubKegiatanResponse, error)
	FindById(ctx context.Context, subKegiatanId string) (subkegiatan.SubKegiatanResponse, error)
	FindAll(ctx context.Context, kodeOpd, rekinId, status string) ([]subkegiatan.SubKegiatanResponse, error)
	Delete(ctx context.Context, subKegiatanId string) error
	CreateRekin(ctx context.Context, request subkegiatan.SubKegiatanCreateRekinRequest) ([]subkegiatan.SubKegiatanResponse, error)
	DeleteSubKegiatanTerpilih(ctx context.Context, idSubKegiatan string) error
}
