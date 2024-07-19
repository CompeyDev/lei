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

const (
	LUA_OK = iota + 1
	LUA_YIELD
	LUA_ERRRUN
	LUA_ERRSYNTAX
	LUA_ERRMEM
	LUA_ERRERR
	LUA_BREAK
)

const (
	LUA_CORUN = iota + 1
	LUA_COSUS
	LUA_CONOR
	LUA_COFIN
	LUA_COERR
)

const (
	LUA_TNIL = iota
	LUA_TBOOLEAN

	LUA_TLIGHTUSERDATA
	LUA_TNUMBER
	LUA_TVECTOR

	LUA_TSTRING

	LUA_TTABLE
	LUA_TFUNCTION
	LUA_TUSERDATA
	LUA_THREAD
	LUA_TBUFFER

	LUA_TPROTO
	LUA_TUPVAL
	LUA_TDEADKEY

	LUA_T_COUNT = LUA_TPROTO
)

type LuaNumber float64
type LuaInteger int32
type LuaUnsigned int32

type LuaState = C.lua_State
type LuaCFunction func(L *LuaState) int32
type LuaContinuation func(L *LuaState, status int32) int32

type LuaUDestructor = func(*C.void)
type LuaDestructor = func(L *LuaState, _ unsafe.Pointer)

type LuaAlloc = func(ud, ptr unsafe.Pointer, osize, nsize C.size_t) *C.void

//
// ==================
// 	    VM state
// ==================
//

func NewState(f LuaAlloc, ud unsafe.Pointer) *LuaState {
	cf := C.malloc(C.size_t(unsafe.Sizeof(f)))
	defer C.free(cf)
	*(*LuaAlloc)(cf) = f

	cud := C.malloc(C.size_t(unsafe.Sizeof(ud)))
	defer C.free(cud)
	cud = ud

	return C.clua_newstate(cf, cud)
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

//
// =======================
// 	 Stack Manipulation
// =======================
//

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

	ccont := C.malloc(C.size_t(unsafe.Sizeof(cont)))
	defer C.free(ccont)
	if cont == nil {
		ccont = C.NULL
	} else {
		*(*LuaContinuation)(ccont) = cont
	}

	cf := C.malloc(C.size_t(unsafe.Sizeof(f)))
	defer C.free(cf)
	*(*LuaCFunction)(cf) = f

	C.clua_pushcclosurek(L, cf, cdebugname, C.int(nup), ccont)
}

func PushBoolean(L *LuaState, b bool) {
	cb := C.int(0)
	if b {
		cb = C.int(1)
	}

	C.lua_pushboolean(L, cb)
}

func PushThread(L *LuaState) bool {
	return C.lua_pushthread(L) != 0
}

func PushLightUserdataTagged(L *LuaState, p unsafe.Pointer, tag int32) {
	C.lua_pushlightuserdatatagged(L, p, C.int(tag))
}

func NewUserdataTagged(L *LuaState, sz uint64, tag int32) unsafe.Pointer {
	return C.lua_newuserdatatagged(L, C.size_t(sz), C.int(tag))
}

func NewUserdataDtor(L *LuaState, sz uint64, dtor LuaUDestructor) unsafe.Pointer {
	cdtor := C.malloc(C.size_t(unsafe.Sizeof(dtor)))
	defer C.free(cdtor)
	*(*LuaUDestructor)(cdtor) = dtor

	return C.clua_newuserdatadtor(L, C.size_t(sz), cdtor)
}

func NewBuffer(L *LuaState, sz uint64) unsafe.Pointer {
	return C.lua_newbuffer(L, C.size_t(sz))
}

//
// =====================
// 	   Get Functions
// =====================
//

func GetTable(L *LuaState, idx int32) int32 {
	return int32(C.lua_gettable(L, C.int(idx)))
}

func GetField(L *LuaState, idx int32, k string) int32 {
	ck := C.CString(k)
	defer C.free(unsafe.Pointer(ck))

	return int32(C.lua_getfield(L, C.int(idx), ck))
}

func RawGetField(L *LuaState, idx int32, k string) int32 {
	ck := C.CString(k)
	defer C.free(unsafe.Pointer(ck))

	return int32(C.lua_rawgetfield(L, C.int(idx), ck))
}

func RawGet(L *LuaState, idx int32) int32 {
	return int32(C.lua_rawget(L, C.int(idx)))
}

func RawGetI(L *LuaState, idx int32, n int32) int32 {
	return int32(C.lua_rawgeti(L, C.int(idx), C.int(n)))
}

func CreateTable(L *LuaState, narr int32, nrec int32) {
	C.lua_createtable(L, C.int(narr), C.int(nrec))
}

func SetReadonly(L *LuaState, idx int32, enabled bool) {
	cenabled := C.int(0)
	if enabled {
		cenabled = C.int(1)
	}

	C.lua_setreadonly(L, C.int(idx), cenabled)
}

func GetReadonly(L *LuaState, idx int32) bool {
	return C.lua_getreadonly(L, C.int(idx)) != 0
}

func SetSafeEnv(L *LuaState, idx int32, enabled bool) {
	cenabled := C.int(0)
	if enabled {
		cenabled = C.int(1)
	}

	C.lua_setsafeenv(L, C.int(idx), cenabled)
}

func Getfenv(L *LuaState, idx int32) {
	C.lua_getfenv(L, C.int(idx))
}

//
// ===================
//    Set Functions
// ===================
//

func SetTable(L *LuaState, idx int32) {
	C.lua_settable(L, C.int(idx))
}

func SetField(L *LuaState, idx int32, k string) {
	ck := C.CString(k)
	defer C.free(unsafe.Pointer(ck))

	C.lua_setfield(L, C.int(idx), ck)
}

func RawSetI(L *LuaState, idx int32, n int32) {
	C.lua_rawseti(L, C.int(idx), C.int(n))
}

func SetMetatable(L *LuaState, objindex int32) int32 {
	return int32(C.lua_setmetatable(L, C.int(objindex)))
}

func Setfenv(L *LuaState, idx int32) int32 {
	return int32(C.lua_setfenv(L, C.int(idx)))
}

//
// =========================
//    Bytecode Functions
// =========================
//

func LuauLoad(L *LuaState, chunkname string, data string, size uint64, env int32) int32 {
	cchunkname := C.CString(chunkname)
	defer C.free(unsafe.Pointer(cchunkname))

	cdata := C.CString(data)
	defer C.free(unsafe.Pointer(cdata))

	return int32(C.luau_load(L, cchunkname, cdata, C.size_t(size), C.int(env)))
}

func LuaCall(L *LuaState, nargs int32, nresults int32) {
	C.lua_call(L, C.int(nargs), C.int(nresults))
}

func LuaPcall(L *LuaState, nargs int32, nresults int32, errfunc int32) int32 {
	return int32(C.lua_pcall(L, C.int(nargs), C.int(nresults), C.int(errfunc)))
}

//
// ========================
//   Coroutine Functions
// ========================
//

func LuaYield(L *LuaState, nresults int32) int32 {
	return int32(C.lua_yield(L, C.int(nresults)))
}

func LuaBreak(L *LuaState) int32 {
	return int32(C.lua_break(L))
}

func LuaResume(L *LuaState, from *LuaState, nargs int32) int32 {
	return int32(C.lua_resume(L, from, C.int(nargs)))
}

func LuaResumeError(L *LuaState, from *LuaState) int32 {
	return int32(C.lua_resumeerror(L, from))
}

func LuaStatus(L *LuaState) int32 {
	return int32(C.lua_status(L))
}

func IsYieldable(L *LuaState) bool {
	return C.lua_isyieldable(L) != 0
}

func GetThreadData(L *LuaState) unsafe.Pointer {
	return C.lua_getthreaddata(L)
}

func SetThreadData(L *LuaState, data unsafe.Pointer) {
	C.lua_setthreaddata(L, data)
}

//
// ======================
//   Garbage Collection
// ======================
//

const (
	LUA_GCSTOP = iota
	LUA_GCRESTART

	LUA_GCCOLLECT

	LUA_GCCOUNT
	LUA_GCCOUNTB

	LUA_GCISRUNNING

	LUA_GCSTEP

	LUA_GCSETGOAL
	LUA_GCSETSTEPMUL
	LUA_GCSETSTEPSIZE
)

func LuaGc(L *LuaState, what int32, data int32) int32 {
	return int32(C.lua_gc(L, C.int(what), C.int(data)))
}

//
// ======================
//   Memory Statistics
// ======================
//

func SetMemCat(L *LuaState, category int32) {
	C.lua_setmemcat(L, C.int(category))
}

func TotalBytes(L *LuaState, category int32) uint64 {
	return uint64(C.lua_totalbytes(L, C.int(category)))
}

//
// ================================
//     Miscellaneous Functions
// ================================
//

func Error(L *LuaState) {
	C.lua_error(L)
}

func Next(L *LuaState, idx int32) int32 {
	return int32(C.lua_next(L, C.int(idx)))
}

func RawIter(L *LuaState, idx int32, iter int32) int32 {
	return int32(C.lua_rawiter(L, C.int(idx), C.int(iter)))
}

func Concat(L *LuaState, n int32) {
	C.lua_concat(L, C.int(n))
}

func Clock() float64 {
	return float64(C.lua_clock())
}

func SetUserdataTag(L *LuaState, idx int32, tag int32) {
	C.lua_setuserdatatag(L, C.int(idx), C.int(tag))
}

func SetUserdataDtor(L *LuaState, tag int32, dtor LuaDestructor) {
	cdtor := C.malloc(C.size_t(unsafe.Sizeof(dtor)))
	defer C.free(cdtor)

	if dtor == nil {
		cdtor = C.NULL
	} else {
		*(*LuaDestructor)(cdtor) = dtor
	}

	C.clua_setuserdatadtor(L, C.int(tag), cdtor)
}

func GetUserdataDtor(L *LuaState, tag int32) LuaDestructor {
	return *(*LuaDestructor)(unsafe.Pointer(C.lua_getuserdatadtor(L, C.int(tag))))
}

func SetUserdataMetatable(L *LuaState, tag int32, idx int32) {
	C.lua_setuserdatametatable(L, C.int(tag), C.int(idx))
}

func GetUserdataMetatable(L *LuaState, tag int32) {
	C.lua_getuserdatametatable(L, C.int(tag))
}

func SetLightUserdataName(L *LuaState, tag int32, name string) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	C.lua_setlightuserdataname(L, C.int(tag), cname)
}

func GetLightUserdataName(L *LuaState, tag int32) string {
	return C.GoString(C.lua_getlightuserdataname(L, C.int(tag)))
}

func CloneFunction(L *LuaState, idx int32) {
	C.lua_clonefunction(L, C.int(idx))
}

func ClearTable(L *LuaState, idx int32) {
	C.lua_cleartable(L, C.int(idx))
}

func GetAllocF(L *LuaState, ud *unsafe.Pointer) LuaAlloc {
	return *(*LuaAlloc)(unsafe.Pointer(C.lua_getallocf(L, ud)))
}

// TODO: Free udtor's after func
// TODO: Rest of it
