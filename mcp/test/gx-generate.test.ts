import { describe, expect, test } from "bun:test";
import { metadata, schema } from "../src/tools/gx-generate";

describe("gx_generate metadata and schema", () => {
  test("registers the generate tool", () => {
    expect(metadata.name).toBe("gx_generate");
    expect(metadata.description).toMatch(/gx generate/i);
    expect(metadata.description).toMatch(/automatically/i);
    expect(metadata.annotations?.readOnlyHint).toBe(false);
  });

  test("accepts filesets, excludes, and session ids", () => {
    expect(schema.filesets.parse(["README.md"])).toEqual(["README.md"]);
    expect(schema.exclude.parse(["docs"])).toEqual(["docs"]);
    expect(schema.session_id.parse("session-1")).toBe("session-1");
    expect(schema.session_ids.parse(["session-1", "session-2"])).toEqual(["session-1", "session-2"]);
  });
});
