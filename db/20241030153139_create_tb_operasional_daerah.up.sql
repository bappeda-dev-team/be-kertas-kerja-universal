CREATE TABLE tb_operasional_daerah (
    id VARCHAR(36) NOT NULL,
    kode_opd VARCHAR(10) NOT NULL,
    nama_opd VARCHAR(255) NOT NULL,
    singkatan VARCHAR(10) NOT NULL,
    alamat TEXT NOT NULL,
    telepon VARCHAR(20) NOT NULL,
    fax VARCHAR(20) NOT NULL,
    email VARCHAR(255) NOT NULL,
    website VARCHAR(255) NOT NULL,
    nama_kepala_opd VARCHAR(255) NOT NULL,
    nip_kepala_opd VARCHAR(20) NOT NULL,
    pangkat_kepala VARCHAR(255) NOT NULL,
    id_lembaga VARCHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB;