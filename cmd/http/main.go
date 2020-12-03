package main

import (
	"net/http"
	"os"

	"hellper/internal/config"
	"hellper/internal/handler"
)

func main() {
	http.HandleFunc("/", handler.NewHandlerRoute())
	http.ListenAndServe(determineListenAddress(), nil)
}

func determineListenAddress() string {
	port := os.Getenv("PORT")

	if port != "" {
		return port
	}

	return config.Env.BindAddress
}
