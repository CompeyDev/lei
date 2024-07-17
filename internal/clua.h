#include <stdlib.h>
#include <lua.h>

void* clua_alloc(void* ud, void *ptr, size_t osize, size_t nsize);
lua_State* clua_newstate(void* goallocf);
l_noret cluaL_errorL(lua_State* L, char* msg);