package main

import (
	"app-server/internal/delivery"
	"app-server/pkg/logger"
	"net/http"

	"go.uber.org/zap"
)

func main() {
	mux := http.NewServeMux()
	port := ":8000"

	// register handlers normal
	mux.Handle("/competitor", logger.LogMiddleware(delivery.JWTRoleMiddleware("user")(http.HandlerFunc(delivery.CreateShot))))
	logger.Logger.Info("Listening on ", zap.String("port", port))

	mux.HandleFunc("/login", delivery.LoginHandler) //local

	err := http.ListenAndServe(port, mux)
	if err != nil {
		logger.Logger.Info("Problem starting the server", zap.Error(err))
	}

}

//mux.HandleFunc("/task/", server.taskHandler)
//mux.HandleFunc("/tag/", server.tagHandler)
//mux.HandleFunc("/due/", server.dueHandler)
//
//handler := middleware.Logging(mux)
//handler = middleware.PanicRecovery(handler)
//
//log.Fatal(http.ListenAndServe("localhost:"+os.Getenv("SERVERPORT"), handler))
