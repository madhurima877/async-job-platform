package worker

import (
	"async-job-platform/internal/job"
	"async-job-platform/internal/queue"
	"log"
	"strconv"
	"time"
)

func Start(id int) {
	for {
		result, err := queue.Client.BLPop(queue.CTX, 0, "jobs").Result()
		if err != nil {
			log.Println(err)
			continue
		}
		jobID, err := strconv.ParseInt(result[1], 10, 64)
		if err != nil {
			log.Println(err)
			continue
		}
		j, err := job.GetByID(jobID)
		if err != nil {
			log.Println(err)
			continue
		}
		if err := job.UpdateStatus(jobID, "processing"); err != nil {
			log.Println(err)
			continue
		}

		log.Printf("Worker-%d executing job ID=%d", id, jobID)
		if jobID%2 == 0 {
			log.Printf("Worker -%d failed job ID=%d", id, jobID)

			if j.RetryCount >= 3 {
				if err := job.UpdateStatus(jobID, "failed"); err != nil {
					log.Println(err)
				}
				log.Printf("Worker-%d permanently failed job ID=%d", id, jobID)
				continue
			}
			if err := job.IncrementRetryCount(jobID); err != nil {
				log.Println(err)
				continue

			}
			time.Sleep(2 * time.Second)

			if err := queue.Enqueue(jobID); err != nil {
				log.Println(err)
				continue
			}
			continue

		}
		if err := job.UpdateStatus(jobID, "completed"); err != nil {
			log.Println(err)
			continue
		}
		log.Printf(
			"Processing job ID=%d Type=%s",
			j.ID,
			j.JobType,
		)

	}
}

// . Recognizing and explaining race conditions i
