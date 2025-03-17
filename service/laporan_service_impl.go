package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/laporan"
	"ekak_kabupaten_madiun/repository"
	"errors"
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

// OPD Supporting pokin
// find strategic opd dibawah sub tema apa
// telusuri hingga ke atas (tematik)
func (service *LaporanServiceImpl) OpdSupportingPokin(ctx context.Context, kodeOpd string, tahun string) (laporan.OpdSupportingPokinResponseData, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return laporan.OpdSupportingPokinResponseData{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return laporan.OpdSupportingPokinResponseData{}, errors.New("kode opd tidak ditemukan")
	}

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

		pohonMap[level][pokins[i].Parent] = append(
			pohonMap[level][pokins[i].Parent],
			pokins[i],
		)
	}

	var pohonKinerjas []laporan.PokinSupporting
	var buildTree func(parentId int, level int) []laporan.PokinSupporting

	buildTree = func(parentId int, level int) []laporan.PokinSupporting {
		var childs []laporan.PokinSupporting

		// Check if the current level exists in pohonMap and has children
		if nodes, exists := pohonMap[level][parentId]; exists {
			for _, node := range nodes {
				nodeResp := laporan.PokinSupporting{
					Id:         node.Id,
					Parent:     node.Parent,
					Tema:       node.NamaPohon,
					JenisPohon: node.JenisPohon,
					Keterangan: node.Keterangan,
					LevelPohon: node.LevelPohon,
					Tahun: node.Tahun,
					Indikators: []laporan.IndikatorResponse{},
					Childs:     buildTree(node.Id, level+1), // Recursively fetch children
				}

				// Add Indikators
				for _, indikator := range node.Indikators {
					indikatorResp := laporan.IndikatorResponse{
						Id:            indikator.Id,
						IdPokin:       indikator.PokinId,
						NamaIndikator: indikator.Indikator,
						Targets:       []laporan.TargetResponse{},
					}

					// Add Targets
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

					nodeResp.Indikators = append(nodeResp.Indikators, indikatorResp)
				}

				childs = append(childs, nodeResp)
			}
		}

		return childs
	}

	// Build Parent-Child
	// Build the tree starting from level 0 (root nodes)
	pohonKinerjas = buildTree(0, 0)

	return laporan.OpdSupportingPokinResponseData{
		Tahun:         tahun,
		KodeOpd:       kodeOpd,
		NamaOpd:       opd.NamaOpd,
		PohonKinerjas: pohonKinerjas,
	}, nil
}
