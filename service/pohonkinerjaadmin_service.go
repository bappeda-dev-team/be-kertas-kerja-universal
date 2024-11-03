package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
)

type PohonKinerjaAdminService interface {
	Create(ctx context.Context, request pohonkinerja.PohonKinerjaAdminCreateRequest) (pohonkinerja.PohonKinerjaAdminResponseData, error)
	Update(ctx context.Context, request pohonkinerja.PohonKinerjaAdminUpdateRequest) error
	Delete(ctx context.Context, id int) error
	FindById(ctx context.Context, id int) (pohonkinerja.PohonKinerjaAdminResponse, error)
	FindAll(ctx context.Context, tahun string) (pohonkinerja.PohonKinerjaAdminResponse, error)
}
