package errs

import (
	"errors"
)

var (
	ErrorInvalidSessionSlice = errors.New("session slice length error")
	ErrorNilConn             = errors.New("conn is nil")
	ErrorInvalidUid          = errors.New("invalid uid")
	ErrorNilSessionManager   = errors.New("session manager is nil")
	ErrorNilContext          = errors.New("context is nil")
	ErrorNil                 = errors.New("is nil")
)
