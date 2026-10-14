# Substrate

Substrate is a lightweight, zero-dependency, native-local invariant execution suite written in Go.
It provides process-level boundary isolation, ephemeral memory security for runtime secrets, and
command safety validation on Linux, macOS and Windows.

This repository contains three small core tools packaged as subcommands:

- ramguard — ephemeral environment injection for secrets without writing plaintext to disk.
- procfreeze — lightweight process grouping and suspend/resume primitives.
- execgate — a safety filter that intercepts high-risk shell commands and displays a dependency "blast radius".

Design goals:
- Zero heavy dependencies; only Go standard library.
- Small, auditable binary and low latency.
- Cross-platform support (Linux, macOS, Windows) with proper build tags and graceful fallbacks on unsupported OS primitives.
- Clear error types and non-panicking behavior.

Quick build
-----------
From repository root:

go build -o substrate ./cmd/substrate

This produces a single `substrate` binary containing three subcommands.

Usage
-----
Run the binary with the subcommand name:

./substrate ramguard --path /usr/bin/printenv --secret SECRET=foo
./substrate procfreeze --manifest workspace_manifest.json freeze
./substrate procfreeze --manifest workspace_manifest.json unfreeze
./substrate execgate --cmd "rm -rf ./dist"           # will prompt for confirmation
./substrate execgate --cmd "rm -rf ./dist" --force-intent  # bypass prompt

Project layout
--------------
- /cmd/substrate/main.go         - CLI entrypoint with three subcommands.
- /pkg/ramguard/                 - RamGuard package (ephemeral env injection).
- /pkg/procfreeze/               - ProcFreeze package (process grouping and freeze/unfreeze).
- /pkg/execgate/                 - ExecGate package (command filtering and blast-radius scan).
- /go.mod
- /README.md

Security notes
--------------
- ramguard attempts to minimize secret lifetime in memory by zeroing owned buffers after child process termination. It cannot (and will not attempt to) modify OS or other process internals to zero memory beyond the process it spawns, and is not a cryptographic secret manager.
- procfreeze uses native OS primitives where available; on platforms without native APIs implemented, it will return descriptive errors.
- execgate is a best-effort local safety filter; it does not replace CI/CD protections or sandboxing for production systems.

Testing
-------
Each package contains unit tests. Run tests with:

go test ./...

Contributions
-------------
This repository is intentionally small and dependency-free. If you extend features that require more complex OS interactions, prefer optional, well-audited packages and keep the core small.

License
-------
(Choose your license and add here.)