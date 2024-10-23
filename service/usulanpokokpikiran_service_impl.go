package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/usulan"
	"ekak_kabupaten_madiun/repository"
	"fmt"

	"github.com/google/uuid"
)

type UsulanPokokPikiranServiceImpl struct {
	UsulanPokokPikiranRepository repository.UsulanPokokPikiranRepository
	DB                           *sql.DB
}

func NewUsulanPokokPikiranServiceImpl(usulanPokokPikiranRepository repository.UsulanPokokPikiranRepository, DB *sql.DB) *UsulanPokokPikiranServiceImpl {
	return &UsulanPokokPikiranServiceImpl{
		UsulanPokokPikiranRepository: usulanPokokPikiranRepository,
		DB:                           DB}
}

func (service *UsulanPokokPikiranServiceImpl) Create(ctx context.Context, request usulan.UsulanPokokPikiranCreateRequest) (usulan.UsulanPokokPikiranResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return usulan.UsulanPokokPikiranResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
	uuId := fmt.Sprintf("USU-POKIR-%s", randomDigits)

	domainUsulanPokokPikiran := domain.UsulanPokokPikiran{
		Id:        uuId,
		Usulan:    request.Usulan,
		Alamat:    request.Alamat,
		Uraian:    request.Uraian,
		Tahun:     request.Tahun,
		RekinId:   request.RekinId,
		PegawaiId: request.PegawaiId,
		KodeOpd:   request.KodeOpd,
		Status:    request.Status,
	}

	usulanPokokPikiran, err := service.UsulanPokokPikiranRepository.Create(ctx, tx, domainUsulanPokokPikiran)
	if err != nil {
		return usulan.UsulanPokokPikiranResponse{}, err
	}

	return helper.ToUsulanPokokPikiranResponse(usulanPokokPikiran), nil
}

func (service *UsulanPokokPikiranServiceImpl) Update(ctx context.Context, request usulan.UsulanPokokPikiranUpdateRequest) (usulan.UsulanPokokPikiranResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return usulan.UsulanPokokPikiranResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	usulans, err := service.UsulanPokokPikiranRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return usulan.UsulanPokokPikiranResponse{}, err
	}

	usulans.Usulan = request.Usulan
	usulans.Alamat = request.Alamat
	usulans.Uraian = request.Uraian
	usulans.Tahun = request.Tahun
	usulans.PegawaiId = request.PegawaiId
	usulans.KodeOpd = request.KodeOpd
	usulans.Status = request.Status

	updatedUsulan, err := service.UsulanPokokPikiranRepository.Update(ctx, tx, usulans)
	if err != nil {
		return usulan.UsulanPokokPikiranResponse{}, err
	}

	return helper.ToUsulanPokokPikiranResponse(updatedUsulan), nil
}

func (service *UsulanPokokPikiranServiceImpl) FindById(ctx context.Context, idUsulan string) (usulan.UsulanPokokPikiranResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return usulan.UsulanPokokPikiranResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	usulanPokokPikiran, err := service.UsulanPokokPikiranRepository.FindById(ctx, tx, idUsulan)
	if err != nil {
		return usulan.UsulanPokokPikiranResponse{}, err
	}

	return helper.ToUsulanPokokPikiranResponse(usulanPokokPikiran), nil
}

func (service *UsulanPokokPikiranServiceImpl) FindAll(ctx context.Context, pegawaiId *string, isActive *bool, rekinId *string) ([]usulan.UsulanPokokPikiranResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []usulan.UsulanPokokPikiranResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	usulanPokokPikirans, err := service.UsulanPokokPikiranRepository.FindAll(ctx, tx, pegawaiId, isActive, rekinId)
	if err != nil {
		return []usulan.UsulanPokokPikiranResponse{}, err
	}

	return helper.ToUsulanPokokPikiranResponses(usulanPokokPikirans), nil
}

func (service *UsulanPokokPikiranServiceImpl) Delete(ctx context.Context, idUsulan string) error {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	_, err = service.UsulanPokokPikiranRepository.FindById(ctx, tx, idUsulan)
	if err != nil {
		return fmt.Errorf("usulan dengan ID %s tidak ditemukan", idUsulan)
	}

	err = service.UsulanPokokPikiranRepository.Delete(ctx, tx, idUsulan)
	helper.PanicIfError(err)

	return nil
}
