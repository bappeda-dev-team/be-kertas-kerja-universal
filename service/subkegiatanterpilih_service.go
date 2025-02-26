package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/subkegiatan"
)

type SubKegiatanTerpilihService interface {
	Update(ctx context.Context, request subkegiatan.SubKegiatanTerpilihUpdateRequest) (subkegiatan.SubKegiatanTerpilihResponse, error)
	FindByKodeSubKegiatan(ctx context.Context, kodeSubKegiatan string) (subkegiatan.SubKegiatanTerpilihResponse, error)
	Delete(ctx context.Context, id string, kodeSubKegiatan string) error
	CreateRekin(ctx context.Context, request subkegiatan.SubKegiatanCreateRekinRequest) ([]subkegiatan.SubKegiatanResponse, error)
	DeleteSubKegiatanTerpilih(ctx context.Context, idSubKegiatan string) error
}
