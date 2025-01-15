package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/rencanakinerja"
	"ekak_kabupaten_madiun/repository"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
)

type ManualIKServiceImpl struct {
	ManualIKRepository repository.ManualIKRepository
	DB                 *sql.DB
	Validate           *validator.Validate
}

func NewManualIKServiceImpl(manualIKRepository repository.ManualIKRepository, db *sql.DB, validate *validator.Validate) *ManualIKServiceImpl {
	return &ManualIKServiceImpl{
		ManualIKRepository: manualIKRepository,
		DB:                 db,
		Validate:           validate,
	}
}

func (service *ManualIKServiceImpl) Create(ctx context.Context, request rencanakinerja.ManualIKCreateRequest, indikatorId string) (rencanakinerja.ManualIKResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return rencanakinerja.ManualIKResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return rencanakinerja.ManualIKResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	request.Id = r.Intn(100000)

	manualIK := domain.ManualIK{
		Id:                  request.Id,
		IndikatorId:         indikatorId,
		Perspektif:          request.Perspektif,
		TujuanRekin:         request.TujuanRekin,
		Definisi:            request.Definisi,
		KeyActivities:       request.KeyActivities,
		Formula:             request.Formula,
		JenisIndikator:      request.JenisIndikator,
		Kinerja:             request.OutputData.Kinerja,
		Penduduk:            request.OutputData.Penduduk,
		Spatial:             request.OutputData.Spatial,
		UnitPenanggungJawab: request.UnitPenanggungJawab,
		UnitPenyediaData:    request.UnitPenyediaData,
		SumberData:          request.SumberData,
		JangkaWaktuAwal:     helper.EmptyStringIfNull(request.JangkaWaktuAwal),
		JangkaWaktuAkhir:    helper.EmptyStringIfNull(request.JangkaWaktuAkhir),
		PeriodePelaporan:    request.PeriodePelaporan,
	}

	manualIK, err = service.ManualIKRepository.Create(ctx, tx, manualIK)
	if err != nil {
		// Check untuk MySQL error
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				return rencanakinerja.ManualIKResponse{}, fmt.Errorf("duplicate entry: indikator ID %s sudah ada", indikatorId)
			}
		}
		return rencanakinerja.ManualIKResponse{}, err
	}

	return helper.ToManualIKResponse(manualIK), nil
}

func (service *ManualIKServiceImpl) Update(ctx context.Context, request rencanakinerja.ManualIKUpdateRequest, indikatorId string) (rencanakinerja.ManualIKResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return rencanakinerja.ManualIKResponse{}, err
	}

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	manualIK := domain.ManualIK{
		IndikatorId:         indikatorId,
		Perspektif:          request.Perspektif,
		TujuanRekin:         request.TujuanRekin,
		Definisi:            request.Definisi,
		KeyActivities:       request.KeyActivities,
		Formula:             request.Formula,
		JenisIndikator:      request.JenisIndikator,
		Kinerja:             request.OutputData.Kinerja,
		Penduduk:            request.OutputData.Penduduk,
		Spatial:             request.OutputData.Spatial,
		UnitPenanggungJawab: request.UnitPenanggungJawab,
		UnitPenyediaData:    request.UnitPenyediaData,
		SumberData:          request.SumberData,
		JangkaWaktuAwal:     request.JangkaWaktuAwal,
		JangkaWaktuAkhir:    request.JangkaWaktuAkhir,
		PeriodePelaporan:    request.PeriodePelaporan,
	}

	manualIK, err = service.ManualIKRepository.Update(ctx, tx, manualIK)
	if err != nil {
		return rencanakinerja.ManualIKResponse{}, err
	}

	return helper.ToManualIKResponse(manualIK), nil
}

// Fungsi untuk mendapatkan data manual IK
// func (service *ManualIKServiceImpl) getManualIKData(ctx context.Context, tx *sql.Tx, indikatorId string) ([]domain.ManualIK, error) {
// 	manualIKs, err := service.ManualIKRepository.GetManualIK(ctx, tx, indikatorId)
// 	if err != nil {
// 		return nil, fmt.Errorf("gagal mengambil data manual IK: %v", err)
// 	}
// 	if len(manualIKs) == 0 {
// 		return nil, fmt.Errorf("manual IK dengan indikator ID %s tidak ditemukan", indikatorId)
// 	}
// 	return manualIKs, nil
// }

// Fungsi untuk mendapatkan data rencana kinerja
func (service *ManualIKServiceImpl) getRencanaKinerjaWithTargetData(ctx context.Context, tx *sql.Tx, indikatorId string) (domain.Indikator, error) {
	indikator, rencanaKinerja, targets, err := service.ManualIKRepository.GetRencanaKinerjaWithTarget(ctx, tx, indikatorId)
	if err != nil {
		return domain.Indikator{}, fmt.Errorf("gagal mengambil data rencana kinerja: %v", err)
	}

	// Set data lengkap
	indikator.RencanaKinerja = rencanaKinerja
	indikator.Target = targets

	return indikator, nil
}

// Fungsi untuk mendapatkan data target
// func (service *ManualIKServiceImpl) getTargetData(ctx context.Context, tx *sql.Tx, indikatorId string) ([]domain.Target, error) {
// 	targets, err := service.ManualIKRepository.GetTargets(ctx, tx, indikatorId)
// 	if err != nil {
// 		return nil, fmt.Errorf("gagal mengambil data target: %v", err)
// 	}
// 	return targets, nil
// }

func (service *ManualIKServiceImpl) FindManualIKByIndikatorId(ctx context.Context, indikatorId string) (rencanakinerja.ManualIKResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return rencanakinerja.ManualIKResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data rencana kinerja terlebih dahulu
	indikator, err := service.getRencanaKinerjaWithTargetData(ctx, tx, indikatorId)
	if err != nil {
		return rencanakinerja.ManualIKResponse{}, err
	}

	// Ambil data manual IK jika ada
	manualIK, err := service.ManualIKRepository.FindByIndikatorId(ctx, tx, indikatorId)
	if err != nil && err != sql.ErrNoRows {
		return rencanakinerja.ManualIKResponse{}, err
	}

	// Set data indikator ke manual IK
	manualIK.DataIndikator = indikator
	manualIK.IndikatorId = indikatorId

	// Convert ke response
	response := helper.ToManualIKResponse(manualIK)
	return response, nil
}
