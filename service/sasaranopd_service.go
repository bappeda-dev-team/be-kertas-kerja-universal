package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"ekak_kabupaten_madiun/model/web/sasaranopd"
)

type SasaranOpdService interface {
	FindAll(ctx context.Context, KodeOpd string, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]sasaranopd.SasaranOpdResponse, error)
	FindByIdRencanaKinerja(ctx context.Context, idRencanaKinerja string) (*sasaranopd.SasaranOpdResponse, error)
	FindIdPokinSasaran(ctx context.Context, id int) (pohonkinerja.PohonKinerjaOpdResponse, error)
}
