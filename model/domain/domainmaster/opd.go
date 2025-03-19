package domainmaster

type Opd struct {
	Id            string
	KodeOpd       string
	NamaOpd       string
	Singkatan     string
	Alamat        string
	Telepon       string
	Fax           string
	Email         string
	Website       string
	NamaKepalaOpd string
	NIPKepalaOpd  string
	PangkatKepala string
	NamaAdmin        string
	NoWaAdmin        string
	IdLembaga     string
}

type OpdWithBidangUrusan struct {
	Id               string
	KodeOpd          string
	NamaOpd          string
	KodeUrusan1 string
	NamaUrusan1 string
	KodeUrusan2 string
	NamaUrusan2 string
	KodeUrusan3 string
	NamaUrusan3 string
	KodeBidangUrusan1 string
	NamaBidangUrusan1 string
	KodeBidangUrusan2 string
	NamaBidangUrusan2 string
	KodeBidangUrusan3 string
	NamaBidangUrusan3 string
	NamaKepalaOpd    string
	NIPKepalaOpd     string
	PangkatKepala    string
	NamaAdmin        string
	NoWaAdmin        string
}
