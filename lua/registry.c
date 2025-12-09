#include <stdlib.h>
#include <stdint.h>
#include <lua.h>
#include <_cgo_export.h>

int registryTrampoline(lua_State* L) {
    uintptr_t* handle_ptr = (uintptr_t*)lua_touserdata(L, lua_upvalueindex(1));
    return registryTrampolineImpl(L, *handle_ptr);
}

void registryTrampolineDtor(lua_State* L) {
    uintptr_t* handle_ptr = (uintptr_t*)lua_touserdata(L, lua_upvalueindex(1));
    registryTrampolineDtorImpl(L, *handle_ptr);
}
