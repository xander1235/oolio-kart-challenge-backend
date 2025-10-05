package exceptions

import (
	"net/http"
	"oolio.com/kart/exceptions/errors"
	"time"
)

func BadRequestException(message string) *errors.ErrorDetails {
	return &errors.ErrorDetails{
		ErrorTimestamp: time.Now().UnixMilli(),
		Message:        message,
		ErrorCode:      http.StatusBadRequest,
	}
}
