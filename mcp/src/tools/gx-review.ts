import { z } from "zod";
import type { InferSchema, ToolMetadata } from "xmcp";
import { formatError, formatResult, runTt } from "../gx";

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
    title: "gx Review",
    readOnlyHint: true,
    destructiveHint: false,
    idempotentHint: false,
  },
};

export default async function gxReview(params: InferSchema<typeof schema>) {
  // --no-comment is not optional here. Without it `gx review` posts a review
  // comment on the matching GitHub PR, so a tool annotated readOnlyHint would
  // write to a pull request other people read — and an agent calling it for
  // context mid-codegen would comment on the PR every time it asked a
  // question. Commenting is a deliberate act that belongs to the CLI, where a
  // human typed the command. The run itself still records to gx Cloud history
  // (unlike --no-publish, which suppresses that too) so MCP reviews count in
  // per-surface usage; --client mcp labels the row.
  //
  // There is deliberately no `fast` parameter. It used to be exposed here, and
  // its description told the caller to use it "for interactive reviews someone
  // is waiting on" — which is every review that arrives through MCP, so the
  // model selected it essentially always, correctly per the description.
  //
  // What it selected was a review with no verification pass, and `kind` is
  // written by the judge: measured across 41 saved reports, 162 of 166
  // unjudged findings carry no kind at all, against 36 of 39 judged ones that
  // do. So --fast silently removes the field findings.go documents as "the
  // strongest single noise separator the review produces", while the report
  // still renders strength and kind as though they meant something. A user
  // triaging one such run by confidence found every Strong/defect finding
  // false and all three real ones filed under "worth exploring".
  //
  // --fast still exists on the CLI, where a human choosing it knows the trade.
  const args = ["review", "--no-comment", "--client", "mcp"];
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
    // No init, no hooks, no gx state, and no PR comment: `gx review
    // --no-comment` reads the repo and reports back, which is what
    // readOnlyHint promises the caller.
    return formatResult(await runTt(args, { cwd: params.cwd, timeoutMs: 300_000 }), {
      action: "review",
      nextActions: ["Use findings as context before committing or pushing changes."],
    });
  } catch (error) {
    return formatError(error, { action: "review" });
  }
}
