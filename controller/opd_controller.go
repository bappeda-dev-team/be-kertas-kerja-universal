package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type OpdController interface {
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	InfoOpd(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
