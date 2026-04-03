# Detailed Design Index

Per-module reference documentation, split by domain. Each file
contains full module sections (purpose, types, API, data flow,
edge cases, danger zones, extension points).

| Domain | File | Modules | Summary |
|--------|------|---------|---------|
| Foundation | [DETAILED_DESIGN-foundation.md](DETAILED_DESIGN-foundation.md) | config/*, assets/*, io, format, parse, sanitize, validate, inspect, flagbind, exec/*, log/*, crypto, sysinfo, rc | Constants, I/O, formatting, external commands |
| Domain | [DETAILED_DESIGN-domain.md](DETAILED_DESIGN-domain.md) | entity, entry, context/*, drift, index, task, tidy, trace, journal/*, memory, notify, claude | Core business logic and data types |
| MCP | [DETAILED_DESIGN-mcp.md](DETAILED_DESIGN-mcp.md) | mcp/proto, mcp/handler, mcp/server/*, mcp/session | JSON-RPC 2.0 server, tools, resources, prompts |
| CLI | [DETAILED_DESIGN-cli.md](DETAILED_DESIGN-cli.md) | bootstrap, cli/parent, 34 command packages | Command registration, groups, taxonomy |
| Output | [DETAILED_DESIGN-output.md](DETAILED_DESIGN-output.md) | write/*, err/*, assets/read/* | Terminal formatting, error constructors, text lookup |

> See individual files for module-level detail.
