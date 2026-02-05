package ffi_test

import (
	"slices"
	"testing"

	"github.com/CompeyDev/lei/ffi"
)

func TestLuauCompile_Basic(t *testing.T) {
	source := `
		local function add(a, b)
			return a + b
		end
		return add(1, 2)
	`

	outsize := 0
	options := &ffi.CompileOptions{
		OptimizationLevel: 1,
		DebugLevel:        1,
		TypeInfoLevel:     0,
		CoverageLevel:     0,
	}

	bytecode := ffi.LuauCompile(source, len(source), options, &outsize)

	if bytecode == nil {
		t.Fatal("LuauCompile returned nil")
	}

	if outsize == 0 {
		t.Fatal("Output size is 0")
	}

	if len(bytecode) != outsize {
		t.Errorf("Expected bytecode length %d, got %d", outsize, len(bytecode))
	}

	t.Logf("Compiled successfully: %d bytes", outsize)
}

func TestLuauCompile_SyntaxError(t *testing.T) {
	source := `
		local function broken(
			-- missing closing parenthesis and end
	`

	outsize := 0
	options := &ffi.CompileOptions{
		OptimizationLevel: 1,
		DebugLevel:        1,
	}

	bytecode := ffi.LuauCompile(source, len(source), options, &outsize)

	// The function should still return bytecode containing the error
	if bytecode == nil {
		t.Fatal("LuauCompile returned nil even for error case")
	}

	t.Logf("Error bytecode: %d bytes", outsize)
}

func TestLuauCompile_WithOptions(t *testing.T) {
	source := `
		local x = vector.create(1, 2, 3)
		return x
	`

	outsize := 0
	options := &ffi.CompileOptions{
		OptimizationLevel: 2,
		DebugLevel:        2,
		TypeInfoLevel:     1,
		CoverageLevel:     1,
		VectorLib:         "vector",
		VectorCtor:        "create",
		VectorType:        "vector",
		MutableGlobals:    []string{"_G"},
		UserdataTypes:     []string{"MyUserdata"},
		DisabledBuiltins:  []string{"math.random"},
	}

	bytecode := ffi.LuauCompile(source, len(source), options, &outsize)

	if bytecode == nil {
		t.Fatal("LuauCompile returned nil")
	}

	if outsize == 0 {
		t.Fatal("Output size is 0")
	}

	t.Logf("Compiled with options: %d bytes", outsize)
}

func TestLuauCompile_EmptySource(t *testing.T) {
	source := ""
	outsize := 0
	options := &ffi.CompileOptions{
		OptimizationLevel: 1,
		DebugLevel:        1,
	}

	bytecode := ffi.LuauCompile(source, len(source), options, &outsize)

	if bytecode == nil {
		t.Fatal("LuauCompile returned nil for empty source")
	}

	t.Logf("Empty source compiled: %d bytes", outsize)
}

func TestLuauCompile_BinaryDataIntegrity(t *testing.T) {
	source := `return "test"`
	outsize := 0
	options := &ffi.CompileOptions{
		OptimizationLevel: 1,
		DebugLevel:        1,
	}

	bytecode := ffi.LuauCompile(source, len(source), options, &outsize)
	hasNullByte := slices.Contains(bytecode, 0)

	t.Logf("Bytecode contains null bytes: %v", hasNullByte)
	t.Logf("Bytecode length: %d, outsize: %d", len(bytecode), outsize)

	if len(bytecode) != outsize {
		t.Errorf("Bytecode length mismatch: expected %d, got %d", outsize, len(bytecode))
	}
}

func TestLuauCompile_ExecuteBytecode(t *testing.T) {
	source := `return 42`
	outsize := 0
	options := &ffi.CompileOptions{OptimizationLevel: 1, DebugLevel: 1}

	bytecode := ffi.LuauCompile(source, len(source), options, &outsize)

	L := ffi.LNewState()
	defer ffi.LuaClose(L)

	ffi.LuauLoad(L, "test", bytecode, uint64(outsize), 0)
	ffi.Pcall(L, 0, 1, 0)

	result := ffi.ToInteger(L, -1)
	if result != 42 {
		t.Error("Executed result did not match")
	}
}
