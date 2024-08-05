package main

type Author struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

func GetAuthors() ([]Author, error) {
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

func GetAuthor(id int) (Author, error) {
	a := Author{Id: id}

	if err := db.QueryRow("SELECT first_name,last_name FROM authors WHERE id=?;", id).Scan(&a.FirstName, &a.LastName); err != nil {
		return Author{}, err
	}

	return a, nil
}

func AddAuthor(a Author) (Author, error) {
	r, err := db.Exec("INSERT INTO authors (first_name,last_name) VALUES (?,?);", a.FirstName, a.LastName)
	if err != nil {
		return Author{}, err
	}
	id, _ := r.LastInsertId()
	a, _ = GetAuthor(int(id))

	return a, nil
}
