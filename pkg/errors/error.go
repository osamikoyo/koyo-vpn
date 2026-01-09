package errors

import "fmt"

type Error struct {
	Sender  string
	Fatal   bool
	Message string
}

func NewError(sender, message string, fatal bool) Error {
	return Error{
		Sender: sender,
		Fatal: fatal,
		Message: message,
	}
}

func (e *Error) Error() string {
	if e.Fatal {
		return fmt.Sprintf("%s FATAL: %s", e.Sender, e.Message)
	} else {
		return fmt.Sprintf("%s: %s", e.Sender, e.Message)
	}
}
