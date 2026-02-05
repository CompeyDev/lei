package lua

//go:generate go tool cgo -- -I../ffi/luau/VM/include $GOFILE

/*
#cgo CFLAGS: -I../ffi/luau/VM/include
#cgo LDFLAGS: -L../ffi/_obj -lLuau.VM -lm -lstdc++

#include <lua.h>
#include <stdlib.h>
#include <stdint.h>

int indexMt(lua_State* L);
void methodMapDtorImpl(lua_State* L, uintptr_t);
void fieldMapDtorImpl(lua_State* L, uintptr_t);
*/
import "C"

import (
	"runtime/cgo"

	"github.com/CompeyDev/lei/ffi"
)

var indexMt = C.indexMt

var methodMapDtor = C.methodMapDtorImpl
var fieldMapDtor = C.fieldMapDtorImpl

type LuaUserData struct {
	vm    *Lua
	index int
	inner IntoUserData
}

func (ud *LuaUserData) Downcast() IntoUserData {
	if ud.inner != nil {
		return ud.inner
	}

	ud.deref(ud.vm)
	ptr := ffi.ToUserdata(ud.vm.state(), -1)

	if ptr != nil {
		return *(*IntoUserData)(ptr)
	} else {
		return nil
	}
}

//
// LuaValue implementation
//

var _ LuaValue = (*LuaUserData)(nil)

func (ud *LuaUserData) lua() *Lua { return ud.vm }
func (ud *LuaUserData) ref() int  { return ud.index }

func (ud *LuaUserData) deref(lua *Lua) int {
	return int(ffi.GetRef(lua.state(), int32(ud.ref())))
}

type IntoUserData interface {
	Methods(*MethodMap)
	MetaMethods(*MethodMap)
	Fields(*FieldMap)
}

type ValueRegistry[T any, U any] struct {
	inner       map[string]T
	transformer func(fn U) T
}

func (vr *ValueRegistry[T, U]) Insert(name string, value any) {
	if getter, ok := value.(T); ok {
		vr.inner[name] = getter
	} else {
		vr.inner[name] = vr.transformer(value.(U))
	}
}

type MethodMap = ValueRegistry[*functionEntry, GoFunction]

func newMethodMap(fnRegistry *functionRegistry) *MethodMap {
	return &MethodMap{
		inner:       make(map[string]*functionEntry),
		transformer: func(fn GoFunction) *functionEntry { return fnRegistry.register(fn) },
	}
}

type FieldGetter = func(*Lua) LuaValue
type FieldMap = ValueRegistry[FieldGetter, LuaValue]

func newFieldMap() *FieldMap {
	return &FieldMap{
		inner: make(map[string]FieldGetter),
		transformer: func(value LuaValue) FieldGetter {
			return func(*Lua) LuaValue { return value }
		},
	}
}

//export indexMtImpl
func indexMtImpl(lua *C.lua_State, fieldHandle, methodHandle uintptr, key *C.char) {
	rawState := (*ffi.LuaState)(lua)
	state := &Lua{
		// FIXME: what about the function registry?
		inner: &StateWithMemory{
			memState: getMemoryState(rawState),
			luaState: rawState,
		},
	}
	keyStr := C.GoString(key)

	// Field lookup
	fields := cgo.Handle(fieldHandle).Value().(*FieldMap)
	if getter := fields.inner[keyStr]; getter != nil {
		value := getter(state)
		value.deref(state)
		return
	}

	// Method lookup
	methods := cgo.Handle(methodHandle).Value().(*MethodMap)
	if method := methods.inner[keyStr]; method != nil {
		pushUpvalue(rawState, method, registryTrampolineDtor)
		ffi.PushCClosureK(rawState, registryTrampoline, nil, 1, nil)
		return
	}

	ffi.PushNil(rawState)
}

func valueRegistryDtorImpl[T any, U any](handle C.uintptr_t) {
	entry := cgo.Handle(handle).Value().(*ValueRegistry[T, U])
	clear(entry.inner)
	cgo.Handle(handle).Delete()
}

//export methodMapDtorImpl
func methodMapDtorImpl(_ *C.lua_State, handle C.uintptr_t) {
	valueRegistryDtorImpl[*functionRegistry, GoFunction](handle)
}

//export fieldMapDtorImpl
func fieldMapDtorImpl(_ *C.lua_State, handle C.uintptr_t) {
	valueRegistryDtorImpl[FieldGetter, LuaValue](handle)
}
