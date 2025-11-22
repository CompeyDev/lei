package lua

import (
	"fmt"

	"github.com/CompeyDev/lei/ffi"
)

type LoadError struct {
	Code    int
	Message string
}

func (e *LoadError) Error() string {
	switch e.Code {
	case ffi.LUA_ERRSYNTAX:
		return "syntax error: " + e.Message
	case ffi.LUA_ERRMEM:
		return "memory allocation error: " + e.Message
	case ffi.LUA_ERRERR:
		return "error handler error: " + e.Message
	default:
		return fmt.Sprintf("load error (code %d): %s", e.Code, e.Message)
	}
}

func newLoadError(state *ffi.LuaState, code int) *LoadError {
	if code != ffi.LUA_OK {
		message := ffi.ToString(state, -1)
		err := &LoadError{Code: code, Message: message}

		ffi.Pop(state, 1)

		return err
	}

	return nil
}
