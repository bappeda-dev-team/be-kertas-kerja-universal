package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/iku"
)

type IkuService interface {
	FindAll(ctx context.Context, tahun string) ([]iku.IkuResponse, error)
}
