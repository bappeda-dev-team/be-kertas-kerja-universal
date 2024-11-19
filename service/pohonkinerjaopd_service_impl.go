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
	for i := 4; i <= 6; i++ {
		pohonMap[i] = make(map[int][]domain.PohonKinerja)
	}

	// Kelompokkan data berdasarkan level dan parent
	for _, p := range pokins {
		if p.LevelPohon >= 4 {
			p.NamaOpd = opd.NamaOpd
			pohonMap[p.LevelPohon][p.Parent] = append(
				pohonMap[p.LevelPohon][p.Parent],
				p,
			)
		}
	}

	// Build response untuk strategic (level 4)
	var strategics []pohonkinerja.StrategicOpdResponse
	for _, strategicList := range pohonMap[4] {
		sort.Slice(strategicList, func(i, j int) bool {
			return strategicList[i].Id < strategicList[j].Id
		})

		for _, strategic := range strategicList {
			var tacticals []pohonkinerja.TacticalOpdResponse
			var strategicPelaksana []pohonkinerja.PelaksanaOpdResponse

			// Get pelaksana untuk strategic
			for _, pelaksana := range strategic.Pelaksana {
				pegawaiPelaksana, err := service.pegawaiRepository.FindById(ctx, tx, pelaksana.PegawaiId)
				if err != nil {
					continue
				}
				strategicPelaksana = append(strategicPelaksana, pohonkinerja.PelaksanaOpdResponse{
					Id:          pelaksana.Id,
					PegawaiId:   pegawaiPelaksana.Id,
					NamaPegawai: pegawaiPelaksana.NamaPegawai,
				})
			}

			// Build tactical (level 5)
			if tacticalList := pohonMap[5][strategic.Id]; len(tacticalList) > 0 {
				sort.Slice(tacticalList, func(i, j int) bool {
					return tacticalList[i].Id < tacticalList[j].Id
				})

				for _, tactical := range tacticalList {
					var operationals []pohonkinerja.OperationalOpdResponse
					var tacticalPelaksana []pohonkinerja.PelaksanaOpdResponse

					// Get pelaksana untuk tactical
					for _, pelaksana := range tactical.Pelaksana {
						pegawaiPelaksana, err := service.pegawaiRepository.FindById(ctx, tx, pelaksana.PegawaiId)
						if err != nil {
							continue
						}
						tacticalPelaksana = append(tacticalPelaksana, pohonkinerja.PelaksanaOpdResponse{
							Id:          pelaksana.Id,
							PegawaiId:   pegawaiPelaksana.Id,
							NamaPegawai: pegawaiPelaksana.NamaPegawai,
						})
					}

					// Build operational (level 6)
					if operationalList := pohonMap[6][tactical.Id]; len(operationalList) > 0 {
						sort.Slice(operationalList, func(i, j int) bool {
							return operationalList[i].Id < operationalList[j].Id
						})

						for _, operational := range operationalList {
							var operationalPelaksana []pohonkinerja.PelaksanaOpdResponse

							// Get pelaksana untuk operational
							for _, pelaksana := range operational.Pelaksana {
								pegawaiPelaksana, err := service.pegawaiRepository.FindById(ctx, tx, pelaksana.PegawaiId)
								if err != nil {
									continue
								}
								operationalPelaksana = append(operationalPelaksana, pohonkinerja.PelaksanaOpdResponse{
									Id:          pelaksana.Id,
									PegawaiId:   pegawaiPelaksana.Id,
									NamaPegawai: pegawaiPelaksana.NamaPegawai,
								})
							}

							operationals = append(operationals, pohonkinerja.OperationalOpdResponse{
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
								Pelaksana: operationalPelaksana,
							})
						}
					}

					tacticals = append(tacticals, pohonkinerja.TacticalOpdResponse{
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
						Pelaksana:    tacticalPelaksana,
						Operationals: operationals,
					})
				}
			}

			strategics = append(strategics, pohonkinerja.StrategicOpdResponse{
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
				Pelaksana: strategicPelaksana,
				Tacticals: tacticals,
			})
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
