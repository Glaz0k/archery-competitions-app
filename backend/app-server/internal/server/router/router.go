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

	CreateCupRoutes(adminRouter)
	CreateCompetitionRoutes(adminRouter)
	CreateIndividualGroupRoutes(adminRouter)
	StartQualificationRoutes(adminRouter)

	EditCompetitionRoutes(adminRouter)
	EditCupRoutes(adminRouter)

	GetCupRoutes(userRouter)
	GetCupRoutes(adminRouter)

	GetAllCupsRoutes(userRouter)
	GetAllCupsRoutes(adminRouter)

	GetAllCompetitionsRoutes(userRouter)
	GetAllCompetitionsRoutes(adminRouter)

	GetIndividualGroupsRoutes(userRouter)
	GetIndividualGroupsRoutes(adminRouter)

	SyncIndividualGroupsRoutes(adminRouter)

	GetIndividualGroupCompetitorsRoutes(userRouter)
	//GetIndividualGroupCompetitorsRoutes(adminRouter)

	GetQualificationTableRoutes(userRouter)
	GetQualificationTableRoutes(adminRouter)
	
	EndCompetitionRoutes(adminRouter)

	DeleteIndividualGroupRoutes(adminRouter)

	EditCompetitorUserRoutes(userRouter)
	// admin router
	GetCompetitorsFromCompetitionUserRoutes(userRouter)

	GetQualificationSectionsRoutes(userRouter)
	GetQualificationSectionsRoutes(adminRouter)

	RegisterCompetitorRoutes(userRouter)

	GetCompetitorRoutes(adminRouter)
	GetCompetitorRoutes(userRouter)

	AddCompetitorCompetitionRoutes(adminRouter)
	EditCompetitorStatusAdminRoutes(adminRouter)

	return router
}
