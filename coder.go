package errors

import (
	"fmt"
	"net/http"
	"sync"
)

var (
	MinCode               = 100000
	UnknownHttpStatus     = http.StatusInternalServerError
	UnknownCode           = UnknownHttpStatus
	SuccessCode           = 200
	SupportHttpStatusList = []int{400, 401, 403, 404, 408, 409, 429, 500, 503, 504}
)

type Coder interface {
	HTTPStatus() int
	Code() int
	Message() string
}

type defaultCoder struct {
	http int
	code int
	msg  string
}

func (coder defaultCoder) Code() int {
	return coder.code
}

func (coder defaultCoder) HTTPStatus() int {
	if coder.http == 0 {
		return UnknownHttpStatus
	}

	return coder.http
}

func (coder defaultCoder) Message() string {
	return coder.msg
}

func newCoder(code int, message string) Coder {
	http := UnknownHttpStatus
	if v, ok := codes[code]; ok {
		http = v
	}

	return &defaultCoder{
		http: http,
		code: code,
		msg:  message,
	}
}

func newHttpCoder(httpSts int, message string) Coder {
	return &defaultCoder{
		http: httpSts,
		code: httpSts,
		msg:  message,
	}
}

func newUnknownCoder(message string) Coder {
	return &defaultCoder{
		http: UnknownHttpStatus,
		code: UnknownCode,
		msg:  message,
	}
}

var codes = map[int]int{} // key是code， value是http码
var codeMux = &sync.Mutex{}

func isSupportHttpStatus(httpStatus int) bool {
	for _, v := range SupportHttpStatusList {
		if v == httpStatus {
			return true
		}
	}

	return false
}

func Register(code int, http int) {
	if code < MinCode {
		panic(fmt.Sprintf("code must be great than %d.", MinCode))
	}

	if !isSupportHttpStatus(http) {
		panic(fmt.Sprintf("http status(%d) is not support. ", http))
	}

	SupportHttpStatusList = []int{400, 401, 403, 404, 408, 409, 429, 500, 503, 504}

	codeMux.Lock()
	defer codeMux.Unlock()

	codes[code] = http
}

func MustRegister(code int, http int) {
	if code < MinCode {
		panic(fmt.Sprintf("code must be great than %d.", MinCode))
	}

	if !isSupportHttpStatus(http) {
		panic(fmt.Sprintf("http status(%d) is not support. ", http))
	}

	codeMux.Lock()
	defer codeMux.Unlock()

	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("code: %d already exist", code))
	}

	codes[code] = http
}

func ParseCoder(err error) Coder {
	if err == nil {
		return nil
	}

	if v, ok := err.(Coder); ok {
		return v
	}

	return newUnknownCoder(err.Error())
}

func Code(err error) int {
	if err == nil {
		return SuccessCode
	}

	if v, ok := err.(Coder); ok {
		return v.Code()
	}

	return UnknownCode
}
