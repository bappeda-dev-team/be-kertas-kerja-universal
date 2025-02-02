package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/tujuanpemda"
)

type TujuanPemdaService interface {
	Create(ctx context.Context, request tujuanpemda.TujuanPemdaCreateRequest) (tujuanpemda.TujuanPemdaResponse, error)
	Update(ctx context.Context, request tujuanpemda.TujuanPemdaUpdateRequest) (tujuanpemda.TujuanPemdaResponse, error)
	Delete(ctx context.Context, id int) error
	FindById(ctx context.Context, tujuanPemdaId int) (tujuanpemda.TujuanPemdaResponse, error)
	FindAll(ctx context.Context, tahun string) ([]tujuanpemda.TujuanPemdaResponse, error)
}
