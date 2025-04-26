package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/gookit/color"
	"golang.org/x/term"
)

const LUAU_VERSION = "0.634"
const ARTIFACT_NAME = "libLuau.VM.a"

func bail(err error) {
	if err != nil {
		panic(err)
	}
}

func buildVm(artifactPath string, cmakeFlags ...string) {
	color.Blue.Println("> Cloning luau-lang/luau")

	dir, homeDirErr := os.MkdirTemp("", "lei-build")
	bail(homeDirErr)

	defer os.RemoveAll(dir)

	// Clone down the Luau repo and checkout the required tag
	Exec("git", "", "clone", "https://github.com/luau-lang/luau.git", dir)
	Exec("git", dir, "checkout", LUAU_VERSION)

	color.Green.Printf("> Cloned repo to%s\n\n", dir)

	color.Blue.Println("> Compile libLuau.VM.a")

	// Build the Luau VM using CMake
	buildDir := path.Join(dir, "cmake")
	bail(os.Mkdir(buildDir, os.ModePerm))

	defaultCmakeFlags := []string{"..", "-DCMAKE_BUILD_TYPE=RelWithDebInfo", "-DLUAU_EXTERN_C=ON", "-DCMAKE_POLICY_VERSION_MINIMUM=3.5"}
	Exec("cmake", buildDir, append(defaultCmakeFlags, cmakeFlags...)...)
	Exec("cmake", buildDir, "--build", ".", "--target Luau.VM", "--config", "RelWithDebInfo")

	color.Green.Println("> Successfully compiled!\n")

	// Copy the artifact to the artifact directory
	artifactFile, artifactErr := os.ReadFile(path.Join(buildDir, ARTIFACT_NAME))
	bail(artifactErr)
	bail(os.WriteFile(artifactPath, artifactFile, os.ModePerm))
}

func main() {
	homeDir, homeDirErr := os.UserHomeDir()
	bail(homeDirErr)

	artifactDir := path.Join(homeDir, ".lei")
	artifactPath := path.Join(artifactDir, ARTIFACT_NAME)
	lockfilePath := path.Join(artifactDir, ".lock")

	bail(os.MkdirAll(artifactDir, os.ModePerm))

	// TODO: Args for clean build
	args := os.Args[1:]

	goArgs := []string{}
	cmakeFlags := []string{}
	features := []string{}

	for _, arg := range args {
		if arg == "--enable-vector4" {
			features = append(features, "LUAU_VECTOR4")
			// FIXME: This flag apparently isn't recognized by cmake for some reason
			cmakeFlags = append(cmakeFlags, "-DLUAU_VECTOR_SIZE=4")
		} else {
			goArgs = append(goArgs, arg)
		}
	}

	lockfileContents, err := os.ReadFile(lockfilePath)
	if !os.IsNotExist(err) {
		bail(err)
	}

	serFeatures := fmt.Sprintf("%v", features)
	toCleanBuild := string(lockfileContents) != serFeatures
	if _, err := os.Stat(artifactPath); err == nil && !toCleanBuild {
		fmt.Printf("[build] Using existing artifact at %s\n", artifactPath)
	} else {
		buildVm(artifactPath, cmakeFlags...)
		bail(os.WriteFile(lockfilePath, []byte(serFeatures), os.ModePerm))
	}

	buildTags := []string{}
	if len(features) > 0 {
		buildTags = append(buildTags, []string{"-tags", strings.Join(features, ",")}...)
	}

	w, _, termErr := term.GetSize(int(os.Stdout.Fd()))
	bail(termErr)
	fmt.Println(strings.Repeat("=", w))

	subcommand := goArgs[0]
	goArgs = goArgs[1:]
	combinedArgs := append(buildTags, goArgs...)
	cmd, _, _, _ := Command("go").
		WithArgs(append([]string{subcommand}, combinedArgs...)...).
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
