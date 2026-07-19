# environment-and-platform

## [2026-07-04-162854] GitNexus prunes registry entries whose repo paths don't resolve in-container

**Context**: Folding the Docker gitnexus wrapper into ctx: indexing ctx silently dropped os from the global registry, and a plain 'gitnexus list' passthrough dropped it again after it was restored. Both 'gitnexus index' and 'gitnexus list' rewrite registry.json minus any entry whose path they cannot stat — and inside a container an unmounted repo is indistinguishable from a deleted one.

**Lesson**: Any registry-touching gitnexus invocation run in Docker must mount EVERY registered repo at its real host path, not just the current one; otherwise it persists a pruned registry. A bash footgun compounded it: a '[ ... ] && cmd' filter as the last command of a while-loop body returns 1 when the filter fails on the final line, and set -e kills the script — use an if-block inside loops under set -e.

**Application**: hack/gitnexus-docker.sh now has registry_mounts() (mount all registered repos, exclude the separately-mounted current one) wired into the index-register and passthrough branches; mcp already did this. The os and orchestrator copies of the wrapper still carry the prune bug — backport before running their index/list targets.

---

## [2026-06-07-170017] Host-pressure alerting: use derivatives, not levels (consolidated)

**Consolidated from**: 2 entries (2026-04-13 to 2026-05-28)

- Swap occupancy is NOT memory pressure: macOS/Windows swap proactively and
  occupancy is a sticky high-water mark that doesn't recede when pressure ends,
  so any alert keyed on SwapUsed/SwapTotal ≥ X% false-positives at session
  start (e.g. after hibernation). Key on OS-native pressure derivatives instead:
  macOS kern.memorystatus_vm_pressure_level (1/2/4 → OK/Warning/Danger), Linux
  PSI /proc/pressure/memory some.avg10/full.avg10; fall back to swap-out RATE
  gated on low available memory, never occupancy.
- Load average measures a queue (runnable + uninterruptible-sleep), not CPU
  utilization — high load with low CPU% means many short-lived/I/O-bound
  processes (e.g. go test spawning hundreds of binaries). For automated alerts
  prefer the 5-minute average over the reactive 1-minute, which fires on normal
  build/test activity.

---

## [2026-05-31-094649] macOS GUI apps inherit a minimal PATH; augment it to find a user-installed CLI

**Context**: A bundled Tauri app launched via Finder/launchd gets a minimal PATH (/usr/bin:/bin:...), so /usr/local/bin/ctx is not found even though it resolves in a terminal-launched dev run.

**Lesson**: Do not rely on inherited PATH when spawning user-installed CLIs from a desktop GUI.

**Application**: ctx_adapter prepends /usr/local/bin:/opt/homebrew/bin to PATH on every std::process::Command invocation.

---

## [2026-05-31-094649] Tauri 2 requires rustc >= 1.88; bump the toolchain before cargo check

**Context**: cargo check for the ctx-desktop Tauri app failed: darling, serde_with, time and plist transitive deps require rustc 1.88, but the local toolchain was 1.87.

**Lesson**: Tauri 2's dependency tree tracks recent rustc releases; the pinned-stable assumption breaks builds.

**Application**: Run rustup update stable before building a Tauri 2 app; this project moved 1.87 -> 1.96.

---

## [2026-05-20-214839] macOS /var symlink trips path-equality; use EvalSymlinks with parent-resolution fallback

**Context**: TestRunInit_EnvCwdMatch_Succeeds in
internal/cli/initialize/init_test.go failed on first run despite a deliberate
setup where the env path and cwd candidate matched. Diagnosis: t.TempDir()
returns paths like /var/folders/..., os.Getwd() after t.Chdir() returns the
canonical /private/var/folders/... (because macOS's /var is a symlink to
/private/var). filepath.Clean preserves the symlink form; equality fails.

**Lesson**: filepath.Clean alone is insufficient for path equality on macOS (and
other systems with symlinked top-level dirs). filepath.EvalSymlinks resolves the
symlinks but fails when the target path does not yet exist — common case for
/Users/volkan/Desktop/WORKSPACE/ctx/.context BEFORE ctx init runs. The right
pattern is a layered fallback: try EvalSymlinks(full), then EvalSymlinks(parent)
+ rejoin basename, then filepath.Clean as last resort.

**Application**: Encapsulated as
internal/cli/initialize/core/envmatch/{envmatch.go,internal.go}. The Same(a, b)
public function calls resolve() on each side; resolve() tries EvalSymlinks on
the full path, falls back to EvalSymlinks on the parent (rejoining the
basename), and falls through to filepath.Clean if both fail. Reusable for any
future env-vs-cwd-style equality check. The package is per-feature
(core/envmatch/) per the cmd/core/ purity rule enforced by
internal/compliance/TestCmdDirPurity.

---

## [2026-02-26-100006] PATH and binary handling (consolidated)

**Consolidated from**: 3 entries (2026-01-21 to 2026-02-17)

- Always use `ctx` from PATH, never `./dist/ctx-linux-arm64` or `go run
  ./cmd/ctx`. Check `which ctx` if unsure.
- Hooks must use PATH, not hardcoded paths. `ctx init` checks if ctx is in PATH
  before proceeding. Tests can skip with `CTX_SKIP_PATH_CHECK=1`.
- Agent must never place binaries in any bin directory (not via cp, mv, or go
  install). Build with `make build`, then ask the user to run the privileged
  install step. Hooks in block-dangerous-commands.sh enforce this.

---

