package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/subkegiatan"
)

type SubKegiatanTerpilihService interface {
	Create(ctx context.Context, request subkegiatan.SubKegiatanTerpilihCreateRequest) (subkegiatan.SubKegiatanTerpilihResponse, error)
	Delete(ctx context.Context, id string) error
}
