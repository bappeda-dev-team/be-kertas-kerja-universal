package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/tujuanpemda"
	"ekak_kabupaten_madiun/repository"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type TujuanPemdaServiceImpl struct {
	TujuanPemdaRepository  repository.TujuanPemdaRepository
	PeriodeRepository      repository.PeriodeRepository
	PohonKinerjaRepository repository.PohonKinerjaRepository
	DB                     *sql.DB
}

func NewTujuanPemdaServiceImpl(tujuanPemdaRepository repository.TujuanPemdaRepository, periodeRepository repository.PeriodeRepository, pohonKinerjaRepository repository.PohonKinerjaRepository, DB *sql.DB) *TujuanPemdaServiceImpl {
	return &TujuanPemdaServiceImpl{
		TujuanPemdaRepository:  tujuanPemdaRepository,
		PeriodeRepository:      periodeRepository,
		PohonKinerjaRepository: pohonKinerjaRepository,
		DB:                     DB,
	}
}

func (service *TujuanPemdaServiceImpl) generateRandomId(ctx context.Context, tx *sql.Tx) int {
	rand.Seed(time.Now().UnixNano())
	for {
		// Generate random number between 10000-99999
		id := rand.Intn(90000) + 10000
		if !service.TujuanPemdaRepository.IsIdExists(ctx, tx, id) {
			return id
		}
	}
}

func generateIndikatorId() string {
	currentYear := time.Now().Format("2006")
	uuid := uuid.New().String()[:5] // Mengambil 5 karakter pertama dari UUID
	return fmt.Sprintf("IND-TJN-PMD-%s-%s", currentYear, uuid)
}

func generateTargetId() string {
	currentYear := time.Now().Format("2006")
	uuid := uuid.New().String()[:5]
	return fmt.Sprintf("TRG-TJN-PMD-%s-%s", currentYear, uuid)
}

func (service *TujuanPemdaServiceImpl) Create(ctx context.Context, request tujuanpemda.TujuanPemdaCreateRequest) (tujuanpemda.TujuanPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi periode
	periode, err := service.PeriodeRepository.FindById(ctx, tx, request.PeriodeId)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf("periode tidak ditemukan: %v", err)
	}

	//validasi pohon kinerja = tematik
	err = service.PohonKinerjaRepository.ValidatePokinLevel(ctx, tx, request.TujuanPemdaId, 0, "tujuan pemda")
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}

	TujuanPemdaId, err := service.PohonKinerjaRepository.FindById(ctx, tx, request.TujuanPemdaId)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}

	// Validasi tahun target untuk setiap indikator
	tahunAwal, _ := strconv.Atoi(periode.TahunAwal)
	tahunAkhir, _ := strconv.Atoi(periode.TahunAkhir)

	for _, indikatorReq := range request.Indikator {
		tahunMap := make(map[string]bool)
		for _, targetReq := range indikatorReq.Target {
			targetTahun, _ := strconv.Atoi(targetReq.Tahun)

			// Validasi rentang tahun
			if targetTahun < tahunAwal || targetTahun > tahunAkhir {
				return tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf(
					"tahun target %d harus berada dalam rentang periode %d-%d",
					targetTahun, tahunAwal, tahunAkhir,
				)
			}

			// Validasi duplikasi tahun
			if tahunMap[targetReq.Tahun] {
				return tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf(
					"duplikasi tahun %s pada indikator %s",
					targetReq.Tahun, indikatorReq.Indikator,
				)
			}
			tahunMap[targetReq.Tahun] = true
		}
	}

	tujuanPemda := domain.TujuanPemda{
		Id:               service.generateRandomId(ctx, tx),
		TujuanPemdaId:    request.TujuanPemdaId,
		RumusPerhitungan: request.RumusPerhitungan,
		SumberData:       request.SumberData,
		PeriodeId:        request.PeriodeId,
	}

	tujuanPemda, err = service.TujuanPemdaRepository.Create(ctx, tx, tujuanPemda)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}

	var indikatorResponses []tujuanpemda.IndikatorResponse
	for _, indikatorRequest := range request.Indikator {
		indikator := domain.Indikator{
			Id:            generateIndikatorId(),
			TujuanPemdaId: tujuanPemda.Id,
			Indikator:     indikatorRequest.Indikator,
		}

		indikator, err = service.TujuanPemdaRepository.CreateIndikator(ctx, tx, indikator)
		if err != nil {
			return tujuanpemda.TujuanPemdaResponse{}, err
		}

		var targetResponses []tujuanpemda.TargetResponse
		for _, targetRequest := range indikatorRequest.Target {
			target := domain.Target{
				Id:          generateTargetId(),
				IndikatorId: indikator.Id,
				Target:      targetRequest.Target,
				Satuan:      targetRequest.Satuan,
				Tahun:       targetRequest.Tahun,
			}

			err = service.TujuanPemdaRepository.CreateTarget(ctx, tx, target)
			if err != nil {
				return tujuanpemda.TujuanPemdaResponse{}, err
			}

			targetResponses = append(targetResponses, tujuanpemda.TargetResponse{
				Id:     target.Id,
				Target: target.Target,
				Satuan: target.Satuan,
				Tahun:  target.Tahun,
			})
		}

		indikatorResponses = append(indikatorResponses, tujuanpemda.IndikatorResponse{
			Id:        indikator.Id,
			Indikator: indikator.Indikator,
			Target:    targetResponses,
		})
	}

	return tujuanpemda.TujuanPemdaResponse{
		Id:               tujuanPemda.Id,
		TujuanPemdaId:    TujuanPemdaId.Id,
		NamaTujuanPemda:  TujuanPemdaId.NamaPohon,
		RumusPerhitungan: tujuanPemda.RumusPerhitungan,
		SumberData:       tujuanPemda.SumberData,
		Periode: tujuanpemda.PeriodeResponse{
			TahunAwal:  periode.TahunAwal,
			TahunAkhir: periode.TahunAkhir,
		},
		Indikator: indikatorResponses,
	}, nil
}

func (service *TujuanPemdaServiceImpl) Update(ctx context.Context, request tujuanpemda.TujuanPemdaUpdateRequest) (tujuanpemda.TujuanPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Convert request.Id string ke int
	tujuanPemdaId := request.Id

	// Validasi tujuan pemda exists
	tujuanPemda, err := service.TujuanPemdaRepository.FindById(ctx, tx, tujuanPemdaId)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}

	//validasi pohon kinerja = tematik
	err = service.PohonKinerjaRepository.ValidatePokinLevel(ctx, tx, request.TujuanPemdaId, 0, "tujuan pemda")
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}

	// Ambil data periode
	periode, err := service.PeriodeRepository.FindById(ctx, tx, tujuanPemda.PeriodeId)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf("periode tidak ditemukan: %v", err)
	}

	// Validasi tahun target untuk setiap indikator
	tahunAwal, _ := strconv.Atoi(periode.TahunAwal)
	tahunAkhir, _ := strconv.Atoi(periode.TahunAkhir)

	for _, indikatorReq := range request.Indikator {
		tahunMap := make(map[string]bool)
		for _, targetReq := range indikatorReq.Target {
			targetTahun, _ := strconv.Atoi(targetReq.Tahun)

			// Validasi rentang tahun
			if targetTahun < tahunAwal || targetTahun > tahunAkhir {
				return tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf(
					"tahun target %d harus berada dalam rentang periode %d-%d",
					targetTahun, tahunAwal, tahunAkhir,
				)
			}

			// Validasi duplikasi tahun
			if tahunMap[targetReq.Tahun] {
				return tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf(
					"duplikasi tahun %s pada indikator %s",
					targetReq.Tahun, indikatorReq.Indikator,
				)
			}
			tahunMap[targetReq.Tahun] = true
		}
	}

	// Update data tujuan pemda
	tujuanPemda.TujuanPemdaId = request.TujuanPemdaId
	tujuanPemda.RumusPerhitungan = request.RumusPerhitungan
	tujuanPemda.SumberData = request.SumberData

	// Proses indikator
	var indikators []domain.Indikator
	for _, indikatorReq := range request.Indikator {
		var targets []domain.Target

		// Proses target untuk setiap indikator
		for _, targetReq := range indikatorReq.Target {
			var targetId string
			if targetReq.Id != "" {
				targetId = targetReq.Id
			} else {
				targetId = generateTargetId()
			}

			target := domain.Target{
				Id:          targetId,
				IndikatorId: indikatorReq.Id,
				Target:      targetReq.Target,
				Satuan:      targetReq.Satuan,
				Tahun:       targetReq.Tahun,
			}
			targets = append(targets, target)
		}

		// Buat atau gunakan ID indikator yang ada
		var indikatorId string
		if indikatorReq.Id != "" {
			indikatorId = indikatorReq.Id
		} else {
			indikatorId = generateIndikatorId()
		}

		indikator := domain.Indikator{
			Id:            indikatorId,
			TujuanPemdaId: tujuanPemdaId,
			Indikator:     indikatorReq.Indikator,
			Target:        targets,
		}
		indikators = append(indikators, indikator)
	}

	tujuanPemda.Indikator = indikators

	// Simpan semua perubahan
	result, err := service.TujuanPemdaRepository.Update(ctx, tx, tujuanPemda)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}

	return service.toTujuanPemdaResponse(result), nil
}
func (service *TujuanPemdaServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Hapus indikator (cascade akan menghapus target)
	err = service.TujuanPemdaRepository.Delete(ctx, tx, id)
	if err != nil {
		return err
	}

	// Hapus tujuan pemda
	return service.TujuanPemdaRepository.Delete(ctx, tx, id)
}

func (service *TujuanPemdaServiceImpl) FindById(ctx context.Context, tujuanPemdaId int) (tujuanpemda.TujuanPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	tujuanPemda, err := service.TujuanPemdaRepository.FindById(ctx, tx, tujuanPemdaId)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}

	pokinData, err := service.PohonKinerjaRepository.FindById(ctx, tx, tujuanPemda.TujuanPemdaId)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf("gagal mengambil data pohon kinerja: %v", err)
	}

	var indikatorResponses []tujuanpemda.IndikatorResponse
	for _, indikator := range tujuanPemda.Indikator {
		indikatorResponse := tujuanpemda.IndikatorResponse{
			Id:        indikator.Id,
			Indikator: indikator.Indikator,
			Target:    []tujuanpemda.TargetResponse{},
		}

		// Tambahkan target hanya jika periode valid
		if tujuanPemda.PeriodeId != 0 && tujuanPemda.Periode.TahunAwal != "Pilih periode" {
			// Urutkan target berdasarkan tahun
			sort.Slice(indikator.Target, func(i, j int) bool {
				return indikator.Target[i].Tahun < indikator.Target[j].Tahun
			})

			indikatorResponse.Target = convertToTargetResponses(indikator.Target)
		}

		indikatorResponses = append(indikatorResponses, indikatorResponse)
	}

	return tujuanpemda.TujuanPemdaResponse{
		Id:               tujuanPemda.Id,
		TujuanPemdaId:    tujuanPemda.TujuanPemdaId,
		NamaTujuanPemda:  pokinData.NamaPohon,
		RumusPerhitungan: tujuanPemda.RumusPerhitungan,
		SumberData:       tujuanPemda.SumberData,
		Periode: tujuanpemda.PeriodeResponse{
			TahunAwal:  tujuanPemda.Periode.TahunAwal,
			TahunAkhir: tujuanPemda.Periode.TahunAkhir,
		},
		Indikator: indikatorResponses,
	}, nil
}

func (service *TujuanPemdaServiceImpl) FindAll(ctx context.Context, tahun string) ([]tujuanpemda.TujuanPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []tujuanpemda.TujuanPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	tujuanPemdaList, err := service.TujuanPemdaRepository.FindAll(ctx, tx, tahun)
	if err != nil {
		return []tujuanpemda.TujuanPemdaResponse{}, err
	}

	tujuanPemdaResponses := make([]tujuanpemda.TujuanPemdaResponse, 0, len(tujuanPemdaList))
	for _, tujuanPemda := range tujuanPemdaList {
		pokinData, err := service.PohonKinerjaRepository.FindById(ctx, tx, tujuanPemda.TujuanPemdaId)
		if err != nil {
			return []tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf("gagal mengambil data pohon kinerja: %v", err)
		}

		var indikatorResponses []tujuanpemda.IndikatorResponse
		for _, indikator := range tujuanPemda.Indikator {
			indikatorResponse := tujuanpemda.IndikatorResponse{
				Id:        indikator.Id,
				Indikator: indikator.Indikator,
				Target:    []tujuanpemda.TargetResponse{},
			}

			// Tambahkan semua target dalam rentang periode
			if tujuanPemda.PeriodeId != 0 && tujuanPemda.Periode.TahunAwal != "" {
				// Urutkan target berdasarkan tahun
				sort.Slice(indikator.Target, func(i, j int) bool {
					return indikator.Target[i].Tahun < indikator.Target[j].Tahun
				})
				indikatorResponse.Target = convertToTargetResponses(indikator.Target)
			}

			indikatorResponses = append(indikatorResponses, indikatorResponse)
		}

		tujuanPemdaResponses = append(tujuanPemdaResponses, tujuanpemda.TujuanPemdaResponse{
			Id:               tujuanPemda.Id,
			TujuanPemdaId:    tujuanPemda.TujuanPemdaId,
			NamaTujuanPemda:  pokinData.NamaPohon,
			RumusPerhitungan: tujuanPemda.RumusPerhitungan,
			SumberData:       tujuanPemda.SumberData,
			Periode: tujuanpemda.PeriodeResponse{
				TahunAwal:  tujuanPemda.Periode.TahunAwal,
				TahunAkhir: tujuanPemda.Periode.TahunAkhir,
			},
			Indikator: indikatorResponses,
		})
	}

	return tujuanPemdaResponses, nil
}

// Fungsi helper untuk konversi target
func convertToTargetResponses(targets []domain.Target) []tujuanpemda.TargetResponse {
	var targetResponses []tujuanpemda.TargetResponse
	for _, target := range targets {
		targetResponses = append(targetResponses, tujuanpemda.TargetResponse{
			Id:     target.Id,
			Target: target.Target,
			Satuan: target.Satuan,
			Tahun:  target.Tahun,
		})
	}
	return targetResponses
}

// Fungsi helper untuk konversi ke response
func (service *TujuanPemdaServiceImpl) toTujuanPemdaResponse(tujuanPemda domain.TujuanPemda) tujuanpemda.TujuanPemdaResponse {
	var indikatorResponses []tujuanpemda.IndikatorResponse
	for _, indikator := range tujuanPemda.Indikator {
		var targetResponses []tujuanpemda.TargetResponse
		for _, target := range indikator.Target {
			targetResponses = append(targetResponses, tujuanpemda.TargetResponse{
				Id:     target.Id,
				Target: target.Target,
				Satuan: target.Satuan,
				Tahun:  target.Tahun,
			})
		}

		indikatorResponses = append(indikatorResponses, tujuanpemda.IndikatorResponse{
			Id:        indikator.Id,
			Indikator: indikator.Indikator,
			Target:    targetResponses,
		})
	}

	return tujuanpemda.TujuanPemdaResponse{
		Id:               tujuanPemda.Id,
		TujuanPemdaId:    tujuanPemda.TujuanPemdaId,
		NamaTujuanPemda:  tujuanPemda.NamaTujuanPemda,
		RumusPerhitungan: tujuanPemda.RumusPerhitungan,
		SumberData:       tujuanPemda.SumberData,
		Periode: tujuanpemda.PeriodeResponse{
			TahunAwal:  tujuanPemda.Periode.TahunAwal,
			TahunAkhir: tujuanPemda.Periode.TahunAkhir,
		},
		Indikator: indikatorResponses,
	}
}

func (service *TujuanPemdaServiceImpl) UpdatePeriode(ctx context.Context, request tujuanpemda.TujuanPemdaUpdateRequest) (tujuanpemda.TujuanPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi tujuan pemda exists
	exists := service.TujuanPemdaRepository.IsIdExists(ctx, tx, request.Id)
	if !exists {
		return tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf("periode pemda dengan id %d tidak ditemukan", request.Id)
	}

	// Validasi periode exists jika periode_id tidak 0
	if request.PeriodeId != 0 {
		exists = service.PeriodeRepository.IsIdExists(ctx, tx, request.PeriodeId)
		if !exists {
			return tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf("periode dengan id %d tidak ditemukan", request.PeriodeId)
		}
	}

	// Persiapkan domain object untuk update (hanya id dan periode_id)
	tujuanPemda := domain.TujuanPemda{
		Id:        request.Id,
		PeriodeId: request.PeriodeId,
	}

	// Update periode
	result, err := service.TujuanPemdaRepository.UpdatePeriode(ctx, tx, tujuanPemda)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, err
	}

	// Ambil data pohon kinerja untuk response
	pokinData, err := service.PohonKinerjaRepository.FindById(ctx, tx, result.TujuanPemdaId)
	if err != nil {
		return tujuanpemda.TujuanPemdaResponse{}, fmt.Errorf("gagal mengambil data pohon kinerja: %v", err)
	}

	// Buat response
	return tujuanpemda.TujuanPemdaResponse{
		Id:              result.Id,
		TujuanPemdaId:   result.TujuanPemdaId,
		NamaTujuanPemda: pokinData.NamaPohon,
		Periode: tujuanpemda.PeriodeResponse{
			TahunAwal:  result.Periode.TahunAwal,
			TahunAkhir: result.Periode.TahunAkhir,
		},
	}, nil
}
