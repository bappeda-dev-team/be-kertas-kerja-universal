package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
)

type CrosscuttingOpdService interface {
	Create(ctx context.Context, request pohonkinerja.CrosscuttingOpdCreateRequest, parentId int) (pohonkinerja.CrosscuttingOpdResponse, error)
	Update(ctx context.Context, request pohonkinerja.CrosscuttingOpdUpdateRequest) (pohonkinerja.CrosscuttingOpdResponse, error)
	FindAllByParent(ctx context.Context, parentId int) ([]pohonkinerja.CrosscuttingOpdResponse, error)
	ApproveOrReject(ctx context.Context, crosscuttingId int, request pohonkinerja.CrosscuttingApproveRequest) (*pohonkinerja.CrosscuttingApproveResponse, error)
	Delete(ctx context.Context, crosscuttingId int) error
}
