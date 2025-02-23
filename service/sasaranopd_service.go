package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/sasaranopd"
)

type SasaranOpdService interface {
	FindAll(ctx context.Context, KodeOpd string, tahunAwal string, tahunAkhir string) ([]sasaranopd.SasaranOpdResponse, error)
}
