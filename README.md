# Async Job Processing Platform

A distributed asynchronous job processing system built with Go, Gin, PostgreSQL, and Redis.

## Features

* Create background jobs via REST APIs
* Persist jobs in PostgreSQL
* Queue jobs using Redis
* Concurrent worker pool using goroutines
* Job status tracking
* Retry mechanism with backoff
* Failed job handling
* Metrics endpoint
* Graceful shutdown

## Tech Stack

* Go
* Gin
* PostgreSQL
* Redis
* Docker
* Git

## APIs

### Create Job

POST /jobs

Request:

{
"job_type": "email",
"payload": {
"to": "[test@example.com](mailto:test@example.com)",
"subject": "Welcome"
}
}

### Get Job

GET /jobs/:id

Response:

{
"id": 33,
"job_type": "email",
"status": "completed",
"retry_count": 0
}

### Metrics

GET /metrics

Response:

{
"queued": 5,
"processing": 3,
"completed": 32,
"failed": 12
}
