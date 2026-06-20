package main

import (
	"async-job-platform/internal/db"
	"async-job-platform/internal/handler"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := db.Connect(); err != nil {
		log.Fatal(err)
	}
	log.Println("POSTGRESQL CONNECTED")
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.POST("/jobs", func(c *gin.Context) {
		var req handler.CreateJobRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, req)
	})

	r.Run(":8080")

}
