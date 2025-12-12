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

func GetMovie(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid ID", 400)
		return
	}

	query := `
		SELECT m.movie_id, m.title, m.year, m.duration_min, m.description,
		       m.poster, c.certificate_name, m.rating
		FROM movie m
		LEFT JOIN certificate c ON m.certificate_id = c.certificate_id
		WHERE m.movie_id = $1
	`

	var m models.Movie
	var year, duration sql.NullInt64
	var rating sql.NullFloat64
	var desc, poster, cert sql.NullString

	err = db.DB.QueryRow(query, id).Scan(&m.ID, &m.Title,
		&year, &duration, &desc, &poster, &cert, &rating)
	if err != nil {
		log.Println("Ошибка запроса к БД:", err)
		http.Error(w, "Movie not found", 404)
		return
	}

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

	m.Directors = utils.FetchStringList(`
			SELECT d.name FROM movie_director md
			JOIN director d ON md.director_id = d.director_id
			WHERE md.movie_id = $1`, m.ID)

	m.Actors = utils.FetchStringList(`
			SELECT a.name FROM movie_cast mc
			JOIN actor a ON mc.actor_id = a.actor_id
			WHERE mc.movie_id = $1`, m.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}
