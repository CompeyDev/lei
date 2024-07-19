#include <stdio.h>
#include <lua.h>
#include <_cgo_export.h>

// void* clua_alloc(void* ud, void *ptr, size_t osize, size_t nsize)
// {
// 	return (void*) go_allocf((GoUintptr) ud,(GoUintptr) ptr, osize, nsize);
// }


lua_State* clua_newstate(void* f, void* ud)
{
	return lua_newstate((lua_Alloc)f, ud);
}

l_noret cluaL_errorL(lua_State* L, char* msg)
{
	luaL_errorL(L, msg);
}

void clua_pushcclosurek(lua_State* L, void* f, char* debugname, int nup, void* cont) {
	return lua_pushcclosurek(L, (lua_CFunction)f, debugname, nup, (lua_Continuation)cont);
}

void* clua_newuserdatadtor(lua_State* L, size_t sz, void* dtor) {
	return lua_newuserdatadtor(L, sz, (void (*)(void*))dtor);
}

void clua_setuserdatadtor(lua_State* L, int tag, void* dtor) {
	return lua_setuserdatadtor(L, tag, (lua_Destructor)dtor);
}