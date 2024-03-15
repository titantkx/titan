package client

import (
	"errors"
	"strings"
)

var ErrDeadlineExceeded = errors.New("context deadline exceeded")

func IsNotFound(err error) bool {
	return strings.Contains(err.Error(), "not found")
}
