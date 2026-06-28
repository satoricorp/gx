import { execFile } from "node:child_process";
import { createHash } from "node:crypto";
import { mkdir, realpath, stat } from "node:fs/promises";
import { homedir } from "node:os";
import { dirname, isAbsolute, join, resolve } from "node:path";
import { commandEnvironment, resolveGxBinary } from "./gx";

export type SessionWorkspaceOptions = {
  cwd?: string;
  sessionId?: string;
  sessionIds?: string[];
};

export type SessionWorkspaceInfo = {
  cwd: string;
  repoRoot: string;
  repoHash: string;
  sessionHash: string;
  workspaceName: string;
  created: boolean;
  autoInitialized: boolean;
  initOutput?: string;
};

const jjBinary = () => process.env.JJ_BINARY || "jj";
const gitBinary = () => process.env.GIT_BINARY || "git";

export function resolvedSessionIds(options: SessionWorkspaceOptions = {}): string[] {
  const raw = [
    ...(options.sessionId ? [options.sessionId] : []),
    ...(options.sessionIds || []),
    ...envSessionIDs(),
  ];
  const out: string[] = [];
  const seen = new Set<string>();
  for (const value of raw.flatMap(splitSessionIDs)) {
    const id = value.trim();
    if (!id || seen.has(id)) {
      continue;
    }
    seen.add(id);
    out.push(id);
  }
  return out;
}

export async function resolveSessionWorkspace(
  options: SessionWorkspaceOptions = {},
): Promise<{ cwd: string; sessionIds: string[]; workspace: SessionWorkspaceInfo }> {
  const sessionIds = resolvedSessionIds(options);
  const primarySessionID = sessionIds[0];
  if (!primarySessionID) {
    throw new Error(
      "gx_generate requires a session id so MCP work can run in an isolated session workspace; pass session_id or set GX_SESSION_ID",
    );
  }

  const startCwd = resolve(options.cwd || process.cwd());
  const initialized = await ensureGxInitialized(startCwd);
  const repoRoot = initialized.repoRoot;
  const repoIdentity = await repoIdentityPath(repoRoot);
  const repoHash = hashID(repoIdentity);
  const sessionHash = hashID(primarySessionID);
  const workspaceName = `gx-${sessionHash}`;
  const cwd = join(sessionWorkspaceRoot(), repoHash, sessionHash);
  const created = await ensureWorkspace(repoRoot, cwd, workspaceName);

  return {
    cwd,
    sessionIds,
    workspace: {
      cwd,
      repoRoot,
      repoHash,
      sessionHash,
      workspaceName,
      created,
      autoInitialized: initialized.autoInitialized,
      initOutput: initialized.initOutput,
    },
  };
}

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

function envSessionIDs(): string[] {
  return [
    process.env.GX_SESSION_ID,
    process.env.GX_SESSION_IDS,
    process.env.CODEX_SESSION_ID,
    process.env.CLAUDE_SESSION_ID,
    process.env.CLAUDE_CODE_SESSION_ID,
    process.env.CURSOR_SESSION_ID,
    process.env.OPENCODE_SESSION_ID,
    process.env.QODO_SESSION_ID,
    process.env.QUAD_CODE_SESSION_ID,
  ].filter((value): value is string => Boolean(value));
}

function splitSessionIDs(value: string): string[] {
  return value.split(/[,\s]+/).filter(Boolean);
}

function hashID(value: string): string {
  return createHash("sha256").update(value).digest("hex").slice(0, 16);
}

function sessionWorkspaceRoot(): string {
  const explicit = process.env.GX_MCP_WORKSPACE_ROOT?.trim();
  if (explicit) {
    return resolve(explicit);
  }
  const gxHome = process.env.GX_HOME?.trim() || join(homedir(), ".gx");
  return join(gxHome, "workspaces");
}

async function repoIdentityPath(repoRoot: string): Promise<string> {
  try {
    const commonDir = await commandOutput(gitBinary(), ["rev-parse", "--git-common-dir"], repoRoot);
    const commonPath = isAbsolute(commonDir) ? commonDir : resolve(repoRoot, commonDir);
    return await realpath(commonPath);
  } catch {
    return await realpath(repoRoot);
  }
}

async function ensureWorkspace(repoRoot: string, workspacePath: string, workspaceName: string): Promise<boolean> {
  if (await isDirectory(workspacePath)) {
    await commandOutput(jjBinary(), ["root"], workspacePath);
    return false;
  }
  await mkdir(dirname(workspacePath), { recursive: true });
  await commandOutput(jjBinary(), ["workspace", "add", "--name", workspaceName, workspacePath], repoRoot);
  return true;
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

async function isDirectory(path: string): Promise<boolean> {
  try {
    return (await stat(path)).isDirectory();
  } catch (error) {
    if ((error as NodeJS.ErrnoException).code === "ENOENT") {
      return false;
    }
    throw error;
  }
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
