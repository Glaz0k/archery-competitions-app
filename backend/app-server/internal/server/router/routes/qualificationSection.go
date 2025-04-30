package routes

import (
	"app-server/internal/server/handlers"

	"github.com/gorilla/mux"
)

func GetQualificationSectionsRoutes(router *mux.Router) {
	router.HandleFunc(GetQualificationSections, handlers.GetQualificationSection).Methods("GET")
}

func GetQualificationRoundsRoutes(router *mux.Router) {
	router.HandleFunc(GetQualificationSectionRounds, handlers.GetQualificationRound).Methods("GET")
}

func GetQualificationSectionRangesRoutes(router *mux.Router) {
	router.HandleFunc(QualificationSectionRangesEndpoint, handlers.GetQualificationSectionRanges).Methods("GET")
}

func EditQualificationSectionRangesRoutes(router *mux.Router) {
	router.HandleFunc(QualificationSectionRangesEndpoint, handlers.EditQualificationSectionRanges).Methods("PUT")
}

func EndRangeRoutes(router *mux.Router) {
	router.HandleFunc(EndRange, handlers.EndRange).Methods("POST")
}
