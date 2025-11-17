package lua

import (
	"unsafe"

	"github.com/CompeyDev/lei/ffi"
)

type LuaString struct {
	lua   *Lua
	index int
}

func (s *LuaString) ToString() string {
	state := s.lua.state
	return ffi.ToString(state, int32(s.index))
}

func (s *LuaString) ToPointer() unsafe.Pointer {
	state := s.lua.state
	return ffi.ToPointer(state, int32(s.index))
}

//
// LuaValue Implementation
//

func (s *LuaString) luaState() *Lua {
	return s.lua
}

func (s *LuaString) stackIndex() int {
	return s.index
}
