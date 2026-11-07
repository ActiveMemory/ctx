// Test double for the `vscode` module (provided by the extension host at
// runtime, so it is not installed). Test files use it via:
//   vi.mock("vscode", async () => (await import("./vscodeMock")).createVscodeMock());
import { vi } from "vitest";

export function createVscodeMock() {
  class Uri {
    constructor(public fsPath: string) {}
    static joinPath = vi.fn();
  }
  class ChatRequestTurn {
    constructor(
      public prompt: string,
      public command?: string
    ) {}
  }
  class ChatResponseMarkdownPart {
    value: { value: string };
    constructor(value: string) {
      this.value = { value };
    }
  }
  class ChatResponseTurn {
    constructor(
      public response: ChatResponseMarkdownPart[],
      public result: { metadata?: { command: string } },
      public command?: string
    ) {}
  }
  return {
    workspace: {
      getConfiguration: vi.fn(() => ({ get: vi.fn(() => undefined) })),
      workspaceFolders: [{ uri: { fsPath: "/test/workspace" } }] as
        | Array<{ uri: { fsPath: string } }>
        | undefined,
      getWorkspaceFolder: vi.fn(() => undefined),
      asRelativePath: vi.fn((uri: Uri) => uri.fsPath),
      fs: { readFile: vi.fn(async () => new TextEncoder().encode("attached text")) },
    },
    window: { activeTextEditor: undefined as unknown },
    env: { sessionId: "0123456789abcdef" },
    lm: { selectChatModels: vi.fn(async () => []) },
    chat: { createChatParticipant: vi.fn(() => ({})) },
    LanguageModelChatMessage: {
      User: (content: string) => ({ role: "user", content }),
      Assistant: (content: string) => ({ role: "assistant", content }),
    },
    Uri,
    ChatRequestTurn,
    ChatResponseTurn,
    ChatResponseMarkdownPart,
  };
}
