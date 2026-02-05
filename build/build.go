package main

import (
	"log"
	"os"
	"strings"
)

func main() {
	usage := func() { log.Fatal("Usage: buildProject <projects...>") }
	if len(os.Args) < 2 {
		usage()
	}

	switch os.Args[1] {
	case "buildProject":
		for _, project := range os.Args[2:] {
			if !strings.HasPrefix(project, "Luau.") {
				log.Fatalf("Invalid project name: %s", project)
			}

			compileLuauProject(project)
		}

	// Display usage menu
	case "-h", "--help":
		fallthrough
	default:
		usage()
	}
}

func compileLuauProject(project string) {
	if err := os.Mkdir("_obj", os.ModePerm); err == nil || !os.IsExist(err) {
		// Directory already exists, i.e., config files generated
		Exec(
			"cmake",
			"-S", "luau",
			"-B", "_obj",
			"-G", "Ninja",

			// Flags
			"-DCMAKE_BUILD_TYPE=RelWithDebInfo",
			"-DLUAU_EXTERN_C=ON",
		)
	}

	Exec("cmake", "--build", "_obj", "-t", project, "--config", "RelWithDebInfo")
}
