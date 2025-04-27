package routes

import (
	"app-server/internal/server/handlers"

	"github.com/gorilla/mux"
)

func GetSparringPlacesRoutes(router *mux.Router) {
	router.HandleFunc(SparringPlace, handlers.GetSparringPlace).Methods("GET")
}

func GetSparringPlaceRangesRoutes(router *mux.Router) {
	router.HandleFunc(Ranges, handlers.GetRanges).Methods("GET")
}

func EditSparringPlaceRangeRoutes(router *mux.Router) {
	router.HandleFunc(Ranges, handlers.EditSparringPlaceRange).Methods("PUT")
}

func EndSparringPlaceRangeRoutes(router *mux.Router) {
	router.HandleFunc(SparringPlaceRangeEnd, handlers.EndSparringPlaceRange).Methods("POST")
}

func EditShootOutRoutes(router *mux.Router) {
	router.HandleFunc(ShootOut, handlers.EditShootOut).Methods("PUT")
}
