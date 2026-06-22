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

export type GxJsonRunResult = GxRunResult & {
  json: unknown;
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

export function parseGxJson(result: GxRunResult): GxJsonRunResult {
  try {
    return { ...result, json: JSON.parse(result.stdout) as unknown };
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error);
    throw Object.assign(new Error(`gx did not return valid JSON: ${message}`), { result });
  }
}

export async function runGxJson(args: string[], options: GxRunOptions = {}): Promise<GxJsonRunResult> {
  return parseGxJson(await runGx(args, options));
}

export function formatJsonResult(result: GxJsonRunResult, options: FormatOptions = {}): string {
  return JSON.stringify(
    envelope(result, {
      ...options,
      result: options.result ?? result.json,
    }),
    null,
    2,
  );
}

export function formatError(error: unknown, options: FormatOptions = {}): string {
  const maybeResult = (error as { result?: GxRunResult }).result;
  if (maybeResult) {
    return formatResult(maybeResult, options);
  }
  const message = error instanceof Error ? error.message : String(error);
  return JSON.stringify(
    {
      ok: false,
      action: options.action,
      command: undefined,
      cwd: options.extra?.cwd ?? undefined,
      exit_code: 1,
      display: message,
      result: null,
      stderr: message,
      next_actions: options.nextActions,
    },
    null,
    2,
  );
}
