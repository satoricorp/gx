import { execFile } from "node:child_process";
import { commandEnvironment, resolveTlBinary } from "./lgtm";

type VersionInfo = {
  version?: string;
  release_version?: string;
  git_sha?: string;
  revision?: string;
};

type Manifest = {
  version?: string;
  git_sha?: string;
  install_url?: string;
  install_command?: string;
};

type UpdateNotice = {
  current_version?: string;
  current_revision?: string;
  latest_version?: string;
  latest_revision?: string;
  install_command: string;
  message: string;
};

let cachedAt = 0;
let cachedNotice: UpdateNotice | undefined;
let inFlight: Promise<UpdateNotice | undefined> | undefined;

const cacheTTL = 60 * 60 * 1000;
const defaultManifestURL = "https://download.lgtm.cx/cli/manifest.json";
const defaultInstallCommand = "curl -fsSL https://download.lgtm.cx/install.sh | sh";

export async function withUpdateNotice(value: unknown): Promise<unknown> {
  const notice = await updateNotice().catch(() => undefined);
  if (!notice || typeof value !== "string") {
    return value;
  }

  let parsed: unknown;
  try {
    parsed = JSON.parse(value);
  } catch {
    return value;
  }
  if (!parsed || typeof parsed !== "object" || Array.isArray(parsed)) {
    return value;
  }

  const output = parsed as Record<string, unknown>;
  output.update_available = true;
  output.update = notice;
  output.next_actions = prependNextAction(
    Array.isArray(output.next_actions) ? output.next_actions.filter((item): item is string => typeof item === "string") : undefined,
    `Update lgtm with \`${notice.install_command}\`.`,
  );
  return JSON.stringify(output, null, 2);
}

async function updateNotice(): Promise<UpdateNotice | undefined> {
  if (process.env.LGTM_MCP_UPDATE_CHECK === "0" || process.env.LGTM_MOCK_LOG) {
    return undefined;
  }
  const now = Date.now();
  if (now - cachedAt < cacheTTL) {
    return cachedNotice;
  }
  if (!inFlight) {
    inFlight = checkForUpdate()
      .then((notice) => {
        cachedNotice = notice;
        cachedAt = Date.now();
        return notice;
      })
      .catch(() => {
        cachedNotice = undefined;
        cachedAt = Date.now();
        return undefined;
      })
      .finally(() => {
        inFlight = undefined;
      });
  }
  return inFlight;
}

async function checkForUpdate(): Promise<UpdateNotice | undefined> {
  const [current, manifest] = await Promise.all([currentVersion(), latestManifest()]);
  if (!current || !manifest || !isNewer(current, manifest)) {
    return undefined;
  }

  const installCommand = manifest.install_command || (manifest.install_url ? `curl -fsSL ${manifest.install_url} | sh` : defaultInstallCommand);
  const latest = manifest.version || shortSHA(manifest.git_sha) || "latest";
  return {
    current_version: current.release_version || current.version,
    current_revision: current.git_sha || current.revision,
    latest_version: manifest.version,
    latest_revision: manifest.git_sha,
    install_command: installCommand,
    message: `lgtm ${latest} is available. Run \`${installCommand}\`.`,
  };
}

function isNewer(current: VersionInfo, manifest: Manifest): boolean {
  const currentSHA = clean(current.git_sha || current.revision);
  const latestSHA = clean(manifest.git_sha);
  if (currentSHA && latestSHA) {
    return !sameRevision(currentSHA, latestSHA);
  }

  const currentVersion = clean(current.release_version) || clean(current.version);
  const latestVersion = clean(manifest.version);
  if (currentVersion && latestVersion) {
    return currentVersion !== latestVersion;
  }
  return false;
}

function sameRevision(left: string, right: string): boolean {
  return left === right || left.startsWith(right) || right.startsWith(left);
}

async function currentVersion(): Promise<VersionInfo | undefined> {
  return new Promise((resolve) => {
    execFile(
      resolveTlBinary(),
      ["version", "--json"],
      {
        env: commandEnvironment(),
        timeout: 1500,
        maxBuffer: 1024 * 1024,
      },
      (error, stdout) => {
        if (error) {
          resolve(undefined);
          return;
        }
        try {
          resolve(JSON.parse(stdout.trim()) as VersionInfo);
        } catch {
          resolve(undefined);
        }
      },
    );
  });
}

async function latestManifest(): Promise<Manifest | undefined> {
  const url = process.env.LGTM_UPDATE_MANIFEST_URL?.trim() || defaultManifestURL;
  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), 1500);
  try {
    const response = await fetch(url, {
      headers: { accept: "application/json" },
      signal: controller.signal,
    });
    if (!response.ok) {
      return undefined;
    }
    return (await response.json()) as Manifest;
  } catch {
    return undefined;
  } finally {
    clearTimeout(timeout);
  }
}

function prependNextAction(existing: string[] | undefined, action: string): string[] {
  const actions = [action, ...(existing ?? [])];
  return actions.filter((item, index) => actions.indexOf(item) === index);
}

function clean(value: string | undefined): string {
  return (value || "").trim();
}

function shortSHA(value: string | undefined): string {
  const cleaned = clean(value);
  return cleaned.length > 8 ? cleaned.slice(0, 8) : cleaned;
}
