import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import type { CallToolResult } from "@modelcontextprotocol/sdk/types.js";
import { z } from "zod";
import lgtmReview, { metadata as reviewMetadata, schema as reviewSchema } from "./tools/lgtm-review";
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
  "lgtm MCP exposes lgtm_review, a read-only tool.",
  "It never initializes a repository; if the repo is not set up for lgtm yet, run lgtm init there first.",
  "Save work with plain Git: git add, then git commit. lgtm installs Git hooks that stamp each commit with its lgtm revision trailer and record it, so no lgtm-specific commit verb is needed.",
  "To amend the latest commit message or restage work, use git commit --amend and preserve the lgtm revision trailer.",
  "Publish with plain git push, then open the PR with gh pr create. The lgtm pre-push hook captures the agent session, links edits to the changed hunks, and publishes the lgtm metadata that becomes the PR summary. Do not run lgtm push or lgtm capture push; they bypass or suppress that hook.",
  "Run lgtm_review for better codegen context from local facts, previous sessions, PRs, and current code changes.",
  "Use git status for inspection.",
].join(" ");

const tools: ToolModule[] = [
  defineTool(reviewMetadata, reviewSchema, lgtmReview),
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
      name: "lgtm MCP",
      version: "0.1.0",
      description: "lgtm MCP stdio server",
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
