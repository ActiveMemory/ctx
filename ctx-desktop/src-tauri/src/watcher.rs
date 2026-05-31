//! File watcher — emits a "ctx-changed" event whenever the active
//! project's `.context/` directory mutates, so the frontend can
//! refetch. The GUI is one writer among several (human CLI + AI
//! agents); this keeps its views from going stale.

use std::path::Path;
use std::sync::Mutex;

use notify::{RecommendedWatcher, RecursiveMode, Watcher};
use tauri::{AppHandle, Emitter, State};

/// Holds the live watcher so it stays alive between calls. Replacing
/// the value drops (and stops) the previous watcher.
#[derive(Default)]
pub struct WatchState(pub Mutex<Option<RecommendedWatcher>>);

/// (Re)starts watching `<dir>/.context`, replacing any previous
/// watcher. Emits "ctx-changed" on any filesystem event there; the
/// frontend debounces and refetches. A missing `.context/` is not
/// an error — the watch simply does not start.
#[tauri::command]
pub fn watch_context(
    app: AppHandle,
    state: State<'_, WatchState>,
    dir: String,
) -> Result<(), String> {
    let ctx_dir = Path::new(&dir).join(".context");
    if !ctx_dir.is_dir() {
        *state.0.lock().map_err(|e| e.to_string())? = None;
        return Ok(());
    }

    let handle = app.clone();
    let mut watcher = notify::recommended_watcher(move |res: notify::Result<notify::Event>| {
        if res.is_ok() {
            let _ = handle.emit("ctx-changed", ());
        }
    })
    .map_err(|e| e.to_string())?;
    watcher
        .watch(&ctx_dir, RecursiveMode::Recursive)
        .map_err(|e| e.to_string())?;

    *state.0.lock().map_err(|e| e.to_string())? = Some(watcher);
    Ok(())
}
