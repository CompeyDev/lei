package ffi

/*
#cgo CFLAGS: -Iluau/VM/include
#cgo LDFLAGS: -Lluau/cmake -lLuau.VM -lm -lstdc++
#include <lua.h>
#include <lualib.h>
#include <stdlib.h>
#include "clua.h"

// From https://golang-nuts.narkive.com/UsNENgyt/cgo-how-to-pass-string-to-char-array
static char** make_char_array(int size) {
	return calloc(sizeof(char*), size);
}

static void set_array_string(char** a, char* s, int n) {
	a[n] = s;
}

static void free_char_array(char** a, int size) {
	int i;
	for (i = 0; i < size; i++)
	free(a[i]);
	free(a);
}
*/
import "C"
import "unsafe"

type LuaLReg C.luaL_Reg

func LRegister(L *LuaState, libname string, l *LuaLReg) {
	clibname := C.CString(libname)
	defer C.free(unsafe.Pointer(clibname))

	C.luaL_register(L, clibname, (*C.luaL_Reg)(l))
}

func LGetMetaField(L *LuaState, obj int32, e string) int32 {
	ce := C.CString(e)
	defer C.free(unsafe.Pointer(ce))

	return int32(C.luaL_getmetafield(L, C.int(obj), ce))
}

func LCallMeta(L *LuaState, obj int32, e string) int32 {
	ce := C.CString(e)
	defer C.free(unsafe.Pointer(ce))

	return int32(C.luaL_callmeta(L, C.int(obj), ce))
}

func LTypeError(L *LuaState, narg int32, tname string) {
	ctname := C.CString(tname)
	defer C.free(unsafe.Pointer(ctname))

	C.luaL_typeerrorL(L, C.int(narg), ctname)
}

func LArgError(L *LuaState, narg int32, extramsg string) {
	cextramsg := C.CString(extramsg)
	defer C.free(unsafe.Pointer(cextramsg))

	C.luaL_argerrorL(L, C.int(narg), cextramsg)
}

func LCheckLString(L *LuaState, narg int32, l *uint64) string {
	p := C.luaL_checklstring(L, C.int(narg), (*C.size_t)(l))
	defer C.free(unsafe.Pointer(p))

	return C.GoString(p)
}

func LOptLString(L *LuaState, narg int32, def string, l *uint64) string {
	cdef := C.CString(def)
	defer C.free(unsafe.Pointer(cdef))

	p := C.luaL_optlstring(L, C.int(narg), cdef, (*C.ulong)(l))
	defer C.free(unsafe.Pointer(p))

	return C.GoString(p)
}

func LCheckNumber(L *LuaState, narg int32) LuaNumber {
	return LuaNumber(C.luaL_checknumber(L, C.int(narg)))
}

func LOptNumber(L *LuaState, narg int32, def LuaNumber) LuaNumber {
	return LuaNumber(C.luaL_optnumber(L, C.int(narg), C.lua_Number(def)))
}

func LCheckBoolean(L *LuaState, narg int32) bool {
	return C.luaL_checkboolean(L, C.int(narg)) != 0
}

func LOptBoolean(L *LuaState, narg int32, def bool) bool {
	cdef := C.int(0)
	if def {
		cdef = C.int(1)
	}

	return C.luaL_optboolean(L, C.int(narg), cdef) != 0
}

func LCheckInteger(L *LuaState, narg int32) LuaInteger {
	return LuaInteger(C.luaL_checkinteger(L, C.int(narg)))
}

func LOptInteger(L *LuaState, narg int32, def LuaInteger) LuaInteger {
	return LuaInteger(C.luaL_optinteger(L, C.int(narg), C.lua_Integer(def)))
}

func LCheckUnsigned(L *LuaState, narg int32) LuaUnsigned {
	return LuaUnsigned(C.luaL_checkunsigned(L, C.int(narg)))
}

func LOptUnsigned(L *LuaState, narg int32, def LuaUnsigned) LuaUnsigned {
	return LuaUnsigned(C.luaL_optunsigned(L, C.int(narg), C.lua_Unsigned(def)))
}

func LCheckVector(L *LuaState, narg int32) *float32 {
	return (*float32)(C.luaL_checkvector(L, C.int(narg)))
}

func LOptVector(L *LuaState, narg int32, def *float32) *float32 {
	return (*float32)(C.luaL_optvector(L, C.int(narg), (*C.float)(def)))
}

func LCheckStack(L *LuaState, sz int32, msg string) {
	cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cmsg))

	C.luaL_checkstack(L, C.int(sz), cmsg)
}

func LCheckType(L *LuaState, narg int32, t int32) {
	C.luaL_checktype(L, C.int(narg), C.int(t))
}

func LCheckAny(L *LuaState, narg int32) {
	C.luaL_checkany(L, C.int(narg))
}

func LNewMetatable(L *LuaState, tname string) bool {
	ctname := C.CString(tname)
	defer C.free(unsafe.Pointer(ctname))

	return C.luaL_newmetatable(L, ctname) != 0
}

func LCheckUdata(L *LuaState, ud int32, tname string) unsafe.Pointer {
	ctname := C.CString(tname)
	defer C.free(unsafe.Pointer(ctname))

	return C.luaL_checkudata(L, C.int(ud), ctname)
}

func LCheckBuffer(L *LuaState, narg int32, len *uint64) unsafe.Pointer {
	return C.luaL_checkbuffer(L, C.int(narg), (*C.size_t)(len))
}

func LWhere(L *LuaState, lvl int32) {
	C.luaL_where(L, C.int(lvl))
}

// NOTE: It's not possible to pass varargs from Go->C via cgo, so instead we
// expect the user to format the message and hand it over to us, which we
// pass to luaL_errorL. This is an inconsistency with the actual C API, but
// there isn't really anything we can do.
func LErrorL(L *LuaState, msg string) {
	cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cmsg))

	PushString(L, msg)
	Error(L)

	// TODO: do we panic on the go side?
}

func LCheckOption(L *LuaState, narg int32, def string, lst []string) int32 {
	cdef := C.CString(def)
	defer C.free(unsafe.Pointer(cdef))

	clst := C.make_char_array(C.int(len(lst)))
	defer C.free_char_array(clst, C.int(len(lst)))
	for i, s := range lst {
		C.set_array_string(clst, C.CString(s), C.int(i))
	}

	return int32(C.luaL_checkoption(L, C.int(narg), cdef, clst))
}

func LToLString(L *LuaState, idx int32, len *uint64) string {
	p := C.luaL_tolstring(L, C.int(idx), (*C.size_t)(len))
	defer C.free(unsafe.Pointer(p))

	return C.GoString(p)
}

func LNewState() *LuaState {
	return C.luaL_newstate()
}

func LTypeName(L *LuaState, idx int32) string {
	return C.GoString(C.luaL_typename(L, C.int(idx)))
}

func LSandbox(L *LuaState) {
	C.luaL_sandbox(L)
}

func LSandboxThread(L *LuaState) {
	C.luaL_sandboxthread(L)
}

//
// Some useful macros
//

func LArgCheck(L *LuaState, cond bool, arg int32, extramsg string) {
	if cond {
		LArgError(L, arg, extramsg)
	}
}

func LArgExpected(L *LuaState, cond bool, arg int32, tname string) {
	if cond {
		LTypeError(L, arg, tname)
	}
}

func LCheckString(L *LuaState, n int32) string {
	return LCheckLString(L, n, nil)
}

func LOptString(L *LuaState, n int32, d string) string {
	return LOptLString(L, n, d, nil)
}

const (
	LUA_COLIBNAME     = "coroutine"
	LUA_TABLIBNAME    = "table"
	LUA_OSLIBNAME     = "os"
	LUA_STRLIBNAME    = "string"
	LUA_BITLIBNAME    = "bit32"
	LUA_BUFFERLIBNAME = "buffer"
	LUA_UTF8LIBNAME   = "utf8"
	LUA_MATHLIBNAME   = "math"
	LUA_DBLIBNAME     = "debug"
)

// DIVERGENCE: We cannot export wrapper functions around C functions if we want to
// pass them to API functions, we preserve the real C pointer by having 'opener'
// functions

func CoroutineOpener() C.lua_CFunction { return C.lua_CFunction(C.luaopen_base) }
func BaseOpener() C.lua_CFunction      { return C.lua_CFunction(C.luaopen_base) }
func TableOpener() C.lua_CFunction     { return C.lua_CFunction(C.luaopen_table) }
func OsOpener() C.lua_CFunction        { return C.lua_CFunction(C.luaopen_os) }
func StringOpener() C.lua_CFunction    { return C.lua_CFunction(C.luaopen_string) }
func Bit32Opener() C.lua_CFunction     { return C.lua_CFunction(C.luaopen_bit32) }
func BufferOpener() C.lua_CFunction    { return C.lua_CFunction(C.luaopen_buffer) }
func Utf8Opener() C.lua_CFunction      { return C.lua_CFunction(C.luaopen_utf8) }
func MathOpener() C.lua_CFunction      { return C.lua_CFunction(C.luaopen_math) }
func DebugOpener() C.lua_CFunction     { return C.lua_CFunction(C.luaopen_debug) }
func LibsOpener() C.lua_CFunction      { return C.lua_CFunction(C.luaL_openlibs) }

// TODO: More utility functions, buffer bindings
