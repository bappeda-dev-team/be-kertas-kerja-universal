package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
	"ekak_kabupaten_madiun/model/web/jabatan"
	"ekak_kabupaten_madiun/repository"
)

type JabatanServiceImpl struct {
	jabatanRepository repository.JabatanRepository
	DB                *sql.DB
}

func NewJabatanServiceImpl(jabatanRepository repository.JabatanRepository, DB *sql.DB) *JabatanServiceImpl {
	return &JabatanServiceImpl{
		jabatanRepository: jabatanRepository,
		DB:                DB,
	}
}

func (service *JabatanServiceImpl) Create(ctx context.Context, request jabatan.JabatanCreateRequest) jabatan.JabatanResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	jabatanDomain := domainmaster.Jabatan{
		NamaJabatan:  request.NamaJabatan,
		KelasJabatan: request.KelasJabatan,
		JenisJabatan: request.JenisJabatan,
		NilaiJabatan: request.NilaiJabatan,
		KodeOpd:      request.KodeOpd,
		IndexJabatan: request.IndexJabatan,
		Tahun:        request.Tahun,
		Esselon:      request.Esselon,
	}

	jabatan := service.jabatanRepository.Create(ctx, tx, jabatanDomain)
	return helper.ToJabatanResponse(jabatan)
}

func (service *JabatanServiceImpl) Update(ctx context.Context, request jabatan.JabatanUpdateRequest) jabatan.JabatanResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	jabatanDomain := domainmaster.Jabatan{
		Id:           request.Id,
		NamaJabatan:  request.NamaJabatan,
		KelasJabatan: request.KelasJabatan,
		JenisJabatan: request.JenisJabatan,
		NilaiJabatan: request.NilaiJabatan,
		KodeOpd:      request.KodeOpd,
		IndexJabatan: request.IndexJabatan,
		Tahun:        request.Tahun,
		Esselon:      request.Esselon,
	}

	jabatan := service.jabatanRepository.Update(ctx, tx, jabatanDomain)
	return helper.ToJabatanResponse(jabatan)
}

func (service *JabatanServiceImpl) Delete(ctx context.Context, id string) error {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	err = service.jabatanRepository.Delete(ctx, tx, id)
	helper.PanicIfError(err)
	return nil
}

func (service *JabatanServiceImpl) FindById(ctx context.Context, id string) (jabatan.JabatanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return jabatan.JabatanResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	jabatans, err := service.jabatanRepository.FindById(ctx, tx, id)
	if err != nil {
		return jabatan.JabatanResponse{}, err
	}

	return helper.ToJabatanResponse(jabatans), nil
}

func (service *JabatanServiceImpl) FindAll(ctx context.Context, kodeOpd string, tahun string) ([]jabatan.JabatanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	jabatans, err := service.jabatanRepository.FindAll(ctx, tx, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}

	return helper.ToJabatanResponses(jabatans), nil
}
