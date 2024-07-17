package internal

/*
#cgo CFLAGS: -Iluau/VM/include -I/usr/lib/gcc/x86_64-pc-linux-gnu/14.1.1/include -I${SRCDIR}
#include <lua.h>
#include <lualib.h>
#include <stdlib.h>
#include <clua.h>

// From https://golang-nuts.narkive.com/UsNENgyt/cgo-how-to-pass-string-to-char-array
static char** makeCharArray(int size) {
	return calloc(sizeof(char*), size);
}

static void setArrayString(char** a, char* s, int n) {
	a[n] = s;
}

static void freeCharArray(char** a, int size) {
	int i;
	for (i = 0; i < size; i++)
	free(a[i]);
	free(a);
}
*/
import "C"
import "unsafe"

type luaL_Reg C.luaL_Reg

func LRegister(L *C.lua_State, libname string, l *luaL_Reg) {
	clibname := C.CString(libname)
	defer C.free(unsafe.Pointer(clibname))

	C.luaL_register(L, clibname, (*C.luaL_Reg)(l))
}

func LGetMetaField(L *C.lua_State, obj int32, e string) int32 {
	ce := C.CString(e)
	defer C.free(unsafe.Pointer(ce))

	return int32(C.luaL_getmetafield(L, C.int(obj), ce))
}

func LCallMeta(L *C.lua_State, obj int32, e string) int32 {
	ce := C.CString(e)
	defer C.free(unsafe.Pointer(ce))

	return int32(C.luaL_callmeta(L, C.int(obj), ce))
}

func LTypeError(L *C.lua_State, narg int32, tname string) {
	ctname := C.CString(tname)
	defer C.free(unsafe.Pointer(ctname))

	C.luaL_typeerrorL(L, C.int(narg), ctname)
}

func LArgError(L *C.lua_State, narg int32, extramsg string) {
	cextramsg := C.CString(extramsg)
	defer C.free(unsafe.Pointer(cextramsg))

	C.luaL_argerrorL(L, C.int(narg), cextramsg)
}

func LCheckLString(L *C.lua_State, narg int32, l *uint64) string {
	p := C.luaL_checklstring(L, C.int(narg), (*C.size_t)(l))
	defer C.free(unsafe.Pointer(p))

	return C.GoString(p)
}

func LOptLString(L *C.lua_State, narg int32, def string, l *uint64) string {
	cdef := C.CString(def)
	defer C.free(unsafe.Pointer(cdef))

	p := C.luaL_optlstring(L, C.int(narg), cdef, (*C.ulong)(l))
	defer C.free(unsafe.Pointer(p))

	return C.GoString(p)
}

func LCheckNumber(L *C.lua_State, narg int32) lua_Number {
	return lua_Number(C.luaL_checknumber(L, C.int(narg)))
}

func LOptNumber(L *C.lua_State, narg int32, def lua_Number) lua_Number {
	return lua_Number(C.luaL_optnumber(L, C.int(narg), C.lua_Number(def)))
}

func LCheckBoolean(L *C.lua_State, narg int32) bool {
	return C.luaL_checkboolean(L, C.int(narg)) != 0
}

func LOptBoolean(L *C.lua_State, narg int32, def bool) bool {
	cdef := C.int(0)
	if def {
		cdef = C.int(1)
	}

	return C.luaL_optboolean(L, C.int(narg), cdef) != 0
}

func LCheckInteger(L *C.lua_State, narg int32) lua_Integer {
	return lua_Integer(C.luaL_checkinteger(L, C.int(narg)))
}

func LOptInteger(L *C.lua_State, narg int32, def lua_Integer) lua_Integer {
	return lua_Integer(C.luaL_optinteger(L, C.int(narg), C.lua_Integer(def)))
}

func LCheckUnsigned(L *C.lua_State, narg int32) lua_Unsigned {
	return lua_Unsigned(C.luaL_checkunsigned(L, C.int(narg)))
}

func LOptUnsigned(L *C.lua_State, narg int32, def lua_Unsigned) lua_Unsigned {
	return lua_Unsigned(C.luaL_optunsigned(L, C.int(narg), C.lua_Unsigned(def)))
}

func LCheckVector(L *C.lua_State, narg int32) *float32 {
	return (*float32)(C.luaL_checkvector(L, C.int(narg)))
}

func LOptVector(L *C.lua_State, narg int32, def *float32) *float32 {
	return (*float32)(C.luaL_optvector(L, C.int(narg), (*C.float)(def)))
}

func LCheckStack(L *C.lua_State, sz int32, msg string) {
	cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cmsg))

	C.luaL_checkstack(L, C.int(sz), cmsg)
}

func LCheckType(L *C.lua_State, narg int32, t int32) {
	C.luaL_checktype(L, C.int(narg), C.int(t))
}

func LCheckAny(L *C.lua_State, narg int32) {
	C.luaL_checkany(L, C.int(narg))
}

func LNewMetatable(L *C.lua_State, tname string) bool {
	ctname := C.CString(tname)
	defer C.free(unsafe.Pointer(ctname))

	return C.luaL_newmetatable(L, ctname) != 0
}

func LCheckUdata(L *C.lua_State, ud int32, tname string) unsafe.Pointer {
	ctname := C.CString(tname)
	defer C.free(unsafe.Pointer(ctname))

	return C.luaL_checkudata(L, C.int(ud), ctname)
}

func LCheckBuffer(L *C.lua_State, narg int32, len *uint64) unsafe.Pointer {
	return C.luaL_checkbuffer(L, C.int(narg), (*C.size_t)(len))
}

func LWhere(L *C.lua_State, lvl int32) {
	C.luaL_where(L, C.int(lvl))
}

// NOTE: It's not possible to pass varargs from Go->C via cgo, so instead we
// expect the user to format the message and hand it over to us, which we
// pass to luaL_errorL. This is an inconsistency with the actual C API, but
// there isn't really anything we can do.
func LErrorL(L *C.lua_State, msg string) {
	cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cmsg))

	C.cluaL_errorL(L, cmsg)
}

func LCheckOption(L *C.lua_State, narg int32, def string, lst []string) int32 {
	cdef := C.CString(def)
	defer C.free(unsafe.Pointer(cdef))

	clst := C.makeCharArray(C.int(len(lst)))
	defer C.freeCharArray(clst, C.int(len(lst)))
	for i, s := range lst {
		C.setArrayString(clst, C.CString(s), C.int(i))
	}

	return int32(C.luaL_checkoption(L, C.int(narg), cdef, clst))
}
