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

type UsulanMandatoriServiceImpl struct {
	usulanMandatoriRepository repository.UsulanMandatoriRepository
	DB                        *sql.DB
}

func NewUsulanMandatoriServiceImpl(usulanMandatoriRepository repository.UsulanMandatoriRepository, DB *sql.DB) *UsulanMandatoriServiceImpl {
	return &UsulanMandatoriServiceImpl{
		usulanMandatoriRepository: usulanMandatoriRepository,
		DB:                        DB,
	}
}

func (service *UsulanMandatoriServiceImpl) Create(ctx context.Context, request usulan.UsulanMandatoriCreateRequest) (usulan.UsulanMandatoriResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return usulan.UsulanMandatoriResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
	uuId := fmt.Sprintf("USU-MAND-%s", randomDigits)

	domainUsulanMandatori := domain.UsulanMandatori{
		Id:               uuId,
		Usulan:           request.Usulan,
		PeraturanTerkait: request.PeraturanTerkait,
		Uraian:           request.Uraian,
		Tahun:            request.Tahun,
		RekinId:          request.RekinId,
		PegawaiId:        request.PegawaiId,
		KodeOpd:          request.KodeOpd,
		Status:           request.Status,
	}

	usulanMandatori, err := service.usulanMandatoriRepository.Create(ctx, tx, domainUsulanMandatori)
	if err != nil {
		return usulan.UsulanMandatoriResponse{}, err
	}

	return helper.ToUsulanMandatoriResponse(usulanMandatori), nil
}

func (service *UsulanMandatoriServiceImpl) Update(ctx context.Context, request usulan.UsulanMandatoriUpdateRequest) (usulan.UsulanMandatoriResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return usulan.UsulanMandatoriResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	existingUsulan, err := service.usulanMandatoriRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return usulan.UsulanMandatoriResponse{}, err
	}

	existingUsulan.Usulan = request.Usulan
	existingUsulan.PeraturanTerkait = request.PeraturanTerkait
	existingUsulan.Uraian = request.Uraian
	existingUsulan.Tahun = request.Tahun
	existingUsulan.PegawaiId = request.PegawaiId
	existingUsulan.KodeOpd = request.KodeOpd
	existingUsulan.Status = request.Status

	updatedUsulan, err := service.usulanMandatoriRepository.Update(ctx, tx, existingUsulan)
	if err != nil {
		return usulan.UsulanMandatoriResponse{}, err
	}

	return helper.ToUsulanMandatoriResponse(updatedUsulan), nil
}

func (service *UsulanMandatoriServiceImpl) FindById(ctx context.Context, idUsulan string) (usulan.UsulanMandatoriResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return usulan.UsulanMandatoriResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	usulanMandatori, err := service.usulanMandatoriRepository.FindById(ctx, tx, idUsulan)
	if err != nil {
		return usulan.UsulanMandatoriResponse{}, err
	}

	return helper.ToUsulanMandatoriResponse(usulanMandatori), nil
}

func (service *UsulanMandatoriServiceImpl) FindAll(ctx context.Context, pegawaiId *string, isActive *bool, rekinId *string) ([]usulan.UsulanMandatoriResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []usulan.UsulanMandatoriResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	usulanMandatoris, err := service.usulanMandatoriRepository.FindAll(ctx, tx, pegawaiId, isActive, rekinId)
	if err != nil {
		return []usulan.UsulanMandatoriResponse{}, err
	}

	return helper.ToUsulanMandatoriResponses(usulanMandatoris), nil
}

func (service *UsulanMandatoriServiceImpl) Delete(ctx context.Context, idUsulan string) error {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	_, err = service.usulanMandatoriRepository.FindById(ctx, tx, idUsulan)
	if err != nil {
		return fmt.Errorf("usulan dengan ID %s tidak ditemukan", idUsulan)
	}

	err = service.usulanMandatoriRepository.Delete(ctx, tx, idUsulan)
	helper.PanicIfError(err)

	return nil
}
