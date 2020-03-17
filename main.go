package main

import (
	"api/model"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

/* **************** CRUD **************** */

// Get all movies
func getMovies(w http.ResponseWriter, r *http.Request) {
	client := newClient()
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
	movies := getAllMovies(client)
	err := tpl.ExecuteTemplate(w, tbe, movies)
	if err != nil {
		log.Fatalln(err)
		http.Error(w, "Unable to create template", http.StatusBadRequest)
	}
}

// Get movie by ID
func getMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	client := newClient()
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

	movie, err := getByID(client, params["id"])
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

//Post now movie
func createMovie(w http.ResponseWriter, r *http.Request) {
	var movie model.Movie
	client := newClient()
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
	err := set(client, movie)
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

//Update exisiting entry
func updateMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updatedMovie model.Movie
	client := newClient()
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
	movie, err := getByID(client, params["id"])
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
	err = set(client, movie)
	if err != nil {
		fmt.Println(err)
	}
	err = tpl.ExecuteTemplate(w, tbe, movie)
	if err != nil {
		log.Fatalln(err)
		http.Error(w, "Unable to update movie", http.StatusBadRequest)
	}
}

//Delete movie
func deleteMovie(w http.ResponseWriter, r *http.Request) {
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

	client := newClient()

	movie, err := delete(client, params["id"])

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

/* **************** MAIN / ROUTER **************** */

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/movies/{type}", getMovies).Methods("GET")
	r.HandleFunc("/movies/id/{id}/type/{type}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	log.Println("Starting Server on port 8080")
	err := http.ListenAndServe(":8080", r)
	log.Fatal(err)
}

/* **************** Connection with REDIS **************** */

func newClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return client
}

func set(client *redis.Client, movie model.Movie) error {
	json, err := json.Marshal(movie)
	if err != nil {
		return err
	}
	err = client.Set("movie:"+movie.ID, json, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func getByID(client *redis.Client, id string) (model.Movie, error) {
	val, err := client.Get("movie:" + id).Result()
	movie := model.Movie{}
	if err != nil {
		log.Println(err)
		return movie, err
	}
	err = json.Unmarshal([]byte(val), &movie)
	if err != nil {
		log.Println(err)
	}
	return movie, err
}

func getAllMovies(client *redis.Client) []model.Movie {
	movieIDs, err := client.SMembers("movies").Result()
	if err != nil {
		log.Println(err)
	}
	var movies []model.Movie
	for _, movieID := range movieIDs {
		movie, _ := getByID(client, movieID)
		movies = append(movies, movie)
	}
	return movies

}

func delete(client *redis.Client, id string) (model.Movie, error) {
	movie, err := getByID(client, id)
	if err == nil {
		_, err = client.SRem("movies", id).Result()
		if err != nil {
			log.Println(err)
		}
		_, err = client.Del("movie:" + id).Result()
		if err != nil {
			log.Println(err)
		}
	}
	return movie, err
}
