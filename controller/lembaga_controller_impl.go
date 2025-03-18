package controller

import (
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/lembaga"
	"ekak_kabupaten_madiun/service"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type LembagaControllerImpl struct {
	LembagaService service.LembagaService
}

func NewLembagaControllerImpl(lembagaService service.LembagaService) *LembagaControllerImpl {
	return &LembagaControllerImpl{
		LembagaService: lembagaService,
	}
}

func (controller *LembagaControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	lembagaCreateRequest := lembaga.LembagaCreateRequest{}
	helper.ReadFromRequestBody(request, &lembagaCreateRequest)

	lembagaResponse, err := controller.LembagaService.Create(request.Context(), lembagaCreateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   500,
			Status: "Internal Server Error",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   lembagaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *LembagaControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	lembagaId := params.ByName("id")

	lembagaUpdateRequest := lembaga.LembagaUpdateRequest{}
	helper.ReadFromRequestBody(request, &lembagaUpdateRequest)

	lembagaUpdateRequest.Id = lembagaId

	lembagaResponse, err := controller.LembagaService.Update(request.Context(), lembagaUpdateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   500,
			Status: "Internal Server Error",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   lembagaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *LembagaControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	lembagaId := params.ByName("id")

	controller.LembagaService.Delete(request.Context(), lembagaId)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}
	helper.WriteToResponseBody(writer, webResponse)

}

func (controller *LembagaControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	lembagaId := params.ByName("id")

	lembagaResponse, err := controller.LembagaService.FindById(request.Context(), lembagaId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   500,
			Status: "Internal Server Error",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   lembagaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (c *LembagaControllerImpl) FindByKode(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	kodeLembaga := params.ByName("kode_lembaga")

	if kodeLembaga == "" {
		webResponse := web.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  http.StatusText(http.StatusBadRequest),
			Message: "Missing requried parameter",
		}
		helper.WriteToResponseBody(w, webResponse)
		return
	}

	lembagaResponse, err := c.LembagaService.FindByKode(r.Context(), kodeLembaga)
	if err != nil {
		if err == sql.ErrNoRows {
			webResponse := web.WebResponse{
				Code:    http.StatusNoContent,
				Status:  http.StatusText(http.StatusNoContent),
				Message: fmt.Sprintf("Lembaga dengan kode: %s tidak ditemukan", kodeLembaga),
			}
			helper.WriteToResponseBody(w, webResponse)
			return
		}
		webResponse := web.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  http.StatusText(http.StatusInternalServerError),
			Message: "Terjadi kesalahan pada server. akan segera ditangani.",
		}
		helper.WriteToResponseBody(w, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Message: "Detail Lembaga by kode",
		Data:    lembagaResponse,
	}
	helper.WriteToResponseBody(w, webResponse)
}

func (controller *LembagaControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	lembagaResponse, err := controller.LembagaService.FindAll(request.Context())
	if err != nil {
		webResponse := web.WebResponse{
			Code:   500,
			Status: "Internal Server Error",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   lembagaResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
