//! File watcher — emits a "ctx-changed" event whenever the active
//! project's `.context/` directory mutates, so the frontend can
//! refetch. The GUI is one writer among several (human CLI + AI
//! agents); this keeps its views from going stale.
//!
//! Both channels carry the project root as the event payload, so
//! listeners can tell WHICH project changed and refresh only that
//! one instead of blanking and refetching everything.

use std::path::Path;
use std::sync::Mutex;

use notify::{RecommendedWatcher, RecursiveMode, Watcher};
use tauri::{AppHandle, Emitter, State};

/// Holds the live watcher so it stays alive between calls. Replacing
/// the value drops (and stops) the previous watcher.
#[derive(Default)]
pub struct WatchState(pub Mutex<Option<RecommendedWatcher>>);

/// Holds one watcher per discovered project so the multi-project
/// dashboard reflects external writes to *any* project, not just the
/// active one. Replacing the vec drops the previous watchers.
#[derive(Default)]
pub struct WorkspaceWatchState(pub Mutex<Vec<RecommendedWatcher>>);

/// Builds a watcher on `<dir>/.context` that emits `event` with the
/// project root (`dir`) as payload on every filesystem event there.
/// Returns Ok(None) when `.context/` is missing — not an error, the
/// watch simply does not start.
fn build_watcher(
    app: &AppHandle,
    dir: &str,
    event: &'static str,
) -> Result<Option<RecommendedWatcher>, String> {
    let ctx_dir = Path::new(dir).join(".context");
    if !ctx_dir.is_dir() {
        return Ok(None);
    }
    let handle = app.clone();
    let root = dir.to_string();
    let mut watcher = notify::recommended_watcher(move |res: notify::Result<notify::Event>| {
        if res.is_ok() {
            let _ = handle.emit(event, root.clone());
        }
    })
    .map_err(|e| e.to_string())?;
    watcher
        .watch(&ctx_dir, RecursiveMode::Recursive)
        .map_err(|e| e.to_string())?;
    Ok(Some(watcher))
}

/// (Re)starts watching `<dir>/.context`, replacing any previous
/// watcher. Emits "ctx-changed" (payload: the project root) on any
/// filesystem event there; the frontend debounces and refetches. A
/// missing `.context/` is not an error — the watch simply does not
/// start. Watcher setup runs on the blocking pool so the stat/kqueue
/// registration never stalls the UI thread.
#[tauri::command]
pub async fn watch_context(
    app: AppHandle,
    state: State<'_, WatchState>,
    dir: String,
) -> Result<(), String> {
    let watcher =
        tauri::async_runtime::spawn_blocking(move || build_watcher(&app, &dir, "ctx-changed"))
            .await
            .map_err(|e| format!("background task failed: {e}"))??;
    *state.0.lock().map_err(|e| e.to_string())? = watcher;
    Ok(())
}

/// (Re)starts watching the `.context/` of every project in `dirs`,
/// replacing any previous set. Any filesystem event in any of them
/// emits "ctx-projects-changed" with that project's root as payload —
/// a SEPARATE channel from the active project's "ctx-changed", so the
/// dashboard refetches on any project's change without making every
/// per-project screen refetch on a foreign project's write. Projects
/// without a `.context/` are skipped.
#[tauri::command]
pub async fn watch_projects(
    app: AppHandle,
    state: State<'_, WorkspaceWatchState>,
    dirs: Vec<String>,
) -> Result<(), String> {
    let watchers = tauri::async_runtime::spawn_blocking(move || {
        let mut watchers = Vec::new();
        for dir in dirs {
            if let Ok(Some(w)) = build_watcher(&app, &dir, "ctx-projects-changed") {
                watchers.push(w);
            }
        }
        watchers
    })
    .await
    .map_err(|e| format!("background task failed: {e}"))?;
    *state.0.lock().map_err(|e| e.to_string())? = watchers;
    Ok(())
}
