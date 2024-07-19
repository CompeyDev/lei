#include <stdlib.h>
#include <lua.h>

lua_State* clua_newstate(void* f, void* ud);
l_noret cluaL_errorL(lua_State* L, char* msg);
void clua_pushcclosurek(lua_State* L, void* f, char* debugname, int nup, void* cont);
void* clua_newuserdatadtor(lua_State* L, size_t sz, void* dtor);
void clua_setuserdatadtor(lua_State* L, int tag, void* dtor);