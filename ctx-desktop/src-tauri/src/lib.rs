mod ctx_adapter;
mod discover;
mod watcher;

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_opener::init())
        .plugin(tauri_plugin_dialog::init())
        .manage(watcher::WatchState::default())
        .manage(watcher::WorkspaceWatchState::default())
        .invoke_handler(tauri::generate_handler![
            ctx_adapter::ctx_info,
            ctx_adapter::ctx_status,
            ctx_adapter::ctx_doctor,
            ctx_adapter::ctx_task_list,
            ctx_adapter::ctx_decision_list,
            ctx_adapter::ctx_learning_list,
            ctx_adapter::ctx_task_add,
            ctx_adapter::ctx_task_complete,
            ctx_adapter::ctx_decision_add,
            ctx_adapter::ctx_learning_add,
            ctx_adapter::ctx_agent_json,
            ctx_adapter::ctx_agent_md,
            ctx_adapter::ctx_read_doc,
            ctx_adapter::kb_info,
            ctx_adapter::kb_read,
            ctx_adapter::ctx_journal,
            ctx_adapter::ctx_journal_show,
            ctx_adapter::ctx_drift,
            ctx_adapter::ctx_compact,
            ctx_adapter::ctx_remind_list,
            ctx_adapter::ctx_remind_add,
            ctx_adapter::ctx_remind_dismiss,
            ctx_adapter::ctx_pad_list,
            ctx_adapter::ctx_pad_add,
            ctx_adapter::ctx_pad_rm,
            ctx_adapter::ctx_pad_show,
            ctx_adapter::ctx_connection_status,
            ctx_adapter::set_ctx_path,
            ctx_adapter::dir_is_ctx_project,
            discover::discover_projects,
            watcher::watch_context,
            watcher::watch_projects,
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
