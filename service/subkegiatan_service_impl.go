package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/subkegiatan"
	"ekak_kabupaten_madiun/repository"
	"errors"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type SubKegiatanServiceImpl struct {
	subKegiatanRepository   repository.SubKegiatanRepository
	opdRepository           repository.OpdRepository
	rencanaKinerjaRepoitory repository.RencanaKinerjaRepository
	DB                      *sql.DB
	validator               *validator.Validate
}

func NewSubKegiatanServiceImpl(subKegiatanRepository repository.SubKegiatanRepository, opdRepository repository.OpdRepository, rencanaKinerjaRepoitory repository.RencanaKinerjaRepository, DB *sql.DB, validator *validator.Validate) *SubKegiatanServiceImpl {
	return &SubKegiatanServiceImpl{
		subKegiatanRepository:   subKegiatanRepository,
		opdRepository:           opdRepository,
		rencanaKinerjaRepoitory: rencanaKinerjaRepoitory,
		DB:                      DB,
		validator:               validator,
	}
}

func (service *SubKegiatanServiceImpl) Create(ctx context.Context, request subkegiatan.SubKegiatanCreateRequest) (subkegiatan.SubKegiatanResponse, error) {
	err := service.validator.Struct(request)
	if err != nil {
		log.Println("Validasi gagal:", err)
		return subkegiatan.SubKegiatanResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("Gagal memulai transaksi:", err)
		return subkegiatan.SubKegiatanResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
	uuId := fmt.Sprintf("SUB-%s", randomDigits)

	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, request.KodeOpd)
	if err != nil {
		return subkegiatan.SubKegiatanResponse{}, fmt.Errorf("kode OPD tidak valid: %v", err)
	}

	if opd.KodeOpd == "" {
		return subkegiatan.SubKegiatanResponse{}, fmt.Errorf("kode OPD tidak ditemukan")
	}

	var indikators []domain.Indikator

	for _, indikatorReq := range request.Indikator {
		indikatorId := indikatorReq.Id
		if indikatorId == "" {
			randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
			indikatorId = fmt.Sprintf("IND-SUB-%s", randomDigits)
		}

		var targets []domain.Target

		for _, targetReq := range indikatorReq.Target {
			targetId := targetReq.Id
			if targetId == "" {
				randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
				targetId = fmt.Sprintf("TRG-SUB-%s", randomDigits)
			}

			target := domain.Target{
				Id:          targetId,
				IndikatorId: indikatorId,
				Target:      targetReq.TargetIndikator,
				Satuan:      targetReq.SatuanIndikator,
				Tahun:       request.Tahun,
			}
			targets = append(targets, target)
		}

		indikator := domain.Indikator{
			Id:            indikatorId,
			SubKegiatanId: uuId,
			Indikator:     indikatorReq.NamaIndikator,
			Tahun:         request.Tahun,
			Target:        targets,
		}
		indikators = append(indikators, indikator)
	}

	subKegiatan := domain.SubKegiatan{
		Id:              uuId,
		NamaSubKegiatan: request.NamaSubKegiatan,
		KodeOpd:         request.KodeOpd,
		Status:          "belum_diambil",
		Tahun:           request.Tahun,
		Indikator:       indikators,
	}

	result, err := service.subKegiatanRepository.Create(ctx, tx, subKegiatan)
	if err != nil {
		log.Println("Gagal membuat data sub kegiatan:", err)
		return subkegiatan.SubKegiatanResponse{}, err
	}

	return helper.ToSubKegiatanResponse(result), nil
}

func (service *SubKegiatanServiceImpl) Update(ctx context.Context, request subkegiatan.SubKegiatanUpdateRequest) (subkegiatan.SubKegiatanResponse, error) {
	err := service.validator.Struct(request)
	if err != nil {
		log.Println("Validasi gagal:", err)
		return subkegiatan.SubKegiatanResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		return subkegiatan.SubKegiatanResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, request.KodeOpd)
	if err != nil {
		return subkegiatan.SubKegiatanResponse{}, fmt.Errorf("kode OPD tidak valid: %v", err)
	}

	if opd.KodeOpd == "" {
		return subkegiatan.SubKegiatanResponse{}, fmt.Errorf("kode OPD tidak ditemukan")
	}

	var indikators []domain.Indikator

	for _, indikatorReq := range request.Indikator {
		indikatorId := indikatorReq.Id
		if indikatorId == "" {
			randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
			indikatorId = fmt.Sprintf("IND-%s", randomDigits)
		}

		var targets []domain.Target

		for _, targetReq := range indikatorReq.Target {
			targetId := targetReq.Id
			if targetId == "" {
				randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
				targetId = fmt.Sprintf("TRG-%s", randomDigits)
			}

			target := domain.Target{
				Id:          targetId,
				IndikatorId: indikatorId,
				Target:      targetReq.TargetIndikator,
				Satuan:      targetReq.SatuanIndikator,
				Tahun:       request.Tahun,
			}
			targets = append(targets, target)
		}

		indikator := domain.Indikator{
			Id:               indikatorId,
			SubKegiatanId:    request.Id,
			RencanaKinerjaId: indikatorReq.RencanaKinerjaId,
			Indikator:        indikatorReq.NamaIndikator,
			Tahun:            request.Tahun,
			Target:           targets,
		}
		indikators = append(indikators, indikator)
	}

	domainSubKegiatan := domain.SubKegiatan{
		Id:              request.Id,
		NamaSubKegiatan: request.NamaSubKegiatan,
		KodeOpd:         request.KodeOpd,
		Tahun:           request.Tahun,
		Indikator:       indikators,
	}

	result, err := service.subKegiatanRepository.Update(ctx, tx, domainSubKegiatan)
	if err != nil {
		return subkegiatan.SubKegiatanResponse{}, fmt.Errorf("gagal mengupdate sub kegiatan: %v", err)
	}

	response := helper.ToSubKegiatanResponse(result)
	return response, nil
}

func (service *SubKegiatanServiceImpl) FindById(ctx context.Context, subKegiatanId string) (subkegiatan.SubKegiatanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("Gagal memulai transaksi:", err)
		return subkegiatan.SubKegiatanResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data SubKegiatan
	subKegiatan, err := service.subKegiatanRepository.FindById(ctx, tx, subKegiatanId)
	if err != nil {
		if err == sql.ErrNoRows {
			return subkegiatan.SubKegiatanResponse{}, fmt.Errorf("sub kegiatan dengan id %s tidak ditemukan", subKegiatanId)
		}
		log.Println("Gagal mencari data sub kegiatan:", err)
		return subkegiatan.SubKegiatanResponse{}, err
	}

	// Ambil data Indikator
	indikators, err := service.subKegiatanRepository.FindIndikatorBySubKegiatanId(ctx, tx, subKegiatanId)
	if err != nil {
		// Jika tidak ada indikator, gunakan array kosong
		if err == sql.ErrNoRows {
			subKegiatan.Indikator = []domain.Indikator{}
			return helper.ToSubKegiatanResponse(subKegiatan), nil
		}
		log.Printf("Gagal mengambil indikator untuk subkegiatan %s: %v", subKegiatanId, err)
		return subkegiatan.SubKegiatanResponse{}, err
	}

	// Untuk setiap Indikator, ambil Target
	for i, indikator := range indikators {
		targets, err := service.subKegiatanRepository.FindTargetByIndikatorId(ctx, tx, indikator.Id)
		if err != nil {
			// Jika tidak ada target, gunakan array kosong
			if err == sql.ErrNoRows {
				indikators[i].Target = []domain.Target{}
				continue
			}
			log.Printf("Gagal mengambil target untuk indikator %s: %v", indikator.Id, err)
			return subkegiatan.SubKegiatanResponse{}, err
		}
		indikators[i].Target = targets
	}

	// Ambil data OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, subKegiatan.KodeOpd)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("OPD tidak ditemukan untuk kode: %s", subKegiatan.KodeOpd)
			subKegiatan.NamaOpd = ""
		} else {
			log.Printf("Error saat mencari OPD dengan kode %s: %v", subKegiatan.KodeOpd, err)
			return subkegiatan.SubKegiatanResponse{}, err
		}
	} else {
		subKegiatan.NamaOpd = opd.NamaOpd
	}

	// Gabungkan data
	subKegiatan.Indikator = indikators

	return helper.ToSubKegiatanResponse(subKegiatan), nil
}

func (service *SubKegiatanServiceImpl) FindAll(ctx context.Context, kodeOpd, rekinId, status string) ([]subkegiatan.SubKegiatanResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("Gagal memulai transaksi:", err)
		return []subkegiatan.SubKegiatanResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data SubKegiatan
	subKegiatans, err := service.subKegiatanRepository.FindAll(ctx, tx, kodeOpd, rekinId, status)
	if err != nil {
		log.Println("Gagal mencari data sub kegiatan:", err)
		return []subkegiatan.SubKegiatanResponse{}, err
	}

	// Untuk setiap SubKegiatan, ambil data Indikator dan Target
	for i, subKegiatan := range subKegiatans {
		// Tambah log untuk debug
		log.Printf("Mencari OPD dengan kode: %s", subKegiatan.KodeOpd)

		opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, subKegiatan.KodeOpd)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("OPD tidak ditemukan untuk kode: %s", subKegiatan.KodeOpd)
				subKegiatans[i].NamaOpd = ""
			} else {
				log.Printf("Error saat mencari OPD dengan kode %s: %v", subKegiatan.KodeOpd, err)
				return []subkegiatan.SubKegiatanResponse{}, err
			}
		} else {
			log.Printf("OPD ditemukan: %+v", opd) // Log seluruh data OPD
			subKegiatans[i].NamaOpd = opd.NamaOpd
		}

		// Ambil Indikator
		indikators, err := service.subKegiatanRepository.FindIndikatorBySubKegiatanId(ctx, tx, subKegiatan.Id)
		if err != nil {
			// Jika tidak ada indikator, lanjutkan dengan array kosong
			if err == sql.ErrNoRows {
				subKegiatans[i].Indikator = []domain.Indikator{}
				continue
			}
			log.Printf("Gagal mengambil indikator untuk subkegiatan %s: %v", subKegiatan.Id, err)
			return []subkegiatan.SubKegiatanResponse{}, err
		}

		// Untuk setiap Indikator, ambil Target
		for j, indikator := range indikators {
			targets, err := service.subKegiatanRepository.FindTargetByIndikatorId(ctx, tx, indikator.Id)
			if err != nil {
				// Jika tidak ada target, lanjutkan dengan array kosong
				if err == sql.ErrNoRows {
					indikators[j].Target = []domain.Target{}
					continue
				}
				log.Printf("Gagal mengambil target untuk indikator %s: %v", indikator.Id, err)
				return []subkegiatan.SubKegiatanResponse{}, err
			}
			indikators[j].Target = targets
		}

		subKegiatans[i].Indikator = indikators
	}

	return helper.ToSubKegiatanResponses(subKegiatans), nil
}

func (service *SubKegiatanServiceImpl) Delete(ctx context.Context, subKegiatanId string) error {
	// Validasi ID
	if subKegiatanId == "" {
		return errors.New("subkegiatan id tidak boleh kosong")
	}

	// Mulai transaksi
	tx, err := service.DB.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Proses delete
	err = service.subKegiatanRepository.Delete(ctx, tx, subKegiatanId)
	if err != nil {
		return fmt.Errorf("gagal menghapus sub kegiatan: %v", err)
	}

	return nil
}

func (service *SubKegiatanServiceImpl) CreateRekin(ctx context.Context, request subkegiatan.SubKegiatanCreateRekinRequest) ([]subkegiatan.SubKegiatanResponse, error) {
	// Konversi single ID ke array
	idSubKegiatanArray := []string{request.IdSubKegiatan}

	tx, err := service.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Cek apakah rencana kinerja dengan ID yang diberikan ada
	_, err = service.rencanaKinerjaRepoitory.FindById(ctx, tx, request.RekinId, "", "")
	if err != nil {
		return nil, fmt.Errorf("rencana kinerja dengan id %s tidak ditemukan: %v", request.RekinId, err)
	}

	var updatedSubKegiatans []domain.SubKegiatan

	// Proses setiap ID usulan
	for _, idSubKegiatan := range idSubKegiatanArray {
		// Cek apakah usulan dengan ID yang diberikan ada
		existingSubKegiatan, err := service.subKegiatanRepository.FindById(ctx, tx, idSubKegiatan)
		if err != nil {
			return nil, fmt.Errorf("subkegiatan dengan id %s tidak ditemukan: %v", idSubKegiatan, err)
		}

		// Cek apakah usulan sudah memiliki rekin_id
		if existingSubKegiatan.RekinId != "" {
			return nil, fmt.Errorf("subkegiatan dengan id %s sudah memiliki rencana kinerja", idSubKegiatan)
		}

		// Update rekin_id dan status
		err = service.subKegiatanRepository.CreateRekin(ctx, tx, idSubKegiatan, request.RekinId)
		if err != nil {
			return nil, fmt.Errorf("gagal mengupdate rekin untuk subkegiatan %s: %v", idSubKegiatan, err)
		}

		// Ambil data usulan yang sudah diupdate
		updatedSubKegiatan, err := service.subKegiatanRepository.FindById(ctx, tx, idSubKegiatan)
		if err != nil {
			return nil, fmt.Errorf("gagal mengambil data subkegiatan yang diupdate: %v", err)
		}

		// Ambil data OPD
		opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, updatedSubKegiatan.KodeOpd)
		if err == nil { // Hanya set jika berhasil mendapatkan data OPD
			updatedSubKegiatan.NamaOpd = opd.NamaOpd
		}

		updatedSubKegiatans = append(updatedSubKegiatans, updatedSubKegiatan)
	}

	// Konversi ke response
	responses := helper.ToSubKegiatanResponses(updatedSubKegiatans)
	return responses, nil
}

func (service *SubKegiatanServiceImpl) DeleteSubKegiatanTerpilih(ctx context.Context, idSubKegiatan string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	err = service.subKegiatanRepository.DeleteSubKegiatanTerpilih(ctx, tx, idSubKegiatan)
	if err != nil {
		return fmt.Errorf("gagal menghapus subkegiatan terpilih: %v", err)
	}

	return nil
}
