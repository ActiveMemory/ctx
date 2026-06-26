# ctx Desktop

A cross-platform desktop client for [`ctx`](https://github.com/ActiveMemory/ctx) —
a calm, local-first window into your project's persistent AI context. It is a
**thin client over the `ctx` CLI**: every read and write shells out to `ctx`, so
the `.context/` files stay the single source of truth.

Stack: **Tauri 2 + React + TypeScript + Tailwind v4** (Vite).

## Prerequisites

| Tool | Version | Notes |
|------|---------|-------|
| **Node.js** | 18+ (built with 22) | https://nodejs.org |
| **Rust** | **1.88+** (built with 1.96) | Tauri 2's deps require ≥1.88. `rustup update stable` |
| **`ctx` CLI** | 0.8.1+ | Must be on your `PATH`. https://ctx.ist |
| **System deps** | per-OS | See Tauri's prerequisites (below) |

Install the per-OS native dependencies (WebView, build tools) by following the
official guide: https://tauri.app/start/prerequisites/

- **macOS:** Xcode Command Line Tools (`xcode-select --install`); WebKit is built in.
- **Linux:** `webkit2gtk`, `libgtk`, `librsvg`, `build-essential` (see the link).
- **Windows:** Microsoft C++ Build Tools + WebView2 (preinstalled on Win 11).

## Setup & run

```bash
git clone https://github.com/ActiveMemory/ctx.git
cd ctx/ctx-desktop
git checkout feat/ctx-desktop      # until merged to main

npm install                        # JS deps
npm run tauri dev                  # build + launch the app (hot-reloads)
```

The first `tauri dev` compiles the Rust backend (a few minutes); later runs are
fast. A native **ctx Desktop** window opens. Click **Add workspace…**, pick a
folder that contains your ctx projects, and switch between them from the dropdown.

## Production build

```bash
npm run tauri build
```

Produces signed-installer artifacts under `src-tauri/target/release/bundle/`
(`.dmg`/`.app` on macOS, `.deb`/`.AppImage` on Linux, `.msi`/`.exe` on Windows).

## The `ctx` CLI dependency

The app calls `ctx` for everything. Most screens work with **released `ctx`
0.8.1**, but the **list/count views** (Overview counts, Tasks/Decisions/Learnings
lists) require the `ctx <artifact> list --json` commands, which currently live on
the `feat/ctx-artifact-list-json` branch (not yet in a release). To enable them,
build and install that branch once:

```bash
git checkout feat/ctx-artifact-list-json
make build && sudo make install    # installs ctx to /usr/local/bin
git checkout feat/ctx-desktop
```

Adding entries, the Context Packet, Journal, and Health screens all work on stock
`ctx` 0.8.1 without this step.

> macOS note: a bundled app launched from Finder inherits a minimal `PATH`, so the
> adapter prepends `/usr/local/bin` and `/opt/homebrew/bin` to find a
> user-installed `ctx`. In `npm run tauri dev` (launched from a terminal) the
> inherited `PATH` already works.

## Architecture

```
ctx-desktop/
├── src/                      React + TS + Tailwind frontend
│   ├── adapter/ctx.ts        typed invoke() wrappers + CLI JSON types (one file)
│   ├── hooks/useReload.ts    debounced "ctx-changed" → reload key
│   ├── lib/markdown.tsx      minimal Markdown → React renderer (no innerHTML)
│   ├── screens/              Projects, Overview, Search, Tasks, Reminders,
│   │                         Decisions, Learnings, Conventions, Constitution
│   │                         (CanonicalDoc), ContextPacket, KnowledgeBase,
│   │                         Scratchpad (Pad), Journal, Drift, Health, Hub
│   └── App.tsx               nav shell, workspace switcher, top bar
└── src-tauri/                Rust host
    └── src/
        ├── ctx_adapter.rs    THE module that spawns `ctx` (reads + writes)
        ├── discover.rs       workspace scan for .context projects
        └── watcher.rs        fs-watch on .context → emits "ctx-changed"
```

- All `ctx` access funnels through `ctx_adapter.rs` (Rust) and `adapter/ctx.ts`
  (TS), so a CLI/output change is a one- or two-file fix.
- Writes synthesize provenance (`--session-id`, and `--branch`/`--commit` from
  git) in the Rust adapter.
- The app spawns `ctx` via `std::process::Command` directly — no
  `tauri-plugin-shell`, so no shell capability wiring.

## Useful scripts

| Command | What |
|---------|------|
| `npm run tauri dev` | Build + launch with hot reload |
| `npm run tauri build` | Production bundles |
| `npm run build` | Type-check + build the frontend only (`tsc && vite build`) |
| `cargo fmt && cargo clippy` (in `src-tauri/`) | Format + lint Rust |
| `cargo test` (in `src-tauri/`) | Run Rust unit tests |

## Privacy

Local-only. No network calls except an optional, explicit update check. Nothing
outside the selected project's `.context/` is read or transmitted; every change
routes through `ctx` to preserve its audit trail.
