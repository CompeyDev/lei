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
*/
import "C"

var registryTrampoline = C.registryTrampoline

//export registryTrampolineImpl
func registryTrampolineImpl(lua *C.lua_State, registryPtr uintptr, funcID uintptr) C.int {
	handle := cgo.Handle(registryPtr)
	reg := handle.Value().(*functionRegistry)
	state := (*ffi.LuaState)(lua)

	fn, ok := reg.get(funcID)
	if !ok {
		ffi.PushString(state, "function not found in registry")
		ffi.Error(state)
		return 0
	}

	return C.int(fn(state))
}

type functionRegistry struct {
	functions map[uintptr]ffi.LuaCFunction
	nextID    uintptr
}

func newFunctionRegistry() *functionRegistry {
	return &functionRegistry{
		functions: make(map[uintptr]ffi.LuaCFunction),
	}
}

func (fr *functionRegistry) register(fn ffi.LuaCFunction) uintptr {
	fr.nextID++
	id := fr.nextID
	fr.functions[id] = fn
	return id
}

func (fr *functionRegistry) get(id uintptr) (ffi.LuaCFunction, bool) {
	fn, ok := fr.functions[id]
	return fn, ok
}

// FIXME: there is a memory leak of function entries here; we need to unregister
// once they are no longer used by Go or Lua. the issue here is that we cannot know
// when Lua is done with the function, so we need some kind of finalizer on the
// Lua side to notify us when it's done. the typical solution would be to use a full
// userdata instead of lightuserdata and set a dtor for that which calls back into Go
// to unregister
