package lua

import (
	"unsafe"

	"github.com/CompeyDev/lei/ffi"
)

type LuaString struct {
	vm    *Lua
	index int
}

func (s *LuaString) ToString() string {
	state := s.vm.state()
	return ffi.ToString(state, int32(s.index))
}

func (s *LuaString) ToPointer() unsafe.Pointer {
	state := s.vm.state()
	return ffi.ToPointer(state, int32(s.index))
}

//
// LuaValue Implementation
//

func (s *LuaString) lua() *Lua {
	return s.vm
}

func (s *LuaString) stackIndex() int {
	return s.index
}
