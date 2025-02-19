package domain

type SasaranPemda struct {
	Id            int
	TujuanPemdaId int
	SubtemaId     int
	NamaSubtema   string
	SasaranPemda  string
	JenisPohon    string
	PeriodeId     int
	Periode       Periode
	Indikator     []Indikator
}

type SasaranPemdaWithPokin struct {
	SubtematikId        int
	JenisPohon          string
	LevelPohon          int
	TematikId           int
	NamaTematik         string
	NamaSubtematik      string
	IdsasaranPemda      int
	SasaranPemda        string
	Keterangan          string
	IndikatorSubtematik []Indikator
	SasaranPemdaList    []SasaranPemda
}
