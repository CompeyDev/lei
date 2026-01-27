package lua

import "github.com/CompeyDev/lei/ffi"

type LuaVector struct{ X, Y, Z float32 }

//
// LuaValue implementation
//

var _ LuaValue = (*LuaVector)(nil)

func (v LuaVector) lua() *Lua { return nil }
func (v LuaVector) ref() int  { return ffi.LUA_NOREF }
func (v LuaVector) deref(lua *Lua) int {
	state := lua.state()
	ffi.PushVector(state, v.X, v.Y, v.Z)
	return int(ffi.GetTop(state))
}
