package middleware

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/model/web"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	Handler http.Handler
}

func NewAuthMiddleware(handler http.Handler) *AuthMiddleware {
	return &AuthMiddleware{Handler: handler}
}

func (middleware *AuthMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/user/login" {
		middleware.Handler.ServeHTTP(writer, request)
		return
	}

	tokenString := request.Header.Get("Authorization")
	if tokenString == "" {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnauthorized)

		webResponse := web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
			Data:   "Token tidak ditemukan",
		}

		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	claims := helper.ValidateJWT(tokenString)
	if claims.UserId == 0 {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnauthorized)

		webResponse := web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
			Data:   "Token tidak valid",
		}

		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	middleware.Handler.ServeHTTP(writer, request)
}
