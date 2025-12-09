#include <stdlib.h>
#include <stdint.h>
#include <lua.h>
#include <_cgo_export.h>

int registryTrampoline(lua_State* L) {
    uintptr_t registry_ptr = (uintptr_t)lua_touserdata(L, lua_upvalueindex(1));
    uintptr_t func_id = (uintptr_t)lua_touserdata(L, lua_upvalueindex(2));
    return registryTrampolineImpl(L, registry_ptr, func_id);
}
