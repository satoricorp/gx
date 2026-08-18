import { describe, expect, test } from "bun:test";
import { resolveCloudURL, cloudURLOverrideEnv } from "../src/buildconfig";
import { commandEnvironment } from "../src/gx";

// An MCP server inherits its environment from the GUI process hosting it, not
// from a shell anyone chose. These tests pin the consequence: a GX_CLOUD_URL
// that arrives that way is discarded, and pointing the MCP somewhere else takes
// the dedicated variable.
describe("cloud URL resolution", () => {
  test("ignores an inherited GX_CLOUD_URL", () => {
    const resolved = resolveCloudURL({ GX_CLOUD_URL: "http://localhost:3201" });
    expect(resolved).not.toBe("http://localhost:3201");
  });

  test("honours the deliberate override", () => {
    const resolved = resolveCloudURL({
      GX_CLOUD_URL: "http://localhost:3201",
      [cloudURLOverrideEnv]: "http://localhost:9999",
    });
    expect(resolved).toBe("http://localhost:9999");
  });

  test("treats a blank override as absent", () => {
    expect(resolveCloudURL({ [cloudURLOverrideEnv]: "   " })).toBe("");
  });
});

describe("child environment", () => {
  test("does not forward an inherited GX_CLOUD_URL to gx", () => {
    const previous = process.env.GX_CLOUD_URL;
    const previousOverride = process.env[cloudURLOverrideEnv];
    try {
      process.env.GX_CLOUD_URL = "http://localhost:3201";
      delete process.env[cloudURLOverrideEnv];
      const env = commandEnvironment();
      expect(env.GX_CLOUD_URL).not.toBe("http://localhost:3201");
      // Absent, never empty: gx reads this with os.LookupEnv, so a
      // present-but-empty value reads as "cloud disabled" and kills the
      // reviewers exactly as dead as a wrong URL would.
      if (env.GX_CLOUD_URL !== undefined) {
        expect(env.GX_CLOUD_URL.length).toBeGreaterThan(0);
      }
    } finally {
      if (previous === undefined) delete process.env.GX_CLOUD_URL;
      else process.env.GX_CLOUD_URL = previous;
      if (previousOverride !== undefined) process.env[cloudURLOverrideEnv] = previousOverride;
    }
  });

  test("passes the deliberate override through to gx", () => {
    const previousOverride = process.env[cloudURLOverrideEnv];
    try {
      process.env[cloudURLOverrideEnv] = "http://localhost:9999";
      expect(commandEnvironment().GX_CLOUD_URL).toBe("http://localhost:9999");
    } finally {
      if (previousOverride === undefined) delete process.env[cloudURLOverrideEnv];
      else process.env[cloudURLOverrideEnv] = previousOverride;
    }
  });

  test("still marks the client surface", () => {
    const env = commandEnvironment();
    expect(env.GX_CLIENT).toBe("mcp");
    expect(env.GX_MCP).toBe("1");
    expect(env.GX_REVIEW_AI).toBe("1");
  });
});
