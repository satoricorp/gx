import { afterEach, beforeEach, describe, expect, test } from "bun:test";
import { mkdir, mkdtemp, readFile, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import gxReview, { metadata as reviewMetadata, schema as reviewSchema } from "../src/tools/gx-review";
import { ensureGxInitialized } from "../src/session-workspace";

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

describe("registered MCP tool names", () => {
  test("exposes exactly gx_review", () => {
    const registered = [reviewMetadata.name].sort();
    expect(registered).toEqual(["gx_review"]);
    expect(registered).not.toContain("gx_commit");
    expect(registered).not.toContain("gx_status");
    expect(registered).not.toContain("gx_push");
    expect(registered).not.toContain("gx_publish");
    expect(registered).not.toContain("gx_edit");
  });
});

describe("gx MCP CLI invocation", () => {
  let mockDir: string;
  let repoRoot: string;
  let callLog: string;
  let previousEnv: Record<string, string | undefined>;

  const envKeys = ["GX_BINARY", "GX_MOCK_LOG", "GX_MOCK_AUTH_ERROR"];

  beforeEach(async () => {
    mockDir = await mkdtemp(join(tmpdir(), "gx-mcp-tools-"));
    repoRoot = join(mockDir, "repo");
    callLog = join(mockDir, "calls.log");
    await mkdir(join(repoRoot, ".git"), { recursive: true });
    await Bun.spawn(["git", "init", repoRoot]).exited;
    previousEnv = Object.fromEntries(envKeys.map((key) => [key, process.env[key]]));

    const mockGx = join(mockDir, "mock-gx.sh");
    await writeFile(
      mockGx,
      `#!/bin/sh
printf 'gx|%s|%s|%s\\n' "$PWD" "$GX_REVIEW_AI" "$*" >> "$GX_MOCK_LOG"
if [ "$1" = init ]; then
  touch "$PWD/.gx-initialized"
  echo "initialized"
  exit 0
fi
if [ "$1" = review ]; then
  if [ "$GX_MOCK_AUTH_ERROR" = "1" ]; then
    echo 'github token is not configured for MCP: run \`gx auth login\` in a terminal, then retry the MCP tool' >&2
    exit 1
  fi
  echo "review ok"
  exit 0
fi
if [ "$1" = status ]; then
  if [ ! -f "$PWD/.gx-initialized" ]; then
    echo "gx not initialized" >&2
    exit 1
  fi
  echo '{"stacks":[],"files":[]}'
  exit 0
fi
echo "unexpected gx args: $@" >&2
exit 1
`,
      "utf8",
    );
    await Bun.spawn(["chmod", "+x", mockGx]).exited;

    process.env.GX_BINARY = mockGx;
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
    expect(calls).toContain("gx|");
    expect(calls).toContain("|1|review --scope architecture --focus internal/authoring --deep --verbose review auth rollback risk");
  });

  test("gx_review surfaces MCP auth login guidance", async () => {
    process.env.GX_MOCK_AUTH_ERROR = "1";

    const output = await gxReview({ cwd: repoRoot });
    const parsed = JSON.parse(output);

    expect(parsed.ok).toBe(false);
    expect(parsed.action).toBe("review");
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

  const envKeys = ["GX_BINARY", "GX_MOCK_LOG", "GX_MCP_INIT_NAME", "GX_MCP_INIT_EMAIL"];

  beforeEach(async () => {
    mockDir = await mkdtemp(join(tmpdir(), "gx-mcp-autoinit-"));
    repoRoot = join(mockDir, "repo");
    callLog = join(mockDir, "calls.log");
    await mkdir(join(repoRoot, ".git"), { recursive: true });
    await Bun.spawn(["git", "init", repoRoot]).exited;
    previousEnv = Object.fromEntries(envKeys.map((key) => [key, process.env[key]]));

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
if [ "$1" = review ]; then
  echo "review ok"
  exit 0
fi
if [ "$1" = status ]; then
  if [ ! -f "$PWD/.gx-initialized" ]; then
    echo "gx not initialized" >&2
    exit 1
  fi
  echo '{"stacks":[]}'
  exit 0
fi
echo "unexpected gx args: $@" >&2
exit 1
`,
      "utf8",
    );
    await Bun.spawn(["chmod", "+x", mockGx]).exited;

    process.env.GX_BINARY = mockGx;
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

  test("gx_review auto-inits an uninitialized repo before reviewing", async () => {
    const output = await gxReview({ cwd: repoRoot });
    const parsed = JSON.parse(output);
    expect(parsed.ok).toBe(true);

    const log = await readFile(callLog, "utf8");
    const lines = log.trim().split("\n");
    const initIndex = lines.findIndex((line) => line.includes("|init --name"));
    const reviewIndex = lines.findIndex((line) => line.includes("|review"));
    expect(initIndex).toBeGreaterThanOrEqual(0);
    expect(reviewIndex).toBeGreaterThan(initIndex);
  });
});
