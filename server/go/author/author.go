package author

import (
	"database/sql"
	"fmt"
)

type Author struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

func GetAuthors(db *sql.DB) ([]Author, error) {
	if db == nil {
		return nil, fmt.Errorf("db is null")
	}

	rows, err := db.Query("SELECT * FROM authors;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []Author

	for rows.Next() {
		var a Author

		if err := rows.Scan(&a.Id, &a.FirstName, &a.LastName); err != nil {
			return nil, err
		}

		authors = append(authors, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return authors, err
}

func GetAuthor(db *sql.DB, id int) (Author, error) {
	if db == nil {
		return Author{}, fmt.Errorf("db is null")
	}

	a := Author{Id: id}

	if err := db.QueryRow("SELECT first_name,last_name FROM authors WHERE id=?;", id).Scan(&a.FirstName, &a.LastName); err != nil {
		return Author{}, err
	}

	return a, nil
}

func AddAuthor(db *sql.DB, a Author) (Author, error) {
	if db == nil {
		return Author{}, fmt.Errorf("db is null")
	}

	r, err := db.Exec("INSERT INTO authors (first_name,last_name) VALUES (?,?);", a.FirstName, a.LastName)
	if err != nil {
		return Author{}, err
	}
	id, _ := r.LastInsertId()
	a, _ = GetAuthor(db, int(id))

	return a, nil
}
