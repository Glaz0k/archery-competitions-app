package routes

import (
	"app-server/internal/server/handlers"

	"github.com/gorilla/mux"
)

func GetQualificationSectionsRoutes(router *mux.Router) {
	router.HandleFunc(GetQualificationSections, handlers.GetQualificationSection).Methods("GET")
}

// TODO: get qualification section round for user and for admin

// TODO: get qualification section ranges for admin for user

// TODO: edit range for admin and for user

// TODO: end range for admin for user
