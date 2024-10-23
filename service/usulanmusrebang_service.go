package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/usulan"
)

type UsulanMusrebangService interface {
	Create(ctx context.Context, request usulan.UsulanMusrebangCreateRequest) (usulan.UsulanMusrebangResponse, error)
	Update(ctx context.Context, request usulan.UsulanMusrebangUpdateRequest) (usulan.UsulanMusrebangResponse, error)
	FindById(ctx context.Context, idUsulan string) (usulan.UsulanMusrebangResponse, error)
	FindAll(ctx context.Context, pegawaiId *string, is_active *bool, rekinId *string) ([]usulan.UsulanMusrebangResponse, error)
	Delete(ctx context.Context, idUsulan string) error
}
