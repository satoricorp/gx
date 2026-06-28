import { z } from "zod";
import type { InferSchema, ToolMetadata } from "xmcp";
import { formatError, formatResult, runGx } from "../gx";
import { ensureGxInitialized } from "../session-workspace";

export const schema = {
  cwd: z.string().optional().describe("Repository working directory. Defaults to the MCP server process cwd."),
  remote: z.string().optional().describe("Optional remote name to sync."),
};

export const metadata: ToolMetadata = {
  name: "gx_sync",
  description: "Run gx sync before generate or push work so remote GitHub merges and GX remote state are reflected locally.",
  annotations: {
    title: "GX Sync",
    readOnlyHint: false,
    destructiveHint: false,
    idempotentHint: false,
  },
};

export default async function gxSync(params: InferSchema<typeof schema>) {
  const args = ["sync"];
  if (params.remote) {
    args.push(params.remote);
  }
  try {
    await ensureGxInitialized(params.cwd);
    return formatResult(await runGx(args, { cwd: params.cwd, timeoutMs: 300_000 }), {
      action: "sync",
      nextActions: ["After sync succeeds, run gx_generate for changes or gx_push for generated features."],
    });
  } catch (error) {
    return formatError(error, { action: "sync" });
  }
}
