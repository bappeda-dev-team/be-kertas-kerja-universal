package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web/sasaranopd"
	"ekak_kabupaten_madiun/repository"
)

type SasaranOpdServiceImpl struct {
	sasaranOpdRepository      repository.SasaranOpdRepository
	opdRepository             repository.OpdRepository
	rencanaKinerjaRepository  repository.RencanaKinerjaRepository
	manualIndikatorRepository repository.ManualIKRepository

	DB *sql.DB
}

func NewSasaranOpdServiceImpl(
	sasaranOpdRepository repository.SasaranOpdRepository,
	opdRepository repository.OpdRepository,
	rencanaKinerjaRepository repository.RencanaKinerjaRepository,
	manualIndikatorRepository repository.ManualIKRepository,
	db *sql.DB) *SasaranOpdServiceImpl {
	return &SasaranOpdServiceImpl{
		sasaranOpdRepository:      sasaranOpdRepository,
		opdRepository:             opdRepository,
		rencanaKinerjaRepository:  rencanaKinerjaRepository,
		manualIndikatorRepository: manualIndikatorRepository,
		DB:                        db,
	}
}

func (service *SasaranOpdServiceImpl) FindAll(ctx context.Context, KodeOpd string, tahunAwal string, tahunAkhir string) ([]sasaranopd.SasaranOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	sasaranOpds, err := service.sasaranOpdRepository.FindAll(ctx, tx, KodeOpd, tahunAwal, tahunAkhir)
	if err != nil {
		return nil, err
	}

	var responses []sasaranopd.SasaranOpdResponse
	for _, sasaranOpd := range sasaranOpds {
		response := sasaranopd.SasaranOpdResponse{
			Id:                sasaranOpd.Id,
			IdPohon:           sasaranOpd.IdPohon,
			NamaPohon:         sasaranOpd.NamaPohon,
			JenisPohon:        sasaranOpd.JenisPohon,
			LevelPohon:        sasaranOpd.LevelPohon,
			TahunPohon:        sasaranOpd.TahunPohon,
			TahunAwalPeriode:  sasaranOpd.TahunAwalPeriode,
			TahunAkhirPeriode: sasaranOpd.TahunAkhirPeriode,
			Pelaksana:         []sasaranopd.PelaksanaOpdResponse{},
		}

		// Convert Pelaksana
		for _, pelaksana := range sasaranOpd.Pelaksana {
			response.Pelaksana = append(response.Pelaksana, sasaranopd.PelaksanaOpdResponse{
				Id:  pelaksana.Id,
				Nip: pelaksana.Nip,
			})
		}

		// Convert RencanaKinerja if exists
		if sasaranOpd.IdRencanaKinerja != "" {
			response.RencanaKinerja = &sasaranopd.RencanaKinerjaOpd{
				Id:                 sasaranOpd.IdRencanaKinerja,
				NamaRencanaKinerja: sasaranOpd.NamaRencanaKinerja,
				Nip:                sasaranOpd.PegawaiId,
				TahunAwal:          sasaranOpd.TahunAwalRencana,
				TahunAkhir:         sasaranOpd.TahunAkhirRencana,
				Indikator:          []sasaranopd.IndikatorResponse{},
			}

			// Convert Indikator
			for _, indikator := range sasaranOpd.IndikatorSasaranOpd {
				indResponse := sasaranopd.IndikatorResponse{
					Id:        indikator.Id,
					Indikator: indikator.Indikator,
					Target:    []sasaranopd.TargetResponse{},
				}

				if indikator.ManualIK != nil {
					indResponse.ManualIK = &sasaranopd.ManualIKResponse{
						Formula:    indikator.ManualIK.Formula,
						SumberData: indikator.ManualIK.SumberData,
					}
				}

				for _, target := range indikator.Target {
					indResponse.Target = append(indResponse.Target, sasaranopd.TargetResponse{
						Tahun:  target.Tahun,
						Target: target.Target,
						Satuan: target.Satuan,
					})
				}

				response.RencanaKinerja.Indikator = append(response.RencanaKinerja.Indikator, indResponse)
			}
		}

		responses = append(responses, response)
	}

	return responses, nil
}

// func (service *SasaranOpdServiceImpl) FindAll(ctx context.Context, KodeOpd string, tahun string) ([]sasaranopd.SasaranOpdResponse, error) {
// 	tx, err := service.DB.Begin()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer helper.CommitOrRollback(tx)

// 	sasaranOpds, err := service.sasaranOpdRepository.FindAll(ctx, tx, KodeOpd, tahun)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var responses []sasaranopd.SasaranOpdResponse
// 	for _, sasaranOpd := range sasaranOpds {
// 		response := sasaranopd.SasaranOpdResponse{
// 			Id:         sasaranOpd.Id,
// 			IdPohon:    sasaranOpd.IdPohon,
// 			NamaPohon:  sasaranOpd.NamaPohon,
// 			JenisPohon: sasaranOpd.JenisPohon,
// 			LevelPohon: sasaranOpd.LevelPohon,
// 			// TahunPohon: sasaranOpd.TahunPohon,
// 			Pelaksana: []sasaranopd.PelaksanaOpdResponse{},
// 		}

// 		// Convert Pelaksana
// 		for _, pelaksana := range sasaranOpd.Pelaksana {
// 			response.Pelaksana = append(response.Pelaksana, sasaranopd.PelaksanaOpdResponse{
// 				Id:  pelaksana.Id,
// 				Nip: pelaksana.PegawaiId,
// 			})
// 		}

// 		// Convert RencanaKinerja if exists
// 		if sasaranOpd.IdRencanaKinerja != "" {
// 			response.RencanaKinerja = &sasaranopd.RencanaKinerjaOpd{
// 				Id:                 sasaranOpd.IdRencanaKinerja,
// 				NamaRencanaKinerja: sasaranOpd.NamaRencanaKinerja,
// 				Nip:                sasaranOpd.PegawaiId,
// 				// TahunAwal:          sasaranOpd.TahunAwalRencana,
// 				// TahunAkhir:         sasaranOpd.TahunAkhirRencana,
// 				Indikator: []sasaranopd.IndikatorResponse{},
// 			}

// 			// Convert Indikator
// 			for _, indikator := range sasaranOpd.IndikatorSasaranOpd {
// 				indResponse := sasaranopd.IndikatorResponse{
// 					Id:        indikator.Id,
// 					Indikator: indikator.Indikator,
// 					Target:    []sasaranopd.TargetResponse{},
// 				}

// 				if indikator.ManualIK != nil {
// 					indResponse.ManualIK = &sasaranopd.ManualIKResponse{
// 						Formula:    indikator.ManualIK.Formula,
// 						SumberData: indikator.ManualIK.SumberData,
// 					}
// 				}

// 				// Convert Target
// 				for _, target := range indikator.Target {
// 					indResponse.Target = append(indResponse.Target, sasaranopd.TargetResponse{
// 						Tahun:  target.Tahun,
// 						Target: target.Target,
// 						Satuan: target.Satuan,
// 					})
// 				}

// 				response.RencanaKinerja.Indikator = append(response.RencanaKinerja.Indikator, indResponse)
// 			}
// 		}

// 		responses = append(responses, response)
// 	}

// 	return responses, nil
// }
