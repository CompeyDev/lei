#include <string.h>
#include <lua.h>
#include <_cgo_export.h>

int indexMt(lua_State* L) {
    const char* key = lua_tostring(L, 2);
    if (key == NULL) {
        lua_pushnil(L);
        return 1;
    }
    
    uintptr_t* fields_handle = (uintptr_t*)lua_touserdata(L, lua_upvalueindex(1));
    uintptr_t* methods_handle = (uintptr_t*)lua_touserdata(L, lua_upvalueindex(2));

    indexMtImpl(L, *fields_handle, *methods_handle, (char*)key);
    return 1;
}
