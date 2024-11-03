package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/jabatan"
)

type JabatanService interface {
	Create(ctx context.Context, request jabatan.JabatanCreateRequest) jabatan.JabatanResponse
	Update(ctx context.Context, request jabatan.JabatanUpdateRequest) jabatan.JabatanResponse
	Delete(ctx context.Context, id string) error
	FindById(ctx context.Context, id string) (jabatan.JabatanResponse, error)
	FindAll(ctx context.Context, kodeOpd string, tahun string) ([]jabatan.JabatanResponse, error)
}
