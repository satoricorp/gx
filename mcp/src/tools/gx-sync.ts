import { z } from "zod";
import type { InferSchema, ToolMetadata } from "xmcp";
import { formatError, formatResult, runGx } from "../gx";
import { ensureGxInitialized } from "../session-workspace";

export const schema = {
  cwd: z.string().optional().describe("Repository working directory. Defaults to the MCP server process cwd."),
  remote: z.string().optional().describe("Optional remote name to sync."),
  capture: z.boolean().optional().describe("Drain pending capture staging uploads instead of syncing remote Git state."),
};

export const metadata: ToolMetadata = {
  name: "gx_sync",
  description:
    "Run gx sync before compose or publish work so remote GitHub merges are reflected locally. Use capture=true only to drain capture uploads.",
  annotations: {
    title: "GX Sync",
    readOnlyHint: false,
    destructiveHint: false,
    idempotentHint: false,
  },
};

export default async function gxSync(params: InferSchema<typeof schema>) {
  const args = ["sync"];
  if (params.capture) {
    args.push("--capture");
  }
  if (params.remote) {
    args.push(params.remote);
  }
  try {
    if (!params.capture) {
      await ensureGxInitialized(params.cwd);
    }
    return formatResult(await runGx(args, { cwd: params.cwd, timeoutMs: 300_000 }), {
      action: params.capture ? "sync.capture" : "sync",
      nextActions: params.capture ? undefined : ["After sync succeeds, run gx_compose for changes or gx_publish for accepted stacks."],
    });
  } catch (error) {
    return formatError(error, { action: params.capture ? "sync.capture" : "sync" });
  }
}
