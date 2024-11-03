package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/jabatan"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type JabatanControllerImpl struct {
	jabatanService service.JabatanService
}

func NewJabatanControllerImpl(jabatanService service.JabatanService) *JabatanControllerImpl {
	return &JabatanControllerImpl{
		jabatanService: jabatanService,
	}
}

func (controller *JabatanControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Decode request body
	jabatanCreateRequest := jabatan.JabatanCreateRequest{}
	helper.ReadFromRequestBody(request, &jabatanCreateRequest)

	jabatanResponse := controller.jabatanService.Create(request.Context(), jabatanCreateRequest)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "success create data jabatan",
		Data:   jabatanResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *JabatanControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	jabatanUpdateRequest := jabatan.JabatanUpdateRequest{}
	helper.ReadFromRequestBody(request, &jabatanUpdateRequest)

	jabatanResponse := controller.jabatanService.Update(request.Context(), jabatanUpdateRequest)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "success update data jabatan",
		Data:   jabatanResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *JabatanControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	jabatanId := params.ByName("jabatanId")
	controller.jabatanService.Delete(request.Context(), jabatanId)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "success delete data jabatan",
		Data:   nil,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *JabatanControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	jabatanId := params.ByName("jabatanId")
	jabatanResponse, err := controller.jabatanService.FindById(request.Context(), jabatanId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   404,
			Status: "not found",
			Data:   nil,
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "success get data jabatan by id",
		Data:   jabatanResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *JabatanControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")

	jabatanResponse, err := controller.jabatanService.FindAll(request.Context(), kodeOpd, tahun)
	if err != nil {
		panic(err)
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "success get data jabatan",
		Data:   jabatanResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
