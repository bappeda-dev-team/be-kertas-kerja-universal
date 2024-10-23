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

type UsulanMusrebangServiceImpl struct {
	usulanMusrebangRepository repository.UsulanMusrebangRepository
	DB                        *sql.DB
}

func NewUsulanMusrebangServiceImpl(usulanMusrebangRepository repository.UsulanMusrebangRepository, DB *sql.DB) *UsulanMusrebangServiceImpl {
	return &UsulanMusrebangServiceImpl{
		usulanMusrebangRepository: usulanMusrebangRepository,
		DB:                        DB,
	}
}

func (service *UsulanMusrebangServiceImpl) Create(ctx context.Context, request usulan.UsulanMusrebangCreateRequest) (usulan.UsulanMusrebangResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return usulan.UsulanMusrebangResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
	uuId := fmt.Sprintf("USU-MUS-%s", randomDigits)

	// Konversi request ke domain.UsulanMusrebang
	domainUsulanMusrebang := domain.UsulanMusrebang{
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

	usulanMusrebang, err := service.usulanMusrebangRepository.Create(ctx, tx, domainUsulanMusrebang)
	if err != nil {
		return usulan.UsulanMusrebangResponse{}, err
	}

	return helper.ToUsulanMusrebangResponse(usulanMusrebang), nil
}

func (service *UsulanMusrebangServiceImpl) Update(ctx context.Context, request usulan.UsulanMusrebangUpdateRequest) (usulan.UsulanMusrebangResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return usulan.UsulanMusrebangResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Cek apakah usulan dengan ID yang diberikan ada
	existingUsulan, err := service.usulanMusrebangRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return usulan.UsulanMusrebangResponse{}, fmt.Errorf("usulan musrebang tidak ditemukan: %v", err)
	}

	// Update data usulan
	existingUsulan.Usulan = request.Usulan
	existingUsulan.Alamat = request.Alamat
	existingUsulan.Uraian = request.Uraian
	existingUsulan.Tahun = request.Tahun
	existingUsulan.PegawaiId = request.PegawaiId
	existingUsulan.KodeOpd = request.KodeOpd
	existingUsulan.Status = request.Status

	updatedUsulan, err := service.usulanMusrebangRepository.Update(ctx, tx, existingUsulan)
	if err != nil {
		return usulan.UsulanMusrebangResponse{}, err
	}

	// Menggunakan helper untuk menghasilkan respons
	return helper.ToUsulanMusrebangResponse(updatedUsulan), nil
}

func (service *UsulanMusrebangServiceImpl) FindById(ctx context.Context, idUsulan string) (usulan.UsulanMusrebangResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return usulan.UsulanMusrebangResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	usulanMusrebang, err := service.usulanMusrebangRepository.FindById(ctx, tx, idUsulan)
	if err != nil {
		return usulan.UsulanMusrebangResponse{}, err
	}

	return helper.ToUsulanMusrebangResponse(usulanMusrebang), nil
}

func (service *UsulanMusrebangServiceImpl) FindAll(ctx context.Context, pegawaiId *string, is_active *bool, rekinId *string) ([]usulan.UsulanMusrebangResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []usulan.UsulanMusrebangResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	usulanMusrebang, err := service.usulanMusrebangRepository.FindAll(ctx, tx, pegawaiId, is_active, rekinId)
	if err != nil {
		return []usulan.UsulanMusrebangResponse{}, err
	}

	usulanMusrebangResponses := helper.ToUsulanMusrebangResponses(usulanMusrebang)

	return usulanMusrebangResponses, nil
}

func (service *UsulanMusrebangServiceImpl) Delete(ctx context.Context, idUsulan string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Cari usulan berdasarkan ID
	_, err = service.usulanMusrebangRepository.FindById(ctx, tx, idUsulan)
	if err != nil {
		return fmt.Errorf("usulan tidak ditemukan: %v", err)
	}

	// Jika usulan ditemukan, lanjutkan dengan penghapusan
	err = service.usulanMusrebangRepository.Delete(ctx, tx, idUsulan)
	if err != nil {
		return fmt.Errorf("gagal menghapus usulan: %v", err)
	}

	return nil
}
