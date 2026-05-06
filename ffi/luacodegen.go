package ffi

//go:generate go run ../build buildProject Luau.VM Luau.CodeGen

/*
#cgo CFLAGS: -Iluau/VM/include -Iluau/CodeGen/include
#cgo LDFLAGS: -L_obj -lLuau.VM -lLuau.CodeGen -lm -lstdc++
#include <stdlib.h>
#include <lua.h>
#include <luacodegen.h>
*/
import "C"

func LuauCodegenSupported() bool {
	return C.luau_codegen_supported() == 1
}

func LuauCodegenCreate(state *C.lua_State) {
	C.luau_codegen_create(state)
}

func LuauCodegenCompile(state *C.lua_State, idx int) {
	C.luau_codegen_compile(state, C.int(idx))
}
