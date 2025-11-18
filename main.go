package main

import (
	"fmt"

	"github.com/CompeyDev/lei/lua"
)

func main() {
	mem := lua.NewMemoryState()
	mem.SetLimit(250 * 1024) // 250KB max
	state := lua.NewWith(lua.StdLibALLSAFE, lua.LuaOptions{InitMemoryState: mem})

	table := state.CreateTable()
	key, value := state.CreateString("hello"), state.CreateString("world")
	table.Set(&key, &value)

	fmt.Printf("Used: %d, Limit: %d\n", mem.Used(), mem.Limit())

	fmt.Println(key.ToString(), table.Get(&key).(*lua.LuaString).ToString())
}
