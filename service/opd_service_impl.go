package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain/domainmaster"

	"ekak_kabupaten_madiun/model/web/lembaga"
	"ekak_kabupaten_madiun/model/web/opdmaster"
	"ekak_kabupaten_madiun/repository"

	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type OpdServiceImpl struct {
	OpdRepository     repository.OpdRepository
	LembagaRepository repository.LembagaRepository
	DB                *sql.DB
	Validator         *validator.Validate
}

func NewOpdServiceImpl(opdRepository repository.OpdRepository, lembagaRepository repository.LembagaRepository, DB *sql.DB, validator *validator.Validate) *OpdServiceImpl {
	return &OpdServiceImpl{
		OpdRepository:     opdRepository,
		LembagaRepository: lembagaRepository,
		DB:                DB,
		Validator:         validator,
	}
}

func (service *OpdServiceImpl) Create(ctx context.Context, request opdmaster.OpdCreateRequest) (opdmaster.OpdResponse, error) {
	err := service.Validator.Struct(request)
	if err != nil {
		return opdmaster.OpdResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return opdmaster.OpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Tambahkan validasi ID Lembaga
	_, err = service.LembagaRepository.FindById(ctx, tx, request.IdLembaga)
	if err != nil {
		return opdmaster.OpdResponse{}, fmt.Errorf("id lembaga tidak ditemukan")
	}

	uuid := uuid.New()
	opdId := fmt.Sprintf("OPD-%s", uuid.String()[:4])

	opd := domainmaster.Opd{
		Id:            opdId,
		KodeOpd:       request.KodeOpd,
		NamaOpd:       request.NamaOpd,
		Singkatan:     helper.EmptyStringIfNull(request.Singkatan),
		Alamat:        helper.EmptyStringIfNull(request.Alamat),
		Telepon:       helper.EmptyStringIfNull(request.Telepon),
		Fax:           helper.EmptyStringIfNull(request.Fax),
		Email:         helper.EmptyStringIfNull(request.Email),
		Website:       helper.EmptyStringIfNull(request.Website),
		NamaKepalaOpd: request.NamaKepalaOpd,
		NIPKepalaOpd:  request.NIPKepalaOpd,
		PangkatKepala: request.PangkatKepala,
		IdLembaga:     request.IdLembaga,
	}

	result, err := service.OpdRepository.Create(ctx, tx, opd)
	if err != nil {
		return opdmaster.OpdResponse{}, err
	}

	lembagaDomain, err := service.LembagaRepository.FindById(ctx, tx, opd.IdLembaga)
	if err != nil {
		return opdmaster.OpdResponse{}, err
	}

	lembagaResponse := lembaga.LembagaResponse{
		Id:          lembagaDomain.Id,
		NamaLembaga: lembagaDomain.NamaLembaga,
	}

	return opdmaster.OpdResponse{
		Id:            result.Id,
		KodeOpd:       result.KodeOpd,
		NamaOpd:       result.NamaOpd,
		Singkatan:     result.Singkatan,
		Alamat:        result.Alamat,
		Telepon:       result.Telepon,
		Fax:           result.Fax,
		Email:         result.Email,
		Website:       result.Website,
		NamaKepalaOpd: result.NamaKepalaOpd,
		NIPKepalaOpd:  result.NIPKepalaOpd,
		PangkatKepala: result.PangkatKepala,
		IdLembaga:     lembagaResponse,
	}, nil
}

func (service *OpdServiceImpl) Update(ctx context.Context, request opdmaster.OpdUpdateRequest) (opdmaster.OpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return opdmaster.OpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi keberadaan OPD
	opd, err := service.OpdRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return opdmaster.OpdResponse{}, err
	}

	// Tambahkan validasi ID Lembaga
	_, err = service.LembagaRepository.FindById(ctx, tx, request.IdLembaga)
	if err != nil {
		return opdmaster.OpdResponse{}, fmt.Errorf("id lembaga tidak ditemukan")
	}

	opd.KodeOpd = request.KodeOpd
	opd.NamaOpd = request.NamaOpd
	opd.Singkatan = helper.EmptyStringIfNull(request.Singkatan)
	opd.Alamat = helper.EmptyStringIfNull(request.Alamat)
	opd.Telepon = helper.EmptyStringIfNull(request.Telepon)
	opd.Fax = helper.EmptyStringIfNull(request.Fax)
	opd.Email = helper.EmptyStringIfNull(request.Email)
	opd.Website = helper.EmptyStringIfNull(request.Website)
	opd.NamaKepalaOpd = request.NamaKepalaOpd
	opd.NIPKepalaOpd = request.NIPKepalaOpd
	opd.PangkatKepala = request.PangkatKepala
	opd.IdLembaga = request.IdLembaga

	result, err := service.OpdRepository.Update(ctx, tx, opd)
	if err != nil {
		return opdmaster.OpdResponse{}, err
	}

	lembagaDomain, err := service.LembagaRepository.FindById(ctx, tx, opd.IdLembaga)
	if err != nil {
		return opdmaster.OpdResponse{}, err
	}

	lembagaResponse := lembaga.LembagaResponse{
		Id:          lembagaDomain.Id,
		NamaLembaga: lembagaDomain.NamaLembaga,
	}

	return opdmaster.OpdResponse{
		Id:            result.Id,
		KodeOpd:       result.KodeOpd,
		NamaOpd:       result.NamaOpd,
		Singkatan:     result.Singkatan,
		Alamat:        result.Alamat,
		Telepon:       result.Telepon,
		Fax:           result.Fax,
		Email:         result.Email,
		Website:       result.Website,
		NamaKepalaOpd: result.NamaKepalaOpd,
		NIPKepalaOpd:  result.NIPKepalaOpd,
		PangkatKepala: result.PangkatKepala,
		IdLembaga:     lembagaResponse,
	}, nil
}

func (service *OpdServiceImpl) Delete(ctx context.Context, opdId string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi keberadaan ID
	_, err = service.OpdRepository.FindById(ctx, tx, opdId)
	if err != nil {
		return err // Akan mengembalikan error jika ID tidak ditemukan
	}

	return service.OpdRepository.Delete(ctx, tx, opdId)
}

func (service *OpdServiceImpl) FindById(ctx context.Context, opdId string) (opdmaster.OpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return opdmaster.OpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	opd, err := service.OpdRepository.FindById(ctx, tx, opdId)
	if err != nil {
		return opdmaster.OpdResponse{}, err
	}

	lembagaDomain, err := service.LembagaRepository.FindById(ctx, tx, opd.IdLembaga)
	if err != nil {
		return opdmaster.OpdResponse{}, err
	}

	// Konversi dari domain ke response
	lembagaResponse := lembaga.LembagaResponse{
		Id:          lembagaDomain.Id,
		NamaLembaga: lembagaDomain.NamaLembaga,
		// ... sesuaikan dengan field lainnya
	}

	return opdmaster.OpdResponse{
		Id:            opd.Id,
		KodeOpd:       opd.KodeOpd,
		NamaOpd:       opd.NamaOpd,
		Singkatan:     opd.Singkatan,
		Alamat:        opd.Alamat,
		Telepon:       opd.Telepon,
		Fax:           opd.Fax,
		Email:         opd.Email,
		Website:       opd.Website,
		NamaKepalaOpd: opd.NamaKepalaOpd,
		NIPKepalaOpd:  opd.NIPKepalaOpd,
		PangkatKepala: opd.PangkatKepala,
		IdLembaga:     lembagaResponse,
	}, nil
}

func (service *OpdServiceImpl) FindAll(ctx context.Context) ([]opdmaster.OpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []opdmaster.OpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	opds, err := service.OpdRepository.FindAll(ctx, tx)
	if err != nil {
		return []opdmaster.OpdResponse{}, err
	}

	var opdResponses []opdmaster.OpdResponse
	for _, opd := range opds {
		lembagaDomain, err := service.LembagaRepository.FindById(ctx, tx, opd.IdLembaga)
		if err != nil {
			return []opdmaster.OpdResponse{}, err
		}

		lembagaResponse := lembaga.LembagaResponse{
			Id:          lembagaDomain.Id,
			NamaLembaga: lembagaDomain.NamaLembaga,
		}

		opdResponses = append(opdResponses, opdmaster.OpdResponse{
			Id:            opd.Id,
			KodeOpd:       opd.KodeOpd,
			NamaOpd:       opd.NamaOpd,
			Singkatan:     opd.Singkatan,
			Alamat:        opd.Alamat,
			Telepon:       opd.Telepon,
			Fax:           opd.Fax,
			Email:         opd.Email,
			Website:       opd.Website,
			NamaKepalaOpd: opd.NamaKepalaOpd,
			NIPKepalaOpd:  opd.NIPKepalaOpd,
			PangkatKepala: opd.PangkatKepala,
			IdLembaga:     lembagaResponse,
		})
	}
	return opdResponses, nil
}

func (service *OpdServiceImpl) FindByKodeOpd(ctx context.Context, kodeOpd string) (opdmaster.OpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return opdmaster.OpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	opd, err := service.OpdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return opdmaster.OpdResponse{}, err
	}

	return opdmaster.OpdResponse{
		Id:            opd.Id,
		KodeOpd:       opd.KodeOpd,
		NamaOpd:       opd.NamaOpd,
		Singkatan:     opd.Singkatan,
		Alamat:        opd.Alamat,
		Telepon:       opd.Telepon,
		Fax:           opd.Fax,
		Email:         opd.Email,
		Website:       opd.Website,
		NamaKepalaOpd: opd.NamaKepalaOpd,
		NIPKepalaOpd:  opd.NIPKepalaOpd,
		PangkatKepala: opd.PangkatKepala,
	}, nil
}
