import { z } from "zod";
import type { InferSchema, ToolMetadata } from "xmcp";
import { formatError, formatJsonResult, runGxJson } from "../gx";
import { ensureGxInitialized } from "../session-workspace";

export const schema = {
  cwd: z.string().optional().describe("Repository working directory. Defaults to the MCP server process cwd."),
  show_all: z.boolean().optional().describe("Include hidden compatibility state where the CLI supports it."),
};

export const metadata: ToolMetadata = {
  name: "gx_status",
  description: "Inspect GX unstaged files, local features, revisions, and remote state with gx status --json.",
  annotations: {
    title: "GX Status",
    readOnlyHint: true,
    destructiveHint: false,
    idempotentHint: true,
  },
};

export default async function gxStatus(params: InferSchema<typeof schema>) {
  const args = ["status", "--json"];
  if (params.show_all) {
    args.push("--show-all");
  }
  try {
    await ensureGxInitialized(params.cwd);
    return formatJsonResult(await runGxJson(args, { cwd: params.cwd, timeoutMs: 120_000 }), {
      action: "status",
      nextActions: ["Run gx_generate for local changes or gx_push for ready features."],
    });
  } catch (error) {
    return formatError(error, { action: "status" });
  }
}
