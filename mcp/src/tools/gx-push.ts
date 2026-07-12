import { z } from "zod";
import type { InferSchema, ToolMetadata } from "xmcp";
import { formatError, formatResult, runGx } from "../gx";
import { ensureGxInitialized } from "../session-workspace";

export const schema = {
  cwd: z.string().optional().describe("Repository working directory. Defaults to the MCP server process cwd."),
  stack: z
    .string()
    .optional()
    .describe(
      "Optional stack/feature name to push. Omit to push the current stack, like git push. Pass all: true to push every accepted stack.",
    ),
  all: z.boolean().optional().describe("Push every accepted stack instead of only the current one."),
};

export const metadata: ToolMetadata = {
  name: "gx_push",
  description:
    "Push GX stacks with gx push. Use after gx_commit/gx_status in the save workflow unless the user asked to keep work local. Syncs code to GitHub/origin, uploads sessions plus GX metadata, and uses guarded lease pushes when edited revisions rewrite a stack.",
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
  } else if (params.all) {
    args.push("--all");
  }
  try {
    await ensureGxInitialized(params.cwd);
    return formatResult(await runGx(args, { cwd: params.cwd, timeoutMs: 300_000 }), {
      action: "push",
      nextActions: ["Run gx_status to verify remote state."],
    });
  } catch (error) {
    return formatError(error, {
      action: "push",
      nextActions: ["Run gx_status to inspect local and remote stack state, then retry gx_push."],
    });
  }
}
