package responses

// APIResponse represents the response for errors
type APIResponse struct {
	Code    int    `json:"code" example:"400" doc:"HTTP status code"`
	Type    string `json:"type" example:"validation_error" doc:"Error type (validation_error, error, etc.)"`
	Message string `json:"message" example:"invalid request" doc:"Human-readable error message"`
} //@name ApiResponse
