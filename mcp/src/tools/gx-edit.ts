import { z } from "zod";
import type { InferSchema, ToolMetadata } from "xmcp";
import { formatError, formatResult, runGx } from "../gx";
import { ensureGxInitialized } from "../session-workspace";

export const schema = {
  cwd: z.string().optional().describe("Repository working directory. Defaults to the MCP server process cwd."),
  revision: z.string().describe("Revision change id, commit id, or short id to re-enter for further edits."),
};

export const metadata: ToolMetadata = {
  name: "gx_edit",
  description:
    "Return to an existing GX revision with gx edit <rev> so further work lands on that change. Use when amending or continuing a recorded revision (gx commit --amend is unsupported). Prefer gx_commit for new staged work. After editing, return to the base branch before normal gx_commit work.",
  annotations: {
    title: "GX Edit",
    readOnlyHint: false,
    destructiveHint: false,
    idempotentHint: false,
  },
};

export default async function gxEdit(params: InferSchema<typeof schema>) {
  const revision = params.revision.trim();
  if (!revision) {
    return formatError(new Error("revision is required"), {
      action: "edit",
      nextActions: ["Run gx_status to list revisions, then call gx_edit with a revision id."],
    });
  }

  try {
    await ensureGxInitialized(params.cwd);
    return formatResult(await runGx(["edit", revision], { cwd: params.cwd, timeoutMs: 120_000 }), {
      action: "edit",
      nextActions: [
        "Stage follow-up changes with git add, then gx_commit. Run gx_status to confirm the active revision. Return to the base branch before starting unrelated work.",
      ],
    });
  } catch (error) {
    return formatError(error, {
      action: "edit",
      nextActions: ["Run gx_status to find a valid revision id, then retry gx_edit."],
    });
  }
}
