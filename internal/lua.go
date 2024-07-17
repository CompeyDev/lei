package internal

/*
#cgo CFLAGS: -Iluau/VM/include -I/usr/lib/gcc/x86_64-pc-linux-gnu/14.1.1/include -I${SRCDIR}
#include <lua.h>
#include <lualib.h>
#include <stdlib.h>
#include <clua.h>
*/
import "C"
import "unsafe"

type lua_Number C.double
type lua_Integer C.int
type lua_Unsigned C.uint

type lua_CFunction func(L *C.lua_State) C.int
type lua_Continuation func(L *C.lua_State, status C.int) C.int

type lua_Udestructor = func(*C.void)
type lua_Destructor = func(L *C.lua_State, _ *C.void)

type lua_Alloc = func(ud, ptr *C.void, osize, nsize C.size_t) *C.void
type clua_Alloc = func(ptr unsafe.Pointer, osize, nsize uint) unsafe.Pointer

//export go_allocf
func go_allocf(fp uintptr, ptr uintptr, osize uint, nsize uint) uintptr {
	p := ((*((*clua_Alloc)(unsafe.Pointer(fp))))(unsafe.Pointer(ptr), osize, nsize))
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

func AbsIndex(L *C.lua_State, idx C.int) C.int {
	return C.lua_absindex(L, idx)
}

func GetTop(L *C.lua_State) C.int {
	return C.lua_gettop(L)
}

func SetTop(L *C.lua_State, idx C.int) {
	C.lua_settop(L, idx)
}

func PushValue(L *C.lua_State, idx C.int) {
	C.lua_pushvalue(L, idx)
}

func Remove(L *C.lua_State, idx C.int) {
	C.lua_remove(L, idx)
}

func Insert(L *C.lua_State, idx C.int) {
	C.lua_insert(L, idx)
}

func Replace(L *C.lua_State, idx C.int) {
	C.lua_replace(L, idx)
}

func CheckStack(L *C.lua_State, sz C.int) bool {
	return C.lua_checkstack(L, sz) != 0
}

func RawCheckStack(L *C.lua_State, sz C.int) {
	C.lua_rawcheckstack(L, sz)
}

func XMove(from, to *C.lua_State, n C.int) {
	C.lua_xmove(from, to, n)
}

func XPush(from, to *C.lua_State, idx C.int) {
	C.lua_xpush(from, to, idx)
}

//
// ======================
// 	    Stack Values
// ======================
//

func IsNumber(L *C.lua_State, idx C.int) bool {
	return C.lua_isnumber(L, idx) != 0
}

func IsString(L *C.lua_State, idx C.int) bool {
	return C.lua_isstring(L, idx) != 0
}

func IsCFunction(L *C.lua_State, idx C.int) bool {
	return C.lua_iscfunction(L, idx) != 0
}

func IsLFunction(L *C.lua_State, idx C.int) bool {
	return C.lua_isLfunction(L, idx) != 0
}

func IsUserData(L *C.lua_State, idx C.int) bool {
	return C.lua_isuserdata(L, idx) != 0
}

func Type(L *C.lua_State, idx C.int) bool {
	return C.lua_type(L, idx) != 0
}

func TypeName(L *C.lua_State, tp C.int) *C.char {
	return C.lua_typename(L, tp)
}

func Equal(L *C.lua_State, idx1, idx2 C.int) bool {
	return C.lua_equal(L, idx1, idx2) != 0
}

func RawEqual(L *C.lua_State, idx1, idx2 C.int) bool {
	return C.lua_rawequal(L, idx1, idx2) != 0
}

func LessThan(L *C.lua_State, idx1, idx2 C.int) bool {
	return C.lua_lessthan(L, idx1, idx2) != 0
}

func ToNumberX(L *C.lua_State, idx C.int, isnum bool) lua_Number {
	isnumInner := C.int(0)
	if isnum {
		isnumInner = C.int(1)
	}

	return lua_Number(C.lua_tonumberx(L, idx, &isnumInner))
}

func ToIntegerX(L *C.lua_State, idx C.int, isnum bool) lua_Integer {
	isnumInner := C.int(0)
	if isnum {
		isnumInner = C.int(1)
	}

	return lua_Integer(C.lua_tointegerx(L, idx, &isnumInner))
}

func ToUnsignedX(L *C.lua_State, idx C.int, isnum bool) lua_Unsigned {
	isnumInner := C.int(0)
	if isnum {
		isnumInner = C.int(1)
	}

	return lua_Unsigned(C.lua_tounsignedx(L, idx, &isnumInner))
}

// TODO: Rest of it
// TODO: Convert C.* types in args to go types
