package controller

import (
	"api/model"
	"api/repo"
	"encoding/json"
	"encoding/xml"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	log.Println("*****INFO*****  - controller.go *** templates initiated")
}

//GetMovies controller
func GetMovies(w http.ResponseWriter, r *http.Request) {
	log.Println("*****INFO*****  - controller.go *** GetMovies function entered")
	client := repo.NewClient()
	params := mux.Vars(r)
	var tbe string // template to be executed name
	if params["type"] != "json" && params["type"] != "xml" {
		http.Error(w, "type param not correct", http.StatusBadRequest)
		log.Println("*****INFO*****  - controller.go *** type param not correct")
		return
	}
	if params["type"] == "json" {
		log.Println("*****INFO*****  - controller.go *** JSON will be returned")
		w.Header().Set("Content-Type", "application/json")
		tbe = "responseTemplateJsonArray.gojson"
	}
	if params["type"] == "xml" {
		log.Println("*****INFO*****  - controller.go *** XML will be returned")
		w.Header().Set("Content-Type", "text/xml")
		tbe = "responseTemplateXmlArray.goxml"
	}

	movies := repo.GetAllMovies(client)
	err := tpl.ExecuteTemplate(w, tbe, movies)
	if err != nil {
		log.Println("*****ERROR*****  - controller.go *** error occured during template execution")
		http.Error(w, "Unable to create template", http.StatusInternalServerError)
		log.Fatalln(err)
	}
	log.Println("*****INFO*****  - controller.go *** movies returned successfully")
}

//GetMovie controller
func GetMovie(w http.ResponseWriter, r *http.Request) {
	log.Println("*****INFO*****  - controller.go *** GetMovie function entered")
	params := mux.Vars(r)
	client := repo.NewClient()
	var tbe string // template to be executed name
	if params["type"] != "json" && params["type"] != "xml" {
		http.Error(w, "type param not correct", http.StatusBadRequest)
		log.Println("*****INFO*****  - controller.go *** type param not correct")
		return
	}
	if params["type"] == "json" {
		log.Println("*****INFO*****  - controller.go *** JSON will be returned")
		w.Header().Set("Content-Type", "application/json")
		tbe = "responseTemplateJson.gojson"
	}
	if params["type"] == "xml" {
		log.Println("*****INFO*****  - controller.go *** XML will be returned")
		w.Header().Set("Content-Type", "text/xml")
		tbe = "responseTemplateXml.goxml"
	}

	movie, err := repo.GetByID(client, params["id"])
	if err == redis.Nil {
		log.Println("*****INFO*****  - controller.go *** movie with the" + params["id"] + "ID not found")
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Println("*****ERROR*****  - controller.go *** error occured")
		log.Fatalln(err)

	}
	err = tpl.ExecuteTemplate(w, tbe, movie)
	if err != nil {
		log.Println("*****ERROR*****  - controller.go *** error occured during template execution")
		http.Error(w, "Unable to create template", http.StatusInternalServerError)
		log.Fatalln(err)
	}
	log.Println("*****INFO*****  - controller.go *** movie returned successfully")

}

//CreateMovie controller
func CreateMovie(w http.ResponseWriter, r *http.Request) {
	log.Println("*****INFO*****  - controller.go *** CreateMovie function entered")
	var movie model.Movie
	client := repo.NewClient()
	var tbe string // template to be executed name
	if r.Header["Content-Type"][0] != "text/xml" && r.Header["Content-Type"][0] != "application/json" {
		http.Error(w, "Content Type not supported", http.StatusBadRequest)
		log.Println("*****INFO*****  - controller.go *** content type not supported")
		return
	}
	if r.Header["Content-Type"][0] == "application/json" {
		log.Println("*****INFO*****  - controller.go *** content type is JSON and JSON will be returned")
		w.Header().Set("Content-Type", "application/json")
		tbe = "responseTemplateJson.gojson"
		err := json.NewDecoder(r.Body).Decode(&movie)
		if err != nil {
			log.Println("*****ERROR*****  - controller.go *** error occured during JSON decoding ")
			http.Error(w, "Unable to create movie", http.StatusBadRequest)
			return
		}
		log.Println("*****INFO*****  - controller.go *** JSON decoded successfully")

	}
	if r.Header["Content-Type"][0] == "text/xml" {
		log.Println("*****INFO*****  - controller.go *** content type is XML and XML will be returned")
		w.Header().Set("Content-Type", "text/xml")
		tbe = "responseTemplateXml.goxml"
		err := xml.NewDecoder(r.Body).Decode(&movie)
		if err != nil {
			log.Println("*****ERROR*****  - controller.go *** error occured during XML decoding ")
			http.Error(w, "Unable to create movie", http.StatusBadRequest)
			return
		}
		log.Println("*****INFO*****  - controller.go *** XML decoded successfully")

	}

	movie.ID = strconv.Itoa(int(time.Now().Unix())) // random ID generator, not safe
	err := repo.Set(client, movie)
	if err != nil {
		log.Println("*****ERROR*****  - controller.go *** error occured during save")
	}
	client.SAdd("movies", movie.ID)
	log.Println("*****INFO*****  - controller.go *** ID of the newly created movie added to set of movie IDs")
	err = tpl.ExecuteTemplate(w, tbe, movie)
	if err != nil {
		log.Println("*****ERROR*****  - controller.go *** error occured during template execution")
		http.Error(w, "Error occured: could not return the created movie, but it was created in DB", http.StatusInternalServerError)

	}
	log.Println("*****INFO*****  - controller.go *** movie created successfully")
}

//UpdateMovie controller
func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	log.Println("*****INFO*****  - controller.go *** UpdateMovie function entered")
	params := mux.Vars(r)
	var updatedMovie model.Movie
	client := repo.NewClient()
	var tbe string // template to be executed name
	if r.Header["Content-Type"][0] != "text/xml" && r.Header["Content-Type"][0] != "application/json" {
		http.Error(w, "Content Type not supported", http.StatusBadRequest)
		log.Println("*****INFO*****  - controller.go *** content type not supported")
		return
	}
	if r.Header["Content-Type"][0] == "application/json" {
		log.Println("*****INFO*****  - controller.go *** content type is JSON and JSON will be returned")
		w.Header().Set("Content-Type", "application/json")
		tbe = "responseTemplateJson.gojson"
		err := json.NewDecoder(r.Body).Decode(&updatedMovie)
		if err != nil {
			log.Println("*****ERROR*****  - controller.go *** error occured during JSON decoding ")
			http.Error(w, "Unable to update movie", http.StatusBadRequest)
			return
		}
		log.Println("*****INFO*****  - controller.go *** JSON decoded successfully")
	}
	if r.Header["Content-Type"][0] == "text/xml" {
		log.Println("*****INFO*****  - controller.go *** content type is XML and XML will be returned")
		w.Header().Set("Content-Type", "text/xml")
		tbe = "responseTemplateXml.goxml"
		err := xml.NewDecoder(r.Body).Decode(&updatedMovie)
		if err != nil {
			log.Println("*****ERROR*****  - controller.go *** error occured during XML decoding ")
			http.Error(w, "Unable to update movie", http.StatusBadRequest)
			return
		}
		log.Println("*****INFO*****  - controller.go *** XML decoded successfully")

	}
	movie, err := repo.GetByID(client, params["id"])
	log.Println("*****INFO*****  - controller.go *** movie return from DB")
	if err == redis.Nil {
		log.Println("*****INFO*****  - controller.go *** movie with the" + params["id"] + "ID not found")
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Println("*****ERROR*****  - controller.go *** error occured")
		log.Fatalln(err)
	}
	if !(updatedMovie.Genre == " " || updatedMovie.Genre == "") {
		movie.Genre = updatedMovie.Genre
	}
	if !(updatedMovie.Title == " " || updatedMovie.Title == "") {
		movie.Title = updatedMovie.Title
	}
	if updatedMovie.Director != nil {
		if !(updatedMovie.Director.Lastname == " " || updatedMovie.Director.Lastname == "") {
			movie.Director.Lastname = updatedMovie.Director.Lastname
		}
		if !(updatedMovie.Director.Firstname == " " || updatedMovie.Director.Firstname == "") {
			movie.Director.Firstname = updatedMovie.Director.Firstname
		}
	}
	err = repo.Set(client, movie)
	if err != nil {
		log.Println("*****ERROR*****  - controller.go *** error occured during save")
	}
	err = tpl.ExecuteTemplate(w, tbe, movie)
	if err != nil {
		log.Println("*****ERROR*****  - controller.go *** error occured during template execution")
		http.Error(w, "Error occured: could not return the updated movie, but it was updated in DB", http.StatusInternalServerError)
	}
	log.Println("*****INFO*****  - controller.go *** movie updated successfully")
}

//DeleteMovie controller
func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	log.Println("*****INFO*****  - controller.go *** DeleteMovie function entered")
	params := mux.Vars(r)
	client := repo.NewClient()
	var tbe string // template to be executed name

	if r.Header["Content-Type"][0] == "application/json" {
		log.Println("*****INFO*****  - controller.go *** content type JSON will be returned")
		w.Header().Set("Content-Type", "application/json")
		tbe = "responseTemplateJson.gojson"
	}
	if r.Header["Content-Type"][0] == "text/xml" {
		log.Println("*****INFO*****  - controller.go *** content type XML will be returned")
		w.Header().Set("Content-Type", "text/xml")
		tbe = "responseTemplateXml.goxml"
	}
	movie, err := repo.Delete(client, params["id"])

	if err == redis.Nil {
		log.Println("*****INFO*****  - controller.go *** movie with the" + params["id"] + "ID not found")
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Println("*****ERROR*****  - controller.go *** error occured")
		log.Fatalln(err)
	}
	err = tpl.ExecuteTemplate(w, tbe, movie)
	if err != nil {
		log.Println("*****ERROR*****  - controller.go *** error occured during template execution")
		http.Error(w, "Error occured: could not return the deleted movie, but it was deleted in DB", http.StatusInternalServerError)
	}
	log.Println("*****INFO*****  - controller.go *** movie deleted successfully")
}
