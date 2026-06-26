package main

import (
	"async-job-platform/internal/db"
	"async-job-platform/internal/handler"
	"async-job-platform/internal/job"
	"async-job-platform/internal/queue"
	"async-job-platform/internal/worker"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

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
	r.GET("/jobs/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid job id"})
			return
		}
		j, err := job.GetJob(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
			return
		}
		c.JSON(http.StatusOK, j)
	})

	r.GET("/metrics", func(c *gin.Context) {
		allow, err := queue.CheckAllowed("ratelimiter")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error: " + err.Error()})
			return
		}
		if !allow {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})

			return
		}

		m, err := job.GetMetrics()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, m)
	})
	for i := 1; i <= 3; i++ {
		go worker.Start(i)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Server exited")

}

// func main() {
// 	var mu sync.Mutex
// 	if err := queue.Connect(); err != nil {
// 		log.Fatal(err)
// 	}
// 	numjobs := 10
// 	var wg sync.WaitGroup
// 	sem := make(chan struct{}, 5)
// 	queue.SetKey("key")
// 	for i := 0; i < numjobs; i++ {
// 		wg.Add(1)
// 		sem <- struct{}{}
// 		go func() {
// 			defer wg.Done()
// 			defer func() {
// 				<-sem
// 			}()
// 			for i := 0; i <= 1000; i++ {
// 				mu.Lock()
// 				queue.Increment("key")
// 				mu.Unlock()
// 			}

// 		}()

// 	}
// 	wg.Wait()
// 	val, _ := queue.GetValue("key")
// 	log.Println(val, "data from redis")

// }

// func main() {
// 	if err := queue.Connect(); err != nil {
// 		log.Fatal(err)
// 	}

// 	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != http.MethodGet {
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 			return
// 		}

// 		allow, err := queue.CheckAllowed("ratelimiter")
// 		if err != nil {
// 			http.Error(w, "internal server error: "+err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		if !allow {
// 			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
// 			return
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		w.Write([]byte(`{"result":"success"}`))
// 	})

// 	log.Println("listening")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }
