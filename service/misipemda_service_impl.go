package service

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/domain"
	visimisipemda "ekak_kabupaten_madiun/model/web/visimisi"
	"ekak_kabupaten_madiun/repository"
	"fmt"
	"sort"

	"github.com/go-playground/validator/v10"
)

type MisiPemdaServiceImpl struct {
	MisiPemdaRepository repository.MisiPemdaRepository
	VisiPemdaRepository repository.VisiPemdaRepository
	Validate            *validator.Validate
	DB                  *sql.DB
}

func NewMisiPemdaServiceImpl(misiPemdaRepository repository.MisiPemdaRepository, visiPemdaRepository repository.VisiPemdaRepository, validate *validator.Validate, DB *sql.DB) *MisiPemdaServiceImpl {
	return &MisiPemdaServiceImpl{
		MisiPemdaRepository: misiPemdaRepository,
		VisiPemdaRepository: visiPemdaRepository,
		Validate:            validate,
		DB:                  DB,
	}
}

func (service *MisiPemdaServiceImpl) Create(ctx context.Context, request visimisipemda.MisiPemdaCreateRequest) (visimisipemda.MisiPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return visimisipemda.MisiPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Cek apakah urutan sudah digunakan
	exists, err := service.MisiPemdaRepository.CheckUrutanExists(ctx, tx, request.IdVisi, request.Urutan)
	if err != nil {
		return visimisipemda.MisiPemdaResponse{}, err
	}
	if exists {
		return visimisipemda.MisiPemdaResponse{}, fmt.Errorf("urutan %d sudah digunakan untuk visi ini", request.Urutan)
	}

	_, err = service.VisiPemdaRepository.FindById(ctx, tx, request.IdVisi)
	if err != nil {
		return visimisipemda.MisiPemdaResponse{}, err
	}

	misiPemda, err := service.MisiPemdaRepository.Create(ctx, tx, domain.MisiPemda{
		IdVisi:            request.IdVisi,
		Misi:              request.Misi,
		Urutan:            request.Urutan,
		TahunAwalPeriode:  request.TahunAwalPeriode,
		TahunAkhirPeriode: request.TahunAkhirPeriode,
		JenisPeriode:      request.JenisPeriode,
		Keterangan:        request.Keterangan,
	})
	if err != nil {
		return visimisipemda.MisiPemdaResponse{}, err
	}

	return helper.ToMisiPemdaResponse(misiPemda), nil
}

func (service *MisiPemdaServiceImpl) Update(ctx context.Context, request visimisipemda.MisiPemdaUpdateRequest) (visimisipemda.MisiPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return visimisipemda.MisiPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	// Cek apakah urutan sudah digunakan oleh misi lain
	exists, err := service.MisiPemdaRepository.CheckUrutanExistsExcept(ctx, tx, request.IdVisi, request.Urutan, request.Id)
	if err != nil {
		return visimisipemda.MisiPemdaResponse{}, err
	}
	if exists {
		return visimisipemda.MisiPemdaResponse{}, fmt.Errorf("urutan %d sudah digunakan untuk visi ini", request.Urutan)
	}

	visiPemda, err := service.VisiPemdaRepository.FindById(ctx, tx, request.IdVisi)
	if err != nil {
		return visimisipemda.MisiPemdaResponse{}, err
	}

	misiPemda, err := service.MisiPemdaRepository.Update(ctx, tx, domain.MisiPemda{
		Id:                request.Id,
		IdVisi:            visiPemda.Id,
		Misi:              request.Misi,
		Urutan:            request.Urutan,
		TahunAwalPeriode:  request.TahunAwalPeriode,
		TahunAkhirPeriode: request.TahunAkhirPeriode,
		JenisPeriode:      request.JenisPeriode,
		Keterangan:        request.Keterangan,
	})
	if err != nil {
		return visimisipemda.MisiPemdaResponse{}, err
	}

	return helper.ToMisiPemdaResponse(misiPemda), nil
}

func (service *MisiPemdaServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := service.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	return service.MisiPemdaRepository.Delete(ctx, tx, id)
}

func (service *MisiPemdaServiceImpl) FindAll(ctx context.Context, tahunAwal string, tahunAkhir string, jenisPeriode string) ([]visimisipemda.VisiPemdaRespons, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	// Ambil semua data misi
	misiPemdaList, err := service.MisiPemdaRepository.FindAll(ctx, tx, tahunAwal, tahunAkhir, jenisPeriode)
	if err != nil {
		return nil, err
	}

	// Buat map untuk mengelompokkan misi berdasarkan id_visi
	visiMisiMap := make(map[int]*visimisipemda.VisiPemdaRespons)

	// Iterasi setiap misi dan kelompokkan berdasarkan id_visi
	for _, misi := range misiPemdaList {
		// Jika id_visi belum ada di map, ambil data visi dan buat entry baru
		if _, exists := visiMisiMap[misi.IdVisi]; !exists {
			// Ambil data visi
			visiPemda, err := service.VisiPemdaRepository.FindByIdWithDefault(ctx, tx, misi.IdVisi)
			if err != nil {
				return nil, fmt.Errorf("gagal mengambil data visi: %v", err)
			}

			visiMisiMap[misi.IdVisi] = &visimisipemda.VisiPemdaRespons{
				IdVisi: visiPemda.Id,
				Visi:   visiPemda.Visi,
				Misi:   make([]visimisipemda.MisiPemdaResponse, 0),
			}
		}

		// Tambahkan misi ke array misi yang sesuai dengan id_visi
		misiResponse := visimisipemda.MisiPemdaResponse{
			Id:                misi.Id,
			Misi:              misi.Misi,
			TahunAwalPeriode:  misi.TahunAwalPeriode,
			TahunAkhirPeriode: misi.TahunAkhirPeriode,
			JenisPeriode:      misi.JenisPeriode,
			Keterangan:        misi.Keterangan,
		}
		visiMisiMap[misi.IdVisi].Misi = append(visiMisiMap[misi.IdVisi].Misi, misiResponse)
	}

	// Konversi map ke slice untuk response
	result := make([]visimisipemda.VisiPemdaRespons, 0, len(visiMisiMap))
	for _, visiMisi := range visiMisiMap {
		result = append(result, *visiMisi)
	}

	// Sort berdasarkan IdVisi
	sort.Slice(result, func(i, j int) bool {
		return result[i].IdVisi < result[j].IdVisi
	})

	return result, nil
}

func (service *MisiPemdaServiceImpl) FindById(ctx context.Context, id int) (visimisipemda.MisiPemdaResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		return visimisipemda.MisiPemdaResponse{}, err
	}
	defer helper.CommitOrRollback(tx)

	misiPemda, err := service.MisiPemdaRepository.FindById(ctx, tx, id)
	if err != nil {
		return visimisipemda.MisiPemdaResponse{}, err
	}

	return helper.ToMisiPemdaResponse(misiPemda), nil
}
