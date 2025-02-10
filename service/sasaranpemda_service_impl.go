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

	//validasi pohon kinerja
	pokinData, err := service.PohonKinerjaRepository.FindById(ctx, tx, request.SubtemaId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("pohon kinerja tidak ditemukan: %v", err)
	}

	// Validasi level pohon kinerja (1-3)
	if pokinData.LevelPohon < 1 || pokinData.LevelPohon > 3 {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("level pohon kinerja harus berada di antara 1-3, level saat ini: %d", pokinData.LevelPohon)
	}
	// Validasi subtema id belum digunakan
	if service.SasaranPemdaRepository.IsSubtemaIdExists(ctx, tx, request.SubtemaId) {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("pohon kinerja dengan id %d sudah digunakan untuk sasaran pemda lain", request.SubtemaId)
	}

	// Validasi tujuan pemda exists
	exists := service.TujuanPemdaRepository.IsIdExists(ctx, tx, request.TujuanPemdaId)
	if !exists {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("tujuan pemda dengan id %d tidak ditemukan", request.TujuanPemdaId)
	}

	sasaranPemda := domain.SasaranPemda{
		Id:            service.generateRandomId(ctx, tx),
		SubtemaId:     request.SubtemaId,
		TujuanPemdaId: request.TujuanPemdaId,
		SasaranPemda:  request.SasaranPemda,
	}

	sasaranPemda, err = service.SasaranPemdaRepository.Create(ctx, tx, sasaranPemda)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("gagal membuat sasaran pemda: %v", err)
	}

	return sasaranpemda.SasaranPemdaResponse{
		Id:            sasaranPemda.Id,
		TujuanPemdaId: sasaranPemda.TujuanPemdaId,
		SubtemaId:     sasaranPemda.SubtemaId,
		NamaSubtema:   pokinData.NamaPohon,
		SasaranPemda:  sasaranPemda.SasaranPemda,
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
	_, err = service.TujuanPemdaRepository.FindById(ctx, tx, request.TujuanPemdaId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("tujuan pemda tidak ditemukan: %v", err)
	}

	if !service.TujuanPemdaRepository.IsIdExists(ctx, tx, request.TujuanPemdaId) {
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

	return sasaranpemda.SasaranPemdaResponse{
		Id:            result.Id,
		TujuanPemdaId: result.TujuanPemdaId,
		SubtemaId:     result.SubtemaId,
		NamaSubtema:   pokinData.NamaPohon,
		SasaranPemda:  result.SasaranPemda,
	}, nil
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

	// Ambil data sasaran pemda
	sasaranPemda, err := service.SasaranPemdaRepository.FindById(ctx, tx, sasaranPemdaId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, err
	}

	// Ambil data pohon kinerja untuk nama subtema
	pokinData, err := service.PohonKinerjaRepository.FindById(ctx, tx, sasaranPemda.SubtemaId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("gagal mengambil data pohon kinerja: %v", err)
	}

	// Ambil data tujuan pemda
	tujuanPemda, err := service.TujuanPemdaRepository.FindById(ctx, tx, sasaranPemda.TujuanPemdaId)
	if err != nil {
		return sasaranpemda.SasaranPemdaResponse{}, fmt.Errorf("gagal mengambil data tujuan pemda: %v", err)
	}

	return sasaranpemda.SasaranPemdaResponse{
		Id:            sasaranPemda.Id,
		TujuanPemdaId: sasaranPemda.TujuanPemdaId,
		TujuanPemda:   tujuanPemda.TujuanPemda, // Tambahkan ini di struct response
		SubtemaId:     sasaranPemda.SubtemaId,
		NamaSubtema:   pokinData.NamaPohon,
		SasaranPemda:  sasaranPemda.SasaranPemda,
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

		sasaranPemdaResponses = append(sasaranPemdaResponses, sasaranpemda.SasaranPemdaResponse{
			Id:           sasaranPemda.Id,
			SubtemaId:    sasaranPemda.SubtemaId,
			NamaSubtema:  pokinData.NamaPohon,
			SasaranPemda: sasaranPemda.SasaranPemda,
		})
	}

	return sasaranPemdaResponses, nil
}

// Fungsi helper untuk konversi target

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
