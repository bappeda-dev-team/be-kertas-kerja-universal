package repository

import (
	"context"
	"database/sql"
	"ekak_kabupaten_madiun/model/domain/domainmaster"
	"fmt"
	"strings"
)

type UrusanRepositoryImpl struct {
}

func NewUrusanRepositoryImpl() *UrusanRepositoryImpl {
	return &UrusanRepositoryImpl{}
}

func (repository *UrusanRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, urusan domainmaster.Urusan) (domainmaster.Urusan, error) {
	script := "INSERT INTO tb_urusan(id, kode_urusan, nama_urusan) VALUES (?, ?, ?)"
	_, err := tx.ExecContext(ctx, script, urusan.Id, urusan.KodeUrusan, urusan.NamaUrusan)
	if err != nil {
		return urusan, err
	}

	return urusan, nil
}

func (repository *UrusanRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, urusan domainmaster.Urusan) (domainmaster.Urusan, error) {
	script := "UPDATE tb_urusan SET kode_urusan = ?, nama_urusan = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, urusan.KodeUrusan, urusan.NamaUrusan, urusan.Id)
	if err != nil {
		return urusan, err
	}

	return urusan, nil
}

func (repository *UrusanRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domainmaster.Urusan, error) {
	script := "SELECT id, kode_urusan, nama_urusan, created_at FROM tb_urusan"
	rows, err := tx.QueryContext(ctx, script)
	if err != nil {
		return []domainmaster.Urusan{}, err
	}

	defer rows.Close()

	var urusans []domainmaster.Urusan
	for rows.Next() {
		urusan := domainmaster.Urusan{}
		rows.Scan(&urusan.Id, &urusan.KodeUrusan, &urusan.NamaUrusan, &urusan.CreatedAt)
		urusans = append(urusans, urusan)
	}

	return urusans, nil
}

func (repository *UrusanRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (domainmaster.Urusan, error) {
	script := "SELECT id, kode_urusan, nama_urusan, created_at FROM tb_urusan WHERE id = ?"
	rows, err := tx.QueryContext(ctx, script, id)
	if err != nil {
		return domainmaster.Urusan{}, err
	}
	defer rows.Close()

	urusan := domainmaster.Urusan{}

	if rows.Next() {
		err := rows.Scan(&urusan.Id, &urusan.KodeUrusan, &urusan.NamaUrusan, &urusan.CreatedAt)
		if err != nil {
			return domainmaster.Urusan{}, err
		}
	} else {
		return domainmaster.Urusan{}, fmt.Errorf("urusan dengan id %s tidak ditemukan", id)
	}

	return urusan, nil
}

func (repository *UrusanRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	script := "DELETE FROM tb_urusan WHERE id = ?"
	_, err := tx.ExecContext(ctx, script, id)
	if err != nil {
		return err
	}

	return nil
}

func (repository *UrusanRepositoryImpl) FindByKodeOpd(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domainmaster.Urusan, error) {
	kodeUrusans := make([]string, 0)

	if len(kodeOpd) >= 1 {
		kode1 := kodeOpd[0:1]
		if kode1 != "0" {
			kodeUrusans = append(kodeUrusans, kode1)
		}
	}

	if len(kodeOpd) >= 6 {
		kode2 := kodeOpd[5:6]
		if kode2 != "0" {
			kodeUrusans = append(kodeUrusans, kode2)
		}
	}

	if len(kodeOpd) >= 11 {
		kode3 := kodeOpd[10:11]
		if kode3 != "0" {
			kodeUrusans = append(kodeUrusans, kode3)
		}
	}

	if len(kodeUrusans) == 0 {
		return []domainmaster.Urusan{}, nil
	}

	query := "SELECT id, kode_urusan, nama_urusan FROM tb_urusan WHERE kode_urusan IN ("
	params := make([]interface{}, len(kodeUrusans))
	for i := range kodeUrusans {
		if i > 0 {
			query += ","
		}
		query += "?"
		params[i] = kodeUrusans[i]
	}
	query += ")"

	rows, err := tx.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urusans []domainmaster.Urusan
	for rows.Next() {
		urusan := domainmaster.Urusan{}
		err := rows.Scan(&urusan.Id, &urusan.KodeUrusan, &urusan.NamaUrusan)
		if err != nil {
			return nil, err
		}
		urusans = append(urusans, urusan)
	}

	return urusans, nil
}

func (repository *UrusanRepositoryImpl) FindUrusanAndBidangByKodeOpd(ctx context.Context, tx *sql.Tx, kodeOpd string) ([]domainmaster.Urusan, error) {
	kodeUrusans := make([]string, 0)
	kodeBidangUrusans := make([]string, 0)

	parts := strings.Split(kodeOpd, ".")
	if len(parts) >= 4 {
		// Urusan dan Bidang Urusan pertama (1 dan 1.01)
		if parts[0] != "0" {
			kodeUrusans = append(kodeUrusans, parts[0])
			kodeBidang := parts[0] + "." + parts[1]
			if parts[1] != "00" {
				kodeBidangUrusans = append(kodeBidangUrusans, kodeBidang)
			}
		}

		// Urusan dan Bidang Urusan kedua (2 dan 2.22)
		if parts[2] != "0" {
			kodeUrusans = append(kodeUrusans, parts[2])
			kodeBidang := parts[2] + "." + parts[3]
			if parts[3] != "00" {
				kodeBidangUrusans = append(kodeBidangUrusans, kodeBidang)
			}
		}

		// Urusan ketiga jika ada
		if len(parts) >= 5 && parts[4] != "0" {
			kodeUrusans = append(kodeUrusans, parts[4])
			if len(parts) >= 6 && parts[5] != "00" {
				kodeBidang := parts[4] + "." + parts[5]
				kodeBidangUrusans = append(kodeBidangUrusans, kodeBidang)
			}
		}
	}

	// Debug: cetak kode yang ditemukan
	fmt.Printf("Kode OPD: %s\n", kodeOpd)
	fmt.Printf("Kode Urusan: %v\n", kodeUrusans)
	fmt.Printf("Kode Bidang Urusan: %v\n", kodeBidangUrusans)

	if len(kodeUrusans) == 0 {
		return []domainmaster.Urusan{}, nil
	}

	// Buat placeholders untuk IN clause
	urusanPlaceholders := createPlaceholders(len(kodeUrusans))
	bidangPlaceholders := createPlaceholders(len(kodeBidangUrusans))

	// Coba query manual dulu untuk memastikan data ada
	checkQuery := `
        SELECT COUNT(*) 
        FROM tb_bidang_urusan 
        WHERE kode_bidang_urusan IN (` + createPlaceholders(len(kodeBidangUrusans)) + `)`

	var count int
	err := tx.QueryRowContext(ctx, checkQuery, interfaceSlice(kodeBidangUrusans)...).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("error checking bidang urusan: %v", err)
	}
	fmt.Printf("Jumlah bidang urusan ditemukan: %d\n", count)

	// Query utama
	query := fmt.Sprintf(`
	SELECT DISTINCT
		u.id,
		u.kode_urusan,
		u.nama_urusan,
		bu.kode_bidang_urusan,
		bu.nama_bidang_urusan
	FROM 
		tb_urusan u
		INNER JOIN tb_bidang_urusan bu ON u.kode_urusan = LEFT(bu.kode_bidang_urusan, 1)
	WHERE 
		u.kode_urusan IN (%s)
		AND bu.kode_bidang_urusan IN (%s)
	ORDER BY 
		u.kode_urusan, bu.kode_bidang_urusan
`, urusanPlaceholders, bidangPlaceholders)

	// Gabungkan parameter
	params := make([]interface{}, 0)
	params = append(params, interfaceSlice(kodeUrusans)...)
	params = append(params, interfaceSlice(kodeBidangUrusans)...)

	rows, err := tx.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var result []domainmaster.Urusan
	urusanMap := make(map[string]*domainmaster.Urusan)

	for rows.Next() {
		var (
			id, kodeUrusan, namaUrusan         string
			kodeBidangUrusan, namaBidangUrusan string
		)

		err := rows.Scan(&id, &kodeUrusan, &namaUrusan, &kodeBidangUrusan, &namaBidangUrusan)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Debug
		fmt.Printf("Scanned data - urusan: %s, bidang: %s\n", kodeUrusan, kodeBidangUrusan)

		// Cek apakah urusan sudah ada di map
		if existingUrusan, exists := urusanMap[kodeUrusan]; exists {
			// Tambahkan bidang urusan ke urusan yang sudah ada
			existingUrusan.BidangUrusan = append(existingUrusan.BidangUrusan, domainmaster.BidangUrusan{
				KodeBidangUrusan: kodeBidangUrusan,
				NamaBidangUrusan: namaBidangUrusan,
			})
		} else {
			// Buat urusan baru
			newUrusan := &domainmaster.Urusan{
				Id:         id,
				KodeUrusan: kodeUrusan,
				NamaUrusan: namaUrusan,
				BidangUrusan: []domainmaster.BidangUrusan{
					{
						KodeBidangUrusan: kodeBidangUrusan,
						NamaBidangUrusan: namaBidangUrusan,
					},
				},
			}
			urusanMap[kodeUrusan] = newUrusan
			result = append(result, *newUrusan)
		}

		// Debug
		fmt.Printf("Current urusan %s has %d bidang\n",
			kodeUrusan, len(urusanMap[kodeUrusan].BidangUrusan))
	}

	// Debug final result
	for _, u := range result {
		fmt.Printf("Final result - urusan %s has %d bidang\n",
			u.KodeUrusan, len(u.BidangUrusan))
		for _, b := range u.BidangUrusan {
			fmt.Printf("  Bidang: %s - %s\n", b.KodeBidangUrusan, b.NamaBidangUrusan)
		}
	}

	return result, nil
}

// Helper function untuk membuat placeholder
func createPlaceholders(n int) string {
	if n <= 0 {
		return ""
	}
	placeholders := make([]string, n)
	for i := range placeholders {
		placeholders[i] = "?"
	}
	return strings.Join(placeholders, ",")
}

// Helper function untuk mengkonversi slice string ke slice interface
func interfaceSlice(slice []string) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}
