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

func EditCompetitorRoutes(router *mux.Router) {
	router.HandleFunc(Competitor, handlers.EditCompetitor).Methods("PUT")
}

func GetAllCompetitorsRoutes(router *mux.Router) {
	router.HandleFunc(AllCompetitorsEndpoint, handlers.GetAllCompetitors).Methods("GET")
}

func DeleteCompetitorRoutes(router *mux.Router) {
	router.HandleFunc(Competitor, handlers.DeleteCompetitor).Methods("DELETE")
}
