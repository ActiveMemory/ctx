//! Project discovery — scans a workspace root for ctx projects.
//!
//! A directory is a ctx project when it contains a `.context/`
//! directory. The walk is depth-bounded and skips heavy/irrelevant
//! directories so a deep tree stays fast.

use std::path::Path;

use serde::Serialize;

/// Directories we never descend into during a scan.
const SKIP_DIRS: &[&str] = &[
    "node_modules",
    "target",
    "dist",
    "build",
    "vendor",
    ".git",
    ".context",
    ".cache",
    ".venv",
    "venv",
];

/// Hard cap on results so a huge tree can't hang the UI.
const MAX_RESULTS: usize = 200;

/// A discovered ctx project.
#[derive(Serialize)]
pub struct Project {
    pub path: String,
    pub name: String,
    pub has_git: bool,
    /// Current git branch, or "" when not a repo / detached / git absent.
    pub branch: String,
}

/// Reads the current branch of the repo at `dir`, or "" on any failure
/// or a detached HEAD (where `--abbrev-ref` yields the literal "HEAD").
fn git_branch(dir: &Path) -> String {
    std::process::Command::new("git")
        .arg("-C")
        .arg(dir)
        .args(["rev-parse", "--abbrev-ref", "HEAD"])
        .output()
        .ok()
        .filter(|o| o.status.success())
        .map(|o| String::from_utf8_lossy(&o.stdout).trim().to_string())
        .filter(|b| b.as_str() != "HEAD")
        .unwrap_or_default()
}

/// Scans `root` (and up to `max_depth` levels below it) for ctx
/// projects, returning them sorted by name. A `.context/` dir marks
/// a project. Stops at MAX_RESULTS.
#[tauri::command]
pub fn discover_projects(root: String, max_depth: u32) -> Result<Vec<Project>, String> {
    let base = Path::new(&root);
    if !base.is_dir() {
        return Err(format!("not a directory: {root}"));
    }
    let mut out = Vec::new();
    walk(base, max_depth, &mut out);
    out.sort_by_key(|p| p.name.to_lowercase());
    Ok(out)
}

/// Depth-first walk collecting projects, descending `depth_left`
/// more levels.
fn walk(dir: &Path, depth_left: u32, out: &mut Vec<Project>) {
    if out.len() >= MAX_RESULTS {
        return;
    }
    if dir.join(".context").is_dir() {
        let name = dir
            .file_name()
            .map(|s| s.to_string_lossy().to_string())
            .unwrap_or_default();
        let has_git = dir.join(".git").exists();
        out.push(Project {
            path: dir.to_string_lossy().to_string(),
            name,
            has_git,
            branch: if has_git {
                git_branch(dir)
            } else {
                String::new()
            },
        });
    }
    if depth_left == 0 {
        return;
    }
    let entries = match std::fs::read_dir(dir) {
        Ok(e) => e,
        Err(_) => return,
    };
    for entry in entries.flatten() {
        let path = entry.path();
        if !path.is_dir() {
            continue;
        }
        let name = entry.file_name().to_string_lossy().to_string();
        if name.starts_with('.') || SKIP_DIRS.contains(&name.as_str()) {
            continue;
        }
        walk(&path, depth_left - 1, out);
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::fs;

    #[test]
    fn finds_context_dirs_and_skips_noise() {
        let tmp = std::env::temp_dir().join(format!("ctxdisc_{}", std::process::id()));
        let _ = fs::remove_dir_all(&tmp);
        fs::create_dir_all(tmp.join("proj-a/.context")).unwrap();
        fs::create_dir_all(tmp.join("nested/proj-b/.context")).unwrap();
        fs::create_dir_all(tmp.join("node_modules/dep/.context")).unwrap();

        let got = discover_projects(tmp.to_string_lossy().to_string(), 4).unwrap();
        let names: Vec<&str> = got.iter().map(|p| p.name.as_str()).collect();

        assert!(
            names.contains(&"proj-a"),
            "depth-1 project found: {names:?}"
        );
        assert!(names.contains(&"proj-b"), "nested project found: {names:?}");
        assert!(
            !names.contains(&"dep"),
            "node_modules must be skipped: {names:?}"
        );

        let _ = fs::remove_dir_all(&tmp);
    }

    #[test]
    fn rejects_non_directory_root() {
        assert!(discover_projects("/definitely/not/here".to_string(), 2).is_err());
    }
}
