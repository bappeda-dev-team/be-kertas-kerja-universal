package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/subkegiatan"
	"ekak_kabupaten_madiun/repository"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type SubKegiatanTerpilihServiceImpl struct {
	SubKegiatanTerpilihRepository repository.SubKegiatanTerpilihRepository
	DB                            *sql.DB
}

func NewSubKegiatanTerpilihServiceImpl(subKegiatanTerpilihRepository repository.SubKegiatanTerpilihRepository, DB *sql.DB) *SubKegiatanTerpilihServiceImpl {
	return &SubKegiatanTerpilihServiceImpl{
		SubKegiatanTerpilihRepository: subKegiatanTerpilihRepository,
		DB:                            DB,
	}
}

func (service *SubKegiatanTerpilihServiceImpl) Create(ctx context.Context, request subkegiatan.SubKegiatanTerpilihCreateRequest) (subkegiatan.SubKegiatanTerpilihResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return subkegiatan.SubKegiatanTerpilihResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi keberadaan SubKegiatanId di tb_subkegiatan
	exists, err := service.SubKegiatanTerpilihRepository.ExistsInSubKegiatan(ctx, tx, request.SubKegiatanId)
	if err != nil {
		return subkegiatan.SubKegiatanTerpilihResponse{}, err
	}
	if !exists {
		return subkegiatan.SubKegiatanTerpilihResponse{}, errors.New("subkegiatan dengan ID tersebut tidak ditemukan di tb_subkegiatan")
	}

	randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
	uuId := fmt.Sprintf("USU-MAND-%s", randomDigits)

	subKegiatanTerpilih := domain.SubKegiatanTerpilih{
		Id:               uuId,
		RencanaKinerjaId: request.RencanaKinerjaId,
		SubKegiatanId:    request.SubKegiatanId,
	}

	// Memeriksa duplikasi berdasarkan rekin_id dan subkegiatan_id
	exists, err = service.SubKegiatanTerpilihRepository.ExistsByRekinAndSubKegiatan(ctx, tx, request.RencanaKinerjaId, request.SubKegiatanId)
	if err != nil {
		return subkegiatan.SubKegiatanTerpilihResponse{}, err
	}
	if exists {
		return subkegiatan.SubKegiatanTerpilihResponse{}, errors.New("subkegiatan ini sudah dipilih untuk rekin yang sama")
	}

	result, err := service.SubKegiatanTerpilihRepository.Create(ctx, tx, subKegiatanTerpilih)
	if err != nil {
		return subkegiatan.SubKegiatanTerpilihResponse{}, err
	}

	return helper.ToSubKegiatanTerpilihResponse(result), nil
}

func (service *SubKegiatanTerpilihServiceImpl) Delete(ctx context.Context, subKegiatanTerpilihId string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	err = service.SubKegiatanTerpilihRepository.Delete(ctx, tx, subKegiatanTerpilihId)
	if err != nil {
		return err
	}

	return nil
}
