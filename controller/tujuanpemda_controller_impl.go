package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/tujuanpemda"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type TujuanPemdaControllerImpl struct {
	TujuanPemdaService service.TujuanPemdaService
}

func NewTujuanPemdaControllerImpl(tujuanPemdaService service.TujuanPemdaService) *TujuanPemdaControllerImpl {
	return &TujuanPemdaControllerImpl{
		TujuanPemdaService: tujuanPemdaService,
	}
}

func (controller *TujuanPemdaControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Decode request body
	tujuanPemdaCreateRequest := tujuanpemda.TujuanPemdaCreateRequest{}
	helper.ReadFromRequestBody(request, &tujuanPemdaCreateRequest)

	// Panggil service create
	tujuanPemdaResponse, err := controller.TujuanPemdaService.Create(request.Context(), tujuanPemdaCreateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusCreated,
		Status: "success create tujuan pemda",
		Data:   tujuanPemdaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *TujuanPemdaControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tujuanPemdaUpdateRequest := tujuanpemda.TujuanPemdaUpdateRequest{}
	helper.ReadFromRequestBody(request, &tujuanPemdaUpdateRequest)

	id := params.ByName("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   "Invalid ID format",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	tujuanPemdaUpdateRequest.Id = idInt

	// Panggil service update
	tujuanPemdaResponse, err := controller.TujuanPemdaService.Update(request.Context(), tujuanPemdaUpdateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "success update tujuan pemda",
		Data:   tujuanPemdaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *TujuanPemdaControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tujuanPemdaId := params.ByName("id")
	id, err := strconv.Atoi(tujuanPemdaId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   "Invalid ID format",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	err = controller.TujuanPemdaService.Delete(request.Context(), id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *TujuanPemdaControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tujuanPemdaId := params.ByName("id")

	id, err := strconv.Atoi(tujuanPemdaId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   "Invalid ID format",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	tujuanPemdaResponse, err := controller.TujuanPemdaService.FindById(request.Context(), id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusNotFound,
			Status: "NOT FOUND",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   tujuanPemdaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *TujuanPemdaControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tahun := params.ByName("tahun")
	tujuanPemdaResponses, err := controller.TujuanPemdaService.FindAll(request.Context(), tahun)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   tujuanPemdaResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *TujuanPemdaControllerImpl) UpdatePeriode(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tujuanPemdaUpdateRequest := tujuanpemda.TujuanPemdaUpdateRequest{}
	helper.ReadFromRequestBody(request, &tujuanPemdaUpdateRequest)

	id := params.ByName("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   "Invalid ID format",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	tujuanPemdaUpdateRequest.Id = idInt

	// Panggil service update
	tujuanPemdaResponse, err := controller.TujuanPemdaService.UpdatePeriode(request.Context(), tujuanPemdaUpdateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "success update periode tujuan pemda",
		Data:   tujuanPemdaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *TujuanPemdaControllerImpl) FindAllWithPokin(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tahun := params.ByName("tahun")
	tujuanPemdaResponses, err := controller.TujuanPemdaService.FindAllWithPokin(request.Context(), tahun)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   tujuanPemdaResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
