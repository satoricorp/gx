package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/satoricorp/gx/internal/cli"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: %s <bash-output> <zsh-output>\n", filepath.Base(os.Args[0]))
		os.Exit(2)
	}

	root := cli.NewRoot(context.Background())
	if err := writeFile(os.Args[1], func(file *os.File) error {
		return root.GenBashCompletion(file)
	}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := writeFile(os.Args[2], func(file *os.File) error {
		return root.GenZshCompletion(file)
	}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func writeFile(path string, generate func(*os.File) error) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create completion dir for %s: %w", path, err)
	}
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer file.Close()
	if err := generate(file); err != nil {
		return fmt.Errorf("generate %s: %w", path, err)
	}
	return nil
}
