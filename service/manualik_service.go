package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/rencanakinerja"
)

type ManualIKService interface {
	Create(ctx context.Context, request rencanakinerja.ManualIKCreateRequest, indikatorId string) (rencanakinerja.ManualIKResponse, error)
	Update(ctx context.Context, request rencanakinerja.ManualIKUpdateRequest, indikatorId string) (rencanakinerja.ManualIKResponse, error)
	FindManualIKByIndikatorId(ctx context.Context, indikatorId string) (rencanakinerja.ManualIKResponse, error)
	// Delete(ctx context.Context, id string) error
	// FindAll(ctx context.Context, rekinId string, pegawaiId string) ([]rencanakinerja.ManualIKResponse, error)
	// FindById(ctx context.Context, id string) (rencanakinerja.ManualIKResponse, error)
}
