package server

import (
	"api/controller"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//StartServer - export
func StartServer() {

	log.Println("*****INFO*****  - server.go *** StartServer function entered")

	r := mux.NewRouter()

	r.HandleFunc("/movies/{type}", controller.GetMovies).Methods("GET")
	r.HandleFunc("/movies/id/{id}/type/{type}", controller.GetMovie).Methods("GET")
	r.HandleFunc("/movies", controller.CreateMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", controller.UpdateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", controller.DeleteMovie).Methods("DELETE")

	log.Println("*****INFO*****  - server.go *** Server started successfully on PORT:8080")
	err := http.ListenAndServe(":8080", r)
	log.Fatal("*****ERROR*****  - server.go *** Issue starting server ", err)
}
