Client
   │
   ▼
Gin APIs
   │
   ├── POST /jobs
   ├── GET /jobs/:id
   └── GET /metrics
   │
   ▼
PostgreSQL
(Source of Truth)
   │
   ▼
Redis Queue
   │
   ▼
Worker Pool
 ┌─────────────┐
 │ Worker-1    │
 │ Worker-2    │
 │ Worker-3    │
 └─────────────┘
   │
   ▼
completed / failed
