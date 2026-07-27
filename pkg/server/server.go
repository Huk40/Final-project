package server

import (
	"net/http"
	"os"

	"goFinalProject/pkg/api"
)

const defaultPort = "7540"

func Run() error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	api.Init()
	http.Handle("/", http.FileServer(http.Dir("web")))
	return http.ListenAndServe(":"+port, nil)
}
