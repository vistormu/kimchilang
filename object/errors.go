package object

import "fmt"


type Error struct {
    Message string
}
func (self *Error) Type() int { return ERROR_OBJ }
func (self *Error) Inspect() string { return self.Message }

func NewError(format string, a ...interface{}) *Error {
    return &Error{Message: fmt.Sprintf(format, a...)}
}
