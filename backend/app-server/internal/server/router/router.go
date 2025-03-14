package router

import (
	"app-server/internal/delivery"
	"app-server/internal/server/handlers"

	"github.com/gorilla/mux"
)

func CreateCupRoutes(router *mux.Router) {
	router.HandleFunc(CreateCup, handlers.CreateCup).Methods("POST")
}

func CreateCompetitionRoutes(router *mux.Router) {
	router.HandleFunc(CreateCompetition, handlers.CreateCompetition).Methods("POST")
}

func Create() *mux.Router {
	router := mux.NewRouter()
	adminRouter := router.NewRoute().Subrouter()
	adminRouter.Use(delivery.JWTRoleMiddleware("admin"))
	CreateCupRoutes(adminRouter)
	CreateCompetitionRoutes(adminRouter)

	return router
}
