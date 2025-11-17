package lua

import (
	"github.com/CompeyDev/lei/ffi"
)

type LuaValue interface {
	luaState() *Lua
	stackIndex() int
}

func TypeName(val LuaValue) string {
	lua := val.luaState()
	return ffi.TypeName(lua.state, ffi.Type(lua.state, int32(val.stackIndex())))
}
