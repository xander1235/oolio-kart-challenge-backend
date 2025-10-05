package exceptions

import (
	"net/http"
	"oolio.com/kart/exceptions/errors"
	"time"
)

func UnprocessableEntityException(message string) *errors.ErrorDetails {
	return &errors.ErrorDetails{
		ErrorTimestamp: time.Now().UnixMilli(),
		Message:        message,
		ErrorCode:      http.StatusUnprocessableEntity,
	}
}
