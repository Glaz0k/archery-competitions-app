package routes

import (
	"app-server/internal/server/handlers"

	"github.com/gorilla/mux"
)

func EditCompetitionRoutes(router *mux.Router) {
	router.HandleFunc(CompetitionEndpoint, handlers.EditCompetition).Methods("PUT")
}

func DeleteCompetitionRoutes(router *mux.Router) {
	router.HandleFunc(CompetitionEndpoint, handlers.DeleteCompetition).Methods("DELETE")
}

func EndCompetitionRoutes(router *mux.Router) {
	router.HandleFunc(EndCompetition, handlers.EndCompetition).Methods("POST")
}

func AddCompetitorCompetitionRoutes(router *mux.Router) {
	router.HandleFunc(CompetitorsCompetitionEndpoint, handlers.AddCompetitorCompetition).Methods("POST")
}

func GetCompetitorsFromCompetitionRoutes(router *mux.Router) {
	router.HandleFunc(CompetitorsCompetitionEndpoint, handlers.GetCompetitorsFromCompetition).Methods("GET")
}

func EditCompetitorStatusRoutes(router *mux.Router) {
	router.HandleFunc(CompetitorCompetitionEndpoint, handlers.EditCompetitorStatus).Methods("PUT")
}

func DeleteCompetitorCompetitorRoutes(router *mux.Router) {
	router.HandleFunc(CompetitorCompetitionEndpoint, handlers.DeleteCompetitorCompetition).Methods("DELETE")
}

func CreateIndividualGroupRoutes(router *mux.Router) {
	router.HandleFunc(IndividualGroupsCompetition, handlers.CreateIndividualGroup).Methods("POST")
}

func GetIndividualGroupsFromCompetitionRoutes(router *mux.Router) {
	router.HandleFunc(IndividualGroupsCompetition, handlers.GetIndividualGroupsFromCompetition).Methods("GET")
}

func GetCompetitionsRoutes(router *mux.Router) {
	router.HandleFunc(CompetitionEndpoint, handlers.GetCompetitions).Methods("GET")
}
