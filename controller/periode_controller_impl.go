package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/periodetahun"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type PeriodeControllerImpl struct {
	PeriodeService service.PeriodeService
}

func NewPeriodeControllerImpl(periodeService service.PeriodeService) *PeriodeControllerImpl {
	return &PeriodeControllerImpl{PeriodeService: periodeService}
}

func (controller *PeriodeControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	periodeCreateRequest := periodetahun.PeriodeCreateRequest{}
	helper.ReadFromRequestBody(request, &periodeCreateRequest)

	periodeResponse, err := controller.PeriodeService.Create(request.Context(), periodeCreateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   periodeResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PeriodeControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	periodeUpdateRequest := periodetahun.PeriodeUpdateRequest{}
	helper.ReadFromRequestBody(request, &periodeUpdateRequest)
	id := params.ByName("id")
	periodeUpdateRequest.Id, _ = strconv.Atoi(id)

	periodeResponse, err := controller.PeriodeService.Update(request.Context(), periodeUpdateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   periodeResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PeriodeControllerImpl) FindByTahun(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tahun := params.ByName("tahun")
	periodeResponse, err := controller.PeriodeService.FindByTahun(request.Context(), tahun)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "failed find data periode",
			Data:   nil,
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	webResponse := web.WebResponse{
		Code:   200,
		Status: "success find data periode",
		Data:   periodeResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
