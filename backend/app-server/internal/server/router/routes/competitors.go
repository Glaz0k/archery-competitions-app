package routes

import (
	"app-server/internal/server/handlers"

	"github.com/gorilla/mux"
)

func RegisterCompetitorRoutes(router *mux.Router) {
	router.HandleFunc(RegisterCompetitor, handlers.RegisterCompetitor).Methods("POST")
}

func GetCompetitorRoutes(router *mux.Router) {
	router.HandleFunc(GetCompetitor, handlers.GetCompetitor).Methods("GET")
}

// TODO: edit for admin
func EditCompetitorUserRoutes(router *mux.Router) {
	router.HandleFunc(Competitor, handlers.UserEditCompetitor).Methods("PUT")
}

// TODO: delete competitor
