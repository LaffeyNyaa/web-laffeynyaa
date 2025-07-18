package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Work struct {
	ID         int
	TitleEN    string
	SubtitleEN string
	TitleZH    string
	SubtitleZH string
	Iter       int
	ContentEN  string
	ContentZH  string
}

func main() {
	pass := os.Getenv("POSTGRES_PASS")
	connStr := fmt.Sprintf("host=localhost port=5432 user=postgres password=%s dbname=web_laffeynyaa", pass)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}

	var works []Work
	db.Find(&works)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.LoadHTMLGlob("tmpls/*")
	router.Static("static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})

	router.GET("/works", func(c *gin.Context) {
		c.HTML(http.StatusOK, "works.tmpl", gin.H{
			"Works": works,
		})
	})

	router.GET("/works/:id", func(c *gin.Context) {
		id := c.Param("id")

		var work Work
		if err := db.Where("id = ?", id).First(&work).Error; err != nil {
			panic("Falied to get data by id")
		}

		c.HTML(http.StatusOK, id+".tmpl", gin.H{
			"Content": work.ContentZH,
		})
	})

	router.Run(":80")
}
