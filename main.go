package main

import lualib "github.com/CompeyDev/lei/internal"

func main() {
	lua := lualib.LNewState()
	lualib.PushCClosureK(lua, func(L *lualib.LuaState) int32 {
		println("hi from closure? ")
		return 0
	}, "test", 0, nil)

	println(lua)
}
