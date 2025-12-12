package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"rest-api/db"
	"rest-api/models"
	"rest-api/utils"
	"strconv"
	"strings"
)

func GetMoviesFiltered(w http.ResponseWriter, r *http.Request) {
	// query params
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	genre := r.URL.Query().Get("genre")
	cert := r.URL.Query().Get("certificate")
	yearFromStr := r.URL.Query().Get("year_from")
	yearToStr := r.URL.Query().Get("year_to")
	sort := r.URL.Query().Get("sort")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// base query with DISTINCT to avoid duplicates when joining genres
	var sb strings.Builder
	sb.WriteString(`
	SELECT DISTINCT m.movie_id, m.title, m.year, m.duration_min, m.description,
	       m.poster, c.certificate_name, m.rating
	FROM movie m
	LEFT JOIN certificate c ON m.certificate_id = c.certificate_id
	LEFT JOIN movie_genre mg ON m.movie_id = mg.movie_id
	LEFT JOIN genre g ON mg.genre_id = g.genre_id
	WHERE 1=1
	`)

	args := []interface{}{}
	argIndex := 1

	if genre != "" {
		sb.WriteString(" AND EXISTS (SELECT 1 FROM movie_genre mg2 JOIN genre g2 ON mg2.genre_id = g2.genre_id WHERE mg2.movie_id = m.movie_id AND g2.genre_name = $" + strconv.Itoa(argIndex) + ")")
		args = append(args, genre)
		argIndex++
	}
	if cert != "" {
		sb.WriteString(" AND c.certificate_name = $" + strconv.Itoa(argIndex))
		args = append(args, cert)
		argIndex++
	}
	if yearFromStr != "" {
		yearFrom, err := strconv.Atoi(yearFromStr)
		if err == nil {
			sb.WriteString(" AND m.year >= $" + strconv.Itoa(argIndex))
			args = append(args, yearFrom)
			argIndex++
		}
	}
	if yearToStr != "" {
		yearTo, err := strconv.Atoi(yearToStr)
		if err == nil {
			sb.WriteString(" AND m.year <= $" + strconv.Itoa(argIndex))
			args = append(args, yearTo)
			argIndex++
		}
	}

	switch sort {
	case "title_asc":
		sb.WriteString(" ORDER BY m.title ASC")
	case "title_desc":
		sb.WriteString(" ORDER BY m.title DESC")
	case "year_asc":
		sb.WriteString(" ORDER BY m.year ASC NULLS LAST")
	case "year_desc":
		sb.WriteString(" ORDER BY m.year DESC NULLS LAST")
	case "rating_asc":
		sb.WriteString(" ORDER BY m.rating ASC NULLS LAST")
	default:
		sb.WriteString(" ORDER BY m.rating DESC NULLS LAST")
	}

	// LIMIT/OFFSET via fmt to avoid messing args indexing
	sb.WriteString(" LIMIT " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(offset))

	query := sb.String()

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		log.Println("Ошибка запроса к базе (filter):", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		var rating sql.NullFloat64
		var certificate sql.NullString
		var description sql.NullString
		var poster sql.NullString
		var duration sql.NullInt64
		var year sql.NullInt64

		if err := rows.Scan(&m.ID, &m.Title, &year, &duration, &description, &poster, &certificate, &rating); err != nil {
			log.Println("Ошибка чтения данных (filter):", err)
			continue
		}

		if year.Valid {
			m.Year = int(year.Int64)
		}
		if duration.Valid {
			m.DurationMin = int(duration.Int64)
		}
		if description.Valid {
			m.Description = description.String
		}
		if poster.Valid {
			m.Poster = poster.String
		}
		if certificate.Valid {
			m.Certificate = certificate.String
		}
		if rating.Valid {
			m.Rating = float32(rating.Float64)
		}

		m.Genres = utils.FetchStringList(`
			SELECT g.genre_name
			FROM movie_genre mg
			JOIN genre g ON mg.genre_id = g.genre_id
			WHERE mg.movie_id = $1
		`, m.ID)

		movies = append(movies, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}
