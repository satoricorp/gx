// Link-time defaults for the MCP server, the Bun analogue of
// internal/buildconfig's -ldflags -X values for the Go CLI.
//
// gx-mcp is compiled by `bun build --compile`, so the Go linker never touches
// it and it cannot carry the CLI's baked endpoints. Without something here it
// had no notion of production at all, and the endpoint it handed the gx child
// was whatever the environment happened to contain.
//
// GX_BAKED_CLOUD_URL is replaced at compile time by `bun build --define`. The
// `typeof` guard is load-bearing rather than defensive: in a build that did not
// define it — `bun run dev`, `bun test`, a checkout with no .env — the
// identifier is undeclared at runtime, and `typeof` is the one operator that
// reads an undeclared name without throwing ReferenceError.
declare const GX_BAKED_CLOUD_URL: string | undefined;

export function bakedCloudURL(): string {
  return typeof GX_BAKED_CLOUD_URL === "string" ? GX_BAKED_CLOUD_URL.trim() : "";
}

/** The env var that points the MCP at a non-production gx server. */
export const cloudURLOverrideEnv = "GX_MCP_CLOUD_URL";

// Which gx server the MCP talks to, and why an inherited GX_CLOUD_URL does not
// get a vote.
//
// An MCP server is not launched from a shell the operator can see. It is spawned
// by whatever GUI process is hosting it, inheriting that process's environment,
// which in turn inherited it from whatever shell opened the app — possibly weeks
// ago, possibly with a dev server URL exported for an afternoon's work. So the
// usual rule that "the environment is the operator's deliberate choice, and beats
// the bake" is simply false here: nobody chose it, and nobody can see it.
//
// Measured consequence, which is what this function exists to prevent: an
// inherited GX_CLOUD_URL=http://localhost:3201 was forwarded to every review the
// MCP ran, so both reviewer legs, the judge, and the rule labeler all failed with
// connection refused against a dev server that was not running. The review still
// returned — degraded, heuristic-only, zero files read — and the reason sat four
// lines down in the run notes.
//
// So production is the default and an inherited value is discarded. Pointing the
// MCP at a local server is still supported, but it takes GX_MCP_CLOUD_URL, a
// variable whose only purpose is that, set deliberately in the MCP server config
// where a reader can see it.
export function resolveCloudURL(env: NodeJS.ProcessEnv = process.env): string {
  const override = (env[cloudURLOverrideEnv] || "").trim();
  if (override) {
    return override;
  }
  return bakedCloudURL();
}
