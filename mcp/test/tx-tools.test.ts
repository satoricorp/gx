import { afterEach, beforeEach, describe, expect, test } from "bun:test";
import { existsSync } from "node:fs";
import { mkdir, mkdtemp, readFile, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import tlReview, { metadata as reviewMetadata, schema as reviewSchema } from "../src/tools/tl-review";

describe("tl_review metadata and schema", () => {
  test("describes the review surface", () => {
    expect(reviewMetadata.name).toBe("tl_review");
    expect(reviewMetadata.description).toMatch(/local facts/i);
    expect(reviewMetadata.annotations?.readOnlyHint).toBe(true);
    expect(reviewSchema.scope.parse("architecture")).toBe("architecture");
    expect(reviewSchema.prompt.parse("review auth rollback risk")).toBe("review auth rollback risk");
    expect(reviewSchema.deep.parse(true)).toBe(true);
    expect(reviewSchema.repo.parse(true)).toBe(true);
  });
});

describe("registered MCP tool names", () => {
  test("exposes exactly tl_review", () => {
    const registered = [reviewMetadata.name].sort();
    expect(registered).toEqual(["tl_review"]);
    expect(registered).not.toContain("tl_commit");
    expect(registered).not.toContain("tl_status");
    expect(registered).not.toContain("tl_push");
    expect(registered).not.toContain("tl_publish");
    expect(registered).not.toContain("tl_edit");
  });
});

describe("tl MCP CLI invocation", () => {
  let mockDir: string;
  let repoRoot: string;
  let callLog: string;
  let previousEnv: Record<string, string | undefined>;

  const envKeys = ["TOTALITY_BINARY", "TOTALITY_MOCK_LOG", "TOTALITY_MOCK_AUTH_ERROR"];

  beforeEach(async () => {
    mockDir = await mkdtemp(join(tmpdir(), "tl-mcp-tools-"));
    repoRoot = join(mockDir, "repo");
    callLog = join(mockDir, "calls.log");
    await mkdir(join(repoRoot, ".git"), { recursive: true });
    await Bun.spawn(["git", "init", repoRoot]).exited;
    previousEnv = Object.fromEntries(envKeys.map((key) => [key, process.env[key]]));

    const mockTt = join(mockDir, "mock-tl.sh");
    await writeFile(
      mockTt,
      `#!/bin/sh
printf 'tl|%s|%s|%s\\n' "$PWD" "$TOTALITY_REVIEW_AI" "$*" >> "$TOTALITY_MOCK_LOG"
if [ "$1" = init ]; then
  touch "$PWD/.totality-initialized"
  echo "initialized"
  exit 0
fi
if [ "$1" = review ]; then
  if [ "$TOTALITY_MOCK_AUTH_ERROR" = "1" ]; then
    echo 'github token is not configured for MCP: run \`tl auth login\` in a terminal, then retry the MCP tool' >&2
    exit 1
  fi
  echo "review ok"
  exit 0
fi
if [ "$1" = doctor ]; then
  if [ ! -f "$PWD/.totality-initialized" ]; then
    echo "tl not initialized" >&2
    exit 1
  fi
  echo '{"doctor":{}}'
  exit 0
fi
echo "unexpected tl args: $@" >&2
exit 1
`,
      "utf8",
    );
    await Bun.spawn(["chmod", "+x", mockTt]).exited;

    process.env.TOTALITY_BINARY = mockTt;
    process.env.TOTALITY_MOCK_LOG = callLog;
    delete process.env.TOTALITY_MOCK_AUTH_ERROR;
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

  test("tl_review passes scope, focus, prompt, deep, and verbose flags", async () => {
    const output = await tlReview({
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
      process.env.TOTALITY_BINARY,
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
    expect(calls).toContain("tl|");
    expect(calls).toContain(
      "|1|review --no-publish --scope architecture --focus internal/authoring --deep --verbose review auth rollback risk",
    );
  });

  test("tl_review can ask for a whole-repo review", async () => {
    // Without this, a dirty working tree makes the diff the review subject, so
    // an agent asking about the codebase gets an answer scoped to the diff.
    const output = await tlReview({ cwd: repoRoot, repo: true });
    const parsed = JSON.parse(output);
    expect(parsed.command).toEqual([process.env.TOTALITY_BINARY, "review", "--no-publish", "--repo"]);

    // Omitting it keeps the patch-focused default.
    const plain = JSON.parse(await tlReview({ cwd: repoRoot }));
    expect(plain.command).toEqual([process.env.TOTALITY_BINARY, "review", "--no-publish"]);
  });

  test("tl_review never publishes, whatever it is asked for", async () => {
    // readOnlyHint is a promise to the calling agent. Without --no-publish,
    // `tl review` posts a comment on the matching GitHub PR and records the
    // run to Totality Cloud — a write other people see, from a tool the agent was
    // told is safe to call freely. No argument combination may drop it.
    const variants = [
      {},
      { repo: true },
      { deep: true, verbose: true },
      { scope: "security" as const, focus: "internal", prompt: "check auth" },
    ];
    for (const variant of variants) {
      const parsed = JSON.parse(await tlReview({ cwd: repoRoot, ...variant }));
      expect(parsed.command).toContain("--no-publish");
    }
  });

  test("tl_review never initializes the repo it reviews", async () => {
    // The mock repo has no `.totality-initialized` marker, so anything that probed
    // or repaired Totality state would show up in the call log. readOnlyHint is only
    // honest if `review` is the single command the tool runs.
    const output = await tlReview({ cwd: repoRoot });
    expect(JSON.parse(output).ok).toBe(true);

    const calls = (await readFile(callLog, "utf8")).trim().split("\n");
    expect(calls).toHaveLength(1);
    expect(calls[0]).toContain("|review");
    expect(calls.some((line) => line.includes("|init"))).toBe(false);
    expect(calls.some((line) => line.includes("|doctor"))).toBe(false);
    expect(existsSync(join(repoRoot, ".totality-initialized"))).toBe(false);
  });

  test("tl_review surfaces MCP auth login guidance", async () => {
    process.env.TOTALITY_MOCK_AUTH_ERROR = "1";

    const output = await tlReview({ cwd: repoRoot });
    const parsed = JSON.parse(output);

    expect(parsed.ok).toBe(false);
    expect(parsed.action).toBe("review");
    expect(parsed.auth_required).toBe(true);
    expect(parsed.display).toBe(
      "Totality cloud authentication is required. Run `tl auth login` in a terminal, then retry the MCP tool.",
    );
    expect(parsed.stderr).toContain("run `tl auth login` in a terminal");
    expect(parsed.next_actions[0]).toBe("Run `tl auth login` in a terminal, then retry the MCP tool.");
  });
});
