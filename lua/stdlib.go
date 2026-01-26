package lua

// StdLib represents flags describing the set of Lua standard libraries to load.
type StdLib uint32

const (
	// COROUTINE library
	// https://www.lua.org/manual/5.4/manual.html#6.2
	StdLibCOROUTINE StdLib = 1 << 0

	// TABLE library
	// https://www.lua.org/manual/5.4/manual.html#6.6
	StdLibTABLE StdLib = 1 << 1

	// OS library
	// https://www.lua.org/manual/5.4/manual.html#6.9
	StdLibOS StdLib = 1 << 3

	// STRING library
	// https://www.lua.org/manual/5.4/manual.html#6.4
	StdLibSTRING StdLib = 1 << 4

	// UTF8 library
	// https://www.lua.org/manual/5.4/manual.html#6.5
	StdLibUTF8 StdLib = 1 << 5

	// BIT library
	// https://www.lua.org/manual/5.2/manual.html#6.7
	StdLibBIT StdLib = 1 << 6

	// MATH library
	// https://www.lua.org/manual/5.4/manual.html#6.7
	StdLibMATH StdLib = 1 << 7

	// BUFFER library (Luau)
	// https://luau.org/library#buffer-library
	StdLibBUFFER StdLib = 1 << 9

	// VECTOR library (Luau)
	// https://luau.org/library#vector-library
	StdLibVECTOR StdLib = 1 << 10

	// DEBUG library (unsafe)
	// https://www.lua.org/manual/5.4/manual.html#6.10
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
		StdLibCOROUTINE: "COROUTINE",
		StdLibTABLE:     "TABLE",
		StdLibOS:        "OS",
		StdLibSTRING:    "STRING",
		StdLibUTF8:      "UTF8",
		StdLibBIT:       "BIT",
		StdLibMATH:      "MATH",
		StdLibBUFFER:    "BUFFER",
		StdLibVECTOR:    "VECTOR",
		StdLibDEBUG:     "DEBUG",
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
