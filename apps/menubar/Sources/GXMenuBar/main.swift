import AppKit
import Foundation
import Sparkle

private let pollInterval: TimeInterval = 60
private let releaseCheckInterval: TimeInterval = 60 * 60
private let defaultConsoleURL = "https://gx.run"
private let documentationURL = "https://docs.gx.run"
private let defaultGitHubLatestReleaseURL = "https://api.github.com/repos/satoricorp/gx/releases/latest"

private struct CommandResult {
    let ok: Bool
    let code: Int32
    let stdout: String
    let stderr: String
    let error: String
}

private enum CommandRunner {
    static func run(_ executable: String, _ args: [String], timeout: TimeInterval = 60) -> CommandResult {
        let process = Process()
        if executable.contains("/") {
            process.executableURL = URL(fileURLWithPath: executable)
            process.arguments = args
        } else {
            process.executableURL = URL(fileURLWithPath: "/usr/bin/env")
            process.arguments = [executable] + args
        }
        process.currentDirectoryURL = URL(fileURLWithPath: NSHomeDirectory())
        process.environment = commandEnvironment()

        let stdout = Pipe()
        let stderr = Pipe()
        process.standardOutput = stdout
        process.standardError = stderr

        let semaphore = DispatchSemaphore(value: 0)
        process.terminationHandler = { _ in semaphore.signal() }

        do {
            try process.run()
        } catch {
            return CommandResult(ok: false, code: 1, stdout: "", stderr: "", error: error.localizedDescription)
        }

        if semaphore.wait(timeout: .now() + timeout) == .timedOut {
            process.terminate()
            _ = semaphore.wait(timeout: .now() + 2)
            return CommandResult(ok: false, code: 124, stdout: "", stderr: "", error: "Command timed out")
        }

        let out = String(data: stdout.fileHandleForReading.readDataToEndOfFile(), encoding: .utf8) ?? ""
        let err = String(data: stderr.fileHandleForReading.readDataToEndOfFile(), encoding: .utf8) ?? ""
        let code = process.terminationStatus
        return CommandResult(
            ok: code == 0,
            code: code,
            stdout: out.trimmingCharacters(in: .whitespacesAndNewlines),
            stderr: err.trimmingCharacters(in: .whitespacesAndNewlines),
            error: code == 0 ? "" : (err.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty ? "gx exited \(code)" : err.trimmingCharacters(in: .whitespacesAndNewlines))
        )
    }

    static func commandEnvironment(extra: [String: String] = [:]) -> [String: String] {
        var env = ProcessInfo.processInfo.environment
        let paths = [
            env["PATH"] ?? "",
            "\(NSHomeDirectory())/.local/bin",
            "/opt/homebrew/bin",
            "/usr/local/bin",
            "/opt/local/bin",
            "/usr/bin",
            "/bin"
        ]
            .flatMap { $0.split(separator: ":").map(String.init) }
        var seen = Set<String>()
        let uniquePaths = paths.filter { !$0.isEmpty && seen.insert($0).inserted }
        env["PATH"] = uniquePaths.joined(separator: ":")
        env["GX_MENUBAR"] = "1"
        for (key, value) in extra {
            env[key] = value
        }
        return env
    }

    static func which(_ executable: String) -> String? {
        let result = run("/usr/bin/env", ["which", executable], timeout: 5)
        guard result.ok else { return nil }
        return result.stdout.split(separator: "\n").first.map(String.init)
    }
}

private enum GXConfig {
    static func cleanURL(_ value: String) -> String {
        value.trimmingCharacters(in: .whitespacesAndNewlines).replacingOccurrences(of: #"/+$"#, with: "", options: .regularExpression)
    }
}

private enum CLIInstaller {
    static var installPath: URL {
        URL(fileURLWithPath: NSHomeDirectory(), isDirectory: true)
            .appendingPathComponent(".local", isDirectory: true)
            .appendingPathComponent("bin", isDirectory: true)
            .appendingPathComponent("gx")
    }

    private static var installDirectory: URL {
        installPath.deletingLastPathComponent()
    }

    static func bundledGXURL() -> URL? {
        guard let resourceURL = Bundle.main.resourceURL else { return nil }
        let candidate = resourceURL.appendingPathComponent("bin", isDirectory: true).appendingPathComponent("gx")
        return FileManager.default.isExecutableFile(atPath: candidate.path) ? candidate : nil
    }

    static func gxExecutable() -> String {
        let env = ProcessInfo.processInfo.environment
        let override = env["GX_BINARY"]?.trimmingCharacters(in: .whitespacesAndNewlines) ?? ""
        if !override.isEmpty {
            return override
        }
        if FileManager.default.isExecutableFile(atPath: installPath.path) {
            return installPath.path
        }
        if let bundled = bundledGXURL() {
            return bundled.path
        }
        return CommandRunner.which("gx") ?? "gx"
    }

    static func installBundledCLI() -> String {
        guard let source = bundledGXURL() else {
            if FileManager.default.isExecutableFile(atPath: installPath.path) {
                do {
                    try installShortcutSymlinks()
                } catch {
                    return "CLI shortcuts failed: \(error.localizedDescription)"
                }
                return "CLI installed at \(installPath.path)"
            }
            return "Bundled gx CLI not found"
        }

        do {
            try FileManager.default.createDirectory(at: installDirectory, withIntermediateDirectories: true)
            if source.standardizedFileURL.path != installPath.standardizedFileURL.path {
                if FileManager.default.fileExists(atPath: installPath.path) {
                    try FileManager.default.removeItem(at: installPath)
                }
                try FileManager.default.copyItem(at: source, to: installPath)
            }
            try FileManager.default.setAttributes([.posixPermissions: 0o755], ofItemAtPath: installPath.path)
            try installShortcutSymlinks()
            return "CLI installed at \(installPath.path) with gxr and gxs shortcuts"
        } catch {
            return "CLI install failed: \(error.localizedDescription)"
        }
    }

    static func installExtractedCLI(from gxRoot: URL) throws {
        let binDirectory = gxRoot.appendingPathComponent("bin", isDirectory: true)
        let gxSource = binDirectory.appendingPathComponent("gx")
        let mcpSource = binDirectory.appendingPathComponent("gx-mcp")
        guard FileManager.default.isExecutableFile(atPath: gxSource.path) else {
            throw NSError(domain: "GXMenuBar", code: 1, userInfo: [NSLocalizedDescriptionKey: "Archive is missing gx"])
        }

        try FileManager.default.createDirectory(at: installDirectory, withIntermediateDirectories: true)
        try installExecutable(gxSource, to: installPath)
        if FileManager.default.isExecutableFile(atPath: mcpSource.path) {
            try installExecutable(mcpSource, to: installDirectory.appendingPathComponent("gx-mcp"))
        }
        try installShortcutSymlinks()
    }

    private static func installExecutable(_ source: URL, to destination: URL) throws {
        let temp = installDirectory.appendingPathComponent(".\(destination.lastPathComponent).\(UUID().uuidString).tmp")
        try FileManager.default.copyItem(at: source, to: temp)
        try FileManager.default.setAttributes([.posixPermissions: 0o755], ofItemAtPath: temp.path)
        if FileManager.default.fileExists(atPath: destination.path) {
            _ = try FileManager.default.replaceItemAt(destination, withItemAt: temp, backupItemName: nil, options: .usingNewMetadataOnly)
        } else {
            try FileManager.default.moveItem(at: temp, to: destination)
        }
    }

    private static func installShortcutSymlinks() throws {
        for name in ["gxr", "gxs"] {
            let shortcut = installDirectory.appendingPathComponent(name)
            try? FileManager.default.removeItem(at: shortcut)
            try FileManager.default.createSymbolicLink(atPath: shortcut.path, withDestinationPath: "gx")
        }
    }
}

private struct CLIVersionInfo: Decodable {
    let version: String?
    let releaseVersion: String?
    let gitSHA: String?
    let revision: String?

    enum CodingKeys: String, CodingKey {
        case version
        case releaseVersion = "release_version"
        case gitSHA = "git_sha"
        case revision
    }

    var display: String {
        clean(releaseVersion) ?? clean(version) ?? clean(revision) ?? "unknown"
    }

    var comparisonVersion: String? {
        if let release = clean(releaseVersion), release != "dev" {
            return release
        }
        if let value = clean(version), value != "dev", SemVer.parse(value) != nil {
            return value
        }
        return nil
    }

    var comparisonRevision: String? {
        clean(gitSHA) ?? clean(revision)
    }
}

private struct GitHubRelease: Decodable {
    let tagName: String
    let name: String?
    let htmlURL: String?
    let assets: [GitHubAsset]

    enum CodingKeys: String, CodingKey {
        case tagName = "tag_name"
        case name
        case htmlURL = "html_url"
        case assets
    }

    var displayVersion: String {
        tagName.trimmingCharacters(in: .whitespacesAndNewlines)
    }

    var normalizedVersion: String {
        normalizeVersion(displayVersion)
    }
}

private struct GitHubAsset: Decodable {
    let name: String
    let browserDownloadURL: String

    enum CodingKeys: String, CodingKey {
        case name
        case browserDownloadURL = "browser_download_url"
    }
}

private struct SemVer: Comparable {
    let major: Int
    let minor: Int
    let patch: Int
    let prerelease: String?

    static func parse(_ value: String) -> SemVer? {
        let normalized = normalizeVersion(value)
        let withoutBuild = normalized.split(separator: "+", maxSplits: 1, omittingEmptySubsequences: false).first.map(String.init) ?? normalized
        let versionParts = withoutBuild.split(separator: "-", maxSplits: 1, omittingEmptySubsequences: false)
        let core = versionParts.first.map(String.init) ?? ""
        let prerelease = versionParts.count > 1 ? String(versionParts[1]) : nil
        let numbers = core.split(separator: ".", omittingEmptySubsequences: false)
        guard numbers.count >= 2, numbers.count <= 3 else { return nil }
        guard let major = Int(numbers[0]), let minor = Int(numbers[1]) else { return nil }
        let patch = numbers.count == 3 ? Int(numbers[2]) : 0
        guard let patch else { return nil }
        return SemVer(major: major, minor: minor, patch: patch, prerelease: prerelease?.isEmpty == true ? nil : prerelease)
    }

    static func < (left: SemVer, right: SemVer) -> Bool {
        if left.major != right.major { return left.major < right.major }
        if left.minor != right.minor { return left.minor < right.minor }
        if left.patch != right.patch { return left.patch < right.patch }
        switch (left.prerelease, right.prerelease) {
        case (nil, nil):
            return false
        case (nil, _):
            return false
        case (_, nil):
            return true
        case let (left?, right?):
            return left.localizedStandardCompare(right) == .orderedAscending
        }
    }
}

private struct ReleaseLoad {
    let release: GitHubRelease?
    let error: String?
}

private enum ReleaseClient {
    static func currentCLI() -> (CLIVersionInfo?, String?) {
        let result = CommandRunner.run(CLIInstaller.gxExecutable(), ["version", "--json"], timeout: 5)
        guard result.ok else {
            return (nil, result.error.isEmpty ? "Unable to read CLI version" : result.error)
        }
        do {
            return (try JSONDecoder().decode(CLIVersionInfo.self, from: Data(result.stdout.utf8)), nil)
        } catch {
            return (nil, "Invalid version JSON: \(error.localizedDescription)")
        }
    }

    static func latestRelease() -> ReleaseLoad {
        let urlString = ProcessInfo.processInfo.environment["GX_MENUBAR_RELEASE_URL"]?.trimmingCharacters(in: .whitespacesAndNewlines)
        guard let url = URL(string: (urlString?.isEmpty == false ? urlString! : defaultGitHubLatestReleaseURL)) else {
            return ReleaseLoad(release: nil, error: "Invalid GitHub release URL")
        }
        var request = URLRequest(url: url, timeoutInterval: 10)
        request.setValue("application/vnd.github+json", forHTTPHeaderField: "Accept")
        request.setValue("GXMenuBar", forHTTPHeaderField: "User-Agent")
        let result = fetch(request)
        guard let data = result.data else {
            return ReleaseLoad(release: nil, error: result.error ?? "GitHub release fetch failed")
        }
        do {
            return ReleaseLoad(release: try JSONDecoder().decode(GitHubRelease.self, from: data), error: nil)
        } catch {
            return ReleaseLoad(release: nil, error: "Invalid GitHub release JSON: \(error.localizedDescription)")
        }
    }

    static func isNewer(current: CLIVersionInfo?, latest: GitHubRelease?) -> Bool {
        guard let latest else { return false }
        guard let current else { return true }
        let latestVersion = latest.normalizedVersion
        if let currentVersion = current.comparisonVersion {
            if let currentSemVer = SemVer.parse(currentVersion), let latestSemVer = SemVer.parse(latestVersion) {
                return latestSemVer > currentSemVer
            }
            return normalizeVersion(currentVersion) != latestVersion
        }
        if let currentRevision = current.comparisonRevision, !currentRevision.isEmpty {
            return !sameRevision(currentRevision, latestVersion)
        }
        return true
    }

    static func update(to release: GitHubRelease, status: (String) -> Void) -> String {
        guard let archive = archiveAsset(in: release) else {
            return "Update failed: release asset not found"
        }

        let fileManager = FileManager.default
        let workRoot = fileManager.temporaryDirectory.appendingPathComponent("gx-update-\(UUID().uuidString)", isDirectory: true)
        do {
            try fileManager.createDirectory(at: workRoot, withIntermediateDirectories: true)
            defer { try? fileManager.removeItem(at: workRoot) }

            status("Downloading GX \(release.displayVersion)...")
            let archiveURL = workRoot.appendingPathComponent(archive.name)
            try download(asset: archive, to: archiveURL)

            if let checksum = checksumAsset(for: archive, in: release) {
                status("Verifying GX \(release.displayVersion)...")
                let checksumURL = workRoot.appendingPathComponent(checksum.name)
                try download(asset: checksum, to: checksumURL)
                try verifyChecksum(archiveURL: archiveURL, checksumURL: checksumURL)
            }

            status("Installing GX \(release.displayVersion)...")
            let extractURL = workRoot.appendingPathComponent("extract", isDirectory: true)
            try fileManager.createDirectory(at: extractURL, withIntermediateDirectories: true)
            let tar = CommandRunner.run("/usr/bin/tar", ["-xzf", archiveURL.path, "-C", extractURL.path], timeout: 30)
            guard tar.ok else {
                return "Update failed: \(tar.error.isEmpty ? "archive extraction failed" : tar.error)"
            }

            let gxRoot = extractURL.appendingPathComponent("gx", isDirectory: true)
            let downloadedGX = gxRoot.appendingPathComponent("bin", isDirectory: true).appendingPathComponent("gx")
            let validation = CommandRunner.run(downloadedGX.path, ["version", "--json"], timeout: 5)
            guard validation.ok else {
                return "Update failed: downloaded gx did not run"
            }
            if let downloaded = try? JSONDecoder().decode(CLIVersionInfo.self, from: Data(validation.stdout.utf8)),
               let downloadedVersion = downloaded.comparisonVersion,
               normalizeVersion(downloadedVersion) != release.normalizedVersion {
                return "Update failed: downloaded \(downloaded.display) did not match \(release.displayVersion)"
            }

            try CLIInstaller.installExtractedCLI(from: gxRoot)
            return "Updated CLI to \(release.displayVersion)"
        } catch {
            return "Update failed: \(error.localizedDescription)"
        }
    }

    private static func archiveAsset(in release: GitHubRelease) -> GitHubAsset? {
        let arch = currentArch()
        let tag = release.displayVersion
        let normalized = release.normalizedVersion
        let candidates = [
            "gx_\(tag)_darwin_\(arch).tar.gz",
            "gx_\(normalized)_darwin_\(arch).tar.gz",
            "gx_darwin_\(arch).tar.gz"
        ]
        for candidate in candidates {
            if let asset = release.assets.first(where: { $0.name == candidate }) {
                return asset
            }
        }
        return release.assets.first { asset in
            asset.name.hasSuffix("_darwin_\(arch).tar.gz") || asset.name == "gx_darwin_\(arch).tar.gz"
        }
    }

    private static func checksumAsset(for archive: GitHubAsset, in release: GitHubRelease) -> GitHubAsset? {
        release.assets.first { $0.name == "\(archive.name).sha256" }
    }

    private static func currentArch() -> String {
        #if arch(arm64)
        return "arm64"
        #else
        return "amd64"
        #endif
    }

    private static func download(asset: GitHubAsset, to destination: URL) throws {
        guard let url = URL(string: asset.browserDownloadURL) else {
            throw NSError(domain: "GXMenuBar", code: 2, userInfo: [NSLocalizedDescriptionKey: "Invalid asset URL for \(asset.name)"])
        }
        var request = URLRequest(url: url, timeoutInterval: 60)
        request.setValue("application/octet-stream", forHTTPHeaderField: "Accept")
        request.setValue("GXMenuBar", forHTTPHeaderField: "User-Agent")
        let result = fetch(request)
        guard let data = result.data else {
            throw NSError(domain: "GXMenuBar", code: 3, userInfo: [NSLocalizedDescriptionKey: result.error ?? "Download failed for \(asset.name)"])
        }
        try data.write(to: destination, options: .atomic)
    }

    private static func verifyChecksum(archiveURL: URL, checksumURL: URL) throws {
        let expectedText = try String(contentsOf: checksumURL, encoding: .utf8)
        guard let expected = expectedText.split(whereSeparator: { $0 == " " || $0 == "\n" || $0 == "\t" }).first.map(String.init), !expected.isEmpty else {
            throw NSError(domain: "GXMenuBar", code: 4, userInfo: [NSLocalizedDescriptionKey: "Checksum file was empty"])
        }
        let result = CommandRunner.run("/usr/bin/shasum", ["-a", "256", archiveURL.path], timeout: 10)
        guard result.ok, let actual = result.stdout.split(separator: " ").first.map(String.init) else {
            throw NSError(domain: "GXMenuBar", code: 5, userInfo: [NSLocalizedDescriptionKey: "Could not compute checksum"])
        }
        if expected.lowercased() != actual.lowercased() {
            throw NSError(domain: "GXMenuBar", code: 6, userInfo: [NSLocalizedDescriptionKey: "Checksum mismatch"])
        }
    }

    private static func fetch(_ request: URLRequest) -> (data: Data?, error: String?) {
        let semaphore = DispatchSemaphore(value: 0)
        var receivedData: Data?
        var receivedResponse: URLResponse?
        var receivedError: Error?
        URLSession.shared.dataTask(with: request) { data, response, error in
            receivedData = data
            receivedResponse = response
            receivedError = error
            semaphore.signal()
        }.resume()
        if semaphore.wait(timeout: .now() + request.timeoutInterval + 2) == .timedOut {
            return (nil, "Request timed out")
        }
        if let receivedError {
            return (nil, receivedError.localizedDescription)
        }
        if let http = receivedResponse as? HTTPURLResponse, !(200..<300).contains(http.statusCode) {
            return (nil, "GitHub returned HTTP \(http.statusCode)")
        }
        return (receivedData, nil)
    }

    private static func sameRevision(_ left: String, _ right: String) -> Bool {
        let cleanLeft = left.trimmingCharacters(in: .whitespacesAndNewlines)
        let cleanRight = right.trimmingCharacters(in: .whitespacesAndNewlines)
        return cleanLeft == cleanRight || cleanLeft.hasPrefix(cleanRight) || cleanRight.hasPrefix(cleanLeft)
    }
}

private func clean(_ value: String?) -> String? {
    let trimmed = value?.trimmingCharacters(in: .whitespacesAndNewlines) ?? ""
    return trimmed.isEmpty ? nil : trimmed
}

private func normalizeVersion(_ value: String) -> String {
    var normalized = value.trimmingCharacters(in: .whitespacesAndNewlines)
    while normalized.hasPrefix("v") || normalized.hasPrefix("V") {
        normalized.removeFirst()
    }
    return normalized
}

private struct DoctorEnvelope: Decodable {
    let doctor: DoctorStatus?
}

private struct DoctorStatus: Decodable {
    let ok: Bool?
    let capture: CaptureStatus?
    let ledger: [LedgerRow]?
    let stats: StatsStatus?
    let diagnose: DiagnoseStatus?
}

private struct CaptureStatus: Decodable {
    let hookInstalled: Bool?
    let hookApplicable: Bool?
    let repoRoot: String?
    let repoHooks: [RepoHookStatus]?
    let repoHooksTotal: Int?
    let repoHooksMissing: Int?
    let repoHooksUnreachable: Int?
    let repoHooksOK: Bool?
    let uploadAuthed: Bool?
    let uploadAPI: String?
    let uploadAuthError: String?
    let pendingExtracts: Int?
    let pendingSessions: Int?
    let cursorReachable: Bool?
    let cursorPath: String?
    let diskFreeGB: Int?
    let diskWarn: Bool?
    let ok: Bool?
}

private struct RepoHookStatus: Decodable {
    let repoRoot: String?
    let hookInstalled: Bool?
    let gitReachable: Bool?
    let error: String?
}

private struct LedgerRow: Decodable {
    let agent: String?
    let filepath: String?
    let status: String?
    let calls: String?
    let tokens: String?
    let files: String?
    let last: String?
}

private struct DiagnoseStatus: Decodable {
    let summary: String?
    let lastCheck: String?
}

private struct StatsStatus: Decodable {
    let approvedStacksWaitingForPublish: Int?
    let publishedStacksWaitingForReview: Int?
    let agents: [AgentStats]?
    let diskUsedBytes: Int64?
}

private struct AgentStats: Decodable {
    let agent: String?
    let label: String?
    let health: String?
    let sessions: Int?
}

private struct DoctorLoad {
    let envelope: DoctorEnvelope?
    let rawJSON: String
    let error: String?
}

private enum HealthState: String {
    case green
    case yellow
    case red
    case unknown
}

private enum DoctorClient {
    static func fetch() -> DoctorLoad {
        let result = CommandRunner.run(CLIInstaller.gxExecutable(), ["doctor", "--json"])
        guard result.ok else {
            return DoctorLoad(envelope: nil, rawJSON: result.stdout, error: result.error.isEmpty ? result.stderr : result.error)
        }
        do {
            let data = Data(result.stdout.utf8)
            let envelope = try JSONDecoder().decode(DoctorEnvelope.self, from: data)
            return DoctorLoad(envelope: envelope, rawJSON: result.stdout, error: nil)
        } catch {
            return DoctorLoad(envelope: nil, rawJSON: result.stdout, error: "Invalid doctor JSON: \(error.localizedDescription)")
        }
    }

    static func state(for doctor: DoctorStatus?) -> HealthState {
        guard let doctor else { return .unknown }
        guard let capture = doctor.capture else {
            if doctor.ok == true { return .green }
            if doctor.ok == false { return .red }
            return .unknown
        }
        let hookMissing = capture.hookApplicable != false && capture.hookInstalled == false
        let registeredHookIssue = (capture.repoHooksMissing ?? 0) > 0 || (capture.repoHooksUnreachable ?? 0) > 0
        let hardFail = hookMissing ||
            registeredHookIssue ||
            capture.uploadAuthed == false ||
            capture.cursorReachable == false ||
            capture.diskWarn == true
        if hardFail {
            return .red
        }
        if doctor.ok == false {
            return .yellow
        }
        return .green
    }
}

private final class GXMenuBarApp: NSObject, NSApplicationDelegate {
    private let statusItem = NSStatusBar.system.statusItem(withLength: NSStatusItem.variableLength)
    private let worker = DispatchQueue(label: "dev.gx.menubar.worker", qos: .utility)
    private let updaterController = SPUStandardUpdaterController(
        startingUpdater: true,
        updaterDelegate: nil,
        userDriverDelegate: nil
    )
    private var timer: Timer?
    private var refreshing = false
    private var updatingCLI = false

    private var doctor: DoctorStatus?
    private var doctorRawJSON = ""
    private var doctorError: String?
    private var lastUpdated: Date?
    private var cliInfo: CLIVersionInfo?
    private var cliError: String?
    private var latestRelease: GitHubRelease?
    private var releaseError: String?
    private var lastReleaseChecked: Date?
    private var updateMessage: String?

    func applicationDidFinishLaunching(_ notification: Notification) {
        NSApp.setActivationPolicy(.accessory)
        configureStatusItem()
        rebuildMenu()
        refreshAll(installCLI: true, forceReleaseCheck: true)
        timer = Timer.scheduledTimer(withTimeInterval: pollInterval, repeats: true) { [weak self] _ in
            self?.refreshAll(installCLI: false)
        }
    }

    func applicationWillTerminate(_ notification: Notification) {
        timer?.invalidate()
    }

    private func configureStatusItem() {
        if let button = statusItem.button {
            if let image = Bundle.main.image(forResource: "trayTemplate") {
                image.isTemplate = true
                image.size = NSSize(width: 19, height: 11)
                button.image = image
                button.imagePosition = .imageOnly
            } else {
                button.title = "GX"
            }
        }
    }

    private func refreshAll(installCLI: Bool, forceReleaseCheck: Bool = false) {
        guard !refreshing, !updatingCLI else { return }
        refreshing = true
        rebuildMenu()
        let shouldCheckRelease = forceReleaseCheck || lastReleaseChecked == nil || Date().timeIntervalSince(lastReleaseChecked!) >= releaseCheckInterval
        worker.async { [weak self] in
            let installMessage: String?
            if installCLI {
                installMessage = CLIInstaller.installBundledCLI()
            } else {
                installMessage = nil
            }
            let currentCLI = ReleaseClient.currentCLI()
            let releaseLoad = shouldCheckRelease ? ReleaseClient.latestRelease() : nil
            let doctorLoad = DoctorClient.fetch()

            DispatchQueue.main.async {
                guard let self else { return }
                self.cliInfo = currentCLI.0
                self.cliError = currentCLI.1
                if let releaseLoad {
                    self.latestRelease = releaseLoad.release
                    self.releaseError = releaseLoad.error
                    self.lastReleaseChecked = Date()
                }
                if let installMessage {
                    self.updateMessage = installMessage
                } else if forceReleaseCheck, let releaseLoad {
                    if let error = releaseLoad.error {
                        self.updateMessage = "Update check failed: \(error)"
                    } else if ReleaseClient.isNewer(current: currentCLI.0, latest: releaseLoad.release) {
                        self.updateMessage = "GX \(releaseLoad.release?.displayVersion ?? "latest") is available"
                    } else {
                        self.updateMessage = "CLI is up to date"
                    }
                }
                self.doctor = doctorLoad.envelope?.doctor
                self.doctorRawJSON = doctorLoad.rawJSON
                self.doctorError = doctorLoad.error
                self.lastUpdated = Date()
                self.refreshing = false
                self.rebuildMenu()
            }
        }
    }

    private func rebuildMenu() {
        let menu = NSMenu()
        menu.addItem(disabled(versionSummary()))
        menu.addItem(disabled(releaseSummary()))
        if let updateMessage, !updateMessage.isEmpty {
            menu.addItem(disabled(updateMessage))
        }
        if let updated = lastUpdated {
            menu.addItem(disabled("Updated \(Self.timeFormatter.string(from: updated))"))
        } else if refreshing {
            menu.addItem(disabled("Updating..."))
        }
        menu.addItem(.separator())

        menu.addItem(doctorSubmenuItem())
        menu.addItem(disabled("Waiting for Publishing: \(doctor?.stats?.approvedStacksWaitingForPublish ?? 0)"))
        menu.addItem(disabled("Waiting for Review: \(doctor?.stats?.publishedStacksWaitingForReview ?? 0)"))
        menu.addItem(submenuItem(title: "Stats", submenu: statsMenu()))
        menu.addItem(.separator())

        if updatingCLI {
            menu.addItem(disabled("Updating CLI..."))
        } else if hasRemoteUpdate() {
            menu.addItem(actionItem("Update CLI to \(latestRelease?.displayVersion ?? "latest")", #selector(updateCLIFromGitHub)))
        } else if latestRelease != nil {
            menu.addItem(disabled("CLI is up to date"))
        }
        menu.addItem(actionItem("Check for Updates", #selector(checkForUpdates)))
        menu.addItem(actionItem("Install Bundled CLI", #selector(installCLI)))
        menu.addItem(actionItem("Check for App Updates", #selector(checkForAppUpdates)))
        menu.addItem(.separator())

        menu.addItem(actionItem("Open https://gx.run", #selector(openConsole)))
        menu.addItem(.separator())

        menu.addItem(submenuItem(title: "MCP", submenu: mcpMenu()))
        menu.addItem(.separator())
        menu.addItem(actionItem("Quit", #selector(quit), keyEquivalent: "q"))
        statusItem.menu = menu
    }

    private func doctorMenu() -> NSMenu {
        let menu = NSMenu()
        if let error = doctorError, doctor == nil {
            menu.addItem(disabled(error))
        }
        if let capture = doctor?.capture {
            menu.addItem(healthMenuItem(label: "Hooks", value: hooksLabel(capture), state: hooksState(capture)))
            menu.addItem(healthMenuItem(label: "Auth", value: authLabel(capture), state: authState(capture)))
            menu.addItem(disabled("Upload: \(uploadEnvironment(capture.uploadAPI))"))
        } else {
            menu.addItem(healthMenuItem(label: "Hooks", value: "unknown", state: .unknown))
            menu.addItem(healthMenuItem(label: "Auth", value: "unknown", state: .unknown))
            menu.addItem(disabled("Upload: unknown"))
        }
        menu.addItem(.separator())
        menu.addItem(actionItem("Run Doctor", #selector(runDoctorNow)))
        menu.addItem(actionItem("Copy Diagnostics", #selector(copyDoctorJSON)))
        return menu
    }

    private func statsMenu() -> NSMenu {
        let menu = NSMenu()
        if let stats = doctor?.stats {
            for agent in statsAgentRows(stats.agents) {
                menu.addItem(agentStatsItem(agent))
            }
            menu.addItem(.separator())
            menu.addItem(disabled("Disk Used: \(formatDiskUsed(stats.diskUsedBytes ?? 0))"))
        } else {
            menu.addItem(disabled("Stats unavailable"))
        }
        return menu
    }

    private func mcpMenu() -> NSMenu {
        let menu = NSMenu()
        menu.addItem(actionItem("Show Instructions", #selector(openMCPDocumentation)))
        return menu
    }

    private func healthMenuItem(label: String, value: String, state: HealthState) -> NSMenuItem {
        let item = NSMenuItem(title: "", action: nil, keyEquivalent: "")
        item.isEnabled = false
        let title = NSMutableAttributedString(
            string: "●  ",
            attributes: [
                .font: NSFont.menuFont(ofSize: 0),
                .foregroundColor: healthColor(state)
            ]
        )
        title.append(NSAttributedString(
            string: "\(label): \(value)",
            attributes: [.font: NSFont.menuFont(ofSize: 0)]
        ))
        item.attributedTitle = title
        return item
    }

    private func hooksState(_ capture: CaptureStatus) -> HealthState {
        let hookMissing = capture.hookApplicable != false && capture.hookInstalled == false
        let registeredHookIssue = (capture.repoHooksMissing ?? 0) > 0 || (capture.repoHooksUnreachable ?? 0) > 0 || capture.repoHooksOK == false
        if hookMissing || registeredHookIssue {
            return .red
        }
        let repoHookTotal = capture.repoHooksTotal ?? capture.repoHooks?.count ?? 0
        if capture.hookApplicable == false && repoHookTotal == 0 {
            return .yellow
        }
        return .green
    }

    private func hooksLabel(_ capture: CaptureStatus) -> String {
        let missing = capture.repoHooksMissing ?? 0
        let unreachable = capture.repoHooksUnreachable ?? 0
        if missing > 0 || unreachable > 0 {
            var parts: [String] = []
            if missing > 0 {
                parts.append("\(missing) missing")
            }
            if unreachable > 0 {
                parts.append("\(unreachable) unreachable")
            }
            return parts.joined(separator: ", ")
        }
        if capture.hookApplicable != false && capture.hookInstalled == false {
            return "missing"
        }
        let repoHookTotal = capture.repoHooksTotal ?? capture.repoHooks?.count ?? 0
        if capture.hookApplicable == false && repoHookTotal == 0 {
            return "not checked"
        }
        return "ok"
    }

    private func authState(_ capture: CaptureStatus) -> HealthState {
        capture.uploadAuthed == true ? .green : .red
    }

    private func authLabel(_ capture: CaptureStatus) -> String {
        if capture.uploadAuthed == true {
            return "ok"
        }
        if let error = capture.uploadAuthError?.trimmingCharacters(in: .whitespacesAndNewlines), !error.isEmpty {
            return error
        }
        return "missing"
    }

    private func uploadEnvironment(_ api: String?) -> String {
        let value = GXConfig.cleanURL(api ?? "")
        if value.isEmpty {
            return "unknown"
        }
        let lower = value.lowercased()
        if lower.contains("staging") || lower.contains("localhost") || lower.contains("127.0.0.1") {
            return "staging"
        }
        return "prod"
    }

    private func doctorSubmenuItem() -> NSMenuItem {
        let item = submenuItem(title: "Status", submenu: doctorMenu())
        item.attributedTitle = doctorAttributedTitle()
        return item
    }

    private func doctorAttributedTitle() -> NSAttributedString {
        let title = NSMutableAttributedString(
            string: "●  ",
            attributes: [
                .font: NSFont.menuFont(ofSize: 0),
                .foregroundColor: doctorDotColor()
            ]
        )
        title.append(NSAttributedString(
            string: "Status",
            attributes: [.font: NSFont.menuFont(ofSize: 0)]
        ))
        return title
    }

    private func doctorDotColor() -> NSColor {
        healthColor(doctorMenuState())
    }

    private func healthColor(_ state: HealthState) -> NSColor {
        switch state {
        case .green:
            return .systemGreen
        case .yellow:
            return .systemYellow
        case .red:
            return .systemRed
        case .unknown:
            return .systemYellow
        }
    }

    private func statsAgentRows(_ agents: [AgentStats]?) -> [AgentStats] {
        let agents = agents ?? []
        return ["cursor", "codex", "claude"].map { id in
            agents.first { ($0.agent ?? "").lowercased() == id } ??
                AgentStats(agent: id, label: Self.agentLabel(id), health: "yellow", sessions: 0)
        }
    }

    private func agentStatsItem(_ agent: AgentStats) -> NSMenuItem {
        let item = NSMenuItem(title: "", action: nil, keyEquivalent: "")
        item.isEnabled = false
        let health = HealthState(rawValue: (agent.health ?? "").lowercased()) ?? .unknown
        let label = agent.label ?? Self.agentLabel(agent.agent ?? "")
        let title = NSMutableAttributedString(
            string: "●  ",
            attributes: [
                .font: NSFont.menuFont(ofSize: 0),
                .foregroundColor: healthColor(health)
            ]
        )
        title.append(NSAttributedString(
            string: "\(label): \(agent.sessions ?? 0) Sessions",
            attributes: [.font: NSFont.menuFont(ofSize: 0)]
        ))
        item.attributedTitle = title
        return item
    }

    private static func agentLabel(_ agent: String) -> String {
        switch agent.lowercased() {
        case "cursor":
            return "Cursor"
        case "codex":
            return "Codex"
        case "claude":
            return "Claude"
        default:
            return agent.isEmpty ? "Agent" : agent
        }
    }

    private func formatDiskUsed(_ bytes: Int64) -> String {
        let gb = Double(bytes) / 1_073_741_824
        if gb >= 1 {
            return String(format: "%.1fGB", gb)
        }
        let mb = Double(bytes) / 1_048_576
        return String(format: "%.1fMB", max(0, mb))
    }

    private func doctorMenuState() -> HealthState {
        if doctor == nil, doctorError != nil {
            return .red
        }
        return DoctorClient.state(for: doctor)
    }

    private func versionSummary() -> String {
        "\(displayVersion()), CLI \(cliVersion())"
    }

    private func releaseSummary() -> String {
        if let releaseError {
            return "Latest CLI: unavailable (\(releaseError))"
        }
        guard let latestRelease else {
            return lastReleaseChecked == nil ? "Latest CLI: checking..." : "Latest CLI: unavailable"
        }
        if hasRemoteUpdate() {
            return "Update available: \(latestRelease.displayVersion)"
        }
        return "Latest CLI: \(latestRelease.displayVersion)"
    }

    private func hasRemoteUpdate() -> Bool {
        ReleaseClient.isNewer(current: cliInfo, latest: latestRelease)
    }

    private func displayVersion() -> String {
        let tag = plistString("GXGitTagVersion")
        let raw = tag == nil || tag == "dev"
            ? plistString("CFBundleShortVersionString") ?? "dev"
            : tag!
        if raw == "dev" || raw.hasPrefix("v") {
            return raw
        }
        return "v\(raw)"
    }

    private func cliVersion() -> String {
        if let cliInfo {
            return abbreviateVersion(cliInfo.display, maxLength: 12)
        }
        if cliError != nil {
            return "unknown"
        }
        return abbreviateVersion(plistString("GXCLIVersion") ?? "dev", maxLength: 12)
    }

    private func abbreviateVersion(_ value: String, maxLength: Int) -> String {
        if value == "dev" || value.count <= maxLength {
            return value
        }
        return String(value.prefix(maxLength))
    }

    private func plistString(_ key: String) -> String? {
        guard let value = Bundle.main.object(forInfoDictionaryKey: key) as? String else {
            return nil
        }
        let trimmed = value.trimmingCharacters(in: .whitespacesAndNewlines)
        return trimmed.isEmpty ? nil : trimmed
    }

    private func disabled(_ title: String) -> NSMenuItem {
        let item = NSMenuItem(title: title, action: nil, keyEquivalent: "")
        item.isEnabled = false
        return item
    }

    private func actionItem(_ title: String, _ selector: Selector, keyEquivalent: String = "", state: NSControl.StateValue = .off) -> NSMenuItem {
        let item = NSMenuItem(title: title, action: selector, keyEquivalent: keyEquivalent)
        item.target = self
        item.state = state
        return item
    }

    private func submenuItem(title: String, submenu: NSMenu) -> NSMenuItem {
        let item = NSMenuItem(title: title, action: nil, keyEquivalent: "")
        item.submenu = submenu
        return item
    }

    @objc private func runDoctorNow() {
        refreshAll(installCLI: false)
    }

    @objc private func installCLI() {
        updateMessage = "Installing bundled CLI..."
        refreshAll(installCLI: true, forceReleaseCheck: true)
    }

    @objc private func checkForUpdates() {
        updateMessage = "Checking for updates..."
        refreshAll(installCLI: false, forceReleaseCheck: true)
    }

    @objc private func checkForAppUpdates() {
        updaterController.checkForUpdates(nil)
    }

    @objc private func updateCLIFromGitHub() {
        guard !updatingCLI, let release = latestRelease else { return }
        updatingCLI = true
        updateMessage = "Preparing GX \(release.displayVersion)..."
        rebuildMenu()
        worker.async { [weak self] in
            let result = ReleaseClient.update(to: release) { message in
                DispatchQueue.main.async { [weak self] in
                    self?.updateMessage = message
                    self?.rebuildMenu()
                }
            }
            let currentCLI = ReleaseClient.currentCLI()
            DispatchQueue.main.async {
                guard let self else { return }
                self.cliInfo = currentCLI.0
                self.cliError = currentCLI.1
                self.updateMessage = result
                self.lastUpdated = Date()
                self.updatingCLI = false
                self.rebuildMenu()
            }
        }
    }

    @objc private func openConsole() {
        NSWorkspace.shared.open(URL(string: defaultConsoleURL)!)
    }

    @objc private func copyDoctorJSON() {
        copyToPasteboard(doctorRawJSON.isEmpty ? "{}" : doctorRawJSON)
    }

    @objc private func openMCPDocumentation() {
        NSWorkspace.shared.open(URL(string: documentationURL)!)
    }

    @objc private func quit() {
        NSApp.terminate(nil)
    }

    private func copyToPasteboard(_ value: String) {
        NSPasteboard.general.clearContents()
        NSPasteboard.general.setString(value, forType: .string)
    }

    private static let timeFormatter: DateFormatter = {
        let formatter = DateFormatter()
        formatter.timeStyle = .short
        formatter.dateStyle = .none
        return formatter
    }()
}

private let app = NSApplication.shared
private let delegate = GXMenuBarApp()
app.delegate = delegate
app.run()
