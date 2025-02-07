package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/periodetahun"
)

type PeriodeService interface {
	Create(ctx context.Context, request periodetahun.PeriodeCreateRequest) (periodetahun.PeriodeResponse, error)
	Update(ctx context.Context, request periodetahun.PeriodeUpdateRequest) (periodetahun.PeriodeResponse, error)
	FindByTahun(ctx context.Context, tahun string) (periodetahun.PeriodeResponse, error)
}
