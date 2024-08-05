package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Blog struct {
	Id        int       `json:"id"`
	Title     string    `json:"title" binding:"required"`
	Content   string    `json:"content" binding:"required"`
	AuthorId  int       `json:"author_id" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Author struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

var db *sql.DB

func main() {
	r := gin.Default()

	dbPath, found := os.LookupEnv("BLOG_DB_PATH")
	if !found {
		dbPath = "./database.db"
	}
	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	db = database

	// Load static files
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	// HTML pages
	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.GET("/blog/:article_id", func(ctx *gin.Context) {
		articleId, _ := ctx.Params.Get("article_id")

		ctx.JSON(http.StatusOK, gin.H{"articleId": articleId})
	})

	// API routes
	r.GET("/api/author", func(ctx *gin.Context) {
		authors, err := authors()
		if err != nil {
			fmt.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}
		ctx.JSON(http.StatusOK, authors)
	})
	r.GET("/api/author/:authorId", func(ctx *gin.Context) {
		idStr, _ := ctx.Params.Get("authorId")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		a, err := author(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, a)
	})
	r.POST("/api/author", func(ctx *gin.Context) {
		var a Author

		if err := ctx.BindJSON(&a); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "bad request",
			})
			return
		}

		if _, err := addAuthor(a); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}

		ctx.JSON(http.StatusCreated, a)
	})

	r.Run(":8080")
}

func authors() ([]Author, error) {
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

	err = rows.Err()
	return authors, err
}

func author(id int) (Author, error) {
	a := Author{}
	a.Id = id

	if err := db.QueryRow("SELECT first_name,last_name FROM authors WHERE id=?;", id).Scan(&a.FirstName, &a.LastName); err != nil {
		return Author{}, err
	}

	return a, nil
}

func addAuthor(a Author) (Author, error) {
	r, err := db.Exec("INSERT INTO authors (first_name,last_name) VALUES (?,?);", a.FirstName, a.LastName)
	if err != nil {
		return Author{}, err
	}
	id, _ := r.LastInsertId()
	a, _ = author(int(id))

	return a, nil
}
