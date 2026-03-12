package lua

import "github.com/CompeyDev/lei/ffi"

type LuaTable struct {
	vm    *Lua
	index int
}

func (t *LuaTable) Set(key LuaValue, value LuaValue) {
	state := t.vm.state()

	t.deref(t.vm)     // table (-3)
	key.deref(t.vm)   // key   (-2)
	value.deref(t.vm) // value (-1)

	// Pop the table off
	defer ffi.Pop(state, 1)

	ffi.SetTable(state, -3)
}

func (t *LuaTable) Get(key LuaValue) LuaValue {
	state := t.vm.state()

	t.deref(t.vm)   //////////////////// table (-3)
	key.deref(t.vm) //////////////////// key   (-2)

	// Pop the table and value off
	defer ffi.Pop(state, 2)

	ffi.GetTable(state, -2)
	val := intoLuaValue(t.vm, -1) ////// value (-1)

	return val
}

func (t *LuaTable) RawSet(key LuaValue, value LuaValue) {
	state := t.vm.state()

	t.deref(t.vm)     // table (-3)
	key.deref(t.vm)   // key   (-2)
	value.deref(t.vm) // value (-1)

	// Pop the table off
	defer ffi.Pop(state, 1)

	ffi.RawSet(state, -3)
}

func (t *LuaTable) RawGet(key LuaValue) LuaValue {
	state := t.vm.state()

	t.deref(t.vm)   // table (-2)
	key.deref(t.vm) // key   (-1)

	// Pop the table and value off
	defer ffi.Pop(state, 2)

	ffi.RawGet(state, -2)
	val := intoLuaValue(t.vm, -1) // value (-1)

	return val
}

func (t *LuaTable) Push(value LuaValue) {
	state := t.vm.state()

	t.deref(t.vm)     // table (-2)
	value.deref(t.vm) // value (-1)

	// Pop the table and key off
	defer ffi.Pop(state, 2)

	// Insert new index and set it to the value
	len := ffi.ObjLen(state, -2)
	ffi.PushInteger(state, ffi.LuaInteger(len+1))
	ffi.Insert(state, -2)
	ffi.SetTable(state, -3)
}

func (t *LuaTable) Pop() LuaValue {
	state := t.vm.state()

	t.deref(t.vm) // table (-1)

	// Pop the table off
	defer ffi.Pop(state, 1)

	// Get the last value and nil it out
	len := ffi.ObjLen(state, -1)
	ffi.PushInteger(state, ffi.LuaInteger(len)) // key   (-1), table (-2)
	ffi.GetTable(state, -2)                     // value (-1), table (-2)
	val := intoLuaValue(t.vm, -1)

	ffi.PushInteger(state, ffi.LuaInteger(len)) // key   (-1), value (-2), table (-3)
	ffi.PushNil(state)                          // nil   (-1), key   (-2), value (-3), table (-4)
	ffi.SetTable(state, -4)

	return val
}

func (t *LuaTable) RawPush(value LuaValue) {
	state := t.vm.state()

	t.deref(t.vm)     // table (-2)
	value.deref(t.vm) // value (-1)

	// Pop the table off
	defer ffi.Pop(state, 1)

	// Insert new index and set it to the value
	len := ffi.ObjLen(state, -2)
	ffi.PushInteger(state, ffi.LuaInteger(len+1)) // key   (-1), value (-2), table (-3)
	ffi.Insert(state, -2)                         // value (-1), key   (-2), table (-3)
	ffi.RawSet(state, -3)
}

func (t *LuaTable) RawPop() LuaValue {
	state := t.vm.state()

	t.deref(t.vm) // table (-1)

	// Pop the table off
	defer ffi.Pop(state, 1)

	// Get the last value and nil it out
	len := ffi.ObjLen(state, -1)
	ffi.PushInteger(state, ffi.LuaInteger(len)) // key   (-1), table (-2)
	ffi.RawGet(state, -2)                       // value (-1), table (-2)
	val := intoLuaValue(t.vm, -1)

	ffi.PushInteger(state, ffi.LuaInteger(len)) // key   (-1), value (-2), table (-3)
	ffi.PushNil(state)                          // nil   (-1), key   (-2), value (-3), table (-4)
	ffi.RawSet(state, -4)

	return val
}

func (t *LuaTable) Equals(other LuaValue) bool {
	state := t.vm.state()

	// Compare by reference first
	otherTable, ok := other.(*LuaTable)
	if !ok {
		return false
	}
	if t.index == otherTable.index {
		return true
	}

	// Compare by value
	t.deref(t.vm)          // table1 (-2)
	otherTable.deref(t.vm) // table2 (-1)

	// Pop off both tables
	defer ffi.Pop(state, 2)

	return ffi.Equal(state, -1, -2)
}

func (t *LuaTable) Clear() {
	state := t.vm.state()

	t.deref(t.vm) // table (-1)

	defer ffi.Pop(state, 1)

	// Iterate and nil out all keys
	ffi.PushNil(state)
	for ffi.Next(state, -2) != 0 {
		ffi.Pop(state, 1)
		ffi.PushValue(state, -1)
		ffi.PushNil(state)
		ffi.SetTable(state, -4)
	}
}

func (t *LuaTable) Len() int {
	state := t.vm.state()

	t.deref(t.vm)
	defer ffi.Pop(state, 1)

	return int(ffi.ObjLen(state, -1))
}

func (t *LuaTable) IsEmpty() bool { return t.Len() == 0 }

func (t *LuaTable) Iterable() map[LuaValue]LuaValue {
	state := t.vm.state()

	t.deref(t.vm)
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

func (t *LuaTable) SetMetatable(metatable *LuaTable) {
	state := t.vm.state()

	t.deref(t.vm)         // table (-2)
	metatable.deref(t.vm) // metatable (-1)

	// Set the metatable for the table
	ffi.SetMetatable(state, -2)

	// Pop metatable, re-ref the table
	ffi.Pop(state, 1)
	t.index = int(ffi.Ref(state, -1))
}

func (t *LuaTable) GetMetatable() *LuaTable {
	state := t.vm.state()

	if ok := ffi.GetMetatable(state, int32(t.index)); ok {
		index := ffi.Ref(state, -1)
		ffi.Pop(state, 1)

		return &LuaTable{vm: t.vm, index: int(index)}
	}

	return nil
}

//
// LuaValue implementation
//

var _ LuaValue = (*LuaTable)(nil)

func (t *LuaTable) lua() *Lua { return t.vm }
func (t *LuaTable) ref() int  { return t.index }

func (t *LuaTable) deref(lua *Lua) int {
	return int(ffi.GetRef(lua.state(), int32(t.ref())))
}
