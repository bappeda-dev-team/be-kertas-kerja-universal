package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"ekak_kabupaten_madiun/repository"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type CrosscuttingOpdServiceImpl struct {
	CrosscuttingOpdRepository repository.CrosscuttingOpdRepository
	PohonKinerjaRepository    repository.PohonKinerjaRepository
	DB                        *sql.DB
}

func NewCrosscuttingOpdServiceImpl(crosscuttingOpdRepository repository.CrosscuttingOpdRepository, pohonKinerjaRepository repository.PohonKinerjaRepository, DB *sql.DB) *CrosscuttingOpdServiceImpl {
	return &CrosscuttingOpdServiceImpl{
		CrosscuttingOpdRepository: crosscuttingOpdRepository,
		PohonKinerjaRepository:    pohonKinerjaRepository,
		DB:                        DB,
	}
}

func (service *CrosscuttingOpdServiceImpl) Create(ctx context.Context, request pohonkinerja.CrosscuttingOpdCreateRequest, parentId int) (pohonkinerja.CrosscuttingOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.CrosscuttingOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Konversi request ke domain
	pokin := domain.PohonKinerja{
		NamaPohon:  request.NamaPohon,
		JenisPohon: request.JenisPohon,
		LevelPohon: request.LevelPohon,
		KodeOpd:    request.KodeOpd,
		Keterangan: request.Keterangan,
		Tahun:      request.Tahun,
		Status:     request.Status,
	}

	for _, indikatorReq := range request.Indikator {
		uuid := uuid.New().String()[:6]
		indikatorId := "IND-CRSS-" + uuid

		indikator := domain.Indikator{
			Id:        indikatorId,
			Indikator: indikatorReq.NamaIndikator,
			Tahun:     request.Tahun,
		}

		// Konversi target dengan generate ID
		for _, targetReq := range indikatorReq.Target {

			targetId := "TRG-CRSS-" + uuid

			target := domain.Target{
				Id:          targetId,
				IndikatorId: indikatorId, // Menggunakan ID indikator yang baru digenerate
				Target:      targetReq.Target,
				Satuan:      targetReq.Satuan,
				Tahun:       request.Tahun,
			}
			indikator.Target = append(indikator.Target, target)
		}
		pokin.Indikator = append(pokin.Indikator, indikator)
	}

	result, err := service.CrosscuttingOpdRepository.CreateCrosscutting(ctx, tx, pokin, parentId)
	if err != nil {
		return pohonkinerja.CrosscuttingOpdResponse{}, err
	}

	// Konversi ke response
	response := pohonkinerja.CrosscuttingOpdResponse{
		Id:         result.Id,
		NamaPohon:  result.NamaPohon,
		JenisPohon: result.JenisPohon,
		LevelPohon: result.LevelPohon,
		KodeOpd:    result.KodeOpd,
		Keterangan: result.Keterangan,
		Tahun:      result.Tahun,
		Status:     result.Status,
	}

	// Konversi indikator untuk response
	for _, indikator := range result.Indikator {
		indikatorResponse := pohonkinerja.IndikatorResponse{
			Id:            indikator.Id,
			NamaIndikator: indikator.Indikator,
		}

		// Konversi target untuk response
		for _, target := range indikator.Target {
			targetResponse := pohonkinerja.TargetResponse{
				Id:              target.Id,
				IndikatorId:     target.IndikatorId,
				TargetIndikator: target.Target,
				SatuanIndikator: target.Satuan,
			}
			indikatorResponse.Target = append(indikatorResponse.Target, targetResponse)
		}
		response.Indikator = append(response.Indikator, indikatorResponse)
	}

	return response, nil
}

func (service *CrosscuttingOpdServiceImpl) Update(ctx context.Context, request pohonkinerja.CrosscuttingOpdUpdateRequest) (pohonkinerja.CrosscuttingOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.CrosscuttingOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Cek apakah data exists
	existing, err := service.PohonKinerjaRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return pohonkinerja.CrosscuttingOpdResponse{}, errors.New("crosscutting tidak ditemukan")
	}

	// Konversi request ke domain
	pokin := domain.PohonKinerja{
		Id:         request.Id,
		NamaPohon:  request.NamaPohon,
		JenisPohon: request.JenisPohon,
		LevelPohon: request.LevelPohon,
		KodeOpd:    request.KodeOpd,
		Keterangan: request.Keterangan,
		Tahun:      request.Tahun,
		Status:     existing.Status, // Gunakan status yang ada
	}

	// Konversi indikator
	for _, indikatorReq := range request.Indikator {
		uuid := uuid.New().String()[:6]
		var indikatorId string
		if indikatorReq.Id != "" {
			indikatorId = indikatorReq.Id
		} else {
			indikatorId = "IND-CRSS-" + uuid
		}

		indikator := domain.Indikator{
			Id:        indikatorId,
			PokinId:   strconv.Itoa(request.Id),
			Indikator: indikatorReq.NamaIndikator,
			Tahun:     request.Tahun,
		}

		// Konversi target
		for _, targetReq := range indikatorReq.Target {
			var targetId string
			if targetReq.Id != "" {
				targetId = targetReq.Id
			} else {
				targetId = "TRG-CRSS-" + uuid
			}

			target := domain.Target{
				Id:          targetId,
				IndikatorId: indikatorId,
				Target:      targetReq.Target,
				Satuan:      targetReq.Satuan,
				Tahun:       request.Tahun,
			}
			indikator.Target = append(indikator.Target, target)
		}
		pokin.Indikator = append(pokin.Indikator, indikator)
	}

	// Update data
	result, err := service.CrosscuttingOpdRepository.UpdateCrosscutting(ctx, tx, pokin, request.ParentId)
	if err != nil {
		return pohonkinerja.CrosscuttingOpdResponse{}, err
	}

	// Konversi ke response
	response := pohonkinerja.CrosscuttingOpdResponse{
		Id:         result.Id,
		NamaPohon:  result.NamaPohon,
		JenisPohon: result.JenisPohon,
		LevelPohon: result.LevelPohon,
		KodeOpd:    result.KodeOpd,
		Keterangan: result.Keterangan,
		Tahun:      result.Tahun,
		Status:     result.Status,
	}

	// Konversi indikator untuk response
	for _, indikator := range result.Indikator {
		indikatorResponse := pohonkinerja.IndikatorResponse{
			Id:            indikator.Id,
			IdPokin:       strconv.Itoa(result.Id),
			NamaIndikator: indikator.Indikator,
		}

		// Konversi target untuk response
		for _, target := range indikator.Target {
			targetResponse := pohonkinerja.TargetResponse{
				Id:              target.Id,
				IndikatorId:     target.IndikatorId,
				TargetIndikator: target.Target,
				SatuanIndikator: target.Satuan,
			}
			indikatorResponse.Target = append(indikatorResponse.Target, targetResponse)
		}
		response.Indikator = append(response.Indikator, indikatorResponse)
	}

	return response, nil
}
func (service *CrosscuttingOpdServiceImpl) FindAllByParent(ctx context.Context, parentId int) ([]pohonkinerja.CrosscuttingOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Get all pohon kinerja
	pokins, err := service.CrosscuttingOpdRepository.FindAllCrosscutting(ctx, tx, parentId)
	if err != nil {
		return nil, err
	}

	// Collect all pokin IDs
	pokinIds := make([]int, len(pokins))
	for i, pokin := range pokins {
		pokinIds[i] = pokin.Id
	}

	// Get all indikator
	indikators, err := service.CrosscuttingOpdRepository.FindIndikatorByPokinId(ctx, tx, pokinIds)
	if err != nil {
		return nil, err
	}

	// Collect all indikator IDs
	indikatorIds := make([]string, len(indikators))
	for i, ind := range indikators {
		indikatorIds[i] = ind.Id
	}

	// Get all targets
	targets, err := service.CrosscuttingOpdRepository.FindTargetByIndikatorIds(ctx, tx, indikatorIds)
	if err != nil {
		return nil, err
	}

	// Create maps for easy lookup
	indikatorMap := make(map[string][]domain.Indikator) // Ubah ke string sebagai key
	for _, ind := range indikators {
		indikatorMap[ind.PokinId] = append(indikatorMap[ind.PokinId], ind)
	}

	targetMap := make(map[string][]domain.Target)
	for _, target := range targets {
		targetMap[target.IndikatorId] = append(targetMap[target.IndikatorId], target)
	}

	// Build response
	var responses []pohonkinerja.CrosscuttingOpdResponse
	for _, pokin := range pokins {
		response := pohonkinerja.CrosscuttingOpdResponse{
			Id:         pokin.Id,
			NamaPohon:  pokin.NamaPohon,
			JenisPohon: pokin.JenisPohon,
			LevelPohon: pokin.LevelPohon,
			KodeOpd:    pokin.KodeOpd,
			Keterangan: pokin.Keterangan,
			Tahun:      pokin.Tahun,
			Status:     pokin.Status,
		}

		// Add indikator
		pokinIdStr := strconv.Itoa(pokin.Id) // Konversi ID pohon kinerja ke string
		for _, indikator := range indikatorMap[pokinIdStr] {
			indikatorResponse := pohonkinerja.IndikatorResponse{
				Id:            indikator.Id,
				IdPokin:       indikator.PokinId, // Sudah string, tidak perlu konversi
				NamaIndikator: indikator.Indikator,
			}

			// Add target
			for _, target := range targetMap[indikator.Id] {
				targetResponse := pohonkinerja.TargetResponse{
					Id:              target.Id,
					IndikatorId:     target.IndikatorId,
					TargetIndikator: target.Target,
					SatuanIndikator: target.Satuan,
				}
				indikatorResponse.Target = append(indikatorResponse.Target, targetResponse)
			}
			response.Indikator = append(response.Indikator, indikatorResponse)
		}
		responses = append(responses, response)
	}

	return responses, nil
}

func (service *CrosscuttingOpdServiceImpl) ApproveOrReject(ctx context.Context, crosscuttingId int, request pohonkinerja.CrosscuttingApproveRequest) (*pohonkinerja.CrosscuttingApproveResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	currentTime := time.Now()
	var pegawaiAction map[string]interface{}

	if request.Approve {
		pegawaiAction = map[string]interface{}{
			"approve_by": request.NipPegawai,
			"approve_at": currentTime,
		}
	} else {
		pegawaiAction = map[string]interface{}{
			"reject_by": request.NipPegawai,
			"reject_at": currentTime,
		}
	}

	err = service.CrosscuttingOpdRepository.ApproveOrRejectCrosscutting(ctx, tx, crosscuttingId, request.Approve, pegawaiAction)
	if err != nil {
		return nil, err
	}

	response := &pohonkinerja.CrosscuttingApproveResponse{
		Id: crosscuttingId,
	}

	if request.Approve {
		response.Status = "crosscutting_disetujui"
		response.Message = "Crosscutting approved successfully"
		response.ApprovedBy = &request.NipPegawai
		response.ApprovedAt = &currentTime
	} else {
		response.Status = "crosscutting_ditolak"
		response.Message = "Crosscutting rejected successfully"
		response.RejectedBy = &request.NipPegawai
		response.RejectedAt = &currentTime
	}

	return response, nil
}

func (service *CrosscuttingOpdServiceImpl) Delete(ctx context.Context, pokinId int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi status sebelum delete
	var status string
	err = tx.QueryRowContext(ctx, "SELECT status FROM tb_pohon_kinerja WHERE id = ?", pokinId).Scan(&status)
	if err != nil {
		return err
	}

	if status != "crosscutting_disetujui" {
		return errors.New("crosscutting hanya dapat dihapus saat status crosscutting_disetujui")
	}

	err = service.CrosscuttingOpdRepository.DeleteCrosscutting(ctx, tx, pokinId)
	if err != nil {
		return err
	}

	return nil
}
