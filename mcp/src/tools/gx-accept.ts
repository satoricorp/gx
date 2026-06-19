import { z } from "zod";
import type { InferSchema, ToolMetadata } from "xmcp";
import { formatError, formatJsonResult, runGxJson } from "../gx";
import { resolveSessionWorkspace } from "../session-workspace";

export const schema = {
  proposal_id: z.string().min(1).describe("GX compose proposal id returned by gx_compose."),
  cwd: z.string().optional().describe("Repository working directory. Defaults to the MCP server process cwd."),
  allow_warnings: z.boolean().optional().describe("Apply despite blocking feasibility warnings after accepting the risk."),
  session_id: z.string().optional().describe("Explicit captured GX session id for the session-isolated workspace."),
  session_ids: z.array(z.string()).optional().describe("Additional captured GX session ids to include as provenance."),
};

export const metadata: ToolMetadata = {
  name: "gx_accept",
  description:
    "Accept a ready GX compose proposal into stacks with gx compose apply. Use when gx_compose was run with auto_accept=false or a user explicitly asks to accept a pending proposal.",
  annotations: {
    title: "GX Accept",
    readOnlyHint: false,
    destructiveHint: false,
    idempotentHint: false,
  },
};

export default async function gxAccept(params: InferSchema<typeof schema>) {
  try {
    const sessionWorkspace = await resolveSessionWorkspace({
      cwd: params.cwd,
      sessionId: params.session_id,
      sessionIds: params.session_ids,
    });
    const args = ["compose", "apply", params.proposal_id, "--json"];
    if (params.allow_warnings) {
      args.push("--allow-warnings");
    }
    const output = JSON.parse(
      formatJsonResult(
        await runGxJson(args, {
          cwd: sessionWorkspace.cwd,
          sessionIds: sessionWorkspace.sessionIds,
        }),
        {
          action: "accept",
          nextActions: ["Run gx_publish to publish accepted stacks."],
        },
      ),
    ) as Record<string, unknown>;
    output.session_workspace = sessionWorkspace.workspace;
    return JSON.stringify(output, null, 2);
  } catch (error) {
    return formatError(error, { action: "accept" });
  }
}
