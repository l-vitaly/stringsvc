package stringsvc

// This is the String service definition.

import (
	"errors"
	"strings"

	"golang.org/x/net/context"
)

var (
	ErrEmpty = errors.New("string empty")
)

type Middleware func(StringSvc) StringSvc

type StringSvc interface {
	Uppercase(context.Context, string) (string, error)
}

type stringSvc struct{}

func NewService() StringSvc {
	return &stringSvc{}
}

func (stringSvc) Uppercase(_ context.Context, s string) (string, error) {
	if s == "" {
		return s, ErrEmpty
	}
	return strings.ToUpper(s), nil
}
