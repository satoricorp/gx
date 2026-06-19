import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import type { CallToolResult } from "@modelcontextprotocol/sdk/types.js";
import { z } from "zod";
import gxAccept, { metadata as acceptMetadata, schema as acceptSchema } from "./tools/gx-accept";
import gxCompose, { metadata as composeMetadata, schema as composeSchema } from "./tools/gx-compose-changes";
import gxFix, { metadata as fixMetadata, schema as fixSchema } from "./tools/gx-fix";
import gxPublish, { metadata as publishMetadata, schema as publishSchema } from "./tools/gx-publish";
import gxReview, { metadata as reviewMetadata, schema as reviewSchema } from "./tools/gx-review";
import gxSetBase, { metadata as setBaseMetadata, schema as setBaseSchema } from "./tools/gx-set-base";
import gxSync, { metadata as syncMetadata, schema as syncSchema } from "./tools/gx-sync";

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
  "GX MCP exposes gx_sync, gx_compose, gx_accept, gx_publish, gx_review, gx_fix, and gx_set_base.",
  "Run gx_sync before gx_compose when remote GitHub merges may have landed.",
  "Use gx_compose with action propose to run gx compose --json in a session-isolated workspace; ready proposals are accepted automatically by default.",
  "When compose returns warnings, repair hints, or proposal issues, revise the proposal with your LLM context and call gx_compose with action review-plan until ready.",
  "Accept held ready revisions with gx_accept or gx_compose action apply; it appends to existing target bookmarks by default and returns the checkout to the repo default branch.",
  "Run gx_publish after accepted stacks are ready to push for review.",
  "Run gx_review for better codegen context from local facts, previous sessions, PRs, and current code changes.",
  "Use gx_fix for compose proposal repair. Do not use a revision move tool; it is intentionally not exposed yet.",
  "Use gx_set_base only to return to the repository default branch.",
  "After propose, summarize the output for the user: state, stacks, revision intents, notable files, warnings, and hunk link evidence when present.",
  "Do not use git commit or git push for normal GX save flows unless the user explicitly asks for raw Git.",
].join(" ");

const tools: ToolModule[] = [
  defineTool(acceptMetadata, acceptSchema, gxAccept),
  defineTool(composeMetadata, composeSchema, gxCompose),
  defineTool(fixMetadata, fixSchema, gxFix),
  defineTool(publishMetadata, publishSchema, gxPublish),
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
      async (args, extra) => toCallToolResult(await tool.handler(args as Record<string, unknown>, extra)),
    );
  }

  const transport = new StdioServerTransport();
  await server.connect(transport);
}

main().catch((error) => {
  console.error(error);
  process.exit(1);
});
