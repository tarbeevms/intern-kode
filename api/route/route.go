package route

import (
	"myapp/api/controllers"

	"github.com/gorilla/mux"
)

func CreateNewRoute(handlerLayer *controllers.HandlerLayer) (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/api/login", handlerLayer.LoginHandler).Methods("POST")
	router.HandleFunc("/api/note", handlerLayer.CreateNote).Methods("POST")
	router.HandleFunc("/api/note/{id}", handlerLayer.UpdateNote).Methods("PUT")
	router.HandleFunc("/api/note/{id}", handlerLayer.DeleteNote).Methods("DELETE")
	router.HandleFunc("/api/note", handlerLayer.GetNotes).Methods("GET")
	return router
}
