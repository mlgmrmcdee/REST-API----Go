package handlers

import (
	"encoding/json"
	"net/http"

	"rest-api/db"
)

func GetCertificates(w http.ResponseWriter, r *http.Request) {
	rows, _ := db.DB.Query(`SELECT certificate_name FROM certificate ORDER BY certificate_name ASC`)
	defer rows.Close()

	var list []string
	for rows.Next() {
		var c string
		rows.Scan(&c)
		list = append(list, c)
	}

	json.NewEncoder(w).Encode(list)
}
