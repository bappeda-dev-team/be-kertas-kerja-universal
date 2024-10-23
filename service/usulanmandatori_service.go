package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/usulan"
)

type UsulanMandatoriService interface {
	Create(ctx context.Context, request usulan.UsulanMandatoriCreateRequest) (usulan.UsulanMandatoriResponse, error)
	Update(ctx context.Context, request usulan.UsulanMandatoriUpdateRequest) (usulan.UsulanMandatoriResponse, error)
	FindById(ctx context.Context, idUsulan string) (usulan.UsulanMandatoriResponse, error)
	FindAll(ctx context.Context, pegawaiId *string, isActive *bool, rekinId *string) ([]usulan.UsulanMandatoriResponse, error)
	Delete(ctx context.Context, idUsulan string) error
}
