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
  scope: reviewScope.optional().describe("Review scope. Defaults to gx review's architecture scope."),
  focus: z.string().optional().describe("Limit review to files under this path prefix."),
  deep: z.boolean().optional().describe("Retrieve more local and indexed context."),
  verbose: z.boolean().optional().describe("Include repo facts, docs, and changed files."),
};

export const metadata: ToolMetadata = {
  name: "gx_review",
  description:
    "Run gx review for local facts, previous-session context, PR/code-change context, and configured AI reviewers. MCP forces GX_REVIEW_AI=1 for this command.",
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
  if (params.deep) {
    args.push("--deep");
  }
  if (params.verbose) {
    args.push("--verbose");
  }
  try {
    return formatResult(await runGx(args, { cwd: params.cwd, timeoutMs: 300_000 }), {
      action: "review",
      nextActions: ["Use findings as context before composing or publishing changes."],
    });
  } catch (error) {
    return formatError(error, { action: "review" });
  }
}
