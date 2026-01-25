![ctx](assets/ctx-banner.png)

## Contributing to Context

Thank you for your interest in contributing to `ctx`.

This document provides guidelines and information for contributors.

## Ways to Contribute

- **Report bugs**: Open an issue describing the problem
- **Suggest features**: Open an issue with your idea
- **Submit fixes**: Fork, fix, and submit a pull request
- **Improve docs**: Documentation improvements are always welcome
- **Share feedback**: Let us know how you use ctx

## Development Setup

### Prerequisites

- [Go 1.25+](https://go.dev/)
- Git

### Building

```bash
# Clone the repository
git clone https://github.com/ActiveMemory/ctx.git
cd ctx

# Build
CGO_ENABLED=0 go build -o ctx ./cmd/ctx

# Run tests
CGO_ENABLED=0 go test ./...

# Install locally
sudo mv ctx /usr/local/bin/
```

### Project Structure

```
ctx/
├── cmd/ctx/           # CLI entry point
├── internal/
│   ├── cli/           # Command implementations
│   ├── context/       # Core context logic
│   ├── drift/         # Drift detection
│   ├── claude/        # Claude Code integration
│   └── templates/     # Embedded templates
├── templates/         # Template source files
├── specs/             # Feature specifications
└── .context/          # ctx's own context (dogfooding)
```

## Submitting Changes

### Before You Start

1. Check existing issues to avoid duplicating work
2. For large changes, open an issue first to discuss the approach
3. Read the specs in `specs/` to understand the design

### Pull Request Process

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes
4. Ensure tests pass: `go test ./...`
5. Ensure code compiles: `go build ./...`
6. Commit with a clear message (see below)
7. Push and open a pull request

### Commit Messages

Follow conventional commits:

```
type(scope): description

[optional body]

[optional footer]
```

Types: `feat`, `fix`, `docs`, `test`, `refactor`, `chore`

Examples:
- `feat(cli): add ctx export command`
- `fix(drift): handle missing files gracefully`
- `docs: update installation instructions`

### Code Style

- Follow Go conventions (`gofmt`, `go vet`)
- Keep functions focused and small
- Add tests for new functionality
- Use meaningful variable names
- Handle errors explicitly

## Developer Certificate of Origin (DCO)

By contributing, you agree to the 
[Developer Certificate of Origin](CONTRIBUTING_DCO.md).

All commits must be signed off:

```bash
git commit -s -m "feat: add new feature"
```

This adds a `Signed-off-by` line to your commit, certifying that you wrote the
code or have the right to submit it under the project's license.

## Code of Conduct

This project follows the [Contributor Covenant Code of 
Conduct](CODE_OF_CONDUCT.md).

Please read it before participating.

## Getting Help

* Open an issue for bugs or feature requests
* Check existing issues and discussions
* Read the specs in `specs/` for design context

## License

By contributing, you agree that your contributions will be licensed under the
[Apache 2.0 License](LICENSE).

