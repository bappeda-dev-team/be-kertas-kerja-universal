package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
)

type PohonKinerjaAdminService interface {
	Create(ctx context.Context, request pohonkinerja.PohonKinerjaAdminCreateRequest) (pohonkinerja.PohonKinerjaAdminResponseData, error)
	Update(ctx context.Context, request pohonkinerja.PohonKinerjaAdminUpdateRequest) (pohonkinerja.PohonKinerjaAdminResponseData, error)
	Delete(ctx context.Context, id int) error
	FindById(ctx context.Context, id int) (pohonkinerja.PohonKinerjaAdminResponseData, error)
	FindAll(ctx context.Context, tahun string) (pohonkinerja.PohonKinerjaAdminResponse, error)
	FindSubTematik(ctx context.Context, tahun string) (pohonkinerja.PohonKinerjaAdminResponse, error)
	FindPokinAdminByIdHierarki(ctx context.Context, idPokin int) (pohonkinerja.PohonKinerjaAdminResponse, error)
}
