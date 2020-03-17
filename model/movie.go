package model

// Movie struct export
type Movie struct {
	ID       string    `json:"id" xml:"id"`
	Genre    string    `json:"genre" xml:"genre"`
	Title    string    `json:"title" xml:"title"`
	Director *Director `json:"director" xml:"director"`
}
