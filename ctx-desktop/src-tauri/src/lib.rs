mod ctx_adapter;

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_opener::init())
        .plugin(tauri_plugin_dialog::init())
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
            ctx_adapter::ctx_journal,
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
