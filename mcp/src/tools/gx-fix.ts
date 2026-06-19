import { z } from "zod";
import type { InferSchema, ToolMetadata } from "xmcp";
import { formatError, formatJsonResult, runGxJson } from "../gx";
import { resolveSessionWorkspace } from "../session-workspace";

export const schema = {
  proposal_id: z.string().min(1).describe("GX compose proposal id to repair."),
  cwd: z.string().optional().describe("Repository working directory. Defaults to the MCP server process cwd."),
  model: z.string().optional().describe("Optional OpenAI model for compose AI repair."),
  max_warnings: z.number().int().positive().optional().describe("Maximum warning diagnostics to send to the model."),
  deterministic_only: z.boolean().optional().describe("Only run deterministic repair; defaults to false so LLM repair is enabled."),
  session_id: z.string().optional().describe("Explicit captured GX session id for the session-isolated workspace."),
  session_ids: z.array(z.string()).optional().describe("Additional captured GX session ids to include as provenance."),
};

export const metadata: ToolMetadata = {
  name: "gx_fix",
  description:
    "Repair a GX compose proposal with gx compose fix. Defaults to deterministic repair plus LLM repair, using user OpenAI tokens first and GX cloud fallback when available.",
  annotations: {
    title: "GX Fix Compose",
    readOnlyHint: false,
    destructiveHint: false,
    idempotentHint: false,
  },
};

export default async function gxFix(params: InferSchema<typeof schema>) {
  try {
    const sessionWorkspace = await resolveSessionWorkspace({
      cwd: params.cwd,
      sessionId: params.session_id,
      sessionIds: params.session_ids,
    });
    const args = ["compose", "fix", params.proposal_id, "--json"];
    if (params.deterministic_only) {
      args.push("--plan");
    }
    if (params.model) {
      args.push("--model", params.model);
    }
    if (params.max_warnings) {
      args.push("--max-warnings", String(params.max_warnings));
    }
    const output = JSON.parse(
      formatJsonResult(
        await runGxJson(args, {
          cwd: sessionWorkspace.cwd,
          sessionIds: sessionWorkspace.sessionIds,
        }),
        {
          action: "fix",
          nextActions: ["If the proposal is ready, call gx_accept with the proposal id."],
        },
      ),
    ) as Record<string, unknown>;
    output.session_workspace = sessionWorkspace.workspace;
    return JSON.stringify(output, null, 2);
  } catch (error) {
    return formatError(error, { action: "fix" });
  }
}
