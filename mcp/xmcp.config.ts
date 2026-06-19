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
    ].join(" "),
  },
};

export default config;
