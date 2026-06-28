import { z } from "zod";
import type { InferSchema, ToolMetadata } from "xmcp";
import { formatError, formatJsonResult, runGxJson } from "../gx";
import { resolveSessionWorkspace } from "../session-workspace";

export const schema = {
  cwd: z.string().optional().describe("Repository working directory. Defaults to the MCP server process cwd."),
  intent: z.string().optional().describe("Optional intent prefix for generated revisions."),
  filesets: z.array(z.string()).optional().describe("Optional changed files or directories to include."),
  exclude: z.array(z.string()).optional().describe("Optional changed files or directories to exclude."),
  session_id: z.string().optional().describe("Explicit captured GX session id to include as revision provenance."),
  session_ids: z.array(z.string()).optional().describe("Additional captured GX session ids to include as provenance."),
};

export const metadata: ToolMetadata = {
  name: "gx_generate",
  description:
    "Generate local GX features and revisions from current work with gx generate. GX applies safe generated revisions automatically, can append to semantically similar stacks, and uses a longer repair budget for MCP.",
  annotations: {
    title: "GX Generate",
    readOnlyHint: false,
    destructiveHint: false,
    idempotentHint: false,
  },
};

export default async function gxGenerate(params: InferSchema<typeof schema>) {
  try {
    const sessionWorkspace = await resolveSessionWorkspace({
      cwd: params.cwd,
      sessionId: params.session_id,
      sessionIds: params.session_ids,
    });
    const args = ["generate", "--json"];
    if (params.intent) {
      args.push("--intent", params.intent);
    }
    for (const excluded of params.exclude ?? []) {
      args.push("--exclude", excluded);
    }
    for (const fileset of params.filesets ?? []) {
      args.push(fileset);
    }
    const output = JSON.parse(
      formatJsonResult(
        await runGxJson(args, {
          cwd: sessionWorkspace.cwd,
          sessionIds: sessionWorkspace.sessionIds,
          timeoutMs: 300_000,
        }),
        {
          action: "generate",
          nextActions: ["Run gx_status to inspect generated features, then gx_push to push ready work."],
        },
      ),
    ) as Record<string, unknown>;
    output.session_workspace = sessionWorkspace.workspace;
    output.routing_decisions = collectRoutingDecisions(output.result);
    return JSON.stringify(output, null, 2);
  } catch (error) {
    return formatError(error, { action: "generate" });
  }
}

function collectRoutingDecisions(result: unknown) {
  const revisions = proposalRevisions(result);
  if (revisions.length === 0) {
    return undefined;
  }
  return revisions
    .map((revision) => {
      const id = stringField(revision, "id");
      const intent = stringField(revision, "intent");
      const targetStack = stringField(revision, "target_stack");
      if (!id && !targetStack) {
        return undefined;
      }
      return {
        revision_id: id || undefined,
        intent: intent || undefined,
        target_stack: targetStack || undefined,
        base_stack: stringField(revision, "base_stack") || undefined,
        route_source: stringField(revision, "route_source") || undefined,
        route_reason: stringField(revision, "route_reason") || undefined,
        route_confidence: numberField(revision, "route_confidence"),
      };
    })
    .filter(Boolean);
}

function proposalRevisions(result: unknown): Record<string, unknown>[] {
  if (!isRecord(result)) {
    return [];
  }
  const proposal = isRecord(result.proposal) ? result.proposal : result;
  const revisions = proposal.revisions;
  if (!Array.isArray(revisions)) {
    return [];
  }
  return revisions.filter(isRecord);
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function stringField(value: Record<string, unknown>, key: string) {
  const field = value[key];
  return typeof field === "string" ? field : "";
}

function numberField(value: Record<string, unknown>, key: string) {
  const field = value[key];
  return typeof field === "number" ? field : undefined;
}
