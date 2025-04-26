//go:build !LUAU_VECTOR4

package ffi

/*
#cgo CFLAGS: -Iluau/VM/include -I/usr/lib/gcc/x86_64-pc-linux-gnu/14.1.1/include
// #cgo LDFLAGS: -L${SRCDIR}/luau/cmake -lLuau.VM -lm -lstdc++
#include <lua.h>
*/
import "C"

func PushVector(L *LuaState, x float32, y float32, z float32) {
	C.lua_pushvector(L, C.float(x), C.float(y), C.float(z))
}
