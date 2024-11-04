package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
	"ekak_kabupaten_madiun/model/web/opdmaster"
	"ekak_kabupaten_madiun/model/web/programkegiatan"
	"ekak_kabupaten_madiun/repository"
	"fmt"

	"github.com/google/uuid"
)

type ProgramServiceImpl struct {
	programRepository repository.ProgramRepository
	opdRepository     repository.OpdRepository
	DB                *sql.DB
}

func NewProgramServiceImpl(programRepository repository.ProgramRepository, opdRepository repository.OpdRepository, DB *sql.DB) *ProgramServiceImpl {
	return &ProgramServiceImpl{
		programRepository: programRepository,
		opdRepository:     opdRepository,
		DB:                DB,
	}
}

func (service *ProgramServiceImpl) Create(ctx context.Context, request programkegiatan.ProgramKegiatanCreateRequest) (programkegiatan.ProgramKegiatanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return programkegiatan.ProgramKegiatanResponse{}, err
	}

	defer helper.CommitOrRollback(tx)
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, request.KodeOPD)
	if err != nil {
		tx.Rollback()
		return programkegiatan.ProgramKegiatanResponse{}, fmt.Errorf("kode OPD tidak valid: %v", err)
	}

	if opd.KodeOpd == "" {
		tx.Rollback()
		return programkegiatan.ProgramKegiatanResponse{}, fmt.Errorf("kode OPD tidak ditemukan")
	}

	uuidPrgm := fmt.Sprintf("PRGM-KGT-%s", uuid.New().String()[:5])

	program := domainmaster.ProgramKegiatan{
		Id:          uuidPrgm,
		KodeProgram: request.KodeProgram,
		NamaProgram: request.NamaProgram,
		KodeOPD:     request.KodeOPD,
		Tahun:       request.Tahun,
		IsActive:    request.IsActive,
	}

	var indikators []domain.Indikator
	for _, indikatorRequest := range request.Indikator {
		uuidIndikator := fmt.Sprintf("IND-KGT-%s", uuid.New().String()[:5])
		indikator := domain.Indikator{
			Id:        uuidIndikator,
			ProgramId: program.Id,
			Indikator: indikatorRequest.Indikator,
			Tahun:     indikatorRequest.Tahun,
		}

		var targets []domain.Target
		for _, targetRequest := range indikatorRequest.Target {
			uuidTarget := fmt.Sprintf("TRGT-KGT-%s", uuid.New().String()[:5])
			target := domain.Target{
				Id:          uuidTarget,
				IndikatorId: indikator.Id,
				Target:      targetRequest.Target,
				Satuan:      targetRequest.Satuan,
				Tahun:       targetRequest.Tahun,
			}
			targets = append(targets, target)
		}
		indikator.Target = targets
		indikators = append(indikators, indikator)
	}
	program.Indikator = indikators

	result, err := service.programRepository.Create(ctx, tx, program)
	if err != nil {
		tx.Rollback()
		return programkegiatan.ProgramKegiatanResponse{}, err
	}
	var indikatorResponses []programkegiatan.IndikatorResponse
	for _, indikator := range result.Indikator {
		var targetResponses []programkegiatan.TargetResponse
		for _, target := range indikator.Target {
			targetResponse := programkegiatan.TargetResponse{
				Id:          target.Id,
				IndikatorId: target.IndikatorId,
				Target:      target.Target,
				Satuan:      target.Satuan,
				Tahun:       target.Tahun,
			}
			targetResponses = append(targetResponses, targetResponse)
		}

		indikatorResponse := programkegiatan.IndikatorResponse{
			Id:        indikator.Id,
			ProgramId: indikator.ProgramId,
			Indikator: indikator.Indikator,
			Tahun:     indikator.Tahun,
			Target:    targetResponses,
		}
		indikatorResponses = append(indikatorResponses, indikatorResponse)
	}

	// if err = tx.Commit(); err != nil {
	// 	tx.Rollback()
	// 	return programkegiatan.ProgramKegiatanResponse{}, err
	// }

	return programkegiatan.ProgramKegiatanResponse{
		Id:          result.Id,
		KodeProgram: result.KodeProgram,
		NamaProgram: result.NamaProgram,
		KodeOPD: opdmaster.OpdResponseForAll{
			KodeOpd: opd.KodeOpd,
			NamaOpd: opd.NamaOpd,
		},
		Tahun:     result.Tahun,
		IsActive:  result.IsActive,
		Indikator: indikatorResponses,
	}, nil
}

func (service *ProgramServiceImpl) Update(ctx context.Context, request programkegiatan.ProgramKegiatanUpdateRequest) (programkegiatan.ProgramKegiatanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return programkegiatan.ProgramKegiatanResponse{}, err
	}
	defer helper.CommitOrRollback(tx)
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, request.KodeOPD)
	if err != nil {
		tx.Rollback()
		return programkegiatan.ProgramKegiatanResponse{}, fmt.Errorf("kode OPD tidak valid: %v", err)
	}

	if opd.KodeOpd == "" {
		tx.Rollback()
		return programkegiatan.ProgramKegiatanResponse{}, fmt.Errorf("kode OPD tidak ditemukan")
	}

	program := domainmaster.ProgramKegiatan{
		Id:          request.Id,
		KodeProgram: request.KodeProgram,
		NamaProgram: request.NamaProgram,
		KodeOPD:     request.KodeOPD,
		Tahun:       request.Tahun,
		IsActive:    request.IsActive,
	}

	var indikators []domain.Indikator
	for _, indikator := range request.Indikator {
		indikatorId := indikator.Id
		if indikatorId == "" {
			indikatorId = fmt.Sprintf("IND-KGT-%s", uuid.New().String()[:5])
		}

		var targets []domain.Target
		for _, target := range indikator.Target {
			targetId := target.Id
			if targetId == "" {
				targetId = fmt.Sprintf("TRGT-KGT-%s", uuid.New().String()[:5])
			}

			targetDomain := domain.Target{
				Id:          targetId,
				IndikatorId: indikatorId,
				Target:      target.Target,
				Satuan:      target.Satuan,
				Tahun:       target.Tahun,
			}
			targets = append(targets, targetDomain)
		}

		indikatorDomain := domain.Indikator{
			Id:        indikatorId,
			ProgramId: request.Id,
			Indikator: indikator.Indikator,
			Tahun:     indikator.Tahun,
			Target:    targets,
		}
		indikators = append(indikators, indikatorDomain)
	}
	program.Indikator = indikators

	result, err := service.programRepository.Update(ctx, tx, program)
	if err != nil {
		return programkegiatan.ProgramKegiatanResponse{}, err
	}

	var indikatorResponses []programkegiatan.IndikatorResponse
	for _, indikator := range result.Indikator {
		var targetResponses []programkegiatan.TargetResponse
		for _, target := range indikator.Target {
			targetResponse := programkegiatan.TargetResponse{
				Id:          target.Id,
				IndikatorId: target.IndikatorId,
				Target:      target.Target,
				Satuan:      target.Satuan,
				Tahun:       target.Tahun,
			}
			targetResponses = append(targetResponses, targetResponse)
		}

		indikatorResponse := programkegiatan.IndikatorResponse{
			Id:        indikator.Id,
			ProgramId: indikator.ProgramId,
			Indikator: indikator.Indikator,
			Tahun:     indikator.Tahun,
			Target:    targetResponses,
		}
		indikatorResponses = append(indikatorResponses, indikatorResponse)
	}

	return programkegiatan.ProgramKegiatanResponse{
		Id:          result.Id,
		KodeProgram: result.KodeProgram,
		NamaProgram: result.NamaProgram,
		KodeOPD: opdmaster.OpdResponseForAll{
			KodeOpd: opd.KodeOpd,
			NamaOpd: opd.NamaOpd,
		},
		Tahun:     result.Tahun,
		IsActive:  result.IsActive,
		Indikator: indikatorResponses,
	}, nil
}

func (service *ProgramServiceImpl) Delete(ctx context.Context, programId string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	err = service.programRepository.Delete(ctx, tx, programId)
	if err != nil {
		return err
	}

	return nil
}

func (service *ProgramServiceImpl) FindById(ctx context.Context, programId string) (programkegiatan.ProgramKegiatanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return programkegiatan.ProgramKegiatanResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Mengambil data program
	result, err := service.programRepository.FindById(ctx, tx, programId)
	if err != nil {
		return programkegiatan.ProgramKegiatanResponse{}, fmt.Errorf("gagal mengambil data program: %v", err)
	}

	// Mengambil data OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, result.KodeOPD)
	if err != nil {
		return programkegiatan.ProgramKegiatanResponse{}, fmt.Errorf("gagal mengambil data OPD: %v", err)
	}

	// Mengambil semua indikator untuk program ini
	indikators, err := service.programRepository.FindIndikatorByProgramId(ctx, tx, result.Id)
	if err != nil {
		return programkegiatan.ProgramKegiatanResponse{}, fmt.Errorf("gagal mengambil data indikator: %v", err)
	}

	var indikatorResponses []programkegiatan.IndikatorResponse
	for _, indikator := range indikators {
		// Mengambil semua target untuk setiap indikator
		targets, err := service.programRepository.FindTargetByIndikatorId(ctx, tx, indikator.Id)
		if err != nil {
			return programkegiatan.ProgramKegiatanResponse{}, fmt.Errorf("gagal mengambil data target: %v", err)
		}

		var targetResponses []programkegiatan.TargetResponse
		for _, target := range targets {
			targetResponse := programkegiatan.TargetResponse{
				Id:          target.Id,
				IndikatorId: target.IndikatorId,
				Target:      target.Target,
				Satuan:      target.Satuan,
				Tahun:       target.Tahun,
			}
			targetResponses = append(targetResponses, targetResponse)
		}

		indikatorResponse := programkegiatan.IndikatorResponse{
			Id:        indikator.Id,
			ProgramId: indikator.ProgramId,
			Indikator: indikator.Indikator,
			Tahun:     indikator.Tahun,
			Target:    targetResponses,
		}
		indikatorResponses = append(indikatorResponses, indikatorResponse)
	}

	return programkegiatan.ProgramKegiatanResponse{
		Id:          result.Id,
		KodeProgram: result.KodeProgram,
		NamaProgram: result.NamaProgram,
		KodeOPD: opdmaster.OpdResponseForAll{
			KodeOpd: opd.KodeOpd,
			NamaOpd: opd.NamaOpd,
		},
		Tahun:     result.Tahun,
		IsActive:  result.IsActive,
		Indikator: indikatorResponses,
	}, nil
}

func (service *ProgramServiceImpl) FindAll(ctx context.Context) ([]programkegiatan.ProgramKegiatanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Mengambil semua program
	results, err := service.programRepository.FindAll(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data program: %v", err)
	}

	var programResponses []programkegiatan.ProgramKegiatanResponse

	for _, program := range results {
		// Mengambil data OPD dari cache atau database
		opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, program.KodeOPD)
		if err != nil {
			return nil, fmt.Errorf("gagal mengambil data OPD: %v", err)
		}

		// Mengambil semua indikator untuk program ini
		indikators, err := service.programRepository.FindIndikatorByProgramId(ctx, tx, program.Id)
		if err != nil {
			return nil, fmt.Errorf("gagal mengambil data indikator untuk program %s: %v", program.Id, err)
		}

		var indikatorResponses []programkegiatan.IndikatorResponse
		for _, indikator := range indikators {
			// Mengambil semua target untuk setiap indikator
			targets, err := service.programRepository.FindTargetByIndikatorId(ctx, tx, indikator.Id)
			if err != nil {
				return nil, fmt.Errorf("gagal mengambil data target untuk indikator %s: %v", indikator.Id, err)
			}

			// Membuat response untuk semua target
			var targetResponses []programkegiatan.TargetResponse
			for _, target := range targets {
				targetResponse := programkegiatan.TargetResponse{
					Id:          target.Id,
					IndikatorId: target.IndikatorId,
					Target:      target.Target,
					Satuan:      target.Satuan,
					Tahun:       target.Tahun,
				}
				targetResponses = append(targetResponses, targetResponse)
			}

			// Membuat response untuk indikator dengan semua targetnya
			indikatorResponse := programkegiatan.IndikatorResponse{
				Id:        indikator.Id,
				ProgramId: indikator.ProgramId,
				Indikator: indikator.Indikator,
				Tahun:     indikator.Tahun,
				Target:    targetResponses,
			}
			indikatorResponses = append(indikatorResponses, indikatorResponse)
		}

		// Membuat response untuk program dengan indikator dan targetnya
		programResponse := programkegiatan.ProgramKegiatanResponse{
			Id:          program.Id,
			KodeProgram: program.KodeProgram,
			NamaProgram: program.NamaProgram,
			KodeOPD: opdmaster.OpdResponseForAll{
				KodeOpd: opd.KodeOpd,
				NamaOpd: opd.NamaOpd,
			},
			Tahun:     program.Tahun,
			IsActive:  program.IsActive,
			Indikator: indikatorResponses,
		}
		programResponses = append(programResponses, programResponse)
	}

	return programResponses, nil
}
