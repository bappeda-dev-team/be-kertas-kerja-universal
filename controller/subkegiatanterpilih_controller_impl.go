package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/subkegiatan"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type SubKegiatanTerpilihControllerImpl struct {
	SubKegiatanTerpilihService service.SubKegiatanTerpilihService
}

func NewSubKegiatanTerpilihControllerImpl(subKegiatanTerpilihService service.SubKegiatanTerpilihService) *SubKegiatanTerpilihControllerImpl {
	return &SubKegiatanTerpilihControllerImpl{
		SubKegiatanTerpilihService: subKegiatanTerpilihService,
	}
}

func (controller *SubKegiatanTerpilihControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	subKegiatanTerpilihCreateRequest := subkegiatan.SubKegiatanTerpilihCreateRequest{}
	helper.ReadFromRequestBody(request, &subKegiatanTerpilihCreateRequest)

	// Ambil rekin dari parameter URL
	rekinId := params.ByName("rencana_kinerja_id")
	subKegiatanTerpilihCreateRequest.RencanaKinerjaId = rekinId

	// Panggil service untuk membuat SubKegiatanTerpilih
	subKegiatanTerpilihResponse, err := controller.SubKegiatanTerpilihService.Create(request.Context(), subKegiatanTerpilihCreateRequest)
	if err != nil {
		webResponse := web.WebSubKegiatanTerpilihResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	// Kirim respons sukses
	webResponse := web.WebSubKegiatanTerpilihResponse{
		Code:   http.StatusCreated,
		Status: "success create sub kegiatan terpilih",
		Data:   subKegiatanTerpilihResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *SubKegiatanTerpilihControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	subKegiatanId := params.ByName("subkegiatan_id")

	err := controller.SubKegiatanTerpilihService.Delete(request.Context(), subKegiatanId)
	if err != nil {
		webResponse := web.WebSubKegiatanTerpilihResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebSubKegiatanTerpilihResponse{
		Code:   http.StatusOK,
		Status: "success delete sub kegiatan terpilih",
	}
	helper.WriteToResponseBody(writer, webResponse)
}
