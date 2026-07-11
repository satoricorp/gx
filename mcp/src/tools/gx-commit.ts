import { mkdtemp, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { z } from "zod";
import type { InferSchema, ToolMetadata } from "xmcp";
import { formatError, formatResult, runGx } from "../gx";
import { ensureGxInitialized } from "../session-workspace";

export const schema = {
  cwd: z.string().optional().describe("Repository working directory. Defaults to the MCP server process cwd."),
  message: z.string().describe("Commit message for the staged revision."),
  branch: z.string().optional().describe("Optional new branch to create at HEAD and record onto. Default: the revision lands on the current branch and HEAD does not move, like git commit."),
  task_summary: z.string().optional().describe("Optional agent-declared summary of the work being committed."),
  commands_run: z.array(z.string()).optional().describe("Optional shell commands the agent ran while doing this work."),
  tests_run: z.array(z.string()).optional().describe("Optional test commands the agent ran while doing this work."),
  session_id: z.string().optional().describe("Optional captured GX session id to attach as provenance."),
};

export const metadata: ToolMetadata = {
  name: "gx_commit",
  description:
    "Record staged Git changes as a GX revision with gx commit. Stage files with git add first. This is the default GX save verb; use gx_status afterward and gx_push when the stack is ready to publish.",
  annotations: {
    title: "GX Commit",
    readOnlyHint: false,
    destructiveHint: false,
    idempotentHint: false,
  },
};

export default async function gxCommit(params: InferSchema<typeof schema>) {
  const message = params.message.trim();
  if (!message) {
    return formatError(new Error("message is required"), {
      action: "commit",
      nextActions: ["Stage files with git add, then call gx_commit with a non-empty message."],
    });
  }

  const args = ["commit", "-m", message];
  const branch = params.branch?.trim();
  if (branch) {
    args.push("--branch", branch);
  }
  const context = {
    task_summary: params.task_summary?.trim() || undefined,
    commands_run: params.commands_run?.map((value) => value.trim()).filter(Boolean),
    tests_run: params.tests_run?.map((value) => value.trim()).filter(Boolean),
  };
  const hasContext =
    Boolean(context.task_summary) ||
    Boolean(context.commands_run?.length) ||
    Boolean(context.tests_run?.length);

  let contextPath: string | undefined;
  if (hasContext) {
    const dir = await mkdtemp(join(tmpdir(), "gx-commit-context-"));
    contextPath = join(dir, "context.json");
    await writeFile(contextPath, `${JSON.stringify(context, null, 2)}\n`, "utf8");
    args.push("--context-file", contextPath);
  }

  try {
    await ensureGxInitialized(params.cwd);
    return formatResult(
      await runGx(args, {
        cwd: params.cwd,
        sessionId: params.session_id,
        timeoutMs: 300_000,
      }),
      {
        action: "commit",
        nextActions: ["Run gx_status to inspect local stack state. Run gx_push when the feature is ready to publish."],
      },
    );
  } catch (error) {
    return formatError(error, {
      action: "commit",
      nextActions: ["Run gx_status to inspect staged files and stack state, then retry gx_commit after git add."],
    });
  }
}
