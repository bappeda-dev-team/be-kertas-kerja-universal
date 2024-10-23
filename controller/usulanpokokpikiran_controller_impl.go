package controller

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/model/web/usulan"
	"ekak_kabupaten_madiun/service"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type UsulanPokokPikiranControllerImpl struct {
	UsulanPokokPikiranService service.UsulanPokokPikiranService
}

func NewUsulanPokokPikiranControllerImpl(usulanPokokPikiranService service.UsulanPokokPikiranService) *UsulanPokokPikiranControllerImpl {
	return &UsulanPokokPikiranControllerImpl{
		UsulanPokokPikiranService: usulanPokokPikiranService,
	}
}

func (controller *UsulanPokokPikiranControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	usulanPokokPikiranCreateRequest := usulan.UsulanPokokPikiranCreateRequest{}
	helper.ReadFromRequestBody(request, &usulanPokokPikiranCreateRequest)

	pegawaiID := params.ByName("pegawai_id")
	if pegawaiID == "" {
		pegawaiID = usulanPokokPikiranCreateRequest.PegawaiId
	}

	if pegawaiID == "" {
		webResponse := web.WebUsulanPokokPikiranResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   "Invalid pegawai_id: tidak ditemukan di parameter URL atau body request",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	usulanPokokPikiranCreateRequest.PegawaiId = pegawaiID

	usulanPokokPikiranResponse, err := controller.UsulanPokokPikiranService.Create(request.Context(), usulanPokokPikiranCreateRequest)
	if err != nil {
		webResponse := web.WebUsulanPokokPikiranResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebUsulanPokokPikiranResponse{
		Code:   http.StatusOK,
		Status: "berhasil membuat usulan pokok pikiran",
		Data:   usulanPokokPikiranResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UsulanPokokPikiranControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	usulanPokokPikiranUpdateRequest := usulan.UsulanPokokPikiranUpdateRequest{}
	helper.ReadFromRequestBody(request, &usulanPokokPikiranUpdateRequest)

	idUsulan := params.ByName("id")
	if idUsulan == "" {
		webResponse := web.WebUsulanPokokPikiranResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   "Invalid id usulan parameter",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}
	usulanPokokPikiranUpdateRequest.Id = idUsulan

	pegawaiID := params.ByName("pegawai_id")
	if pegawaiID != "" {
		if usulanPokokPikiranUpdateRequest.PegawaiId != pegawaiID {
			webResponse := web.WebUsulanPokokPikiranResponse{
				Code:   http.StatusForbidden,
				Status: "FORBIDDEN",
				Data:   "Tidak dapat mengedit usulan pegawai lain",
			}
			helper.WriteToResponseBody(writer, webResponse)
			return
		}
	}

	usulanPokokPikiranResponse, err := controller.UsulanPokokPikiranService.Update(request.Context(), usulanPokokPikiranUpdateRequest)
	if err != nil {
		webResponse := web.WebUsulanPokokPikiranResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebUsulanPokokPikiranResponse{
		Code:   http.StatusOK,
		Status: "berhasil mengupdate usulan pokok pikiran",
		Data:   usulanPokokPikiranResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UsulanPokokPikiranControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pegawaiID := params.ByName("pegawai_id")
	rekinID := params.ByName("rencana_kinerja_id")
	isActive := request.URL.Query().Get("is_active")

	var pegawaiIDPtr *string
	if pegawaiID != "" {
		pegawaiIDPtr = &pegawaiID
	}

	var rekinIDPtr *string
	if rekinID != "" {
		rekinIDPtr = &rekinID
	}

	var isActivePtr *bool
	if isActive != "" {
		isActiveBool, err := strconv.ParseBool(isActive)
		if err != nil {
			webResponse := web.WebUsulanPokokPikiranResponse{
				Code:   http.StatusBadRequest,
				Status: "BAD REQUEST",
				Data:   "Parameter is_active harus berupa boolean",
			}
			helper.WriteToResponseBody(writer, webResponse)
			return
		}
		isActivePtr = &isActiveBool
	}
	usulanPokokPikiranResponses, err := controller.UsulanPokokPikiranService.FindAll(request.Context(), pegawaiIDPtr, isActivePtr, rekinIDPtr)
	if err != nil {
		webResponse := web.WebUsulanPokokPikiranResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebUsulanPokokPikiranResponse{
		Code:   http.StatusOK,
		Status: "berhasil mendapatkan semua usulan pokok pikiran",
		Data:   usulanPokokPikiranResponses,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UsulanPokokPikiranControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	idUsulan := params.ByName("id")

	usulanPokokPikiranResponse, err := controller.UsulanPokokPikiranService.FindById(request.Context(), idUsulan)
	if err != nil {
		webResponse := web.WebUsulanPokokPikiranResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebUsulanPokokPikiranResponse{
		Code:   http.StatusOK,
		Status: "berhasil mendapatkan usulan pokok pikiran berdasarkan id",
		Data:   usulanPokokPikiranResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UsulanPokokPikiranControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	idUsulan := params.ByName("id")
	if idUsulan == "" {
		webResponse := web.WebUsulanPokokPikiranResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   "Invalid id usulan parameter",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	err := controller.UsulanPokokPikiranService.Delete(request.Context(), idUsulan)
	if err != nil {
		webResponse := web.WebUsulanPokokPikiranResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebUsulanPokokPikiranResponse{
		Code:   http.StatusOK,
		Status: "berhasil menghapus usulan pokok pikiran",
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UsulanPokokPikiranControllerImpl) FindAllByRekin(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pegawaiID := params.ByName("pegawai_id")
	rekinID := params.ByName("rencana_kinerja_id")
	isActive := request.URL.Query().Get("is_active")

	var pegawaiIDPtr *string
	if pegawaiID != "" {
		pegawaiIDPtr = &pegawaiID
	}

	var rekinIDPtr *string
	if rekinID != "" {
		rekinIDPtr = &rekinID
	}

	var isActivePtr *bool
	if isActive != "" {
		isActiveBool, err := strconv.ParseBool(isActive)
		if err != nil {
			webResponse := web.WebUsulanPokokPikiranResponse{
				Code:   http.StatusBadRequest,
				Status: "BAD REQUEST",
				Data:   "Parameter is_active harus berupa boolean",
			}
			helper.WriteToResponseBody(writer, webResponse)
			return
		}
		isActivePtr = &isActiveBool
	}

	host := os.Getenv("host")
	port := os.Getenv("port")

	buttonActions := []web.ActionButton{
		{
			NameAction: "Create Usulan Pokok Pikiran",
			Method:     "POST",
			Url:        fmt.Sprintf("%s:%s/usulan_pokok_pikiran/create", host, port),
		},
		{
			NameAction: "Update Usulan Pokok Pikiran",
			Method:     "PUT",
			Url:        fmt.Sprintf("%s:%s/usulan_pokok_pikiran/update/:id", host, port),
		},
		{
			NameAction: "Delete Usulan Pokok Pikiran",
			Method:     "DELETE",
			Url:        fmt.Sprintf("%s:%s/usulan_pokok_pikiran/delete/:id", host, port),
		},
		{
			NameAction: "Pilihan Usulan Musrebang",
			Method:     "GET",
			Url:        fmt.Sprintf("%s:%s/usulan_pokok_pikiran/pilihan", host, port),
		},
		{
			NameAction:  "Create Usulan Yang Dipilih",
			Method:      "POST",
			Url:         fmt.Sprintf("%s:%s/usulan_terpilih/create", host, port),
			JenisUsulan: "pokok_pikiran",
		},
	}
	usulanPokokPikiranResponses, err := controller.UsulanPokokPikiranService.FindAll(request.Context(), pegawaiIDPtr, isActivePtr, rekinIDPtr)
	if err != nil {
		webResponse := web.WebUsulanPokokPikiranResponse{
			Code:        http.StatusBadRequest,
			Status:      "BAD REQUEST",
			DataPilihan: err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebUsulanPokokPikiranResponse{
		Code:        http.StatusOK,
		Status:      "berhasil mendapatkan semua usulan pokok pikiran",
		DataPilihan: usulanPokokPikiranResponses,
		Action:      buttonActions,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
