package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Blog struct {
	Id          int
	Title       string
	ContentHTML string
}

func main() {
	r := gin.Default()

	// Load static files
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.GET("/blog/:article_id", func(ctx *gin.Context) {
		articleId, _ := ctx.Params.Get("article_id")

		ctx.JSON(http.StatusOK, gin.H{"articleId": articleId})
	})

	r.Run(":8080")
}
