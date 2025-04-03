package routes

import (
	"app-server/internal/server/handlers"

	"github.com/gorilla/mux"
)

func CreateCupRoutes(router *mux.Router) {
	router.HandleFunc(CupsEndpoint, handlers.CreateCup).Methods("POST")
}

func GetCupRoutes(router *mux.Router) {
	router.HandleFunc(CupEndpoint, handlers.GetCup).Methods("GET")
}

func GetAllCupsRoutes(router *mux.Router) {
	router.HandleFunc(CupsEndpoint, handlers.GetAllCups).Methods("GET")
}

func EditCupRoutes(router *mux.Router) {
	router.HandleFunc(CupEndpoint, handlers.EditCup).Methods("PUT")
}

// TODO: DeleteCupRoutes

func CreateCompetitionRoutes(router *mux.Router) {
	router.HandleFunc(CompetitionsEndpoint, handlers.CreateCompetition).Methods("POST")
}

func GetAllCompetitionsRoutes(router *mux.Router) {
	router.HandleFunc(CompetitionsEndpoint, handlers.GetAllCompetitions).Methods("GET")
}
