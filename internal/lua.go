package internal

/*
#cgo CFLAGS: -Iluau/VM/include -I/usr/lib/gcc/x86_64-pc-linux-gnu/14.1.1/include
// #cgo LDFLAGS: -L${SRCDIR}/luau/cmake -lLuau.VM -lm -lstdc++
#include <lua.h>
#include <lualib.h>
#include <stdlib.h>
#include "clua.h"
*/
import "C"
import "unsafe"

type LuaNumber float64
type LuaInteger int32
type LuaUnsigned int32

type LuaState = C.lua_State
type LuaCFunction func(L *LuaState) int32
type LuaContinuation func(L *LuaState, status int32) int32

type LuaUDestructor = func(*C.void)
type LuaDestructor = func(L *LuaState, _ unsafe.Pointer)

// type LuaAlloc = func(ud, ptr *C.void, osize, nsize C.size_t) *C.void
type LuaAlloc = func(ptr unsafe.Pointer, osize, nsize uint64) unsafe.Pointer

//export go_allocf
func go_allocf(fp uintptr, ptr uintptr, osize uint64, nsize uint64) uintptr {
	p := ((*((*LuaAlloc)(unsafe.Pointer(fp))))(unsafe.Pointer(ptr), osize, nsize))
	return uintptr(p)
}

//
// ==================
// 	    VM state
// ==================
//

func NewState(ud unsafe.Pointer) *LuaState {
	return C.clua_newstate(unsafe.Pointer(ud))
}

func LuaClose(L *LuaState) {
	C.lua_close(L)
}

func NewThread(L *LuaState) *LuaState {
	return C.lua_newthread(L)
}

func MainThread(L *LuaState) *LuaState {
	return C.lua_mainthread(L)
}

func ResetThread(L *LuaState) {
	C.lua_resetthread(L)
}

func IsThreadReset(L *LuaState) bool {
	return C.lua_isthreadreset(L) != 0
}

//
// ==================
// 	    VM Stack
// ==================
//

func AbsIndex(L *LuaState, idx int32) int32 {
	return int32(C.lua_absindex(L, C.int(idx)))
}

func GetTop(L *LuaState) int32 {
	return int32(C.lua_gettop(L))
}

func SetTop(L *LuaState, idx int32) {
	C.lua_settop(L, C.int(idx))
}

func PushValue(L *LuaState, idx int32) {
	C.lua_pushvalue(L, C.int(idx))
}

func Remove(L *LuaState, idx int32) {
	C.lua_remove(L, C.int(idx))
}

func Insert(L *LuaState, idx int32) {
	C.lua_insert(L, C.int(idx))
}

func Replace(L *LuaState, idx int32) {
	C.lua_replace(L, C.int(idx))
}

func CheckStack(L *LuaState, sz int32) bool {
	return C.lua_checkstack(L, C.int(sz)) != 0
}

func RawCheckStack(L *LuaState, sz int32) {
	C.lua_rawcheckstack(L, C.int(sz))
}

func XMove(from, to *LuaState, n int32) {
	C.lua_xmove(from, to, C.int(n))
}

func XPush(from, to *LuaState, idx int32) {
	C.lua_xpush(from, to, C.int(idx))
}

//
// ======================
// 	    Stack Values
// ======================
//

func IsNumber(L *LuaState, idx int32) bool {
	return C.lua_isnumber(L, C.int(idx)) != 0
}

func IsString(L *LuaState, idx int32) bool {
	return C.lua_isstring(L, C.int(idx)) != 0
}

func IsCFunction(L *LuaState, idx int32) bool {
	return C.lua_iscfunction(L, C.int(idx)) != 0
}

func IsLFunction(L *LuaState, idx int32) bool {
	return C.lua_isLfunction(L, C.int(idx)) != 0
}

func IsUserData(L *LuaState, idx int32) bool {
	return C.lua_isuserdata(L, C.int(idx)) != 0
}

func Type(L *LuaState, idx int32) bool {
	return C.lua_type(L, C.int(idx)) != 0
}

func TypeName(L *LuaState, tp int32) string {
	return C.GoString(C.lua_typename(L, C.int(tp)))
}

func Equal(L *LuaState, idx1, idx2 int32) bool {
	return C.lua_equal(L, C.int(idx1), C.int(idx2)) != 0
}

func RawEqual(L *LuaState, idx1, idx2 int32) bool {
	return C.lua_rawequal(L, C.int(idx1), C.int(idx2)) != 0
}

func LessThan(L *LuaState, idx1, idx2 int32) bool {
	return C.lua_lessthan(L, C.int(idx1), C.int(idx2)) != 0
}

func ToNumberX(L *LuaState, idx int32, isnum bool) LuaNumber {
	isnumInner := C.int(0)
	if isnum {
		isnumInner = C.int(1)
	}

	return LuaNumber(C.lua_tonumberx(L, C.int(idx), &isnumInner))
}

func ToIntegerX(L *LuaState, idx int32, isnum bool) LuaInteger {
	isnumInner := C.int(0)
	if isnum {
		isnumInner = C.int(1)
	}

	return LuaInteger(C.lua_tointegerx(L, C.int(idx), &isnumInner))
}

func ToUnsignedX(L *LuaState, idx int32, isnum bool) LuaUnsigned {
	isnumInner := C.int(0)
	if isnum {
		isnumInner = C.int(1)
	}

	return LuaUnsigned(C.lua_tounsignedx(L, C.int(idx), &isnumInner))
}

func ToVector(L *LuaState, idx int32) {
	C.lua_tovector(L, C.int(idx))
}

func ToBoolean(L *LuaState, idx int32) bool {
	return C.lua_toboolean(L, C.int(idx)) != 0
}

func ToLString(L *LuaState, idx int32, len *uint64) string {
	return C.GoString(C.lua_tolstring(L, C.int(idx), (*C.size_t)(len)))
}

func ToStringAtom(L *LuaState, idx int32, atom *int32) string {
	return C.GoString(C.lua_tostringatom(L, C.int(idx), (*C.int)(atom)))
}

func NameCallAtom(L *LuaState, atom *int32) string {
	return C.GoString(C.lua_namecallatom(L, (*C.int)(atom)))
}

func ObjLen(L *LuaState, idx int32) uint64 {
	return uint64(C.lua_objlen(L, C.int(idx)))
}

func ToCFunction(L *LuaState, idx int32) LuaCFunction {
	p := unsafe.Pointer(C.lua_tocfunction(L, C.int(idx)))
	if p == C.NULL {
		return nil
	}

	return *(*LuaCFunction)(p)
}

func ToLightUserdata(L *LuaState, idx int32) unsafe.Pointer {
	return unsafe.Pointer(C.lua_tolightuserdata(L, C.int(idx)))
}

func ToLightUserdataTagged(L *LuaState, idx int32, tag int32) unsafe.Pointer {
	return unsafe.Pointer(C.lua_tolightuserdatatagged(L, C.int(idx), C.int(tag)))
}

func ToUserdata(L *LuaState, idx int32) unsafe.Pointer {
	return unsafe.Pointer(C.lua_touserdata(L, C.int(idx)))
}

func ToUserdataTagged(L *LuaState, idx int32, tag int32) unsafe.Pointer {
	return unsafe.Pointer(C.lua_touserdatatagged(L, C.int(idx), C.int(tag)))
}

func UserdataTag(L *LuaState, idx int32) int32 {
	return int32(C.lua_userdatatag(L, C.int(idx)))
}

func LightUserdataTag(L *LuaState, idx int32) int32 {
	return int32(C.lua_lightuserdatatag(L, C.int(idx)))
}

func ToThread(L *LuaState, idx int32) *LuaState {
	return C.lua_tothread(L, C.int(idx))
}

func ToBuffer(L *LuaState, idx int32, len *uint64) unsafe.Pointer {
	return unsafe.Pointer(C.lua_tobuffer(L, C.int(idx), (*C.size_t)(len)))
}

func ToPointer(L *LuaState, idx int32) unsafe.Pointer {
	return unsafe.Pointer(C.lua_topointer(L, C.int(idx)))
}

// =======================
// 	 Stack Manipulation
// =======================

func PushNil(L *LuaState) {
	C.lua_pushnil(L)
}

func PushNumber(L *LuaState, n LuaNumber) {
	C.lua_pushnumber(L, C.lua_Number(n))
}

func PushInteger(L *LuaState, n LuaInteger) {
	C.lua_pushinteger(L, C.lua_Integer(n))
}

func PushUnsigned(L *LuaState, n LuaUnsigned) {
	C.lua_pushunsigned(L, C.lua_Unsigned(n))
}

func PushLString(L *LuaState, s string, l uint64) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))

	C.lua_pushlstring(L, cs, C.size_t(l))
}

func PushString(L *LuaState, s string) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))

	C.lua_pushstring(L, cs)
}

// NOTE: We can't have lua_pushfstringL, since varadic
// arguments from Go->C isn't something that is possible.
// func PushFStringL(L *lua_State, fmt string) {}

func PushCClosureK(L *LuaState, f LuaCFunction, debugname string, nup int32, cont LuaContinuation) {
	cdebugname := C.CString(debugname)
	defer C.free(unsafe.Pointer(cdebugname))

	var ccont unsafe.Pointer
	if cont == nil {
		ccont = C.NULL
	} else {
		ccont = unsafe.Pointer(&cont)
	}

	C.clua_pushcclosurek(L, unsafe.Pointer(&f), cdebugname, C.int(nup), ccont)
}

// TODO: Rest of it
