package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

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
	jenisPeriode := params.ByName("jenis_periode")

	sasaranOpdResponse, err := controller.SasaranOpdService.FindAll(request.Context(), KodeOpd, tahunAwal, tahunAkhir, jenisPeriode)
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
			Status: "get all sasaran opd",
			Data:   sasaranOpdResponse,
		}
		helper.WriteToResponseBody(writer, webResponse)
	}
}

func (controller *SasaranOpdControllerImpl) FindByIdRencanaKinerja(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	idRencanaKinerja := params.ByName("id_rencana_kinerja")

	sasaranOpdResponse, err := controller.SasaranOpdService.FindByIdRencanaKinerja(request.Context(), idRencanaKinerja)
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
			Status: "get all sasaran opd by id rencana kinerja",
			Data:   sasaranOpdResponse,
		}
		helper.WriteToResponseBody(writer, webResponse)
	}
}

func (controller *SasaranOpdControllerImpl) FindIdPokinSasaran(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	idPokinStr := params.ByName("id")
	idPokin, err := strconv.Atoi(idPokinStr)
	helper.PanicIfError(err)

	sasaranOpdResponse, err := controller.SasaranOpdService.FindIdPokinSasaran(request.Context(), idPokin)
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
			Status: "get all sasaran opd by id rencana kinerja",
			Data:   sasaranOpdResponse,
		}
		helper.WriteToResponseBody(writer, webResponse)
	}
}
