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
	"sort"
	"strconv"
)

type PohonKinerjaOpdServiceImpl struct {
	pohonKinerjaOpdRepository repository.PohonKinerjaRepository
	opdRepository             repository.OpdRepository
	DB                        *sql.DB
}

func NewPohonKinerjaOpdServiceImpl(pohonKinerjaOpdRepository repository.PohonKinerjaRepository, opdRepository repository.OpdRepository, DB *sql.DB) *PohonKinerjaOpdServiceImpl {
	return &PohonKinerjaOpdServiceImpl{
		pohonKinerjaOpdRepository: pohonKinerjaOpdRepository,
		opdRepository:             opdRepository,
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

	// Validasi kode OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, request.KodeOpd)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("kode opd tidak ditemukan")
	}
	if opd.KodeOpd == "" {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("kode opd tidak valid")
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

	// Validasi request
	if request.NamaPohon == "" {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("nama program tidak boleh kosong")
	}

	// Validasi kode OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, request.KodeOpd)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("kode opd tidak ditemukan")
	}
	if opd.KodeOpd == "" {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("kode opd tidak valid")
	}

	// Validasi data yang akan diupdate
	_, err = service.pohonKinerjaOpdRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("data pohon kinerja tidak ditemukan")
	}

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
		Parent:     strconv.Itoa(pokin.Parent),
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

	// Konversi id string ke int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("id tidak valid")
	}

	// Cek apakah data exists
	_, err = service.pohonKinerjaOpdRepository.FindById(ctx, tx, idInt)
	if err != nil {
		return errors.New("data pohon kinerja tidak ditemukan")
	}

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

	// Ambil data OPD berdasarkan kode OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, pokin.KodeOpd)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("data opd tidak ditemukan")
	}

	return pohonkinerja.PohonKinerjaOpdResponse{
		Id:         pokin.Id,
		Parent:     strconv.Itoa(pokin.Parent),
		NamaPohon:  pokin.NamaPohon,
		JenisPohon: pokin.JenisPohon,
		LevelPohon: pokin.LevelPohon,
		KodeOpd:    pokin.KodeOpd,
		NamaOpd:    opd.NamaOpd,
		Keterangan: pokin.Keterangan,
		Tahun:      pokin.Tahun,
	}, nil
}

// func (service *PohonKinerjaOpdServiceImpl) FindAll(ctx context.Context, kodeOpd, tahun string) (pohonkinerja.PohonKinerjaOpdAllResponse, error) {
// 	tx, err := service.DB.Begin()
// 	if err != nil {
// 		return pohonkinerja.PohonKinerjaOpdAllResponse{}, err
// 	}
// 	defer helper.CommitOrRollback(tx)

// 	// Validasi kode OPD
// 	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
// 	if err != nil {
// 		return pohonkinerja.PohonKinerjaOpdAllResponse{}, errors.New("kode opd tidak ditemukan")
// 	}
// 	if opd.KodeOpd == "" {
// 		return pohonkinerja.PohonKinerjaOpdAllResponse{}, errors.New("kode opd tidak valid")
// 	}

// 	// Ambil semua pohon kinerja
// 	pokins, err := service.pohonKinerjaOpdRepository.FindAll(ctx, tx, kodeOpd, tahun)
// 	if err != nil {
// 		return pohonkinerja.PohonKinerjaOpdAllResponse{}, err
// 	}

// 	// Validasi jika data tidak ditemukan
// 	if len(pokins) == 0 {
// 		return pohonkinerja.PohonKinerjaOpdAllResponse{}, errors.New("data tidak ditemukan untuk kode OPD dan tahun yang diberikan")
// 	}

// 	// Buat response sesuai format yang diinginkan
// 	response := pohonkinerja.PohonKinerjaOpdAllResponse{
// 		Tahun:         tahun,
// 		KodeOpd:       kodeOpd,
// 		NamaOpd:       opd.NamaOpd,
// 		PohonKinerjas: []pohonkinerja.PohonKinerjaOpdResponse{},
// 	}

// 	// Urutkan berdasarkan level dan parent
// 	var strategics, tacticals, operationals []pohonkinerja.PohonKinerjaOpdResponse

// 	// Kelompokkan berdasarkan level
// 	for _, pokin := range pokins {
// 		parent := ""
// 		if pokin.Parent != 0 {
// 			parent = strconv.Itoa(pokin.Parent)
// 		}

// 		detail := pohonkinerja.PohonKinerjaOpdResponse{
// 			Id:         pokin.Id,
// 			Parent:     parent,
// 			JenisPohon: helper.GetJenisPohon(pokin.LevelPohon),
// 			LevelPohon: pokin.LevelPohon,
// 			NamaPohon:  pokin.NamaPohon,
// 		}

// 		switch pokin.LevelPohon {
// 		case 4:
// 			strategics = append(strategics, detail)
// 		case 5:
// 			tacticals = append(tacticals, detail)
// 		case 6:
// 			operationals = append(operationals, detail)
// 		}
// 	}

// 	// Susun ulang berdasarkan hierarki
// 	response.PohonKinerjas = []pohonkinerja.PohonKinerjaOpdResponse{}

// 	// Proses setiap Strategic
// 	for _, strategic := range strategics {
// 		response.PohonKinerjas = append(response.PohonKinerjas, strategic)

// 		// Cari Tactical yang terkait
// 		for _, tactical := range tacticals {
// 			if tactical.Parent == strconv.Itoa(strategic.Id) {
// 				response.PohonKinerjas = append(response.PohonKinerjas, tactical)

// 				// Cari Operational yang terkait dengan Tactical ini
// 				for _, operational := range operationals {
// 					if operational.Parent == strconv.Itoa(tactical.Id) {
// 						response.PohonKinerjas = append(response.PohonKinerjas, operational)
// 					}
// 				}
// 			}
// 		}
// 	}

// 	return response, nil
// }

func (service *PohonKinerjaOpdServiceImpl) FindAll(ctx context.Context, kodeOpd, tahun string) (pohonkinerja.PohonKinerjaOpdAllResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdAllResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdAllResponse{}, errors.New("kode opd tidak ditemukan")
	}

	pokins, err := service.pohonKinerjaOpdRepository.FindAll(ctx, tx, kodeOpd, tahun)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdAllResponse{}, err
	}

	// Urutkan data berdasarkan LevelPohon dan Id
	sort.Slice(pokins, func(i, j int) bool {
		if pokins[i].LevelPohon == pokins[j].LevelPohon {
			return pokins[i].Id < pokins[j].Id
		}
		return pokins[i].LevelPohon < pokins[j].LevelPohon
	})

	pohonMap := make(map[int]map[int][]domain.PohonKinerja)
	for i := 0; i <= 6; i++ {
		pohonMap[i] = make(map[int][]domain.PohonKinerja)
	}

	for i := range pokins {
		pokins[i].NamaOpd = opd.NamaOpd
		pohonMap[pokins[i].LevelPohon][pokins[i].Parent] = append(
			pohonMap[pokins[i].LevelPohon][pokins[i].Parent],
			pokins[i],
		)
	}

	var allStrategics []domain.PohonKinerja
	for _, strategics := range pohonMap[4] {
		allStrategics = append(allStrategics, strategics...)
	}

	// Urutkan strategics berdasarkan Id
	sort.Slice(allStrategics, func(i, j int) bool {
		return allStrategics[i].Id < allStrategics[j].Id
	})

	strategics := helper.BuildStrategicOpdResponses(pohonMap, allStrategics)

	return pohonkinerja.PohonKinerjaOpdAllResponse{
		KodeOpd:    kodeOpd,
		NamaOpd:    opd.NamaOpd,
		Tahun:      tahun,
		Strategics: strategics,
	}, nil
}
