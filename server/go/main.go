package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"mg-blog/author"
	"mg-blog/blog"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	r := gin.Default()

	dbPathPtr := flag.String(
		"database-path",
		"database.db",
		"path to db, defaults to 'database.db'",
	)

	flag.Parse()

	database, err := sql.Open("sqlite3", *dbPathPtr)
	if err != nil {
		log.Fatal(err)
	}
	db = database

	// Load static files
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "static")

	// Middleware
	r.Use(gin.Logger())

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
	r.GET("/blog/:id", func(c *gin.Context) {
		idStr := c.Params.ByName("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
			return
		}

		b, err := blog.GetBlog(db, id)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
			return
		}

		c.HTML(http.StatusOK, "blog.html", b)
	})

	// API Routes (removing post for safety)
	r.GET("/api/author", apiGetAuthors)
	r.GET("/api/author/:authorId", apiGetAuthor)
	// r.POST("/api/author", apiPostAuthor)
	r.GET("/api/blog", apiGetBlogs)
	r.GET("/api/blog/:blogId", apiGetBlog)
	// r.POST("/api/blog", apiPostBlog)

	r.Run(":8080")
}

func apiGetAuthors(c *gin.Context) {
	authors, err := author.GetAuthors(db)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, authors)
}

func apiPostAuthor(c *gin.Context) {
	var a author.Author

	if err := c.BindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if _, err := author.AddAuthor(db, a); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
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
			"error": err.Error(),
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
			"error": err.Error(),
		})
		return
	}

	b, err := blog.AddBlog(db, b)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
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
