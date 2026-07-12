import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import type { CallToolResult } from "@modelcontextprotocol/sdk/types.js";
import { z } from "zod";
import gxCommit, { metadata as commitMetadata, schema as commitSchema } from "./tools/gx-commit";
import gxEdit, { metadata as editMetadata, schema as editSchema } from "./tools/gx-edit";
import gxReview, { metadata as reviewMetadata, schema as reviewSchema } from "./tools/gx-review";
import gxStatus, { metadata as statusMetadata, schema as statusSchema } from "./tools/gx-status";
import { withUpdateNotice } from "./update";

type ToolModule = {
  metadata: {
    name: string;
    description?: string;
    annotations?: Record<string, unknown>;
  };
  schema: Record<string, z.ZodType>;
  handler: (args: Record<string, unknown>, extra: unknown) => unknown | Promise<unknown>;
};

const instructions = [
  "GX MCP exposes gx_commit, gx_status, gx_edit, and gx_review.",
  "If a repo is not initialized for GX, MCP runs gx init non-interactively before repository tools continue.",
  "When the user says save work, save using gx, or save with gx, run the GX save workflow: git add, gx_commit, gx_status, then publish ready work with plain git push (the GX pre-push hook captures the session and publishes) unless the user asks to keep work local.",
  "gx_commit is the default verb. Use it after git add to record staged work as a GX revision.",
  "Use gx_edit only to re-enter an existing revision for further edits (gx commit --amend is unsupported).",
  "Run gx_status after commits or edits to inspect local/remote stack state before publishing.",
  "Run gx_review for better codegen context from local facts, previous sessions, PRs, and current code changes.",
  "Publish with plain git push, then open the PR with gh pr create. Do not run gx push or gx capture push; they bypass or suppress the pre-push hook that GX publishes through.",
  "Use GX MCP tools before the CLI for GX workflows; use the CLI only as fallback when MCP is unavailable.",
  "Do not use git commit for normal GX save flows unless the user explicitly asks for raw Git. Prefer denying or requiring approval for raw git commit, reset, and branch deletion in clients that support tool policies.",
].join(" ");

const tools: ToolModule[] = [
  defineTool(commitMetadata, commitSchema, gxCommit),
  defineTool(statusMetadata, statusSchema, gxStatus),
  defineTool(editMetadata, editSchema, gxEdit),
  defineTool(reviewMetadata, reviewSchema, gxReview),
];

function defineTool(
  metadata: ToolModule["metadata"],
  schema: ToolModule["schema"],
  handler: unknown,
): ToolModule {
  return {
    metadata,
    schema,
    handler: handler as ToolModule["handler"],
  };
}

function redirectConsoleToStderr() {
  const stderrConsole = new console.Console(process.stderr, process.stderr);
  for (const method of ["log", "debug", "info", "warn", "error", "dir", "table", "trace"] as const) {
    console[method] = stderrConsole[method].bind(stderrConsole) as never;
  }
}

function toCallToolResult(value: unknown): CallToolResult {
  if (typeof value === "string" || typeof value === "number") {
    return {
      content: [{ type: "text", text: String(value) }],
    };
  }
  if (value && typeof value === "object" && ("content" in value || "structuredContent" in value || "isError" in value)) {
    return value as CallToolResult;
  }
  return {
    content: [{ type: "text", text: JSON.stringify(value ?? null, null, 2) }],
  };
}

async function main() {
  redirectConsoleToStderr();
  const server = new McpServer(
    {
      name: "GX MCP",
      version: "0.1.0",
      description: "GX MCP stdio server",
    },
    { instructions },
  );

  for (const tool of tools) {
    server.registerTool(
      tool.metadata.name,
      {
        title: tool.metadata.annotations?.title as string | undefined,
        description: tool.metadata.description,
        inputSchema: z.object(tool.schema),
        annotations: tool.metadata.annotations,
      },
      async (args, extra) => toCallToolResult(await withUpdateNotice(await tool.handler(args as Record<string, unknown>, extra))),
    );
  }

  const transport = new StdioServerTransport();
  await server.connect(transport);
}

main().catch((error) => {
  console.error(error);
  process.exit(1);
});
