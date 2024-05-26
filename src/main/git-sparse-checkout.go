package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"sync"

	"github.com/zfzackfrost/git-sparse-checkout/src/cli"
)

func main() {
	args := cli.ParseCmdArgs()

	var err error
	runGitCommand := gitCommandRunner()

	err = os.MkdirAll(args.LocalDir, os.FileMode(0o644))
	handleError(err, "Could not make directory: %s", args.LocalDir)

	err = os.Chdir(args.LocalDir)
	handleError(err, "Could not cd to directory: %s", args.LocalDir)

	runGitCommand("init")
	runGitCommand("remote", "add", "-f", "origin", args.Url)
	runGitCommand("config", "core.sparseCheckout", "true")

	writeSparsePaths(args.Paths)

	runGitCommand("pull", "origin", args.Branch)

}
func writeSparsePaths(paths []string) {
	f, err := os.OpenFile(path.Join(".git", "info", "sparse-checkout"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer f.Close()

	handleError(err, "")
	for i := 0; i < len(paths); i++ {
		f.WriteString(paths[i] + "\n")
	}
}
func gitCommandRunner() func(args ...string) {
	gitPath, err := exec.LookPath("git")
	handleError(err, "Executable `git` not found!")

	return func(args ...string) {
		cmd := exec.Command(gitPath, args...)

		var wg sync.WaitGroup
		if stdout, err := cmd.StdoutPipe(); err == nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				io.Copy(os.Stdout, stdout)
			}()
		}
		if stderr, err := cmd.StderrPipe(); err == nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				io.Copy(os.Stderr, stderr)
			}()
		}

		cmd.Run()
		wg.Wait()
	}
}
func handleError(err error, pattern string, args ...any) {
	if err != nil {
		if len(pattern) > 0 {
			fmt.Fprintf(os.Stderr, pattern+"\n", args...)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
