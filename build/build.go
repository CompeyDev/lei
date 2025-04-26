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

func cloneSrc() string {
	color.Blue.Println("> Cloning luau-lang/luau")

	dir, tempDirErr := os.MkdirTemp("", "lei-build")
	bail(tempDirErr)

	// Clone down the Luau repo and checkout the required tag
	Exec("git", "", "clone", "https://github.com/luau-lang/luau.git", dir)
	Exec("git", dir, "checkout", LUAU_VERSION)

	color.Green.Printf("> Cloned repo to%s\n\n", dir)
	return dir
}

func buildVm(srcPath string, artifactPath string, includesDir string, cmakeFlags ...string) {
	color.Blue.Println("> Compile libLuau.VM.a")

	// Build the Luau VM using CMake
	buildDir := path.Join(srcPath, "cmake")
	buildDirErr := os.Mkdir(buildDir, os.ModePerm)
	if !os.IsExist(buildDirErr) {
		bail(buildDirErr)
	}

	defaultCmakeFlags := []string{"..", "-DCMAKE_BUILD_TYPE=RelWithDebInfo", "-DLUAU_EXTERN_C=ON", "-DCMAKE_POLICY_VERSION_MINIMUM=3.5"}
	Exec("cmake", buildDir, append(defaultCmakeFlags, cmakeFlags...)...)
	Exec("cmake", buildDir, "--build", ".", "--target Luau.VM", "--config", "RelWithDebInfo")

	color.Green.Println("> Successfully compiled!\n")

	// Copy the artifact to the artifact directory
	artifactFile, artifactErr := os.ReadFile(path.Join(buildDir, ARTIFACT_NAME))
	bail(artifactErr)
	bail(os.WriteFile(artifactPath, artifactFile, os.ModePerm))

	// Copy the header files into the includes directory
	headerDir := path.Join(srcPath, "VM", "include")
	headerFiles, headerErr := os.ReadDir(headerDir)
	bail(headerErr)
	for _, file := range headerFiles {
		src := path.Join(headerDir, file.Name())
		dest := path.Join(includesDir, file.Name())

		headerContents, headerReadErr := os.ReadFile(src)
		bail(headerReadErr)

		os.WriteFile(dest, headerContents, os.ModePerm)
	}
}

func main() {
	workDir, workDirErr := os.Getwd()
	bail(workDirErr)

	artifactDir := path.Join(workDir, ".lei")
	artifactPath := path.Join(artifactDir, ARTIFACT_NAME)
	lockfilePath := path.Join(artifactDir, ".lock")
	includesDir := path.Join(artifactDir, "includes")

	bail(os.MkdirAll(includesDir, os.ModePerm)) // includesDir is the deepest dir, creates all

	gitignore, gitignoreErr := os.ReadFile(".gitignore")
	if gitignoreErr == nil && !strings.Contains(string(gitignore), ".lei") {
		color.Yellow.Println("> WARN: The gitignore in the CWD does not include `.lei`, consider adding it")
	}

	// TODO: Args for clean build
	args := os.Args[1:]

	goArgs := []string{}
	cmakeFlags := []string{}
	features := []string{}

	// TODO: maybe use env vars for this config instead
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
	toCleanBuild := (string(lockfileContents) != serFeatures) || os.Getenv("LEI_CLEAN_BUILD") == "true"
	if _, err := os.Stat(artifactPath); err == nil && !toCleanBuild {
		fmt.Printf("[build] Using existing artifact at %s\n", artifactPath)
	} else {
		srcPath, notUnset := os.LookupEnv("LEI_LUAU_SRC")
		if !notUnset {
			srcPath = cloneSrc()
			defer os.RemoveAll(srcPath)
		}

		buildVm(srcPath, artifactPath, includesDir, cmakeFlags...)
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
		WithVar("CGO_CFLAGS", fmt.Sprintf("-I%s", includesDir)).
		WithVar("CGO_ENABLED", "1").
		PipeAll(Forward).
		ToCommand()

	bail(cmd.Start())
	bail(cmd.Wait())
}
