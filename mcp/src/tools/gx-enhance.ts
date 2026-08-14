import { z } from "zod";
import type { InferSchema, ToolMetadata } from "xmcp";
import { formatError, formatResult, runTt } from "../gx";

const enhanceScope = z.enum([
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
  scope: enhanceScope.optional().describe("Optional focused scope. Omit for gx enhance's patch-focused default."),
  focus: z.string().optional().describe("Limit the run to files under this path prefix."),
  prompt: z.string().optional().describe("Optional steering prompt passed as the positional gx enhance prompt."),
  repo: z
    .boolean()
    .optional()
    .describe(
      "Assess the whole repository instead of just the current change. Use this to ask about the codebase itself, or when uncommitted work would otherwise narrow the run to the diff.",
    ),
  deep: z.boolean().optional().describe("Run full-spectrum analysis with more local and indexed context."),
  fast: z
    .boolean()
    .optional()
    .describe(
      "Optimize for wall clock: one reviewer instead of two, no verification pass, no project test suite, and findings written without code examples. Use for interactive runs someone is waiting on; omit when thoroughness matters more than latency.",
    ),
  verbose: z.boolean().optional().describe("Include repo facts, docs, and changed files."),
};

export const metadata: ToolMetadata = {
  name: "gx_enhance",
  description:
    "Run gx enhance to get issues and tips for improving the code you are currently working on, grounded in local facts, patch facts, previous-session context, PR/code-change context, indexed review resources, and configured AI reviewers. Defaults to enhancing the current change; pass repo=true to assess the whole repository instead, which is what questions about the codebase need when the working tree is dirty. Pass deep=true for full-spectrum analysis. MCP forces GX_REVIEW_AI=1 for this command.",
  annotations: {
    title: "gx Enhance",
    readOnlyHint: true,
    destructiveHint: false,
    idempotentHint: false,
  },
};

export default async function gxEnhance(params: InferSchema<typeof schema>) {
  // --no-comment is not optional here. Without it `gx enhance` posts a comment
  // on the matching GitHub PR, so a tool annotated readOnlyHint would write to
  // a pull request other people read — and an agent calling it for context
  // mid-codegen would comment on the PR every time it asked a question.
  // Commenting is a deliberate act that belongs to the CLI, where a human
  // typed the command. The run itself still records to gx Cloud history
  // (unlike --no-publish, which suppresses that too) so MCP runs count in
  // per-surface usage; --client mcp labels the row.
  const args = ["enhance", "--no-comment", "--client", "mcp"];
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
    // No init, no hooks, no gx state, and no PR comment: `gx enhance
    // --no-comment` reads the repo and reports back, which is what
    // readOnlyHint promises the caller.
    return formatResult(await runTt(args, { cwd: params.cwd, timeoutMs: 300_000 }), {
      action: "enhance",
      nextActions: ["Use findings as context before committing or pushing changes."],
    });
  } catch (error) {
    return formatError(error, { action: "enhance" });
  }
}
