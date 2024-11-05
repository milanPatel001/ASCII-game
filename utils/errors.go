package utils

import (
	"errors"
	"fmt"
)

var (
	AUTH_ERROR = errors.New("Not Authenticated")
)

type PacketError struct {
	Code    int
	Message string
}

func (e *PacketError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}
