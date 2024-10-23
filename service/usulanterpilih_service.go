package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/usulan"
)

type UsulanTerpilihService interface {
	Create(ctx context.Context, request usulan.UsulanTerpilihCreateRequest) (usulan.UsulanTerpilihResponse, error)
	Delete(ctx context.Context, idUsulan string) error
}
