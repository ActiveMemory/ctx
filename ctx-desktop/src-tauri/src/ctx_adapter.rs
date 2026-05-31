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
