// ctx adapter (frontend half) — typed wrappers over the Rust
// commands in src-tauri/src/ctx_adapter.rs. All ctx access funnels
// through here so a CLI/output change is a one-file fix.
import { invoke } from "@tauri-apps/api/core";

export interface CtxInfo {
  found: boolean;
  version: string;
  error: string | null;
}

// Mirrors `ctx status --json`.
export interface CtxStatus {
  context_dir: string;
  total_files: number;
  total_tokens: number;
  total_size: number;
  files: { name: string; tokens: number; is_empty: boolean }[];
}

// Mirrors `ctx task list --json`.
export interface Task {
  text: string;
  status: "pending" | "done";
  section: string;
  priority: string;
  is_sub: boolean;
  session: string;
  branch: string;
  commit: string;
  added: string;
}

// Mirrors `ctx decision list --json`.
export interface Decision {
  timestamp: string;
  date: string;
  title: string;
  status: string;
  context: string;
  rationale: string;
  consequence: string;
  superseded: boolean;
}

// Mirrors `ctx learning list --json`.
export interface Learning {
  timestamp: string;
  date: string;
  title: string;
  context: string;
  lesson: string;
  application: string;
  superseded: boolean;
}

/** Detect the ctx binary and read its version. */
export function ctxInfo(): Promise<CtxInfo> {
  return invoke<CtxInfo>("ctx_info");
}

/** `ctx status --json` for the project at `dir`. */
export async function ctxStatus(dir: string): Promise<CtxStatus> {
  return JSON.parse(await invoke<string>("ctx_status", { dir }));
}

/** Raw `ctx doctor --json` for the project at `dir`. */
export async function ctxDoctor(dir: string): Promise<unknown> {
  return JSON.parse(await invoke<string>("ctx_doctor", { dir }));
}

/** `ctx task list --json` for the project at `dir`. */
export async function ctxTasks(dir: string): Promise<Task[]> {
  const out = JSON.parse(await invoke<string>("ctx_task_list", { dir }));
  return out.tasks ?? [];
}

/** `ctx decision list --json` for the project at `dir`. */
export async function ctxDecisions(dir: string): Promise<Decision[]> {
  const out = JSON.parse(await invoke<string>("ctx_decision_list", { dir }));
  return out.decisions ?? [];
}

/** `ctx learning list --json` for the project at `dir`. */
export async function ctxLearnings(dir: string): Promise<Learning[]> {
  const out = JSON.parse(await invoke<string>("ctx_learning_list", { dir }));
  return out.learnings ?? [];
}
