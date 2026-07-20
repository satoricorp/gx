import { execFile } from "node:child_process";
import { resolve } from "node:path";
import { commandEnvironment, resolveGxBinary } from "./gx";

const gitBinary = () => process.env.GIT_BINARY || "git";

export async function ensureGxInitialized(rawCwd?: string): Promise<{
  repoRoot: string;
  autoInitialized: boolean;
  initOutput?: string;
}> {
  const cwd = resolve(rawCwd || process.cwd());
  const repoRoot = await gitRepoRoot(cwd);
  if (await gxRepoReady(cwd)) {
    return { repoRoot, autoInitialized: false };
  }
  const initOutput = await initializeWithGx(cwd);
  return { repoRoot, autoInitialized: true, initOutput };
}

async function gitRepoRoot(cwd: string): Promise<string> {
  return commandOutput(gitBinary(), ["rev-parse", "--show-toplevel"], cwd);
}

async function gxRepoReady(cwd: string): Promise<boolean> {
  try {
    await commandOutput(resolveGxBinary(), ["status", "--json"], cwd);
    return true;
  } catch {
    return false;
  }
}

async function initializeWithGx(cwd: string): Promise<string> {
  const name = await initName(cwd);
  const email = await initEmail(cwd);
  return commandOutput(resolveGxBinary(), ["init", "--name", name, "--email", email], cwd);
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
