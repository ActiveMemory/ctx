//! ctx_adapter — the single module that invokes the `ctx` CLI.
//!
//! Every read goes through `ctx` so the `.context/` files stay the
//! source of truth. Commands resolve their context dir from the
//! process working directory (`$PWD/.context`), so each call runs
//! `ctx` with `current_dir` set to the selected project root.
//!
//! We use std::process::Command directly (no tauri-plugin-shell), so
//! no shell capability/permission wiring is required.
//!
//! Every command is async and pushes its process spawns / fs walks
//! onto the blocking pool via [`blocking`], so a slow `ctx` (or a
//! huge tree) never stalls the UI thread.

use std::io::Read;
use std::process::{Command, Stdio};
use std::sync::{Mutex, OnceLock};

use serde::Serialize;

/// Name of the ctx binary; resolved from PATH unless overridden via
/// [`set_ctx_path`].
const CTX_BIN: &str = "ctx";

/// User-set absolute path to the ctx binary (for non-standard installs).
/// `None` means "resolve `ctx` from PATH".
fn ctx_bin_cell() -> &'static Mutex<Option<String>> {
    static CELL: OnceLock<Mutex<Option<String>>> = OnceLock::new();
    CELL.get_or_init(|| Mutex::new(None))
}

/// The ctx executable to spawn: the user override if set, else `ctx`.
fn ctx_bin() -> String {
    ctx_bin_cell()
        .lock()
        .ok()
        .and_then(|g| g.clone())
        .unwrap_or_else(|| CTX_BIN.to_string())
}

/// PATH prepended with the dirs a GUI launch (minimal launchd PATH on
/// macOS) needs to find a user-installed ctx — Homebrew, MacPorts, and
/// common Go / local-bin install locations.
fn search_path() -> String {
    let mut extra = vec![
        "/usr/local/bin".to_string(),
        "/opt/homebrew/bin".to_string(),
        "/opt/local/bin".to_string(),
    ];
    if let Ok(home) = std::env::var("HOME") {
        if !home.is_empty() {
            extra.push(format!("{home}/go/bin"));
            extra.push(format!("{home}/.local/bin"));
        }
    }
    let prefix = extra.join(":");
    match std::env::var("PATH") {
        Ok(existing) => format!("{prefix}:{existing}"),
        Err(_) => prefix,
    }
}

/// Runs `f` on the blocking pool — the bridge that keeps process
/// spawns and fs walks off the async runtime the commands run on.
async fn blocking<T, F>(f: F) -> Result<T, String>
where
    T: Send + 'static,
    F: FnOnce() -> Result<T, String> + Send + 'static,
{
    tauri::async_runtime::spawn_blocking(f)
        .await
        .map_err(|e| format!("background task failed: {e}"))?
}

/// Overrides the ctx binary path (persisted by the frontend after a
/// successful call). An empty string clears the override, returning to
/// PATH resolution. A non-empty path is validated first: it must be an
/// existing file whose `--version` output mentions ctx, so a stale or
/// arbitrary saved path can never silently become the spawn target.
#[tauri::command]
pub async fn set_ctx_path(path: String) -> Result<(), String> {
    blocking(move || {
        let trimmed = path.trim().to_string();
        if trimmed.is_empty() {
            if let Ok(mut cell) = ctx_bin_cell().lock() {
                *cell = None;
            }
            return Ok(());
        }
        if !std::path::Path::new(&trimmed).is_file() {
            return Err(format!("not a file: {trimmed}"));
        }
        let out = run_bin(&trimmed, "", &["--version"])?;
        if !out.to_lowercase().contains("ctx") {
            return Err(format!(
                "`{trimmed}` does not look like the ctx binary (--version said: {})",
                out.trim()
            ));
        }
        if let Ok(mut cell) = ctx_bin_cell().lock() {
            *cell = Some(trimmed);
        }
        Ok(())
    })
    .await
}

/// True when `<dir>/.context` exists — lets the UI validate a restored
/// active project before trusting it (a moved/deleted project should
/// fall back to the chooser, not error on every screen).
#[tauri::command]
pub async fn dir_is_ctx_project(dir: String) -> bool {
    blocking(move || {
        Ok(!dir.is_empty() && std::path::Path::new(&dir).join(".context").is_dir())
    })
    .await
    .unwrap_or(false)
}

/// Reports whether the ctx binary is available and its version.
#[derive(Serialize)]
pub struct CtxInfo {
    pub found: bool,
    pub version: String,
    pub error: Option<String>,
}

/// Hard ceiling on a single `ctx` invocation. The CLI is local and
/// fast; anything past this is a hang, so we surface a timeout error
/// rather than block the UI command forever.
const CTX_TIMEOUT: std::time::Duration = std::time::Duration::from_secs(30);

/// How often we poll a running child for completion.
const POLL_INTERVAL: std::time::Duration = std::time::Duration::from_millis(25);

/// Runs `bin <args>` in `dir` (empty dir = inherit cwd) and returns
/// stdout on success, or a trimmed stderr/error string on failure.
///
/// The child is spawned with piped output and polled via `try_wait`;
/// on the [`CTX_TIMEOUT`] deadline it is killed and reaped, so a
/// wedged process can neither block the caller nor linger as an
/// orphan. The pipes are drained on threads (so a chatty child can't
/// deadlock on a full pipe); killing the child closes its pipes,
/// which lets those threads terminate.
fn run_bin(bin: &str, dir: &str, args: &[&str]) -> Result<String, String> {
    let mut cmd = Command::new(bin);
    cmd.args(args);
    if !dir.is_empty() {
        cmd.current_dir(dir);
    }
    cmd.env("PATH", search_path());
    cmd.stdin(Stdio::null());
    cmd.stdout(Stdio::piped());
    cmd.stderr(Stdio::piped());

    let mut child = cmd
        .spawn()
        .map_err(|e| format!("could not run `{bin}`: {e}"))?;

    let mut stdout_pipe = child.stdout.take();
    let mut stderr_pipe = child.stderr.take();
    let out_thread = std::thread::spawn(move || {
        let mut buf = Vec::new();
        if let Some(p) = stdout_pipe.as_mut() {
            let _ = p.read_to_end(&mut buf);
        }
        buf
    });
    let err_thread = std::thread::spawn(move || {
        let mut buf = Vec::new();
        if let Some(p) = stderr_pipe.as_mut() {
            let _ = p.read_to_end(&mut buf);
        }
        buf
    });

    let deadline = std::time::Instant::now() + CTX_TIMEOUT;
    let status = loop {
        match child.try_wait() {
            Ok(Some(status)) => break status,
            Ok(None) => {
                if std::time::Instant::now() >= deadline {
                    let _ = child.kill();
                    let _ = child.wait(); // reap — no zombie
                    let _ = out_thread.join();
                    let _ = err_thread.join();
                    return Err(format!(
                        "`{bin} {}` timed out after {}s",
                        args.join(" "),
                        CTX_TIMEOUT.as_secs()
                    ));
                }
                std::thread::sleep(POLL_INTERVAL);
            }
            Err(e) => {
                let _ = child.kill();
                let _ = child.wait();
                let _ = out_thread.join();
                let _ = err_thread.join();
                return Err(format!("could not wait on `{bin}`: {e}"));
            }
        }
    };

    let stdout = out_thread.join().unwrap_or_default();
    let stderr = err_thread.join().unwrap_or_default();

    if !status.success() {
        let stderr = String::from_utf8_lossy(&stderr).trim().to_string();
        if stderr.is_empty() {
            return Err(format!("`{bin} {}` exited with an error", args.join(" ")));
        }
        return Err(stderr);
    }
    Ok(String::from_utf8_lossy(&stdout).to_string())
}

/// Runs `ctx <args>` in `dir` — see [`run_bin`] for the execution and
/// timeout contract.
fn run_ctx(dir: &str, args: &[&str]) -> Result<String, String> {
    run_bin(&ctx_bin(), dir, args)
}

/// [`run_ctx`] for owned argument vectors (the write paths build their
/// args dynamically).
fn run_ctx_owned(dir: &str, args: &[String]) -> Result<String, String> {
    let refs: Vec<&str> = args.iter().map(|s| s.as_str()).collect();
    run_ctx(dir, &refs)
}

/// Detects the ctx binary and returns its version string.
#[tauri::command]
pub async fn ctx_info() -> CtxInfo {
    let res = blocking(move || run_ctx("", &["--version"])).await;
    match res {
        Ok(out) => CtxInfo {
            found: true,
            version: out.trim().to_string(),
            error: None,
        },
        Err(e) => CtxInfo {
            found: false,
            version: String::new(),
            error: Some(e),
        },
    }
}

/// Which optional CLI contracts the installed ctx supports.
#[derive(Serialize)]
pub struct CtxCapabilities {
    /// True when `ctx task list --json` exists (the structured-list
    /// contract the Tasks/Decisions/Learnings/Search screens need).
    pub list_json: bool,
}

/// Probes the installed ctx for the `list --json` contract by actually
/// running `ctx task list --json` in `dir` once. Only an unknown
/// command/flag, a spawn failure, or a timeout means "unsupported" —
/// any other error (no .context, parse trouble, …) proves the
/// subcommand exists. We deliberately do NOT probe via `--help`:
/// cobra prints help and exits 0 for unknown subcommands, so a help
/// probe would falsely report support.
#[tauri::command]
pub async fn ctx_capabilities(dir: String) -> CtxCapabilities {
    let list_json = blocking(move || {
        Ok(match run_ctx(&dir, &["task", "list", "--json"]) {
            Ok(_) => true,
            Err(e) => {
                let msg = e.to_lowercase();
                !(msg.contains("unknown command")
                    || msg.contains("unknown flag")
                    || msg.contains("unknown shorthand flag")
                    || msg.contains("could not run")
                    || msg.contains("timed out"))
            }
        })
    })
    .await
    .unwrap_or(false);
    CtxCapabilities { list_json }
}

/// Returns `ctx status --json` for the project at `dir`.
#[tauri::command]
pub async fn ctx_status(dir: String) -> Result<String, String> {
    blocking(move || run_ctx(&dir, &["status", "--json"])).await
}

/// Returns `ctx doctor --json` for the project at `dir`.
#[tauri::command]
pub async fn ctx_doctor(dir: String) -> Result<String, String> {
    blocking(move || run_ctx(&dir, &["doctor", "--json"])).await
}

/// Returns `ctx task list --json` for the project at `dir`.
#[tauri::command]
pub async fn ctx_task_list(dir: String) -> Result<String, String> {
    blocking(move || run_ctx(&dir, &["task", "list", "--json"])).await
}

/// Returns `ctx decision list --json` for the project at `dir`.
#[tauri::command]
pub async fn ctx_decision_list(dir: String) -> Result<String, String> {
    blocking(move || run_ctx(&dir, &["decision", "list", "--json"])).await
}

/// Returns `ctx learning list --json` for the project at `dir`.
#[tauri::command]
pub async fn ctx_learning_list(dir: String) -> Result<String, String> {
    blocking(move || run_ctx(&dir, &["learning", "list", "--json"])).await
}

/// Returns `ctx journal source --limit N` for `dir`. The CLI has
/// no JSON mode for the journal yet, so this is the raw table
/// text; the UI renders it verbatim rather than risk a fragile
/// column parse.
#[tauri::command]
pub async fn ctx_journal(dir: String, limit: u32) -> Result<String, String> {
    blocking(move || {
        let l = limit.to_string();
        run_ctx(&dir, &["journal", "source", "--limit", l.as_str()])
    })
    .await
}

/// Returns `ctx journal source --show <session>` for `dir` — the raw
/// text of a single session, used by the multi-project dashboard's
/// per-task drill-down to show the journal entry a task's `session`
/// id points at. Rendered verbatim (no JSON mode for the journal yet).
#[tauri::command]
pub async fn ctx_journal_show(dir: String, session: String) -> Result<String, String> {
    if session.trim().is_empty() {
        return Err("session id is empty".to_string());
    }
    blocking(move || run_ctx(&dir, &["journal", "source", "--show", session.as_str()])).await
}

/// Runs `ctx drift` (or `ctx drift --fix`) for `dir` and returns
/// the report. With `fix`, auto-corrects supported issues
/// (staleness, missing files).
#[tauri::command]
pub async fn ctx_drift(dir: String, fix: bool) -> Result<String, String> {
    blocking(move || {
        if fix {
            run_ctx(&dir, &["drift", "--fix"])
        } else {
            run_ctx(&dir, &["drift"])
        }
    })
    .await
}

/// Runs `ctx compact` (or `ctx compact --archive`) for `dir`.
/// With `archive`, moves completed/old content into
/// `.context/archive/`. Mutating — the UI confirms first.
#[tauri::command]
pub async fn ctx_compact(dir: String, archive: bool) -> Result<String, String> {
    blocking(move || {
        if archive {
            run_ctx(&dir, &["compact", "--archive"])
        } else {
            run_ctx(&dir, &["compact"])
        }
    })
    .await
}

/// Canonical context files this app may read directly. The
/// allowlist is the path-traversal guard: `name` must match one of
/// these exactly, so no caller-supplied path ever reaches the
/// filesystem.
const READABLE_DOCS: &[&str] = &["CONSTITUTION.md", "CONVENTIONS.md"];

/// Reads a canonical `.context/<name>.md` file verbatim for a
/// read-only viewer.
///
/// This is the one read that does not funnel through `ctx`: there is
/// no `ctx` command that returns a single file's full content (the
/// agent packet is budget-trimmed and omits the constitution), and
/// these files ARE the source of truth, so reading them directly is
/// both accurate and safe. `name` is allowlisted (see
/// [`READABLE_DOCS`]). A missing file returns an empty string so the
/// screen can show a friendly "not present" state instead of an
/// error.
#[tauri::command]
pub async fn ctx_read_doc(dir: String, name: String) -> Result<String, String> {
    if !READABLE_DOCS.contains(&name.as_str()) {
        return Err(format!("`{name}` is not a readable context document"));
    }
    blocking(move || {
        let path = std::path::Path::new(&dir).join(".context").join(&name);
        match std::fs::read_to_string(&path) {
            Ok(content) => Ok(content),
            Err(e) if e.kind() == std::io::ErrorKind::NotFound => Ok(String::new()),
            Err(e) => Err(format!("could not read {}: {e}", path.display())),
        }
    })
    .await
}

/// What the KB browser needs to render its left rail: whether a kb
/// exists, which top-level docs are present, and the topic slugs.
#[derive(Serialize)]
pub struct KbInfo {
    pub exists: bool,
    pub docs: Vec<String>,
    pub topics: Vec<String>,
}

/// Top-level `.context/kb/*.md` files the browser surfaces, in
/// display order. Only those that exist are returned.
const KB_DOCS: &[&str] = &[
    "index.md",
    "evidence-index.md",
    "source-map.md",
    "source-coverage.md",
    "grounding-sources.md",
    "outstanding-questions.md",
];

/// Hard bound on topic-tree recursion: a symlink cycle (or an absurdly
/// deep tree) must not walk forever.
const MAX_TOPIC_DEPTH: usize = 16;

/// Recursively collects every directory under `root` that contains
/// an `index.md`, as a slash-joined slug relative to `root`. Stops
/// descending past [`MAX_TOPIC_DEPTH`] levels.
fn collect_topics(
    root: &std::path::Path,
    cur: &std::path::Path,
    depth: usize,
    out: &mut Vec<String>,
) {
    if depth > MAX_TOPIC_DEPTH {
        return;
    }
    let entries = match std::fs::read_dir(cur) {
        Ok(e) => e,
        Err(_) => return,
    };
    let mut has_index = false;
    let mut subdirs = Vec::new();
    for entry in entries.flatten() {
        let p = entry.path();
        if p.is_dir() {
            subdirs.push(p);
        } else if p.file_name().is_some_and(|n| n == "index.md") {
            has_index = true;
        }
    }
    if has_index && cur != root {
        if let Ok(rel) = cur.strip_prefix(root) {
            out.push(rel.to_string_lossy().replace('\\', "/"));
        }
    }
    for sub in subdirs {
        collect_topics(root, &sub, depth + 1, out);
    }
}

/// Inventories the `.context/kb/` of the project at `dir`.
#[tauri::command]
pub async fn kb_info(dir: String) -> KbInfo {
    blocking(move || {
        let kb = std::path::Path::new(&dir).join(".context").join("kb");
        if !kb.is_dir() {
            return Ok(KbInfo {
                exists: false,
                docs: vec![],
                topics: vec![],
            });
        }
        let docs = KB_DOCS
            .iter()
            .filter(|f| kb.join(f).is_file())
            .map(|s| s.to_string())
            .collect();
        let mut topics = Vec::new();
        let topics_root = kb.join("topics");
        collect_topics(&topics_root, &topics_root, 0, &mut topics);
        topics.sort();
        Ok(KbInfo {
            exists: true,
            docs,
            topics,
        })
    })
    .await
    .unwrap_or(KbInfo {
        exists: false,
        docs: vec![],
        topics: vec![],
    })
}

/// Reads a file under `.context/kb/` by its kb-relative path.
///
/// `rel` is validated segment-by-segment: no empty, `.`, `..`, or
/// backslash segments are allowed, which prevents both absolute
/// paths and traversal out of the kb. The resolved path is then
/// canonicalized and prefix-checked against the canonical kb root,
/// so a symlink inside the kb cannot escape it either. A missing
/// file returns "".
#[tauri::command]
pub async fn kb_read(dir: String, rel: String) -> Result<String, String> {
    if rel.is_empty() {
        return Err("empty kb path".to_string());
    }
    for seg in rel.split('/') {
        if seg.is_empty() || seg == "." || seg == ".." || seg.contains('\\') {
            return Err(format!("invalid kb path: {rel}"));
        }
    }
    blocking(move || {
        let kb_root = std::path::Path::new(&dir).join(".context").join("kb");
        let canon_root = match kb_root.canonicalize() {
            Ok(p) => p,
            // No kb at all — same friendly empty state as a missing file.
            Err(_) => return Ok(String::new()),
        };
        let path = kb_root.join(&rel);
        let real = match path.canonicalize() {
            Ok(p) => p,
            Err(e) if e.kind() == std::io::ErrorKind::NotFound => return Ok(String::new()),
            Err(e) => return Err(format!("could not read {}: {e}", path.display())),
        };
        if !real.starts_with(&canon_root) {
            return Err(format!("invalid kb path: {rel}"));
        }
        match std::fs::read_to_string(&real) {
            Ok(content) => Ok(content),
            Err(e) if e.kind() == std::io::ErrorKind::NotFound => Ok(String::new()),
            Err(e) => Err(format!("could not read {}: {e}", real.display())),
        }
    })
    .await
}

/// Returns `ctx agent --format json --budget N` for `dir` — the
/// structured context packet used by the budget-preview screen.
#[tauri::command]
pub async fn ctx_agent_json(dir: String, budget: u32) -> Result<String, String> {
    blocking(move || {
        let b = budget.to_string();
        run_ctx(&dir, &["agent", "--format", "json", "--budget", b.as_str()])
    })
    .await
}

/// Returns `ctx agent --budget N` markdown for `dir` — the
/// paste-ready packet for the "copy packet" action.
#[tauri::command]
pub async fn ctx_agent_md(dir: String, budget: u32) -> Result<String, String> {
    blocking(move || {
        let b = budget.to_string();
        run_ctx(&dir, &["agent", "--budget", b.as_str()])
    })
    .await
}

/// `ctx remind list` — raw text list of pending reminders.
#[tauri::command]
pub async fn ctx_remind_list(dir: String) -> Result<String, String> {
    blocking(move || run_ctx(&dir, &["remind", "list"])).await
}

/// `ctx remind add <text>` — adds a session reminder. The text is
/// passed after `--` so user input can never be parsed as a flag.
#[tauri::command]
pub async fn ctx_remind_add(dir: String, text: String) -> Result<String, String> {
    if text.trim().is_empty() {
        return Err("reminder text is empty".to_string());
    }
    blocking(move || run_ctx(&dir, &["remind", "add", "--", text.as_str()])).await
}

/// `ctx remind dismiss <target>` — `target` is a number or "all",
/// passed after `--` so it can never be parsed as a flag.
#[tauri::command]
pub async fn ctx_remind_dismiss(dir: String, target: String) -> Result<String, String> {
    if target.trim().is_empty() {
        return Err("dismiss target is empty".to_string());
    }
    blocking(move || run_ctx(&dir, &["remind", "dismiss", "--", target.as_str()])).await
}

/// `ctx pad` — raw text list of encrypted scratchpad entries
/// (decrypted for display by the CLI).
#[tauri::command]
pub async fn ctx_pad_list(dir: String) -> Result<String, String> {
    blocking(move || run_ctx(&dir, &["pad"])).await
}

/// `ctx pad add <text>` — appends a scratchpad entry. The text is
/// passed after `--` so user input can never be parsed as a flag.
#[tauri::command]
pub async fn ctx_pad_add(dir: String, text: String) -> Result<String, String> {
    if text.trim().is_empty() {
        return Err("entry text is empty".to_string());
    }
    blocking(move || run_ctx(&dir, &["pad", "add", "--", text.as_str()])).await
}

/// `ctx pad rm <n>` — removes the scratchpad entry numbered `n`.
#[tauri::command]
pub async fn ctx_pad_rm(dir: String, n: String) -> Result<String, String> {
    if n.trim().is_empty() {
        return Err("entry number is empty".to_string());
    }
    blocking(move || run_ctx(&dir, &["pad", "rm", "--", n.as_str()])).await
}

/// `ctx pad show <n>` — raw text of a single scratchpad entry.
#[tauri::command]
pub async fn ctx_pad_show(dir: String, n: String) -> Result<String, String> {
    if n.trim().is_empty() {
        return Err("entry number is empty".to_string());
    }
    blocking(move || run_ctx(&dir, &["pad", "show", "--", n.as_str()])).await
}

/// `ctx connection status` — hub connection status. Errors when the
/// project has no hub configured; the UI surfaces that as a
/// "not connected" state.
#[tauri::command]
pub async fn ctx_connection_status(dir: String) -> Result<String, String> {
    blocking(move || run_ctx(&dir, &["connection", "status"])).await
}

/// Reads a single trimmed git value from the repo at `dir`, or "".
fn git_field(dir: &str, args: &[&str]) -> String {
    Command::new("git")
        .arg("-C")
        .arg(dir)
        .args(args)
        .output()
        .ok()
        .filter(|o| o.status.success())
        .map(|o| String::from_utf8_lossy(&o.stdout).trim().to_string())
        .unwrap_or_default()
}

/// One 8-char session id per app launch, for write provenance. Minted
/// once (not per call) so every write from this desktop session shares
/// a coherent id in the journal, and rapid writes can't collide. The
/// CLI truncates to 8 chars anyway; derived from the wall clock to
/// avoid a rand dependency.
fn session_id() -> String {
    static SID: OnceLock<String> = OnceLock::new();
    SID.get_or_init(|| {
        let nanos = std::time::SystemTime::now()
            .duration_since(std::time::UNIX_EPOCH)
            .map(|d| d.as_nanos())
            .unwrap_or(0);
        format!("{:08x}", (nanos as u64) & 0xffff_ffff)
    })
    .clone()
}

/// Provenance flags for the write commands: the session id always,
/// `--branch`/`--commit` only when git actually reports a value — a
/// project without git (or with no commits yet) must omit the flags
/// entirely rather than send empty values the CLI rejects.
fn provenance_args(dir: &str) -> Vec<String> {
    let mut args = vec!["--session-id".to_string(), session_id()];
    let branch = git_field(dir, &["rev-parse", "--abbrev-ref", "HEAD"]);
    if !branch.is_empty() {
        args.push("--branch".to_string());
        args.push(branch);
    }
    let commit = git_field(dir, &["rev-parse", "--short", "HEAD"]);
    if !commit.is_empty() {
        args.push("--commit".to_string());
        args.push(commit);
    }
    args
}

/// Adds a task via `ctx task add`, synthesizing the required
/// provenance flags (session id minted here, branch/commit read
/// from git; see [`provenance_args`]). Empty `priority`/`section`
/// are omitted. Flags come first, then `--`, then the user text, so
/// the text can never be parsed as a flag.
#[tauri::command]
pub async fn ctx_task_add(
    dir: String,
    text: String,
    priority: String,
    section: String,
) -> Result<String, String> {
    if text.trim().is_empty() {
        return Err("task text is empty".to_string());
    }
    blocking(move || {
        let mut args: Vec<String> = vec!["task".to_string(), "add".to_string()];
        args.extend(provenance_args(&dir));
        if !priority.is_empty() {
            args.push("--priority".to_string());
            args.push(priority);
        }
        if !section.is_empty() {
            args.push("--section".to_string());
            args.push(section);
        }
        args.push("--".to_string());
        args.push(text);
        run_ctx_owned(&dir, &args)
    })
    .await
}

/// Marks a task complete via `ctx task complete <id-or-text>`. The
/// target goes after `--` so it can never be parsed as a flag.
#[tauri::command]
pub async fn ctx_task_complete(dir: String, target: String) -> Result<String, String> {
    if target.trim().is_empty() {
        return Err("task identifier is empty".to_string());
    }
    blocking(move || run_ctx(&dir, &["task", "complete", "--", target.as_str()])).await
}

/// Adds an ADR-style decision via `ctx decision add`, synthesizing
/// provenance (see [`provenance_args`]). All three ADR fields are
/// required by the CLI, so we reject empties up front with a clear
/// message. Flags first, then `--`, then the title.
#[tauri::command]
pub async fn ctx_decision_add(
    dir: String,
    title: String,
    context: String,
    rationale: String,
    consequence: String,
) -> Result<String, String> {
    if title.trim().is_empty()
        || context.trim().is_empty()
        || rationale.trim().is_empty()
        || consequence.trim().is_empty()
    {
        return Err("title, context, rationale and consequence are all required".to_string());
    }
    blocking(move || {
        let mut args: Vec<String> = vec![
            "decision".to_string(),
            "add".to_string(),
            "--context".to_string(),
            context,
            "--rationale".to_string(),
            rationale,
            "--consequence".to_string(),
            consequence,
        ];
        args.extend(provenance_args(&dir));
        args.push("--".to_string());
        args.push(title);
        run_ctx_owned(&dir, &args)
    })
    .await
}

/// Adds a learning via `ctx learning add`, synthesizing provenance
/// (see [`provenance_args`]). Context, lesson, and application are
/// all required by the CLI. Flags first, then `--`, then the title.
#[tauri::command]
pub async fn ctx_learning_add(
    dir: String,
    title: String,
    context: String,
    lesson: String,
    application: String,
) -> Result<String, String> {
    if title.trim().is_empty()
        || context.trim().is_empty()
        || lesson.trim().is_empty()
        || application.trim().is_empty()
    {
        return Err("title, context, lesson and application are all required".to_string());
    }
    blocking(move || {
        let mut args: Vec<String> = vec![
            "learning".to_string(),
            "add".to_string(),
            "--context".to_string(),
            context,
            "--lesson".to_string(),
            lesson,
            "--application".to_string(),
            application,
        ];
        args.extend(provenance_args(&dir));
        args.push("--".to_string());
        args.push(title);
        run_ctx_owned(&dir, &args)
    })
    .await
}
