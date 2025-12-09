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
	state := (*ffi.LuaState)(lua)
	entry := cgo.Handle(handle).Value().(*functionEntry)

	fn, ok := entry.registry.get(entry.id)
	if !ok {
		ffi.PushString(state, "function not found in registry")
		ffi.Error(state)
		return 0
	}

	return C.int(fn(state))
}

//export registryTrampolineDtorImpl
func registryTrampolineDtorImpl(_ *C.lua_State, handle C.uintptr_t) {
	entry := cgo.Handle(handle).Value().(*functionEntry)
	delete(entry.registry.functions, entry.id)
	cgo.Handle(handle).Delete()
}

type functionRegistry struct {
	functions map[uintptr]ffi.LuaCFunction
	nextID    uintptr
}

type functionEntry struct {
	registry *functionRegistry
	id       uintptr
}

func newFunctionRegistry() *functionRegistry {
	return &functionRegistry{
		functions: make(map[uintptr]ffi.LuaCFunction),
	}
}

func (fr *functionRegistry) register(fn ffi.LuaCFunction) *functionEntry {
	fr.nextID++
	id := fr.nextID
	fr.functions[id] = fn
	return &functionEntry{registry: fr, id: id}
}

func (fr *functionRegistry) get(id uintptr) (ffi.LuaCFunction, bool) {
	fn, ok := fr.functions[id]
	return fn, ok
}
