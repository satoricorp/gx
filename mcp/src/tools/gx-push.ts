import { z } from "zod";
import type { InferSchema, ToolMetadata } from "xmcp";
import { formatError, formatResult, runGx } from "../gx";
import { ensureGxInitialized } from "../session-workspace";

export const schema = {
  cwd: z.string().optional().describe("Repository working directory. Defaults to the MCP server process cwd."),
  stack: z.string().optional().describe("Optional stack/feature name to push. Omit to push all eligible generated work."),
};

export const metadata: ToolMetadata = {
  name: "gx_push",
  description:
    "Push generated GX features with gx push. Use after gx_generate/gx_status in the GX save workflow unless the user asked to keep work local. Syncs code to GitHub/origin, uploads sessions plus GX metadata, and uses guarded lease pushes when edited revisions rewrite a stack.",
  annotations: {
    title: "GX Push",
    readOnlyHint: false,
    destructiveHint: false,
    idempotentHint: false,
  },
};

export default async function gxPush(params: InferSchema<typeof schema>) {
  const args = ["push"];
  if (params.stack) {
    args.push(params.stack);
  }
  try {
    await ensureGxInitialized(params.cwd);
    return formatResult(await runGx(args, { cwd: params.cwd, timeoutMs: 300_000 }), {
      action: "push",
      nextActions: ["Run gx_status to verify remote state. Run gx_sync after remote merges land."],
    });
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error);
    const nextActions = message.includes("gx sync")
      ? ["Run gx_sync, resolve any divergence, then retry gx_push for that stack."]
      : ["Run gx_status to inspect local and remote stack state."];
    return formatError(error, { action: "push", nextActions });
  }
}
