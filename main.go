package main

import "api/server"

func main() {
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
