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

	s.deref()
	defer ffi.Pop(state, 1)

	return ffi.ToString(state, -1)
}

func (s *LuaString) ToPointer() unsafe.Pointer {
	state := s.vm.state()

	s.deref()
	defer ffi.Pop(state, 1)

	return ffi.ToPointer(state, -1)
}

//
// LuaValue implementation
//

var _ LuaValue = (*LuaString)(nil)

func (s *LuaString) lua() *Lua { return s.vm }
func (s *LuaString) ref() int  { return s.index }

func (s *LuaString) deref() int {
	return int(ffi.RawGetI(s.lua().state(), ffi.LUA_REGISTRYINDEX, int32(s.ref())))
}
