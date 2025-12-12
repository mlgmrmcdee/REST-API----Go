package models

type Movie struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Year        int      `json:"year,omitempty"`
	DurationMin int      `json:"duration_min,omitempty"`
	Description string   `json:"description,omitempty"`
	Poster      string   `json:"poster,omitempty"`
	Certificate string   `json:"certificate,omitempty"`
	Rating      float32  `json:"rating,omitempty"`
	Genres      []string `json:"genres,omitempty"`
	Actors      []string `json:"actors,omitempty"`
	Directors   []string `json:"directors,omitempty"`
}
