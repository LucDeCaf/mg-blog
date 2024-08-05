package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

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

	// Page Routes
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// Special Routes
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", gin.H{})
	})

	// API Routes
	r.GET("/api/author", apiGetAuthors)
	r.POST("/api/author", apiPostAuthor)
	r.GET("/api/author/:authorId", apiGetAuthor)
	r.GET("/api/blog", apiGetBlogs)
	r.GET("/api/blog/:blogId", apiGetBlog)

	r.Run(":8080")
}

func apiGetAuthors(c *gin.Context) {
	authors, err := GetAuthors()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, authors)
}

func apiPostAuthor(c *gin.Context) {
	var a Author

	if err := c.BindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}

	if _, err := AddAuthor(a); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, a)
}

func apiGetAuthor(c *gin.Context) {
	idStr, _ := c.Params.Get("authorId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	a, err := GetAuthor(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, a)
}

func apiGetBlog(c *gin.Context) {
	idStr, _ := c.Params.Get("blogId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("bad value for param blogId '%v'", idStr),
		})
	}

	b, err := GetBlog(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, b)
}

func apiGetBlogs(c *gin.Context) {
	blogs, err := GetBlogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, blogs)
}
