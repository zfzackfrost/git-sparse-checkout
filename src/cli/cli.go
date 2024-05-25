package cli

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

type CmdArgs struct {
	Repo     string
	LocalDir string
	Branch   string
	Paths    []string
}

func ParseCmdArgs() CmdArgs {
	var args CmdArgs
	flag.StringVar(&args.Repo, "repo", "", "repository to sparsely clone")
	flag.StringVar(&args.LocalDir, "dir", "", "local directory to clone repo to")
	flag.StringVar(&args.Branch, "branch", "main", "branch of repository to sparsely clone from")
	flag.Parse()
	args.Paths = flag.Args()
	if len(args.Repo) == 0 {
		fmt.Fprintln(os.Stderr, "repo must not be empty!")
		os.Exit(1)
	}
	if len(args.LocalDir) == 0 {
		repoParts := strings.Split(args.Repo, "/")
		repoName := repoParts[len(repoParts)-1]
		if cwd, err := os.Getwd(); err == nil {
			args.LocalDir = path.Join(cwd, repoName)
		} else {
			fmt.Fprintln(os.Stderr, "Error getting CWD:")
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	}
	return args
}
