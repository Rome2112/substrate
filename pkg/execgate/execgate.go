// Package execgate inspects command strings for high-risk signatures and
// optionally executes them after providing a "blast radius" summary.
// The package is intentionally small and uses only the standard library.
package execgate

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ErrSafetyAbort indicates that the user aborted execution.
var ErrSafetyAbort = errors.New("execution aborted by safety gate")

// riskSignatures is a list of high-risk substrings to detect.
var riskSignatures = []string{
	"rm -rf",
	"rm -r",
	"DROP TABLE",
	"DROP DATABASE",
	"migrate --force",
	"migrate --unsafe",
	"PGPASSWORD=",
	"postgres://",
	"mysql://",
}

// InspectAndMaybeRun inspects cmdStr for risky signatures. If a risk is found,
// a directory scan for dependency manifests is printed to stdout, and the function
// will prompt the user for explicit confirmation unless force is true.
// If execution proceeds, the underlying shell will run the command and this
// function returns (true, nil). If the function refuses to run, returns (false, nil).
func InspectAndMaybeRun(ctx context.Context, cmdStr string, force bool) (bool, error) {
	lc := strings.ToLower(cmdStr)
	isRisk := false
	for _, sig := range riskSignatures {
		if strings.Contains(strings.ToLower(sig), "drop ") {
			// case-sensitive SQL keywords are relevant; check case-insensitively.
		}
		if strings.Contains(lc, strings.ToLower(sig)) {
			isRisk = true
			break
		}
	}

	if isRisk {
		fmt.Fprintln(os.Stdout, "High-risk command detected. Computing blast radius...")
		deps := scanDependencyManifests(".")
		printBlastRadiusSummary(deps)
		if !force {
			if ok := promptForConfirmation(); !ok {
				return false, ErrSafetyAbort
			}
		}
	}

	// Execute the command using the system shell
	var exe *exec.Cmd
	if isWindows() {
		exe = exec.CommandContext(ctx, "cmd.exe", "/C", cmdStr)
	} else {
		exe = exec.CommandContext(ctx, "sh", "-c", cmdStr)
	}
	exe.Stdin = os.Stdin
	exe.Stdout = os.Stdout
	exe.Stderr = os.Stderr

	if err := exe.Run(); err != nil {
		return true, fmt.Errorf("underlying command failed: %w", err)
	}
	return true, nil
}

// scanDependencyManifests walks starting at root and discovers files that are common dependency and repo indicators.
func scanDependencyManifests(root string) []string {
	found := []string{}
	// small list of filenames to look for
	targets := map[string]struct{}{
		"go.mod":      {},
		"package.json": {},
		"requirements.txt": {},
		"Pipfile":        {},
		"mix.exs":        {},
		".git":           {},
	}

	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			// ignore access errors
			return nil
		}
		name := d.Name()
		if _, ok := targets[name]; ok {
			found = append(found, path)
			// if .git directory, don't descend into it
			if name == ".git" && d.IsDir() {
				return filepath.SkipDir
			}
		}
		return nil
	})
	return found
}

func printBlastRadiusSummary(paths []string) {
	if len(paths) == 0 {
		fmt.Fprintln(os.Stdout, "No dependency manifests (.git, go.mod, package.json etc.) found in current tree.")
		return
	}
	fmt.Fprintln(os.Stdout, "Affected dependency manifests:")
	for _, p := range paths {
		fmt.Fprintf(os.Stdout, " - %s\n", p)
	}
}

// promptForConfirmation prints a [y/N] prompt and returns true only for a leading 'y' or 'Y' input.
func promptForConfirmation() bool {
	fmt.Fprint(os.Stdout, "Do you intend to continue? [y/N]: ")
	r := bufio.NewReader(os.Stdin)
	resp, err := r.ReadString('\n')
	if err != nil {
		return false
	}
	resp = strings.TrimSpace(resp)
	if resp == "y" || resp == "Y" {
		return true
	}
	return false
}

func isWindows() bool {
	return strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") || (os.PathSeparator == '\\')
}