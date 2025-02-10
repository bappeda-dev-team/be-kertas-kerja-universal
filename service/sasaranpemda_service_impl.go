package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/sasaranpemda"
	"ekak_kabupaten_madiun/repository"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

type SasaranPemdaServiceImpl struct {
	SasaranPemdaRepository repository.SasaranPemdaRepository
	PeriodeRepository      repository.PeriodeRepository
	PohonKinerjaRepository repository.PohonKinerjaRepository
	TujuanPemdaRepository  repository.TujuanPemdaRepository
	DB                     *sql.DB
}

func NewSasaranPemdaServiceImpl(sasaranPemdaRepository repository.SasaranPemdaRepository, periodeRepository repository.PeriodeRepository, pohonKinerjaRepository repository.PohonKinerjaRepository, tujuanPemdaRepository repository.TujuanPemdaRepository, DB *sql.DB) *SasaranPemdaServiceImpl {
	return &SasaranPemdaServiceImpl{
		SasaranPemdaRepository: sasaranPemdaRepository,
		PeriodeRepository:      periodeRepository,
		PohonKinerjaRepository: pohonKinerjaRepository,
		TujuanPemdaRepository:  tujuanPemdaRepository,
		DB:                     DB,
	}
}

func (service *SasaranPemdaServiceImpl) generateRandomId(ctx context.Context, tx *sql.Tx) int {
	rand.Seed(time.Now().UnixNano())
	for {
		// Generate random number between 10000-99999
		id := rand.Intn(90000) + 10000
		if !service.SasaranPemdaRepository.IsIdExists(ctx, tx, id) {
			return id
		}
	}
}

func (service *SasaranPemdaServiceImpl) Create(ctx context.Context, request sasaranpemda.SasaranPemdaCreateRequest) (sasaranpemda.SasaranPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi periode
	// periode, err := service.PeriodeRepository.FindById(ctx, tx, request.PeriodeId)
	// if err != nil {
	// 	return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("periode tidak ditemukan: %v", err)
	// }

	//validasi pohon kinerja
	pokinData, err := service.PohonKinerjaRepository.FindById(ctx, tx, request.SubtemaId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("pohon kinerja tidak ditemukan: %v", err)
	}

	// Validasi level pohon kinerja (1)
	err = service.PohonKinerjaRepository.ValidatePokinLevel(ctx, tx, request.SubtemaId, 1, "sasaran pemda")
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}

	// // Validasi tahun target untuk setiap indikator
	// tahunAwal, _ := strconv.Atoi(periode.TahunAwal)
	// tahunAkhir, _ := strconv.Atoi(periode.TahunAkhir)

	// for _, indikatorReq := range request.Indikator {
	// 	tahunMap := make(map[string]bool)
	// 	for _, targetReq := range indikatorReq.Target {
	// 		targetTahun, _ := strconv.Atoi(targetReq.Tahun)

	// 		// Validasi rentang tahun
	// 		if targetTahun < tahunAwal || targetTahun > tahunAkhir {
	// 			return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf(
	// 				"tahun target %d harus berada dalam rentang periode %d-%d",
	// 				targetTahun, tahunAwal, tahunAkhir,
	// 			)
	// 		}

	// 		// Validasi duplikasi tahun
	// 		if tahunMap[targetReq.Tahun] {
	// 			return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf(
	// 				"duplikasi tahun %s pada indikator %s",
	// 				targetReq.Tahun, indikatorReq.Indikator,
	// 			)
	// 		}
	// 		tahunMap[targetReq.Tahun] = true
	// 	}
	// }

	if service.SasaranPemdaRepository.IsSubtemaIdExists(ctx, tx, request.SubtemaId) {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("pohon kinerja sudah digunakan untuk sasaran pemda lain")
	}

	sasaranPemda := domain.SasaranPemda{
		Id:            service.generateRandomId(ctx, tx),
		SubtemaId:     request.SubtemaId,
		TujuanPemdaId: request.TujuanPemdaId,
		SasaranPemda:  request.SasaranPemda,
		// PeriodeId:    request.PeriodeId,
	}

	sasaranPemda, err = service.SasaranPemdaRepository.Create(ctx, tx, sasaranPemda)
	if err != nil {
		if strings.Contains(err.Error(), "pohon kinerja dengan id") {
			return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("pohon kinerja sudah digunakan untuk sasaran pemda lain")
		}
		return sasaranpemda.SasaranPemdaResponse{}, err
	}

	return sasaranpemda.SasaranPemdaResponse{
		Id:            sasaranPemda.Id,
		TujuanPemdaId: sasaranPemda.TujuanPemdaId,
		SubtemaId:     sasaranPemda.SubtemaId,
		NamaSubtema:   pokinData.NamaPohon,
		SasaranPemda:  sasaranPemda.SasaranPemda,
		// Periode: sasaranpemda.PeriodeResponse{
		// 	TahunAwal:  periode.TahunAwal,
		// 	TahunAkhir: periode.TahunAkhir,
		// },

	}, nil
}

func (service *SasaranPemdaServiceImpl) Update(ctx context.Context, request sasaranpemda.SasaranPemdaUpdateRequest) (sasaranpemda.SasaranPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Convert request.Id string ke int
	sasaranPemdaId := request.Id

	// Validasi tujuan pemda exists
	tujuanPemda := service.TujuanPemdaRepository.IsIdExists(ctx, tx, request.TujuanPemdaId)
	if !tujuanPemda {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("tujuan pemda dengan id %d tidak ditemukan", request.TujuanPemdaId)
	}

	// Validasi sasaran pemda pemda exists
	sasaranPemda, err := service.SasaranPemdaRepository.FindById(ctx, tx, sasaranPemdaId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}
	//validasi pohon kinerja
	pokinData, err := service.PohonKinerjaRepository.FindById(ctx, tx, request.SubtemaId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("pohon kinerja tidak ditemukan: %v", err)
	}

	// Validasi level pohon kinerja (1-3)
	if pokinData.LevelPohon < 1 || pokinData.LevelPohon > 3 {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("level pohon kinerja harus berada di antara 1-3, level saat ini: %d", pokinData.LevelPohon)
	}

	// Update data sasaran pemda
	sasaranPemda.TujuanPemdaId = request.TujuanPemdaId
	sasaranPemda.SubtemaId = request.SubtemaId
	sasaranPemda.SasaranPemda = request.SasaranPemda

	// Simpan semua perubahan
	result, err := service.SasaranPemdaRepository.Update(ctx, tx, sasaranPemda)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}

	return service.toSasaranPemdaResponse(result), nil
}

func (service *SasaranPemdaServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	return service.SasaranPemdaRepository.Delete(ctx, tx, id)
}

func (service *SasaranPemdaServiceImpl) FindById(ctx context.Context, sasaranPemdaId int) (sasaranpemda.SasaranPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	sasaranPemda, err := service.SasaranPemdaRepository.FindById(ctx, tx, sasaranPemdaId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}

	pokinData, err := service.PohonKinerjaRepository.FindById(ctx, tx, sasaranPemda.SubtemaId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("gagal mengambil data pohon kinerja: %v", err)
	}

	var indikatorResponses []sasaranpemda.IndikatorResponse
	for _, indikator := range sasaranPemda.Indikator {
		indikatorResponse := sasaranpemda.IndikatorResponse{
			Id:        indikator.Id,
			Indikator: indikator.Indikator,
			Target:    []sasaranpemda.TargetResponse{},
		}

		// Tambahkan target hanya jika periode valid
		if sasaranPemda.PeriodeId != 0 && sasaranPemda.Periode.TahunAwal != "Pilih periode" {
			// Urutkan target berdasarkan tahun
			sort.Slice(indikator.Target, func(i, j int) bool {
				return indikator.Target[i].Tahun < indikator.Target[j].Tahun
			})

			indikatorResponse.Target = convertToTargetSasaranPemdaResponses(indikator.Target)
		}

		indikatorResponses = append(indikatorResponses, indikatorResponse)
	}

	return sasaranpemda.SasaranPemdaResponse{
		Id:           sasaranPemda.Id,
		SubtemaId:    sasaranPemda.SubtemaId,
		NamaSubtema:  pokinData.NamaPohon,
		SasaranPemda: sasaranPemda.SasaranPemda,
		Periode: sasaranpemda.PeriodeResponse{
			TahunAwal:  sasaranPemda.Periode.TahunAwal,
			TahunAkhir: sasaranPemda.Periode.TahunAkhir,
		},
		Indikator: indikatorResponses,
	}, nil
}

func (service *SasaranPemdaServiceImpl) FindAll(ctx context.Context, tahun string) ([]sasaranpemda.SasaranPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return []sasaranpemda.SasaranPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	sasaranPemdaList, err := service.SasaranPemdaRepository.FindAll(ctx, tx, tahun)
	if err != nil {
		return []sasaranpemda.SasaranPemdaResponse{}, err
	}

	sasaranPemdaResponses := make([]sasaranpemda.SasaranPemdaResponse, 0, len(sasaranPemdaList))
	for _, sasaranPemda := range sasaranPemdaList {
		pokinData, err := service.PohonKinerjaRepository.FindById(ctx, tx, sasaranPemda.SubtemaId)
		if err != nil {
			return []sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("gagal mengambil data pohon kinerja: %v", err)
		}

		var indikatorResponses []sasaranpemda.IndikatorResponse
		for _, indikator := range sasaranPemda.Indikator {
			indikatorResponse := sasaranpemda.IndikatorResponse{
				Id:        indikator.Id,
				Indikator: indikator.Indikator,
				Target:    []sasaranpemda.TargetResponse{},
			}

			// Tambahkan semua target dalam rentang periode
			if sasaranPemda.PeriodeId != 0 && sasaranPemda.Periode.TahunAwal != "" {
				// Urutkan target berdasarkan tahun
				sort.Slice(indikator.Target, func(i, j int) bool {
					return indikator.Target[i].Tahun < indikator.Target[j].Tahun
				})
				indikatorResponse.Target = convertToTargetSasaranPemdaResponses(indikator.Target)
			}

			indikatorResponses = append(indikatorResponses, indikatorResponse)
		}

		sasaranPemdaResponses = append(sasaranPemdaResponses, sasaranpemda.SasaranPemdaResponse{
			Id:           sasaranPemda.Id,
			SubtemaId:    sasaranPemda.SubtemaId,
			NamaSubtema:  pokinData.NamaPohon,
			SasaranPemda: sasaranPemda.SasaranPemda,
			Periode: sasaranpemda.PeriodeResponse{
				TahunAwal:  sasaranPemda.Periode.TahunAwal,
				TahunAkhir: sasaranPemda.Periode.TahunAkhir,
			},
			Indikator: indikatorResponses,
		})
	}

	return sasaranPemdaResponses, nil
}

// Fungsi helper untuk konversi target
func convertToTargetSasaranPemdaResponses(targets []domain.Target) []sasaranpemda.TargetResponse {
	var targetResponses []sasaranpemda.TargetResponse
	for _, target := range targets {
		targetResponses = append(targetResponses, sasaranpemda.TargetResponse{
			Id:     target.Id,
			Target: target.Target,
			Satuan: target.Satuan,
			Tahun:  target.Tahun,
		})
	}
	return targetResponses
}

func (service *SasaranPemdaServiceImpl) toSasaranPemdaResponse(sasaranPemda domain.SasaranPemda) sasaranpemda.SasaranPemdaResponse {
	var indikatorResponses []sasaranpemda.IndikatorResponse
	for _, indikator := range sasaranPemda.Indikator {
		var targetResponses []sasaranpemda.TargetResponse
		for _, target := range indikator.Target {
			targetResponses = append(targetResponses, sasaranpemda.TargetResponse{
				Id:     target.Id,
				Target: target.Target,
				Satuan: target.Satuan,
				Tahun:  target.Tahun,
			})
		}

		indikatorResponses = append(indikatorResponses, sasaranpemda.IndikatorResponse{
			Id:        indikator.Id,
			Indikator: indikator.Indikator,
			Target:    targetResponses,
		})
	}

	return sasaranpemda.SasaranPemdaResponse{
		Id:           sasaranPemda.Id,
		SubtemaId:    sasaranPemda.SubtemaId,
		SasaranPemda: sasaranPemda.SasaranPemda,
		Periode: sasaranpemda.PeriodeResponse{
			TahunAwal:  sasaranPemda.Periode.TahunAwal,
			TahunAkhir: sasaranPemda.Periode.TahunAkhir,
		},
		Indikator: indikatorResponses,
	}
}

func (service *SasaranPemdaServiceImpl) FindAllWithPokin(ctx context.Context, tahun string) ([]sasaranpemda.SasaranPemdaWithPokinResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data dari repository
	items, err := service.SasaranPemdaRepository.FindAllWithPokin(ctx, tx, tahun)
	if err != nil {
		return nil, err
	}

	// Transform ke response
	var responses []sasaranpemda.SasaranPemdaWithPokinResponse
	for _, item := range items {
		var indikatorResponses []sasaranpemda.IndikatorSubtematikResponse
		for _, indikator := range item.IndikatorSubtematik {
			var targetResponses []sasaranpemda.TargetResponse
			for _, target := range indikator.Target {
				targetResponses = append(targetResponses, sasaranpemda.TargetResponse{
					Target: target.Target,
					Satuan: target.Satuan,
				})
			}

			indikatorResponses = append(indikatorResponses, sasaranpemda.IndikatorSubtematikResponse{
				Indikator: indikator.Indikator,
				Target:    targetResponses,
			})
		}

		responses = append(responses, sasaranpemda.SasaranPemdaWithPokinResponse{
			IdsasaranPemda:      item.IdsasaranPemda,
			TematikId:           item.TematikId,
			NamaTematik:         item.NamaTematik,
			SubtematikId:        item.SubtematikId,
			NamaSubtematik:      item.NamaSubtematik,
			JenisPohon:          item.JenisPohon,
			LevelPohon:          item.LevelPohon,
			SasaranPemda:        item.SasaranPemda,
			Keterangan:          item.Keterangan,
			IndikatorSubtematik: indikatorResponses,
		})
	}

	return responses, nil
}
