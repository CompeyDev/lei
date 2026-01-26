package lua

import "github.com/CompeyDev/lei/ffi"

// StdLib represents flags describing the set of Lua standard libraries to load.
type StdLib uint32

const (
	// COROUTINE library
	// https://luau.org/library#coroutine-library
	StdLibCOROUTINE StdLib = 1 << 0

	// TABLE library
	// https://luau.org/library#table-library
	StdLibTABLE StdLib = 1 << 1

	// OS library
	// https://luau.org/library#os-library
	StdLibOS StdLib = 1 << 3

	// STRING library
	// https://luau.org/library#string-library
	StdLibSTRING StdLib = 1 << 4

	// UTF8 library
	// https://luau.org/library#utf8-library
	StdLibUTF8 StdLib = 1 << 5

	// BIT library
	// https://luau.org/library#bit32-library
	StdLibBIT StdLib = 1 << 6

	// MATH library
	// https://luau.org/library#math-library
	StdLibMATH StdLib = 1 << 7

	// BUFFER library (Luau)
	// https://luau.org/library#buffer-library
	StdLibBUFFER StdLib = 1 << 9

	// VECTOR library (Luau)
	// https://luau.org/library#vector-library
	StdLibVECTOR StdLib = 1 << 10

	// DEBUG library (unsafe)
	// https://luau.org/library#debug-library
	StdLibDEBUG StdLib = 1 << 31

	// StdLibNONE represents no libraries
	StdLibNONE StdLib = 0

	// StdLibALL represents all standard libraries (unsafe)
	StdLibALL StdLib = ^StdLib(0) // equivalent to uint32 max

	// StdLibALLSAFE represents the safe subset of standard libraries
	StdLibALLSAFE StdLib = (1 << 30) - 1
)

func (s StdLib) Contains(lib StdLib) bool {
	return (s & lib) != 0
}

func (s StdLib) And(lib StdLib) StdLib {
	return s & lib
}

func (s StdLib) Or(lib StdLib) StdLib {
	return s | lib
}

func (s StdLib) Xor(lib StdLib) StdLib {
	return s ^ lib
}

func (s *StdLib) Add(lib StdLib) {
	*s |= lib
}

func (s *StdLib) Remove(lib StdLib) {
	*s &^= lib
}

func (s *StdLib) Toggle(lib StdLib) {
	*s ^= lib
}

func (s StdLib) String() string {
	if s == StdLibNONE {
		return "NONE"
	}
	if s == StdLibALL {
		return "ALL"
	}

	var libs []string
	flags := map[StdLib]string{
		StdLibCOROUTINE: ffi.LUA_COLIBNAME,
		StdLibTABLE:     ffi.LUA_TABLIBNAME,
		StdLibOS:        ffi.LUA_OSLIBNAME,
		StdLibSTRING:    ffi.LUA_STRLIBNAME,
		StdLibUTF8:      ffi.LUA_UTF8LIBNAME,
		StdLibBIT:       ffi.LUA_BITLIBNAME,
		StdLibMATH:      ffi.LUA_MATHLIBNAME,
		StdLibBUFFER:    ffi.LUA_BUFFERLIBNAME,
		StdLibVECTOR:    ffi.LUA_VECLIBNAME,
		StdLibDEBUG:     ffi.LUA_VECLIBNAME,
	}

	for flag, name := range flags {
		if s.Contains(flag) {
			libs = append(libs, name)
		}
	}

	if len(libs) == 0 {
		return "NONE"
	}

	result := ""
	for i, lib := range libs {
		if i > 0 {
			result += "|"
		}
		result += lib
	}
	return result
}
