package controller

import (
	"database/sql"
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type LaporanControllerImpl struct {
	LaporanService service.LaporanService
}

func NewLaporanControllerImpl(s service.LaporanService) *LaporanControllerImpl {
	return &LaporanControllerImpl{
		LaporanService: s,
	}
}

func (controller *LaporanControllerImpl) OpdSupportingPokin(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	kodeOpd := params.ByName("kode_opd")
	tahun := params.ByName("tahun")

	if kodeOpd == "" || tahun == "" {
		webResponse := web.WebResponse{
			Code:    http.StatusBadRequest,
			Status:  http.StatusText(http.StatusBadRequest),
			Message: "Missing requried parameter",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	supportingPokinResponse, err := controller.LaporanService.OpdSupportingPokin(request.Context(), kodeOpd, tahun)
	if err != nil {
		if err == sql.ErrNoRows {
			webResponse := web.WebResponse{
				Code:    http.StatusNoContent,
				Status:  http.StatusText(http.StatusNoContent),
				Message: "OPD tidak memiliki support ke Pohon Kinerja Pemda",
			}
			helper.WriteToResponseBody(writer, webResponse)
			return
		}

		webResponse := web.WebResponse{
			Code:    http.StatusInternalServerError,
			Status:  http.StatusText(http.StatusInternalServerError),
			Message: "Terjadi kesalahan pada server. akan segera ditangani.",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:    200,
		Status:  http.StatusText(http.StatusOK),
		Message: "List Supporting OPD Terhadap Pohon Kinerja Pemda",
		Data:    supportingPokinResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
