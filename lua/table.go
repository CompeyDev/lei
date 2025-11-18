package lua

import "github.com/CompeyDev/lei/ffi"

type LuaTable struct {
	vm    *Lua
	index int
}

func (t *LuaTable) Set(key LuaValue, value LuaValue) {
	state := t.vm.state()

	ffi.PushValue(state, int32(key.stackIndex()))
	ffi.PushValue(state, int32(value.stackIndex()))
	ffi.SetTable(state, int32(t.index))
}

func (t *LuaTable) Get(key LuaValue) LuaValue {
	state := t.vm.state()

	ffi.PushValue(state, int32(key.stackIndex()))
	valueType := ffi.GetTable(state, int32(t.index))

	switch valueType {
	// TODO: other types
	case ffi.LUA_TSTRING:
		return &LuaString{vm: t.vm, index: int(ffi.GetTop(state))}
	default:
		panic("Unknown type")
	}
}

//
// LuaValue Implementation
//

func (s *LuaTable) lua() *Lua {
	return s.vm
}

func (s *LuaTable) stackIndex() int {
	return s.index
}
