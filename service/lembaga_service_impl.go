package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
	"ekak_kabupaten_madiun/model/web/lembaga"
	"ekak_kabupaten_madiun/repository"
	"fmt"

	"github.com/google/uuid"
)

type LembagaServiceImpl struct {
	LembagaRepository repository.LembagaRepository
	DB                *sql.DB
}

func NewLembagaServiceImpl(lembagaRepository repository.LembagaRepository, DB *sql.DB) *LembagaServiceImpl {
	return &LembagaServiceImpl{
		LembagaRepository: lembagaRepository,
		DB:                DB,
	}
}

func (service *LembagaServiceImpl) Create(ctx context.Context, request lembaga.LembagaCreateRequest) (lembaga.LembagaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return lembaga.LembagaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Generate ID unik dengan prefix
	uuid := uuid.New().String()
	lembagaId := fmt.Sprintf("LMBG-%v", uuid[:4])

	// Buat objek domain lembaga
	lembagaDomain := domainmaster.Lembaga{
		Id:          lembagaId,
		NamaLembaga: request.NamaLembaga,
		IsActive:    true, // Tambahkan respons isactive true
	}

	result := service.LembagaRepository.Create(ctx, tx, lembagaDomain)

	response := lembaga.LembagaResponse{
		Id:          result.Id,
		NamaLembaga: result.NamaLembaga,
		IsActive:    result.IsActive, // Tambahkan respons isactive true
	}

	return response, nil
}

func (service *LembagaServiceImpl) Update(ctx context.Context, request lembaga.LembagaUpdateRequest) (lembaga.LembagaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return lembaga.LembagaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Cek apakah data exists
	_, err = service.LembagaRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return lembaga.LembagaResponse{}, err
	}

	lembagaDomain := domainmaster.Lembaga{
		Id:          request.Id,
		NamaLembaga: request.NamaLembaga,
		IsActive:    request.IsActive,
	}

	result := service.LembagaRepository.Update(ctx, tx, lembagaDomain)

	response := lembaga.LembagaResponse{
		Id:          result.Id,
		NamaLembaga: result.NamaLembaga,
		IsActive:    result.IsActive,
	}

	return response, nil
}

func (service *LembagaServiceImpl) FindById(ctx context.Context, id string) (lembaga.LembagaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return lembaga.LembagaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	result, err := service.LembagaRepository.FindById(ctx, tx, id)
	if err != nil {
		return lembaga.LembagaResponse{}, err
	}

	response := lembaga.LembagaResponse{
		Id:          result.Id,
		NamaLembaga: result.NamaLembaga,
		IsActive:    result.IsActive,
	}

	return response, nil
}

func (service *LembagaServiceImpl) FindAll(ctx context.Context) ([]lembaga.LembagaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []lembaga.LembagaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	result, err := service.LembagaRepository.FindAll(ctx, tx)
	if err != nil {
		return []lembaga.LembagaResponse{}, err
	}

	response := []lembaga.LembagaResponse{}
	for _, value := range result {
		response = append(response, lembaga.LembagaResponse{
			Id:          value.Id,
			NamaLembaga: value.NamaLembaga,
			IsActive:    value.IsActive,
		})
	}
	return response, nil
}

func (service *LembagaServiceImpl) Delete(ctx context.Context, id string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	return service.LembagaRepository.Delete(ctx, tx, id)
}
