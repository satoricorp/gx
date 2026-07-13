import { afterEach, beforeEach, describe, expect, test } from "bun:test";
import { mkdir, mkdtemp, readFile, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import gxCommit, { metadata as commitMetadata, schema as commitSchema } from "../src/tools/gx-commit";
import gxEdit, { metadata as editMetadata, schema as editSchema } from "../src/tools/gx-edit";
import gxReview, { metadata as reviewMetadata, schema as reviewSchema } from "../src/tools/gx-review";
import gxStatus, { metadata as statusMetadata, schema as statusSchema } from "../src/tools/gx-status";
import { ensureGxInitialized } from "../src/session-workspace";

describe("gx_commit metadata and schema", () => {
  test("describes the commit surface", () => {
    expect(commitMetadata.name).toBe("gx_commit");
    expect(commitMetadata.description).toMatch(/gx commit/i);
    expect(commitMetadata.annotations?.readOnlyHint).toBe(false);
    expect(commitSchema.message.parse("record staged work")).toBe("record staged work");
    expect(commitSchema.commands_run.parse(["go test ./..."])).toEqual(["go test ./..."]);
  });
});

describe("gx_edit metadata and schema", () => {
  test("describes the edit surface", () => {
    expect(editMetadata.name).toBe("gx_edit");
    expect(editMetadata.description).toMatch(/gx edit/i);
    expect(editMetadata.annotations?.readOnlyHint).toBe(false);
    expect(editSchema.revision.parse("abc123")).toBe("abc123");
  });
});

describe("gx_review metadata and schema", () => {
  test("describes the review surface", () => {
    expect(reviewMetadata.name).toBe("gx_review");
    expect(reviewMetadata.description).toMatch(/local facts/i);
    expect(reviewMetadata.annotations?.readOnlyHint).toBe(true);
    expect(reviewSchema.scope.parse("architecture")).toBe("architecture");
    expect(reviewSchema.prompt.parse("review auth rollback risk")).toBe("review auth rollback risk");
    expect(reviewSchema.deep.parse(true)).toBe(true);
  });
});

describe("gx_status metadata and schema", () => {
  test("describes status as read-only", () => {
    expect(statusMetadata.name).toBe("gx_status");
    expect(statusMetadata.description).toMatch(/gx status --json/i);
    expect(statusMetadata.annotations?.readOnlyHint).toBe(true);
    expect(statusSchema.show_all.parse(true)).toBe(true);
  });
});

describe("registered MCP tool names", () => {
  test("exposes exactly gx_commit, gx_status, gx_edit, gx_review", () => {
    const registered = [commitMetadata.name, statusMetadata.name, editMetadata.name, reviewMetadata.name].sort();
    expect(registered).toEqual(["gx_commit", "gx_edit", "gx_review", "gx_status"]);
    expect(registered).not.toContain("gx_push");
    expect(registered).not.toContain("gx_publish");
  });
});

describe("gx MCP CLI invocation", () => {
  let mockDir: string;
  let repoRoot: string;
  let callLog: string;
  let previousEnv: Record<string, string | undefined>;

  const envKeys = ["GX_BINARY", "JJ_BINARY", "GX_MOCK_LOG", "GX_MOCK_AUTH_ERROR"];

  beforeEach(async () => {
    mockDir = await mkdtemp(join(tmpdir(), "gx-mcp-tools-"));
    repoRoot = join(mockDir, "repo");
    callLog = join(mockDir, "calls.log");
    await mkdir(join(repoRoot, ".git"), { recursive: true });
    previousEnv = Object.fromEntries(envKeys.map((key) => [key, process.env[key]]));

    const mockGx = join(mockDir, "mock-gx.sh");
    await writeFile(
      mockGx,
      `#!/bin/sh
printf 'gx|%s|%s|%s\\n' "$PWD" "$GX_REVIEW_AI" "$*" >> "$GX_MOCK_LOG"
if [ "$1" = commit ]; then
  if [ "$GX_MOCK_AUTH_ERROR" = "1" ]; then
    echo 'github token is not configured for MCP: run \`gx auth login\` in a terminal, then retry the MCP tool' >&2
    exit 1
  fi
  echo "commit ok"
  exit 0
fi
if [ "$1" = edit ]; then
  echo "edit ok"
  exit 0
fi
if [ "$1" = review ]; then
  echo "review ok"
  exit 0
fi
if [ "$1" = status ]; then
  echo '{"stacks":[],"files":[]}'
  exit 0
fi
echo "unexpected gx args: $@" >&2
exit 1
`,
      "utf8",
    );
    await Bun.spawn(["chmod", "+x", mockGx]).exited;

    const mockJj = join(mockDir, "mock-jj.sh");
    await writeFile(
      mockJj,
      `#!/bin/sh
printf 'jj|%s|%s\\n' "$PWD" "$*" >> "$GX_MOCK_LOG"
if [ "$1" = root ]; then
  echo "$PWD"
  exit 0
fi
echo "unexpected jj args: $@" >&2
exit 1
`,
      "utf8",
    );
    await Bun.spawn(["chmod", "+x", mockJj]).exited;

    process.env.GX_BINARY = mockGx;
    process.env.JJ_BINARY = mockJj;
    process.env.GX_MOCK_LOG = callLog;
    delete process.env.GX_MOCK_AUTH_ERROR;
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

  test("gx_commit passes message and context file", async () => {
    const output = await gxCommit({
      cwd: repoRoot,
      message: "record staged work",
      task_summary: "add commit surface",
      commands_run: ["go test ./internal/commitcontext"],
      tests_run: ["mcp/test/gx-tools.test.ts"],
      session_id: "session-commit",
    });
    const parsed = JSON.parse(output);
    expect(parsed.action).toBe("commit");
    expect(parsed.display).toBe("commit ok");
    expect(parsed.command[0]).toBe(process.env.GX_BINARY);
    expect(parsed.command).toContain("commit");
    expect(parsed.command).toContain("-m");
    expect(parsed.command).toContain("record staged work");
    expect(parsed.command).toContain("--context-file");

    const calls = await readFile(callLog, "utf8");
    expect(calls).toContain("commit -m");
    expect(calls).toContain("--context-file");
  });

  test("gx_edit passes the revision id", async () => {
    const output = await gxEdit({ cwd: repoRoot, revision: "abc123def" });
    const parsed = JSON.parse(output);
    expect(parsed.action).toBe("edit");
    expect(parsed.display).toBe("edit ok");
    expect(parsed.command).toEqual([process.env.GX_BINARY, "edit", "abc123def"]);
  });

  test("gx_review passes scope, focus, prompt, deep, and verbose flags", async () => {
    const output = await gxReview({
      cwd: repoRoot,
      scope: "architecture",
      focus: "internal/authoring",
      prompt: "review auth rollback risk",
      deep: true,
      verbose: true,
    });
    const parsed = JSON.parse(output);
    expect(parsed.action).toBe("review");
    expect(parsed.display).toBe("review ok");
    expect(parsed.command).toEqual([
      process.env.GX_BINARY,
      "review",
      "--scope",
      "architecture",
      "--focus",
      "internal/authoring",
      "--deep",
      "--verbose",
      "review auth rollback risk",
    ]);

    const calls = await readFile(callLog, "utf8");
    expect(calls).toContain("jj|");
    expect(calls).toContain("gx|");
    expect(calls).toContain("|1|review --scope architecture --focus internal/authoring --deep --verbose review auth rollback risk");
  });

  test("gx_status runs status json", async () => {
    const output = await gxStatus({ cwd: repoRoot });
    const parsed = JSON.parse(output);
    expect(parsed.action).toBe("status");
    expect(parsed.result).toEqual({ stacks: [], files: [] });
    expect(parsed.command).toEqual([process.env.GX_BINARY, "status", "--json"]);
    expect(parsed.next_actions).toEqual([
      "Stage with git add and gx_commit for new work, gx_edit to continue a revision, or git push to publish when a stack is ready.",
    ]);
  });

  test("gx_status passes --all when show_all is set", async () => {
    const output = await gxStatus({ cwd: repoRoot, show_all: true });
    const parsed = JSON.parse(output);
    expect(parsed.command).toEqual([process.env.GX_BINARY, "status", "--json", "--all"]);
  });

  test("gx_commit surfaces MCP auth login guidance", async () => {
    process.env.GX_MOCK_AUTH_ERROR = "1";

    const output = await gxCommit({ cwd: repoRoot, message: "record staged work" });
    const parsed = JSON.parse(output);

    expect(parsed.ok).toBe(false);
    expect(parsed.action).toBe("commit");
    expect(parsed.auth_required).toBe(true);
    expect(parsed.display).toBe(
      "GX cloud authentication is required. Run `gx auth login` in a terminal, then retry the MCP tool.",
    );
    expect(parsed.stderr).toContain("run `gx auth login` in a terminal");
    expect(parsed.next_actions[0]).toBe("Run `gx auth login` in a terminal, then retry the MCP tool.");
  });
});

describe("ensureGxInitialized auto-init", () => {
  let mockDir: string;
  let repoRoot: string;
  let callLog: string;
  let previousEnv: Record<string, string | undefined>;

  const envKeys = ["GX_BINARY", "JJ_BINARY", "GX_MOCK_LOG", "GX_MCP_INIT_NAME", "GX_MCP_INIT_EMAIL"];

  beforeEach(async () => {
    mockDir = await mkdtemp(join(tmpdir(), "gx-mcp-autoinit-"));
    repoRoot = join(mockDir, "repo");
    callLog = join(mockDir, "calls.log");
    await mkdir(join(repoRoot, ".git"), { recursive: true });
    previousEnv = Object.fromEntries(envKeys.map((key) => [key, process.env[key]]));

    const mockJj = join(mockDir, "mock-jj.sh");
    await writeFile(
      mockJj,
      `#!/bin/sh
printf 'jj|%s|%s\\n' "$PWD" "$*" >> "$GX_MOCK_LOG"
if [ "$1" = root ]; then
  if [ -f "$PWD/.gx-initialized" ]; then
    echo "$PWD"
    exit 0
  fi
  echo "There is no jj repo in the current directory" >&2
  exit 1
fi
echo "unexpected jj args: $@" >&2
exit 1
`,
      "utf8",
    );
    await Bun.spawn(["chmod", "+x", mockJj]).exited;

    const mockGx = join(mockDir, "mock-gx.sh");
    await writeFile(
      mockGx,
      `#!/bin/sh
printf 'gx|%s|%s\\n' "$PWD" "$*" >> "$GX_MOCK_LOG"
if [ "$1" = init ]; then
  touch "$PWD/.gx-initialized"
  echo "initialized"
  exit 0
fi
if [ "$1" = commit ]; then
  echo "commit ok"
  exit 0
fi
echo "unexpected gx args: $@" >&2
exit 1
`,
      "utf8",
    );
    await Bun.spawn(["chmod", "+x", mockGx]).exited;

    process.env.GX_BINARY = mockGx;
    process.env.JJ_BINARY = mockJj;
    process.env.GX_MOCK_LOG = callLog;
    process.env.GX_MCP_INIT_NAME = "Autoinit Test";
    process.env.GX_MCP_INIT_EMAIL = "autoinit@test.local";
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

  test("runs gx init when the repo is uninitialized, then is a no-op", async () => {
    const first = await ensureGxInitialized(repoRoot);
    expect(first.autoInitialized).toBe(true);

    const afterFirst = await readFile(callLog, "utf8");
    expect(afterFirst).toContain("init --name Autoinit Test --email autoinit@test.local");

    const second = await ensureGxInitialized(repoRoot);
    expect(second.autoInitialized).toBe(false);

    const afterSecond = await readFile(callLog, "utf8");
    const initCalls = afterSecond.split("\n").filter((line) => line.includes("|init --name"));
    expect(initCalls).toHaveLength(1);
  });

  test("gx_commit auto-inits an uninitialized repo before committing", async () => {
    const output = await gxCommit({ cwd: repoRoot, message: "record staged work" });
    const parsed = JSON.parse(output);
    expect(parsed.ok).toBe(true);

    const log = await readFile(callLog, "utf8");
    const lines = log.trim().split("\n");
    const initIndex = lines.findIndex((line) => line.includes("|init --name"));
    const commitIndex = lines.findIndex((line) => line.includes("|commit -m"));
    expect(initIndex).toBeGreaterThanOrEqual(0);
    expect(commitIndex).toBeGreaterThan(initIndex);
  });
});
