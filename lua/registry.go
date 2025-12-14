package lua

import (
	"runtime/cgo"

	"github.com/CompeyDev/lei/ffi"
)

/*
#cgo CFLAGS: -I../ffi/luau/VM/include
#cgo LDFLAGS: -L../ffi/luau/cmake -lLuau.VM -lm -lstdc++

#include <lua.h>
#include <stdlib.h>
#include <stdint.h>

int registryTrampoline(lua_State* L);
void registryTrampolineDtor(lua_State* L);
*/
import "C"

var registryTrampoline = C.registryTrampoline
var registryTrampolineDtor = C.registryTrampolineDtor

//export registryTrampolineImpl
func registryTrampolineImpl(lua *C.lua_State, handle C.uintptr_t) C.int {
	rawState := (*ffi.LuaState)(lua)
	state := &Lua{
		inner: &StateWithMemory{
			memState: getMemoryState(rawState),
			luaState: rawState,
		},
	}

	entry := cgo.Handle(handle).Value().(*functionEntry)

	fn, ok := entry.registry.get(entry.id)
	if !ok {
		ffi.PushString(rawState, "function not found in registry")
		ffi.Error(rawState)
		return 0
	}

	argsCount := int(ffi.ToNumber(rawState, 1))
	args := make([]LuaValue, argsCount)

	for i := range argsCount {
		// Lua stack is 1-based, and the first argument is at index 2 (since index 1 is the count)
		stackIndex := int32(i + 2)
		args[i] = intoLuaValue(state, stackIndex)
	}

	returns, err := fn(state, args...)

	if err != nil {
		ffi.PushString(rawState, err.Error())
		ffi.Error(rawState)
		return 0
	}

	for _, ret := range returns {
		ret.deref(state)
	}

	return C.int(len(returns))
}

//export registryTrampolineDtorImpl
func registryTrampolineDtorImpl(_ *C.lua_State, handle C.uintptr_t) {
	entry := cgo.Handle(handle).Value().(*functionEntry)
	delete(entry.registry.functions, entry.id)
	cgo.Handle(handle).Delete()
}

type GoFunction func(lua *Lua, args ...LuaValue) ([]LuaValue, error)

type functionRegistry struct {
	functions map[uintptr]GoFunction
	nextID    uintptr
}

type functionEntry struct {
	registry *functionRegistry
	id       uintptr
}

func newFunctionRegistry() *functionRegistry {
	return &functionRegistry{
		functions: make(map[uintptr]GoFunction),
	}
}

func (fr *functionRegistry) register(fn GoFunction) *functionEntry {
	fr.nextID++
	id := fr.nextID
	fr.functions[id] = fn
	return &functionEntry{registry: fr, id: id}
}

func (fr *functionRegistry) get(id uintptr) (GoFunction, bool) {
	fn, ok := fr.functions[id]
	return fn, ok
}
