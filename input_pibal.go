package main

import (
	"fmt"
	"net/http"

	//	"context"
	//	"time"

	//	"regexp"

	"github.com/gin-gonic/gin"
	//	"go.mongodb.org/mongo-driver/bson"
	//	"go.mongodb.org/mongo-driver/mongo"
	//	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("template/*")

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/index")
	})

	r.GET("/index", func(c *gin.Context) {
		fmt.Println("ok")
		c.HTML(http.StatusOK, "index.html", nil)
	})

	return r
}

func main() {
	r := setupRouter()

	// Listen and Server in 0.0.0.0:8080
	r.Run(":6868")
}
