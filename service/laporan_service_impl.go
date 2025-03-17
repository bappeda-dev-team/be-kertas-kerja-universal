package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/laporan"
	"ekak_kabupaten_madiun/repository"
)

type LaporanServiceImpl struct {
	laporanRepository      repository.LaporanRepository
	opdRepository          repository.OpdRepository
	pohonKinerjaRepository repository.PohonKinerjaRepository
	DB                     *sql.DB
}

func NewLaporanServiceImpl(laporanRepository repository.LaporanRepository, opdRepository repository.OpdRepository, pohonKinerjaRepository repository.PohonKinerjaRepository, DB *sql.DB) *LaporanServiceImpl {
	return &LaporanServiceImpl{
		laporanRepository:      laporanRepository,
		opdRepository:          opdRepository,
		pohonKinerjaRepository: pohonKinerjaRepository,
		DB:                     DB,
	}
}

func (service *LaporanServiceImpl) OpdSupportingPokin(ctx context.Context, kodeOpd string, tahun string) (laporan.OpdSupportingPokinResponseData, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return laporan.OpdSupportingPokinResponseData{}, err
	}
	defer helper.CommitOrRollback(tx)

	pokins, err := service.laporanRepository.OpdSupportingPokin(ctx, tx, kodeOpd, tahun)
	if err != nil {
		return laporan.OpdSupportingPokinResponseData{}, err
	}

	pohonMap := make(map[int]map[int][]domain.OpdSupportingPokin)

	for i := range pokins {
		level := pokins[i].LevelPohon

		if pohonMap[level] == nil {
			pohonMap[level] = make(map[int][]domain.OpdSupportingPokin)
		}

		if pokins[i].KodeOpd != "" {
			opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, pokins[i].KodeOpd)
			if err != nil {
				return laporan.OpdSupportingPokinResponseData{}, err
			}
			pokins[i].NamaOpd = opd.NamaOpd
		}

		pohonMap[level][pokins[i].Parent] = append(
			pohonMap[level][pokins[i].Parent],
			pokins[i],
		)
	}

	var pohonKinerjas []laporan.PokinSupporting
	for _, tematik := range pohonMap[0][0] {
		tematikResp := laporan.PokinSupporting{
			Id:         tematik.Id,
			Parent:     tematik.Parent,
			Tema:       tematik.NamaPohon,
			JenisPohon: tematik.JenisPohon,
			Keterangan: tematik.Keterangan,
			LevelPohon: tematik.LevelPohon,
			Indikators: []laporan.IndikatorResponse{},
		}

		for _, indikator := range tematik.Indikators {
			indikatorResp := laporan.IndikatorResponse{
				Id:            indikator.Id,
				IdPokin:       indikator.PokinId,
				NamaIndikator: indikator.Indikator,
				Targets:       []laporan.TargetResponse{},
			}

			for _, target := range indikator.Target {
				targetResp := laporan.TargetResponse{
					Id:              target.Id,
					IndikatorId:     target.IndikatorId,
					TargetIndikator: target.Target,
					SatuanIndikator: target.Satuan,
					TahunSasaran:    target.Tahun,
				}
				indikatorResp.Targets = append(indikatorResp.Targets, targetResp)
			}

			tematikResp.Indikators = append(tematikResp.Indikators, indikatorResp)
		}

		pohonKinerjas = append(pohonKinerjas, tematikResp)
	}

	return laporan.OpdSupportingPokinResponseData{
		Tahun:         tahun,
		PohonKinerjas: pohonKinerjas,
	}, nil
}
