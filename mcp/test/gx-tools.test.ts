import { afterEach, beforeEach, describe, expect, test } from "bun:test";
import { mkdir, mkdtemp, readFile, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import gxPush, { metadata as pushMetadata, schema as pushSchema } from "../src/tools/gx-push";
import gxReview, { metadata as reviewMetadata, schema as reviewSchema } from "../src/tools/gx-review";
import gxStatus, { metadata as statusMetadata, schema as statusSchema } from "../src/tools/gx-status";

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

describe("gx_push metadata and schema", () => {
  test("describes pushing generated features", () => {
    expect(pushMetadata.name).toBe("gx_push");
    expect(pushMetadata.description).toMatch(/generated GX features/i);
    expect(pushMetadata.annotations?.readOnlyHint).toBe(false);
    expect(pushSchema.stack.parse("feature/review")).toBe("feature/review");
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

describe("gx_review and gx_push CLI invocation", () => {
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
if [ "$1" = review ]; then
  echo "review ok"
  exit 0
fi
if [ "$1" = push ]; then
  if [ "$GX_MOCK_AUTH_ERROR" = "1" ]; then
    echo 'github token is not configured for MCP: run \`gx auth login\` in a terminal, then retry the MCP tool' >&2
    exit 1
  fi
  echo "push ok"
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

  test("gx_push passes an optional stack name", async () => {
    const output = await gxPush({ cwd: repoRoot, stack: "feature/review" });
    const parsed = JSON.parse(output);
    expect(parsed.action).toBe("push");
    expect(parsed.display).toBe("push ok");
    expect(parsed.command).toEqual([process.env.GX_BINARY, "push", "feature/review"]);
    expect(parsed.next_actions).toEqual(["Run gx_status to verify remote state. Run gx_sync after remote merges land."]);
  });

  test("gx_status runs status json", async () => {
    const output = await gxStatus({ cwd: repoRoot });
    const parsed = JSON.parse(output);
    expect(parsed.action).toBe("status");
    expect(parsed.result).toEqual({ stacks: [], files: [] });
    expect(parsed.command).toEqual([process.env.GX_BINARY, "status", "--json"]);
    expect(parsed.next_actions).toEqual(["Run gx_generate for local changes or gx_push for ready features."]);
  });

  test("gx_push surfaces MCP auth login guidance", async () => {
    process.env.GX_MOCK_AUTH_ERROR = "1";

    const output = await gxPush({ cwd: repoRoot });
    const parsed = JSON.parse(output);

    expect(parsed.ok).toBe(false);
    expect(parsed.action).toBe("push");
    expect(parsed.auth_required).toBe(true);
    expect(parsed.display).toBe(
      "GX cloud authentication is required. Run `gx auth login` in a terminal, then retry the MCP tool.",
    );
    expect(parsed.stderr).toContain("run `gx auth login` in a terminal");
    expect(parsed.next_actions[0]).toBe("Run `gx auth login` in a terminal, then retry the MCP tool.");
  });
});
