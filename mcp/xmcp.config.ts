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
    name: "lgtm MCP",
    instructions: [
      "lgtm MCP exposes lgtm_review.",
      "Save work with plain Git: git add, then git commit. lgtm installs Git hooks that stamp each commit with its lgtm revision trailer and record it, so no lgtm-specific commit verb is needed.",
      "Publish with plain git push, then open the PR with gh pr create. The lgtm pre-push hook captures the agent session, links edits to the changed hunks, and publishes the lgtm metadata that becomes the PR summary. Do not run lgtm push or lgtm capture push; they bypass or suppress that hook.",
      "Run lgtm_review for better codegen context from local facts, previous sessions, PRs, and current code changes.",
      "Use git status for inspection.",
    ].join(" "),
  },
};

export default config;
