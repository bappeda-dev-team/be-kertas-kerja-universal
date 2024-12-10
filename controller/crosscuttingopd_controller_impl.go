package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"ekak_kabupaten_madiun/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type CrosscuttingOpdControllerImpl struct {
	CrosscuttingOpdService service.CrosscuttingOpdService
}

func NewCrosscuttingOpdControllerImpl(crosscuttingOpdService service.CrosscuttingOpdService) *CrosscuttingOpdControllerImpl {
	return &CrosscuttingOpdControllerImpl{
		CrosscuttingOpdService: crosscuttingOpdService,
	}
}

func (controller *CrosscuttingOpdControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	parentId, err := strconv.Atoi(params.ByName("parentId"))
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   "Invalid parent ID",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	crosscuttingCreateRequest := pohonkinerja.CrosscuttingOpdCreateRequest{}
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&crosscuttingCreateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	crosscuttingResponse, err := controller.CrosscuttingOpdService.Create(request.Context(), crosscuttingCreateRequest, parentId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   500,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   crosscuttingResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *CrosscuttingOpdControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Ambil crosscuttingId dari params
	crosscuttingId, err := strconv.Atoi(params.ByName("crosscuttingId"))
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   "Invalid crosscutting ID",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Decode request body
	crosscuttingUpdateRequest := pohonkinerja.CrosscuttingOpdUpdateRequest{}
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&crosscuttingUpdateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Set ID dari params ke request
	crosscuttingUpdateRequest.Id = crosscuttingId

	// Panggil service
	crosscuttingResponse, err := controller.CrosscuttingOpdService.Update(request.Context(), crosscuttingUpdateRequest)
	if err != nil {
		// Handle specific errors
		if err.Error() == "crosscutting tidak ditemukan" {
			webResponse := web.WebResponse{
				Code:   404,
				Status: "NOT FOUND",
				Data:   err.Error(),
			}
			helper.WriteToResponseBody(writer, webResponse)
			return
		}
		if err.Error() == "kode OPD hanya dapat diubah saat status crosscutting_menunggu" {
			webResponse := web.WebResponse{
				Code:   400,
				Status: "BAD REQUEST",
				Data:   err.Error(),
			}
			helper.WriteToResponseBody(writer, webResponse)
			return
		}

		// Default error response
		webResponse := web.WebResponse{
			Code:   500,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Success response
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   crosscuttingResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *CrosscuttingOpdControllerImpl) FindAllByParent(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	parentId, err := strconv.Atoi(params.ByName("parentId"))
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   "Invalid parent ID",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	crosscuttingResponses, err := controller.CrosscuttingOpdService.FindAllByParent(request.Context(), parentId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   500,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   crosscuttingResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *CrosscuttingOpdControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

}

func (controller *CrosscuttingOpdControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	parentId, err := strconv.Atoi(params.ByName("parentId"))
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   "Invalid parent ID",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Panggil service untuk mendapatkan data
	crosscuttingResponses, err := controller.CrosscuttingOpdService.FindAllByParent(request.Context(), parentId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   500,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Jika data kosong, return array kosong
	if len(crosscuttingResponses) == 0 {
		webResponse := web.WebResponse{
			Code:   200,
			Status: "OK",
			Data:   []pohonkinerja.CrosscuttingOpdResponse{},
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Return response sukses dengan data
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   crosscuttingResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *CrosscuttingOpdControllerImpl) ApproveOrReject(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	crosscuttingId, err := strconv.Atoi(params.ByName("crosscuttingId"))
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   "Invalid crosscutting ID",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	parentId, err := strconv.Atoi(params.ByName("parentId"))
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   "Invalid parent ID",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	var approveRequest pohonkinerja.CrosscuttingApproveRequest
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&approveRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   400,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Set parent ID dari parameter URL
	approveRequest.ParentId = parentId

	response, err := controller.CrosscuttingOpdService.ApproveOrReject(request.Context(), crosscuttingId, approveRequest)
	if err != nil {
		if err.Error() == "crosscutting sudah disetujui atau ditolak" {
			webResponse := web.WebResponse{
				Code:   400,
				Status: "BAD REQUEST",
				Data:   err.Error(),
			}
			helper.WriteToResponseBody(writer, webResponse)
			return
		}

		webResponse := web.WebResponse{
			Code:   500,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   response,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
