package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"ekak_kabupaten_madiun/repository"

	"log"

	"fmt"

	"github.com/google/uuid"
)

type PohonKinerjaAdminServiceImpl struct {
	pohonKinerjaRepository repository.PohonKinerjaRepository
	opdRepository          repository.OpdRepository
	DB                     *sql.DB
}

func NewPohonKinerjaAdminServiceImpl(pohonKinerjaRepository repository.PohonKinerjaRepository, opdRepository repository.OpdRepository, DB *sql.DB) *PohonKinerjaAdminServiceImpl {
	return &PohonKinerjaAdminServiceImpl{
		pohonKinerjaRepository: pohonKinerjaRepository,
		opdRepository:          opdRepository,
		DB:                     DB,
	}
}

func (service *PohonKinerjaAdminServiceImpl) Create(ctx context.Context, request pohonkinerja.PohonKinerjaAdminCreateRequest) (pohonkinerja.PohonKinerjaAdminResponseData, error) {
	log.Printf("Memulai proses pembuatan PohonKinerja untuk tahun: %s", request.Tahun)

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Error memulai transaksi: %v", err)
		return pohonkinerja.PohonKinerjaAdminResponseData{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Logging persiapan indikator
	log.Printf("Mempersiapkan %d indikator", len(request.Indikator))

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
		KodeOpd:    helper.EmptyStringIfNull(request.KodeOpd),
		Keterangan: request.Keterangan,
		Tahun:      request.Tahun,
		Indikator:  indikators,
	}

	log.Printf("Menyimpan PohonKinerja dengan NamaPohon: %s, LevelPohon: %d", request.NamaPohon, request.LevelPohon)
	result, err := service.pohonKinerjaRepository.CreatePokinAdmin(ctx, tx, pohonKinerja)
	if err != nil {
		log.Printf("Error saat menyimpan PohonKinerja: %v", err)
		return pohonkinerja.PohonKinerjaAdminResponseData{}, err
	}

	log.Printf("Berhasil membuat PohonKinerja dengan ID: %d", result.Id)

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

	log.Printf("Proses pembuatan PohonKinerja selesai")
	return response, nil
}

func (service *PohonKinerjaAdminServiceImpl) Update(ctx context.Context, request pohonkinerja.PohonKinerjaAdminUpdateRequest) (pohonkinerja.PohonKinerjaAdminResponseData, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponseData{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Cek apakah data exists
	_, err = service.pohonKinerjaRepository.FindPokinAdminById(ctx, tx, request.Id)
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponseData{}, err
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
		KodeOpd:    helper.EmptyStringIfNull(request.KodeOpd),
		Keterangan: request.Keterangan,
		Tahun:      request.Tahun,
		Indikator:  indikators,
	}

	result, err := service.pohonKinerjaRepository.UpdatePokinAdmin(ctx, tx, pohonKinerja)
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

func (service *PohonKinerjaAdminServiceImpl) Delete(ctx context.Context, id int) error {
	// Mulai transaksi
	tx, err := service.DB.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer helper.CommitOrRollback(tx)

	// Cek apakah data exists sebelum dihapus
	pokin, err := service.pohonKinerjaRepository.FindPokinAdminById(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("data tidak ditemukan: %v", err)
	}

	// Validasi tambahan: pastikan data yang akan dihapus memiliki level yang sesuai
	// Ini opsional, tergantung kebutuhan bisnis
	if pokin.LevelPohon < 0 || pokin.LevelPohon > 6 {
		return fmt.Errorf("level pohon kinerja tidak valid")
	}

	// Lakukan penghapusan secara hierarki
	err = service.pohonKinerjaRepository.DeletePokinAdmin(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus data: %v", err)
	}

	return nil
}

func (service *PohonKinerjaAdminServiceImpl) FindById(ctx context.Context, id int) (pohonkinerja.PohonKinerjaAdminResponseData, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponseData{}, err
	}
	defer helper.CommitOrRollback(tx)

	pokin, err := service.pohonKinerjaRepository.FindPokinAdminById(ctx, tx, id)
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponseData{}, err
	}

	log.Printf("Pohon Kinerja ditemukan: %+v", pokin)

	// Ambil data OPD jika kode OPD ada
	if pokin.KodeOpd != "" {
		opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, pokin.KodeOpd)
		if err == nil {
			pokin.NamaOpd = opd.NamaOpd
		}
	}

	// Konversi pokin.Id dari int ke string
	pokinIdStr := fmt.Sprint(pokin.Id)

	// Ambil indikator berdasarkan pokin ID
	indikators, err := service.pohonKinerjaRepository.FindIndikatorByPokinId(ctx, tx, pokinIdStr)
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponseData{}, err
	}

	log.Printf("Indikator ditemukan: %+v", indikators)

	// Konversi indikator domain ke IndikatorResponse
	var indikatorResponses []pohonkinerja.IndikatorResponse
	for _, ind := range indikators {
		// Ambil target berdasarkan indikator ID
		targets, err := service.pohonKinerjaRepository.FindTargetByIndikatorId(ctx, tx, ind.Id)
		if err != nil {
			return pohonkinerja.PohonKinerjaAdminResponseData{}, err
		}

		var targetResponses []pohonkinerja.TargetResponse
		for _, t := range targets {
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
			IdPokin:       ind.PokinId,
			NamaIndikator: ind.Indikator,
			Target:        targetResponses,
		}
		indikatorResponses = append(indikatorResponses, indikatorResponse)
	}

	response := pohonkinerja.PohonKinerjaAdminResponseData{
		Id:         pokin.Id,
		Parent:     pokin.Parent,
		NamaPohon:  pokin.NamaPohon,
		NamaOpd:    pokin.NamaOpd,
		JenisPohon: pokin.JenisPohon,
		LevelPohon: pokin.LevelPohon,
		KodeOpd:    pokin.KodeOpd,
		Keterangan: pokin.Keterangan,
		Tahun:      pokin.Tahun,
		Indikators: indikatorResponses,
	}

	return response, nil
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

	// Kelompokkan data dan ambil data OPD untuk setiap pohon kinerja
	for i := range pokins {
		if pokins[i].KodeOpd != "" {
			opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, pokins[i].KodeOpd)
			if err == nil {
				// Update data pohon kinerja dengan nama OPD
				pokins[i].NamaOpd = opd.NamaOpd
			}
		}
		pohonMap[pokins[i].LevelPohon][pokins[i].Parent] = append(
			pohonMap[pokins[i].LevelPohon][pokins[i].Parent],
			pokins[i],
		)
	}

	// Bangun response dimulai dari Tematik (level 0)
	var tematiks []pohonkinerja.TematikResponse
	for _, tematik := range pohonMap[0][0] {
		tematikResp := helper.BuildTematikResponse(pohonMap, tematik)
		tematiks = append(tematiks, tematikResp)
	}

	return pohonkinerja.PohonKinerjaAdminResponse{
		Tahun:   tahun,
		Tematik: tematiks,
	}, nil
}

func (service *PohonKinerjaAdminServiceImpl) FindSubTematik(ctx context.Context, tahun string) (pohonkinerja.PohonKinerjaAdminResponse, error) {
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

	// Inisialisasi map dengan benar
	pohonMap := make(map[int]map[int][]domain.PohonKinerja)
	pohonMap[0] = make(map[int][]domain.PohonKinerja) // Inisialisasi untuk level 0
	pohonMap[1] = make(map[int][]domain.PohonKinerja) // Inisialisasi untuk level 1

	// Filter dan kelompokkan data hanya untuk level 0 dan 1
	for _, pokin := range pokins {
		if pokin.LevelPohon <= 1 { // Hanya proses level 0 dan 1
			pohonMap[pokin.LevelPohon][pokin.Parent] = append(pohonMap[pokin.LevelPohon][pokin.Parent], pokin)
		}
	}

	// Bangun response dimulai dari Tematik (level 0)
	var tematiks []pohonkinerja.TematikResponse
	for _, tematik := range pohonMap[0][0] {
		tematikResp := helper.BuildTematikResponseLimited(pohonMap, tematik) // Gunakan BuildTematikResponseLimited
		tematiks = append(tematiks, tematikResp)
	}

	return pohonkinerja.PohonKinerjaAdminResponse{
		Tahun:   tahun,
		Tematik: tematiks,
	}, nil
}

func (service *PohonKinerjaAdminServiceImpl) FindPokinAdminByIdHierarki(ctx context.Context, idPokin int) (pohonkinerja.TematikResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.TematikResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Validasi level_pohon
	pokin, err := service.pohonKinerjaRepository.FindPokinAdminById(ctx, tx, idPokin)
	if err != nil {
		return pohonkinerja.TematikResponse{}, err
	}

	if pokin.LevelPohon != 0 {
		return pohonkinerja.TematikResponse{}, fmt.Errorf("ID yang diberikan bukan merupakan Tematik (level 0)")
	}

	// Ambil semua data pohon kinerja
	pokins, err := service.pohonKinerjaRepository.FindPokinAdminByIdHierarki(ctx, tx, idPokin)
	if err != nil {
		return pohonkinerja.TematikResponse{}, err
	}

	// Buat map untuk menyimpan data berdasarkan level dan parent
	pohonMap := make(map[int]map[int][]domain.PohonKinerja)
	for i := 0; i <= 6; i++ {
		pohonMap[i] = make(map[int][]domain.PohonKinerja)
	}

	// Kelompokkan data dan ambil data OPD untuk setiap pohon kinerja
	for i := range pokins {
		if pokins[i].KodeOpd != "" {
			opd, err := service.opdRepository.FindByKodeOpd(ctx, tx, pokins[i].KodeOpd)
			if err == nil {
				// Update data pohon kinerja dengan nama OPD
				pokins[i].NamaOpd = opd.NamaOpd
			}
		}
		pohonMap[pokins[i].LevelPohon][pokins[i].Parent] = append(
			pohonMap[pokins[i].LevelPohon][pokins[i].Parent],
			pokins[i],
		)
	}

	// Bangun response dimulai dari Tematik (level 0)
	var tematiks []pohonkinerja.TematikResponse
	for _, tematik := range pohonMap[0][0] {
		tematikResp := helper.BuildTematikResponse(pohonMap, tematik)
		tematiks = append(tematiks, tematikResp)
	}

	return tematiks[0], nil
}

func (service *PohonKinerjaAdminServiceImpl) CreateStrategicAdmin(ctx context.Context, request pohonkinerja.PohonKinerjaAdminStrategicCreateRequest) (pohonkinerja.PohonKinerjaAdminResponseData, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponseData{}, err
	}
	defer helper.CommitOrRollback(tx)

	existingPokin, err := service.pohonKinerjaRepository.FindPokinToClone(ctx, tx, request.IdToClone)
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponseData{}, err
	}

	err = service.pohonKinerjaRepository.ValidateParentLevel(ctx, tx, request.Parent, existingPokin.LevelPohon)
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponseData{}, err
	}

	// 3. Siapkan data baru
	newPokin := domain.PohonKinerja{
		Parent:     request.Parent,
		NamaPohon:  existingPokin.NamaPohon,
		JenisPohon: existingPokin.JenisPohon,
		LevelPohon: existingPokin.LevelPohon,
		KodeOpd:    existingPokin.KodeOpd,
		Keterangan: existingPokin.Keterangan,
		Tahun:      existingPokin.Tahun,
	}

	newPokinId, err := service.pohonKinerjaRepository.InsertClonedPokin(ctx, tx, newPokin)
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponseData{}, err
	}

	indikators, err := service.pohonKinerjaRepository.FindIndikatorToClone(ctx, tx, request.IdToClone)
	if err != nil {
		return pohonkinerja.PohonKinerjaAdminResponseData{}, err
	}

	var indikatorResponses []pohonkinerja.IndikatorResponse

	for _, indikator := range indikators {
		newIndikatorId := "IND-POKIN-" + uuid.New().String()[:6]

		err = service.pohonKinerjaRepository.InsertClonedIndikator(ctx, tx, newIndikatorId, newPokinId, indikator)
		if err != nil {
			return pohonkinerja.PohonKinerjaAdminResponseData{}, err
		}

		targets, err := service.pohonKinerjaRepository.FindTargetToClone(ctx, tx, indikator.Id)
		if err != nil {
			return pohonkinerja.PohonKinerjaAdminResponseData{}, err
		}

		var targetResponses []pohonkinerja.TargetResponse

		for _, target := range targets {
			newTargetId := "TRGT-IND-POKIN-" + uuid.New().String()[:5]
			err = service.pohonKinerjaRepository.InsertClonedTarget(ctx, tx, newTargetId, newIndikatorId, target)
			if err != nil {
				return pohonkinerja.PohonKinerjaAdminResponseData{}, err
			}

			targetResponses = append(targetResponses, pohonkinerja.TargetResponse{
				Id:              newTargetId,
				IndikatorId:     newIndikatorId,
				TargetIndikator: target.Target,
				SatuanIndikator: target.Satuan,
			})
		}

		indikatorResponses = append(indikatorResponses, pohonkinerja.IndikatorResponse{
			Id:            newIndikatorId,
			IdPokin:       fmt.Sprint(newPokinId),
			NamaIndikator: indikator.Indikator,
			Target:        targetResponses,
		})
	}

	response := pohonkinerja.PohonKinerjaAdminResponseData{
		Id:         int(newPokinId),
		Parent:     existingPokin.Parent,
		NamaPohon:  existingPokin.NamaPohon,
		JenisPohon: existingPokin.JenisPohon,
		LevelPohon: existingPokin.LevelPohon,
		KodeOpd:    existingPokin.KodeOpd,
		Keterangan: existingPokin.Keterangan,
		Tahun:      existingPokin.Tahun,
		Indikators: indikatorResponses,
	}

	return response, nil
}
