package lua

import "github.com/CompeyDev/lei/ffi"

type LuaTable struct {
	vm    *Lua
	index int
}

func (t *LuaTable) Set(key LuaValue, value LuaValue) {
	state := t.vm.state()

	t.deref()     // table (-3)
	key.deref()   // key   (-2)
	value.deref() // value (-1)

	ffi.SetTable(state, -3)
	ffi.Pop(state, 1)
}

func (t *LuaTable) Get(key LuaValue) LuaValue {
	state := t.vm.state()

	t.deref()   //////////////////// table (-3)
	key.deref() //////////////////// key   (-2)
	ffi.GetTable(state, -2)

	val := intoLuaValue(t.vm, -1) // value (-1)
	ffi.Pop(state, 2)

	return val
}

func (t *LuaTable) Iterable() map[LuaValue]LuaValue {
	state := t.vm.state()

	t.deref()
	tableIndex := ffi.GetTop(state)
	ffi.PushNil(state)

	obj := make(map[LuaValue]LuaValue)
	for ffi.Next(state, tableIndex) != 0 {
		key, value := intoLuaValue(t.vm, -2), intoLuaValue(t.vm, -1)
		obj[key] = value

		ffi.Pop(state, 1) // only pop value, leave key in place
	}

	ffi.Pop(state, 1)
	return obj
}

//
// LuaValue implementation
//

func (t *LuaTable) lua() *Lua { return t.vm }
func (t *LuaTable) ref() int  { return t.index }

func (t *LuaTable) deref() int {
	return int(ffi.RawGetI(t.lua().state(), ffi.LUA_REGISTRYINDEX, int32(t.ref())))
}
