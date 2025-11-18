package lua

import "github.com/CompeyDev/lei/ffi"

type LuaValue interface {
	lua() *Lua
	stackIndex() int
}

func TypeName(val LuaValue) string {
	lua := val.lua().state()
	return ffi.TypeName(lua, ffi.Type(lua, int32(val.stackIndex())))
}
