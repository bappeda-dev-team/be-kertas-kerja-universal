package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/inovasi"
)

type InovasiService interface {
	Create(ctx context.Context, request inovasi.InovasiCreateRequest) (inovasi.InovasiResponse, error)
	FindById(ctx context.Context, inovasiId string) (inovasi.InovasiResponse, error)
	FindAll(ctx context.Context, rekinId string, pegawaiId string) ([]inovasi.InovasiResponse, error)
	Update(ctx context.Context, request inovasi.InovasiUpdateRequest) (inovasi.InovasiResponse, error)
	Delete(ctx context.Context, inovasiId string) error
}
