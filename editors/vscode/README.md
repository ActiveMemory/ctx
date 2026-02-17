# ctx — VS Code Chat Extension

A VS Code Chat Participant that brings [ctx](https://ctx.ist) — persistent project context for AI coding sessions — directly into GitHub Copilot Chat.

## Usage

Type `@ctx` in the VS Code Chat view, then use slash commands:

| Command | Description |
|---------|-------------|
| `@ctx /init` | Initialize a `.context/` directory with template files |
| `@ctx /status` | Show context summary with token estimate |
| `@ctx /agent` | Print AI-ready context packet |
| `@ctx /drift` | Detect stale or invalid context |
| `@ctx /recall` | Browse and search AI session history |
| `@ctx /hook` | Generate AI tool integration configs |
| `@ctx /add` | Add a task, decision, or learning |
| `@ctx /load` | Output assembled context Markdown |
| `@ctx /compact` | Archive completed tasks and clean up |
| `@ctx /sync` | Reconcile context with codebase |

## Prerequisites

- [ctx](https://ctx.ist) CLI installed and available on PATH (or configure `ctx.executablePath`)
- VS Code 1.93+ with GitHub Copilot Chat

## Configuration

| Setting | Default | Description |
|---------|---------|-------------|
| `ctx.executablePath` | `ctx` | Path to the ctx executable |

## Development

```bash
cd editors/vscode
npm install
npm run watch   # Watch mode
npm run build   # Production build
```

## License

Apache-2.0
