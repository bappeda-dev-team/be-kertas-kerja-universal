package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/tujuanopd"
	"ekak_kabupaten_madiun/repository"
	"fmt"
	"strconv"

	"github.com/google/uuid"
)

type TujuanOpdServiceImpl struct {
	TujuanOpdRepository repository.TujuanOpdRepository
	OpdRepository       repository.OpdRepository
	DB                  *sql.DB
}

func NewTujuanOpdServiceImpl(tujuanOpdRepository repository.TujuanOpdRepository, opdRepository repository.OpdRepository, DB *sql.DB) *TujuanOpdServiceImpl {
	return &TujuanOpdServiceImpl{
		TujuanOpdRepository: tujuanOpdRepository,
		OpdRepository:       opdRepository,
		DB:                  DB,
	}
}

func (service *TujuanOpdServiceImpl) Create(ctx context.Context, request tujuanopd.TujuanOpdCreateRequest) (tujuanopd.TujuanOpdResponse, error) {

	tx, err := service.DB.Begin()
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	tujuanOpdDomain := domain.TujuanOpd{
		KodeOpd:          request.KodeOpd,
		Tujuan:           request.Tujuan,
		RumusPerhitungan: request.RumusPerhitungan,
		SumberData:       request.SumberData,
		TahunAwal:        request.TahunAwal,
		TahunAkhir:       request.TahunAkhir,
	}

	// Convert indikator request to domain
	for _, indikatorReq := range request.Indikator {
		// Generate ID indikator dengan format IND-TJN-XXXXX
		uuidInd := uuid.New().String()[:5]
		indikatorId := fmt.Sprintf("IND-TJN-%s", uuidInd)

		indikatorDomain := domain.Indikator{
			Id:        indikatorId,
			Indikator: indikatorReq.Indikator,
		}

		// Convert target request to domain
		for _, targetReq := range indikatorReq.Target {
			// Generate ID target dengan format TRG-TJN-XXXXX
			uuidTrg := uuid.New().String()[:5]
			targetId := fmt.Sprintf("TRG-TJN-%s", uuidTrg)

			targetDomain := domain.Target{
				Id:     targetId,
				Target: targetReq.Target,
				Satuan: targetReq.Satuan,
				Tahun:  targetReq.Tahun,
			}
			indikatorDomain.Target = append(indikatorDomain.Target, targetDomain)
		}

		tujuanOpdDomain.Indikator = append(tujuanOpdDomain.Indikator, indikatorDomain)
	}

	// Panggil repository dan terima hasil domain yang sudah ada ID-nya
	tujuanOpdResult, err := service.TujuanOpdRepository.Create(ctx, tx, tujuanOpdDomain)
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, err
	}

	return helper.ToTujuanOpdResponse(tujuanOpdResult), nil
}

func (service *TujuanOpdServiceImpl) Update(ctx context.Context, request tujuanopd.TujuanOpdUpdateRequest) (tujuanopd.TujuanOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Cek apakah data exists dan ambil data lama
	existingTujuanOpd, err := service.TujuanOpdRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, err
	}

	// Buat map untuk menyimpan ID lama
	existingIndikatorMap := make(map[string]domain.Indikator)
	existingTargetMap := make(map[string]domain.Target)

	// Simpan ID lama ke map
	for _, indikator := range existingTujuanOpd.Indikator {
		existingIndikatorMap[indikator.Id] = indikator
		for _, target := range indikator.Target {
			existingTargetMap[target.Id] = target
		}
	}

	// Update data utama
	tujuanOpd := domain.TujuanOpd{
		Id:               request.Id,
		KodeOpd:          request.KodeOpd,
		Tujuan:           request.Tujuan,
		RumusPerhitungan: request.RumusPerhitungan,
		SumberData:       request.SumberData,
		TahunAwal:        request.TahunAwal,
		TahunAkhir:       request.TahunAkhir,
		Indikator:        []domain.Indikator{},
	}

	// Convert indikator dan target request ke domain
	for _, indikatorReq := range request.Indikator {
		var indikatorId string

		// Gunakan ID lama jika ada, generate baru jika tidak ada
		if indikatorReq.Id != "" {
			indikatorId = indikatorReq.Id
		} else {
			uuidInd := uuid.New().String()[:5]
			indikatorId = fmt.Sprintf("IND-TJN-%s", uuidInd)
		}

		indikatorDomain := domain.Indikator{
			Id:        indikatorId,
			Indikator: indikatorReq.Indikator,
			Target:    []domain.Target{},
		}

		// Convert target
		for _, targetReq := range indikatorReq.Target {
			var targetId string

			// Gunakan ID lama jika ada, generate baru jika tidak ada
			if targetReq.Id != "" {
				targetId = targetReq.Id
			} else {
				uuidTrg := uuid.New().String()[:5]
				targetId = fmt.Sprintf("TRG-TJN-%s", uuidTrg)
			}

			targetDomain := domain.Target{
				Id:          targetId,
				IndikatorId: indikatorId,
				Target:      targetReq.Target,
				Satuan:      targetReq.Satuan,
				Tahun:       targetReq.Tahun,
			}
			indikatorDomain.Target = append(indikatorDomain.Target, targetDomain)
		}

		tujuanOpd.Indikator = append(tujuanOpd.Indikator, indikatorDomain)
	}

	// Update semua data ke database
	err = service.TujuanOpdRepository.Update(ctx, tx, tujuanOpd)
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, err
	}

	return helper.ToTujuanOpdResponse(tujuanOpd), nil
}

func (service *TujuanOpdServiceImpl) Delete(ctx context.Context, tujuanOpdId int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	_, err = service.TujuanOpdRepository.FindById(ctx, tx, tujuanOpdId)
	if err != nil {
		return err
	}

	return service.TujuanOpdRepository.Delete(ctx, tx, tujuanOpdId)
}

func (service *TujuanOpdServiceImpl) FindById(ctx context.Context, tujuanOpdId int) (tujuanopd.TujuanOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data tujuan OPD
	tujuanOpd, err := service.TujuanOpdRepository.FindById(ctx, tx, tujuanOpdId)
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, err
	}

	// Ambil data indikator
	indikators, err := service.TujuanOpdRepository.FindIndikatorByTujuanId(ctx, tx, tujuanOpdId)
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, err
	}

	// Map untuk menyimpan indikator
	indikatorMap := make(map[string]*domain.Indikator)

	// Proses setiap indikator
	for _, ind := range indikators {
		indikatorMap[ind.Id] = &domain.Indikator{
			Id:        ind.Id,
			Indikator: ind.Indikator,
			Target:    []domain.Target{},
		}

		targets, err := service.TujuanOpdRepository.FindTargetByIndikatorId(ctx, tx, ind.Id, "9999") // Gunakan tahun yang besar untuk mendapatkan semua target
		if err != nil {
			return tujuanopd.TujuanOpdResponse{}, err
		}

		// Tambahkan target ke indikator
		indikatorMap[ind.Id].Target = targets // Langsung assign targets yang didapat
	}

	// Convert map ke slice untuk response
	var indikatorResponses []domain.Indikator
	for _, indikator := range indikatorMap {
		indikatorResponses = append(indikatorResponses, *indikator)
	}

	// Set indikator ke tujuan OPD
	tujuanOpd.Indikator = indikatorResponses

	return helper.ToTujuanOpdResponse(tujuanOpd), nil
}

func (service *TujuanOpdServiceImpl) FindAll(ctx context.Context, kodeOpd string, tahun string) ([]tujuanopd.TujuanOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi tahun
	if len(tahun) != 4 {
		return nil, fmt.Errorf("format tahun tidak valid: %s", tahun)
	}
	if _, err := strconv.Atoi(tahun); err != nil {
		return nil, fmt.Errorf("tahun harus berupa angka: %s", tahun)
	}

	tujuanOpds, err := service.TujuanOpdRepository.FindAll(ctx, tx, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}

	opd, err := service.OpdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return nil, err
	}

	for i := range tujuanOpds {
		tujuanOpds[i].NamaOpd = opd.NamaOpd
		indikators, err := service.TujuanOpdRepository.FindIndikatorByTujuanId(ctx, tx, tujuanOpds[i].Id)
		if err != nil {
			return nil, err
		}

		// Untuk setiap indikator, ambil targetnya
		for j := range indikators {
			targets, err := service.TujuanOpdRepository.FindTargetByIndikatorId(ctx, tx, indikators[j].Id, tahun)
			if err != nil {
				return nil, err
			}
			indikators[j].Target = targets
		}

		tujuanOpds[i].Indikator = indikators
	}

	return helper.ToTujuanOpdResponses(tujuanOpds), nil
}
