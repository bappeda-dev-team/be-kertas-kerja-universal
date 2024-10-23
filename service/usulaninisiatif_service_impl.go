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

type UsulanInisiatifServiceImpl struct {
	UsulanInisiatifRepository repository.UsulanInisiatifRepository
	DB                        *sql.DB
}

func NewUsulanInisiatifServiceImpl(usulanInisiatifRepository repository.UsulanInisiatifRepository, DB *sql.DB) *UsulanInisiatifServiceImpl {
	return &UsulanInisiatifServiceImpl{
		UsulanInisiatifRepository: usulanInisiatifRepository,
		DB:                        DB,
	}
}

func (service *UsulanInisiatifServiceImpl) Create(ctx context.Context, request usulan.UsulanInisiatifCreateRequest) (usulan.UsulanInisiatifResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return usulan.UsulanInisiatifResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
	uuId := fmt.Sprintf("USU-INIS-%s", randomDigits)

	domainUsulanInisiatif := domain.UsulanInisiatif{
		Id:        uuId,
		Usulan:    request.Usulan,
		Manfaat:   request.Manfaat,
		Uraian:    request.Uraian,
		Tahun:     request.Tahun,
		RekinId:   request.RekinId,
		PegawaiId: request.PegawaiId,
		KodeOpd:   request.KodeOpd,
		Status:    request.Status,
	}

	usulanInisiatif, err := service.UsulanInisiatifRepository.Create(ctx, tx, domainUsulanInisiatif)
	if err != nil {
		return usulan.UsulanInisiatifResponse{}, err
	}

	response := helper.ToUsulanInisiatifResponse(usulanInisiatif)
	return response, nil
}

func (service *UsulanInisiatifServiceImpl) Update(ctx context.Context, request usulan.UsulanInisiatifUpdateRequest) (usulan.UsulanInisiatifResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return usulan.UsulanInisiatifResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Cari usulan inisiatif berdasarkan ID
	existingUsulan, err := service.UsulanInisiatifRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return usulan.UsulanInisiatifResponse{}, fmt.Errorf("usulan inisiatif dengan ID %s tidak ditemukan", request.Id)
	}

	existingUsulan.Usulan = request.Usulan
	existingUsulan.Manfaat = request.Manfaat
	existingUsulan.Uraian = request.Uraian
	existingUsulan.Tahun = request.Tahun
	existingUsulan.PegawaiId = request.PegawaiId
	existingUsulan.KodeOpd = request.KodeOpd
	existingUsulan.Status = request.Status
	usulanInisiatif, err := service.UsulanInisiatifRepository.Update(ctx, tx, existingUsulan)
	if err != nil {
		return usulan.UsulanInisiatifResponse{}, err
	}

	response := helper.ToUsulanInisiatifResponse(usulanInisiatif)
	return response, nil
}

func (service *UsulanInisiatifServiceImpl) FindById(ctx context.Context, idUsulan string) (usulan.UsulanInisiatifResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return usulan.UsulanInisiatifResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	usulanInisiatif, err := service.UsulanInisiatifRepository.FindById(ctx, tx, idUsulan)
	if err != nil {
		return usulan.UsulanInisiatifResponse{}, fmt.Errorf("usulan inisiatif dengan ID %s tidak ditemukan", idUsulan)
	}

	response := helper.ToUsulanInisiatifResponse(usulanInisiatif)
	return response, nil
}

func (service *UsulanInisiatifServiceImpl) FindAll(ctx context.Context, pegawaiId *string, isActive *bool, rekinId *string) ([]usulan.UsulanInisiatifResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []usulan.UsulanInisiatifResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	usulanInisiatif, err := service.UsulanInisiatifRepository.FindAll(ctx, tx, pegawaiId, isActive, rekinId)
	if err != nil {
		return []usulan.UsulanInisiatifResponse{}, err
	}

	response := helper.ToUsulanInisiatifResponses(usulanInisiatif)
	return response, nil
}

func (service *UsulanInisiatifServiceImpl) Delete(ctx context.Context, idUsulan string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Cari usulan inisiatif berdasarkan ID
	_, err = service.UsulanInisiatifRepository.FindById(ctx, tx, idUsulan)
	if err != nil {
		return fmt.Errorf("usulan inisiatif dengan ID %s tidak ditemukan", idUsulan)
	}

	// Jika usulan inisiatif ditemukan, lanjutkan dengan penghapusan
	err = service.UsulanInisiatifRepository.Delete(ctx, tx, idUsulan)
	if err != nil {
		return fmt.Errorf("gagal menghapus usulan inisiatif: %v", err)
	}

	return nil
}
