package handlers

import (
	"encoding/json"
	"net/http"

	"rest-api/db"
)

func GetGenres(w http.ResponseWriter, r *http.Request) {
	rows, _ := db.DB.Query(`SELECT genre_name FROM genre ORDER BY genre_name ASC`)
	defer rows.Close()

	var list []string
	for rows.Next() {
		var g string
		rows.Scan(&g)
		list = append(list, g)
	}

	json.NewEncoder(w).Encode(list)
}
