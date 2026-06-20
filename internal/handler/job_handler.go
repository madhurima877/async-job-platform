package handler

type CreateJobRequest struct {
	JobType string                 `json:"job_type" binding:"required"`
	Payload map[string]interface{} `json:"payload" binding:"required"`
}
