# Plan: `ctx system` — System Resource Monitoring

## Context

OOM event (~17GB used on 16GB system) made the box unresponsive and prevented
ctx from persisting session data. A `ctx system` command surfaces resource
pressure at two severity tiers (WARNING, DANGER) so users can act before it's
too late. A companion hook (`check-resources`) proactively warns during sessions
when resources hit DANGER level.

## Approach

**Repurpose the existing hidden `system` parent command.** Currently it's
`Hidden: true` with no `RunE`, acting only as a namespace for hook subcommands.
We un-hide it, add a `RunE` that shows resource stats, and keep all hook
subcommands individually hidden. Net effect: `ctx system` appears in help and
shows stats; `ctx system check-context-size` etc. remain invisible.

**New `internal/sysinfo/` package** for OS-level resource gathering with build
tags (Linux primary, macOS secondary, Windows graceful fallback). This is the
first use of build tags in the project — the right Go pattern for
platform-specific code.

## Thresholds

| Resource | WARNING | DANGER |
|----------|---------|--------|
| Memory | >= 80% used | >= 90% used |
| Swap | >= 50% used | >= 75% used |
| Disk (cwd) | >= 85% full | >= 95% full |
| Load (1m) | >= 0.8x CPUs | >= 1.5x CPUs |

## Files to Create

### `internal/sysinfo/` — Resource gathering (new package)

| File | Purpose |
|------|---------|
| `doc.go` | Package documentation |
| `resources.go` | Types (`MemInfo`, `DiskInfo`, `LoadInfo`, `Snapshot`, `ResourceAlert`, `Severity`), `Collect()` entry point, `MaxSeverity()` |
| `threshold.go` | `evaluate()` — checks snapshot against thresholds, returns alerts; `formatGiB()` helper |
| `memory_linux.go` | `//go:build linux` — parses `/proc/meminfo` via `parseMeminfo(io.Reader)` |
| `memory_darwin.go` | `//go:build darwin` — `sysctl hw.memsize` + `vm_stat` + `sysctl vm.swapusage` |
| `memory_other.go` | `//go:build !linux && !darwin` — returns `MemInfo{Supported: false}` |
| `load_linux.go` | `//go:build linux` — parses `/proc/loadavg` via `parseLoadavg(io.Reader)` |
| `load_darwin.go` | `//go:build darwin` — `sysctl` for load averages |
| `load_other.go` | `//go:build !linux && !darwin` — stub |
| `disk.go` | `syscall.Statfs` — works on Linux + macOS without build tags |
| `disk_windows.go` | `//go:build windows` — stub (Statfs unavailable on Windows) |
| `resources_test.go` | Tests for `evaluate()`, `MaxSeverity()`, `Collect()` with constructed data |
| `threshold_test.go` | Boundary tests (exactly at 80%, 79.9%, 90.1%, etc.) |
| `parse_linux_test.go` | `//go:build linux` — tests `parseMeminfo()` and `parseLoadavg()` with fake input |

### `internal/cli/system/` — Command + hook (modified package)

| File | Purpose |
|------|---------|
| `resources.go` | NEW: `runResources()`, `outputResourcesText()`, `outputResourcesJSON()` |
| `check_resources.go` | NEW: Hidden hook subcommand — VERBATIM relay on DANGER only |
| `resources_test.go` | NEW: Output formatting tests with constructed snapshots |
| `check_resources_test.go` | NEW: Hook output tests |

## Files to Modify

| File | Change |
|------|--------|
| `internal/cli/system/system.go` | Un-hide, add `RunE` + `--json` flag, register `checkResourcesCmd()` |

## Output UX

**`ctx system` (all OK):**
```
System Resources
====================

Memory:    4.2 / 16.0 GB (26%)                     ✓ ok
Swap:      0.0 /  8.0 GB (0%)                      ✓ ok
Disk:    180.2 / 500.0 GB (36%)                     ✓ ok
Load:     0.52 / 0.41 / 0.38  (8 CPUs, ratio 0.07) ✓ ok

All clear — no resource warnings.
```

**`ctx system` (DANGER):**
```
System Resources
====================

Memory:   14.7 / 16.0 GB (92%)                     ✖ DANGER
Swap:      6.2 /  8.0 GB (78%)                     ✖ DANGER
Disk:    180.2 / 500.0 GB (36%)                     ✓ ok
Load:    12.50 / 9.30 / 6.10  (8 CPUs, ratio 1.56) ✖ DANGER

Alerts:
  ✖ Memory 92% used (14.7 / 16.0 GB)
  ✖ Swap 78% used (6.2 / 8.0 GB)
  ✖ Load 1.56x CPU count
```

**Hook output (DANGER only, VERBATIM relay):**
```
IMPORTANT: Relay this resource warning to the user VERBATIM.

┌─ Resource Alert ──────────────────────────────────
│ ✖ Memory 92% used (14.7 / 16.0 GB)
│ ✖ Swap 78% used (6.2 / 8.0 GB)
│
│ System resources are critically low.
│ Persist unsaved context NOW with /ctx-wrap-up
│ and consider ending this session.
└───────────────────────────────────────────────────
```

## Testing strategy

- **Threshold/evaluation logic**: Table-driven tests with constructed `Snapshot` values — no OS access
- **Parsing logic**: Extract `parseMeminfo(io.Reader)` and `parseLoadavg(io.Reader)` — test with fake `/proc` content
- **Output formatting**: Construct snapshots, capture `cmd.SetOut(&buf)` output, assert strings
- **Cross-compilation**: `./hack/build-all.sh dev` must pass for all 6 targets

## Implementation order

1. `internal/sysinfo/` package — types, threshold logic, platform collectors, tests
2. `internal/cli/system/resources.go` — user-facing output formatting
3. `internal/cli/system/system.go` — un-hide, add RunE
4. `internal/cli/system/check_resources.go` — hook subcommand
5. Cross-compilation check + manual test on dev machine

## Verification

```bash
CGO_ENABLED=0 go test ./internal/sysinfo/...
CGO_ENABLED=0 go test ./internal/cli/system/...
CGO_ENABLED=0 go test ./...
./hack/build-all.sh dev          # all 6 targets compile
./ctx system                     # manual check
./ctx system --json | jq .       # valid JSON
```
