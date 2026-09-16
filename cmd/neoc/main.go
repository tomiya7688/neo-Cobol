package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/tomiya7688/neo-Cobol/internal/compiler"
)

const version = "0.1.0-dev"

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "neoc:", err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		printUsage(stderr)
		return errors.New("missing command")
	}

	switch args[0] {
	case "version", "--version", "-version":
		fmt.Fprintf(stdout, "neoc %s\n", version)
		return nil
	case "check":
		return commandCheck(args[1:], stdout, stderr)
	case "emit-c":
		return commandEmitC(args[1:], stdout, stderr)
	case "build":
		return commandBuild(args[1:], stdout, stderr)
	case "run":
		return commandRun(args[1:], stdout, stderr)
	case "help", "--help", "-h":
		printUsage(stdout)
		return nil
	default:
		printUsage(stderr)
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func commandCheck(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	fs.SetOutput(stderr)
	if err := fs.Parse(args); err != nil {
		return err
	}
	path, err := requireSource(fs.Args())
	if err != nil {
		return err
	}
	source, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if _, err := compiler.Check(string(source)); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "%s: ok\n", path)
	return nil
}

func commandEmitC(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("emit-c", flag.ContinueOnError)
	fs.SetOutput(stderr)
	output := fs.String("o", "", "write generated C to this file")
	if err := fs.Parse(args); err != nil {
		return err
	}
	path, err := requireSource(fs.Args())
	if err != nil {
		return err
	}
	generated, err := compileFile(path)
	if err != nil {
		return err
	}
	if *output == "" {
		_, err = io.WriteString(stdout, generated)
		return err
	}
	return os.WriteFile(*output, []byte(generated), 0o644)
}

func commandBuild(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("build", flag.ContinueOnError)
	fs.SetOutput(stderr)
	output := fs.String("o", "", "output executable path")
	if err := fs.Parse(args); err != nil {
		return err
	}
	path, err := requireSource(fs.Args())
	if err != nil {
		return err
	}
	out := *output
	if out == "" {
		out = defaultExecutable(path)
	}
	if err := build(path, out, stderr); err != nil {
		return err
	}
	fmt.Fprintln(stdout, out)
	return nil
}

func commandRun(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	fs.SetOutput(stderr)
	if err := fs.Parse(args); err != nil {
		return err
	}
	path, err := requireSource(fs.Args())
	if err != nil {
		return err
	}
	dir, err := os.MkdirTemp("", "neoc-run-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	exe := filepath.Join(dir, "program")
	if runtime.GOOS == "windows" {
		exe += ".exe"
	}
	if err := build(path, exe, stderr); err != nil {
		return err
	}
	cmd := exec.Command(exe)
	cmd.Stdout, cmd.Stderr, cmd.Stdin = stdout, stderr, os.Stdin
	return cmd.Run()
}

func build(path, output string, stderr io.Writer) error {
	generated, err := compileFile(path)
	if err != nil {
		return err
	}
	dir, err := os.MkdirTemp("", "neoc-build-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	cfile := filepath.Join(dir, "program.c")
	if err := os.WriteFile(cfile, []byte(generated), 0o644); err != nil {
		return err
	}
	cc, err := findCCompiler()
	if err != nil {
		return err
	}
	cmd := exec.Command(cc, "-std=c11", cfile, "-o", output)
	cmd.Stdout = stderr
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("C compiler failed: %w", err)
	}
	return nil
}

func compileFile(path string) (string, error) {
	source, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return compiler.EmitC(string(source))
}

func findCCompiler() (string, error) {
	if configured := os.Getenv("CC"); configured != "" {
		if path, err := exec.LookPath(configured); err == nil {
			return path, nil
		}
		return "", fmt.Errorf("CC=%q was not found", configured)
	}
	for _, candidate := range []string{"cc", "clang", "gcc"} {
		if path, err := exec.LookPath(candidate); err == nil {
			return path, nil
		}
	}
	return "", errors.New("no C compiler found; install GCC/Clang or set CC")
}

func requireSource(args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("expected exactly one source file")
	}
	return args[0], nil
}

func defaultExecutable(source string) string {
	ext := filepath.Ext(source)
	out := source[:len(source)-len(ext)]
	if runtime.GOOS == "windows" {
		out += ".exe"
	}
	return out
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "usage: neoc <command> [options] <source.ncob>")
	fmt.Fprintln(w, "commands: check, emit-c, build, run, version")
}
