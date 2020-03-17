package controller

import (
	"api/model"
	"api/repo"
	"encoding/json"
	"encoding/xml"
	"fmt"
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
}

//GetMovies controller
func GetMovies(w http.ResponseWriter, r *http.Request) {
	client := repo.NewClient()
	params := mux.Vars(r)
	var tbe string // template to be executed name
	if params["type"] != "json" && params["type"] != "xml" {
		http.Error(w, "type param not correct", http.StatusBadRequest)
		return
	}
	if params["type"] == "json" {
		w.Header().Set("Content-Type", "application/json")
		tbe = "responseTemplateJsonArray.gojson"
	}
	if params["type"] == "xml" {
		w.Header().Set("Content-Type", "text/xml")
		tbe = "responseTemplateXmlArray.goxml"
	}
	movies := repo.GetAllMovies(client)
	err := tpl.ExecuteTemplate(w, tbe, movies)
	if err != nil {
		log.Fatalln(err)
		http.Error(w, "Unable to create template", http.StatusBadRequest)
	}
}

//GetMovie controller
func GetMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	client := repo.NewClient()
	var tbe string // template to be executed name
	if params["type"] != "json" && params["type"] != "xml" {
		http.Error(w, "type param not correct", http.StatusBadRequest)
		return
	}
	if params["type"] == "json" {
		w.Header().Set("Content-Type", "application/json")
		tbe = "responseTemplateJson.gojson"
	}
	if params["type"] == "xml" {
		w.Header().Set("Content-Type", "text/xml")
		tbe = "responseTemplateXml.goxml"
	}

	movie, err := repo.GetByID(client, params["id"])
	if err == redis.Nil {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Fatalln(err)
	}
	err = tpl.ExecuteTemplate(w, tbe, movie)
	if err != nil {
		log.Fatalln(err)
		http.Error(w, "Unable to create template", http.StatusBadRequest)
	}

}

//CreateMovie controller
func CreateMovie(w http.ResponseWriter, r *http.Request) {
	var movie model.Movie
	client := repo.NewClient()
	var tbe string // template to be executed name
	if r.Header["Content-Type"][0] != "text/xml" && r.Header["Content-Type"][0] != "application/json" {
		http.Error(w, "Content Type not supported", http.StatusBadRequest)
		return
	}
	if r.Header["Content-Type"][0] == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		tbe = "responseTemplateJson.gojson"
		err := json.NewDecoder(r.Body).Decode(&movie)
		if err != nil {
			log.Println(err)
			http.Error(w, "Unable to create movie", http.StatusBadRequest)
			return
		}

	}
	if r.Header["Content-Type"][0] == "text/xml" {
		w.Header().Set("Content-Type", "text/xml")
		tbe = "responseTemplateXml.goxml"
		err := xml.NewDecoder(r.Body).Decode(&movie)
		if err != nil {
			log.Println(err)
			http.Error(w, "Unable to create movie", http.StatusBadRequest)
			return
		}
	}

	movie.ID = strconv.Itoa(int(time.Now().Unix()))
	err := repo.Set(client, movie)
	if err != nil {
		fmt.Println(err)
	}
	client.SAdd("movies", movie.ID)

	err = tpl.ExecuteTemplate(w, tbe, movie)
	if err != nil {
		log.Fatalln(err)
		http.Error(w, "Unable to create movie", http.StatusBadRequest)
	}
}

//UpdateMovie controller
func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updatedMovie model.Movie
	client := repo.NewClient()
	var tbe string // template to be executed name
	if r.Header["Content-Type"][0] != "text/xml" && r.Header["Content-Type"][0] != "application/json" {
		http.Error(w, "Content Type not supported", http.StatusBadRequest)
		return
	}
	if r.Header["Content-Type"][0] == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		tbe = "responseTemplateJson.gojson"
		err := json.NewDecoder(r.Body).Decode(&updatedMovie)
		if err != nil {
			log.Println(err)
			http.Error(w, "Unable to update movie", http.StatusBadRequest)
			return
		}
	}
	if r.Header["Content-Type"][0] == "text/xml" {
		w.Header().Set("Content-Type", "text/xml")
		tbe = "responseTemplateXml.goxml"
		err := xml.NewDecoder(r.Body).Decode(&updatedMovie)
		if err != nil {
			log.Println(err)
			http.Error(w, "Unable to update movie", http.StatusBadRequest)
			return
		}
	}
	movie, err := repo.GetByID(client, params["id"])
	if err == redis.Nil {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}
	if err != nil {
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
		fmt.Println(err)
	}
	err = tpl.ExecuteTemplate(w, tbe, movie)
	if err != nil {
		log.Fatalln(err)
		http.Error(w, "Unable to update movie", http.StatusBadRequest)
	}
}

//DeleteMovie controller
func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var tbe string // template to be executed name

	if r.Header["Content-Type"][0] == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		tbe = "responseTemplateJson.gojson"
	}
	if r.Header["Content-Type"][0] == "text/xml" {
		w.Header().Set("Content-Type", "text/xml")
		tbe = "responseTemplateXml.goxml"
	}

	client := repo.NewClient()

	movie, err := repo.Delete(client, params["id"])

	if err == redis.Nil {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Fatalln(err)
	}
	err = tpl.ExecuteTemplate(w, tbe, movie)
	if err != nil {
		log.Fatalln(err)
		http.Error(w, "Unable to create template", http.StatusBadRequest)
	}

}
