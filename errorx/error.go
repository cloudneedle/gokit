package errorx

import "fmt"

type Error struct {
	code   ErrorCode
	err    error
	detail string
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) Code() ErrorCode {
	return e.code
}

func (e *Error) DetailString() string {
	return e.detail
}

func New(code ErrorCode, detail string) error {
	return &Error{
		code:   code,
		err:    fmt.Errorf("%v", code),
		detail: detail,
	}
}

func Code(code ErrorCode) error {
	return &Error{
		code:   code,
		err:    fmt.Errorf("%v", code),
		detail: "",
	}
}
