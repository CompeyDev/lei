package lua

import "github.com/CompeyDev/lei/ffi"

type LuaNil struct{}

//
// LuaValue Implementation
//

var _ LuaValue = (*LuaNil)(nil)

func (n *LuaNil) lua() *Lua        { return nil }
func (n *LuaNil) ref() int         { return ffi.LUA_REFNIL }
func (n *LuaNil) deref(_ *Lua) int { return 0 }
