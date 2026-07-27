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
      "GX MCP exposes gx_review.",
      "Save work with plain Git: git add, then git commit. GX installs Git hooks that stamp each commit with its GX revision trailer and record it, so no GX-specific commit verb is needed.",
      "Publish with plain git push, then open the PR with gh pr create. The GX pre-push hook captures the agent session, links edits to the changed hunks, and publishes the GX metadata that becomes the PR summary. Do not run gx push or gx capture push; they bypass or suppress that hook.",
      "Run gx_review for better codegen context from local facts, previous sessions, PRs, and current code changes.",
      "Use git status for inspection.",
    ].join(" "),
  },
};

export default config;
