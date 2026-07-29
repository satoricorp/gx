import { execFile } from "node:child_process";
import { existsSync, statSync } from "node:fs";
import { homedir } from "node:os";
import { join, resolve } from "node:path";

export type GxRunOptions = {
  cwd?: string;
  sessionId?: string;
  sessionIds?: string[];
  timeoutMs?: number;
};

export type GxRunResult = {
  command: string[];
  cwd: string;
  exitCode: number;
  stdout: string;
  stderr: string;
};

export type FormatOptions = {
  action?: string;
  display?: string;
  result?: unknown;
  stderr?: string;
  nextActions?: string[];
  extra?: Record<string, unknown>;
};

export function resolveGxBinary() {
  const override = (process.env.GX_BINARY || "").trim();
  if (override) {
    return override;
  }
  const installed = join(homedir(), ".local", "bin", "gx");
  if (existsSync(installed)) {
    return installed;
  }
  return "gx";
}

function resolveCwd(raw?: string) {
  const cwd = resolve(raw || process.cwd());
  let stat;
  try {
    stat = statSync(cwd);
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error);
    throw new Error(`cwd does not exist: ${cwd}: ${message}`);
  }
  if (!stat.isDirectory()) {
    throw new Error(`cwd is not a directory: ${cwd}`);
  }
  return cwd;
}

export function commandEnvironment() {
  const env = { ...process.env };
  env.GX_REVIEW_AI = "1";
  env.GX_MCP = "1";
  const pathEntries = [
    join(homedir(), ".local", "bin"),
    env.PATH || "",
    "/opt/homebrew/bin",
    "/usr/local/bin",
    "/opt/local/bin",
    "/usr/bin",
    "/bin",
  ].flatMap((entry) => entry.split(":").filter(Boolean));
  env.PATH = [...new Set(pathEntries)].join(":");
  return env;
}

export function runGx(args: string[], options: GxRunOptions = {}): Promise<GxRunResult> {
  const cwd = resolveCwd(options.cwd);
  const env = commandEnvironment();
  const binary = resolveGxBinary();
  const sessionIds = [...(options.sessionIds || [])];
  if (options.sessionId) {
    sessionIds.unshift(options.sessionId);
  }
  const cleanedSessionIds = [...new Set(sessionIds.map((id) => id.trim()).filter(Boolean))];
  if (cleanedSessionIds.length === 1) {
    env.GX_SESSION_ID = cleanedSessionIds[0];
  } else if (cleanedSessionIds.length > 1) {
    env.GX_SESSION_IDS = cleanedSessionIds.join(",");
  }

  return new Promise((resolve, reject) => {
    execFile(
      binary,
      args,
      {
        cwd,
        env,
        timeout: options.timeoutMs ?? 120_000,
        maxBuffer: 10 * 1024 * 1024,
      },
      (error, stdout, stderr) => {
        const exitCode =
          typeof (error as NodeJS.ErrnoException | null)?.code === "number"
            ? Number((error as NodeJS.ErrnoException).code)
            : error
              ? 1
              : 0;
        const result = {
          command: [binary, ...args],
          cwd,
          exitCode,
          stdout: redact(stdout.trim(), env),
          stderr: redact(stderr.trim(), env),
        };
        if (error) {
          reject(Object.assign(new Error(result.stderr || result.stdout || error.message), { result }));
          return;
        }
        resolve(result);
      },
    );
  });
}

function redact(value: string, env: NodeJS.ProcessEnv) {
  let out = value;
  const names = [
    "OPENAI_API_KEY",
    "GX_OPENAI_API_KEY",
    "ANTHROPIC_API_KEY",
    "GX_UPLOAD_TOKEN",
    "AWS_SECRET_ACCESS_KEY",
  ];
  for (const name of names) {
    const secret = env[name]?.trim();
    if (secret && secret.length >= 8) {
      out = out.split(secret).join(`[REDACTED:${name}]`);
    }
  }
  return out;
}

function envelope(result: GxRunResult, options: FormatOptions = {}) {
  return {
    ok: result.exitCode === 0,
    action: options.action,
    command: result.command,
    cwd: result.cwd,
    exit_code: result.exitCode,
    display: options.display ?? result.stdout,
    result: options.result ?? (result.stdout || null),
    stderr: options.stderr ?? (result.stderr || undefined),
    next_actions: options.nextActions,
    ...(options.extra ?? {}),
  };
}

export function formatResult(result: GxRunResult, options: FormatOptions = {}): string {
  return JSON.stringify(envelope(result, options), null, 2);
}

export function formatError(error: unknown, options: FormatOptions = {}): string {
  const maybeResult = (error as { result?: GxRunResult }).result;
  if (maybeResult) {
    return formatResult(maybeResult, withAuthGuidance(`${maybeResult.stderr}\n${maybeResult.stdout}`, options));
  }
  const message = error instanceof Error ? error.message : String(error);
  const guided = withAuthGuidance(message, options);
  return JSON.stringify(
    {
      ok: false,
      action: guided.action,
      command: undefined,
      cwd: guided.extra?.cwd ?? undefined,
      exit_code: 1,
      display: guided.display ?? message,
      result: null,
      stderr: guided.stderr ?? message,
      next_actions: guided.nextActions,
      ...(guided.extra ?? {}),
    },
    null,
    2,
  );
}

function withAuthGuidance(message: string, options: FormatOptions): FormatOptions {
  const guidance = authGuidance(message);
  if (!guidance) {
    return options;
  }
  return {
    ...options,
    display: options.display ?? guidance.display,
    nextActions: prependNextAction(options.nextActions, guidance.nextAction),
    extra: {
      ...(options.extra ?? {}),
      auth_required: true,
    },
  };
}

function authGuidance(message: string): { display: string; nextAction: string } | undefined {
  const text = message.toLowerCase();
  if (!text.trim()) {
    return undefined;
  }
  const mentionsLogin = text.includes("gx auth login");
  const missingToken = text.includes("github token is not configured") || text.includes("not logged in");
  if (!mentionsLogin && !missingToken) {
    return undefined;
  }
  if (text.includes("gx auth logout") || text.includes("session expired")) {
    const nextAction = "Run `gx auth logout` then `gx auth login` in a terminal, then retry the MCP tool.";
    return {
      display: `GX cloud authentication needs to be refreshed. ${nextAction}`,
      nextAction,
    };
  }
  const nextAction = "Run `gx auth login` in a terminal, then retry the MCP tool.";
  return {
    display: `GX cloud authentication is required. ${nextAction}`,
    nextAction,
  };
}

function prependNextAction(existing: string[] | undefined, action: string): string[] {
  const actions = [action, ...(existing ?? [])];
  return actions.filter((item, index) => actions.indexOf(item) === index);
}
