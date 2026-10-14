package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Rome2112/substrate/pkg/execgate"
	"github.com/Rome2112/substrate/pkg/procfreeze"
	"github.com/Rome2112/substrate/pkg/ramguard"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	switch strings.TrimSpace(os.Args[1]) {
	case "ramguard":
		ramguardCmd(ctx, os.Args[2:])
	case "procfreeze":
		procfreezeCmd(ctx, os.Args[2:])
	case "execgate":
		execgateCmd(ctx, os.Args[2:])
	case "-h", "--help", "help":
		printUsage()
	case "version":
		fmt.Println("substrate version 0.1.0")
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand: %s\n", os.Args[1])
		printUsage()
		os.Exit(2)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage:
  substrate <subcommand> [flags]

Subcommands:
  ramguard    Run a child process with in-memory secret injection
  procfreeze  Freeze or unfreeze processes from a workspace manifest
  execgate    Check command safety and execute with confirmation

Examples:
  substrate ramguard --path /bin/echo --secret "TOKEN=abc" --args "hello"
  substrate procfreeze --manifest ./examples/workspace_manifest.json freeze
  substrate execgate --cmd "rm -rf ./dist" --force-intent
`)
}

func ramguardCmd(ctx context.Context, args []string) {
	fs := flag.NewFlagSet("ramguard", flag.ContinueOnError)
	path := fs.String("path", "", "path to executable")
	secret := fs.String("secret", "", "KEY=VALUE environment variable to inject")
	argString := fs.String("args", "", "space-delimited arguments for the child process")
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	if *path == "" {
		fmt.Fprintln(os.Stderr, "--path is required")
		os.Exit(2)
	}
	if *secret == "" {
		fmt.Fprintln(os.Stderr, "--secret is required")
		os.Exit(2)
	}

	secrets, err := parseSecrets(*secret)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	if err := ramguard.RunWithSecrets(ctx, *path, splitArgs(*argString), secrets); err != nil {
		fmt.Fprintln(os.Stderr, "ramguard error:", err)
		os.Exit(1)
	}
}

func procfreezeCmd(ctx context.Context, args []string) {
	fs := flag.NewFlagSet("procfreeze", flag.ContinueOnError)
	manifestPath := fs.String("manifest", "workspace_manifest.json", "path to manifest JSON")
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	remaining := fs.Args()
	if len(remaining) != 1 {
		fmt.Fprintln(os.Stderr, "expected action: freeze|unfreeze")
		os.Exit(2)
	}

	mgr := procfreeze.NewManager()
	if err := mgr.LoadManifest(*manifestPath); err != nil {
		fmt.Fprintln(os.Stderr, "manifest error:", err)
		os.Exit(1)
	}

	switch remaining[0] {
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
		fmt.Fprintf(os.Stderr, "unknown action %q\n", remaining[0])
		os.Exit(2)
	}
}

func execgateCmd(ctx context.Context, args []string) {
	fs := flag.NewFlagSet("execgate", flag.ContinueOnError)
	cmdString := fs.String("cmd", "", "command to inspect and maybe execute")
	forceIntent := fs.Bool("force-intent", false, "skip confirmation prompt")
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	if *cmdString == "" {
		fmt.Fprintln(os.Stderr, "--cmd is required")
		os.Exit(2)
	}

	ok, err := execgate.InspectAndMaybeRun(ctx, *cmdString, *forceIntent)
	if err != nil {
		if errors.Is(err, execgate.ErrSafetyAbort) {
			fmt.Fprintln(os.Stderr, "execgate aborted by safety policy")
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, "execgate error:", err)
		os.Exit(1)
	}
	if !ok {
		fmt.Fprintln(os.Stdout, "execution skipped")
		os.Exit(0)
	}
}

func parseSecrets(raw string) (map[string]string, error) {
	secrets := map[string]string{}
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 || strings.TrimSpace(kv[0]) == "" {
			return nil, fmt.Errorf("invalid secret %q: expected KEY=VALUE", part)
		}
		secrets[kv[0]] = kv[1]
	}
	return secrets, nil
}

func splitArgs(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	return strings.Fields(raw)
}
