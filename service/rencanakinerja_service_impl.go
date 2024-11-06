package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/opdmaster"
	"ekak_kabupaten_madiun/model/web/rencanakinerja"
	"ekak_kabupaten_madiun/repository"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type RencanaKinerjaServiceImpl struct {
	rencanaKinerjaRepository repository.RencanaKinerjaRepository
	DB                       *sql.DB
	Validate                 *validator.Validate
	opdRepository            repository.OpdRepository
}

func NewRencanaKinerjaServiceImpl(rencanaKinerjaRepository repository.RencanaKinerjaRepository, DB *sql.DB, validate *validator.Validate, opdRepository repository.OpdRepository) *RencanaKinerjaServiceImpl {
	return &RencanaKinerjaServiceImpl{
		rencanaKinerjaRepository: rencanaKinerjaRepository,
		DB:                       DB,
		Validate:                 validate,
		opdRepository:            opdRepository,
	}
}

func (service *RencanaKinerjaServiceImpl) Create(ctx context.Context, request rencanakinerja.RencanaKinerjaCreateRequest) (rencanakinerja.RencanaKinerjaResponse, error) {
	log.Println("Memulai proses Create RencanaKinerja")

	err := service.Validate.Struct(request)
	if err != nil {
		log.Printf("Validasi gagal: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("validasi gagal: %v", err)
	}

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Gagal memulai transaksi: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Perbaikan pengecekan kode OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, request.KodeOpd)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Kode OPD %s tidak ditemukan", request.KodeOpd)
			return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("kode OPD %s tidak ditemukan", request.KodeOpd)
		}
		log.Printf("Gagal memeriksa kode OPD: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal memeriksa kode OPD: %v", err)
	}

	if opd.KodeOpd == "" {
		log.Printf("Kode OPD %s tidak valid", request.KodeOpd)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("kode OPD %s tidak valid", request.KodeOpd)
	}

	randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
	customId := fmt.Sprintf("REKIN-PEG-%s", randomDigits)

	rencanaKinerja := domain.RencanaKinerja{
		Id:                   customId,
		NamaRencanaKinerja:   request.NamaRencanaKinerja,
		Tahun:                request.Tahun,
		StatusRencanaKinerja: request.StatusRencanaKinerja,
		Catatan:              request.Catatan,
		KodeOpd:              request.KodeOpd,
		PegawaiId:            request.PegawaiId,
		KodeSubKegiatan:      "",
		Indikator:            make([]domain.Indikator, len(request.Indikator)),
	}

	log.Printf("RencanaKinerja dibuat dengan ID: %s", customId)

	for i, indikatorRequest := range request.Indikator {
		indikatorRandomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
		indikatorId := fmt.Sprintf("IND-REKIN-%s", indikatorRandomDigits)
		indikator := domain.Indikator{
			Id:               indikatorId,
			Indikator:        indikatorRequest.NamaIndikator,
			Tahun:            request.Tahun,
			Target:           make([]domain.Target, len(indikatorRequest.Target)),
			RencanaKinerjaId: rencanaKinerja.Id,
		}

		if indikator.Indikator == "" {
			log.Printf("Indikator kosong ditemukan: %+v", indikator)
		}

		log.Printf("Indikator dibuat: %+v", indikator)

		for j, targetRequest := range indikatorRequest.Target {
			targetRandomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
			targetId := fmt.Sprintf("TRGT-IND-REKIN-%s", targetRandomDigits)
			target := domain.Target{
				Id:          targetId,
				Tahun:       request.Tahun,
				Target:      targetRequest.Target,
				Satuan:      targetRequest.SatuanIndikator,
				IndikatorId: indikator.Id,
			}
			indikator.Target[j] = target
			log.Printf("Target dibuat dengan ID: %s", targetId)
		}

		rencanaKinerja.Indikator[i] = indikator
	}

	log.Println("Memanggil repository.Create")
	rencanaKinerja, err = service.rencanaKinerjaRepository.Create(ctx, tx, rencanaKinerja)
	if err != nil {
		log.Printf("Gagal menyimpan RencanaKinerja: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal menyimpan RencanaKinerja: %v", err)
	}

	rencanaKinerja.NamaOpd = opd.NamaOpd

	log.Println("RencanaKinerja berhasil disimpan")
	response := helper.ToRencanaKinerjaResponse(rencanaKinerja)
	log.Printf("Response: %+v", response)

	return response, nil
}

func (service *RencanaKinerjaServiceImpl) Update(ctx context.Context, request rencanakinerja.RencanaKinerjaUpdateRequest) (rencanakinerja.RencanaKinerjaResponse, error) {
	log.Println("Memulai proses Update RencanaKinerja")

	err := service.Validate.Struct(request)
	if err != nil {
		log.Printf("Validasi gagal: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("validasi gagal: %v", err)
	}

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Gagal memulai transaksi: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	var rencanaKinerja domain.RencanaKinerja
	if request.Id != "" {
		rencanaKinerja, err = service.rencanaKinerjaRepository.FindById(ctx, tx, request.Id, "", "")
		if err != nil {
			log.Printf("Gagal menemukan RencanaKinerja: %v", err)
			return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal menemukan RencanaKinerja: %v", err)
		}
	} else {
		randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
		rencanaKinerja.Id = fmt.Sprintf("REKIN-PEG-%s", randomDigits)
		log.Printf("Membuat RencanaKinerja baru dengan ID: %s", rencanaKinerja.Id)
	}

	rencanaKinerja.NamaRencanaKinerja = request.NamaRencanaKinerja
	rencanaKinerja.Tahun = request.Tahun
	rencanaKinerja.StatusRencanaKinerja = request.StatusRencanaKinerja
	rencanaKinerja.Catatan = request.Catatan
	rencanaKinerja.KodeOpd = request.KodeOpd
	rencanaKinerja.PegawaiId = request.PegawaiId

	rencanaKinerja.Indikator = make([]domain.Indikator, len(request.Indikator))
	for i, indikatorRequest := range request.Indikator {
		var indikatorId string
		if indikatorRequest.Id != "" {
			indikatorId = indikatorRequest.Id
		} else {
			randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
			indikatorId = fmt.Sprintf("IND-REKIN-%s", randomDigits)
			log.Printf("Membuat Indikator baru dengan ID: %s", indikatorId)
		}

		indikator := domain.Indikator{
			Id:               indikatorId,
			Indikator:        indikatorRequest.Indikator,
			Tahun:            rencanaKinerja.Tahun,
			RencanaKinerjaId: rencanaKinerja.Id,
		}

		indikator.Target = make([]domain.Target, len(indikatorRequest.Target))
		for j, targetRequest := range indikatorRequest.Target {
			var targetId string
			if targetRequest.Id != "" {
				targetId = targetRequest.Id
			} else {
				randomDigits := fmt.Sprintf("%05d", uuid.New().ID()%100000)
				targetId = fmt.Sprintf("TRGT-IND-REKIN-%s", randomDigits)
				log.Printf("Membuat Target baru dengan ID: %s", targetId)
			}

			target := domain.Target{
				Id:          targetId,
				Tahun:       rencanaKinerja.Tahun,
				Target:      targetRequest.Target,
				Satuan:      targetRequest.SatuanIndikator,
				IndikatorId: indikator.Id,
			}
			indikator.Target[j] = target
		}

		rencanaKinerja.Indikator[i] = indikator
	}

	log.Println("Memanggil repository.Update")
	rencanaKinerja, err = service.rencanaKinerjaRepository.Update(ctx, tx, rencanaKinerja)
	if err != nil {
		log.Printf("Gagal memperbarui RencanaKinerja: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal memperbarui RencanaKinerja: %v", err)
	}

	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, rencanaKinerja.KodeOpd)
	if err != nil {
		log.Printf("Gagal mengambil data OPD: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal mengambil data OPD: %v", err)
	}

	// Set nama OPD ke dalam rencanaKinerja
	rencanaKinerja.NamaOpd = opd.NamaOpd
	log.Println("RencanaKinerja berhasil diperbarui")
	response := helper.ToRencanaKinerjaResponse(rencanaKinerja)
	log.Printf("Response: %+v", response)

	return response, nil
}

func (service *RencanaKinerjaServiceImpl) FindAll(ctx context.Context, pegawaiId string, kodeOPD string, tahun string) ([]rencanakinerja.RencanaKinerjaResponse, error) {
	log.Println("Memulai proses FindAll RencanaKinerja")

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Gagal memulai transaksi: %v", err)
		return nil, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	log.Printf("Mencari RencanaKinerja dengan pegawaiId: %s, kodeOPD: %s, tahun: %s", pegawaiId, kodeOPD, tahun)
	rencanaKinerjaList, err := service.rencanaKinerjaRepository.FindAll(ctx, tx, pegawaiId, kodeOPD, tahun)
	if err != nil {
		log.Printf("Gagal mencari RencanaKinerja: %v", err)
		return nil, fmt.Errorf("gagal mencari RencanaKinerja: %v", err)
	}
	log.Printf("Ditemukan %d RencanaKinerja", len(rencanaKinerjaList))

	var responses []rencanakinerja.RencanaKinerjaResponse
	for _, rencana := range rencanaKinerjaList {
		log.Printf("Memproses RencanaKinerja dengan ID: %s", rencana.Id)

		indikators, err := service.rencanaKinerjaRepository.FindIndikatorbyRekinId(ctx, tx, rencana.Id)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("Gagal mencari Indikator: %v", err)
			return nil, fmt.Errorf("gagal mencari Indikator: %v", err)
		}

		var indikatorResponses []rencanakinerja.IndikatorResponse
		for _, indikator := range indikators {
			targets, err := service.rencanaKinerjaRepository.FindTargetByIndikatorId(ctx, tx, indikator.Id)
			if err != nil && err != sql.ErrNoRows {
				log.Printf("Gagal mencari Target: %v", err)
				return nil, fmt.Errorf("gagal mencari Target: %v", err)
			}

			var targetResponses []rencanakinerja.TargetResponse
			for _, target := range targets {
				targetResponses = append(targetResponses, rencanakinerja.TargetResponse{
					Id:              target.Id,
					IndikatorId:     target.IndikatorId,
					TargetIndikator: target.Target,
					SatuanIndikator: target.Satuan,
				})
			}

			indikatorResponses = append(indikatorResponses, rencanakinerja.IndikatorResponse{
				Id:               indikator.Id,
				RencanaKinerjaId: indikator.RencanaKinerjaId,
				NamaIndikator:    indikator.Indikator,
				Target:           targetResponses,
			})
		}

		ActionButton := []web.ActionButton{
			{
				NameAction: "Find By Id Rencana Kinerja",
				Method:     "GET",
				Url:        "/detail-rencana_kinerja/:rencana_kinerja_id",
			},
			{
				NameAction: "Update Rencana Kinerja",
				Method:     "PUT",
				Url:        "/rencana_kinerja/update/:id",
			},
			{
				NameAction: "Delete Rencana Kinerja",
				Method:     "DELETE",
				Url:        "/rencana_kinerja/delete/:id",
			},
		}

		opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, rencana.KodeOpd)
		if err != nil {
			log.Printf("Gagal mencari OPD: %v", err)
			return nil, fmt.Errorf("gagal mencari OPD: %v", err)
		}

		responses = append(responses, rencanakinerja.RencanaKinerjaResponse{
			Id:                   rencana.Id,
			NamaRencanaKinerja:   rencana.NamaRencanaKinerja,
			Tahun:                rencana.Tahun,
			StatusRencanaKinerja: rencana.StatusRencanaKinerja,
			Catatan:              rencana.Catatan,
			KodeOpd: opdmaster.OpdResponseForAll{
				KodeOpd: opd.KodeOpd,
				NamaOpd: opd.NamaOpd,
			},
			PegawaiId: rencana.PegawaiId,
			Indikator: indikatorResponses,
			Action:    ActionButton,
		})
		log.Printf("RencanaKinerja Response ditambahkan untuk ID: %s", rencana.Id)
	}

	return responses, nil
}

func (service *RencanaKinerjaServiceImpl) FindById(ctx context.Context, id string, kodeOPD string, tahun string) (rencanakinerja.RencanaKinerjaResponse, error) {
	log.Printf("Mencari RencanaKinerja dengan ID: %s, KodeOPD: %s, Tahun: %s", id, kodeOPD, tahun)

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Gagal memulai transaksi: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	rencanaKinerja, err := service.rencanaKinerjaRepository.FindById(ctx, tx, id, kodeOPD, tahun)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("RencanaKinerja tidak ditemukan untuk ID: %s", id)
			return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("rencana kinerja tidak ditemukan")
		}
		log.Printf("Gagal menemukan rencana kinerja: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal menemukan rencana kinerja: %v", err)
	}

	log.Printf("RencanaKinerja ditemukan: %+v", rencanaKinerja)

	indikators, err := service.rencanaKinerjaRepository.FindIndikatorbyRekinId(ctx, tx, rencanaKinerja.Id)
	if err != nil {
		log.Printf("Gagal menemukan indikator: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal menemukan indikator: %v", err)
	}
	rencanaKinerja.Indikator = indikators

	log.Printf("Jumlah indikator ditemukan: %d", len(indikators))

	for i, indikator := range rencanaKinerja.Indikator {
		targets, err := service.rencanaKinerjaRepository.FindTargetByIndikatorId(ctx, tx, indikator.Id)
		if err != nil {
			log.Printf("Gagal menemukan target untuk indikator %s: %v", indikator.Id, err)
			return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal menemukan target untuk indikator %s: %v", indikator.Id, err)
		}
		rencanaKinerja.Indikator[i].Target = targets
		log.Printf("Jumlah target ditemukan untuk indikator %s: %d", indikator.Id, len(targets))
	}

	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, rencanaKinerja.KodeOpd)
	if err != nil {
		log.Printf("Gagal mengambil data OPD: %v", err)
		return rencanakinerja.RencanaKinerjaResponse{}, fmt.Errorf("gagal mengambil data OPD: %v", err)
	}

	// Set nama OPD ke dalam rencanaKinerja
	rencanaKinerja.NamaOpd = opd.NamaOpd

	response := helper.ToRencanaKinerjaResponse(rencanaKinerja)
	log.Printf("Response: %+v", response)

	return response, nil
}

func (service *RencanaKinerjaServiceImpl) Delete(ctx context.Context, id string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	rencanaKinerja, err := service.rencanaKinerjaRepository.FindById(ctx, tx, id, "", "")
	if err != nil {
		return err
	}

	return service.rencanaKinerjaRepository.Delete(ctx, tx, rencanaKinerja.Id)
}
