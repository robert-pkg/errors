package errors

import (
	"fmt"
)

func IsHttpStatus(err error, httpStatus int) bool {
	return Code(err) == httpStatus
}

func NewHttpCode(http int, message string) error {
	if !isSupportHttpStatus(http) {
		panic(fmt.Sprintf("http status(%d) is not support. ", http))
	}

	return &withCode{
		err:   fmt.Errorf(message),
		Coder: newHttpCoder(http, message),
		stack: callers(),
	}
}
