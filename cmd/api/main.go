package main

import (
	"async-job-platform/internal/db"
	"async-job-platform/internal/handler"
	"async-job-platform/internal/job"
	"async-job-platform/internal/queue"
	"async-job-platform/internal/worker"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := db.Connect(); err != nil {
		log.Fatal(err)
	}
	if err := queue.Connect(); err != nil {
		log.Fatal(err)
	}
	log.Println("REDIS CONNECTED")
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

		id, err := job.Create(req.JobType, req.Payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := queue.Enqueue(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": id})
	})
	for i := 1; i <= 3; i++ {
		go worker.Start(i)
	}
	r.Run(":8080")

}
