package lua

import (
	"fmt"
	"runtime/cgo"

	"github.com/CompeyDev/lei/ffi"
)

//go:generate go tool cgo $GOFILE

/*
#cgo CFLAGS: -I../ffi/luau/VM/include
#cgo LDFLAGS: -L../ffi/_obj -lLuau.VM -lm -lstdc++

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
func registryTrampolineImpl(lua *C.lua_State, handle uintptr) (C.int, *C.char) {
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
		return C.int(-1), C.CString("function not found in registry")
	}

	argsCount := int(ffi.ToNumber(rawState, 1))
	args := make([]LuaValue, argsCount)

	for i := range argsCount {
		// Lua stack is 1-based, and the first argument is at index 2 (since index 1 is the count)
		stackIndex := int32(i + 2)
		args[i] = intoLuaValue(state, stackIndex)
	}

	returns, callErr := fn(state, args...)

	// SAFETY: This must be caught elsewhere to avoid the longjmp
	if callErr != nil {
		return C.int(-1), C.CString(callErr.Error())
	}

	for _, ret := range returns {
		ret.deref(state)
	}

	return C.int(len(returns)), nil
}

//export registryTrampolineDtorImpl
func registryTrampolineDtorImpl(_ *C.lua_State, handle C.uintptr_t) {
	entry := cgo.Handle(handle).Value().(*functionEntry)
	delete(entry.registry.functions, entry.id)
	cgo.Handle(handle).Delete()
}

type GoFunction func(lua *Lua, args ...LuaValue) ([]LuaValue, error)

type functionRegistry struct {
	recoverPanics bool
	functions     map[uintptr]GoFunction
	nextID        uintptr
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

	if fr.recoverPanics {
		rawFn := fn
		fn = func(lua *Lua, args ...LuaValue) (result []LuaValue, err error) {
			defer func() {
				// Deferred panic handler
				if recv := recover(); recv != nil {
					switch v := recv.(type) {
					case error:
						err = v
					default:
						err = fmt.Errorf("go panic: %v", v)
					}
				}
			}()

			result, err = rawFn(lua, args...)

			return result, err
		}
	}

	return fn, ok
}
