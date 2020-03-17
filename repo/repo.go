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
	return client
}

//Set a new movie
func Set(client *redis.Client, movie model.Movie) error {
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

//GetByID - return a movie by ID
func GetByID(client *redis.Client, id string) (model.Movie, error) {
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

//GetAllMovies - return all movies from DB
func GetAllMovies(client *redis.Client) []model.Movie {
	movieIDs, err := client.SMembers("movies").Result()
	if err != nil {
		log.Println(err)
	}
	var movies []model.Movie
	for _, movieID := range movieIDs {
		movie, _ := GetByID(client, movieID)
		movies = append(movies, movie)
	}
	return movies

}

//Delete from database
func Delete(client *redis.Client, id string) (model.Movie, error) {
	movie, err := GetByID(client, id)
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
