package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/iku"
)

type IkuService interface {
	FindAll(ctx context.Context, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]iku.IkuResponse, error)
}
