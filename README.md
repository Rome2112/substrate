# Substrate

Substrate is a lightweight, zero-dependency, native-local invariant execution suite written in Go.
It provides lightweight process isolation, ephemeral memory-resident secrets, and command safety
validation for local developer workflows across Linux, macOS, and Windows.

The repository includes three primary tools:

- ramguard: ephemeral injection of secret environment variables without writing plaintext files
- procfreeze: suspend/resume workspace processes using native OS process primitives and optional cgroup v2 support
- execgate: detect high-risk commands, display a dependency blast radius, and require confirmation before execution

Project layout
--------------

- `cmd/substrate/main.go` — CLI entrypoint for the three subcommands
- `pkg/ramguard/` — secret injection and secure execution
- `pkg/procfreeze/` — process grouping, freeze/unfreeze, and manifest-driven workspace control
- `pkg/execgate/` — risky command scanning and blast-radius analysis
- `examples/` — sample manifests and usage scripts
- `.github/workflows/ci.yml` — CI workflow covering Go tests on Ubuntu, macOS, and Windows

Build
-----

From the repository root:

    go build -o substrate ./cmd/substrate

Usage
-----

    ./substrate ramguard --path /usr/bin/env --secret "API_TOKEN=s3cr3t" --args ""
    ./substrate procfreeze --manifest ./examples/workspace_manifest.json freeze
    ./substrate procfreeze --manifest ./examples/workspace_manifest.json unfreeze
    ./substrate execgate --cmd "rm -rf ./dist" --force-intent

Security model
--------------

- `ramguard` injects variables directly into the child process environment without writing them to `.env` files.
- It attempts to zero the memory slice that owns the secret values immediately after the child exits.
- `procfreeze` uses native OS primitives for suspend/resume and supports optional Linux cgroup v2 grouping.
- `execgate` blocks unsafe command execution unless the caller provides confirmation or passes `--force-intent`.

Testing
-------

Run the full test suite:

    go test ./...

License
-------

MIT
