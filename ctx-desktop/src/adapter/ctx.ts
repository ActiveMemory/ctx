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

// Mirrors `ctx agent --format json`. Section arrays render
// generically; an empty array means that section was dropped.
export interface AgentPacket {
  generated: string;
  budget: number;
  tokens_used: number;
  read_order: string[];
  constitution: string[];
  tasks: string[];
  conventions: string[];
  decisions: string[];
  learnings: string[];
  summaries: string[];
  instruction: string;
}

/** Detect the ctx binary and read its version. */
export function ctxInfo(): Promise<CtxInfo> {
  return invoke<CtxInfo>("ctx_info");
}

/** Structured context packet from `ctx agent --format json --budget N`. */
export async function ctxAgentPacket(
  dir: string,
  budget: number,
): Promise<AgentPacket> {
  return JSON.parse(await invoke<string>("ctx_agent_json", { dir, budget }));
}

/** Paste-ready markdown packet from `ctx agent --budget N`. */
export function ctxAgentMarkdown(dir: string, budget: number): Promise<string> {
  return invoke<string>("ctx_agent_md", { dir, budget });
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

/**
 * `ctx task add` — provenance (session id, branch, commit) is
 * synthesized in the Rust adapter. Empty priority/section omitted.
 */
export function ctxTaskAdd(
  dir: string,
  text: string,
  priority = "",
  section = "",
): Promise<string> {
  return invoke<string>("ctx_task_add", { dir, text, priority, section });
}

/** `ctx task complete <id-or-text>`. */
export function ctxTaskComplete(dir: string, target: string): Promise<string> {
  return invoke<string>("ctx_task_complete", { dir, target });
}

/**
 * `ctx decision add` — all three ADR fields required; provenance
 * (session id, branch, commit) synthesized in the Rust adapter.
 */
export function ctxDecisionAdd(
  dir: string,
  title: string,
  context: string,
  rationale: string,
  consequence: string,
): Promise<string> {
  return invoke<string>("ctx_decision_add", {
    dir,
    title,
    context,
    rationale,
    consequence,
  });
}

/**
 * `ctx learning add` — all three fields required; provenance
 * synthesized in the Rust adapter.
 */
export function ctxLearningAdd(
  dir: string,
  title: string,
  context: string,
  lesson: string,
  application: string,
): Promise<string> {
  return invoke<string>("ctx_learning_add", {
    dir,
    title,
    context,
    lesson,
    application,
  });
}
