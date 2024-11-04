package service

import (
	"context"
	"ekak_kabupaten_madiun/model/web/bidangurusanresponse"
)

type BidangUrusanService interface {
	Create(ctx context.Context, request bidangurusanresponse.BidangUrusanCreateRequest) (bidangurusanresponse.BidangUrusanResponse, error)
	Update(ctx context.Context, request bidangurusanresponse.BidangUrusanUpdateRequest) (bidangurusanresponse.BidangUrusanResponse, error)
	Delete(ctx context.Context, id string) error
	FindById(ctx context.Context, id string) (bidangurusanresponse.BidangUrusanResponse, error)
	FindAll(ctx context.Context) ([]bidangurusanresponse.BidangUrusanResponse, error)
}
