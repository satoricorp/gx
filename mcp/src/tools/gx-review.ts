import { z } from "zod";
import type { InferSchema, ToolMetadata } from "xmcp";
import { formatError, formatResult, runGx } from "../gx";

const reviewScope = z.enum([
  "architecture",
  "security",
  "performance",
  "onboarding",
  "docs",
  "dependencies",
  "testing",
  "maintainability",
]);

export const schema = {
  cwd: z.string().optional().describe("Repository working directory. Defaults to the MCP server process cwd."),
  scope: reviewScope.optional().describe("Optional focused review scope. Omit for gx review's patch-focused default."),
  focus: z.string().optional().describe("Limit review to files under this path prefix."),
  prompt: z.string().optional().describe("Optional reviewer prompt passed as the positional gx review prompt."),
  repo: z
    .boolean()
    .optional()
    .describe(
      "Review the whole repository instead of just the current change. Use this to ask about the codebase itself, or when uncommitted work would otherwise narrow the review to the diff.",
    ),
  deep: z.boolean().optional().describe("Run full-spectrum review with more local and indexed context."),
  verbose: z.boolean().optional().describe("Include repo facts, docs, and changed files."),
};

export const metadata: ToolMetadata = {
  name: "gx_review",
  description:
    "Run gx review for local facts, patch facts, previous-session context, PR/code-change context, indexed review resources, and configured AI reviewers. Defaults to reviewing the current change; pass repo=true to review the whole repository instead, which is what questions about the codebase need when the working tree is dirty. Pass deep=true for full-spectrum review. MCP forces GX_REVIEW_AI=1 for this command.",
  annotations: {
    title: "GX Review",
    readOnlyHint: true,
    destructiveHint: false,
    idempotentHint: false,
  },
};

export default async function gxReview(params: InferSchema<typeof schema>) {
  const args = ["review"];
  if (params.scope) {
    args.push("--scope", params.scope);
  }
  if (params.focus) {
    args.push("--focus", params.focus);
  }
  if (params.repo) {
    args.push("--repo");
  }
  if (params.deep) {
    args.push("--deep");
  }
  if (params.verbose) {
    args.push("--verbose");
  }
  const prompt = params.prompt?.trim();
  if (prompt) {
    args.push(prompt);
  }
  try {
    // No init, no hooks, no GX state: `gx review` reads the repo and nothing
    // else, which is what readOnlyHint promises the caller.
    return formatResult(await runGx(args, { cwd: params.cwd, timeoutMs: 300_000 }), {
      action: "review",
      nextActions: ["Use findings as context before committing or pushing changes."],
    });
  } catch (error) {
    return formatError(error, { action: "review" });
  }
}
