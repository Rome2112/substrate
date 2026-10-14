# Examples for Substrate

This directory contains example artifacts to help test Substrate locally.

- `workspace_manifest.json` — sample manifest for procfreeze (replace pid values with real PIDs)
- `run_ramguard.sh` — small example script demonstrating ramguard usage

Usage:
1. Build the CLI:
   go build -o substrate ./cmd/substrate

2. Run the ramguard example:
   chmod +x run_ramguard.sh
   ./run_ramguard.sh

3. Test procfreeze by replacing example pids with live process IDs and running:
   ./substrate procfreeze --manifest workspace_manifest.json freeze
   ./substrate procfreeze --manifest workspace_manifest.json unfreeze