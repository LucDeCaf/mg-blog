package main

import "time"

type Blog struct {
	Id        int       `json:"id"`
	Title     string    `json:"title" binding:"required"`
	Content   string    `json:"content" binding:"required"`
	AuthorId  int       `json:"author_id" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetBlogs() ([]Blog, error) {
	rows, err := db.Query("SELECT * FROM blog_posts;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []Blog

	for rows.Next() {
		var b Blog

		if err := rows.Scan(
			&b.Id,
			&b.Title,
			&b.Content,
			&b.AuthorId,
			&b.CreatedAt,
			&b.UpdatedAt,
		); err != nil {
			return nil, err
		}

		blogs = append(blogs, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return blogs, nil
}

func GetBlog(id int) (Blog, error) {
	b := Blog{Id: id}

	if err := db.QueryRow("SELECT (title,content,author_id,created_at,updated_at) FROM blog_posts WHERE id=?;", id).Scan(
		&b.Title,
		&b.Content,
		&b.AuthorId,
		&b.CreatedAt,
		&b.UpdatedAt,
	); err != nil {
		return Blog{}, err
	}

	return b, nil
}
