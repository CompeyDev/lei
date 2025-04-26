package ffi

//#include <stdlib.h>
import "C"

func GetSubtable(L *LuaState, idx int32, fname string) bool {
	absIdx := AbsIndex(L, idx)
	if !CheckStack(L, 3+20) {
		panic("stack overflow")
	}

	PushString(L, fname)
	if GetTable(L, absIdx) == LUA_TTABLE {
		return true
	}

	Pop(L, 1)
	NewTable(L)
	PushString(L, fname)
	PushValue(L, -2)
	SetTable(L, absIdx)
	return false
}

func RequireLib(L *LuaState, modName string, openF LuaCFunction, isGlobal bool) {
	if !CheckStack(L, 3+20) {
		LErrorL(L, "stack overflow")
	}

	GetSubtable(L, LUA_REGISTRYINDEX, "_LOADED")
	if GetField(L, -1, modName) == LUA_TNIL {
		Pop(L, 1)
		PushCFunction(L, openF)
		PushString(L, modName)
		Call(L, 1, 1)
		PushValue(L, -1)
		SetField(L, -3, modName)
	}

	if isGlobal {
		PushNil(L)
		SetGlobal(L, modName)
	} else {
		PushValue(L, -1)
		SetGlobal(L, modName)
	}

	Replace(L, -2)
}
