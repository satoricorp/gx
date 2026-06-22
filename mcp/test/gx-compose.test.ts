import { afterEach, beforeEach, describe, expect, test } from "bun:test";
import { createHash } from "node:crypto";
import { mkdir, mkdtemp, readFile, realpath, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { metadata, schema } from "../src/tools/gx-compose-changes";
import gxCompose from "../src/tools/gx-compose-changes";

const composeFixture = {
  state: "ready",
  proposal: {
    id: "demux-1",
    revisions: [{ id: "u1", intent: "test", target_stack: "feature/test" }],
    hunks: [],
    hunk_links: [],
  },
};

describe("gx_compose metadata", () => {
  test("registers canonical V1 tool name", () => {
    expect(metadata.name).toBe("gx_compose");
    expect(metadata.description).not.toMatch(/brief/i);
    expect(metadata.description).toMatch(/session-isolated JJ workspace/i);
    expect(metadata.description).toMatch(/Requires a GX session id/i);
    expect(metadata.description).toMatch(/appends to existing target bookmarks/i);
    expect(metadata.description).toMatch(/default branch/i);
    expect(metadata.description).not.toMatch(/\/v1\//);
  });
});

describe("gx_compose schema", () => {
  test("accepts propose defaults", () => {
    const parsed = schema.action?.parse(undefined);
    expect(parsed).toBeUndefined();
    expect(schema.cwd.parse(undefined)).toBeUndefined();
    expect(schema.proposal.parse(undefined)).toBeUndefined();
  });

  test("accepts review-plan and apply actions", () => {
    expect(schema.action.parse("review-plan")).toBe("review-plan");
    expect(schema.action.parse("apply")).toBe("apply");
    expect(schema.proposal.parse({ revisions: [] })).toEqual({ revisions: [] });
  });
});

describe("gx_compose CLI invocation", () => {
  let mockDir: string;
  let repoRoot: string;
  let workspaceRoot: string;
  let callLog: string;
  let previousEnv: Record<string, string | undefined>;

  const envKeys = [
    "GX_BINARY",
    "JJ_BINARY",
    "GIT_BINARY",
    "GX_MCP_WORKSPACE_ROOT",
    "GX_SESSION_ID",
    "GX_SESSION_IDS",
    "CODEX_SESSION_ID",
    "CLAUDE_SESSION_ID",
    "CLAUDE_CODE_SESSION_ID",
    "CURSOR_SESSION_ID",
    "OPENCODE_SESSION_ID",
    "QODO_SESSION_ID",
    "QUAD_CODE_SESSION_ID",
    "MOCK_REPO_ROOT",
    "MOCK_GIT_COMMON_DIR",
    "MOCK_JJ_ROOT_FAIL_ONCE",
    "MOCK_JJ_ROOT_COUNT",
    "GX_MOCK_LOG",
  ];

  beforeEach(async () => {
    mockDir = await mkdtemp(join(tmpdir(), "gx-mcp-mock-"));
    repoRoot = join(mockDir, "repo");
    workspaceRoot = join(mockDir, "workspaces");
    callLog = join(mockDir, "calls.log");
    await mkdir(join(repoRoot, ".git"), { recursive: true });
    previousEnv = Object.fromEntries(envKeys.map((key) => [key, process.env[key]]));

    const mockScript = join(mockDir, "mock-gx.sh");
    await writeFile(
      mockScript,
      `#!/bin/sh
printf 'gx|%s|%s|%s\\n' "$PWD" "$GX_SESSION_ID$GX_SESSION_IDS" "$*" >> "$GX_MOCK_LOG"
if [ "$1" = init ]; then
  echo "initialized"
  exit 0
fi
if [ "$1" = compose ] && [ "$2" = --json ]; then
  echo '${JSON.stringify(composeFixture)}'
  exit 0
fi
if [ "$1" = compose ] && [ "$2" = apply ]; then
  echo '{"state":"applied","accepted":1}'
  exit 0
fi
if [ "$1" = compose ] && [ "$2" = review-plan ]; then
  echo '{"state":"reviewed","ready":true}'
  exit 0
fi
if [ "$1" = compose ] && [ "$2" = apply-plan ]; then
  echo '{"state":"applied","accepted":1}'
  exit 0
fi
echo "unexpected args: $@" >&2
exit 1
`,
      "utf8",
    );
    await Bun.spawn(["chmod", "+x", mockScript]).exited;

    const mockJJ = join(mockDir, "mock-jj.sh");
    await writeFile(
      mockJJ,
      `#!/bin/sh
printf 'jj|%s|%s\\n' "$PWD" "$*" >> "$GX_MOCK_LOG"
if [ "$1" = root ]; then
  if [ "$MOCK_JJ_ROOT_FAIL_ONCE" = "1" ]; then
    if [ ! -f "$MOCK_JJ_ROOT_COUNT" ]; then
      echo 1 > "$MOCK_JJ_ROOT_COUNT"
      echo "No jj repo found" >&2
      exit 1
    fi
  fi
  echo "$MOCK_REPO_ROOT"
  exit 0
fi
if [ "$1" = workspace ] && [ "$2" = add ]; then
  last=""
  for arg in "$@"; do
    last="$arg"
  done
  mkdir -p "$last"
  exit 0
fi
echo "unexpected jj args: $@" >&2
exit 1
`,
      "utf8",
    );
    await Bun.spawn(["chmod", "+x", mockJJ]).exited;

    const mockGit = join(mockDir, "mock-git.sh");
    await writeFile(
      mockGit,
      `#!/bin/sh
printf 'git|%s|%s\\n' "$PWD" "$*" >> "$GX_MOCK_LOG"
if [ "$1" = rev-parse ] && [ "$2" = --git-common-dir ]; then
  echo "$MOCK_GIT_COMMON_DIR"
  exit 0
fi
if [ "$1" = config ] && [ "$2" = user.name ]; then
  echo "Test User"
  exit 0
fi
if [ "$1" = config ] && [ "$2" = user.email ]; then
  echo "test@example.com"
  exit 0
fi
echo "unexpected git args: $@" >&2
exit 1
`,
      "utf8",
    );
    await Bun.spawn(["chmod", "+x", mockGit]).exited;

    process.env.GX_BINARY = mockScript;
    process.env.JJ_BINARY = mockJJ;
    process.env.GIT_BINARY = mockGit;
    process.env.GX_MCP_WORKSPACE_ROOT = workspaceRoot;
    process.env.MOCK_REPO_ROOT = repoRoot;
    process.env.MOCK_GIT_COMMON_DIR = join(repoRoot, ".git");
    process.env.MOCK_JJ_ROOT_COUNT = join(mockDir, "jj-root-count");
    process.env.GX_MOCK_LOG = callLog;
    delete process.env.MOCK_JJ_ROOT_FAIL_ONCE;
    delete process.env.GX_SESSION_ID;
    delete process.env.GX_SESSION_IDS;
    delete process.env.CODEX_SESSION_ID;
    delete process.env.CLAUDE_SESSION_ID;
    delete process.env.CLAUDE_CODE_SESSION_ID;
    delete process.env.CURSOR_SESSION_ID;
    delete process.env.OPENCODE_SESSION_ID;
    delete process.env.QODO_SESSION_ID;
    delete process.env.QUAD_CODE_SESSION_ID;
  });

  afterEach(() => {
    for (const key of envKeys) {
      const value = previousEnv[key];
      if (value === undefined) {
        delete process.env[key];
      } else {
        process.env[key] = value;
      }
    }
  });

  test("propose runs gx compose --json in the session workspace and auto-accepts ready proposals", async () => {
    const output = await gxCompose({ cwd: repoRoot, intent: "ship it", session_id: "codex-session-1" });
    const parsed = JSON.parse(output);
    const repoHash = hash(await realpath(join(repoRoot, ".git")));
    const sessionHash = hash("codex-session-1");
    const workspacePath = join(workspaceRoot, repoHash, sessionHash);
    const realRepoRoot = await realpath(repoRoot);
    const realWorkspacePath = await realpath(workspacePath);

    expect(parsed.exit_code).toBe(0);
    expect(parsed.action).toBe("compose.auto_accept");
    expect(parsed.command).toEqual([process.env.GX_BINARY, "compose", "apply", "demux-1", "--json"]);
    expect(parsed.cwd).toBe(workspacePath);
    expect(parsed.session_workspace).toMatchObject({
      cwd: workspacePath,
      repoRoot,
      repoHash,
      sessionHash,
      workspaceName: `gx-${sessionHash}`,
      created: true,
    });
    expect(parsed.result.compose.state).toBe("ready");
    expect(parsed.result.accept.state).toBe("applied");
    expect(parsed.overview.hunk_link_count).toBe(0);

    const calls = await readFile(callLog, "utf8");
    expect(calls).toContain(`jj|${realRepoRoot}|workspace add --name gx-${sessionHash} ${workspacePath}`);
    expect(calls).toContain(`gx|${realWorkspacePath}|codex-session-1|compose --json --intent ship it`);
    expect(calls).toContain(`gx|${realWorkspacePath}|codex-session-1|compose apply demux-1 --json`);
  });

  test("propose can leave a ready proposal pending", async () => {
    const output = await gxCompose({
      cwd: repoRoot,
      intent: "ship it",
      session_id: "codex-session-1",
      auto_accept: false,
    });
    const parsed = JSON.parse(output);
    expect(parsed.action).toBe("compose.propose");
    expect(parsed.command).toEqual([process.env.GX_BINARY, "compose", "--json", "--intent", "ship it"]);
    expect(parsed.result.state).toBe("ready");
    const calls = await readFile(callLog, "utf8");
    expect(calls).not.toContain("compose apply demux-1");
  });

  test("auto-initializes a git repo before creating the session workspace", async () => {
    process.env.MOCK_JJ_ROOT_FAIL_ONCE = "1";
    const output = await gxCompose({ cwd: repoRoot, intent: "ship it", session_id: "codex-session-1" });
    const parsed = JSON.parse(output);

    expect(parsed.exit_code).toBe(0);
    expect(parsed.session_workspace.autoInitialized).toBe(true);
    expect(parsed.session_workspace.initOutput).toBe("initialized");

    const calls = await readFile(callLog, "utf8");
    const realRepoRoot = await realpath(repoRoot);
    expect(calls).toContain(`gx|${realRepoRoot}||init --name Test User --email test@example.com`);
    expect(calls).toContain("compose apply demux-1 --json");
  });

  test("propose refuses without a session id", async () => {
    const output = await gxCompose({ cwd: repoRoot, intent: "ship it" });
    expect(output).toContain("gx_compose requires a session id");
  });

  test("review-plan runs gx compose review-plan", async () => {
    const output = await gxCompose({
      action: "review-plan",
      proposal: { revisions: [], hunks: [] },
      cwd: repoRoot,
      session_id: "codex-session-2",
    });
    const parsed = JSON.parse(output);
    expect(parsed.exit_code).toBe(0);
    expect(parsed.command.slice(0, 4)).toEqual([
      process.env.GX_BINARY,
      "compose",
      "review-plan",
      "--json",
    ]);
    expect(parsed.result.state).toBe("reviewed");
  });

  test("apply runs gx compose apply-plan", async () => {
    const output = await gxCompose({
      action: "apply",
      proposal: { revisions: [], hunks: [] },
      allow_warnings: true,
      cwd: repoRoot,
      session_id: "codex-session-3",
    });
    const parsed = JSON.parse(output);
    expect(parsed.exit_code).toBe(0);
    expect(parsed.command).toContain("--allow-warnings");
    expect(parsed.result.state).toBe("applied");
  });
});

function hash(value: string): string {
  return createHash("sha256").update(value).digest("hex").slice(0, 16);
}
