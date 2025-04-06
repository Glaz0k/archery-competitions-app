package routes

import (
	"app-server/internal/server/handlers"

	"github.com/gorilla/mux"
)

func GetIndividualGroupsRoutes(router *mux.Router) {
	router.HandleFunc(IndividualGroupEndpoint, handlers.GetIndividualGroups).Methods("GET")
}

func DeleteIndividualGroupRoutes(router *mux.Router) {
	router.HandleFunc(IndividualGroupEndpoint, handlers.DeleteIndividualGroup).Methods("DELETE")
}

// TODO: Разбить на админа и юзера

func GetIndividualGroupCompetitorsRoutes(router *mux.Router) {
	router.HandleFunc(IndividualGroupCompetitorsEndpoint, handlers.GetCompetitor).Methods("GET")
}

func SyncIndividualGroupsRoutes(router *mux.Router) {
	router.HandleFunc(UpdateIndividualGroup, handlers.SyncIndividualGroup).Methods("POST")
}

func GetCompetitorFromIndividualGroupRoutes(router *mux.Router) {
	router.HandleFunc(IndividualGroupCompetitorsEndpoint, handlers.GetCompetitorsFromGroup).Methods("GET")
}

// TODO: разбить на юзера и админа

func GetQualificationTableRoutes(router *mux.Router) {
	router.HandleFunc(QualificationTable, handlers.GetQualifications).Methods("GET")
}

func StartQualificationRoutes(router *mux.Router) {
	router.HandleFunc(StartQualification, handlers.StartQualification).Methods("POST")
}

// TODO: end qualification

// TODO: final grid

// TODO: quarterfinal start

// TODO: semifinal start

// TODO: final start

// TODO: final end
