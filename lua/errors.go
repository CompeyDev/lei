package lua

import (
	"fmt"

	"github.com/CompeyDev/lei/ffi"
)

type LuaError struct {
	Code    int
	Message string
}

func (e *LuaError) Error() string {
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

func newLuaError(state *ffi.LuaState, code int) *LuaError {
	if code != ffi.LUA_OK {
		message := ffi.ToString(state, -1)
		err := &LuaError{Code: code, Message: message}

		ffi.Pop(state, 1)

		return err
	}

	return nil
}
