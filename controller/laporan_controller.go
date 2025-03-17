package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type LaporanController interface {
	OpdSupportingPokin(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
