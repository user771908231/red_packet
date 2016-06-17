// error.go
// created by kory 2014-02-10
//

package Error

import (
	"fmt"
)

//func New(errCode int, errStr string) Error {
//	return Error{errCode, errStr}
//}

// New returns an error that formats as the given text.
func New(vlist ...interface{}) Error {
	errno := 0
	errstr := ""
	for _, vv := range vlist {
		switch v := vv.(type) {
		case int:
			{
				errno = v //vlist[0].(int)
			}
		case string:
			{
				if errno == 0 {
					errno = -1
				}
				errstr = v //vlist[0].(string)
			}
		case error:
			{
				if errno == 0 {
					errno = -1
				}
				if v != nil {
					errstr = v.Error()
				}
			}
		}
	}
	return Error{errno, errstr}
}

func NewError(errCode int, errStr string) error {
	return &Error{errCode, errStr}
}

func OK() Error {
	return Error{0, ""}
}

// errorString is a trivial implementation of error.
type Error struct {
	errCode int
	errStr  string
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.errCode == 0 {
		e.errStr = "OK"
	}
	return fmt.Sprintf("[%v] %v", e.errCode, e.errStr)
}

func (e *Error) Code() int {
	if e == nil {
		return 0
	}

	return e.errCode
}

func (e *Error) IsError() bool {
	if e == nil {
		return false
	}

	return e.errCode != 0
}

func (e *Error) Assign(err error) Error {
	if err == nil {
		e.errCode = 0
		e.errStr = ""
	} else {
		e.errCode = -1
		e.errStr = err.Error()
	}
	return *e
}

func (e *Error) SetError(code int, err error) {
	e.errCode = code
	e.errStr = ""
}
