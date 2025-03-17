package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/laporan"
)

type LaporanService interface {
	OpdSupportingPokin(ctx context.Context, kodeOpd string, tahun string) (laporan.OpdSupportingPokinResponseData, error)
}
