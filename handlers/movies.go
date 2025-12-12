package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"rest-api/db"
	"rest-api/models"
	"rest-api/utils"
)

func GetMovies(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 25
	}

	offset := (page - 1) * limit

	query := `
		SELECT m.movie_id, m.title, m.year, m.duration_min, m.description,
		       m.poster, c.certificate_name, m.rating
		FROM movie m
		LEFT JOIN certificate c ON m.certificate_id = c.certificate_id
		ORDER BY m.rating DESC NULLS LAST
		LIMIT $1 OFFSET $2
	`

	rows, err := db.DB.Query(query, limit, offset)
	if err != nil {
		log.Println("Ошибка запроса к БД:", err)
		http.Error(w, "DB error", 500)
		return
	}
	defer rows.Close()

	var movies []models.Movie

	for rows.Next() {
		var m models.Movie
		var year, duration sql.NullInt64
		var rating sql.NullFloat64
		var desc, poster, cert sql.NullString

		rows.Scan(&m.ID, &m.Title, &year, &duration, &desc, &poster, &cert, &rating)

		if year.Valid {
			m.Year = int(year.Int64)
		}
		if duration.Valid {
			m.DurationMin = int(duration.Int64)
		}
		if desc.Valid {
			m.Description = desc.String
		}
		if poster.Valid {
			m.Poster = poster.String
		}
		if cert.Valid {
			m.Certificate = cert.String
		}
		if rating.Valid {
			m.Rating = float32(rating.Float64)
		}

		m.Genres = utils.FetchStringList(`
			SELECT g.genre_name FROM movie_genre mg
			JOIN genre g ON mg.genre_id = g.genre_id
			WHERE mg.movie_id = $1`, m.ID)

		movies = append(movies, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}
