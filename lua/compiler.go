package lua

import "C"
import (
	"github.com/CompeyDev/lei/ffi"
)

type Compiler struct{ options *ffi.CompileOptions }

func (c *Compiler) WithOptimizationLevel(lvl int) *Compiler {
	c.options.OptimizationLevel = lvl
	return c
}

func (c *Compiler) WithDebugLevel(lvl int) *Compiler {
	c.options.DebugLevel = lvl
	return c
}

func (c *Compiler) WithTypeInfoLevel(lvl int) *Compiler {
	c.options.TypeInfoLevel = lvl
	return c
}

func (c *Compiler) WithCoverageLevel(lvl int) *Compiler {
	c.options.CoverageLevel = lvl
	return c
}

func (c *Compiler) WithMutableGlobals(globals []string) *Compiler {
	c.options.MutableGlobals = append(c.options.MutableGlobals, globals...)
	return c
}

func (c *Compiler) WithUserdataTypes(types []string) *Compiler {
	c.options.UserdataTypes = append(c.options.UserdataTypes, types...)
	return c
}

func (c *Compiler) WithConstantLibraries(libs []string) *Compiler {
	c.options.LibrariesWithKnownMembers = append(c.options.LibrariesWithKnownMembers, libs...)
	return c
}

func (c *Compiler) WithDisabledBuiltins(builtins []string) *Compiler {
	c.options.DisabledBuiltins = append(c.options.DisabledBuiltins, builtins...)
	return c
}

func (c *Compiler) Compile(source string) ([]byte, error) {
	outsize := 0
	bytecode := ffi.LuauCompile(source, len(source), c.options, &outsize)

	// Check for compilation error
	// If bytecode starts with 0, the rest is an error message starting with ':'
	// See https://github.com/luau-lang/luau/blob/0.671/Compiler/src/Compiler.cpp#L4410
	if outsize > 0 && bytecode[0] == 0 {
		// Extract error message (skip the 0 byte and ':' character)
		message := ""
		if outsize > 2 {
			message = string(bytecode[2:])
		}

		// Check if input is incomplete (ends with <eof>)
		incompleteInput := len(message) > 0 &&
			(len(message) >= 5 && message[len(message)-5:] == "<eof>")

		return nil, &SyntaxError{
			IncompleteInput: incompleteInput,
			Message:         message,
		}
	}

	return bytecode, nil
}

func DefaultCompiler() *Compiler {
	return &Compiler{options: &ffi.CompileOptions{
		OptimizationLevel: 1,
		DebugLevel:        1,
		TypeInfoLevel:     0,
		CoverageLevel:     0,
		MutableGlobals:    make([]string, 0),
	}}
}

type SyntaxError struct {
	IncompleteInput bool
	Message         string
}

func (e *SyntaxError) Error() string {
	if e.IncompleteInput {
		return "incomplete input: " + e.Message
	}

	return "syntax error: " + e.Message
}

func isBytecode(data []byte) bool {
	// Luau bytecode starts with a version byte (currently 0-5 range)
	// See: https://github.com/luau-lang/luau/blob/0.671/Compiler/src/BytecodeBuilder.cpp#L13
	if len(data) == 0 {
		return false
	}

	// Check if the first byte is within the bytecode versionByte range (source code starting with
	// these bytes would be extremely rare)
	versionByte := data[0]
	return versionByte >= ffi.LBC_VERSION_MIN && versionByte <= ffi.LBC_VERSION_MAX
}
