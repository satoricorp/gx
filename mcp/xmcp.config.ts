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
      "GX MCP exposes gx_commit, gx_status, and gx_review.",
      "When the user says save work, save using gx, or save with gx, run the GX save workflow: git add, gx_commit, gx_status, then publish ready work with plain git push (the GX pre-push hook captures the session and publishes) unless the user asks to keep work local.",
      "gx_commit is the default verb. Use it after git add to record staged work as a GX revision.",
      "Run gx_status after commits to inspect local/remote stack state before publishing.",
      "Run gx_review for better codegen context from local facts, previous sessions, PRs, and current code changes.",
      "Publish with plain git push, then open the PR with gh pr create. Do not run gx push or gx capture push; they bypass or suppress the pre-push hook that GX publishes through.",
      "Use GX MCP tools before the CLI for GX workflows; use the CLI only as fallback when MCP is unavailable.",
      "Do not use git commit for normal GX save flows unless the user explicitly asks for raw Git. Prefer denying or requiring approval for raw git commit, reset, and branch deletion in clients that support tool policies.",
    ].join(" "),
  },
};

export default config;
