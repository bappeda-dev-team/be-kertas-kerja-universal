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
		NamaAdmin:     helper.EmptyStringIfNull(request.NamaAdmin),
		NoWaAdmin:     helper.EmptyStringIfNull(request.NoWaAdmin),
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
		NamaAdmin:     result.NamaAdmin,
		NoWaAdmin:     result.NoWaAdmin,
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
	opd.NamaAdmin = request.NamaAdmin
	opd.NoWaAdmin = request.NoWaAdmin
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
		NamaAdmin:     result.NamaAdmin,
		NoWaAdmin:     result.NoWaAdmin,
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
		fmt.Printf("Error saat mencari Lembaga: %v\n", err)
		return opdmaster.OpdResponse{}, err
	}

	// Konversi dari domain ke response
	lembagaResponse := lembaga.LembagaResponse{
		Id:          lembagaDomain.Id,
		NamaLembaga: lembagaDomain.NamaLembaga,
		KodeLembaga: lembagaDomain.KodeLembaga,
	}

	response := opdmaster.OpdResponse{
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
		NamaAdmin:     opd.NamaAdmin,
		NoWaAdmin:     opd.NoWaAdmin,
		IdLembaga:     lembagaResponse,
	}

	return response, nil
}

// TODO: add kode lembaga filter
func (service *OpdServiceImpl) FindAll(ctx context.Context) ([]opdmaster.OpdWithBidangUrusan, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []opdmaster.OpdWithBidangUrusan{}, err
	}
	defer helper.CommitOrRollback(tx)

	opds, err := service.OpdRepository.FindAll(ctx, tx)
	if err != nil {
		return []opdmaster.OpdWithBidangUrusan{}, err
	}

	var opdResponses []opdmaster.OpdWithBidangUrusan
	for _, opd := range opds {
		opdResponses = append(opdResponses, opdmaster.OpdWithBidangUrusan{
			Id:                opd.Id,
			KodeOpd:           opd.KodeOpd,
			NamaOpd:           opd.NamaOpd,
			KodeUrusan1:       opd.KodeUrusan1,
			NamaUrusan1:       opd.NamaUrusan1,
			KodeUrusan2:       opd.KodeUrusan2,
			NamaUrusan2:       opd.NamaUrusan2,
			KodeUrusan3:       opd.KodeUrusan3,
			NamaUrusan3:       opd.NamaUrusan3,
			KodeBidangUrusan1: opd.KodeBidangUrusan1,
			NamaBidangUrusan1: opd.NamaBidangUrusan1,
			KodeBidangUrusan2: opd.KodeBidangUrusan2,
			NamaBidangUrusan2: opd.NamaBidangUrusan2,
			KodeBidangUrusan3: opd.KodeBidangUrusan3,
			NamaBidangUrusan3: opd.NamaBidangUrusan3,
			NamaAdmin:         opd.NamaAdmin,
			NoWaAdmin:         opd.NoWaAdmin,
			NamaKepalaOpd:     opd.NamaKepalaOpd,
			NIPKepalaOpd:      opd.NIPKepalaOpd,
			PangkatKepala:     opd.PangkatKepala,
		})
	}
	return opdResponses, nil
}

func (service *OpdServiceImpl) InfoOpd(ctx context.Context, kodeOpd string, kodeLembaga string) (opdmaster.OpdWithBidangUrusan, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return opdmaster.OpdWithBidangUrusan{}, err
	}
	defer helper.CommitOrRollback(tx)

	opd, err := service.OpdRepository.InfoOpd(ctx, tx, kodeOpd, kodeLembaga)
	if err != nil {
		return opdmaster.OpdWithBidangUrusan{}, err
	}

	opdResponse := opdmaster.OpdWithBidangUrusan{
		Id:                opd.Id,
		KodeOpd:           opd.KodeOpd,
		NamaOpd:           opd.NamaOpd,
		KodeUrusan1:       opd.KodeUrusan1,
		NamaUrusan1:       opd.NamaUrusan1,
		KodeUrusan2:       opd.KodeUrusan2,
		NamaUrusan2:       opd.NamaUrusan2,
		KodeUrusan3:       opd.KodeUrusan3,
		NamaUrusan3:       opd.NamaUrusan3,
		KodeBidangUrusan1: opd.KodeBidangUrusan1,
		NamaBidangUrusan1: opd.NamaBidangUrusan1,
		KodeBidangUrusan2: opd.KodeBidangUrusan2,
		NamaBidangUrusan2: opd.NamaBidangUrusan2,
		KodeBidangUrusan3: opd.KodeBidangUrusan3,
		NamaBidangUrusan3: opd.NamaBidangUrusan3,
		NamaAdmin:         opd.NamaAdmin,
		NoWaAdmin:         opd.NoWaAdmin,
		NamaKepalaOpd:     opd.NamaKepalaOpd,
		NIPKepalaOpd:      opd.NIPKepalaOpd,
		PangkatKepala:     opd.PangkatKepala,
	}

	return opdResponse, nil
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
