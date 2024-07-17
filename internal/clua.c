#include <stdio.h>
#include <lua.h>
#include <_cgo_export.h>

void* clua_alloc(void* ud, void *ptr, size_t osize, size_t nsize)
{
	return (void*) go_allocf((GoUintptr) ud,(GoUintptr) ptr, osize, nsize);
}

lua_State* clua_newstate(void* goallocf)
{
	return lua_newstate(&clua_alloc, goallocf);
}

l_noret cluaL_errorL(lua_State* L, char* msg)
{
	return luaL_error(L, msg);
}
