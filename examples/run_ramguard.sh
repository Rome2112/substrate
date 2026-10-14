#!/usr/bin/env bash
set -euo pipefail

# Example usage: run a harmless env print with a secret injected in-memory.
# Make sure ./substrate is built first: go build -o substrate ./cmd/substrate
./substrate ramguard --path /usr/bin/env --secret "API_TOKEN=demo_secret" --args ""