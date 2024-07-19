package main

import lualib "github.com/CompeyDev/lei/internal"

func main() {
	lua := lualib.LNewState()
	println("Lua VM Address: ", lua)

	lualib.PushCClosureK(lua, func(L *lualib.LuaState) int32 {
		println("hi from closure?")
		return 0
	}, "test", 0, nil)

	if !lualib.IsCFunction(lua, 1) {
		panic("CFunction was not correctly pushed onto stack")
	}
}
