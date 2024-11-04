package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/opdmaster"
)

type OpdService interface {
	Create(ctx context.Context, request opdmaster.OpdCreateRequest) (opdmaster.OpdResponse, error)
	Update(ctx context.Context, request opdmaster.OpdUpdateRequest) (opdmaster.OpdResponse, error)
	Delete(ctx context.Context, id string) error
	FindById(ctx context.Context, id string) (opdmaster.OpdResponse, error)
	FindByKodeOpd(ctx context.Context, kodeOpd string) (opdmaster.OpdResponse, error)
	FindAll(ctx context.Context) ([]opdmaster.OpdResponse, error)
}
