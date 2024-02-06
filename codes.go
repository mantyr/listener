package listener

import (
	"github.com/mantyr/listener/codes"
)

const (
	OK              codes.Code = 0
	Canceled        codes.Code = 1
	Unknown         codes.Code = 2
	InvalidArgument codes.Code = 3
	Invalid         codes.Code = 4
	NotFound        codes.Code = 5
	AlreadyExists   codes.Code = 6
	Aborted         codes.Code = 7
	Unimplemented   codes.Code = 8
	Internal        codes.Code = 9
	Unavailable     codes.Code = 10

	InvalidEvent              codes.Code = 11
	EventDecoderAlreadyExists codes.Code = 12
	EventDecoderNotFound      codes.Code = 13
)

func Code(err error) codes.Code {
	if err == nil {
		return OK
	}
	e, ok := err.(*Error)
	if !ok {
		return Unknown
	}
	return e.Code
}

func AlreadyExistsType(err error) string {
	if err == nil {
		return ""
	}
	e, ok := err.(*Error)
	if !ok {
		return ""
	}
	return e.AlreadyExistsType
}

type Error struct {
	Code              codes.Code
	Err               error
	AlreadyExistsType string
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Err == nil {
		return "empty error"
	}
	return e.Err.Error()
}
