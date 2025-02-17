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
	PeriodeRepository   repository.PeriodeRepository
	DB                  *sql.DB
}

func NewTujuanOpdServiceImpl(tujuanOpdRepository repository.TujuanOpdRepository, opdRepository repository.OpdRepository, periodeRepository repository.PeriodeRepository, DB *sql.DB) *TujuanOpdServiceImpl {
	return &TujuanOpdServiceImpl{
		TujuanOpdRepository: tujuanOpdRepository,
		OpdRepository:       opdRepository,
		PeriodeRepository:   periodeRepository,
		DB:                  DB,
	}
}

func (service *TujuanOpdServiceImpl) Create(ctx context.Context, request tujuanopd.TujuanOpdCreateRequest) (tujuanopd.TujuanOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi periode
	periode, err := service.PeriodeRepository.FindById(ctx, tx, request.PeriodeId)
	if err != nil {
		if err == sql.ErrNoRows {
			return tujuanopd.TujuanOpdResponse{}, fmt.Errorf("periode dengan id %d tidak ditemukan", request.PeriodeId)
		}
		return tujuanopd.TujuanOpdResponse{}, err
	}

	// Convert tahun awal dan akhir ke integer untuk validasi
	tahunAwal, err := strconv.Atoi(periode.TahunAwal)
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, fmt.Errorf("format tahun awal periode tidak valid: %s", periode.TahunAwal)
	}

	tahunAkhir, err := strconv.Atoi(periode.TahunAkhir)
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, fmt.Errorf("format tahun akhir periode tidak valid: %s", periode.TahunAkhir)
	}

	tujuanOpdDomain := domain.TujuanOpd{
		KodeOpd:          request.KodeOpd,
		KodeBidangUrusan: request.KodeBidangUrusan,
		Tujuan:           request.Tujuan,
		RumusPerhitungan: request.RumusPerhitungan,
		SumberData:       request.SumberData,
		PeriodeId:        domain.Periode{Id: request.PeriodeId},
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

		// Map untuk mengecek duplikasi tahun
		tahunMap := make(map[string]bool)

		// Convert target request to domain
		for _, targetReq := range indikatorReq.Target {
			// Validasi format tahun target
			tahunTarget, err := strconv.Atoi(targetReq.Tahun)
			if err != nil {
				return tujuanopd.TujuanOpdResponse{}, fmt.Errorf("format tahun target tidak valid: %s", targetReq.Tahun)
			}

			// Validasi tahun target berada dalam range periode
			if tahunTarget < tahunAwal || tahunTarget > tahunAkhir {
				return tujuanopd.TujuanOpdResponse{}, fmt.Errorf(
					"tahun target %d harus berada dalam rentang periode %d-%d",
					tahunTarget,
					tahunAwal,
					tahunAkhir,
				)
			}

			// Validasi duplikasi tahun
			if tahunMap[targetReq.Tahun] {
				return tujuanopd.TujuanOpdResponse{}, fmt.Errorf(
					"tahun target %s duplikat dalam indikator yang sama",
					targetReq.Tahun,
				)
			}
			tahunMap[targetReq.Tahun] = true

			// Validasi target dan satuan tidak boleh kosong
			if targetReq.Target == "" {
				return tujuanopd.TujuanOpdResponse{}, fmt.Errorf(
					"target untuk tahun %s tidak boleh kosong",
					targetReq.Tahun,
				)
			}
			if targetReq.Satuan == "" {
				return tujuanopd.TujuanOpdResponse{}, fmt.Errorf(
					"satuan untuk tahun %s tidak boleh kosong",
					targetReq.Tahun,
				)
			}

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

		// Validasi minimal harus ada 1 target
		if len(indikatorReq.Target) == 0 {
			return tujuanopd.TujuanOpdResponse{}, fmt.Errorf(
				"indikator harus memiliki minimal 1 target dalam rentang periode %d-%d",
				tahunAwal,
				tahunAkhir,
			)
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

	// Validasi periode
	periode, err := service.PeriodeRepository.FindById(ctx, tx, request.PeriodeId)
	if err != nil {
		if err == sql.ErrNoRows {
			return tujuanopd.TujuanOpdResponse{}, fmt.Errorf("periode dengan id %d tidak ditemukan", request.PeriodeId)
		}
		return tujuanopd.TujuanOpdResponse{}, err
	}

	// Convert tahun awal dan akhir ke integer untuk validasi
	tahunAwal, err := strconv.Atoi(periode.TahunAwal)
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, fmt.Errorf("format tahun awal periode tidak valid: %s", periode.TahunAwal)
	}

	tahunAkhir, err := strconv.Atoi(periode.TahunAkhir)
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, fmt.Errorf("format tahun akhir periode tidak valid: %s", periode.TahunAkhir)
	}

	// Cek apakah data exists
	_, err = service.TujuanOpdRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, err
	}

	// Update data utama
	tujuanOpd := domain.TujuanOpd{
		Id:               request.Id,
		KodeOpd:          request.KodeOpd,
		KodeBidangUrusan: request.KodeBidangUrusan,
		Tujuan:           request.Tujuan,
		RumusPerhitungan: request.RumusPerhitungan,
		SumberData:       request.SumberData,
		PeriodeId:        domain.Periode{Id: request.PeriodeId},
	}

	// Convert indikator request to domain
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
		}

		// Map untuk mengecek duplikasi tahun
		tahunMap := make(map[string]bool)

		// Convert target request to domain
		for _, targetReq := range indikatorReq.Target {
			// Validasi format tahun target
			tahunTarget, err := strconv.Atoi(targetReq.Tahun)
			if err != nil {
				return tujuanopd.TujuanOpdResponse{}, fmt.Errorf("format tahun target tidak valid: %s", targetReq.Tahun)
			}

			// Validasi tahun target berada dalam range periode
			if tahunTarget < tahunAwal || tahunTarget > tahunAkhir {
				return tujuanopd.TujuanOpdResponse{}, fmt.Errorf(
					"tahun target %d harus berada dalam rentang periode %d-%d",
					tahunTarget,
					tahunAwal,
					tahunAkhir,
				)
			}

			// Validasi duplikasi tahun
			if tahunMap[targetReq.Tahun] {
				return tujuanopd.TujuanOpdResponse{}, fmt.Errorf(
					"tahun target %s duplikat dalam indikator yang sama",
					targetReq.Tahun,
				)
			}
			tahunMap[targetReq.Tahun] = true

			// Validasi target dan satuan tidak boleh kosong
			if targetReq.Target == "" {
				return tujuanopd.TujuanOpdResponse{}, fmt.Errorf(
					"target untuk tahun %s tidak boleh kosong",
					targetReq.Tahun,
				)
			}
			if targetReq.Satuan == "" {
				return tujuanopd.TujuanOpdResponse{}, fmt.Errorf(
					"satuan untuk tahun %s tidak boleh kosong",
					targetReq.Tahun,
				)
			}

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

		// Validasi minimal harus ada 1 target
		if len(indikatorReq.Target) == 0 {
			return tujuanopd.TujuanOpdResponse{}, fmt.Errorf(
				"indikator harus memiliki minimal 1 target dalam rentang periode %d-%d",
				tahunAwal,
				tahunAkhir,
			)
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

	tujuanOpd, err := service.TujuanOpdRepository.FindById(ctx, tx, tujuanOpdId)
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, err
	}

	// Ambil data OPD
	opd, err := service.OpdRepository.FindByKodeOpd(ctx, tx, tujuanOpd.KodeOpd)
	if err != nil {
		return tujuanopd.TujuanOpdResponse{}, err
	}

	response := tujuanopd.TujuanOpdResponse{
		Id:               tujuanOpd.Id,
		KodeBidangUrusan: tujuanOpd.KodeBidangUrusan,
		KodeOpd:          tujuanOpd.KodeOpd,
		NamaOpd:          opd.NamaOpd,
		Tujuan:           tujuanOpd.Tujuan,
		RumusPerhitungan: tujuanOpd.RumusPerhitungan,
		SumberData:       tujuanOpd.SumberData,
		Periode: tujuanopd.PeriodeResponse{
			Id:         tujuanOpd.PeriodeId.Id,
			TahunAwal:  tujuanOpd.PeriodeId.TahunAwal,
			TahunAkhir: tujuanOpd.PeriodeId.TahunAkhir,
		},
		Indikator: make([]tujuanopd.IndikatorResponse, 0),
	}

	for _, indikator := range tujuanOpd.Indikator {
		indikatorResponse := tujuanopd.IndikatorResponse{
			Id:            indikator.Id,
			IdTujuanOpd:   tujuanOpd.Id,
			NamaIndikator: indikator.Indikator,
			Target:        make([]tujuanopd.TargetResponse, 0),
		}

		for _, target := range indikator.Target {
			targetResponse := tujuanopd.TargetResponse{
				Id:              target.Id,
				IndikatorId:     target.IndikatorId,
				Tahun:           target.Tahun,
				TargetIndikator: target.Target,
				SatuanIndikator: target.Satuan,
			}
			indikatorResponse.Target = append(indikatorResponse.Target, targetResponse)
		}

		response.Indikator = append(response.Indikator, indikatorResponse)
	}

	return response, nil
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

	// Ambil data OPD
	opd, err := service.OpdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return nil, err
	}

	// Ambil semua tujuan OPD
	tujuanOpds, err := service.TujuanOpdRepository.FindAll(ctx, tx, kodeOpd, tahun)
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]tujuanopd.TujuanOpdResponse, 0), nil
		}
		return nil, err
	}

	var responses []tujuanopd.TujuanOpdResponse
	for _, tujuan := range tujuanOpds {
		tujuanResponse := tujuanopd.TujuanOpdResponse{
			Id:               tujuan.Id,
			KodeBidangUrusan: tujuan.KodeBidangUrusan,
			KodeOpd:          tujuan.KodeOpd,
			NamaOpd:          opd.NamaOpd,
			Tujuan:           tujuan.Tujuan,
			RumusPerhitungan: tujuan.RumusPerhitungan,
			SumberData:       tujuan.SumberData,
			Periode: tujuanopd.PeriodeResponse{
				Id:         tujuan.PeriodeId.Id,
				TahunAwal:  tujuan.PeriodeId.TahunAwal,
				TahunAkhir: tujuan.PeriodeId.TahunAkhir,
			},
			Indikator: make([]tujuanopd.IndikatorResponse, 0),
		}

		// Proses setiap indikator
		for _, indikator := range tujuan.Indikator {
			indikatorResponse := tujuanopd.IndikatorResponse{
				Id:            indikator.Id,
				IdTujuanOpd:   tujuan.Id,
				NamaIndikator: indikator.Indikator,
				Target:        make([]tujuanopd.TargetResponse, 0),
			}

			// Buat map untuk target yang ada
			targetMap := make(map[string]domain.Target)
			for _, t := range indikator.Target {
				targetMap[t.Tahun] = t
			}

			// Generate target untuk setiap tahun dalam periode
			tahunAwal, _ := strconv.Atoi(tujuan.PeriodeId.TahunAwal)
			tahunAkhir, _ := strconv.Atoi(tujuan.PeriodeId.TahunAkhir)

			for year := tahunAwal; year <= tahunAkhir; year++ {
				tahunStr := strconv.Itoa(year)

				if target, exists := targetMap[tahunStr]; exists && target.Id != "" {
					// Jika target ada dan memiliki ID
					targetResponse := tujuanopd.TargetResponse{
						Id:              target.Id,
						IndikatorId:     indikator.Id,
						Tahun:           tahunStr,
						TargetIndikator: target.Target,
						SatuanIndikator: target.Satuan,
					}
					indikatorResponse.Target = append(indikatorResponse.Target, targetResponse)
				} else {
					// Jika tidak ada target atau target tidak memiliki ID
					targetResponse := tujuanopd.TargetResponse{
						Id:              "",
						IndikatorId:     "",
						Tahun:           tahunStr,
						TargetIndikator: "",
						SatuanIndikator: "",
					}
					indikatorResponse.Target = append(indikatorResponse.Target, targetResponse)
				}
			}

			tujuanResponse.Indikator = append(tujuanResponse.Indikator, indikatorResponse)
		}

		responses = append(responses, tujuanResponse)
	}

	if len(responses) == 0 {
		responses = make([]tujuanopd.TujuanOpdResponse, 0)
	}

	return responses, nil
}
