package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type IkuControllerImpl struct {
	IkuService service.IkuService
}

func NewIkuControllerImpl(ikuService service.IkuService) *IkuControllerImpl {
	return &IkuControllerImpl{
		IkuService: ikuService,
	}
}

func (controller *IkuControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tahun := params.ByName("tahun")
	if tahun == "" {
		// Handle error jika tahun tidak ada
		helper.WriteToResponseBody(writer, "Tahun harus diisi")
		return
	}

	ikuResponses, err := controller.IkuService.FindAll(request.Context(), tahun)
	if err != nil {
		helper.WriteToResponseBody(writer, err.Error())
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   ikuResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
