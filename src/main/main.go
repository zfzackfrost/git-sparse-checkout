package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/zfzackfrost/git-sparse-checkout/src/cli"
)

func main() {
	args := cli.ParseCmdArgs()

	var err error
	var gitPath string

	gitPath, err = exec.LookPath("git")
	handle_error(err, "Executable `git` not found!")

	err = os.MkdirAll(args.LocalDir, os.FileMode(0o644))
	handle_error(err, "Could not make directory: %s", args.LocalDir)

	err = os.Chdir(args.LocalDir)
	handle_error(err, "Could not cd to directory: %s", args.LocalDir)

	run_git_command(gitPath, "init")
	run_git_command(gitPath, "remote", "add", "-f", "origin", args.Repo)
	run_git_command(gitPath, "config", "core.sparseCheckout", "true")
	run_git_command(gitPath, "config", "core.sparseCheckout", "true")

	write_sparse_paths(args.Paths)

	run_git_command(gitPath, "pull", "origin", args.Branch)

}
func write_sparse_paths(paths []string) {
	f, err := os.OpenFile(path.Join(".git", "info", "sparse-checkout"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer f.Close()

	handle_error(err, "")
	for i := 0; i < len(paths); i++ {
		f.WriteString(paths[i] + "\n")
	}
}
func run_git_command(gitPath string, args ...string) {
	cmd := exec.Command(gitPath, args...)
	cmd.Run()
}
func handle_error(err error, pattern string, args ...any) {
	if err != nil {
		if len(pattern) > 0 {
			fmt.Fprintf(os.Stderr, pattern+"\n", args...)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
