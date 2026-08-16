import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import type { CallToolResult } from "@modelcontextprotocol/sdk/types.js";
import { z } from "zod";
import gxConstraints, { metadata as constraintsMetadata, schema as constraintsSchema } from "./tools/gx-constraints";
import gxReview, { metadata as reviewMetadata, schema as reviewSchema } from "./tools/gx-review";
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
  "gx MCP exposes gx_review and gx_constraints, both read-only tools.",
  "It never initializes a repository; if the repo is not set up for gx yet, run gx init there first.",
  "Save work with plain Git: git add, then git commit. gx installs Git hooks that stamp each commit with its gx revision trailer and record it, so no gx-specific commit verb is needed.",
  "To amend the latest commit message or restage work, use git commit --amend and preserve the gx revision trailer.",
  "Publish with plain git push, then open the PR with gh pr create. The gx pre-push hook captures the agent session, links edits to the changed hunks, and publishes the gx metadata that becomes the PR summary. Do not run gx push or gx capture push; they bypass or suppress that hook.",
  "Run gx_review for better codegen context from local facts, previous sessions, PRs, and current code changes.",
  "Run gx_constraints before shipping to check the change against the pre-ship constraint gates and relay its ship/no-ship verdict.",
  "Use git status for inspection.",
].join(" ");

const tools: ToolModule[] = [
  defineTool(reviewMetadata, reviewSchema, gxReview),
  defineTool(constraintsMetadata, constraintsSchema, gxConstraints),
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
      name: "gx MCP",
      version: "0.1.0",
      description: "gx MCP stdio server",
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
