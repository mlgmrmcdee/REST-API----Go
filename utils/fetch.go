package utils

import (
	"log"

	"rest-api/db"
)

func FetchStringList(query string, id int) []string {
	rows, err := db.DB.Query(query, id)
	if err != nil {
		log.Println("Ошибка запроса к БД:", err)
		return []string{}
	}
	defer rows.Close()

	var list []string
	for rows.Next() {
		var s string
		rows.Scan(&s)
		list = append(list, s)
	}
	return list
}
