package exceptions

import (
	"oolio.com/kart/exceptions/errors"
	"time"
)

func GenericException(message string, errorCode int) *errors.ErrorDetails {
	return &errors.ErrorDetails{
		ErrorTimestamp: time.Now().UnixMilli(),
		Message:        message,
		ErrorCode:      errorCode,
	}
}
