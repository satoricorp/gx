import { afterEach, beforeEach, describe, expect, test } from "bun:test";
import { existsSync } from "node:fs";
import { mkdir, mkdtemp, readFile, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import gxReview, { metadata as reviewMetadata, schema as reviewSchema } from "../src/tools/gx-review";

describe("gx_review metadata and schema", () => {
  test("describes the review surface", () => {
    expect(reviewMetadata.name).toBe("gx_review");
    expect(reviewMetadata.description).toMatch(/local facts/i);
    expect(reviewMetadata.annotations?.readOnlyHint).toBe(true);
    expect(reviewSchema.scope.parse("architecture")).toBe("architecture");
    expect(reviewSchema.prompt.parse("review auth rollback risk")).toBe("review auth rollback risk");
    expect(reviewSchema.deep.parse(true)).toBe(true);
    expect(reviewSchema.repo.parse(true)).toBe(true);
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

    const mockTt = join(mockDir, "mock-gx.sh");
    await writeFile(
      mockTt,
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
if [ "$1" = doctor ]; then
  if [ ! -f "$PWD/.gx-initialized" ]; then
    echo "gx not initialized" >&2
    exit 1
  fi
  echo '{"doctor":{}}'
  exit 0
fi
echo "unexpected gx args: $@" >&2
exit 1
`,
      "utf8",
    );
    await Bun.spawn(["chmod", "+x", mockTt]).exited;

    process.env.GX_BINARY = mockTt;
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
      "--no-publish",
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
    expect(calls).toContain(
      "|1|review --no-publish --scope architecture --focus internal/authoring --deep --verbose review auth rollback risk",
    );
  });

  test("gx_review can ask for a whole-repo review", async () => {
    // Without this, a dirty working tree makes the diff the review subject, so
    // an agent asking about the codebase gets an answer scoped to the diff.
    const output = await gxReview({ cwd: repoRoot, repo: true });
    const parsed = JSON.parse(output);
    expect(parsed.command).toEqual([process.env.GX_BINARY, "review", "--no-publish", "--repo"]);

    // Omitting it keeps the patch-focused default.
    const plain = JSON.parse(await gxReview({ cwd: repoRoot }));
    expect(plain.command).toEqual([process.env.GX_BINARY, "review", "--no-publish"]);
  });

  test("gx_review never publishes, whatever it is asked for", async () => {
    // readOnlyHint is a promise to the calling agent. Without --no-publish,
    // `gx review` posts a comment on the matching GitHub PR and records the
    // run to gx Cloud — a write other people see, from a tool the agent was
    // told is safe to call freely. No argument combination may drop it.
    const variants = [
      {},
      { repo: true },
      { deep: true, verbose: true },
      { scope: "security" as const, focus: "internal", prompt: "check auth" },
    ];
    for (const variant of variants) {
      const parsed = JSON.parse(await gxReview({ cwd: repoRoot, ...variant }));
      expect(parsed.command).toContain("--no-publish");
    }
  });

  test("gx_review never initializes the repo it reviews", async () => {
    // The mock repo has no `.gx-initialized` marker, so anything that probed
    // or repaired gx state would show up in the call log. readOnlyHint is only
    // honest if `review` is the single command the tool runs.
    const output = await gxReview({ cwd: repoRoot });
    expect(JSON.parse(output).ok).toBe(true);

    const calls = (await readFile(callLog, "utf8")).trim().split("\n");
    expect(calls).toHaveLength(1);
    expect(calls[0]).toContain("|review");
    expect(calls.some((line) => line.includes("|init"))).toBe(false);
    expect(calls.some((line) => line.includes("|doctor"))).toBe(false);
    expect(existsSync(join(repoRoot, ".gx-initialized"))).toBe(false);
  });

  test("gx_review surfaces MCP auth login guidance", async () => {
    process.env.GX_MOCK_AUTH_ERROR = "1";

    const output = await gxReview({ cwd: repoRoot });
    const parsed = JSON.parse(output);

    expect(parsed.ok).toBe(false);
    expect(parsed.action).toBe("review");
    expect(parsed.auth_required).toBe(true);
    expect(parsed.display).toBe(
      "gx cloud authentication is required. Run `gx auth login` in a terminal, then retry the MCP tool.",
    );
    expect(parsed.stderr).toContain("run `gx auth login` in a terminal");
    expect(parsed.next_actions[0]).toBe("Run `gx auth login` in a terminal, then retry the MCP tool.");
  });
});
