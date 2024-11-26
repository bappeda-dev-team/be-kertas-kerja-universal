package dataseeder

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
	"ekak_kabupaten_madiun/repository"
	"log"

	"github.com/google/uuid"
)

type PegawaiSeederImpl struct {
	DB                *sql.DB
	PegawaiRepository repository.PegawaiRepository
}

func NewPegawaiSeederImpl(db *sql.DB, pegawaiRepository repository.PegawaiRepository) *PegawaiSeederImpl {
	return &PegawaiSeederImpl{
		DB:                db,
		PegawaiRepository: pegawaiRepository,
	}
}

func (pegawai *PegawaiSeederImpl) Seed(ctx context.Context, tx *sql.Tx) error {
	// Cek dan buat pegawai pertama
	_, err := pegawai.PegawaiRepository.FindByNip(ctx, tx, "admin1")
	if err == sql.ErrNoRows {
		superAdmin1 := domainmaster.Pegawai{
			Id:          "ADMIN-" + uuid.New().String()[:4],
			NamaPegawai: "super admin satu",
			Nip:         "admin1",
			KodeOpd:     "",
		}

		_, err = pegawai.PegawaiRepository.Create(ctx, tx, superAdmin1)
		if err != nil {
			return err
		}
		log.Println("Pegawai pertama berhasil di-seed")
	}

	// Cek dan buat pegawai kedua
	_, err = pegawai.PegawaiRepository.FindByNip(ctx, tx, "admin2")
	if err == sql.ErrNoRows {
		superAdmin2 := domainmaster.Pegawai{
			Id:          "ADMIN-" + uuid.New().String()[:4],
			NamaPegawai: "super admin dua",
			Nip:         "admin2",
			KodeOpd:     "",
		}

		_, err = pegawai.PegawaiRepository.Create(ctx, tx, superAdmin2)
		if err != nil {
			return err
		}
		log.Println("Pegawai kedua berhasil di-seed")
	}

	return nil
}
