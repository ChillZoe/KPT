package errors

import (
	"fmt"
)

type RuntimeError struct {
	Err       error
	CustomMsg string
}

func New(text string) error {
	return &RuntimeError{nil, text}
}

func (e *RuntimeError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s:\n%s", e.CustomMsg, e.Err)
	} else {
		return e.CustomMsg
	}
}
