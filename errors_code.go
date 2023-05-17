package errors

import (
	"fmt"
	"io"
)

// WithCode annotates err with a new message.
// If err is nil, WithCode returns nil.
func WithCode(code int, message string) error {
	return &withCode{
		err:   fmt.Errorf(message),
		Coder: newCoder(code, message),
		stack: callers(),
	}
}

func WithCodef(code int, format string, args ...interface{}) error {
	return &withCode{
		err:   fmt.Errorf(format, args...),
		Coder: newCoder(code, fmt.Sprintf(format, args...)),
		stack: callers(),
	}
}

func WrapCode(err error, code int, message string) error {
	if err == nil {
		return nil
	}

	return &withCode{
		err:   fmt.Errorf(message),
		cause: err,
		Coder: newCoder(code, message),
		stack: callers(),
	}
}

func WrapCodef(err error, code int, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	message := fmt.Sprintf(format, args...)
	return &withCode{
		err:   fmt.Errorf(message),
		cause: err,
		Coder: newCoder(code, message),
		stack: callers(),
	}
}

type withCode struct {
	err   error
	cause error
	Coder
	*stack
}

func (w *withCode) Error() string { return w.Message() }
func (w *withCode) Cause() error  { return w.cause }

// Unwrap provides compatibility for Go 1.13 error chains.
func (w *withCode) Unwrap() error { return w.cause }

func (w *withCode) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			if w.cause != nil {
				fmt.Fprintf(s, "%+v\n", w.cause)
			}

			fmt.Fprintf(s, "%s (%d)", w.Message(), w.Code())
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Message())
	}
}
