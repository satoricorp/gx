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
    name: "Totality MCP",
    instructions: [
      "Totality MCP exposes tl_review.",
      "Save work with plain Git: git add, then git commit. Totality installs Git hooks that stamp each commit with its Totality revision trailer and record it, so no Totality-specific commit verb is needed.",
      "Publish with plain git push, then open the PR with gh pr create. The Totality pre-push hook captures the agent session, links edits to the changed hunks, and publishes the Totality metadata that becomes the PR summary. Do not run tl push or tl capture push; they bypass or suppress that hook.",
      "Run tl_review for better codegen context from local facts, previous sessions, PRs, and current code changes.",
      "Use git status for inspection.",
    ].join(" "),
  },
};

export default config;
