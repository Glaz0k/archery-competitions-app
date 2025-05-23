package routes

import (
	"app-server/internal/server/handlers"

	"github.com/gorilla/mux"
)

func GetIndividualGroupRoutes(router *mux.Router) {
	router.HandleFunc(IndividualGroupEndpoint, handlers.GetIndividualGroup).Methods("GET")
}

func DeleteIndividualGroupRoutes(router *mux.Router) {
	router.HandleFunc(IndividualGroupEndpoint, handlers.DeleteIndividualGroup).Methods("DELETE")
}

// TODO: Разбить на админа и юзера

func GetIndividualGroupCompetitorsRoutes(router *mux.Router) {
	router.HandleFunc(IndividualGroupCompetitorsEndpoint, handlers.GetCompetitorsFromGroup).Methods("GET")
}

func SyncIndividualGroupsRoutes(router *mux.Router) {
	router.HandleFunc(UpdateIndividualGroup, handlers.SyncIndividualGroup).Methods("POST")
}

func GetCompetitorFromIndividualGroupRoutes(router *mux.Router) {
	router.HandleFunc(IndividualGroupCompetitorsEndpoint, handlers.GetCompetitorsFromGroup).Methods("GET")
}

// TODO: разбить на юзера и админа

func GetQualificationTableRoutes(router *mux.Router) {
	router.HandleFunc(QualificationTable, handlers.GetQualification).Methods("GET")
}

func StartQualificationRoutes(router *mux.Router) {
	router.HandleFunc(StartQualification, handlers.StartQualification).Methods("POST")
}

func EndQualificationRoutes(router *mux.Router) {
	router.HandleFunc(EndQualification, handlers.EndQualification).Methods("POST")
}

func GetFinalGridRoutes(router *mux.Router) {
	router.HandleFunc(FinalGrid, handlers.GetFinalGrid).Methods("GET")
}

func StartQuarterfinalRoutes(router *mux.Router) {
	router.HandleFunc(StartQuarterfinal, handlers.StartQuarterfinal).Methods("POST")
}

func StartSemifinalRoutes(router *mux.Router) {
	router.HandleFunc(StartSemifinal, handlers.StartSemifinal).Methods("POST")
}

func StartFinalRoutes(router *mux.Router) {
	router.HandleFunc(StartFinal, handlers.StartFinal).Methods("POST")
}

func EndFinalRoutes(router *mux.Router) {
	router.HandleFunc(EndFinal, handlers.EndFinal).Methods("POST")
}
