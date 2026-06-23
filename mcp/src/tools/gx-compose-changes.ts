import { mkdtemp, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { z } from "zod";
import type { InferSchema, ToolMetadata } from "xmcp";
import { formatError, formatJsonResult, runGxJson, type GxJsonRunResult } from "../gx";
import { resolveSessionWorkspace, type SessionWorkspaceInfo } from "../session-workspace";

const composeAction = z.enum(["propose", "review-plan", "apply"]);

export const schema = {
  action: composeAction
    .optional()
    .describe(
      "Compose workflow step. propose: run gx compose --json on the working copy (default). review-plan: validate an LLM-authored proposal without applying. apply: accept a reviewed proposal into GX stacks, appending to existing target bookmarks and returning the checkout to the repo default branch.",
    ),
  cwd: z.string().optional().describe("Repository working directory. Defaults to the MCP server process cwd."),
  intent: z.string().optional().describe("Optional intent prefix for proposed revisions (propose only)."),
  filesets: z
    .array(z.string())
    .optional()
    .describe("Optional changed files or directories to include in compose (propose only)."),
  exclude: z
    .array(z.string())
    .optional()
    .describe("Optional changed files or directories to exclude from compose (propose only)."),
  proposal: z
    .record(z.string(), z.unknown())
    .optional()
    .describe(
      "LLM-authored GX compose proposal JSON (review-plan and apply). Preserve source proposal id/change ids. Assign every top-level hunk exactly once.",
    ),
  allow_warnings: z
    .boolean()
    .optional()
    .describe("Apply despite blocking feasibility warnings after intentionally accepting them (apply only)."),
  auto_accept: z
    .boolean()
    .optional()
    .describe("For action=propose, automatically accept an apply-ready compose proposal into GX stacks. Defaults to true."),
  session_id: z
    .string()
    .optional()
    .describe("Explicit captured GX session id to include as revision provenance."),
  session_ids: z
    .array(z.string())
    .optional()
    .describe("Explicit captured GX session ids to include as revision provenance."),
};

export const metadata: ToolMetadata = {
  name: "gx_compose",
  description:
    "Propose, review, or apply ordered GX revisions from a session-isolated JJ workspace via gx compose. If needed, MCP runs gx init before creating the workspace. action=propose layers working-copy changes into the pending compose proposal and auto-accepts apply-ready proposals by default; action=review-plan validates an LLM-authored plan; action=apply accepts a reviewed plan into stacks, appends to existing target bookmarks by default, and returns the checkout to the repo default branch. Requires a GX session id so concurrent MCP sessions do not share a dirty checkout. Run gx_sync before compose when remote merges may have landed. Use gx_review for local review.",
  annotations: {
    title: "GX Compose",
    readOnlyHint: false,
    destructiveHint: false,
    idempotentHint: false,
  },
};

export default async function gxCompose(params: InferSchema<typeof schema>) {
  const action = params.action ?? "propose";
  const requestedRunOptions = {
    cwd: params.cwd,
    sessionId: params.session_id,
    sessionIds: params.session_ids,
  };

  try {
    const sessionWorkspace = await resolveSessionWorkspace(requestedRunOptions);
    const runOptions = {
      cwd: sessionWorkspace.cwd,
      sessionIds: sessionWorkspace.sessionIds,
    };

    if (action === "propose") {
      const args = ["compose", "--json"];
      if (params.intent) {
        args.push("--intent", params.intent);
      }
      for (const excluded of params.exclude ?? []) {
        args.push("--exclude", excluded);
      }
      for (const fileset of params.filesets ?? []) {
        args.push(fileset);
      }
      const compose = await runGxJson(args, runOptions);
      if (params.auto_accept !== false) {
        const proposalID = proposalIDFromCompose(compose.json);
        if (proposalID && composeStateIsReady(compose.json)) {
          const accept = await runGxJson(["compose", "apply", proposalID, "--json"], runOptions);
          return formatComposeAutoAcceptResult(compose, accept, sessionWorkspace.workspace);
        }
      }
      return formatComposeResult(compose, sessionWorkspace.workspace);
    }

    if (!params.proposal) {
      return formatError(new Error(`gx_compose action "${action}" requires proposal`), { action: `compose.${action}` });
    }

    let dir: string | undefined;
    try {
      dir = await mkdtemp(join(tmpdir(), "gx-compose-plan-"));
      const planFile = join(dir, "plan.json");
      await writeFile(planFile, JSON.stringify(params.proposal), "utf8");
      const args =
        action === "review-plan"
          ? ["compose", "review-plan", "--json", "--plan-file", planFile]
          : ["compose", "apply-plan", "--json", "--plan-file", planFile];
      if (action === "apply" && params.allow_warnings) {
        args.push("--allow-warnings");
      }
      return formatWorkspaceJsonResult(await runGxJson(args, runOptions), sessionWorkspace.workspace, {
        action: `compose.${action}`,
        nextActions:
          action === "apply"
            ? ["Run gx_publish to publish accepted stacks."]
            : ["Call gx_compose action=apply when the reviewed plan is ready."],
      });
    } finally {
      if (dir) {
        await rm(dir, { recursive: true, force: true });
      }
    }
  } catch (error) {
    return formatError(error, { action: `compose.${action}` });
  }
}

type ComposeRevision = {
  id?: unknown;
  intent?: unknown;
  files?: unknown;
  hunk_ids?: unknown;
  target_stack?: unknown;
  base_stack?: unknown;
};

type ComposeHunk = {
  id?: unknown;
  file?: unknown;
  symbol?: unknown;
  symbol_kind?: unknown;
  patch?: unknown;
};

type ComposeProposal = {
  id?: unknown;
  revisions?: unknown;
  hunks?: unknown;
  hunk_links?: unknown;
};

type ComposeResult = {
  state?: unknown;
  proposal?: ComposeProposal;
};

type RevisionOverview = {
  id: string;
  stack: string;
  base_stack: string | null;
  intent: string;
  hunk_ids: string[];
  hunk_count: number;
  summary: string;
  overview: string;
};

type StackOverview = {
  name: string;
  base_stack: string | null;
  revision_count: number;
  revision_ids: string[];
  overview: string;
  revision_overview: string;
  revisions: RevisionOverview[];
};

function formatComposeResult(result: GxJsonRunResult, workspace: SessionWorkspaceInfo): string {
  const overview = composeOverview(result.json);
  return formatJsonResult(result, {
    action: "compose.propose",
    display: overview.proposal_overview,
    nextActions: composeStateIsReady(result.json)
      ? ["Ready proposal was left pending because auto_accept=false. Call gx_accept with the proposal id to accept it."]
      : ["If compose reports repair hints or blocking warnings, call gx_fix or revise the proposal and run review-plan."],
    extra: {
      session_workspace: workspace,
      overview,
    },
  });
}

function formatComposeAutoAcceptResult(
  compose: GxJsonRunResult,
  accept: GxJsonRunResult,
  workspace: SessionWorkspaceInfo,
): string {
  const overview = composeOverview(compose.json);
  return JSON.stringify(
    {
      ok: accept.exitCode === 0,
      action: "compose.auto_accept",
      command: accept.command,
      cwd: accept.cwd,
      session_workspace: workspace,
      exit_code: accept.exitCode,
      display: overview.proposal_overview,
      overview,
      result: {
        compose: compose.json,
        accept: accept.json,
      },
      stderr: [compose.stderr, accept.stderr].filter(Boolean).join("\n") || undefined,
      next_actions: ["Run gx_publish to publish accepted stacks."],
    },
    null,
    2,
  );
}

function formatWorkspaceJsonResult(
  result: GxJsonRunResult,
  workspace: SessionWorkspaceInfo,
  options: Parameters<typeof formatJsonResult>[1] = {},
): string {
  const parsed = JSON.parse(formatJsonResult(result, options)) as Record<string, unknown>;
  parsed.session_workspace = workspace;
  return JSON.stringify(parsed, null, 2);
}

function proposalIDFromCompose(json: unknown) {
  const proposal = (json as ComposeResult).proposal;
  return nonEmptyString(proposal?.id);
}

function composeStateIsReady(json: unknown) {
  const state = nonEmptyString((json as ComposeResult).state)?.toLowerCase();
  return state === "ready" || state === "ready_to_apply";
}

function composeOverview(json: unknown) {
  const result = json as ComposeResult;
  const revisions = Array.isArray(result.proposal?.revisions)
    ? (result.proposal.revisions as ComposeRevision[])
    : [];
  const hunks = Array.isArray(result.proposal?.hunks)
    ? (result.proposal.hunks as ComposeHunk[])
    : [];
  const hunkLinks = Array.isArray(result.proposal?.hunk_links) ? result.proposal.hunk_links : [];
  const hunksById = hunkLookup(hunks);
  const revisionRows = revisionOverviews(revisions, hunksById);
  const stacks = stackOverviews(revisionRows);
  const state = typeof result.state === "string" ? result.state : "unknown";
  const stackTable = formatStackTable(stacks);
  const revisionTable = formatRevisionTable(revisionRows);
  return {
    state,
    revision_count: revisions.length,
    stack_count: stacks.length,
    hunk_link_count: hunkLinks.length,
    proposal_overview: formatProposalOverview(state, stacks.length, revisions.length, stackTable, revisionTable),
    stack_table: stackTable,
    revision_table: revisionTable,
    stacks,
    revisions: revisionRows,
  };
}

function hunkLookup(hunks: ComposeHunk[]): Map<string, ComposeHunk> {
  const out = new Map<string, ComposeHunk>();
  for (const hunk of hunks) {
    const id = nonEmptyString(hunk.id);
    if (id) {
      out.set(id, hunk);
    }
  }
  return out;
}

function revisionOverviews(revisions: ComposeRevision[], hunksById: Map<string, ComposeHunk>): RevisionOverview[] {
  return revisions.map((revision) => {
    const id = nonEmptyString(revision.id) || "(unknown)";
    const stack = nonEmptyString(revision.target_stack) || "(current stack)";
    const baseStack = nonEmptyString(revision.base_stack) || null;
    const intent = nonEmptyString(revision.intent) || "(no intent)";
    const hunkIds = stringArray(revision.hunk_ids);
    const hunks = hunkIds.map((hunkId) => hunksById.get(hunkId)).filter((hunk): hunk is ComposeHunk => Boolean(hunk));
    const summary = summarizeRevisionChange(intent, hunks);
    return {
      id,
      stack,
      base_stack: baseStack,
      intent,
      hunk_ids: hunkIds,
      hunk_count: hunkIds.length,
      summary,
      overview: summarizeRevision(id, summary, hunkIds),
    };
  });
}

function stackOverviews(revisions: RevisionOverview[]): StackOverview[] {
  const byStack = new Map<string, StackOverview>();
  for (const revision of revisions) {
    const name = revision.stack;
    const baseStack = revision.base_stack;
    let stack = byStack.get(name);
    if (!stack) {
      stack = {
        name,
        base_stack: baseStack,
        revision_count: 0,
        revision_ids: [],
        overview: "",
        revision_overview: "",
        revisions: [],
      };
      byStack.set(name, stack);
    } else if (!stack.base_stack && baseStack) {
      stack.base_stack = baseStack;
    }
    stack.revisions.push(revision);
  }

  return [...byStack.values()].map((stack) => {
    const revisionIds = stack.revisions.map((revision) => formatRevisionLabel(revision.id));
    return {
      ...stack,
      revision_count: stack.revisions.length,
      revision_ids: revisionIds,
      overview: summarizeStack(stack.revisions),
      revision_overview: formatRevisionOverview(stack.revisions),
    };
  });
}

function formatProposalOverview(state: string, stackCount: number, revisionCount: number, stackTable: string, revisionTable: string): string {
  return [
    `STATE: ${state}`,
    `STACKS: ${stackCount}`,
    `REVISIONS: ${revisionCount}`,
    "",
    stackTable,
    "",
    revisionTable,
  ].join("\n");
}

function summarizeStack(revisions: RevisionOverview[]): string {
  const summaries = revisions.map((revision) => revision.summary).filter(Boolean);
  if (summaries.length === 0) {
    return "No revision intent provided.";
  }
  if (summaries.length === 1) {
    return summaries[0];
  }
  return summaries.slice(0, 4).join("; ") + (summaries.length > 4 ? `; plus ${summaries.length - 4} more revisions` : "");
}

function summarizeRevision(id: string, summary: string, hunkIds: string[]): string {
  const hunkSummary = hunkIds.length === 1 ? "1 hunk" : `${hunkIds.length} hunks`;
  return `${formatRevisionLabel(id)}: ${summary}${hunkIds.length > 0 ? ` (${hunkSummary})` : ""}`;
}

function formatRevisionOverview(revisions: RevisionOverview[]): string {
  if (revisions.length === 0) {
    return "(no revisions)";
  }
  return revisions.map((revision) => `- ${revision.overview}`).join("\n");
}

function summarizeRevisionChange(intent: string, hunks: ComposeHunk[]): string {
  const files = uniqueStrings(hunks.map((hunk) => nonEmptyString(hunk.file)).filter(Boolean));
  const patches = hunks.map((hunk) => nonEmptyString(hunk.patch)).filter(Boolean);
  const addedLines = meaningfulChangedLines(patches, "+");
  const removedLines = meaningfulChangedLines(patches, "-");
  const symbols = uniqueStrings(hunks.map((hunk) => nonEmptyString(hunk.symbol)).filter(Boolean));

  const semanticSummary = summarizeSemanticChange(intent, files, addedLines, removedLines, symbols);
  if (semanticSummary) {
    return semanticSummary;
  }

  const dependencySummary = summarizeDependencyChange(files, addedLines);
  if (dependencySummary) {
    return dependencySummary;
  }

  const testSummary = summarizeTestChange(files, addedLines);
  if (testSummary) {
    return testSummary;
  }

  if (files.length > 0 && files.every(isDocumentationFile)) {
    const prose = summarizeProseLines(addedLines);
    if (prose) {
      return `Document ${prose}`;
    }
  }

  const declarationSummary = summarizeDeclarations(addedLines);
  if (declarationSummary) {
    return declarationSummary;
  }

  if (symbols.length > 0) {
    const added = summarizeCodeLines(addedLines);
    if (added) {
      return `Update ${formatList(symbols, 3)} to ${added}`;
    }
    return `Update ${formatList(symbols, 3)}`;
  }

  const prose = summarizeProseLines(addedLines);
  if (prose) {
    const removed = summarizeProseLines(removedLines);
    if (removed && removed !== prose) {
      return `Replace "${truncate(removed, 72)}" with "${truncate(prose, 96)}"`;
    }
    return prose;
  }

  return cleanIntent(intent);
}

function summarizeSemanticChange(
  intent: string,
  files: string[],
  addedLines: string[],
  removedLines: string[],
  symbols: string[],
): string {
  const allText = [...addedLines, ...removedLines, intent, ...symbols].join("\n").toLowerCase();

  if (files.every((file) => file.startsWith("mcp/src/tools/"))) {
    const toolSummary = summarizeMcpToolChange(files, allText, symbols);
    if (toolSummary) {
      return toolSummary;
    }
  }

  if (files.includes("mcp/xmcp.config.ts")) {
    return "Update MCP server instructions and exposed tool guidance";
  }

  if (files.includes("mcp/package.json")) {
    return "Update MCP package metadata and scripts";
  }

  if (files.includes("apps/desktop/src/main.js") && /gx_daemon_addr|43124|daemoncontrolurls/.test(allText)) {
    return "Point the desktop app at the stable GX daemon control port";
  }

  if (files.some((file) => file === "internal/cli/service.go" || file.endsWith("internal/cli/service_test.go")) && /controlurl|43124|gxservice/.test(allText)) {
    return "Use the shared stable daemon control URL in CLI service startup";
  }

  if (files.some((file) => file === "internal/launcher/daemon_mgr.go" || file.endsWith("internal/launcher/daemon_mgr_test.go")) && /controlurl|43124|gxservice/.test(allText)) {
    return "Use the shared stable daemon control URL in launcher-managed daemons";
  }

  if (files.includes("internal/daemon/mcp.go") && /resolve.*node|resolve.*mcp|nodepath|mcpserverpath/.test(allText)) {
    return "Resolve bundled Node and MCP server paths before launching provider runtime";
  }

  if (files.some((file) => file.startsWith("internal/vcs/")) && /resolve.*executable|executablefile|execrunner/.test(allText)) {
    return "Resolve executable paths before running VCS commands";
  }

  if (files.some((file) => file.startsWith("internal/authoring/"))) {
    const authoringSummary = summarizeAuthoringChange(files, allText, symbols);
    if (authoringSummary) {
      return authoringSummary;
    }
  }

  if (files.length > 0 && files.every(isDocumentationFile)) {
    const prose = summarizeProseLines(addedLines);
    if (prose && /mcp|demux|compose|gx/.test(allText)) {
      return `Document ${truncate(prose, 110)}`;
    }
  }

  return "";
}

function summarizeMcpToolChange(files: string[], allText: string, symbols: string[]): string {
  if (files.some((file) => file.endsWith("gx-compose-changes.ts"))) {
    if (/proposal_overview|stack_table|revision_table|formatcomposeresult|composeoverview|display/.test(allText)) {
      return "Show the full GX compose proposal with stable stack and revision overviews";
    }
    if (/--json|--plan/.test(allText)) {
      return "Run GX compose through the normal JSON workflow";
    }
    return "Update GX compose MCP behavior";
  }

  if (files.some((file) => file.endsWith("gx-push.ts"))) {
    if (/push|publish with gx|save to cloud|natural requests/.test(allText)) {
      return "Teach MCP publish tool natural-language trigger phrases";
    }
    return "Update GX publish MCP behavior";
  }

  if (files.some((file) => file.endsWith("gx-apply-revision-plan.ts"))) {
    if (/save using gx|save with gx|make commits|make revisions|natural/.test(allText)) {
      return "Teach MCP apply tool natural save and revision trigger phrases";
    }
    return "Update GX apply MCP behavior";
  }

  if (/metadata|description|schema|annotations/.test(allText)) {
    return `Update MCP tool metadata for ${formatList(toolNamesFromFiles(files), 3)}`;
  }

  if (symbols.length > 0) {
    return `Update MCP tool helpers for ${formatList(symbols.map(humanizeIdentifier), 3)}`;
  }

  return "";
}

function summarizeAuthoringChange(files: string[], allText: string, symbols: string[]): string {
  if (files.some((file) => file.endsWith("demux_ai.go")) || /openai|cloud demux|reviewdemuxproposal/.test(allText)) {
    return "Route compose AI review through the GX Cloud OpenAI endpoint";
  }

  if (files.some((file) => file.endsWith("demux_review_shape.go")) || /shape demux|coalescetiny|coalescesemantic|targetminloc/.test(allText)) {
    return "Coalesce tiny compose revisions into semantic review-sized changes";
  }

  if (files.some((file) => file.endsWith("demux_routing.go")) || /feature\/|bug\/|docs\/|chore\/|conventional/.test(allText)) {
    return "Use conventional stack names for compose routing";
  }

  if (symbols.length > 0) {
    return `Update compose authoring around ${formatList(symbols.map(humanizeIdentifier), 3)}`;
  }

  return "";
}

function summarizeDependencyChange(files: string[], addedLines: string[]): string {
  if (!files.some((file) => ["go.mod", "go.sum", "package.json", "package-lock.json", "bun.lock", "bun.lockb", "pnpm-lock.yaml", "yarn.lock"].includes(file))) {
    return "";
  }
  const modules = uniqueStrings(
    addedLines
      .map((line) => line.match(/(?:github\.com|golang\.org|modernc\.org|npm:|@?[\w.-]+\/[\w.-]+)[^\s"',)]*/)?.[0] ?? "")
      .filter(Boolean)
      .map((module) => module.replace(/\/go\.mod$/, "")),
  );
  if (modules.some((module) => module.includes("openai"))) {
    return "Add the OpenAI Go SDK dependency and checksums";
  }
  if (modules.length > 0) {
    return `Add dependency entries for ${formatList(modules, 3)}`;
  }
  return "Update dependency metadata";
}

function summarizeTestChange(files: string[], addedLines: string[]): string {
  if (files.length > 0 && !files.every(isTestFile)) {
    return "";
  }
  const tests = uniqueStrings(
    addedLines
      .map((line) => line.match(/\bfunc\s+(Test[A-Za-z0-9_]+)/)?.[1] ?? line.match(/\b(?:it|test)\(['"`]([^'"`]+)['"`]/)?.[1] ?? "")
      .filter(Boolean),
  );
  if (tests.length === 0) {
    return "";
  }
  return `Add coverage for ${formatList(tests.map(humanizeIdentifier), 3)}`;
}

function summarizeDeclarations(addedLines: string[]): string {
  const functions = uniqueStrings(
    addedLines
      .map((line) =>
        line.match(/\bfunc\s+(?:\([^)]+\)\s*)?([A-Za-z_][A-Za-z0-9_]*)/)?.[1] ??
        line.match(/\bfunction\s+([A-Za-z_][A-Za-z0-9_]*)/)?.[1] ??
        "",
      )
      .filter(Boolean)
      .filter((name) => !name.startsWith("Test")),
  );
  if (functions.length > 0) {
    return `Add ${formatList(functions.map(humanizeIdentifier), 3)}`;
  }

  const declarations = uniqueStrings(
    addedLines
      .map((line) =>
        line.match(/\b(?:function|const|let|var|class|type|interface)\s+([A-Za-z_][A-Za-z0-9_]*)/)?.[1] ??
        "",
      )
      .filter(Boolean),
  );
  if (declarations.length === 0) {
    return "";
  }
  return `Add ${formatList(declarations.map(humanizeIdentifier), 3)}`;
}

function summarizeProseLines(lines: string[]): string {
  const prose = lines
    .map(cleanChangedLine)
    .filter((line) => line.length >= 12)
    .find((line) => /[a-zA-Z]{4}/.test(line));
  return prose ? truncate(sentenceCase(prose), 140) : "";
}

function summarizeCodeLines(lines: string[]): string {
  const code = lines
    .map(cleanChangedLine)
    .filter((line) => line.length >= 6)
    .find((line) => !line.startsWith("//") && !line.startsWith("/*"));
  if (!code) {
    return "";
  }
  const cleaned = truncate(code.replace(/;$/, ""), 90);
  return cleaned.startsWith("return ") ? cleaned : `use ${cleaned}`;
}

function meaningfulChangedLines(patches: string[], marker: "+" | "-"): string[] {
  const skipPrefix = marker === "+" ? "+++" : "---";
  const out: string[] = [];
  for (const patch of patches) {
    for (const line of patch.split("\n")) {
      if (!line.startsWith(marker) || line.startsWith(skipPrefix)) {
        continue;
      }
      const value = line.slice(1).trim();
      if (!value || value === "{" || value === "}" || value === ");" || value === "}," || value === "},") {
        continue;
      }
      out.push(value);
    }
  }
  return out;
}

function cleanChangedLine(line: string): string {
  return line
    .replace(/^[-*]\s+/, "")
    .replace(/[`*#>]/g, "")
    .replace(/\s+/g, " ")
    .trim()
    .replace(/[.;:]$/, "");
}

function isDocumentationFile(file: string): boolean {
  return file.endsWith(".md") || file.endsWith(".mdx") || file.startsWith("docs/") || file.startsWith("skills/");
}

function isTestFile(file: string): boolean {
  return file.endsWith("_test.go") || file.includes(".test.") || file.includes(".spec.") || file.startsWith("test/") || file.includes("/test/") || file.includes("/tests/");
}

function cleanIntent(intent: string): string {
  return nonEmptyString(intent) || "Update revision";
}

function humanizeIdentifier(value: string): string {
  return value
    .replace(/^Test/, "")
    .replace(/([A-Z]+)([A-Z][a-z])/g, "$1 $2")
    .replace(/([a-z0-9])([A-Z])/g, "$1 $2")
    .replace(/[_-]+/g, " ")
    .trim()
    .toLowerCase();
}

function formatList(values: string[], max: number): string {
  const shown = values.slice(0, max);
  return shown.join(", ") + (values.length > max ? `, and ${values.length - max} more` : "");
}

function toolNamesFromFiles(files: string[]): string[] {
  return uniqueStrings(
    files
      .map((file) => file.split("/").pop() ?? "")
      .filter(Boolean)
      .map((file) => file.replace(/\.ts$/, "").replace(/^gx-/, "gx ")),
  );
}

function uniqueStrings(values: string[]): string[] {
  const out: string[] = [];
  const seen = new Set<string>();
  for (const value of values) {
    if (seen.has(value)) {
      continue;
    }
    seen.add(value);
    out.push(value);
  }
  return out;
}

function sentenceCase(value: string): string {
  if (!value) {
    return value;
  }
  return value[0].toUpperCase() + value.slice(1);
}

function formatStackTable(stacks: StackOverview[]): string {
  if (stacks.length === 0) {
    return "STACK  BASE  REVISIONS  OVERVIEW\n(no stacks)";
  }
  const rows = stacks.map((stack) => ({
    stack: stack.name,
    base: stack.base_stack ?? "-",
    revisions: stack.revision_ids.join(", "),
    overview: stack.overview,
  }));
  const stackWidth = columnWidth("STACK", rows.map((row) => row.stack), 32);
  const baseWidth = columnWidth("BASE", rows.map((row) => row.base), 24);
  const revisionsWidth = columnWidth("REVISIONS", rows.map((row) => row.revisions), 24);
  const header = [
    padRight("STACK", stackWidth),
    padRight("BASE", baseWidth),
    padRight("REVISIONS", revisionsWidth),
    "OVERVIEW",
  ].join("  ");
  const divider = [
    "-".repeat(stackWidth),
    "-".repeat(baseWidth),
    "-".repeat(revisionsWidth),
    "-".repeat(8),
  ].join("  ");
  const body = rows.map((row) =>
    [
      padRight(truncate(row.stack, stackWidth), stackWidth),
      padRight(truncate(row.base, baseWidth), baseWidth),
      padRight(truncate(row.revisions, revisionsWidth), revisionsWidth),
      row.overview,
    ].join("  "),
  );
  return [header, divider, ...body].join("\n");
}

function formatRevisionTable(revisions: RevisionOverview[]): string {
  if (revisions.length === 0) {
    return "REVISION  STACK  BASE  HUNKS  OVERVIEW\n(no revisions)";
  }
  const rows = revisions.map((revision) => ({
    revision: formatRevisionLabel(revision.id),
    stack: revision.stack,
    base: revision.base_stack ?? "-",
    hunks: String(revision.hunk_count),
    overview: revision.summary,
  }));
  const revisionWidth = columnWidth("REVISION", rows.map((row) => row.revision), 8);
  const stackWidth = columnWidth("STACK", rows.map((row) => row.stack), 32);
  const baseWidth = columnWidth("BASE", rows.map((row) => row.base), 24);
  const hunksWidth = columnWidth("HUNKS", rows.map((row) => row.hunks), 5);
  const header = [
    padRight("REVISION", revisionWidth),
    padRight("STACK", stackWidth),
    padRight("BASE", baseWidth),
    padRight("HUNKS", hunksWidth),
    "OVERVIEW",
  ].join("  ");
  const divider = [
    "-".repeat(revisionWidth),
    "-".repeat(stackWidth),
    "-".repeat(baseWidth),
    "-".repeat(hunksWidth),
    "-".repeat(8),
  ].join("  ");
  const body = rows.map((row) =>
    [
      padRight(truncate(row.revision, revisionWidth), revisionWidth),
      padRight(truncate(row.stack, stackWidth), stackWidth),
      padRight(truncate(row.base, baseWidth), baseWidth),
      padRight(row.hunks, hunksWidth),
      row.overview,
    ].join("  "),
  );
  return [header, divider, ...body].join("\n");
}

function columnWidth(header: string, values: string[], max: number): number {
  return Math.min(
    Math.max(header.length, ...values.map((value) => value.length)),
    max,
  );
}

function padRight(value: string, width: number): string {
  return value + " ".repeat(Math.max(0, width - value.length));
}

function truncate(value: string, width: number): string {
  if (value.length <= width) {
    return value;
  }
  if (width <= 1) {
    return value.slice(0, width);
  }
  return value.slice(0, width - 1) + "…";
}

function formatRevisionLabel(id: string): string {
  const trimmed = id.trim();
  const match = trimmed.match(/^u([0-9]+)$/);
  return match ? match[1] : trimmed;
}

function nonEmptyString(value: unknown): string {
  return typeof value === "string" ? value.trim() : "";
}

function stringArray(value: unknown): string[] {
  if (!Array.isArray(value)) {
    return [];
  }
  return value.map(nonEmptyString).filter(Boolean);
}
