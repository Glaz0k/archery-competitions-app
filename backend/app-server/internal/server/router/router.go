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

func CreateIndividualGroupRoutes(router *mux.Router) {
	router.HandleFunc(CreateIndividualGroup, handlers.CreateIndividualGroup).Methods("POST")
}

func CreateRangeGroupRoutes(router *mux.Router) {
	router.HandleFunc(CreateRangeGroup, handlers.CreateRangeGroup).Methods("POST")
}

func CreateQualificationRoutes(router *mux.Router) {
	router.HandleFunc(CreateQualification, handlers.CreateQualification).Methods("POST")
}

func CreateQualificationRoundRoutes(router *mux.Router) {
	router.HandleFunc(CreateQualificationRound, handlers.CreateQualificationRound).Methods("POST")
}

func CreateQualificationSectionRoutes(router *mux.Router) {
	router.HandleFunc(CreateQualificationSection, handlers.CreateQualificationSection).Methods("POST")
}

func CreateRangeRoutes(router *mux.Router) {
	router.HandleFunc(CreateRange, handlers.CreateRange).Methods("POST")
}

func CreateShotRoutes(router *mux.Router) {
	router.HandleFunc(CreateShot, handlers.CreateShot).Methods("POST")
}

func Create() *mux.Router {
	router := mux.NewRouter()
	adminRouter := router.NewRoute().Subrouter()
	adminRouter.Use(delivery.JWTRoleMiddleware("admin"))

	CreateCupRoutes(adminRouter)
	CreateCompetitionRoutes(adminRouter)
	CreateIndividualGroupRoutes(adminRouter)
	CreateRangeGroupRoutes(adminRouter)
	CreateQualificationRoutes(adminRouter)
	CreateQualificationRoundRoutes(adminRouter)
	CreateQualificationSectionRoutes(adminRouter)
	CreateRangeRoutes(adminRouter)

	userRouter := router.NewRoute().Subrouter()
	userRouter.Use(delivery.JWTRoleMiddleware("user"))
	CreateShotRoutes(userRouter)
	CreateShotRoutes(adminRouter)

	return router
}
