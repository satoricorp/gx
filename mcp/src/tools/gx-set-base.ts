import { z } from "zod";
import type { InferSchema, ToolMetadata } from "xmcp";
import { formatError, formatJsonResult, runGxJson } from "../gx";
import { ensureGxInitialized } from "../session-workspace";

type BaseResult = {
  DefaultBranch?: unknown;
  default_branch?: unknown;
};

export const schema = {
  cwd: z.string().optional().describe("Repository working directory. Defaults to the MCP server process cwd."),
};

export const metadata: ToolMetadata = {
  name: "gx_set_base",
  description:
    "Return the GX authoring base to the repository default branch. This tool intentionally does not accept an arbitrary branch target.",
  annotations: {
    title: "GX Set Base",
    readOnlyHint: false,
    destructiveHint: false,
    idempotentHint: true,
  },
};

export default async function gxSetBase(params: InferSchema<typeof schema>) {
  try {
    await ensureGxInitialized(params.cwd);
    const current = await runGxJson(["base", "--json"], { cwd: params.cwd });
    const branch = defaultBranchFromBase(current.json);
    if (!branch) {
      return formatError(new Error("gx base did not report a default branch"), { action: "set_base" });
    }
    return formatJsonResult(await runGxJson(["base", "--set", branch, "--json"], { cwd: params.cwd }), {
      action: "set_base",
      nextActions: ["Run gx_generate after returning to the default branch."],
    });
  } catch (error) {
    return formatError(error, { action: "set_base" });
  }
}

function defaultBranchFromBase(json: unknown) {
  const result = json as BaseResult;
  const value = result.DefaultBranch ?? result.default_branch;
  return typeof value === "string" ? value.trim() : "";
}
