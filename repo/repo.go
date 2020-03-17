package repo

import (
	"api/model"
	"encoding/json"
	"log"

	"github.com/go-redis/redis"
)

//NewClient - connect with DB
func NewClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	log.Println("*****INFO*****  - repo.go *** redis client created")
	return client
}

//Set a new movie
func Set(client *redis.Client, movie model.Movie) error {
	log.Println("*****INFO*****  - repo.go *** set method initiated")
	json, err := json.Marshal(movie)
	if err != nil {
		log.Println("*****ERROR*****  - repo.go *** error during json.Marshal")
		return err
	}
	err = client.Set("movie:"+movie.ID, json, 0).Err()
	if err != nil {
		log.Println("*****ERROR*****  - repo.go *** error during saving movie into base")
		return err
	}
	log.Println("*****INFO*****  - repo.go *** movie save completed successfully ")
	return nil
}

//GetByID - return a movie by ID
func GetByID(client *redis.Client, id string) (model.Movie, error) {
	log.Println("*****INFO*****  - repo.go *** GetByID method initiated")
	val, err := client.Get("movie:" + id).Result()
	movie := model.Movie{}
	if err != nil {
		log.Println("*****ERROR*****  - repo.go ***  error retrieving movie from base")
		return movie, err
	}
	err = json.Unmarshal([]byte(val), &movie)
	if err != nil {
		log.Println("*****ERROR*****  - repo.go *** error during json.Unmarshal")
		return movie, err
	}

	log.Println("*****INFO*****  - repo.go *** movie returned successfully ")
	return movie, err
}

//GetAllMovies - return all movies from DB
func GetAllMovies(client *redis.Client) []model.Movie {
	log.Println("*****INFO*****  - repo.go *** GetAllID method initiated")
	movieIDs, err := client.SMembers("movies").Result()
	if err != nil {
		log.Println("*****ERROR*****  - repo.go ***  error retrieving movie IDs from base")
	}
	var movies []model.Movie
	for _, movieID := range movieIDs {
		movie, _ := GetByID(client, movieID)
		movies = append(movies, movie)
	}
	log.Println("*****INFO*****  - repo.go *** movies returned successfully ")
	return movies

}

//Delete from database
func Delete(client *redis.Client, id string) (model.Movie, error) {
	log.Println("*****INFO*****  - repo.go *** Delete method initiated ")
	movie, err := GetByID(client, id)
	if err == nil {
		_, err = client.SRem("movies", id).Result()
		if err != nil {
			log.Println("*****ERROR*****  - repo.go ***  error removing movie id from base")
		}
		_, err = client.Del("movie:" + id).Result()
		if err != nil {
			log.Println("*****ERROR*****  - repo.go ***  error removing movie from base")
		}
	}
	log.Println("*****INFO*****  - repo.go *** returning movie and error ")
	return movie, err
}
