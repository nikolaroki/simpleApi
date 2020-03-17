package main

import (
	"api/server"
	"log"
)

func main() {
	log.Println("*****INFO*****  - main.go *** App started")
	server.StartServer()
}

// Example for POST
// {
// 	"genre": "Dal",
// 	"title": "Komedija",
// 	"director": {
// 		"firstname":"Tragedija",
// 		"lastname":"Drama"
// 	}
// }
