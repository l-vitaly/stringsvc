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

type StringService interface {
	Uppercase(context.Context, string) (string, error)
}

type stringService struct{}

func NewService() *stringService {
	return &stringService{}
}

func (*stringService) Uppercase(_ context.Context, s string) (string, error) {
	if s == "" {
		return s, ErrEmpty
	}
	return strings.ToUpper(s), nil
}
