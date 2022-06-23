package models

import (
	"errors"
	"fmt"
)

func NewError(format string, v ...interface{}) error {
	msg := fmt.Sprintf(format, v...)
	return errors.New(msg)
}
