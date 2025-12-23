package lua

import "github.com/CompeyDev/lei/ffi"

type LuaNumber float64

//
// LuaValue implementation
//

var _ LuaValue = (*LuaNumber)(nil)

// Numbers are cheap to copy, so we don't store the reference index

func (n *LuaNumber) lua() *Lua { return nil }
func (n *LuaNumber) ref() int  { return ffi.LUA_NOREF }
func (n *LuaNumber) deref(lua *Lua) int {
	state := lua.state()
	ffi.PushNumber(state, ffi.LuaNumber(*n))
	return int(ffi.GetTop(state))
}
