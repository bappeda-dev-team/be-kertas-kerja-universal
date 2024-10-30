package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"ekak_kabupaten_madiun/repository"
	"errors"
	"fmt"
	"strconv"
)

type PohonKinerjaOpdServiceImpl struct {
	pohonKinerjaOpdRepository repository.PohonKinerjaOpdRepository
	DB                        *sql.DB
}

func NewPohonKinerjaOpdServiceImpl(pohonKinerjaOpdRepository repository.PohonKinerjaOpdRepository, DB *sql.DB) *PohonKinerjaOpdServiceImpl {
	return &PohonKinerjaOpdServiceImpl{
		pohonKinerjaOpdRepository: pohonKinerjaOpdRepository,
		DB:                        DB,
	}
}

func (service *PohonKinerjaOpdServiceImpl) Create(ctx context.Context, request pohonkinerja.PohonKinerjaCreateRequest) (pohonkinerja.PohonKinerjaOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi request
	if request.NamaPohon == "" {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("nama program tidak boleh kosong")
	}

	pohonKinerja := domain.PohonKinerja{
		NamaPohon:  request.NamaPohon,
		Parent:     request.Parent,
		JenisPohon: request.JenisPohon,
		LevelPohon: request.LevelPohon,
		KodeOpd:    request.KodeOpd,
		Keterangan: request.Keterangan,
		Tahun:      request.Tahun,
	}

	result, err := service.pohonKinerjaOpdRepository.Create(ctx, tx, pohonKinerja)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, err
	}

	// Tambahkan log untuk debug tipe data dan nilai
	fmt.Printf("ID Type: %T, Value: %v\n", result.Id, result.Id)

	response := pohonkinerja.PohonKinerjaOpdResponse{
		Id:         result.Id,
		Parent:     strconv.Itoa(result.Parent),
		NamaPohon:  result.NamaPohon,
		JenisPohon: result.JenisPohon,
		LevelPohon: result.LevelPohon,
		KodeOpd:    result.KodeOpd,
		Keterangan: result.Keterangan,
		Tahun:      result.Tahun,
	}

	// Log response untuk memastikan
	fmt.Printf("Response: %+v\n", response)

	return response, nil
}

func (service *PohonKinerjaOpdServiceImpl) Update(ctx context.Context, request pohonkinerja.PohonKinerjaUpdateRequest) (pohonkinerja.PohonKinerjaOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	pohonKinerja := domain.PohonKinerja{
		Id:         request.Id,
		NamaPohon:  request.NamaPohon,
		JenisPohon: request.JenisPohon,
		LevelPohon: request.LevelPohon,
		KodeOpd:    request.KodeOpd,
		Keterangan: request.Keterangan,
		Tahun:      request.Tahun,
	}

	pokin, err := service.pohonKinerjaOpdRepository.Update(ctx, tx, pohonKinerja)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, err
	}

	return pohonkinerja.PohonKinerjaOpdResponse{
		Id:         pokin.Id,
		NamaPohon:  pokin.NamaPohon,
		JenisPohon: pokin.JenisPohon,
		LevelPohon: pokin.LevelPohon,
		KodeOpd:    pokin.KodeOpd,
		Keterangan: pokin.Keterangan,
		Tahun:      pokin.Tahun,
	}, nil
}

func (service *PohonKinerjaOpdServiceImpl) Delete(ctx context.Context, id string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	err = service.pohonKinerjaOpdRepository.Delete(ctx, tx, id)
	if err != nil {
		return err
	}

	return nil
}

func (service *PohonKinerjaOpdServiceImpl) FindById(ctx context.Context, id int) (pohonkinerja.PohonKinerjaOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	pokin, err := service.pohonKinerjaOpdRepository.FindById(ctx, tx, id)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, err
	}

	// Tambahkan validasi jika data tidak ditemukan
	if pokin.Id == 0 {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("data tidak ditemukan")
	}

	return pohonkinerja.PohonKinerjaOpdResponse{
		Id:         pokin.Id,
		Parent:     strconv.Itoa(pokin.Parent),
		NamaPohon:  pokin.NamaPohon,
		JenisPohon: pokin.JenisPohon,
		LevelPohon: pokin.LevelPohon,
		KodeOpd:    pokin.KodeOpd,
		Keterangan: pokin.Keterangan,
		Tahun:      pokin.Tahun,
	}, nil
}

func (service *PohonKinerjaOpdServiceImpl) FindAll(ctx context.Context, kodeOpd, tahun string) (pohonkinerja.PohonKinerjaOpdAllResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdAllResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil semua pohon kinerja
	pokins, err := service.pohonKinerjaOpdRepository.FindAll(ctx, tx, kodeOpd, tahun)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdAllResponse{}, err
	}

	// Validasi jika data tidak ditemukan
	if len(pokins) == 0 {
		return pohonkinerja.PohonKinerjaOpdAllResponse{}, errors.New("data tidak ditemukan untuk kode OPD dan tahun yang diberikan")
	}

	// Buat response sesuai format yang diinginkan
	response := pohonkinerja.PohonKinerjaOpdAllResponse{
		Tahun:         tahun,
		KodeOpd:       kodeOpd,
		PohonKinerjas: []pohonkinerja.PohonKinerjaOpdResponse{},
	}

	// Urutkan berdasarkan level dan parent
	var strategics, tacticals, operationals []pohonkinerja.PohonKinerjaOpdResponse

	// Kelompokkan berdasarkan level
	for _, pokin := range pokins {
		parent := ""
		if pokin.Parent != 0 {
			parent = strconv.Itoa(pokin.Parent)
		}

		detail := pohonkinerja.PohonKinerjaOpdResponse{
			Id:         pokin.Id,
			Parent:     parent,
			JenisPohon: getJenisPohon(pokin.LevelPohon),
			LevelPohon: pokin.LevelPohon,
			NamaPohon:  pokin.NamaPohon,
		}

		switch pokin.LevelPohon {
		case 4:
			strategics = append(strategics, detail)
		case 5:
			tacticals = append(tacticals, detail)
		case 6:
			operationals = append(operationals, detail)
		}
	}

	// Susun ulang berdasarkan hierarki
	response.PohonKinerjas = []pohonkinerja.PohonKinerjaOpdResponse{}

	// Proses setiap Strategic
	for _, strategic := range strategics {
		response.PohonKinerjas = append(response.PohonKinerjas, strategic)

		// Cari Tactical yang terkait
		for _, tactical := range tacticals {
			if tactical.Parent == strconv.Itoa(strategic.Id) {
				response.PohonKinerjas = append(response.PohonKinerjas, tactical)

				// Cari Operational yang terkait dengan Tactical ini
				for _, operational := range operationals {
					if operational.Parent == strconv.Itoa(tactical.Id) {
						response.PohonKinerjas = append(response.PohonKinerjas, operational)
					}
				}
			}
		}
	}

	return response, nil
}

// Helper function untuk menentukan jenis pohon berdasarkan level
func getJenisPohon(level int) string {
	switch level {
	case 4:
		return "Strategic"
	case 5:
		return "Tactical"
	case 6:
		return "Operational"
	default:
		return ""
	}
}
