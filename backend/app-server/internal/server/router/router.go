package router

import (
	"app-server/internal/delivery"
	"app-server/internal/server/handlers"
	"app-server/pkg/logger"

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

func StartQualificationRoutes(router *mux.Router) {
	router.HandleFunc(StartQualification, handlers.StartQualification).Methods("POST")
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

func EditCompetitionRoutes(router *mux.Router) {
	router.HandleFunc(CompetitionEndpoint, handlers.EditCompetition).Methods("PUT")
}

func GetCupRoutes(router *mux.Router) {
	router.HandleFunc(CupEndpoint, handlers.GetCup).Methods("GET")
}

func GetIndividualGroupsRoutes(router *mux.Router) {
	router.HandleFunc(IndividualGroupEndpoint, handlers.GetIndividualGroups).Methods("GET")
}

func GetIndividualGroupCompetitorsRoutes(router *mux.Router) {
	router.HandleFunc(IndividualGroupCompetitorsEndpoint, handlers.GetCompetitors).Methods("GET")
}

func EditCupRoutes(router *mux.Router) {
	router.HandleFunc(CupEndpoint, handlers.EditCup).Methods("PUT")
}
func Create() *mux.Router {
	router := mux.NewRouter()
	router.Use(logger.LogMiddleware)

	adminRouter := router.NewRoute().Subrouter()
	adminRouter.Use(delivery.JWTRoleMiddleware("admin"))

	CreateCupRoutes(adminRouter)
	CreateCompetitionRoutes(adminRouter)
	CreateIndividualGroupRoutes(adminRouter)
	CreateRangeGroupRoutes(adminRouter)
	StartQualificationRoutes(adminRouter)
	CreateQualificationRoundRoutes(adminRouter)
	CreateQualificationSectionRoutes(adminRouter)
	CreateRangeRoutes(adminRouter)

	EditCompetitionRoutes(adminRouter)
	EditCupRoutes(adminRouter)

	userRouter := router.NewRoute().Subrouter()
	userRouter.Use(delivery.JWTRoleMiddleware("user"))
	CreateShotRoutes(userRouter)
	CreateShotRoutes(adminRouter)

	GetCupRoutes(userRouter)
	GetCupRoutes(adminRouter)

	GetIndividualGroupsRoutes(adminRouter)
	GetIndividualGroupsRoutes(userRouter)
	GetIndividualGroupCompetitorsRoutes(userRouter)
	GetIndividualGroupCompetitorsRoutes(adminRouter)

	return router
}
