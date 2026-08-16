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
    name: "gx MCP",
    instructions: [
      "gx MCP exposes gx_enhance and gx_review.",
      "Save work with plain Git: git add, then git commit. gx installs Git hooks that stamp each commit with its gx revision trailer and record it, so no gx-specific commit verb is needed.",
      "Publish with plain git push, then open the PR with gh pr create. The gx pre-push hook captures the agent session, links edits to the changed hunks, and publishes the gx metadata that becomes the PR summary. Do not run gx push or gx capture push; they bypass or suppress that hook.",
      "Run gx_enhance while coding for issues and tips on how to improve the code you are dealing with — it enhances your changes with context from local facts, previous sessions, PRs, and current code changes.",
      "Run gx_review before shipping to check the change against the pre-ship gates.",
      "Use git status for inspection.",
    ].join(" "),
  },
};

export default config;
