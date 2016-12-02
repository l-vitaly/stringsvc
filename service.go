package stringsvc

// This is the String service definition.

import (
	"errors"
	"strings"

	"golang.org/x/net/context"
)

var (
	// ErrEmpty string empty error
	ErrEmpty = errors.New("string empty")
)

// Middleware for StringSvc
type Middleware func(StringSvc) StringSvc

// StringSvc service interface
type StringSvc interface {
	Uppercase(context.Context, string) (string, error)
}

type stringSvc struct{}

// NewService create new StringSvc service
func NewService() StringSvc {
	return &stringSvc{}
}

func (stringSvc) Uppercase(_ context.Context, s string) (string, error) {
	if s == "" {
		return s, ErrEmpty
	}
	return strings.ToUpper(s), nil
}
