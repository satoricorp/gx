import { execFile } from "node:child_process";
import { resolve } from "node:path";
import { commandEnvironment, resolveGxBinary } from "./gx";

const jjBinary = () => process.env.JJ_BINARY || "jj";
const gitBinary = () => process.env.GIT_BINARY || "git";

export async function ensureGxInitialized(rawCwd?: string): Promise<{
  repoRoot: string;
  autoInitialized: boolean;
  initOutput?: string;
}> {
  const cwd = resolve(rawCwd || process.cwd());
  try {
    return {
      repoRoot: await commandOutput(jjBinary(), ["root"], cwd),
      autoInitialized: false,
    };
  } catch (firstError) {
    const initOutput = await initializeWithGx(cwd, firstError);
    return {
      repoRoot: await commandOutput(jjBinary(), ["root"], cwd),
      autoInitialized: true,
      initOutput,
    };
  }
}

async function initializeWithGx(cwd: string, cause: unknown): Promise<string> {
  const name = await initName(cwd);
  const email = await initEmail(cwd);
  try {
    return await commandOutput(resolveGxBinary(), ["init", "--name", name, "--email", email], cwd);
  } catch (error) {
    const causeText = cause instanceof Error ? cause.message : String(cause);
    const initText = error instanceof Error ? error.message : String(error);
    throw new Error(`repo is not initialized for GX and automatic gx init failed: ${initText}; initial jj root error: ${causeText}`);
  }
}

async function initName(cwd: string): Promise<string> {
  return firstNonEmpty(
    process.env.GX_MCP_INIT_NAME,
    process.env.GIT_AUTHOR_NAME,
    await commandOutputOrEmpty(gitBinary(), ["config", "user.name"], cwd),
    "GX MCP",
  );
}

async function initEmail(cwd: string): Promise<string> {
  return firstNonEmpty(
    process.env.GX_MCP_INIT_EMAIL,
    process.env.GIT_AUTHOR_EMAIL,
    await commandOutputOrEmpty(gitBinary(), ["config", "user.email"], cwd),
    "gx-mcp@example.invalid",
  );
}

function firstNonEmpty(...values: Array<string | undefined>): string {
  for (const value of values) {
    const trimmed = value?.trim();
    if (trimmed) {
      return trimmed;
    }
  }
  return "";
}

async function commandOutputOrEmpty(command: string, args: string[], cwd: string): Promise<string> {
  try {
    return await commandOutput(command, args, cwd);
  } catch {
    return "";
  }
}

function commandOutput(command: string, args: string[], cwd: string): Promise<string> {
  return new Promise((resolveCommand, reject) => {
    execFile(command, args, { cwd, env: commandEnvironment(), timeout: 120_000, maxBuffer: 10 * 1024 * 1024 }, (error, stdout, stderr) => {
      if (error) {
        reject(new Error((stderr || stdout || error.message).trim()));
        return;
      }
      resolveCommand(stdout.trim());
    });
  });
}
