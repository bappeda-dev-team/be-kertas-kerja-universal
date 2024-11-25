package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/opdmaster"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"ekak_kabupaten_madiun/repository"
	"errors"
	"fmt"
	"sort"
	"strconv"

	"github.com/google/uuid"
)

type PohonKinerjaOpdServiceImpl struct {
	pohonKinerjaOpdRepository repository.PohonKinerjaRepository
	opdRepository             repository.OpdRepository
	pegawaiRepository         repository.PegawaiRepository
	DB                        *sql.DB
}

func NewPohonKinerjaOpdServiceImpl(pohonKinerjaOpdRepository repository.PohonKinerjaRepository, opdRepository repository.OpdRepository, pegawaiRepository repository.PegawaiRepository, DB *sql.DB) *PohonKinerjaOpdServiceImpl {
	return &PohonKinerjaOpdServiceImpl{
		pohonKinerjaOpdRepository: pohonKinerjaOpdRepository,
		opdRepository:             opdRepository,
		pegawaiRepository:         pegawaiRepository,
		DB:                        DB,
	}
}

func (service *PohonKinerjaOpdServiceImpl) Create(ctx context.Context, request pohonkinerja.PohonKinerjaCreateRequest) (pohonkinerja.PohonKinerjaOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi request
	if request.NamaPohon == "" {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("nama program tidak boleh kosong")
	}

	// Validasi kode OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, request.KodeOpd)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("kode opd tidak ditemukan")
	}
	if opd.KodeOpd == "" {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("kode opd tidak valid")
	}

	// Validasi dan persiapan data pelaksana
	var pelaksanaList []domain.PelaksanaPokin
	var pelaksanaResponses []pohonkinerja.PelaksanaOpdResponse

	for _, pelaksanaReq := range request.PelaksanaId {
		// Generate ID untuk pelaksana_pokin
		pelaksanaId := fmt.Sprintf("PLKS-%s", uuid.New().String()[:8])

		// Validasi setiap pelaksana
		pegawaiPelaksana, err := service.pegawaiRepository.FindById(ctx, tx, pelaksanaReq.PegawaiId)
		if err != nil {
			return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("pelaksana tidak ditemukan")
		}
		if pegawaiPelaksana.Id == "" {
			return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("pelaksana tidak ditemukan")
		}

		// Tambahkan ke list pelaksana
		pelaksanaList = append(pelaksanaList, domain.PelaksanaPokin{
			Id:        pelaksanaId,
			PegawaiId: pelaksanaReq.PegawaiId,
		})

		// Siapkan response pelaksana
		pelaksanaResponses = append(pelaksanaResponses, pohonkinerja.PelaksanaOpdResponse{
			Id:          pelaksanaId,
			PegawaiId:   pegawaiPelaksana.Id,
			NamaPegawai: pegawaiPelaksana.NamaPegawai,
		})
	}

	pohonKinerja := domain.PohonKinerja{
		NamaPohon:  request.NamaPohon,
		Parent:     request.Parent,
		JenisPohon: request.JenisPohon,
		LevelPohon: request.LevelPohon,
		KodeOpd:    request.KodeOpd,
		Keterangan: request.Keterangan,
		Tahun:      request.Tahun,
		Pelaksana:  pelaksanaList,
	}

	result, err := service.pohonKinerjaOpdRepository.Create(ctx, tx, pohonKinerja)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, err
	}

	response := pohonkinerja.PohonKinerjaOpdResponse{
		Id:         result.Id,
		Parent:     strconv.Itoa(result.Parent),
		NamaPohon:  result.NamaPohon,
		JenisPohon: result.JenisPohon,
		LevelPohon: result.LevelPohon,
		KodeOpd:    result.KodeOpd,
		NamaOpd:    opd.NamaOpd,
		Keterangan: result.Keterangan,
		Tahun:      result.Tahun,
		Pelaksana:  pelaksanaResponses,
	}

	return response, nil
}

func (service *PohonKinerjaOpdServiceImpl) Update(ctx context.Context, request pohonkinerja.PohonKinerjaUpdateRequest) (pohonkinerja.PohonKinerjaOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi request
	if request.NamaPohon == "" {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("nama program tidak boleh kosong")
	}

	// Validasi kode OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, request.KodeOpd)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("kode opd tidak ditemukan")
	}
	if opd.KodeOpd == "" {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("kode opd tidak valid")
	}

	// Validasi data yang akan diupdate
	_, err = service.pohonKinerjaOpdRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("data pohon kinerja tidak ditemukan")
	}

	// Validasi dan persiapan data pelaksana
	var pelaksanaList []domain.PelaksanaPokin
	var pelaksanaResponses []pohonkinerja.PelaksanaOpdResponse

	for _, pelaksanaReq := range request.PelaksanaId {
		// Generate ID untuk pelaksana_pokin
		pelaksanaId := fmt.Sprintf("PLKS-%s", uuid.New().String()[:8])

		// Validasi setiap pelaksana
		pegawaiPelaksana, err := service.pegawaiRepository.FindById(ctx, tx, pelaksanaReq.PegawaiId)
		if err != nil {
			return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("pelaksana tidak ditemukan")
		}
		if pegawaiPelaksana.Id == "" {
			return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("pelaksana tidak ditemukan")
		}

		// Tambahkan ke list pelaksana
		pelaksanaList = append(pelaksanaList, domain.PelaksanaPokin{
			Id:        pelaksanaId,
			PegawaiId: pelaksanaReq.PegawaiId,
		})

		// Siapkan response pelaksana
		pelaksanaResponses = append(pelaksanaResponses, pohonkinerja.PelaksanaOpdResponse{
			Id:          pelaksanaId,
			PegawaiId:   pegawaiPelaksana.Id,
			NamaPegawai: pegawaiPelaksana.NamaPegawai,
		})
	}

	pohonKinerja := domain.PohonKinerja{
		Id:         request.Id,
		NamaPohon:  request.NamaPohon,
		Parent:     request.Parent,
		JenisPohon: request.JenisPohon,
		LevelPohon: request.LevelPohon,
		KodeOpd:    request.KodeOpd,
		Keterangan: request.Keterangan,
		Tahun:      request.Tahun,
		Pelaksana:  pelaksanaList,
	}

	result, err := service.pohonKinerjaOpdRepository.Update(ctx, tx, pohonKinerja)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, err
	}

	return pohonkinerja.PohonKinerjaOpdResponse{
		Id:         result.Id,
		Parent:     strconv.Itoa(result.Parent),
		NamaPohon:  result.NamaPohon,
		JenisPohon: result.JenisPohon,
		LevelPohon: result.LevelPohon,
		KodeOpd:    result.KodeOpd,
		NamaOpd:    opd.NamaOpd,
		Keterangan: result.Keterangan,
		Tahun:      result.Tahun,
		Pelaksana:  pelaksanaResponses,
	}, nil
}

func (service *PohonKinerjaOpdServiceImpl) Delete(ctx context.Context, id string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	// Konversi id string ke int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("id tidak valid")
	}

	// Cek apakah data exists
	_, err = service.pohonKinerjaOpdRepository.FindById(ctx, tx, idInt)
	if err != nil {
		return errors.New("data pohon kinerja tidak ditemukan")
	}

	err = service.pohonKinerjaOpdRepository.Delete(ctx, tx, id)
	if err != nil {
		return err
	}

	return nil
}

func (service *PohonKinerjaOpdServiceImpl) FindById(ctx context.Context, id int) (pohonkinerja.PohonKinerjaOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// 1. Ambil data pohon kinerja
	pokin, err := service.pohonKinerjaOpdRepository.FindById(ctx, tx, id)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, err
	}

	// 2. Validasi data pohon kinerja
	if pokin.Id == 0 {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("data tidak ditemukan")
	}

	// 3. Ambil data OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, pokin.KodeOpd)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("data opd tidak ditemukan")
	}

	// 4. Ambil data pelaksana
	pelaksanaList, err := service.pohonKinerjaOpdRepository.FindPelaksanaPokin(ctx, tx, fmt.Sprint(pokin.Id))
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdResponse{}, errors.New("gagal mengambil data pelaksana")
	}

	// 5. Proses data pelaksana
	var pelaksanaResponses []pohonkinerja.PelaksanaOpdResponse
	for _, pelaksana := range pelaksanaList {
		pegawaiPelaksana, err := service.pegawaiRepository.FindById(ctx, tx, pelaksana.PegawaiId)
		if err != nil {
			continue // Skip jika pegawai tidak ditemukan
		}

		pelaksanaResponses = append(pelaksanaResponses, pohonkinerja.PelaksanaOpdResponse{
			Id:          pelaksana.Id,
			PegawaiId:   pegawaiPelaksana.Id,
			NamaPegawai: pegawaiPelaksana.NamaPegawai,
		})
	}

	// 6. Susun response
	response := pohonkinerja.PohonKinerjaOpdResponse{
		Id:         pokin.Id,
		Parent:     strconv.Itoa(pokin.Parent),
		NamaPohon:  pokin.NamaPohon,
		JenisPohon: pokin.JenisPohon,
		LevelPohon: pokin.LevelPohon,
		KodeOpd:    pokin.KodeOpd,
		NamaOpd:    opd.NamaOpd,
		Keterangan: pokin.Keterangan,
		Tahun:      pokin.Tahun,
		Pelaksana:  pelaksanaResponses,
	}

	return response, nil
}

func (service *PohonKinerjaOpdServiceImpl) FindAll(ctx context.Context, kodeOpd, tahun string) (pohonkinerja.PohonKinerjaOpdAllResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdAllResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdAllResponse{}, errors.New("kode opd tidak ditemukan")
	}

	// Ambil semua data pohon kinerja
	pokins, err := service.pohonKinerjaOpdRepository.FindAll(ctx, tx, kodeOpd, tahun)
	if err != nil {
		return pohonkinerja.PohonKinerjaOpdAllResponse{}, err
	}

	// Buat map untuk menyimpan data berdasarkan level dan parent
	pohonMap := make(map[int]map[int][]domain.PohonKinerja)
	pelaksanaMap := make(map[int][]pohonkinerja.PelaksanaOpdResponse)

	// Kelompokkan data dan ambil data pelaksana
	maxLevel := 0
	for _, p := range pokins {
		if p.LevelPohon >= 4 {
			// Update max level jika ditemukan level yang lebih tinggi
			if p.LevelPohon > maxLevel {
				maxLevel = p.LevelPohon
			}

			// Inisialisasi map untuk level jika belum ada
			if pohonMap[p.LevelPohon] == nil {
				pohonMap[p.LevelPohon] = make(map[int][]domain.PohonKinerja)
			}

			p.NamaOpd = opd.NamaOpd
			pohonMap[p.LevelPohon][p.Parent] = append(
				pohonMap[p.LevelPohon][p.Parent],
				p,
			)

			// Ambil data pelaksana
			pelaksanaList, err := service.pohonKinerjaOpdRepository.FindPelaksanaPokin(ctx, tx, fmt.Sprint(p.Id))
			if err == nil {
				var pelaksanaResponses []pohonkinerja.PelaksanaOpdResponse
				for _, pelaksana := range pelaksanaList {
					pegawaiPelaksana, err := service.pegawaiRepository.FindById(ctx, tx, pelaksana.PegawaiId)
					if err != nil {
						continue
					}
					pelaksanaResponses = append(pelaksanaResponses, pohonkinerja.PelaksanaOpdResponse{
						Id:          pelaksana.Id,
						PegawaiId:   pegawaiPelaksana.Id,
						NamaPegawai: pegawaiPelaksana.NamaPegawai,
					})
				}
				pelaksanaMap[p.Id] = pelaksanaResponses
			}
		}
	}

	// Build response untuk strategic (level 4)
	var strategics []pohonkinerja.StrategicOpdResponse
	if strategicList := pohonMap[4]; len(strategicList) > 0 {
		for _, strategicsByParent := range strategicList {
			sort.Slice(strategicsByParent, func(i, j int) bool {
				return strategicsByParent[i].Id < strategicsByParent[j].Id
			})

			for _, strategic := range strategicsByParent {
				strategicResp := service.buildStrategicResponse(ctx, tx, pohonMap, strategic, pelaksanaMap)
				strategics = append(strategics, strategicResp)
			}
		}
	}

	// Urutkan strategics berdasarkan Id
	sort.Slice(strategics, func(i, j int) bool {
		return strategics[i].Id < strategics[j].Id
	})

	return pohonkinerja.PohonKinerjaOpdAllResponse{
		KodeOpd:    kodeOpd,
		NamaOpd:    opd.NamaOpd,
		Tahun:      tahun,
		Strategics: strategics,
	}, nil
}

func (service *PohonKinerjaOpdServiceImpl) FindStrategicNoParent(ctx context.Context, kodeOpd, tahun string) ([]pohonkinerja.StrategicOpdResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi kode OPD
	opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, kodeOpd)
	if err != nil {
		return nil, errors.New("kode opd tidak ditemukan")
	}

	// Ambil data strategic dengan level pohon 4
	pokins, err := service.pohonKinerjaOpdRepository.FindStrategicNoParent(ctx, tx, 4, 0, kodeOpd, tahun)
	if err != nil {
		return nil, err
	}

	// Urutkan data berdasarkan ID
	sort.Slice(pokins, func(i, j int) bool {
		return pokins[i].Id < pokins[j].Id
	})

	// Konversi ke response format
	var strategics []pohonkinerja.StrategicOpdResponse
	for _, pokin := range pokins {
		strategic := pohonkinerja.StrategicOpdResponse{
			Id: pokin.Id,
			KodeOpd: opdmaster.OpdResponseForAll{
				KodeOpd: kodeOpd,
				NamaOpd: opd.NamaOpd,
			},
			Strategi:   pokin.NamaPohon,
			Keterangan: pokin.Keterangan,
		}
		strategics = append(strategics, strategic)
	}

	return strategics, nil
}

func (service *PohonKinerjaOpdServiceImpl) DeletePelaksana(ctx context.Context, pelaksanaId string) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)
	return service.pohonKinerjaOpdRepository.DeletePelaksanaPokin(ctx, tx, pelaksanaId)
}

// Tambahkan fungsi helper untuk membangun OperationalN response
func (service *PohonKinerjaOpdServiceImpl) buildOperationalNResponse(ctx context.Context, tx *sql.Tx, pohonMap map[int]map[int][]domain.PohonKinerja, operationalN domain.PohonKinerja, pelaksanaMap map[int][]pohonkinerja.PelaksanaOpdResponse) pohonkinerja.OperationalNOpdResponse {
	operationalNResp := pohonkinerja.OperationalNOpdResponse{
		Id:         operationalN.Id,
		Parent:     operationalN.Parent,
		Strategi:   operationalN.NamaPohon,
		JenisPohon: operationalN.JenisPohon,
		LevelPohon: operationalN.LevelPohon,
		Keterangan: operationalN.Keterangan,
		KodeOpd: opdmaster.OpdResponseForAll{
			KodeOpd: operationalN.KodeOpd,
			NamaOpd: operationalN.NamaOpd,
		},
		Pelaksana: pelaksanaMap[operationalN.Id],
	}

	// Build child nodes secara rekursif
	nextLevel := operationalN.LevelPohon + 1
	if nextOperationalNList := pohonMap[nextLevel][operationalN.Id]; len(nextOperationalNList) > 0 {
		var childs []pohonkinerja.OperationalNOpdResponse
		sort.Slice(nextOperationalNList, func(i, j int) bool {
			return nextOperationalNList[i].Id < nextOperationalNList[j].Id
		})

		for _, nextOpN := range nextOperationalNList {
			childResp := service.buildOperationalNResponse(ctx, tx, pohonMap, nextOpN, pelaksanaMap)
			childs = append(childs, childResp)
		}
		operationalNResp.Childs = childs
	}

	return operationalNResp
}

// Helper functions untuk membangun response
func (service *PohonKinerjaOpdServiceImpl) buildStrategicResponse(ctx context.Context, tx *sql.Tx, pohonMap map[int]map[int][]domain.PohonKinerja, strategic domain.PohonKinerja, pelaksanaMap map[int][]pohonkinerja.PelaksanaOpdResponse) pohonkinerja.StrategicOpdResponse {
	strategicResp := pohonkinerja.StrategicOpdResponse{
		Id:         strategic.Id,
		Parent:     nil,
		Strategi:   strategic.NamaPohon,
		JenisPohon: strategic.JenisPohon,
		LevelPohon: strategic.LevelPohon,
		Keterangan: strategic.Keterangan,
		KodeOpd: opdmaster.OpdResponseForAll{
			KodeOpd: strategic.KodeOpd,
			NamaOpd: strategic.NamaOpd,
		},
		Pelaksana: pelaksanaMap[strategic.Id],
	}

	// Build tactical (level 5)
	if tacticalList := pohonMap[5][strategic.Id]; len(tacticalList) > 0 {
		var tacticals []pohonkinerja.TacticalOpdResponse
		sort.Slice(tacticalList, func(i, j int) bool {
			return tacticalList[i].Id < tacticalList[j].Id
		})

		for _, tactical := range tacticalList {
			tacticalResp := service.buildTacticalResponse(ctx, tx, pohonMap, tactical, pelaksanaMap)
			tacticals = append(tacticals, tacticalResp)
		}
		strategicResp.Tacticals = tacticals
	}

	return strategicResp
}

func (service *PohonKinerjaOpdServiceImpl) buildTacticalResponse(ctx context.Context, tx *sql.Tx, pohonMap map[int]map[int][]domain.PohonKinerja, tactical domain.PohonKinerja, pelaksanaMap map[int][]pohonkinerja.PelaksanaOpdResponse) pohonkinerja.TacticalOpdResponse {
	tacticalResp := pohonkinerja.TacticalOpdResponse{
		Id:         tactical.Id,
		Parent:     tactical.Parent,
		Strategi:   tactical.NamaPohon,
		JenisPohon: tactical.JenisPohon,
		LevelPohon: tactical.LevelPohon,
		Keterangan: tactical.Keterangan,
		KodeOpd: opdmaster.OpdResponseForAll{
			KodeOpd: tactical.KodeOpd,
			NamaOpd: tactical.NamaOpd,
		},
		Pelaksana: pelaksanaMap[tactical.Id],
	}

	// Build operational (level 6)
	if operationalList := pohonMap[6][tactical.Id]; len(operationalList) > 0 {
		var operationals []pohonkinerja.OperationalOpdResponse
		sort.Slice(operationalList, func(i, j int) bool {
			return operationalList[i].Id < operationalList[j].Id
		})

		for _, operational := range operationalList {
			operationalResp := service.buildOperationalResponse(ctx, tx, pohonMap, operational, pelaksanaMap)
			operationals = append(operationals, operationalResp)
		}
		tacticalResp.Operationals = operationals
	}

	return tacticalResp
}

func (service *PohonKinerjaOpdServiceImpl) buildOperationalResponse(ctx context.Context, tx *sql.Tx, pohonMap map[int]map[int][]domain.PohonKinerja, operational domain.PohonKinerja, pelaksanaMap map[int][]pohonkinerja.PelaksanaOpdResponse) pohonkinerja.OperationalOpdResponse {
	operationalResp := pohonkinerja.OperationalOpdResponse{
		Id:         operational.Id,
		Parent:     operational.Parent,
		Strategi:   operational.NamaPohon,
		JenisPohon: operational.JenisPohon,
		LevelPohon: operational.LevelPohon,
		Keterangan: operational.Keterangan,
		KodeOpd: opdmaster.OpdResponseForAll{
			KodeOpd: operational.KodeOpd,
			NamaOpd: operational.NamaOpd,
		},
		Pelaksana: pelaksanaMap[operational.Id],
	}

	// Build operational-n untuk level > 6
	nextLevel := operational.LevelPohon + 1
	if operationalNList := pohonMap[nextLevel][operational.Id]; len(operationalNList) > 0 {
		var childs []pohonkinerja.OperationalNOpdResponse
		sort.Slice(operationalNList, func(i, j int) bool {
			return operationalNList[i].Id < operationalNList[j].Id
		})

		for _, opN := range operationalNList {
			childResp := service.buildOperationalNResponse(ctx, tx, pohonMap, opN, pelaksanaMap)
			childs = append(childs, childResp)
		}
		operationalResp.Childs = childs
	}

	return operationalResp
}
