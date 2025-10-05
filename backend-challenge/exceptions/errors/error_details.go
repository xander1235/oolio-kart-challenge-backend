package errors

type ErrorDetails struct {
	ErrorTimestamp int64  `json:"timestamp"`
	Message        string `json:"error_message"`
	ErrorCode      int    `json:"error_code"`
}
