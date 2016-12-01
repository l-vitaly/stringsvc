package stringsvctest

import (
	"testing"

	"github.com/l-vitaly/stringsvc"
)

var uppercaseProvider = []struct {
	value    string
	expected string
}{
	{"test", "TEST"},
	{"foo", "FOO"},
	{"BAR", "BAR"},
}

func TestUppercase(t *testing.T) {
	svc := stringsvc.NewService()

	for _, tt := range uppercaseProvider {
		actual, err := svc.Uppercase(nil, tt.value)
		if err != nil {
			t.Errorf("Uppercase: %s", err)
			return
		}
		if actual != tt.expected {
			t.Errorf("Uppercase(%s): expected %s, actual %s", tt.value, tt.expected, actual)
			return
		}
	}
}
