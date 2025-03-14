package api

import (
	"app-server/internal/functions"

	"github.com/gorilla/mux"
)

func Create(router *mux.Router) {
	router.HandleFunc("/cup/register", functions.AddCup).Methods("POST")
}
