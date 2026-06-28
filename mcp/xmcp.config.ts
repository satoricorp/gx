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
      "GX MCP exposes gx_sync, gx_generate, gx_status, gx_push, gx_review, and gx_set_base.",
      "Run gx_sync before gx_generate when remote GitHub merges may have landed.",
      "Use gx_generate to create local features and revisions; GX applies safe generated revisions automatically and may append to semantically similar stacks.",
      "Run gx_status after generation to inspect local/remote stack state, then gx_push when features are ready.",
      "If gx_push reports remote divergence, run gx_sync before retrying that stack.",
      "Run gx_review for better codegen context from local facts, previous sessions, PRs, and current code changes.",
      "Use gx_set_base only to return to the repository default branch.",
      "After generate, summarize the output for the user: features, revision intents, notable files, and warnings when present.",
      "Do not use git commit or git push for normal GX save flows unless the user explicitly asks for raw Git.",
    ].join(" "),
  },
};

export default config;
