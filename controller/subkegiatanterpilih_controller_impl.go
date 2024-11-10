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

func (controller *SubKegiatanTerpilihControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	subKegiatanTerpilihUpdateRequest := subkegiatan.SubKegiatanTerpilihUpdateRequest{}
	helper.ReadFromRequestBody(request, &subKegiatanTerpilihUpdateRequest)
	rencanaKinerjaId := params.ByName("rencana_kinerja_id")
	subKegiatanTerpilihUpdateRequest.Id = rencanaKinerjaId

	// Panggil service untuk membuat SubKegiatanTerpilih
	subKegiatanTerpilihResponse, err := controller.SubKegiatanTerpilihService.Update(request.Context(), subKegiatanTerpilihUpdateRequest)
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
		Status: "success add sub kegiatan terpilih",
		Data:   subKegiatanTerpilihResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *SubKegiatanTerpilihControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeSubKegiatan := params.ByName("kode_subkegiatan")
	rencanaKinerjaId := params.ByName("rencana_kinerja_id")

	err := controller.SubKegiatanTerpilihService.Delete(request.Context(), rencanaKinerjaId, kodeSubKegiatan)
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

func (controller *SubKegiatanTerpilihControllerImpl) FindByKodeSubKegiatan(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeSubKegiatan := params.ByName("kode_subkegiatan")

	subKegiatanTerpilihResponse, err := controller.SubKegiatanTerpilihService.FindByKodeSubKegiatan(request.Context(), kodeSubKegiatan)
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
		Status: "success find sub kegiatan terpilih",
		Data:   subKegiatanTerpilihResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
