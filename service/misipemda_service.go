package service

import (
	"context"
	visimisipemda "ekak_kabupaten_madiun/model/web/visimisi"
)

type MisiPemdaService interface {
	Create(ctx context.Context, request visimisipemda.MisiPemdaCreateRequest) (visimisipemda.MisiPemdaResponse, error)
	Update(ctx context.Context, request visimisipemda.MisiPemdaUpdateRequest) (visimisipemda.MisiPemdaResponse, error)
	Delete(ctx context.Context, id int) error
	FindAll(ctx context.Context, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]visimisipemda.VisiPemdaRespons, error)
	FindById(ctx context.Context, id int) (visimisipemda.MisiPemdaResponse, error)
}
