package main

import (
	"log"
	"net/http"

	"rest-api/db"
	"rest-api/handlers"
)

func main() {
	// Подключение БД
	err := db.Connect()
	if err != nil {
		log.Fatalf("Ошибка БД: %v", err)
	}

	mux := http.NewServeMux()

	// Раздача статики
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Главная
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	// Эндпоинты
	mux.HandleFunc("/movies", handlers.GetMovies)
	mux.HandleFunc("/movie", handlers.GetMovie)
	mux.HandleFunc("/movies/filter", handlers.GetMoviesFiltered)
	mux.HandleFunc("/genres", handlers.GetGenres)
	mux.HandleFunc("/certificates", handlers.GetCertificates)

	log.Println("Сервер запущен на :8080")
	err = http.ListenAndServe(":8080", handlers.Logging(mux))
	if err != nil {
		log.Fatalf("Ошибка сервера: %v", err)
	}
}
