package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type SubKegiatanTerpilihController interface {
	Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByKodeSubKegiatan(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CreateRekin(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteSubKegiatanTerpilih(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
