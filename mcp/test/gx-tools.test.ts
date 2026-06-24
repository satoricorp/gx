import { afterEach, beforeEach, describe, expect, test } from "bun:test";
import { mkdir, mkdtemp, readFile, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import gxPublish, { metadata as publishMetadata, schema as publishSchema } from "../src/tools/gx-publish";
import gxReview, { metadata as reviewMetadata, schema as reviewSchema } from "../src/tools/gx-review";

describe("gx_review metadata and schema", () => {
  test("describes the review surface", () => {
    expect(reviewMetadata.name).toBe("gx_review");
    expect(reviewMetadata.description).toMatch(/local facts/i);
    expect(reviewMetadata.annotations?.readOnlyHint).toBe(true);
    expect(reviewSchema.scope.parse("architecture")).toBe("architecture");
    expect(reviewSchema.deep.parse(true)).toBe(true);
  });
});

describe("gx_publish metadata and schema", () => {
  test("describes publishing accepted stacks", () => {
    expect(publishMetadata.name).toBe("gx_publish");
    expect(publishMetadata.description).toMatch(/accepted GX stacks/i);
    expect(publishMetadata.annotations?.readOnlyHint).toBe(false);
    expect(publishSchema.stack.parse("feature/review")).toBe("feature/review");
  });
});

describe("gx_review and gx_publish CLI invocation", () => {
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
if [ "$1" = publish ]; then
  if [ "$GX_MOCK_AUTH_ERROR" = "1" ]; then
    echo 'github token is not configured for MCP: run \`gx auth login\` in a terminal, then retry the MCP tool' >&2
    exit 1
  fi
  echo "publish ok"
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

  test("gx_review passes scope, focus, deep, and verbose flags", async () => {
    const output = await gxReview({
      cwd: repoRoot,
      scope: "architecture",
      focus: "internal/authoring",
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
    ]);

    const calls = await readFile(callLog, "utf8");
    expect(calls).toContain("jj|");
    expect(calls).toContain("gx|");
    expect(calls).toContain("|1|review --scope architecture --focus internal/authoring --deep --verbose");
  });

  test("gx_publish passes an optional stack name", async () => {
    const output = await gxPublish({ cwd: repoRoot, stack: "feature/review" });
    const parsed = JSON.parse(output);
    expect(parsed.action).toBe("publish");
    expect(parsed.display).toBe("publish ok");
    expect(parsed.command).toEqual([process.env.GX_BINARY, "publish", "feature/review"]);
    expect(parsed.next_actions).toEqual(["Run gx_sync after GitHub merges land."]);
  });

  test("gx_publish surfaces MCP auth login guidance", async () => {
    process.env.GX_MOCK_AUTH_ERROR = "1";

    const output = await gxPublish({ cwd: repoRoot });
    const parsed = JSON.parse(output);

    expect(parsed.ok).toBe(false);
    expect(parsed.action).toBe("publish");
    expect(parsed.auth_required).toBe(true);
    expect(parsed.display).toBe(
      "GX cloud authentication is required. Run `gx auth login` in a terminal, then retry the MCP tool.",
    );
    expect(parsed.stderr).toContain("run `gx auth login` in a terminal");
    expect(parsed.next_actions[0]).toBe("Run `gx auth login` in a terminal, then retry the MCP tool.");
  });
});
