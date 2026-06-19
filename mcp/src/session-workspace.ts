import { execFile } from "node:child_process";
import { createHash } from "node:crypto";
import { mkdir, realpath, stat } from "node:fs/promises";
import { homedir } from "node:os";
import { dirname, isAbsolute, join, resolve } from "node:path";

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
      "gx_compose requires a session id so MCP work can run in an isolated session workspace; pass session_id or set GX_SESSION_ID",
    );
  }

  const startCwd = resolve(options.cwd || process.cwd());
  const repoRoot = await commandOutput(jjBinary(), ["root"], startCwd);
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
    },
  };
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

function commandOutput(command: string, args: string[], cwd: string): Promise<string> {
  return new Promise((resolveCommand, reject) => {
    execFile(command, args, { cwd, timeout: 120_000, maxBuffer: 10 * 1024 * 1024 }, (error, stdout, stderr) => {
      if (error) {
        reject(new Error((stderr || stdout || error.message).trim()));
        return;
      }
      resolveCommand(stdout.trim());
    });
  });
}
