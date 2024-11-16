package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type PohonKinerjaAdminController interface {
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindSubTematik(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindPokinAdminByIdHierarki(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CreateStrategicAdmin(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindPokinByTematik(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindPokinByStrategic(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindPokinByTactical(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindPokinByOperational(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
