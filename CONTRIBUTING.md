# Contributing to kubetray

Thanks for your interest in contributing! The project is in early stages and there's plenty of room to help.

## Getting Started

```bash
git clone https://github.com/ugurozkn/kubetray.git
cd kubetray
make build    # Build binary
make test     # Run tests
```

### Prerequisites

- Go 1.24+
- A container runtime (Colima, Docker Desktop, or OrbStack)
- helm, kubectl

## How to Contribute

1. Check [open issues](https://github.com/ugurozkn/kubetray/issues) for something to work on
2. Fork the repo
3. Create a branch: `git checkout -b my-feature`
4. Make your changes
5. Run tests: `make test`
6. Commit with a clear message
7. Open a pull request

## Project Structure

```
kubetray/
├── cmd/           # CLI commands (start, stop, clean, mcp)
├── pkg/
│   ├── config/    # Config and state management
│   ├── k8s/       # k3d cluster operations
│   ├── mcp/       # MCP server for AI integration
│   ├── platform/  # OS detection, dependency checking
│   └── ui/        # Terminal output (spinners, tables, colors)
├── npm/           # npm package wrapper
├── main.go
├── Makefile
└── .goreleaser.yaml
```

## Guidelines

- **Keep it simple** — no over-engineering, only add what's needed
- **One feature per PR** — small, focused pull requests are easier to review
- **Test your changes** — run `make test` and manually verify with `./kubetray`
- **Error messages should help** — suggest what the user can do to fix the problem
- **Follow existing patterns** — look at how current commands are structured in `cmd/`

## Adding a New Command

1. Create `cmd/yourcommand.go`
2. Define the cobra command and register it with `rootCmd.AddCommand()` in `init()`
3. Use `pkg/ui` for terminal output (spinners, colors, tables)
4. Use `pkg/config` to load config and `pkg/k8s` for cluster operations

## Reporting Bugs

Open an [issue](https://github.com/ugurozkn/kubetray/issues/new) with:
- What you expected to happen
- What actually happened
- Your OS and architecture (`kubetray version`)
- Steps to reproduce

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
