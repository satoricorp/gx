import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import type { CallToolResult } from "@modelcontextprotocol/sdk/types.js";
import { z } from "zod";
import gxGenerate, { metadata as generateMetadata, schema as generateSchema } from "./tools/gx-generate";
import gxPush, { metadata as pushMetadata, schema as pushSchema } from "./tools/gx-push";
import gxReview, { metadata as reviewMetadata, schema as reviewSchema } from "./tools/gx-review";
import gxSetBase, { metadata as setBaseMetadata, schema as setBaseSchema } from "./tools/gx-set-base";
import gxStatus, { metadata as statusMetadata, schema as statusSchema } from "./tools/gx-status";
import gxSync, { metadata as syncMetadata, schema as syncSchema } from "./tools/gx-sync";
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
  "GX MCP exposes gx_sync, gx_generate, gx_status, gx_push, gx_review, and gx_set_base.",
  "If a repo is not initialized for GX, MCP runs gx init non-interactively before repository tools continue.",
  "When the user says save work, save using gx, or save with gx, run the GX save workflow: gx_generate, gx_status, then gx_push for ready stacks unless the user asks to keep work local.",
  "Run gx_sync before gx_generate when remote GitHub merges may have landed.",
  "Use gx_generate to create local features and revisions; GX applies safe generated revisions automatically and may append to semantically similar stacks.",
  "Run gx_status after generation to inspect local/remote stack state, then gx_push when features are ready.",
  "If gx_push reports remote divergence, run gx_sync before retrying that stack.",
  "Run gx_review for better codegen context from local facts, previous sessions, PRs, and current code changes.",
  "Use gx_set_base only to return to the repository default branch.",
  "After generate, summarize the output for the user: features, revision intents, notable files, and warnings when present.",
  "Use GX MCP tools before the CLI for GX workflows; use the CLI only as fallback when MCP is unavailable.",
  "Do not use git commit or git push for normal GX save flows unless the user explicitly asks for raw Git. Prefer denying or requiring approval for raw git commit, push, reset, and branch deletion in clients that support tool policies.",
].join(" ");

const tools: ToolModule[] = [
  defineTool(generateMetadata, generateSchema, gxGenerate),
  defineTool(statusMetadata, statusSchema, gxStatus),
  defineTool(pushMetadata, pushSchema, gxPush),
  defineTool(reviewMetadata, reviewSchema, gxReview),
  defineTool(setBaseMetadata, setBaseSchema, gxSetBase),
  defineTool(syncMetadata, syncSchema, gxSync),
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
