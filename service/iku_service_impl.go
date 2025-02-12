package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web/iku"
	"ekak_kabupaten_madiun/repository"
)

type IkuServiceImpl struct {
	IkuRepository repository.IkuRepository
	DB            *sql.DB
}

func NewIkuServiceImpl(ikuRepository repository.IkuRepository, db *sql.DB) *IkuServiceImpl {
	return &IkuServiceImpl{
		IkuRepository: ikuRepository,
		DB:            db,
	}
}

func (service *IkuServiceImpl) FindAll(ctx context.Context, tahun string) ([]iku.IkuResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data dari repository
	indikatorTargets, err := service.IkuRepository.FindAll(ctx, tx, tahun)
	if err != nil {
		return nil, err
	}

	// Transform ke response
	var responses []iku.IkuResponse
	for _, item := range indikatorTargets {
		var targetResponses []iku.TargetResponse
		for _, target := range item.Target {
			targetResponses = append(targetResponses, iku.TargetResponse{
				Target: target.Target,
				Satuan: target.Satuan,
			})
		}

		responses = append(responses, iku.IkuResponse{
			IndikatorId: item.Id,
			Sumber:      item.Sumber,
			Indikator:   item.Indikator,
			CreatedAt:   item.CreatedAt,
			Target:      targetResponses,
		})
	}

	return responses, nil
}
