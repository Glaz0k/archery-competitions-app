package router

import (
	"app-server/pkg/logger"

	"github.com/gorilla/mux"

	"app-server/internal/delivery"
	. "app-server/internal/server/router/routes"
)

func Create() *mux.Router {
	router := mux.NewRouter()
	router.Use(logger.LogMiddleware)

	adminRouter := router.NewRoute().Subrouter()
	adminRouter.Use(delivery.JWTRoleMiddleware("admin"))

	userRouter := router.NewRoute().Subrouter()
	userRouter.Use(delivery.JWTRoleMiddleware("user"))

	commonRouter := router.NewRoute().Subrouter()
	commonRouter.Use(delivery.JWTRoleMiddleware("admin, user"))

	CreateCupRoutes(adminRouter)
	GetCupRoutes(commonRouter)
	GetAllCupsRoutes(commonRouter)
	EditCupRoutes(adminRouter)
	CreateCompetitionRoutes(adminRouter)
	GetAllCompetitionsRoutes(commonRouter)

	CreateIndividualGroupRoutes(adminRouter)
	StartQualificationRoutes(adminRouter)
	EndQualificationRoutes(adminRouter)

	EditCompetitionRoutes(adminRouter)

	GetIndividualGroupsRoutes(commonRouter)

	SyncIndividualGroupsRoutes(adminRouter)

	GetCompetitorsFromCompetitionRoutes(commonRouter)
	GetIndividualGroupCompetitorsRoutes(commonRouter)
	GetFinalGrid(commonRouter)

	GetQualificationTableRoutes(userRouter)
	GetQualificationTableRoutes(adminRouter)

	EndCompetitionRoutes(adminRouter)

	DeleteIndividualGroupRoutes(adminRouter)

	EditCompetitorRoutes(commonRouter)

	GetCompetitorsFromCompetitionRoutes(commonRouter)

	GetQualificationSectionsRoutes(userRouter)
	GetQualificationSectionsRoutes(adminRouter)

	RegisterCompetitorRoutes(userRouter)

	GetCompetitorRoutes(commonRouter)

	AddCompetitorCompetitionRoutes(adminRouter)
	EditCompetitorStatusRoutes(commonRouter)
	DeleteCompetitorCompetitorRoutes(adminRouter)
	GetIndividualGroupsFromCompetitionRoutes(commonRouter)

	return router
}
