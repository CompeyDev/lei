package main

import (
	"fmt"
	"os"
	"path"
	"strings"
)

const LUAU_VERSION = "0.634"
const ARTIFACT_NAME = "libLuau.VM.a"

func bail(err error) {
	if err != nil {
		panic(err)
	}
}

func buildVm(artifactDir string, artifactPath string, cmakeFlags ...string) {
	dir, homeDirErr := os.MkdirTemp("", "lei-build")
	bail(homeDirErr)

	defer os.RemoveAll(dir)

	// Clone down the Luau repo and checkout the required tag
	Exec("git", "", "clone", "https://github.com/luau-lang/luau.git", dir)
	Exec("git", dir, "checkout", LUAU_VERSION)

	// Build the Luau VM using CMake
	buildDir := path.Join(dir, "cmake")
	bail(os.Mkdir(buildDir, os.ModePerm))

	defaultCmakeFlags := []string{"..", "-DCMAKE_BUILD_TYPE=RelWithDebInfo", "-DLUAU_EXTERN_C=ON"}
	Exec("cmake", buildDir, append(defaultCmakeFlags, cmakeFlags...)...)
	Exec("cmake", buildDir, "--build", ".", "--target Luau.VM", "--config", "RelWithDebInfo")

	// Copy the artifact to the artifact directory
	artifactFile, artifactErr := os.ReadFile(artifactPath)
	bail(artifactErr)
	bail(os.WriteFile(path.Join(artifactDir, ARTIFACT_NAME), artifactFile, os.ModePerm))
}

func main() {
	homeDir, homeDirErr := os.UserHomeDir()
	bail(homeDirErr)

	artifactDir := path.Join(homeDir, ".lei")
	artifactPath := path.Join(artifactDir, ARTIFACT_NAME)

	bail(os.MkdirAll(artifactDir, os.ModePerm))

	// TODO: Args for clean build
	args := os.Args[1:]

	goArgs := []string{}
	cmakeFlags := []string{}
	features := []string{}

	for _, arg := range args {
		if arg == "--enable-vector4" {
			features = append(features, "LUAU_VECTOR4")
			cmakeFlags = append(cmakeFlags, "-DLUAU_VECTOR_SIZE=4")
		} else {
			goArgs = append(goArgs, arg)
		}
	}

	if _, err := os.Stat(artifactPath); err == nil {
		fmt.Printf("[build] Using existing artifact at %s\n", artifactPath)
	} else {
		buildVm(artifactDir, artifactPath, cmakeFlags...)
	}

	buildArgs := []string{"build"}
	if len(features) > 0 {
		buildArgs = append(buildArgs, []string{"-tags", strings.Join(features, ",")}...)
	}

	cmd, _, _, _ := Command("go").
		WithArgs(append(
			buildArgs, goArgs...,
		)...).
		WithVar(
			"CGO_LDFLAGS",
			fmt.Sprintf("-L %s -lLuau.VM -lm -lstdc++", artifactDir),
		).
		WithVar("CGO_ENABLED", "1").
		PipeAll(Forward).
		ToCommand()

	bail(cmd.Start())
	bail(cmd.Wait())
}