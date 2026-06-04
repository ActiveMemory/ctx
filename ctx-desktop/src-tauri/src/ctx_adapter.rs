//! ctx_adapter — the single module that invokes the `ctx` CLI.
//!
//! Every read goes through `ctx` so the `.context/` files stay the
//! source of truth. Commands resolve their context dir from the
//! process working directory (`$PWD/.context`), so each call runs
//! `ctx` with `current_dir` set to the selected project root.
//!
//! We use std::process::Command directly (no tauri-plugin-shell), so
//! no shell capability/permission wiring is required.

use std::process::Command;

use serde::Serialize;

/// Name of the ctx binary; resolved from PATH.
const CTX_BIN: &str = "ctx";

/// Directories prepended to PATH so a GUI launch (which inherits a
/// minimal launchd PATH on macOS) can still find a user-installed ctx.
const EXTRA_PATH: &str = "/usr/local/bin:/opt/homebrew/bin";

/// Reports whether the ctx binary is available and its version.
#[derive(Serialize)]
pub struct CtxInfo {
    pub found: bool,
    pub version: String,
    pub error: Option<String>,
}

/// Runs `ctx <args>` in `dir` (empty dir = inherit cwd) and returns
/// stdout on success, or a trimmed stderr/error string on failure.
fn run_ctx(dir: &str, args: &[&str]) -> Result<String, String> {
    let mut cmd = Command::new(CTX_BIN);
    cmd.args(args);
    if !dir.is_empty() {
        cmd.current_dir(dir);
    }
    let path = match std::env::var("PATH") {
        Ok(existing) => format!("{EXTRA_PATH}:{existing}"),
        Err(_) => EXTRA_PATH.to_string(),
    };
    cmd.env("PATH", path);

    let output = cmd
        .output()
        .map_err(|e| format!("could not run `{CTX_BIN}`: {e}"))?;
    if !output.status.success() {
        let stderr = String::from_utf8_lossy(&output.stderr).trim().to_string();
        if stderr.is_empty() {
            return Err(format!("`ctx {}` exited with an error", args.join(" ")));
        }
        return Err(stderr);
    }
    Ok(String::from_utf8_lossy(&output.stdout).to_string())
}

/// Detects the ctx binary and returns its version string.
#[tauri::command]
pub fn ctx_info() -> CtxInfo {
    match run_ctx("", &["--version"]) {
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

/// Returns `ctx status --json` for the project at `dir`.
#[tauri::command]
pub fn ctx_status(dir: String) -> Result<String, String> {
    run_ctx(&dir, &["status", "--json"])
}

/// Returns `ctx doctor --json` for the project at `dir`.
#[tauri::command]
pub fn ctx_doctor(dir: String) -> Result<String, String> {
    run_ctx(&dir, &["doctor", "--json"])
}

/// Returns `ctx task list --json` for the project at `dir`.
#[tauri::command]
pub fn ctx_task_list(dir: String) -> Result<String, String> {
    run_ctx(&dir, &["task", "list", "--json"])
}

/// Returns `ctx decision list --json` for the project at `dir`.
#[tauri::command]
pub fn ctx_decision_list(dir: String) -> Result<String, String> {
    run_ctx(&dir, &["decision", "list", "--json"])
}

/// Returns `ctx learning list --json` for the project at `dir`.
#[tauri::command]
pub fn ctx_learning_list(dir: String) -> Result<String, String> {
    run_ctx(&dir, &["learning", "list", "--json"])
}

/// Returns `ctx journal source --limit N` for `dir`. The CLI has
/// no JSON mode for the journal yet, so this is the raw table
/// text; the UI renders it verbatim rather than risk a fragile
/// column parse.
#[tauri::command]
pub fn ctx_journal(dir: String, limit: u32) -> Result<String, String> {
    let l = limit.to_string();
    run_ctx(&dir, &["journal", "source", "--limit", l.as_str()])
}

/// Runs `ctx drift` (or `ctx drift --fix`) for `dir` and returns
/// the report. With `fix`, auto-corrects supported issues
/// (staleness, missing files).
#[tauri::command]
pub fn ctx_drift(dir: String, fix: bool) -> Result<String, String> {
    if fix {
        run_ctx(&dir, &["drift", "--fix"])
    } else {
        run_ctx(&dir, &["drift"])
    }
}

/// Runs `ctx compact` (or `ctx compact --archive`) for `dir`.
/// With `archive`, moves completed/old content into
/// `.context/archive/`. Mutating — the UI confirms first.
#[tauri::command]
pub fn ctx_compact(dir: String, archive: bool) -> Result<String, String> {
    if archive {
        run_ctx(&dir, &["compact", "--archive"])
    } else {
        run_ctx(&dir, &["compact"])
    }
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
pub fn ctx_read_doc(dir: String, name: String) -> Result<String, String> {
    if !READABLE_DOCS.contains(&name.as_str()) {
        return Err(format!("`{name}` is not a readable context document"));
    }
    let path = std::path::Path::new(&dir).join(".context").join(&name);
    match std::fs::read_to_string(&path) {
        Ok(content) => Ok(content),
        Err(e) if e.kind() == std::io::ErrorKind::NotFound => Ok(String::new()),
        Err(e) => Err(format!("could not read {}: {e}", path.display())),
    }
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

/// Recursively collects every directory under `root` that contains
/// an `index.md`, as a slash-joined slug relative to `root`.
fn collect_topics(root: &std::path::Path, cur: &std::path::Path, out: &mut Vec<String>) {
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
        collect_topics(root, &sub, out);
    }
}

/// Inventories the `.context/kb/` of the project at `dir`.
#[tauri::command]
pub fn kb_info(dir: String) -> KbInfo {
    let kb = std::path::Path::new(&dir).join(".context").join("kb");
    if !kb.is_dir() {
        return KbInfo {
            exists: false,
            docs: vec![],
            topics: vec![],
        };
    }
    let docs = KB_DOCS
        .iter()
        .filter(|f| kb.join(f).is_file())
        .map(|s| s.to_string())
        .collect();
    let mut topics = Vec::new();
    let topics_root = kb.join("topics");
    collect_topics(&topics_root, &topics_root, &mut topics);
    topics.sort();
    KbInfo {
        exists: true,
        docs,
        topics,
    }
}

/// Reads a file under `.context/kb/` by its kb-relative path.
///
/// `rel` is validated segment-by-segment: no empty, `.`, `..`, or
/// backslash segments are allowed, which prevents both absolute
/// paths and traversal out of the kb. A missing file returns "".
#[tauri::command]
pub fn kb_read(dir: String, rel: String) -> Result<String, String> {
    if rel.is_empty() {
        return Err("empty kb path".to_string());
    }
    for seg in rel.split('/') {
        if seg.is_empty() || seg == "." || seg == ".." || seg.contains('\\') {
            return Err(format!("invalid kb path: {rel}"));
        }
    }
    let path = std::path::Path::new(&dir)
        .join(".context")
        .join("kb")
        .join(&rel);
    match std::fs::read_to_string(&path) {
        Ok(content) => Ok(content),
        Err(e) if e.kind() == std::io::ErrorKind::NotFound => Ok(String::new()),
        Err(e) => Err(format!("could not read {}: {e}", path.display())),
    }
}

/// Returns `ctx agent --format json --budget N` for `dir` — the
/// structured context packet used by the budget-preview screen.
#[tauri::command]
pub fn ctx_agent_json(dir: String, budget: u32) -> Result<String, String> {
    let b = budget.to_string();
    run_ctx(&dir, &["agent", "--format", "json", "--budget", b.as_str()])
}

/// Returns `ctx agent --budget N` markdown for `dir` — the
/// paste-ready packet for the "copy packet" action.
#[tauri::command]
pub fn ctx_agent_md(dir: String, budget: u32) -> Result<String, String> {
    let b = budget.to_string();
    run_ctx(&dir, &["agent", "--budget", b.as_str()])
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

/// Generates an 8-char session id for write provenance. The CLI
/// truncates to 8 chars anyway; we derive from the wall clock so
/// it varies per call without pulling in a rand dependency.
fn session_id() -> String {
    let nanos = std::time::SystemTime::now()
        .duration_since(std::time::UNIX_EPOCH)
        .map(|d| d.as_nanos())
        .unwrap_or(0);
    format!("{:08x}", (nanos as u64) & 0xffff_ffff)
}

/// Adds a task via `ctx task add`, synthesizing the required
/// provenance flags (session id minted here, branch/commit read
/// from git). Empty `priority`/`section` are omitted.
#[tauri::command]
pub fn ctx_task_add(
    dir: String,
    text: String,
    priority: String,
    section: String,
) -> Result<String, String> {
    if text.trim().is_empty() {
        return Err("task text is empty".to_string());
    }
    let branch = git_field(&dir, &["rev-parse", "--abbrev-ref", "HEAD"]);
    let commit = git_field(&dir, &["rev-parse", "--short", "HEAD"]);
    let session = session_id();

    let mut args: Vec<&str> = vec![
        "task",
        "add",
        text.as_str(),
        "--session-id",
        session.as_str(),
        "--branch",
        branch.as_str(),
        "--commit",
        commit.as_str(),
    ];
    if !priority.is_empty() {
        args.push("--priority");
        args.push(priority.as_str());
    }
    if !section.is_empty() {
        args.push("--section");
        args.push(section.as_str());
    }
    run_ctx(&dir, &args)
}

/// Marks a task complete via `ctx task complete <id-or-text>`.
#[tauri::command]
pub fn ctx_task_complete(dir: String, target: String) -> Result<String, String> {
    if target.trim().is_empty() {
        return Err("task identifier is empty".to_string());
    }
    run_ctx(&dir, &["task", "complete", target.as_str()])
}

/// Adds an ADR-style decision via `ctx decision add`, synthesizing
/// provenance. All three ADR fields are required by the CLI, so we
/// reject empties up front with a clear message.
#[tauri::command]
pub fn ctx_decision_add(
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
    let branch = git_field(&dir, &["rev-parse", "--abbrev-ref", "HEAD"]);
    let commit = git_field(&dir, &["rev-parse", "--short", "HEAD"]);
    let session = session_id();

    run_ctx(
        &dir,
        &[
            "decision",
            "add",
            title.as_str(),
            "--context",
            context.as_str(),
            "--rationale",
            rationale.as_str(),
            "--consequence",
            consequence.as_str(),
            "--session-id",
            session.as_str(),
            "--branch",
            branch.as_str(),
            "--commit",
            commit.as_str(),
        ],
    )
}

/// Adds a learning via `ctx learning add`, synthesizing provenance.
/// Context, lesson, and application are all required by the CLI.
#[tauri::command]
pub fn ctx_learning_add(
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
    let branch = git_field(&dir, &["rev-parse", "--abbrev-ref", "HEAD"]);
    let commit = git_field(&dir, &["rev-parse", "--short", "HEAD"]);
    let session = session_id();

    run_ctx(
        &dir,
        &[
            "learning",
            "add",
            title.as_str(),
            "--context",
            context.as_str(),
            "--lesson",
            lesson.as_str(),
            "--application",
            application.as_str(),
            "--session-id",
            session.as_str(),
            "--branch",
            branch.as_str(),
            "--commit",
            commit.as_str(),
        ],
    )
}
