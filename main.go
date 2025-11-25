package main

import (
	"fmt"

	"github.com/CompeyDev/lei/lua"
)

func main() {
	mem := lua.NewMemoryState()
	// mem.SetLimit(250 * 1024) // 250KB max
	state := lua.NewWith(lua.StdLibALLSAFE, lua.LuaOptions{InitMemoryState: mem})

	table := state.CreateTable()
	key, value := state.CreateString("hello"), state.CreateString("lei")
	table.Set(key, value)

	fmt.Printf("Used: %d, Limit: %d\n", mem.Used(), mem.Limit())

	fmt.Println(key.ToString(), table.Get(key).(*lua.LuaString).ToString())
	values, err := state.Execute("main", []byte("print('hello, lei!'); return {['mrrp'] = 'foo', ['meow'] = 'bar'}, 'baz'"))
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, value := range values {
		fmt.Print(i, ": ")

		switch v := value.(type) {
		case *lua.LuaString:
			fmt.Println(v.ToString())
		case *lua.LuaTable:
			fmt.Println()

			for key, val := range v.Iterable() {
				k, kErr := lua.As[string](key)
				v, vErr := lua.As[string](val)

				if kErr != nil || vErr != nil {
					fmt.Println("  (non-string key or value)")
				}

				fmt.Printf("  %v: %v\n", k, v)
			}
		}
	}

	iterable, iterErr := lua.As[map[string]string](table)
	if iterErr != nil {
		fmt.Println(iterErr)
		return
	}

	for k, v := range iterable { // or, we can use `.Iterable`
		fmt.Printf("%s %s\n", k, v)
	}
}
