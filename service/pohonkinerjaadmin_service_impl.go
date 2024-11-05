package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"ekak_kabupaten_madiun/repository"

	"github.com/google/uuid"
)

type PohonKinerjaAdminServiceImpl struct {
	pohonKinerjaRepository repository.PohonKinerjaRepository
	DB                     *sql.DB
}

func NewPohonKinerjaAdminServiceImpl(pohonKinerjaRepository repository.PohonKinerjaRepository, DB *sql.DB) *PohonKinerjaAdminServiceImpl {
	return &PohonKinerjaAdminServiceImpl{
		pohonKinerjaRepository: pohonKinerjaRepository,
		DB:                     DB,
	}
}

func (service *PohonKinerjaAdminServiceImpl) Create(ctx context.Context, request pohonkinerja.PohonKinerjaAdminCreateRequest) (pohonkinerja.PohonKinerjaAdminResponseData, error) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// Persiapkan data indikator dan target
	var indikators []domain.Indikator
	for _, ind := range request.Indikator {
		indikatorId := "IND-POKIN-" + uuid.New().String()

		var targets []domain.Target
		for _, t := range ind.Target {
			targetId := "TRGT-IND-POKIN-" + uuid.New().String()
			target := domain.Target{
				Id:          targetId,
				IndikatorId: indikatorId,
				Target:      t.Target,
				Satuan:      t.Satuan,
				Tahun:       request.Tahun,
			}
			targets = append(targets, target)
		}

		indikator := domain.Indikator{
			Id:        indikatorId,
			Indikator: ind.NamaIndikator,
			Tahun:     request.Tahun,
			Target:    targets,
		}
		indikators = append(indikators, indikator)
	}

	pohonKinerja := domain.PohonKinerja{
		Parent:     request.Parent,
		NamaPohon:  request.NamaPohon,
		JenisPohon: request.JenisPohon,
		LevelPohon: request.LevelPohon,
		KodeOpd:    request.KodeOpd,
		Keterangan: request.Keterangan,
		Tahun:      request.Tahun,
		Indikator:  indikators,
	}

	result, err := service.pohonKinerjaRepository.CreatePokinAdmin(ctx, tx, pohonKinerja)
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponseData{}, err
	}

	// Konversi indikator domain ke IndikatorResponse
	var indikatorResponses []pohonkinerja.IndikatorResponse
	for _, ind := range result.Indikator {
		var targetResponses []pohonkinerja.TargetResponse
		for _, t := range ind.Target {
			targetResponse := pohonkinerja.TargetResponse{
				Id:              t.Id,
				IndikatorId:     t.IndikatorId,
				TargetIndikator: t.Target,
				SatuanIndikator: t.Satuan,
			}
			targetResponses = append(targetResponses, targetResponse)
		}

		indikatorResponse := pohonkinerja.IndikatorResponse{
			Id:            ind.Id,
			NamaIndikator: ind.Indikator,
			Target:        targetResponses,
		}
		indikatorResponses = append(indikatorResponses, indikatorResponse)
	}

	response := pohonkinerja.PohonKinerjaAdminResponseData{
		Id:         result.Id,
		Parent:     result.Parent,
		NamaPohon:  result.NamaPohon,
		JenisPohon: result.JenisPohon,
		LevelPohon: result.LevelPohon,
		KodeOpd:    result.KodeOpd,
		Keterangan: result.Keterangan,
		Tahun:      result.Tahun,
		Indikators: indikatorResponses,
	}

	return response, nil
}

func (service *PohonKinerjaAdminServiceImpl) Update(ctx context.Context, request pohonkinerja.PohonKinerjaAdminUpdateRequest) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Cek apakah data exists
	_, err = service.pohonKinerjaRepository.FindPokinAdminById(ctx, tx, request.Id)
	if err != nil {
		return err
	}

	// Persiapkan data indikator dan target
	var indikators []domain.Indikator
	for _, ind := range request.Indikator {
		indikatorId := ""
		if ind.Id != "" {
			indikatorId = ind.Id
		} else {
			indikatorId = "IND-POKIN-" + uuid.New().String()[:4]
		}

		var targets []domain.Target
		for _, t := range ind.Target {
			targetId := ""
			if t.Id != "" {
				targetId = t.Id
			} else {
				targetId = "TRGT-IND-POKIN-" + uuid.New().String()[:4]
			}

			target := domain.Target{
				Id:          targetId,
				IndikatorId: indikatorId,
				Target:      t.Target,
				Satuan:      t.Satuan,
				Tahun:       request.Tahun,
			}
			targets = append(targets, target)
		}

		indikator := domain.Indikator{
			Id:        indikatorId,
			Indikator: ind.NamaIndikator,
			Tahun:     request.Tahun,
			Target:    targets,
		}
		indikators = append(indikators, indikator)
	}

	pohonKinerja := domain.PohonKinerja{
		Id:         request.Id,
		Parent:     request.Parent,
		NamaPohon:  request.NamaPohon,
		JenisPohon: request.JenisPohon,
		LevelPohon: request.LevelPohon,
		KodeOpd:    request.KodeOpd,
		Keterangan: request.Keterangan,
		Tahun:      request.Tahun,
		Indikator:  indikators,
	}

	_, err = service.pohonKinerjaRepository.UpdatePokinAdmin(ctx, tx, pohonKinerja)
	if err != nil {
		return err
	}

	return nil
}

func (service *PohonKinerjaAdminServiceImpl) Delete(ctx context.Context, id int) error {
	// Mulai transaksi
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Cek apakah data exists sebelum dihapus
	_, err = service.pohonKinerjaRepository.FindPokinAdminById(ctx, tx, id)
	if err != nil {
		return err
	}

	// Lakukan penghapusan
	err = service.pohonKinerjaRepository.DeletePokinAdmin(ctx, tx, id)
	if err != nil {
		return err
	}

	return nil
}

func (service *PohonKinerjaAdminServiceImpl) FindById(ctx context.Context, id int) (pohonkinerja.PohonKinerjaAdminResponse, error) {
	return pohonkinerja.PohonKinerjaAdminResponse{}, nil
}

func (service *PohonKinerjaAdminServiceImpl) FindAll(ctx context.Context, tahun string) (pohonkinerja.PohonKinerjaAdminResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil semua data pohon kinerja
	pokins, err := service.pohonKinerjaRepository.FindPokinAdminAll(ctx, tx, tahun)
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponse{}, err
	}

	// Buat map untuk menyimpan data berdasarkan level dan parent
	pohonMap := make(map[int]map[int][]domain.PohonKinerja)
	for i := 0; i <= 6; i++ {
		pohonMap[i] = make(map[int][]domain.PohonKinerja)
	}

	// Kelompokkan data berdasarkan level dan parent
	for _, pokin := range pokins {
		pohonMap[pokin.LevelPohon][pokin.Parent] = append(pohonMap[pokin.LevelPohon][pokin.Parent], pokin)
	}

	// Bangun response dimulai dari Tematik (level 0)
	var tematiks []pohonkinerja.TematikResponse
	for _, tematik := range pohonMap[0][0] {
		tematikResp := service.buildTematikResponse(pohonMap, tematik)
		tematiks = append(tematiks, tematikResp)
	}

	return pohonkinerja.PohonKinerjaAdminResponse{
		Tahun:   tahun,
		Tematik: tematiks,
	}, nil
}

func (service *PohonKinerjaAdminServiceImpl) buildTematikResponse(pohonMap map[int]map[int][]domain.PohonKinerja, tematik domain.PohonKinerja) pohonkinerja.TematikResponse {
	tematikResp := pohonkinerja.TematikResponse{
		Id:         tematik.Id,
		Parent:     nil,
		Tema:       tematik.NamaPohon,
		Keterangan: tematik.Keterangan,
		Indikators: convertToIndikatorResponses(tematik.Indikator),
	}

	// Cek dan tambahkan subtematik jika ada
	if subTematiks := pohonMap[1][tematik.Id]; len(subTematiks) > 0 {
		var subTematikResponses []pohonkinerja.SubtematikResponse
		for _, subTematik := range subTematiks {
			subTematikResp := service.buildSubTematikResponse(pohonMap, subTematik)
			subTematikResponses = append(subTematikResponses, subTematikResp)
		}
		tematikResp.SubTematiks = subTematikResponses
	}

	return tematikResp
}

func (service *PohonKinerjaAdminServiceImpl) buildSubTematikResponse(pohonMap map[int]map[int][]domain.PohonKinerja, subTematik domain.PohonKinerja) pohonkinerja.SubtematikResponse {
	subTematikResp := pohonkinerja.SubtematikResponse{
		Id:         subTematik.Id,
		Parent:     subTematik.Parent,
		Tema:       subTematik.NamaPohon,
		Keterangan: subTematik.Keterangan,
		Indikators: convertToIndikatorResponses(subTematik.Indikator),
	}

	// Cek dan tambahkan subsubtematik jika ada
	if subSubTematiks := pohonMap[2][subTematik.Id]; len(subSubTematiks) > 0 {
		var subSubTematikResponses []pohonkinerja.SubSubTematikResponse
		for _, subSubTematik := range subSubTematiks {
			subSubTematikResp := service.buildSubSubTematikResponse(pohonMap, subSubTematik)
			subSubTematikResponses = append(subSubTematikResponses, subSubTematikResp)
		}
		subTematikResp.SubSubTematiks = subSubTematikResponses
	} else {
		// Jika tidak ada subsubtematik, cek strategic langsung
		if strategics := pohonMap[4][subTematik.Id]; len(strategics) > 0 {
			subTematikResp.Strategics = service.buildStrategicResponses(pohonMap, strategics)
		}
	}

	return subTematikResp
}

func (service *PohonKinerjaAdminServiceImpl) buildSubSubTematikResponse(pohonMap map[int]map[int][]domain.PohonKinerja, subSubTematik domain.PohonKinerja) pohonkinerja.SubSubTematikResponse {
	subSubTematikResp := pohonkinerja.SubSubTematikResponse{
		Id:         subSubTematik.Id,
		Parent:     subSubTematik.Parent,
		Tema:       subSubTematik.NamaPohon,
		Keterangan: subSubTematik.Keterangan,
		Indikators: convertToIndikatorResponses(subSubTematik.Indikator),
	}

	// Cek dan tambahkan supersubtematik jika ada
	if superSubTematiks := pohonMap[3][subSubTematik.Id]; len(superSubTematiks) > 0 {
		var superSubTematikResponses []pohonkinerja.SuperSubTematikResponse
		for _, superSubTematik := range superSubTematiks {
			superSubTematikResp := service.buildSuperSubTematikResponse(pohonMap, superSubTematik)
			superSubTematikResponses = append(superSubTematikResponses, superSubTematikResp)
		}
		subSubTematikResp.SuperSubTematiks = superSubTematikResponses
	} else {
		// Jika tidak ada supersubtematik, cek strategic langsung
		if strategics := pohonMap[4][subSubTematik.Id]; len(strategics) > 0 {
			subSubTematikResp.Strategics = service.buildStrategicResponses(pohonMap, strategics)
		}
	}

	return subSubTematikResp
}

func (service *PohonKinerjaAdminServiceImpl) buildSuperSubTematikResponse(pohonMap map[int]map[int][]domain.PohonKinerja, superSubTematik domain.PohonKinerja) pohonkinerja.SuperSubTematikResponse {
	superSubTematikResp := pohonkinerja.SuperSubTematikResponse{
		Id:         superSubTematik.Id,
		Parent:     superSubTematik.Parent,
		Tema:       superSubTematik.NamaPohon,
		Keterangan: superSubTematik.Keterangan,
		Indikators: convertToIndikatorResponses(superSubTematik.Indikator),
	}

	// Cek dan tambahkan strategic
	if strategics := pohonMap[4][superSubTematik.Id]; len(strategics) > 0 {
		superSubTematikResp.Strategics = service.buildStrategicResponses(pohonMap, strategics)
	}

	return superSubTematikResp
}

func (service *PohonKinerjaAdminServiceImpl) buildStrategicResponses(pohonMap map[int]map[int][]domain.PohonKinerja, strategics []domain.PohonKinerja) []pohonkinerja.StrategicResponse {
	var responses []pohonkinerja.StrategicResponse
	for _, strategic := range strategics {
		strategicResp := pohonkinerja.StrategicResponse{
			Id:              strategic.Id,
			Parent:          strategic.Parent,
			Strategi:        strategic.NamaPohon,
			Keterangan:      strategic.Keterangan,
			KodeOpd:         strategic.KodeOpd,
			PerangkatDaerah: "", // Sesuaikan dengan data yang tersedia
			Indikators:      convertToIndikatorResponses(strategic.Indikator),
		}

		// Cek dan tambahkan tactical
		if tacticals := pohonMap[5][strategic.Id]; len(tacticals) > 0 {
			strategicResp.Tacticals = service.buildTacticalResponses(pohonMap, tacticals)
		}

		responses = append(responses, strategicResp)
	}
	return responses
}

func (service *PohonKinerjaAdminServiceImpl) buildTacticalResponses(pohonMap map[int]map[int][]domain.PohonKinerja, tacticals []domain.PohonKinerja) []pohonkinerja.TacticalResponse {
	var responses []pohonkinerja.TacticalResponse
	for _, tactical := range tacticals {
		keterangan := &tactical.Keterangan
		tacticalResp := pohonkinerja.TacticalResponse{
			Id:              tactical.Id,
			Parent:          tactical.Parent,
			Strategi:        tactical.NamaPohon,
			Keterangan:      keterangan,
			KodeOpd:         tactical.KodeOpd,
			PerangkatDaerah: "", // Sesuaikan dengan data yang tersedia
			Indikators:      convertToIndikatorResponses(tactical.Indikator),
		}

		// Cek dan tambahkan operational
		if operationals := pohonMap[6][tactical.Id]; len(operationals) > 0 {
			tacticalResp.Operationals = service.buildOperationalResponses(operationals)
		}

		responses = append(responses, tacticalResp)
	}
	return responses
}

func (service *PohonKinerjaAdminServiceImpl) buildOperationalResponses(operationals []domain.PohonKinerja) []pohonkinerja.OperationalResponse {
	var responses []pohonkinerja.OperationalResponse
	for _, operational := range operationals {
		keterangan := &operational.Keterangan
		operationalResp := pohonkinerja.OperationalResponse{
			Id:              operational.Id,
			Parent:          operational.Parent,
			Strategi:        operational.NamaPohon,
			Keterangan:      keterangan,
			KodeOpd:         operational.KodeOpd,
			PerangkatDaerah: "", // Sesuaikan dengan data yang tersedia
			Indikators:      convertToIndikatorResponses(operational.Indikator),
		}
		responses = append(responses, operationalResp)
	}
	return responses
}

func convertToIndikatorResponses(indikators []domain.Indikator) []pohonkinerja.IndikatorResponse {
	var responses []pohonkinerja.IndikatorResponse
	for _, indikator := range indikators {
		var targetResponses []pohonkinerja.TargetResponse
		for _, target := range indikator.Target {
			targetResp := pohonkinerja.TargetResponse{
				Id:              target.Id,
				IndikatorId:     target.IndikatorId,
				TargetIndikator: target.Target,
				SatuanIndikator: target.Satuan,
			}
			targetResponses = append(targetResponses, targetResp)
		}

		indikatorResp := pohonkinerja.IndikatorResponse{
			Id:            indikator.Id,
			IdPokin:       indikator.PokinId,
			NamaIndikator: indikator.Indikator,
			Target:        targetResponses,
		}
		responses = append(responses, indikatorResp)
	}
	return responses
}

// Fungsi helper untuk debug
