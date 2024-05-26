package cli

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

type CmdArgs struct {
	Url      string
	LocalDir string
	Branch   string
	Paths    []string
}

func ParseCmdArgs() CmdArgs {
	var args CmdArgs
	flag.StringVar(&args.Url, "url", "", "url of the repository to sparsely clone")
	flag.StringVar(&args.LocalDir, "dir", "", "local directory for cloned repo")
	flag.StringVar(&args.Branch, "branch", "main", "branch of repository to sparsely clone from")
	flag.Parse()
	args.Paths = flag.Args()
	if len(args.Url) == 0 {
		fmt.Fprintln(os.Stderr, "Missing required argument: -url")
		fmt.Fprintln(os.Stderr, "Use flag -help/-h for more info")
		os.Exit(1)
	}
	if len(args.LocalDir) == 0 {
		repoParts := strings.Split(args.Url, "/")
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
