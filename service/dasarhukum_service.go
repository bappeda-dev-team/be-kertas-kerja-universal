package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/dasarhukum"
)

type DasarHukumService interface {
	Create(ctx context.Context, request dasarhukum.DasarHukumCreateRequest) (dasarhukum.DasarHukumResponse, error)
	Update(ctx context.Context, request dasarhukum.DasarHukumUpdateRequest) (dasarhukum.DasarHukumResponse, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, rekinId string, pegawaiId string) ([]dasarhukum.DasarHukumResponse, error)
	FindById(ctx context.Context, id string) (dasarhukum.DasarHukumResponse, error)
}
