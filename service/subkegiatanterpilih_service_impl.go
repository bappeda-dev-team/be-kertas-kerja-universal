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
	"log"
)

type SubKegiatanTerpilihServiceImpl struct {
	RencanaKinerjaRepository      repository.RencanaKinerjaRepository
	SubKegiatanRepository         repository.SubKegiatanRepository
	SubKegiatanTerpilihRepository repository.SubKegiatanTerpilihRepository
	DB                            *sql.DB
}

func NewSubKegiatanTerpilihServiceImpl(rencanaKinerjaRepository repository.RencanaKinerjaRepository, subKegiatanRepository repository.SubKegiatanRepository, subKegiatanTerpilihRepository repository.SubKegiatanTerpilihRepository, DB *sql.DB) *SubKegiatanTerpilihServiceImpl {
	return &SubKegiatanTerpilihServiceImpl{
		RencanaKinerjaRepository:      rencanaKinerjaRepository,
		SubKegiatanRepository:         subKegiatanRepository,
		SubKegiatanTerpilihRepository: subKegiatanTerpilihRepository,
		DB:                            DB,
	}
}

func (service *SubKegiatanTerpilihServiceImpl) Update(ctx context.Context, request subkegiatan.SubKegiatanTerpilihUpdateRequest) (subkegiatan.SubKegiatanTerpilihResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return subkegiatan.SubKegiatanTerpilihResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	var rencanaKinerja domain.RencanaKinerja
	if request.Id != "" {
		rencanaKinerja, err = service.RencanaKinerjaRepository.FindById(ctx, tx, request.Id, "", "")
		if err != nil {
			log.Printf("Gagal menemukan RencanaKinerja: %v", err)
			return subkegiatan.SubKegiatanTerpilihResponse{}, fmt.Errorf("gagal menemukan RencanaKinerja: %v", err)
		}
	} else {
		return subkegiatan.SubKegiatanTerpilihResponse{}, errors.New("id rencana kinerja tidak boleh kosong")
	}

	// Cek apakah data dengan kode_subkegiatan tersebut ada
	_, err = service.SubKegiatanRepository.FindByKodeSubKegiatan(ctx, tx, request.KodeSubKegiatan)
	if err != nil {
		return subkegiatan.SubKegiatanTerpilihResponse{}, errors.New("subkegiatan tidak ditemukan")
	}

	subKegiatanTerpilih := domain.SubKegiatanTerpilih{
		Id:              rencanaKinerja.Id,
		KodeSubKegiatan: request.KodeSubKegiatan,
	}

	result, err := service.SubKegiatanTerpilihRepository.Update(ctx, tx, subKegiatanTerpilih)
	if err != nil {
		return subkegiatan.SubKegiatanTerpilihResponse{}, err
	}

	return subkegiatan.SubKegiatanTerpilihResponse{
		KodeSubKegiatan: subkegiatan.SubKegiatanResponse{
			KodeSubKegiatan: result.KodeSubKegiatan,
		},
	}, nil
}

func (service *SubKegiatanTerpilihServiceImpl) FindByKodeSubKegiatan(ctx context.Context, kodeSubKegiatan string) (subkegiatan.SubKegiatanTerpilihResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return subkegiatan.SubKegiatanTerpilihResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Cek apakah data dengan kode_subkegiatan tersebut ada

	result, err := service.SubKegiatanRepository.FindByKodeSubKegiatan(ctx, tx, kodeSubKegiatan)
	if err != nil {
		return subkegiatan.SubKegiatanTerpilihResponse{}, errors.New("subkegiatan tidak ditemukan")
	}

	return subkegiatan.SubKegiatanTerpilihResponse{
		KodeSubKegiatan: subkegiatan.SubKegiatanResponse{
			KodeSubKegiatan: result.KodeSubKegiatan,
			NamaSubKegiatan: result.NamaSubKegiatan,
		},
	}, nil
}

func (service *SubKegiatanTerpilihServiceImpl) Delete(ctx context.Context, id string, kodeSubKegiatan string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi: Cek apakah data dengan id dan kodeSubKegiatan ada
	_, err = service.SubKegiatanTerpilihRepository.FindByIdAndKodeSubKegiatan(ctx, tx, id, kodeSubKegiatan)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("data tidak ditemukan")
		}
		return err
	}

	// Lanjutkan dengan penghapusan jika data ditemukan
	err = service.SubKegiatanTerpilihRepository.Delete(ctx, tx, id, kodeSubKegiatan)
	if err != nil {
		return err
	}

	return nil
}
