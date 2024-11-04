package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
	"ekak_kabupaten_madiun/model/web/urusanrespon"
	"ekak_kabupaten_madiun/repository"
	"fmt"

	"github.com/google/uuid"
)

type UrusanServiceImpl struct {
	UrusanRepository repository.UrusanRepository
	DB               *sql.DB
}

func NewUrusanServiceImpl(urusanRepository repository.UrusanRepository, DB *sql.DB) *UrusanServiceImpl {
	return &UrusanServiceImpl{
		UrusanRepository: urusanRepository,
		DB:               DB,
	}
}

func (service *UrusanServiceImpl) Create(ctx context.Context, request urusanrespon.UrusanCreateRequest) (urusanrespon.UrusanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return urusanrespon.UrusanResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
	uuId := fmt.Sprintf("URU-%s", randomDigits)

	domainUrusan := domainmaster.Urusan{
		Id:         uuId,
		KodeUrusan: request.KodeUrusan,
		NamaUrusan: request.NamaUrusan,
	}

	urusan, err := service.UrusanRepository.Create(ctx, tx, domainUrusan)
	if err != nil {
		return urusanrespon.UrusanResponse{}, err
	}

	return urusanrespon.UrusanResponse{
		Id:         urusan.Id,
		KodeUrusan: urusan.KodeUrusan,
		NamaUrusan: urusan.NamaUrusan,
	}, nil
}

func (service *UrusanServiceImpl) Update(ctx context.Context, request urusanrespon.UrusanUpdateRequest) (urusanrespon.UrusanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return urusanrespon.UrusanResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	domainUrusan := domainmaster.Urusan{
		Id:         request.Id,
		KodeUrusan: request.KodeUrusan,
		NamaUrusan: request.NamaUrusan,
	}

	urusan, err := service.UrusanRepository.Update(ctx, tx, domainUrusan)
	if err != nil {
		return urusanrespon.UrusanResponse{}, err
	}

	return urusanrespon.UrusanResponse{
		Id:         urusan.Id,
		KodeUrusan: urusan.KodeUrusan,
		NamaUrusan: urusan.NamaUrusan,
	}, nil
}

func (service *UrusanServiceImpl) FindById(ctx context.Context, id string) (urusanrespon.UrusanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return urusanrespon.UrusanResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	urusan, err := service.UrusanRepository.FindById(ctx, tx, id)
	if err != nil {
		return urusanrespon.UrusanResponse{}, err
	}

	return urusanrespon.UrusanResponse{
		Id:         urusan.Id,
		KodeUrusan: urusan.KodeUrusan,
		NamaUrusan: urusan.NamaUrusan,
	}, nil
}

func (service *UrusanServiceImpl) FindAll(ctx context.Context) ([]urusanrespon.UrusanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []urusanrespon.UrusanResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	urusans, err := service.UrusanRepository.FindAll(ctx, tx)
	if err != nil {
		return []urusanrespon.UrusanResponse{}, err
	}

	var urusanResponses []urusanrespon.UrusanResponse
	for _, urusan := range urusans {
		urusanResponses = append(urusanResponses, urusanrespon.UrusanResponse{
			Id:         urusan.Id,
			KodeUrusan: urusan.KodeUrusan,
			NamaUrusan: urusan.NamaUrusan,
		})
	}

	return urusanResponses, nil
}

func (service *UrusanServiceImpl) Delete(ctx context.Context, id string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Cek apakah urusan dengan ID tersebut ada
	_, err = service.UrusanRepository.FindById(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("urusan dengan id %s tidak ditemukan", id)
	}

	return service.UrusanRepository.Delete(ctx, tx, id)
}
