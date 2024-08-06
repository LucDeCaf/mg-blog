package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"mg-blog/author"
	"mg-blog/blog"

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
		blogs, err := blog.GetBlogs(db)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
			return
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"blogs": blogs,
		})
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
	r.POST("/api/blog", apiPostBlog)
	r.GET("/api/blog/:blogId", apiGetBlog)

	r.Run(":8080")
}

func apiGetAuthors(c *gin.Context) {
	authors, err := author.GetAuthors(db)
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
	var a author.Author

	if err := c.BindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}

	if _, err := author.AddAuthor(db, a); err != nil {
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

	a, err := author.GetAuthor(db, id)
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
		return
	}

	b, err := blog.GetBlog(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, b)
}

func apiPostBlog(c *gin.Context) {
	var b blog.Blog

	if err := c.BindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}

	if _, err := blog.AddBlog(db, b); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, b)
}

func apiGetBlogs(c *gin.Context) {
	blogs, err := blog.GetBlogs(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, blogs)
}
