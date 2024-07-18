//go:build LUAU_VECTOR4

package internal

/*
#cgo CFLAGS: -Iluau/VM/include -I/usr/lib/gcc/x86_64-pc-linux-gnu/14.1.1/include -DLUA_VECTOR_SIZE=4
// #cgo LDFLAGS: -L${SRCDIR}/luau/cmake -lLuau.VM -lm -lstdc++
#include <lua.h>
*/
import "C"

func PushVector(L *C.lua_State, x float32, y float32, z float32, w float32) {
	C.lua_pushvector(L, C.float(x), C.float(y), C.float(z), C.float(w))
}
