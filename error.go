package gservice

import "fmt"

var (
	ErrConn         = fmt.Errorf("conn fail")
	ErrErrNoHandler = fmt.Errorf("manager didn't handle Errs()")
)
