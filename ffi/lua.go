package ffi

/*
#cgo CFLAGS: -Iluau/VM/include
#cgo LDFLAGS: -Lluau/cmake -lLuau.VM -lm -lstdc++
#include <lua.h>
#include <lualib.h>
#include <stdlib.h>
#include "clua.h"
*/
import "C"
import (
	"errors"
	"strconv"
	"unsafe"
)

//
// ===========================
//     luaconf.h Constants
// ===========================
//

const (
	LUAI_MAXCSTACK  = 8000 // Max allowed values in lua stack
	LUA_UTAG_LIMIT  = 128  // Max number of lua userdata tags
	LUA_LUTAG_LIMIT = 128  // Max number of light lua userdata tags
)

const LUA_MULTRET = -1

//
// ====================
//    Pseudo Indices
// ====================
//

const (
	LUA_REGISTRYINDEX = C.LUA_REGISTRYINDEX
	LUA_ENVIRONINDEX  = C.LUA_ENVIRONINDEX
	LUA_GLOBALSINDEX  = C.LUA_GLOBALSINDEX
)

//
// ======================
//     Thread Status
// ======================
//

const (
	LUA_OK = iota
	LUA_YIELD
	LUA_ERRRUN
	LUA_ERRSYNTAX
	LUA_ERRMEM
	LUA_ERRERR
	LUA_BREAK
)

//
// =========================
//     Coroutine Status
// =========================
//

const (
	LUA_CORUN = iota
	LUA_COSUS
	LUA_CONOR
	LUA_COFIN
	LUA_COERR
)

//
// ===================
//     Basic Types
// ===================
//

const LUA_TNONE = -1
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
	LUA_TTHREAD
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

func NewState(f unsafe.Pointer, ud unsafe.Pointer) *LuaState {
	// cf := C.malloc(C.size_t(unsafe.Sizeof(f)))
	// defer C.free(cf)
	// *(*LuaAlloc)(cf) = f

	// cud := C.malloc(C.size_t(unsafe.Sizeof(ud)))
	// defer C.free(cud)
	// cud = ud

	// ^^^^^^^^^^^^^^^^
	// we no longer support a go pointer for safety, consumers should instead use
	// cgo generated trampolines

	return C.clua_newstate(f, ud)
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

func Type(L *LuaState, idx int32) int32 {
	return int32(C.lua_type(L, C.int(idx)))
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

func ToNumberX(L *LuaState, idx int32, isnum *bool) LuaNumber {
	cisnumber := C.int(0)
	if *isnum {
		cisnumber = C.int(1)
	}

	num := LuaNumber(C.lua_tonumberx(L, C.int(idx), &cisnumber))
	*isnum = cisnumber != C.int(0)

	return num
}

func ToIntegerX(L *LuaState, idx int32, isnum *bool) LuaInteger {
	cisnumber := C.int(0)
	if *isnum {
		cisnumber = C.int(1)
	}

	integer := LuaInteger(C.lua_tointegerx(L, C.int(idx), &cisnumber))
	*isnum = cisnumber != C.int(0)

	return integer
}

func ToUnsignedX(L *LuaState, idx int32, isnum *bool) LuaUnsigned {
	cisnumber := C.int(0)
	if *isnum {
		cisnumber = C.int(1)
	}

	unsigned := LuaUnsigned(C.lua_tounsignedx(L, C.int(idx), &cisnumber))
	*isnum = cisnumber != C.int(0)

	return unsigned
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

	// NOTE: CStrings are null-terminated, and hence one longer than Go strings
	C.lua_pushlstring(L, cs, C.size_t(l+1))
}

func PushString(L *LuaState, s string) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))

	C.lua_pushstring(L, cs)
}

// NOTE: We can't have lua_pushfstringL, since varadic
// arguments from Go->C isn't something that is possible.
// func PushFStringL(L *lua_State, fmt string) {}

func PushCClosureK(L *LuaState, f unsafe.Pointer, debugname *string, nup int32, cont unsafe.Pointer) {
	var cdebugname *C.char
	if debugname != nil && *debugname != "" {
		cdebugname = C.CString(*debugname)
		defer C.free(unsafe.Pointer(cdebugname))
	}

	C.clua_pushcclosurek(L, f, cdebugname, C.int(nup), cont)
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

func NewUserdataDtor(L *LuaState, sz uint64, dtor unsafe.Pointer) unsafe.Pointer {
	return C.clua_newuserdatadtor(L, C.size_t(sz), dtor)
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

func GetMetatable(L *LuaState, objindex int32) int32 {
	return int32(C.lua_getmetatable(L, C.int(objindex)))
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

func RawSet(L *LuaState, idx int32) {
	C.lua_rawset(L, C.int(idx))
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

func LuauLoad(L *LuaState, chunkname string, data []byte, size uint64, env int32) bool {
	cchunkname := C.CString(chunkname)
	defer C.free(unsafe.Pointer(cchunkname))

	var cdata *C.char
	if size == 0 {
		// NULL for empty slices
		cdata = (*C.char)(C.NULL)
	} else {
		cdata = (*C.char)(unsafe.Pointer(&data[0]))
	}

	// NOTE: We don't free the bytecode after it's loaded

	return C.luau_load(L, cchunkname, cdata, C.size_t(size), C.int(env)) == 0
}

func Call(L *LuaState, nargs int32, nresults int32) {
	C.lua_call(L, C.int(nargs), C.int(nresults))
}

func Pcall(L *LuaState, nargs int32, nresults int32, errfunc int32) int32 {
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

func SetUserdataDtor(L *LuaState, tag int32, dtor unsafe.Pointer) {
	C.clua_setuserdatadtor(L, C.int(tag), dtor)
}

func GetUserdataDtor(L *LuaState, tag int32) LuaDestructor {
	return *(*LuaDestructor)(unsafe.Pointer(C.lua_getuserdatadtor(L, C.int(tag))))
}

func SetUserdataMetatable(L *LuaState, tag int32, idx int32) {
	C.lua_setuserdatametatable(L, C.int(idx))
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

func GetAllocF(L *LuaState, ud *unsafe.Pointer) unsafe.Pointer {
	// SAFETY: we cannot call this as a Go function, must be treated as an opaque
	// C pointer with unsafe.Pointer, previously this used to be casted
	return unsafe.Pointer(C.lua_getallocf(L, ud))
}

//
// ==========================
//      Reference System
// ==========================
//

const (
	LUA_NOREF = iota - 1
	LUA_REFNIL
)

func Ref(L *LuaState, idx int32) int32 {
	return int32(C.lua_ref(L, C.int(idx)))
}

func Unref(L *LuaState, ref int32) {
	C.lua_unref(L, C.int(ref))
}

//
// ==================
//     Debug API
// ==================
//

const LUA_IDSIZE = 256

type LuaDebug struct {
	Name        string
	What        string
	Source      string
	ShortSrc    string
	LineDefined int8
	CurrentLine int8
	NUpVals     uint8
	NParams     uint8
	IsVarArg    int8
	Userdata    unsafe.Pointer
	SSbuf       string // size = LUA_IDSIZE
}
type LuaHook = func(L *LuaState, ar *LuaDebug)
type LuaCoverage = func(
	context unsafe.Pointer,
	function string,
	linedefined int32,
	depth int32,
	hits *int32,
	size uint64,
)

func StackDepth(L *LuaState) int32 {
	return int32(C.lua_stackdepth(L))
}

func (ar *LuaDebug) toCLuaDebug() (*C.lua_Debug, error) {
	cname := C.CString(ar.Name)
	defer C.free(unsafe.Pointer(cname))
	carwhat := C.CString(ar.What)
	defer C.free(unsafe.Pointer(carwhat))
	csource := C.CString(ar.Source)
	defer C.free(unsafe.Pointer(csource))
	cshortsrc := C.CString(ar.ShortSrc)
	defer C.free(unsafe.Pointer(cshortsrc))
	if len(ar.SSbuf)+1 != LUA_IDSIZE { // contains null delimeter, so LUA_IDSIZE is one greater than the string len
		return nil, errors.New("lua.GetInfo: SSbuf must be exactly " + strconv.Itoa(LUA_IDSIZE) + " bytes long")
	}
	cssbuf := C.CString(ar.SSbuf)
	defer C.free(unsafe.Pointer(cssbuf))

	return &C.lua_Debug{
		name:        cname,
		what:        carwhat,
		source:      csource,
		short_src:   cshortsrc,
		linedefined: C.int(ar.LineDefined),
		currentline: C.int(ar.CurrentLine),
		nupvals:     C.uchar(ar.NUpVals),
		nparams:     C.uchar(ar.NParams),
		isvararg:    C.char(ar.IsVarArg),
		userdata:    ar.Userdata,
		ssbuf:       *(*[LUA_IDSIZE]C.char)(unsafe.Pointer(cssbuf)),
	}, nil
}

// Errors if invalid LuaDebug provided, returns -1
func GetInfo(L *LuaState, level int32, what string, ar *LuaDebug) (int32, error) {
	cwhat := C.CString(what)
	defer C.free(unsafe.Pointer(cwhat))

	car, err := ar.toCLuaDebug()
	if err != nil {
		return -1, err
	}

	carp := C.malloc(C.size_t(unsafe.Sizeof(*ar)))
	defer C.free(carp)
	*(**C.lua_Debug)(carp) = car

	return int32(C.lua_getinfo(L, C.int(level), cwhat, (*C.lua_Debug)(carp))), nil
}

func GetArgument(L *LuaState, level int32, n int32) int32 {
	return int32(C.lua_getargument(L, C.int(level), C.int(n)))
}

func GetLocal(L *LuaState, level int32, n int32) string {
	return C.GoString(C.lua_getlocal(L, C.int(level), C.int(n)))
}

func SetLocal(L *LuaState, level int32, n int32) string {
	return C.GoString(C.lua_setlocal(L, C.int(level), C.int(n)))
}

func GetUpvalue(L *LuaState, funcindex int32, n int32) string {
	return C.GoString(C.lua_getupvalue(L, C.int(funcindex), C.int(n)))
}

func SetUpvalue(L *LuaState, funcindex int32, n int32) string {
	return C.GoString(C.lua_setupvalue(L, C.int(funcindex), C.int(n)))
}

func SingleStep(L *LuaState, enabled bool) {
	cenabled := C.int(0)
	if enabled {
		cenabled = C.int(1)
	}

	C.lua_singlestep(L, cenabled)
}

func Breakpoint(L *LuaState, funcindex int32, line int32, enabled bool) int32 {
	cenabled := C.int(0)
	if enabled {
		cenabled = C.int(1)
	}

	return int32(C.lua_breakpoint(L, C.int(funcindex), C.int(line), cenabled))
}

func GetCoverage(L *LuaState, funcindex int32, context unsafe.Pointer, callback LuaCoverage) {
	ccallback := C.malloc(C.size_t(unsafe.Sizeof(callback)))
	defer C.free(ccallback)
	*(*LuaCoverage)(ccallback) = callback

	C.clua_getcoverage(L, C.int(funcindex), context, ccallback)
}

func DebugTrace(L *LuaState) string {
	return C.GoString(C.lua_debugtrace(L))
}

type LuaCallbacks struct {
	Userdata            unsafe.Pointer
	Interrupt           func(L *LuaState, gc int32)
	Panic               func(L *LuaState, errcode int32)
	UserThread          func(LP *LuaState, L *LuaState)
	UserAtom            func(s string, l uint64) int16
	DebugBreak          func(L *LuaState, ar *LuaDebug)
	DebugStep           func(L *LuaState, ar *LuaDebug)
	DebugInterrupt      func(L *LuaState, ar *LuaDebug)
	DebugProtectedError func(L *LuaState)
	OnAllocate          func(L *LuaState, osize uint64, nsize uint64)
}

func Callbacks(L *LuaState) *LuaCallbacks {
	ccallbacks := C.lua_callbacks(L)

	return &LuaCallbacks{
		Userdata:            ccallbacks.userdata,
		Interrupt:           *(*func(L *LuaState, gc int32))(unsafe.Pointer(ccallbacks.interrupt)),
		Panic:               *(*func(L *LuaState, errcode int32))(unsafe.Pointer(ccallbacks.panic)),
		UserThread:          *(*func(LP *LuaState, L *LuaState))(unsafe.Pointer(ccallbacks.userthread)),
		UserAtom:            *(*func(s string, l uint64) int16)(unsafe.Pointer(ccallbacks.useratom)),
		DebugBreak:          *(*func(L *LuaState, ar *LuaDebug))(unsafe.Pointer(ccallbacks.debugbreak)),
		DebugStep:           *(*func(L *LuaState, ar *LuaDebug))(unsafe.Pointer(ccallbacks.debugstep)),
		DebugInterrupt:      *(*func(L *LuaState, ar *LuaDebug))(unsafe.Pointer(ccallbacks.debuginterrupt)),
		DebugProtectedError: *(*func(L *LuaState))(unsafe.Pointer(ccallbacks.debugprotectederror)),
		OnAllocate:          *(*func(L *LuaState, osize uint64, nsize uint64))(unsafe.Pointer(ccallbacks.onallocate)),
	}
}

func ToNumber(L *LuaState, i int32) LuaNumber {
	return ToNumberX(L, i, new(bool))
}

func ToInteger(L *LuaState, i int32) LuaInteger {
	return ToIntegerX(L, i, new(bool))
}

func ToUnsigned(L *LuaState, i int32) LuaUnsigned {
	return ToUnsignedX(L, i, new(bool))
}

func Pop(L *LuaState, n int32) {
	SetTop(L, -n-1)
}

func NewTable(L *LuaState) {
	CreateTable(L, 0, 0)
}

func NewUserdata(L *LuaState, sz uint64) unsafe.Pointer {
	return NewUserdataTagged(L, sz, 0)
}

func IsFunction(L *LuaState, n int32) bool {
	return Type(L, n) == LUA_TFUNCTION
}

func IsTable(L *LuaState, n int32) bool {
	return Type(L, n) == LUA_TTABLE
}

func IsLightUserdata(L *LuaState, n int32) bool {
	return Type(L, n) == LUA_TLIGHTUSERDATA
}

func IsNil(L *LuaState, n int32) bool {
	return Type(L, n) == LUA_TNIL
}

func IsBoolean(L *LuaState, n int32) bool {
	return Type(L, n) == LUA_TBOOLEAN
}

func IsVector(L *LuaState, n int32) bool {
	return Type(L, n) == LUA_TVECTOR
}

func IsThread(L *LuaState, n int32) bool {
	return Type(L, n) == LUA_TTHREAD
}

func IsBuffer(L *LuaState, n int32) bool {
	return Type(L, n) == LUA_TBUFFER
}

func IsNone(L *LuaState, n int32) bool {
	return Type(L, n) == LUA_TNONE
}

func IsNoneOrNil(L *LuaState, n int32) bool {
	return Type(L, n) <= LUA_TNIL
}

func PushLiteral(L *LuaState, s string) {
	PushLString(L, s, uint64(len(s)))
}

func PushCFunction(L *LuaState, f unsafe.Pointer) {
	PushCClosureK(L, f, new(string), 0, nil)
}

func PushCFunctionD(L *LuaState, f unsafe.Pointer, debugname *string) {
	PushCClosureK(L, f, debugname, 0, nil)
}

func PushCClosure(L *LuaState, f unsafe.Pointer, nup int32) {
	PushCClosureK(L, f, new(string), nup, nil)
}

func PushCClosureD(L *LuaState, f unsafe.Pointer, debugname *string, nup int32) {
	PushCClosureK(L, f, debugname, nup, nil)
}

func PushLightUserdata(L *LuaState, p unsafe.Pointer) {
	PushLightUserdataTagged(L, p, 0)
}

func SetGlobal(L *LuaState, global string) {
	SetField(L, LUA_GLOBALSINDEX, global)
}

func GetGlobal(L *LuaState, global string) int32 {
	return GetField(L, LUA_GLOBALSINDEX, global)
}

func ToString(L *LuaState, i int32) string {
	return ToLString(L, i, new(uint64))
}
