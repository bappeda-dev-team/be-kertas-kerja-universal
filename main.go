package main

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/middleware"
	"fmt"
	"log"
	"net/http"
	"os"
)

func NewServer(authMiddleware *middleware.AuthMiddleware) *http.Server {
	host := os.Getenv("host")
	port := os.Getenv("port")
	addr := fmt.Sprintf("%s:%s", host, port)

	cors := helper.NewCORSMiddleware()

	if addr == ":" {
		addr = "localhost:8080"
	}

	return &http.Server{
		Addr:    addr,
		Handler: cors.Handler(authMiddleware),
	}
}

func main() {
	server := InitializeServer()

	log.Printf("Server berjalan di %s", server.Addr)

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
