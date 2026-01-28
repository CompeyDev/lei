package lua_test

import (
	"testing"

	"github.com/CompeyDev/lei/lua"
)

func TestAs(t *testing.T) {
	state := lua.New()

	// 1. Tag match
	t.Run("Tag match", func(t *testing.T) {
		type Person struct {
			Name string `lua:"username"`
		}

		table := state.CreateTable()
		table.Set(state.CreateString("username"), state.CreateString("Alice"))

		res, err := lua.As[Person](table)
		if err != nil {
			t.Fatal(err)
		}

		if res.Name != "Alice" {
			t.Fatalf("expected Alice, got %v", res.Name)
		}
	})

	// 2. Exact field name match
	t.Run("Exact match", func(t *testing.T) {
		type Person struct {
			Age string // TODO: make this int once we have numbers
		}

		table := state.CreateTable()
		table.Set(state.CreateString("Age"), state.CreateString("30"))

		res, err := lua.As[Person](table)
		if err != nil {
			t.Fatal(err)
		}

		if res.Age != "30" {
			t.Fatalf("expected '30', got %v", res.Age)
		}
	})

	// 3. Lowercase-first-letter fallback
	t.Run("Lowercase fallback", func(t *testing.T) {
		type Person struct{ Country string }

		table := state.CreateTable()
		table.Set(state.CreateString("country"), state.CreateString("Germany"))

		res, err := lua.As[Person](table)
		if err != nil {
			t.Fatal(err)
		}

		if res.Country != "Germany" {
			t.Fatalf("expected 'Germany', got %v", res.Country)
		}
	})

	// 4. Unexported field ignored
	t.Run("Unexported field", func(t *testing.T) {
		type Box struct{ secret string }

		table := state.CreateTable()
		table.Set(state.CreateString("secret"), state.CreateString("trains r cool"))

		res, err := lua.As[Box](table)
		if err != nil {
			t.Fatal(err)
		}

		if res.secret != "" {
			t.Fatalf("expected empty, got %v", res.secret)
		}
	})

	// 5. Mixed fields
	t.Run("Mixed fields", func(t *testing.T) {
		type Person struct {
			Name  string `lua:"username"`
			Age   string // TODO: use int once LuaNumber is implemented
			Email string
		}

		table := state.CreateTable()
		table.Set(state.CreateString("username"), state.CreateString("Bob"))
		table.Set(state.CreateString("Age"), state.CreateString("25"))
		table.Set(state.CreateString("email"), state.CreateString("bobby@example.com"))

		res, err := lua.As[Person](table)
		if err != nil {
			t.Fatal(err)
		}

		if res.Name != "Bob" || res.Age != "25" || res.Email != "bobby@example.com" {
			t.Fatalf("unexpected result: %+v", res)
		}
	})

	// 6. Missing Lua key
	t.Run("Missing key", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		table := state.CreateTable()
		table.Set(state.CreateString("Name"), state.CreateString("Johnny"))

		res, err := lua.As[Person](table)
		if err != nil {
			t.Fatal(err)
		}

		if res.Name != "Johnny" || res.Age != 0 {
			t.Fatalf("unexpected result: %+v", res)
		}
	})

	// 7. Extra Lua key ignored
	t.Run("Extra key ignored", func(t *testing.T) {
		type Person struct{ Name string }

		table := state.CreateTable()
		table.Set(state.CreateString("unknown"), state.CreateTable())

		res, err := lua.As[Person](table)
		if err != nil {
			t.Fatal(err)
		}

		if res.Name != "" {
			t.Fatalf("expected Name empty, got %v", res.Name)
		}
	})

	// 8. Tag overrides lowercase fallback
	t.Run("Tag overrides fallback", func(t *testing.T) {
		type Person struct {
			Name string `lua:"user"`
		}

		table := state.CreateTable()
		table.Set(state.CreateString("name"), state.CreateString("Dave"))
		table.Set(state.CreateString("user"), state.CreateString("Eve"))

		res, err := lua.As[Person](table)
		if err != nil {
			t.Fatal(err)
		}

		if res.Name != "Eve" {
			t.Fatalf("expected 'Eve', got %v", res.Name)
		}
	})
}
