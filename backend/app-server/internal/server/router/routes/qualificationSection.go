package routes

import (
	"app-server/internal/server/handlers"

	"github.com/gorilla/mux"
)

func GetQualificationSectionsRoutes(router *mux.Router) {
	router.HandleFunc(GetQualificationSections, handlers.GetQualificationSection).Methods("GET")
}

// TODO: get qualification section round for user and for admin

func GetQualificationSectionRangesRoutes(router *mux.Router) {
	router.HandleFunc(QualificationSectionRangesEndpoint, handlers.GetQualificationSectionRanges).Methods("GET")
}

func EditQualificationSectionRangesRoutes(router *mux.Router) {
	router.HandleFunc(QualificationSectionRangesEndpoint, handlers.EditQualificationSectionRanges).Methods("PUT")
}

// TODO: end range for admin for user
