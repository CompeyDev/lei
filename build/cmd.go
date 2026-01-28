package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

type CommandPipeMode int

const (
	Forward CommandPipeMode = iota + 1
	Consume
)

type CommandBuilder struct {
	cmd      string
	args     []string
	env      map[string]string
	dir      string
	stdin    string
	pipeMode struct {
		stdin  CommandPipeMode
		stdout CommandPipeMode
		stderr CommandPipeMode
	}
}

func Command(cmd string) CommandBuilder {
	return CommandBuilder{cmd: cmd}
}

func (c CommandBuilder) WithArgs(args ...string) CommandBuilder {
	c.args = args
	return c
}

func (c CommandBuilder) WithStdin(stdin string) CommandBuilder {
	c.stdin = stdin
	return c
}

func (c CommandBuilder) WithVar(env string, val string) CommandBuilder {
	if c.env == nil {
		c.env = make(map[string]string)
	}

	c.env[env] = val
	return c
}

func (c CommandBuilder) Dir(dir string) CommandBuilder {
	c.dir = dir
	return c
}

func (c CommandBuilder) PipeStdin(mode CommandPipeMode) CommandBuilder {
	c.pipeMode.stdin = mode
	return c
}

func (c CommandBuilder) PipeStdout(mode CommandPipeMode) CommandBuilder {
	c.pipeMode.stdout = mode
	return c
}

func (c CommandBuilder) PipeStderr(mode CommandPipeMode) CommandBuilder {
	c.pipeMode.stderr = mode
	return c
}

func (c CommandBuilder) PipeAll(mode CommandPipeMode) CommandBuilder {
	c.pipeMode.stdin = mode
	c.pipeMode.stdout = mode
	c.pipeMode.stderr = mode
	return c
}

func (c CommandBuilder) ToCommand() (exec.Cmd, io.Reader, io.Writer, io.Writer) {
	stdinReader := pipeModeToReader(c.pipeMode.stdin, os.Stdin, c.stdin)
	stdoutWriter := pipeModeToWriter(c.pipeMode.stdout, os.Stdout)
	stderrWriter := pipeModeToWriter(c.pipeMode.stderr, os.Stderr)

	cmdPath, err := exec.LookPath(c.cmd)
	if err != nil {
		panic(err)
	}

	var env []string = os.Environ()
	for envVar, val := range c.env {
		env = append(env, envVar+"="+val)
	}

	cmd := exec.Cmd{
		Path:   cmdPath,
		Args:   append([]string{cmdPath}, c.args...),
		Dir:    c.dir,
		Env:    env,
		Stdin:  stdinReader,
		Stdout: stdoutWriter,
		Stderr: stderrWriter,
	}

	return cmd, stdinReader, stdoutWriter, stderrWriter
}

func pipeModeToReader(mode CommandPipeMode, def io.Reader, input string) io.Reader {
	switch mode {
	case Forward:
		return def
	case Consume:
		return bytes.NewReader([]byte(input))
	default:
		panic("invalid pipe mode")
	}
}

func pipeModeToWriter(mode CommandPipeMode, def io.Writer) io.Writer {
	switch mode {
	case Forward:
		return def
	case Consume:
		return bytes.NewBuffer([]byte{})
	default:
		panic("invalid pipe mode")
	}
}

func Exec(exe string, args ...string) {
	cmd, _, _, _ := Command(exe).
		WithArgs(args...).
		WithVar("CLICOLOR_FORCE", "1").
		PipeAll(Forward).
		ToCommand()

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start command %s: %v", exe, err)
	}

	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			commandStr := strings.Join(append([]string{exe}, args...), " ")
			log.Fatalf("'%s' exited with %d", commandStr, exitErr.ExitCode())
		}

		log.Fatalf("%s command failed: %v", exe, err)
	}
}
