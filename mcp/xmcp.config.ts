import type { XmcpConfig } from "xmcp";

const config: XmcpConfig = {
  stdio: {
    silent: true,
  },
  http: false,
  paths: {
    tools: "src/tools",
    prompts: false,
    resources: false,
  },
  template: {
    name: "GX MCP",
    instructions: [
      "GX MCP exposes gx_commit, gx_status, gx_push, gx_sync, gx_generate, gx_review, and gx_set_base.",
      "When the user says save work, save using gx, or save with gx, run the GX save workflow: git add, gx_commit, gx_status, then gx_push for ready stacks unless the user asks to keep work local.",
      "gx_commit is the default verb. Use it after git add to record staged work as a GX revision.",
      "Run gx_sync before gx_generate or gx_push when remote GitHub merges may have landed.",
      "Use gx_generate only for bulk organization of large working copies; GX applies safe generated revisions automatically and may append to semantically similar stacks.",
      "Run gx_status after commits or generation to inspect local/remote stack state, then gx_push when features are ready.",
      "If gx_push reports remote divergence, run gx_sync before retrying that stack.",
      "Run gx_review for better codegen context from local facts, previous sessions, PRs, and current code changes.",
      "Use gx_set_base only to return to the repository default branch.",
      "After generate, summarize the output for the user: features, revision intents, notable files, and warnings when present.",
      "Use GX MCP tools before the CLI for GX workflows; use the CLI only as fallback when MCP is unavailable.",
      "Do not use git commit or git push for normal GX save flows unless the user explicitly asks for raw Git. Prefer denying or requiring approval for raw git commit, push, reset, and branch deletion in clients that support tool policies.",
    ].join(" "),
  },
};

export default config;
