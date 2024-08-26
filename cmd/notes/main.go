package main

import (
	"fmt"
	"log"
	"myapp/api/controllers"
	"myapp/api/middleware"
	"myapp/api/route"
	"myapp/internal/logic"
	"myapp/internal/notes"
	"myapp/internal/session"
	"myapp/internal/users"
	"myapp/pkg/dbconnect"
	"net/http"

	"go.uber.org/zap"
)

func main() {
	mysqlDB, err := dbconnect.ConnectToMySQL()
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
		return
	}
	defer mysqlDB.Close()

	noteRepo := notes.NewNoteRepository(mysqlDB)
	userRepo := users.NewUserRepo(mysqlDB)
	sessionRepo := session.NewSessionRepo(mysqlDB)

	logicLayer := logic.NewLogicLayer(noteRepo, userRepo, sessionRepo)

	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	handlerLayer := controllers.NewHandlerLayer(logicLayer, logger)

	router := route.CreateNewRoute(handlerLayer)

	mux := middleware.Logging(logger, router)
	mux = middleware.Authenticated(mux, handlerLayer.LogicLayer)
	mux = middleware.Panic(mux)

	fmt.Println("Server is listening on port 8080...")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
		return
	}
}
