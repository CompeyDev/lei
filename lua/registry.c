#include <stdlib.h>
#include <stdint.h>
#include <lua.h>
#include <_cgo_export.h>

typedef struct registryTrampolineImpl_return trampolineResult;

int registryTrampoline(lua_State* L) {
    uintptr_t* handle_ptr = (uintptr_t*)lua_touserdata(L, lua_upvalueindex(1));
    trampolineResult result = registryTrampolineImpl(L, *handle_ptr);
   
    // Handle errors after crossing the C boundary to prevent a longjmp triggered
    // from the Go side, which would violate Go's stack winding rules 
    
    int status = result.r0;
    char* err  = result.r1;
   
    // TODO: Figure out what happens if some Lua code calls this without a pcall, longjmp?
    if (err != NULL) {
        lua_pushstring(L, err);
        lua_error(L);
    }

    return status;
}

void registryTrampolineDtor(lua_State* L) {
    uintptr_t* handle_ptr = (uintptr_t*)lua_touserdata(L, lua_upvalueindex(1));
    registryTrampolineDtorImpl(L, *handle_ptr);
}
