package sys

import (
	"errors"

	"github.com/arifullov/auth/internal/sys/codes"
)

type commonError struct {
	msg  string
	code codes.Code
}

func NewCommonError(code codes.Code, msg string) *commonError {
	return &commonError{code: code, msg: msg}
}

func (e *commonError) Error() string {
	return e.msg
}

func (e *commonError) Code() codes.Code {
	return e.code
}

func IsCommonError(err error) bool {
	var ce *commonError
	return errors.As(err, &ce)
}

func GetCommonError(err error) *commonError {
	var ce *commonError
	if !errors.As(err, &ce) {
		return nil
	}
	return ce
}
