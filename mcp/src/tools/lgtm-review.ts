import { z } from "zod";
import type { InferSchema, ToolMetadata } from "xmcp";
import { formatError, formatResult, runTt } from "../lgtm";

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
  scope: reviewScope.optional().describe("Optional focused review scope. Omit for lgtm review's patch-focused default."),
  focus: z.string().optional().describe("Limit review to files under this path prefix."),
  prompt: z.string().optional().describe("Optional reviewer prompt passed as the positional lgtm review prompt."),
  repo: z
    .boolean()
    .optional()
    .describe(
      "Review the whole repository instead of just the current change. Use this to ask about the codebase itself, or when uncommitted work would otherwise narrow the review to the diff.",
    ),
  deep: z.boolean().optional().describe("Run full-spectrum review with more local and indexed context."),
  fast: z
    .boolean()
    .optional()
    .describe(
      "Optimize for wall clock: one reviewer instead of two, no verification pass, no project test suite, and findings written without code examples. Use for interactive reviews someone is waiting on; omit when thoroughness matters more than latency.",
    ),
  verbose: z.boolean().optional().describe("Include repo facts, docs, and changed files."),
};

export const metadata: ToolMetadata = {
  name: "lgtm_review",
  description:
    "Run lgtm review for local facts, patch facts, previous-session context, PR/code-change context, indexed review resources, and configured AI reviewers. Defaults to reviewing the current change; pass repo=true to review the whole repository instead, which is what questions about the codebase need when the working tree is dirty. Pass deep=true for full-spectrum review. MCP forces LGTM_REVIEW_AI=1 for this command.",
  annotations: {
    title: "lgtm Review",
    readOnlyHint: true,
    destructiveHint: false,
    idempotentHint: false,
  },
};

export default async function lgtmReview(params: InferSchema<typeof schema>) {
  // --no-publish is not optional here. Without it `lgtm review` posts a review
  // comment on the matching GitHub PR and records the run to lgtm Cloud, so a
  // tool annotated readOnlyHint would write to a pull request other people
  // read — and an agent calling it for context mid-codegen would comment on
  // the PR every time it asked a question. Publishing is a deliberate act that
  // belongs to the CLI, where a human typed the command.
  const args = ["review", "--no-publish"];
  if (params.scope) {
    args.push("--scope", params.scope);
  }
  if (params.focus) {
    args.push("--focus", params.focus);
  }
  if (params.repo) {
    args.push("--repo");
  }
  if (params.fast) {
    args.push("--fast");
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
    // No init, no hooks, no lgtm state, and nothing published: `lgtm review
    // --no-publish` reads the repo and reports back, which is what
    // readOnlyHint promises the caller.
    return formatResult(await runTt(args, { cwd: params.cwd, timeoutMs: 300_000 }), {
      action: "review",
      nextActions: ["Use findings as context before committing or pushing changes."],
    });
  } catch (error) {
    return formatError(error, { action: "review" });
  }
}
