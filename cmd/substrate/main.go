package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/rome2112/substrate/pkg/execgate"
	"github.com/rome2112/substrate/pkg/procfreeze"
	"github.com/rome2112/substrate/pkg/ramguard"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	switch os.Args[1] {
	case "ramguard":
		ramguardCmd(ctx, os.Args[2:])
	case "procfreeze":
		procfreezeCmd(ctx, os.Args[2:])
	case "execgate":
		execgateCmd(ctx, os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand: %s\n", os.Args[1])
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage:
  substrate <subcommand> [options]

Subcommands:
  ramguard    - run a command with ephemeral secret injection
  procfreeze  - freeze/unfreeze process groups from a manifest
  execgate    - inspect a command for high-risk patterns and run with confirmation

Run "substrate <subcommand> --help" for subcommand options.
`)
}

func ramguardCmd(ctx context.Context, args []string) {
	fs := flag.NewFlagSet("ramguard", flag.ExitOnError)
	path := fs.String("path", "", "executable path to run (required)")
	arglist := fs.String("args", "", "space-separated args (optional)")
	// secret can be provided multiple times; but for simplicity provide once as KEY=VALUE
	secret := fs.String("secret", "", "sensitive env pair (KEY=VALUE). Do NOT use files.")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, "parse:", err)
		os.Exit(2)
	}
	if *path == "" {
		fmt.Fprintln(os.Stderr, "path is required")
		fs.Usage()
		os.Exit(2)
	}

	secrets := map[string]string{}
	if *secret != "" {
		// simple KEY=VALUE parser
		p := *secret
		for i := 0; i < 1; i++ { // placeholder for potential extension to multiple secrets
			var k, v string
			n := 0
			for i := 0; i < len(p); i++ {
				if p[i] == '=' {
					k = p[:i]
					v = p[i+1:]
					n = i
					break
				}
			}
			if n == 0 || k == "" {
				fmt.Fprintln(os.Stderr, "secret must be KEY=VALUE")
				os.Exit(2)
			}
			secrets[k] = v
		}
	}

	argSlice := []string{}
	if *arglist != "" {
		// naive split by spaces; recommended to use proper quoting in production
		argSlice = append(argSlice, (*arglist))
	}

	if err := ramguard.RunWithSecrets(ctx, *path, argSlice, secrets); err != nil {
		fmt.Fprintln(os.Stderr, "ramguard error:", err)
		os.Exit(1)
	}
}

func procfreezeCmd(ctx context.Context, args []string) {
	fs := flag.NewFlagSet("procfreeze", flag.ExitOnError)
	manifest := fs.String("manifest", "manifest.json", "path to workspace state manifest JSON")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, "parse:", err)
		os.Exit(2)
	}

	// remaining args: action
	rem := fs.Args()
	if len(rem) < 1 {
		fmt.Fprintln(os.Stderr, "expected action: freeze|unfreeze")
		os.Exit(2)
	}

	mgr := procfreeze.NewManager()
	if err := mgr.LoadManifest(*manifest); err != nil {
		fmt.Fprintln(os.Stderr, "load manifest:", err)
		os.Exit(1)
	}

	switch rem[0] {
	case "freeze":
		if err := mgr.Freeze(ctx); err != nil {
			fmt.Fprintln(os.Stderr, "freeze error:", err)
			os.Exit(1)
		}
	case "unfreeze":
		if err := mgr.Unfreeze(ctx); err != nil {
			fmt.Fprintln(os.Stderr, "unfreeze error:", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "unknown action:", rem[0])
		os.Exit(2)
	}
}

func execgateCmd(ctx context.Context, args []string) {
	fs := flag.NewFlagSet("execgate", flag.ExitOnError)
	cmdStr := fs.String("cmd", "", "command string to inspect and (maybe) run")
	force := fs.Bool("force-intent", false, "bypass interactive confirmation for high-risk commands")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, "parse:", err)
		os.Exit(2)
	}
	if *cmdStr == "" {
		fmt.Fprintln(os.Stderr, "cmd is required")
		fs.Usage()
		os.Exit(2)
	}

	executed, err := execgate.InspectAndMaybeRun(ctx, *cmdStr, *force)
	if err != nil {
		fmt.Fprintln(os.Stderr, "execgate error:", err)
		os.Exit(1)
	}
	if !executed {
		fmt.Println("not executed")
	}
}