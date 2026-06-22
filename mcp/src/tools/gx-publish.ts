import { z } from "zod";
import type { InferSchema, ToolMetadata } from "xmcp";
import { formatError, formatResult, runGx } from "../gx";
import { ensureGxInitialized } from "../session-workspace";

export const schema = {
  cwd: z.string().optional().describe("Repository working directory. Defaults to the MCP server process cwd."),
  stack: z.string().optional().describe("Optional stack name to publish. Omit to publish all accepted stacks."),
};

export const metadata: ToolMetadata = {
  name: "gx_publish",
  description:
    "Publish accepted GX stacks with gx publish. Use after gx_compose auto-accepts or gx_accept applies a ready proposal.",
  annotations: {
    title: "GX Publish",
    readOnlyHint: false,
    destructiveHint: false,
    idempotentHint: false,
  },
};

export default async function gxPublish(params: InferSchema<typeof schema>) {
  const args = ["publish"];
  if (params.stack) {
    args.push(params.stack);
  }
  try {
    await ensureGxInitialized(params.cwd);
    return formatResult(await runGx(args, { cwd: params.cwd, timeoutMs: 300_000 }), {
      action: "publish",
      nextActions: ["Run gx_sync after GitHub merges land."],
    });
  } catch (error) {
    return formatError(error, { action: "publish" });
  }
}
