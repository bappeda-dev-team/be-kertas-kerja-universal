package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"ekak_kabupaten_madiun/repository"
	"fmt"

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
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data pohon kinerja berdasarkan ID
	pokin, err := service.pohonKinerjaRepository.FindPokinAdminById(ctx, tx, id)
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponse{}, err
	}

	// Ambil semua data terkait untuk membangun hierarki
	allPokins, err := service.pohonKinerjaRepository.FindPokinAdminAll(ctx, tx, pokin.Tahun)
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponse{}, err
	}

	// Buat map untuk menyimpan pohon berdasarkan level dan parent
	pohonMap := make(map[int]map[int][]domain.PohonKinerja)
	for i := 1; i <= 7; i++ {
		pohonMap[i] = make(map[int][]domain.PohonKinerja)
	}

	for _, p := range allPokins {
		pohonMap[p.LevelPohon][p.Parent] = append(pohonMap[p.LevelPohon][p.Parent], p)
	}

	var tematiks []pohonkinerja.TematikResponse

	// Tentukan level dari data yang dicari
	switch pokin.LevelPohon {
	case 1: // Tematik
		tematikResp := convertToTematikResponse(pokin)
		// Cari sub-tematik
		var subTematiks []pohonkinerja.SubtematikResponse
		for _, subTematik := range pohonMap[2][pokin.Id] {
			subTematikResp := convertToSubTematikResponse(subTematik)
			// Cari sub-sub-tematik
			subTematikResp.SubSubTematiks = findSubSubTematiks(pohonMap, subTematik.Id)
			subTematiks = append(subTematiks, subTematikResp)
		}
		tematikResp.SubTematiks = subTematiks
		tematiks = append(tematiks, tematikResp)

	case 2: // Sub Tematik
		// Cari tematik parent
		for _, tematik := range pohonMap[1][pokin.Parent] {
			tematikResp := convertToTematikResponse(tematik)
			var subTematiks []pohonkinerja.SubtematikResponse
			subTematikResp := convertToSubTematikResponse(pokin)
			subTematikResp.SubSubTematiks = findSubSubTematiks(pohonMap, pokin.Id)
			subTematiks = append(subTematiks, subTematikResp)
			tematikResp.SubTematiks = subTematiks
			tematiks = append(tematiks, tematikResp)
		}

		// ... tambahkan case untuk level lainnya ...

	case 3: // Sub Sub Tematik
		// Cari tematik dan sub-tematik parent
		for _, subTematik := range pohonMap[2][pokin.Parent] {
			for _, tematik := range pohonMap[1][subTematik.Parent] {
				tematikResp := convertToTematikResponse(tematik)
				var subTematiks []pohonkinerja.SubtematikResponse

				subTematikResp := convertToSubTematikResponse(subTematik)
				var subSubTematiks []pohonkinerja.SubSubTematikResponse

				subSubTematikResp := convertToSubSubTematikResponse(pokin)
				subSubTematikResp.SuperSubTematiks = findSuperSubTematiks(pohonMap, pokin.Id)

				subSubTematiks = append(subSubTematiks, subSubTematikResp)
				subTematikResp.SubSubTematiks = subSubTematiks
				subTematiks = append(subTematiks, subTematikResp)
				tematikResp.SubTematiks = subTematiks
				tematiks = append(tematiks, tematikResp)
			}
		}

	case 4: // Super Sub Tematik
		for _, subSubTematik := range pohonMap[3][pokin.Parent] {
			for _, subTematik := range pohonMap[2][subSubTematik.Parent] {
				for _, tematik := range pohonMap[1][subTematik.Parent] {
					tematikResp := convertToTematikResponse(tematik)
					var subTematiks []pohonkinerja.SubtematikResponse

					subTematikResp := convertToSubTematikResponse(subTematik)
					var subSubTematiks []pohonkinerja.SubSubTematikResponse

					subSubTematikResp := convertToSubSubTematikResponse(subSubTematik)
					var superSubTematiks []pohonkinerja.SuperSubTematikResponse

					superSubTematikResp := convertToSuperSubTematikResponse(pokin)
					superSubTematikResp.Strategics = findStrategics(pohonMap, pokin.Id)

					superSubTematiks = append(superSubTematiks, superSubTematikResp)
					subSubTematikResp.SuperSubTematiks = superSubTematiks
					subSubTematiks = append(subSubTematiks, subSubTematikResp)
					subTematikResp.SubSubTematiks = subSubTematiks
					subTematiks = append(subTematiks, subTematikResp)
					tematikResp.SubTematiks = subTematiks
					tematiks = append(tematiks, tematikResp)
				}
			}
		}

	case 5: // Strategic
		for _, superSubTematik := range pohonMap[4][pokin.Parent] {
			for _, subSubTematik := range pohonMap[3][superSubTematik.Parent] {
				for _, subTematik := range pohonMap[2][subSubTematik.Parent] {
					for _, tematik := range pohonMap[1][subTematik.Parent] {
						tematikResp := buildCompleteHierarchy(tematik, subTematik, subSubTematik, superSubTematik, pokin, pohonMap)
						tematiks = append(tematiks, tematikResp)
					}
				}
			}
		}

	case 6: // Tactical
		for _, strategic := range pohonMap[5][pokin.Parent] {
			for _, superSubTematik := range pohonMap[4][strategic.Parent] {
				for _, subSubTematik := range pohonMap[3][superSubTematik.Parent] {
					for _, subTematik := range pohonMap[2][subSubTematik.Parent] {
						for _, tematik := range pohonMap[1][subTematik.Parent] {
							tematikResp := buildCompleteHierarchy(tematik, subTematik, subSubTematik, superSubTematik, strategic, pohonMap)
							tematiks = append(tematiks, tematikResp)
						}
					}
				}
			}
		}

	case 7: // Operational
		for _, tactical := range pohonMap[6][pokin.Parent] {
			for _, strategic := range pohonMap[5][tactical.Parent] {
				for _, superSubTematik := range pohonMap[4][strategic.Parent] {
					for _, subSubTematik := range pohonMap[3][superSubTematik.Parent] {
						for _, subTematik := range pohonMap[2][subSubTematik.Parent] {
							for _, tematik := range pohonMap[1][subTematik.Parent] {
								tematikResp := buildCompleteHierarchy(tematik, subTematik, subSubTematik, superSubTematik, strategic, pohonMap)
								tematiks = append(tematiks, tematikResp)
							}
						}
					}
				}
			}
		}
	}

	return pohonkinerja.PohonKinerjaAdminResponse{
		Tahun:   pokin.Tahun,
		Tematik: tematiks,
	}, nil
}

// Helper function untuk membangun hierarki lengkap
func buildCompleteHierarchy(tematik, subTematik, subSubTematik, superSubTematik, currentNode domain.PohonKinerja, pohonMap map[int]map[int][]domain.PohonKinerja) pohonkinerja.TematikResponse {
	tematikResp := convertToTematikResponse(tematik)
	subTematikResp := convertToSubTematikResponse(subTematik)
	subSubTematikResp := convertToSubSubTematikResponse(subSubTematik)
	superSubTematikResp := convertToSuperSubTematikResponse(superSubTematik)

	// Sesuaikan dengan level dari currentNode
	switch currentNode.LevelPohon {
	case 5:
		strategicResp := convertToStrategicResponse(currentNode)
		strategicResp.Tacticals = findTacticals(pohonMap, currentNode.Id)
		superSubTematikResp.Strategics = []pohonkinerja.StrategicResponse{strategicResp}
	case 6:
		tacticalResp := convertToTacticalResponse(currentNode)
		tacticalResp.Operationals = findOperationals(pohonMap, currentNode.Id)
		strategicResp := convertToStrategicResponse(pohonMap[5][currentNode.Parent][0])
		strategicResp.Tacticals = []pohonkinerja.TacticalResponse{tacticalResp}
		superSubTematikResp.Strategics = []pohonkinerja.StrategicResponse{strategicResp}
	}

	subSubTematikResp.SuperSubTematiks = []pohonkinerja.SuperSubTematikResponse{superSubTematikResp}
	subTematikResp.SubSubTematiks = []pohonkinerja.SubSubTematikResponse{subSubTematikResp}
	tematikResp.SubTematiks = []pohonkinerja.SubtematikResponse{subTematikResp}

	return tematikResp
}

// Helper functions untuk mencari data di setiap level
func findSubSubTematiks(pohonMap map[int]map[int][]domain.PohonKinerja, parentId int) []pohonkinerja.SubSubTematikResponse {
	var subSubTematiks []pohonkinerja.SubSubTematikResponse
	for _, subSubTematik := range pohonMap[3][parentId] {
		subSubTematikResp := convertToSubSubTematikResponse(subSubTematik)
		subSubTematikResp.SuperSubTematiks = findSuperSubTematiks(pohonMap, subSubTematik.Id)
		subSubTematiks = append(subSubTematiks, subSubTematikResp)
	}
	return subSubTematiks
}

func findSuperSubTematiks(pohonMap map[int]map[int][]domain.PohonKinerja, parentId int) []pohonkinerja.SuperSubTematikResponse {
	var superSubTematiks []pohonkinerja.SuperSubTematikResponse
	for _, superSubTematik := range pohonMap[4][parentId] {
		superSubTematikResp := convertToSuperSubTematikResponse(superSubTematik)
		superSubTematikResp.Strategics = findStrategics(pohonMap, superSubTematik.Id)
		superSubTematiks = append(superSubTematiks, superSubTematikResp)
	}
	return superSubTematiks
}

func findStrategics(pohonMap map[int]map[int][]domain.PohonKinerja, parentId int) []pohonkinerja.StrategicResponse {
	var strategics []pohonkinerja.StrategicResponse
	for _, strategic := range pohonMap[5][parentId] {
		strategicResp := convertToStrategicResponse(strategic)
		strategicResp.Tacticals = findTacticals(pohonMap, strategic.Id)
		strategics = append(strategics, strategicResp)
	}
	return strategics
}

func findTacticals(pohonMap map[int]map[int][]domain.PohonKinerja, parentId int) []pohonkinerja.TacticalResponse {
	var tacticals []pohonkinerja.TacticalResponse
	for _, tactical := range pohonMap[6][parentId] {
		tacticalResp := convertToTacticalResponse(tactical)
		tacticalResp.Operationals = findOperationals(pohonMap, tactical.Id)
		tacticals = append(tacticals, tacticalResp)
	}
	return tacticals
}

func findOperationals(pohonMap map[int]map[int][]domain.PohonKinerja, parentId int) []pohonkinerja.OperationalResponse {
	var operationals []pohonkinerja.OperationalResponse
	for _, operational := range pohonMap[7][parentId] {
		operationalResp := convertToOperationalResponse(operational)
		operationals = append(operationals, operationalResp)
	}
	return operationals
}

func (service *PohonKinerjaAdminServiceImpl) FindAll(ctx context.Context, tahun string) (pohonkinerja.PohonKinerjaAdminResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	pokins, err := service.pohonKinerjaRepository.FindPokinAdminAll(ctx, tx, tahun)
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponse{}, err
	}

	// Buat map untuk menyimpan pohon berdasarkan level dan parent
	pohonMap := make(map[int]map[int][]domain.PohonKinerja)
	for i := 1; i <= 7; i++ {
		pohonMap[i] = make(map[int][]domain.PohonKinerja)
	}

	for _, pokin := range pokins {
		pohonMap[pokin.LevelPohon][pokin.Parent] = append(pohonMap[pokin.LevelPohon][pokin.Parent], pokin)
	}

	// Mulai dari level tematik (level 1)
	var tematiks []pohonkinerja.TematikResponse
	for _, tematik := range pohonMap[1][0] {
		tematikResp := convertToTematikResponse(tematik)

		// Cari sub-tematik (level 2)
		var subTematiks []pohonkinerja.SubtematikResponse
		for _, subTematik := range pohonMap[2][tematik.Id] {
			subTematikResp := convertToSubTematikResponse(subTematik)

			// Cari sub-sub-tematik (level 3)
			var subSubTematiks []pohonkinerja.SubSubTematikResponse
			for _, subSubTematik := range pohonMap[3][subTematik.Id] {
				subSubTematikResp := convertToSubSubTematikResponse(subSubTematik)

				// Cari super-sub-tematik (level 4)
				var superSubTematiks []pohonkinerja.SuperSubTematikResponse
				for _, superSubTematik := range pohonMap[4][subSubTematik.Id] {
					superSubTematikResp := convertToSuperSubTematikResponse(superSubTematik)

					// Cari strategic (level 5)
					var strategics []pohonkinerja.StrategicResponse
					for _, strategic := range pohonMap[5][superSubTematik.Id] {
						strategicResp := convertToStrategicResponse(strategic)

						// Cari tactical (level 6)
						var tacticals []pohonkinerja.TacticalResponse
						for _, tactical := range pohonMap[6][strategic.Id] {
							tacticalResp := convertToTacticalResponse(tactical)

							// Cari operational (level 7)
							var operationals []pohonkinerja.OperationalResponse
							for _, operational := range pohonMap[7][tactical.Id] {
								operationalResp := convertToOperationalResponse(operational)
								operationals = append(operationals, operationalResp)
							}

							tacticalResp.Operationals = operationals
							tacticals = append(tacticals, tacticalResp)
						}

						strategicResp.Tacticals = tacticals
						strategics = append(strategics, strategicResp)
					}

					superSubTematikResp.Strategics = strategics
					superSubTematiks = append(superSubTematiks, superSubTematikResp)
				}

				subSubTematikResp.SuperSubTematiks = superSubTematiks
				subSubTematiks = append(subSubTematiks, subSubTematikResp)
			}

			subTematikResp.SubSubTematiks = subSubTematiks
			subTematiks = append(subTematiks, subTematikResp)
		}

		tematikResp.SubTematiks = subTematiks
		tematiks = append(tematiks, tematikResp)
	}

	return pohonkinerja.PohonKinerjaAdminResponse{
		Tahun:   tahun,
		Tematik: tematiks,
	}, nil
}

// Helper functions untuk konversi
func convertToTematikResponse(pokin domain.PohonKinerja) pohonkinerja.TematikResponse {
	return pohonkinerja.TematikResponse{
		Id:         pokin.Id,
		Parent:     nil,
		Tema:       pokin.NamaPohon,
		Keterangan: pokin.Keterangan,
		Indikators: convertToSimpleIndikators(pokin.Indikator),
	}
}

func convertToSubTematikResponse(pokin domain.PohonKinerja) pohonkinerja.SubtematikResponse {
	return pohonkinerja.SubtematikResponse{
		Id:         pokin.Id,
		Parent:     pokin.Parent,
		Tema:       pokin.NamaPohon,
		Keterangan: pokin.Keterangan,
		Indikators: convertToSimpleIndikators(pokin.Indikator),
	}
}

func convertToSubSubTematikResponse(pokin domain.PohonKinerja) pohonkinerja.SubSubTematikResponse {
	return pohonkinerja.SubSubTematikResponse{
		Id:         pokin.Id,
		Parent:     pokin.Parent,
		Tema:       pokin.NamaPohon,
		Keterangan: pokin.Keterangan,
		Indikators: convertToSimpleIndikators(pokin.Indikator),
	}
}

func convertToSuperSubTematikResponse(pokin domain.PohonKinerja) pohonkinerja.SuperSubTematikResponse {
	return pohonkinerja.SuperSubTematikResponse{
		Id:         pokin.Id,
		Parent:     pokin.Parent,
		Tema:       pokin.NamaPohon,
		Keterangan: pokin.Keterangan,
		Indikators: convertToSimpleIndikators(pokin.Indikator),
	}
}

func convertToStrategicResponse(pokin domain.PohonKinerja) pohonkinerja.StrategicResponse {
	return pohonkinerja.StrategicResponse{
		Id:              pokin.Id,
		Parent:          pokin.Parent,
		Strategi:        pokin.NamaPohon,
		Keterangan:      pokin.Keterangan,
		KodeOpd:         pokin.KodeOpd,
		PerangkatDaerah: "", // Sesuaikan dengan data yang tersedia
		Indikators:      convertToSimpleIndikators(pokin.Indikator),
	}
}

func convertToTacticalResponse(pokin domain.PohonKinerja) pohonkinerja.TacticalResponse {
	return pohonkinerja.TacticalResponse{
		Id:              pokin.Id,
		Parent:          pokin.Parent,
		Strategi:        pokin.NamaPohon,
		Keterangan:      &pokin.Keterangan,
		KodeOpd:         pokin.KodeOpd,
		PerangkatDaerah: "", // Sesuaikan dengan data yang tersedia
		Indikators:      convertToSimpleIndikators(pokin.Indikator),
	}
}

func convertToOperationalResponse(pokin domain.PohonKinerja) pohonkinerja.OperationalResponse {
	return pohonkinerja.OperationalResponse{
		Id:              pokin.Id,
		Parent:          pokin.Parent,
		Strategi:        pokin.NamaPohon,
		Keterangan:      &pokin.Keterangan,
		KodeOpd:         pokin.KodeOpd,
		PerangkatDaerah: "", // Sesuaikan dengan data yang tersedia
		Indikators:      convertToSimpleIndikators(pokin.Indikator),
	}
}

func convertToSimpleIndikators(indikators []domain.Indikator) []pohonkinerja.IndikatorSimpleResponse {
	var result []pohonkinerja.IndikatorSimpleResponse
	for _, ind := range indikators {
		if len(ind.Target) > 0 {
			result = append(result, pohonkinerja.IndikatorSimpleResponse{
				Indikator: ind.Indikator,
				Target:    fmt.Sprintf("%d", ind.Target[0].Target),
				Satuan:    ind.Target[0].Satuan,
			})
		}
	}
	return result
}
