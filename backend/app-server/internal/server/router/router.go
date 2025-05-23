package router

import (
	"app-server/pkg/logger"

	"github.com/gorilla/mux"

	"app-server/internal/delivery"
	. "app-server/internal/server/router/routes"
)

func Create(secretKey string) *mux.Router {
	router := mux.NewRouter()
	router.Use(logger.LogMiddleware)

	adminRouter := router.NewRoute().Subrouter()
	adminRouter.Use(delivery.JWTRoleMiddleware("admin", secretKey))

	userRouter := router.NewRoute().Subrouter()
	userRouter.Use(delivery.JWTRoleMiddleware("user", secretKey))

	commonRouter := router.NewRoute().Subrouter()
	commonRouter.Use(delivery.JWTRoleMiddleware("admin, user", secretKey))

	CreateCupRoutes(adminRouter)
	GetCupRoutes(commonRouter)
	GetAllCupsRoutes(commonRouter)
	EditCupRoutes(adminRouter)
	CreateCompetitionRoutes(adminRouter)
	GetAllCompetitionsRoutes(commonRouter)
	GetCompetitionsRoutes(commonRouter)

	CreateIndividualGroupRoutes(adminRouter)
	StartQualificationRoutes(adminRouter)
	EndQualificationRoutes(adminRouter)

	EditCompetitionRoutes(adminRouter)

	GetIndividualGroupRoutes(commonRouter)
	GetCompetitorFromIndividualGroupRoutes(commonRouter)

	SyncIndividualGroupsRoutes(adminRouter)

	GetCompetitorsFromCompetitionRoutes(commonRouter)
	GetIndividualGroupCompetitorsRoutes(commonRouter)
	GetFinalGridRoutes(commonRouter)
	StartQuarterfinalRoutes(adminRouter)
	StartSemifinalRoutes(adminRouter)
	StartFinalRoutes(adminRouter)
	EndFinalRoutes(adminRouter)

	GetQualificationTableRoutes(userRouter)
	GetQualificationTableRoutes(adminRouter)

	EndCompetitionRoutes(adminRouter)

	DeleteIndividualGroupRoutes(adminRouter)

	EditCompetitorRoutes(commonRouter)

	GetCompetitorsFromCompetitionRoutes(commonRouter)

	GetQualificationSectionsRoutes(commonRouter)
	GetQualificationRoundsRoutes(commonRouter)

	RegisterCompetitorRoutes(userRouter)

	GetCompetitorRoutes(commonRouter)

	AddCompetitorCompetitionRoutes(adminRouter)
	EditCompetitorStatusRoutes(commonRouter)
	DeleteCompetitorCompetitorRoutes(adminRouter)
	GetIndividualGroupsFromCompetitionRoutes(commonRouter)
	GetQualificationSectionRangesRoutes(commonRouter)
	EditQualificationSectionRangesRoutes(commonRouter)
	EndRangeRoutes(commonRouter)
	DeleteCompetitionRoutes(adminRouter)
	GetAllCompetitorsRoutes(adminRouter)

	DeleteCompetitorRoutes(adminRouter)

	GetSparringPlacesRoutes(commonRouter)
	GetSparringPlaceRangesRoutes(commonRouter)
	EditSparringPlaceRangeRoutes(commonRouter)
	EndSparringPlaceRangeRoutes(commonRouter)
	EditShootOutRoutes(commonRouter)
	DeleteCupRoutes(adminRouter)

	return router
}
