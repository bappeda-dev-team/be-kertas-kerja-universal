package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/pohonkinerja"
	"ekak_kabupaten_madiun/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type PohonKinerjaAdminControllerImpl struct {
	pohonKinerjaAdminService service.PohonKinerjaAdminService
}

func NewPohonKinerjaAdminControllerImpl(pohonKinerjaAdminService service.PohonKinerjaAdminService) *PohonKinerjaAdminControllerImpl {
	return &PohonKinerjaAdminControllerImpl{pohonKinerjaAdminService: pohonKinerjaAdminService}
}

func (controller *PohonKinerjaAdminControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Decode request body
	pohonKinerjaCreateRequest := pohonkinerja.PohonKinerjaAdminCreateRequest{}
	helper.ReadFromRequestBody(request, &pohonKinerjaCreateRequest)

	// Panggil service create
	pohonKinerjaResponse, err := controller.pohonKinerjaAdminService.Create(request.Context(), pohonKinerjaCreateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Buat response sukses
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   pohonKinerjaResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PohonKinerjaAdminControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Decode request body
	pohonKinerjaUpdateRequest := pohonkinerja.PohonKinerjaAdminUpdateRequest{}
	helper.ReadFromRequestBody(request, &pohonKinerjaUpdateRequest)

	// Ambil ID dari parameter URL
	pohonKinerjaId := params.ByName("pohonKinerjaId")
	pohonKinerjaUpdateRequest.Id, _ = strconv.Atoi(pohonKinerjaId)

	// Panggil service update
	pohonKinerjaResponse := controller.pohonKinerjaAdminService.Update(request.Context(), pohonKinerjaUpdateRequest)
	// if err != nil {
	// 	webResponse := web.WebResponse{
	// 		Code:   http.StatusBadRequest,
	// 		Status: "BAD REQUEST",
	// 		Data:   err.Error(),
	// 	}
	// 	helper.WriteToResponseBody(writer, webResponse)
	// 	return
	// }

	// Buat response sukses
	webResponse := web.WebResponse{
		Code:   200,
		Status: "Success Update Pohon Kinerja",
		Data:   pohonKinerjaResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PohonKinerjaAdminControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Ambil ID dari parameter URL
	pohonKinerjaId := params.ByName("pohonKinerjaId")
	id, err := strconv.Atoi(pohonKinerjaId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   "ID tidak valid",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Panggil service delete
	err = controller.pohonKinerjaAdminService.Delete(request.Context(), id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Buat response sukses
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   "Pohon Kinerja berhasil dihapus",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PohonKinerjaAdminControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Ambil ID dari parameter URL
	pohonKinerjaId := params.ByName("id")
	id, err := strconv.Atoi(pohonKinerjaId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   "ID tidak valid",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Panggil service findById
	result, err := controller.pohonKinerjaAdminService.FindById(request.Context(), id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusNotFound,
			Status: "NOT FOUND",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Buat response sukses
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   result,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PohonKinerjaAdminControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tahun := params.ByName("tahun")

	// Panggil service findAll
	result, err := controller.pohonKinerjaAdminService.FindAll(request.Context(), tahun)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Buat response sukses
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   result,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *PohonKinerjaAdminControllerImpl) FindSubTematik(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tahun := params.ByName("tahun")

	// Panggil service findAll
	result, err := controller.pohonKinerjaAdminService.FindSubTematik(request.Context(), tahun)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Buat response sukses
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   result,
	}

	helper.WriteToResponseBody(writer, webResponse)

}
