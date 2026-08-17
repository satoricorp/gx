import { z } from "zod";
import type { InferSchema, ToolMetadata } from "xmcp";
import { formatError, formatResult, runTt } from "../gx";

export const schema = {
  cwd: z.string().optional().describe("Repository working directory. Defaults to the MCP server process cwd."),
  prompt: z
    .string()
    .optional()
    .describe(
      "Optional intent hint passed as the positional gx gates argument: what the change is meant to do, so the back-pressure gate can judge scope against it.",
    ),
  verbose: z.boolean().optional().describe("List the source files behind each gate."),
};

export const metadata: ToolMetadata = {
  name: "gx_gates",
  description:
    "Run the gx gates exit gate on the current change: six gate checks — correctness, security, code health, back-pressure, accessibility, performance — each PASS/FAIL/SKIPPED, plus a ship/no-ship verdict. Deterministic checks (tests, linters, secrets scan, dependency audits) decide what they can; one AI judgment call covers the rest. Run it before opening a PR. A no-ship verdict is a successful run, not an error.",
  annotations: {
    title: "gx Gates",
    readOnlyHint: true,
    destructiveHint: false,
    idempotentHint: false,
  },
};

export default async function gxGates(params: InferSchema<typeof schema>) {
  // --report-only is not optional here. `gx gates` exits non-zero on a
  // no-ship verdict — it IS the gate — but runTt treats any non-zero exit as a
  // tool failure, and a no-ship verdict is the tool doing its job. The verdict
  // travels in the markdown report (--md), whose last two lines are the
  // machine seam: "Verdict: ..." and "Next: ...". --client mcp labels the run.
  const args = ["gates", "--report-only", "--md", "--client", "mcp"];
  if (params.verbose) {
    args.push("--verbose");
  }
  const prompt = params.prompt?.trim();
  if (prompt) {
    args.push(prompt);
  }
  try {
    return formatResult(await runTt(args, { cwd: params.cwd, timeoutMs: 360_000 }), {
      action: "gates",
      nextActions: [
        "Relay the full per-gate report. If any gate FAILED, suggest a concrete fix per failed gate; if every gate passed or was skipped, tell the user the change is clear to ship and ask whether to continue.",
      ],
    });
  } catch (error) {
    return formatError(error, { action: "gates" });
  }
}
