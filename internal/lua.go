package internal

/*
#cgo CFLAGS: -Iluau/VM/include -I/usr/lib/gcc/x86_64-pc-linux-gnu/14.1.1/include
#cgo LDFLAGS: -L${SRCDIR}/luau/cmake -lLuau.VM
#include <lua.h>
#include <lualib.h>
#include <stdlib.h>
#include "clua.h"
*/
import "C"
import "unsafe"

type lua_Number float64
type lua_Integer int32
type lua_Unsigned int32

type lua_CFunction func(L *C.lua_State) int32
type lua_Continuation func(L *C.lua_State, status int32) int32

type lua_Udestructor = func(*C.void)
type lua_Destructor = func(L *C.lua_State, _ unsafe.Pointer)

// type lua_Alloc = func(ud, ptr *C.void, osize, nsize C.size_t) *C.void
type lua_Alloc = func(ptr unsafe.Pointer, osize, nsize uint64) unsafe.Pointer

//export go_allocf
func go_allocf(fp uintptr, ptr uintptr, osize uint64, nsize uint64) uintptr {
	p := ((*((*lua_Alloc)(unsafe.Pointer(fp))))(unsafe.Pointer(ptr), osize, nsize))
	return uintptr(p)
}

//
// ==================
// 	    VM state
// ==================
//

func NewState(ud unsafe.Pointer) *C.lua_State {
	return C.clua_newstate(unsafe.Pointer(ud))
}

func LuaClose(L *C.lua_State) {
	C.lua_close(L)
}

func NewThread(L *C.lua_State) *C.lua_State {
	return C.lua_newthread(L)
}

func MainThread(L *C.lua_State) *C.lua_State {
	return C.lua_mainthread(L)
}

func ResetThread(L *C.lua_State) {
	C.lua_resetthread(L)
}

func IsThreadReset(L *C.lua_State) bool {
	return C.lua_isthreadreset(L) != 0
}

//
// ==================
// 	    VM Stack
// ==================
//

func AbsIndex(L *C.lua_State, idx int32) int32 {
	return int32(C.lua_absindex(L, C.int(idx)))
}

func GetTop(L *C.lua_State) int32 {
	return int32(C.lua_gettop(L))
}

func SetTop(L *C.lua_State, idx int32) {
	C.lua_settop(L, C.int(idx))
}

func PushValue(L *C.lua_State, idx int32) {
	C.lua_pushvalue(L, C.int(idx))
}

func Remove(L *C.lua_State, idx int32) {
	C.lua_remove(L, C.int(idx))
}

func Insert(L *C.lua_State, idx int32) {
	C.lua_insert(L, C.int(idx))
}

func Replace(L *C.lua_State, idx int32) {
	C.lua_replace(L, C.int(idx))
}

func CheckStack(L *C.lua_State, sz int32) bool {
	return C.lua_checkstack(L, C.int(sz)) != 0
}

func RawCheckStack(L *C.lua_State, sz int32) {
	C.lua_rawcheckstack(L, C.int(sz))
}

func XMove(from, to *C.lua_State, n int32) {
	C.lua_xmove(from, to, C.int(n))
}

func XPush(from, to *C.lua_State, idx int32) {
	C.lua_xpush(from, to, C.int(idx))
}

//
// ======================
// 	    Stack Values
// ======================
//

func IsNumber(L *C.lua_State, idx int32) bool {
	return C.lua_isnumber(L, C.int(idx)) != 0
}

func IsString(L *C.lua_State, idx int32) bool {
	return C.lua_isstring(L, C.int(idx)) != 0
}

func IsCFunction(L *C.lua_State, idx int32) bool {
	return C.lua_iscfunction(L, C.int(idx)) != 0
}

func IsLFunction(L *C.lua_State, idx int32) bool {
	return C.lua_isLfunction(L, C.int(idx)) != 0
}

func IsUserData(L *C.lua_State, idx int32) bool {
	return C.lua_isuserdata(L, C.int(idx)) != 0
}

func Type(L *C.lua_State, idx int32) bool {
	return C.lua_type(L, C.int(idx)) != 0
}

func TypeName(L *C.lua_State, tp int32) string {
	return C.GoString(C.lua_typename(L, C.int(tp)))
}

func Equal(L *C.lua_State, idx1, idx2 int32) bool {
	return C.lua_equal(L, C.int(idx1), C.int(idx2)) != 0
}

func RawEqual(L *C.lua_State, idx1, idx2 int32) bool {
	return C.lua_rawequal(L, C.int(idx1), C.int(idx2)) != 0
}

func LessThan(L *C.lua_State, idx1, idx2 int32) bool {
	return C.lua_lessthan(L, C.int(idx1), C.int(idx2)) != 0
}

func ToNumberX(L *C.lua_State, idx int32, isnum bool) lua_Number {
	isnumInner := C.int(0)
	if isnum {
		isnumInner = C.int(1)
	}

	return lua_Number(C.lua_tonumberx(L, C.int(idx), &isnumInner))
}

func ToIntegerX(L *C.lua_State, idx int32, isnum bool) lua_Integer {
	isnumInner := C.int(0)
	if isnum {
		isnumInner = C.int(1)
	}

	return lua_Integer(C.lua_tointegerx(L, C.int(idx), &isnumInner))
}

func ToUnsignedX(L *C.lua_State, idx int32, isnum bool) lua_Unsigned {
	isnumInner := C.int(0)
	if isnum {
		isnumInner = C.int(1)
	}

	return lua_Unsigned(C.lua_tounsignedx(L, C.int(idx), &isnumInner))
}

// TODO: Rest of it
