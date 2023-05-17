package errors

import (
	"fmt"
	"net/http"
	"testing"
)

func printCode(err error) {
	// %s: Returns the user-safe error string mapped to the error code or the error message if none is specified.
	fmt.Printf("%s\n", "====================> %s <====================")
	fmt.Printf("%s\n\n", err)

	// %v: Alias for %s.
	fmt.Printf("%s\n", "====================> %v <====================")
	fmt.Printf("%v\n\n", err)

	// %-v: Output caller details, useful for troubleshooting.
	fmt.Printf("%s\n", "====================> %-v <====================")
	fmt.Printf("%-v\n\n", err)

	fmt.Printf("%s\n", "====================> %+v <====================")
	str := fmt.Sprintf("%+v", err)
	fmt.Printf("%s\n\n", str)

}

func getSimpleCode() error {
	return WithCode(100, "simple fail.")
}

func TestSimpleCoder(t *testing.T) {
	err := getSimpleCode()
	printCode(err)
}

func getSimpleError() error {
	return New("simple error")
}

func getWrappError() error {
	err := getSimpleError()

	//return Wrap(err, "check Info fail")
	return WrapCode(err, 100000, "check Info fail")
}

func TestWrapCode(t *testing.T) {
	err := getWrappError()
	printCode(err)
}

func TestMakeResponse(t *testing.T) {

	type Response struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	makeFunc := func(err error, data interface{}) (int, *Response) {

		if err == nil {
			return http.StatusOK, &Response{
				Data: data,
			}
		}

		coder := ParseCoder(err)
		return coder.HTTPStatus(), &Response{
			Code:    coder.Code(),
			Message: coder.Message(),
			Data:    data,
		}
	}

	if true {
		err := getSimpleCode()
		httpStatus, resp := makeFunc(err, nil)
		fmt.Printf("http:%d \n", httpStatus)
		fmt.Printf("resp:%v \n", resp)
	}

	if true {

		Register(100000, 400)
		err := getWrappError()
		httpStatus, resp := makeFunc(err, nil)
		fmt.Printf("http:%d \n", httpStatus)
		fmt.Printf("resp:%v \n", resp)
	}

}
