package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type SasaranOpdControllerImpl struct {
	SasaranOpdService service.SasaranOpdService
}

func NewSasaranOpdControllerImpl(SasaranOpdService service.SasaranOpdService) *SasaranOpdControllerImpl {
	return &SasaranOpdControllerImpl{
		SasaranOpdService: SasaranOpdService,
	}
}

func (controller *SasaranOpdControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	KodeOpd := params.ByName("kode_opd")
	tahunAwal := params.ByName("tahun_awal")
	tahunAkhir := params.ByName("tahun_akhir")

	sasaranOpdResponse, err := controller.SasaranOpdService.FindAll(request.Context(), KodeOpd, tahunAwal, tahunAkhir)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD_REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
	} else {
		webResponse := web.WebResponse{
			Code:   200,
			Status: "OK",
			Data:   sasaranOpdResponse,
		}
		helper.WriteToResponseBody(writer, webResponse)
	}
}
