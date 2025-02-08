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
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SasaranPemdaServiceImpl struct {
	SasaranPemdaRepository repository.SasaranPemdaRepository
	PeriodeRepository      repository.PeriodeRepository
	PohonKinerjaRepository repository.PohonKinerjaRepository
	DB                     *sql.DB
}

func NewSasaranPemdaServiceImpl(sasaranPemdaRepository repository.SasaranPemdaRepository, periodeRepository repository.PeriodeRepository, pohonKinerjaRepository repository.PohonKinerjaRepository, tujuanPemdaRepository repository.TujuanPemdaRepository, DB *sql.DB) *SasaranPemdaServiceImpl {
	return &SasaranPemdaServiceImpl{
		SasaranPemdaRepository: sasaranPemdaRepository,
		PeriodeRepository:      periodeRepository,
		PohonKinerjaRepository: pohonKinerjaRepository,
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

func generateIndikatorSasaranPemda() string {
	currentYear := time.Now().Format("2006")
	uuid := uuid.New().String()[:5] // Mengambil 5 karakter pertama dari UUID
	return fmt.Sprintf("IND-SSRN-PMD-%s-%s", currentYear, uuid)
}

func generateTargetSasaranPemda() string {
	currentYear := time.Now().Format("2006")
	uuid := uuid.New().String()[:5]
	return fmt.Sprintf("TRG-SSRN-PMD-%s-%s", currentYear, uuid)
}

func (service *SasaranPemdaServiceImpl) Create(ctx context.Context, request sasaranpemda.SasaranPemdaCreateRequest) (sasaranpemda.SasaranPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi periode
	periode, err := service.PeriodeRepository.FindById(ctx, tx, request.PeriodeId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("periode tidak ditemukan: %v", err)
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

	// Validasi tahun target untuk setiap indikator
	tahunAwal, _ := strconv.Atoi(periode.TahunAwal)
	tahunAkhir, _ := strconv.Atoi(periode.TahunAkhir)

	for _, indikatorReq := range request.Indikator {
		tahunMap := make(map[string]bool)
		for _, targetReq := range indikatorReq.Target {
			targetTahun, _ := strconv.Atoi(targetReq.Tahun)

			// Validasi rentang tahun
			if targetTahun < tahunAwal || targetTahun > tahunAkhir {
				return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf(
					"tahun target %d harus berada dalam rentang periode %d-%d",
					targetTahun, tahunAwal, tahunAkhir,
				)
			}

			// Validasi duplikasi tahun
			if tahunMap[targetReq.Tahun] {
				return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf(
					"duplikasi tahun %s pada indikator %s",
					targetReq.Tahun, indikatorReq.Indikator,
				)
			}
			tahunMap[targetReq.Tahun] = true
		}
	}

	if service.SasaranPemdaRepository.IsSubtemaIdExists(ctx, tx, request.SubtemaId) {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("pohon kinerja sudah digunakan untuk sasaran pemda lain")
	}

	sasaranPemda := domain.SasaranPemda{
		Id:               service.generateRandomId(ctx, tx),
		SubtemaId:        request.SubtemaId,
		SasaranPemda:     request.SasaranPemda,
		RumusPerhitungan: request.RumusPerhitungan,
		SumberData:       request.SumberData,
		PeriodeId:        request.PeriodeId,
	}

	sasaranPemda, err = service.SasaranPemdaRepository.Create(ctx, tx, sasaranPemda)
	if err != nil {
		if strings.Contains(err.Error(), "pohon kinerja dengan id") {
			return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("pohon kinerja sudah digunakan untuk sasaran pemda lain")
		}
		return sasaranpemda.SasaranPemdaResponse{}, err
	}

	var indikatorResponses []sasaranpemda.IndikatorResponse
	for _, indikatorRequest := range request.Indikator {
		indikator := domain.Indikator{
			Id:             generateIndikatorSasaranPemda(),
			SasaranPemdaId: sasaranPemda.Id,
			Indikator:      indikatorRequest.Indikator,
		}

		indikator, err = service.SasaranPemdaRepository.CreateIndikator(ctx, tx, indikator)
		if err != nil {
			return sasaranpemda.SasaranPemdaResponse{}, err
		}

		var targetResponses []sasaranpemda.TargetResponse
		for _, targetRequest := range indikatorRequest.Target {
			target := domain.Target{
				Id:          generateTargetSasaranPemda(),
				IndikatorId: indikator.Id,
				Target:      targetRequest.Target,
				Satuan:      targetRequest.Satuan,
				Tahun:       targetRequest.Tahun,
			}

			err = service.SasaranPemdaRepository.CreateTarget(ctx, tx, target)
			if err != nil {
				return sasaranpemda.SasaranPemdaResponse{}, err
			}

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
		Id:               sasaranPemda.Id,
		SubtemaId:        sasaranPemda.SubtemaId,
		NamaSubtema:      pokinData.NamaPohon,
		SasaranPemda:     sasaranPemda.SasaranPemda,
		RumusPerhitungan: sasaranPemda.RumusPerhitungan,
		SumberData:       sasaranPemda.SumberData,
		Periode: sasaranpemda.PeriodeResponse{
			TahunAwal:  periode.TahunAwal,
			TahunAkhir: periode.TahunAkhir,
		},
		Indikator: indikatorResponses,
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
	sasaranPemda, err := service.SasaranPemdaRepository.FindById(ctx, tx, sasaranPemdaId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}

	// Ambil data periode
	periode, err := service.PeriodeRepository.FindById(ctx, tx, sasaranPemda.PeriodeId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("periode tidak ditemukan: %v", err)
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

	// Validasi tahun target untuk setiap indikator
	tahunAwal, _ := strconv.Atoi(periode.TahunAwal)
	tahunAkhir, _ := strconv.Atoi(periode.TahunAkhir)

	for _, indikatorReq := range request.Indikator {
		tahunMap := make(map[string]bool)
		for _, targetReq := range indikatorReq.Target {
			targetTahun, _ := strconv.Atoi(targetReq.Tahun)

			// Validasi rentang tahun
			if targetTahun < tahunAwal || targetTahun > tahunAkhir {
				return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf(
					"tahun target %d harus berada dalam rentang periode %d-%d",
					targetTahun, tahunAwal, tahunAkhir,
				)
			}

			// Validasi duplikasi tahun
			if tahunMap[targetReq.Tahun] {
				return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf(
					"duplikasi tahun %s pada indikator %s",
					targetReq.Tahun, indikatorReq.Indikator,
				)
			}
			tahunMap[targetReq.Tahun] = true
		}
	}

	// Update data sasaran pemda
	sasaranPemda.SubtemaId = request.SubtemaId
	sasaranPemda.SasaranPemda = request.SasaranPemda
	sasaranPemda.RumusPerhitungan = request.RumusPerhitungan
	sasaranPemda.SumberData = request.SumberData

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
			Id:             indikatorId,
			SasaranPemdaId: sasaranPemdaId,
			Indikator:      indikatorReq.Indikator,
			Target:         targets,
		}
		indikators = append(indikators, indikator)
	}

	sasaranPemda.Indikator = indikators

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
		Id:               sasaranPemda.Id,
		SubtemaId:        sasaranPemda.SubtemaId,
		NamaSubtema:      pokinData.NamaPohon,
		SasaranPemda:     sasaranPemda.SasaranPemda,
		RumusPerhitungan: sasaranPemda.RumusPerhitungan,
		SumberData:       sasaranPemda.SumberData,
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
			Id:               sasaranPemda.Id,
			SubtemaId:        sasaranPemda.SubtemaId,
			NamaSubtema:      pokinData.NamaPohon,
			SasaranPemda:     sasaranPemda.SasaranPemda,
			RumusPerhitungan: sasaranPemda.RumusPerhitungan,
			SumberData:       sasaranPemda.SumberData,
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
		Id:               sasaranPemda.Id,
		SubtemaId:        sasaranPemda.SubtemaId,
		SasaranPemda:     sasaranPemda.SasaranPemda,
		RumusPerhitungan: sasaranPemda.RumusPerhitungan,
		SumberData:       sasaranPemda.SumberData,
		Periode: sasaranpemda.PeriodeResponse{
			TahunAwal:  sasaranPemda.Periode.TahunAwal,
			TahunAkhir: sasaranPemda.Periode.TahunAkhir,
		},
		Indikator: indikatorResponses,
	}
}

func (service *SasaranPemdaServiceImpl) UpdatePeriode(ctx context.Context, request sasaranpemda.SasaranPemdaUpdateRequest) (sasaranpemda.SasaranPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi tujuan pemda exists
	exists := service.SasaranPemdaRepository.IsIdExists(ctx, tx, request.Id)
	if !exists {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("periode pemda dengan id %d tidak ditemukan", request.Id)
	}

	// Validasi periode exists jika periode_id tidak 0
	if request.PeriodeId != 0 {
		exists = service.PeriodeRepository.IsIdExists(ctx, tx, request.PeriodeId)
		if !exists {
			return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("periode dengan id %d tidak ditemukan", request.PeriodeId)
		}
	}

	// Persiapkan domain object untuk update (hanya id dan periode_id)
	sasaranPemda := domain.SasaranPemda{
		Id:        request.Id,
		PeriodeId: request.PeriodeId,
	}

	// Update periode
	result, err := service.SasaranPemdaRepository.UpdatePeriode(ctx, tx, sasaranPemda)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}

	// Ambil data pohon kinerja untuk response
	pokinData, err := service.PohonKinerjaRepository.FindById(ctx, tx, result.SubtemaId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("gagal mengambil data pohon kinerja: %v", err)
	}

	// Buat response
	return sasaranpemda.SasaranPemdaResponse{
		Id:           result.Id,
		SubtemaId:    result.SubtemaId,
		SasaranPemda: result.SasaranPemda,
		NamaSubtema:  pokinData.NamaPohon,
		Periode: sasaranpemda.PeriodeResponse{
			TahunAwal:  result.Periode.TahunAwal,
			TahunAkhir: result.Periode.TahunAkhir,
		},
	}, nil
}

func (service *SasaranPemdaServiceImpl) FindAllWithPokin(ctx context.Context, tahun string) ([]sasaranpemda.SasaranPemdaWithPokinResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi apakah tahun berada dalam periode yang valid
	periode, err := service.PeriodeRepository.FindByTahun(ctx, tx, tahun)
	if err != nil {
		return nil, fmt.Errorf("tahun %s tidak berada dalam periode manapun", tahun)
	}

	sasaranPemdaList, err := service.SasaranPemdaRepository.FindAllWithPokin(ctx, tx, tahun)
	if err != nil {
		return nil, err
	}

	var responses []sasaranpemda.SasaranPemdaWithPokinResponse
	for _, item := range sasaranPemdaList {
		var sasaranPemdaResponse *sasaranpemda.SasaranPemdaResponse
		if item.SasaranPemda != nil {
			var indikatorResponses []sasaranpemda.IndikatorResponse
			for _, indikator := range item.SasaranPemda.Indikator {
				var targetResponses []sasaranpemda.TargetResponse

				// Generate target untuk semua tahun dalam periode
				tahunAwal, _ := strconv.Atoi(periode.TahunAwal)
				tahunAkhir, _ := strconv.Atoi(periode.TahunAkhir)

				targetMap := make(map[string]domain.Target)
				for _, target := range indikator.Target {
					targetMap[target.Tahun] = target
				}

				for tahun := tahunAwal; tahun <= tahunAkhir; tahun++ {
					tahunStr := strconv.Itoa(tahun)
					if target, exists := targetMap[tahunStr]; exists {
						targetResponses = append(targetResponses, sasaranpemda.TargetResponse{
							Id:     target.Id,
							Target: target.Target,
							Satuan: target.Satuan,
							Tahun:  target.Tahun,
						})
					} else {
						targetResponses = append(targetResponses, sasaranpemda.TargetResponse{
							Id:     "-",
							Target: "-",
							Satuan: "-",
							Tahun:  tahunStr,
						})
					}
				}

				indikatorResponses = append(indikatorResponses, sasaranpemda.IndikatorResponse{
					Id:        indikator.Id,
					Indikator: indikator.Indikator,
					Target:    targetResponses,
				})
			}

			sasaranPemdaResponse = &sasaranpemda.SasaranPemdaResponse{
				Id:           item.SasaranPemda.Id,
				SasaranPemda: item.SasaranPemda.SasaranPemda,
				// TematikId:        item.TujuanPemda.TematikId,
				RumusPerhitungan: item.SasaranPemda.RumusPerhitungan,
				SumberData:       item.SasaranPemda.SumberData,
				Periode: sasaranpemda.PeriodeResponse{
					TahunAwal:  periode.TahunAwal,
					TahunAkhir: periode.TahunAkhir,
				},
				Indikator: indikatorResponses,
			}
		}

		responses = append(responses, sasaranpemda.SasaranPemdaWithPokinResponse{
			PokinId:      item.PokinId,
			NamaPohon:    item.NamaPohon,
			JenisPohon:   item.JenisPohon,
			LevelPohon:   item.LevelPohon,
			Keterangan:   item.Keterangan,
			TahunPokin:   item.TahunPokin,
			SasaranPemda: sasaranPemdaResponse,
		})
	}

	return responses, nil
}

func (service *SasaranPemdaServiceImpl) FindPokinWithPeriode(ctx context.Context, pokinId int) (sasaranpemda.PokinWithPeriodeResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return sasaranpemda.PokinWithPeriodeResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi pokin ID
	err = service.PohonKinerjaRepository.ValidatePokinId(ctx, tx, pokinId)
	if err != nil {
		return sasaranpemda.PokinWithPeriodeResponse{}, err
	}

	// Ambil data pokin dengan periode
	pokin, err := service.PohonKinerjaRepository.FindPokinWithPeriode(ctx, tx, pokinId)
	if err != nil {
		return sasaranpemda.PokinWithPeriodeResponse{}, err
	}

	// Transform ke response
	response := sasaranpemda.PokinWithPeriodeResponse{
		Id:         pokin.Id,
		NamaPohon:  pokin.NamaPohon,
		JenisPohon: pokin.JenisPohon,
		LevelPohon: pokin.LevelPohon,
		Tahun:      pokin.Tahun,
	}

	// Transform indikator dan target
	for _, indikator := range pokin.Indikator {
		indikatorResponse := sasaranpemda.PokinIndikatorResponse{
			Id:        indikator.Id,
			Indikator: indikator.Indikator,
		}

		// Transform target untuk setiap tahun dalam periode
		for _, target := range indikator.Target {
			targetResponse := sasaranpemda.PokinTargetResponse{
				Id:     target.Id,
				Target: target.Target,
				Satuan: target.Satuan,
				Tahun:  target.Tahun,
			}
			indikatorResponse.Target = append(indikatorResponse.Target, targetResponse)
		}

		response.Indikator = append(response.Indikator, indikatorResponse)
	}

	// Cari periode yang sesuai dengan tahun pokin
	periode, err := service.PeriodeRepository.FindByTahun(ctx, tx, pokin.Tahun)
	if err == nil {
		// Jika periode ditemukan
		response.Periode = sasaranpemda.PokinPeriodeResponse{
			Id:         periode.Id,
			TahunAwal:  periode.TahunAwal,
			TahunAkhir: periode.TahunAkhir,
		}
	} else {
		// Jika periode tidak ditemukan, set nilai default
		response.Periode = sasaranpemda.PokinPeriodeResponse{
			Id:         0,
			TahunAwal:  "",
			TahunAkhir: "",
		}
	}

	return response, nil
}
