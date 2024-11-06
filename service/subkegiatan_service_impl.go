package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/subkegiatan"
	"ekak_kabupaten_madiun/repository"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type SubKegiatanServiceImpl struct {
	subKegiatanRepository repository.SubKegiatanRepository
	DB                    *sql.DB
	validator             *validator.Validate
}

func NewSubKegiatanServiceImpl(subKegiatanRepository repository.SubKegiatanRepository, DB *sql.DB, validator *validator.Validate) *SubKegiatanServiceImpl {
	return &SubKegiatanServiceImpl{
		subKegiatanRepository: subKegiatanRepository,
		DB:                    DB,
		validator:             validator,
	}
}

func (service *SubKegiatanServiceImpl) Create(ctx context.Context, request subkegiatan.SubKegiatanCreateRequest) (subkegiatan.SubKegiatanResponse, error) {
	err := service.validator.Struct(request)
	if err != nil {
		log.Println("Validasi gagal:", err)
		return subkegiatan.SubKegiatanResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("Gagal memulai transaksi:", err)
		return subkegiatan.SubKegiatanResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
	uuId := fmt.Sprintf("SUB-%s", randomDigits)

	subKegiatan := domain.SubKegiatan{
		Id:              uuId,
		PegawaiId:       request.PegawaiId,
		NamaSubKegiatan: request.NamaSubKegiatan,
		KodeOpd:         request.KodeOpd,
		Tahun:           request.Tahun,
	}

	result, err := service.subKegiatanRepository.Create(ctx, tx, subKegiatan)
	if err != nil {
		log.Println("Gagal membuat data sub kegiatan:", err)
		return subkegiatan.SubKegiatanResponse{}, err
	}

	return helper.ToSubKegiatanResponse(result), nil
}

func (service *SubKegiatanServiceImpl) Update(ctx context.Context, request subkegiatan.SubKegiatanUpdateRequest) (subkegiatan.SubKegiatanResponse, error) {
	err := service.validator.Struct(request)
	if err != nil {
		log.Println("Validasi gagal:", err)
		return subkegiatan.SubKegiatanResponse{}, err
	}
	tx, err := service.DB.Begin()
	if err != nil {
		return subkegiatan.SubKegiatanResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	domainSubKegiatan := domain.SubKegiatan{
		Id:              request.Id,
		NamaSubKegiatan: request.NamaSubKegiatan,
		KodeOpd:         request.KodeOpd,
		Tahun:           request.Tahun,
	}

	subkegiatans, err := service.subKegiatanRepository.Update(ctx, tx, domainSubKegiatan)
	if err != nil {
		return subkegiatan.SubKegiatanResponse{}, err
	}

	response := helper.ToSubKegiatanResponse(subkegiatans)
	return response, nil

}

func (service *SubKegiatanServiceImpl) FindById(ctx context.Context, subKegiatanId string) (subkegiatan.SubKegiatanResponse, error) {

	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("Gagal memulai transaksi:", err)
		return subkegiatan.SubKegiatanResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	result, err := service.subKegiatanRepository.FindById(ctx, tx, subKegiatanId)
	if err != nil {
		log.Println("Gagal mencari data sub kegiatan:", err)
		return subkegiatan.SubKegiatanResponse{}, err
	}

	return helper.ToSubKegiatanResponse(result), nil
}

func (service *SubKegiatanServiceImpl) FindAll(ctx context.Context, kodeOpd, pegawaiId string) ([]subkegiatan.SubKegiatanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("Gagal memulai transaksi:", err)
		return []subkegiatan.SubKegiatanResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	result, err := service.subKegiatanRepository.FindAll(ctx, tx, kodeOpd, pegawaiId)
	if err != nil {
		log.Println("Gagal mencari data sub kegiatan:", err)
		return []subkegiatan.SubKegiatanResponse{}, err
	}

	return helper.ToSubKegiatanResponses(result), nil
}

func (service *SubKegiatanServiceImpl) Delete(ctx context.Context, subKegiatanId string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	err = service.subKegiatanRepository.Delete(ctx, tx, subKegiatanId)
	if err != nil {
		return fmt.Errorf("gagal menghapus data sub kegiatan: %v", err)
	}
	return nil
}
