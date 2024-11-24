package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
	"ekak_kabupaten_madiun/model/web/pegawai"
	"ekak_kabupaten_madiun/repository"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PegawaiServiceImpl struct {
	pegawaiRepository repository.PegawaiRepository
	DB                *sql.DB
}

func NewPegawaiServiceImpl(pegawaiRepository repository.PegawaiRepository, DB *sql.DB) *PegawaiServiceImpl {
	return &PegawaiServiceImpl{
		pegawaiRepository: pegawaiRepository,
		DB:                DB,
	}
}

func (service *PegawaiServiceImpl) Create(ctx context.Context, request pegawai.PegawaiCreateRequest) (pegawai.PegawaiResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pegawai.PegawaiResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Generate UUID dan timestamp untuk ID yang lebih unik
	currentTime := time.Now().Format("20060102")
	uuid := uuid.New().String()
	pegawaiId := fmt.Sprintf("PEG-%s-%s", currentTime, uuid[:8])

	// Debug log untuk memastikan ID ter-generate
	fmt.Printf("Generated ID: %s\n", pegawaiId)

	pegawaiDomain := domainmaster.Pegawai{
		Id:          pegawaiId,
		NamaPegawai: request.NamaPegawai,
		Nip:         request.Nip,
		KodeOpd:     request.KodeOpd,
	}

	pegawais, err := service.pegawaiRepository.Create(ctx, tx, pegawaiDomain)
	if err != nil {
		return pegawai.PegawaiResponse{}, err
	}

	return helper.ToPegawaiResponse(pegawais), nil
}

func (service *PegawaiServiceImpl) Update(ctx context.Context, request pegawai.PegawaiUpdateRequest) (pegawai.PegawaiResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pegawai.PegawaiResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	pegawaiData, err := service.pegawaiRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return pegawai.PegawaiResponse{}, err
	}

	pegawaiData.NamaPegawai = request.NamaPegawai
	pegawaiData.Nip = request.Nip
	pegawaiData.KodeOpd = request.KodeOpd
	updatedPegawai := service.pegawaiRepository.Update(ctx, tx, pegawaiData)
	return helper.ToPegawaiResponse(updatedPegawai), nil
}

func (service *PegawaiServiceImpl) Delete(ctx context.Context, id string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Tambahkan validasi jika id tidak ada
	pegawais, err := service.pegawaiRepository.FindById(ctx, tx, id)
	if err != nil {
		return err
	}
	if pegawais.Id == "" {
		return fmt.Errorf("pegawai dengan id %s tidak ditemukan", id)
	}

	err = service.pegawaiRepository.Delete(ctx, tx, id)
	if err != nil {
		return err
	}

	return nil
}

func (service *PegawaiServiceImpl) FindById(ctx context.Context, id string) (pegawai.PegawaiResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pegawai.PegawaiResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	pegawais, err := service.pegawaiRepository.FindById(ctx, tx, id)
	if err != nil {
		return pegawai.PegawaiResponse{}, err
	}

	return helper.ToPegawaiResponse(pegawais), nil
}

func (service *PegawaiServiceImpl) FindAll(ctx context.Context) ([]pegawai.PegawaiResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []pegawai.PegawaiResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	pegawais, err := service.pegawaiRepository.FindAll(ctx, tx)
	if err != nil {
		return []pegawai.PegawaiResponse{}, err
	}

	return helper.ToPegawaiResponses(pegawais), nil
}
