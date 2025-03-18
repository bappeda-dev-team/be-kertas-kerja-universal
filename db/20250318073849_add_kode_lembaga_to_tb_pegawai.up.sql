ALTER TABLE tb_pegawai ADD COLUMN kode_lembaga VARCHAR(255);
ALTER TABLE tb_pegawai ADD CONSTRAINT fk_kode_lembaga FOREIGN KEY (kode_lembaga) REFERENCES tb_lembaga(kode_lembaga) ON DELETE SET NULL;
