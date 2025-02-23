package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	visimisipemda "ekak_kabupaten_madiun/model/web/visimisi"
	"ekak_kabupaten_madiun/repository"

	"github.com/go-playground/validator/v10"
)

type VisiPemdaServiceImpl struct {
	VisiPemdaRepository repository.VisiPemdaRepository
	Validate            *validator.Validate
	DB                  *sql.DB
}

func NewVisiPemdaServiceImpl(visiPemdaRepository repository.VisiPemdaRepository, validate *validator.Validate, DB *sql.DB) *VisiPemdaServiceImpl {
	return &VisiPemdaServiceImpl{
		VisiPemdaRepository: visiPemdaRepository,
		Validate:            validate,
		DB:                  DB,
	}
}

func (service *VisiPemdaServiceImpl) Create(ctx context.Context, request visimisipemda.VisiPemdaCreateRequest) (visimisipemda.VisiPemdaResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return visimisipemda.VisiPemdaResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return visimisipemda.VisiPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	visiPemda := domain.VisiPemda{
		Visi:              request.Visi,
		TahunAwalPeriode:  request.TahunAwalPeriode,
		TahunAkhirPeriode: request.TahunAkhirPeriode,
		JenisPeriode:      request.JenisPeriode,
		Keterangan:        request.Keterangan,
	}

	visiPemda, err = service.VisiPemdaRepository.Create(ctx, tx, visiPemda)
	if err != nil {
		return visimisipemda.VisiPemdaResponse{}, err
	}

	return helper.ToVisiPemdaResponse(visiPemda), nil
}

func (service *VisiPemdaServiceImpl) Update(ctx context.Context, request visimisipemda.VisiPemdaUpdateRequest) (visimisipemda.VisiPemdaResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return visimisipemda.VisiPemdaResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return visimisipemda.VisiPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	visiPemda, err := service.VisiPemdaRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return visimisipemda.VisiPemdaResponse{}, err
	}

	visiPemda.Visi = request.Visi
	visiPemda.TahunAwalPeriode = request.TahunAwalPeriode
	visiPemda.TahunAkhirPeriode = request.TahunAkhirPeriode
	visiPemda.JenisPeriode = request.JenisPeriode
	visiPemda.Keterangan = request.Keterangan

	visiPemda, err = service.VisiPemdaRepository.Update(ctx, tx, visiPemda)
	if err != nil {
		return visimisipemda.VisiPemdaResponse{}, err
	}

	return helper.ToVisiPemdaResponse(visiPemda), nil
}

func (service *VisiPemdaServiceImpl) Delete(ctx context.Context, visiPemdaId int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	visiPemda, err := service.VisiPemdaRepository.FindById(ctx, tx, visiPemdaId)
	if err != nil {
		return err
	}

	return service.VisiPemdaRepository.Delete(ctx, tx, visiPemda.Id)
}

func (service *VisiPemdaServiceImpl) FindById(ctx context.Context, visiPemdaId int) (visimisipemda.VisiPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return visimisipemda.VisiPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	visiPemda, err := service.VisiPemdaRepository.FindById(ctx, tx, visiPemdaId)
	if err != nil {
		return visimisipemda.VisiPemdaResponse{}, err
	}

	return helper.ToVisiPemdaResponse(visiPemda), nil
}

func (service *VisiPemdaServiceImpl) FindAll(ctx context.Context, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]visimisipemda.VisiPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	visiPemdaList, err := service.VisiPemdaRepository.FindAll(ctx, tx, tahunAwal, tahunAkhir, jenisPeriode)
	if err != nil {
		return nil, err
	}

	return helper.ToVisiPemdaResponses(visiPemdaList), nil
}
