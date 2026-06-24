import AppKit
import Foundation

private let pollInterval: TimeInterval = 60
private let defaultAPIURL = "http://localhost:3201"
private let defaultConsoleURL = "https://gx.dev/console"

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

private struct UploadCredentials {
    let apiURL: String
    let token: String
    let orgID: String
}

private enum GXConfig {
    static var gxHome: URL {
        let custom = ProcessInfo.processInfo.environment["GX_HOME"]?.trimmingCharacters(in: .whitespacesAndNewlines) ?? ""
        if !custom.isEmpty {
            return URL(fileURLWithPath: custom, isDirectory: true)
        }
        return URL(fileURLWithPath: NSHomeDirectory(), isDirectory: true).appendingPathComponent(".gx", isDirectory: true)
    }

    static var uploadCredentialsPath: URL {
        gxHome.appendingPathComponent("upload.json")
    }

    static func readUploadCredentials() -> UploadCredentials {
        let env = ProcessInfo.processInfo.environment
        let envAPI = cleanURL(env["GX_API_URL"] ?? "")
        let envToken = (env["GX_UPLOAD_TOKEN"] ?? "").trimmingCharacters(in: .whitespacesAndNewlines)
        let envOrg = (env["GX_ORG_ID"] ?? "").trimmingCharacters(in: .whitespacesAndNewlines)

        var fileAPI = ""
        var fileToken = ""
        var fileOrg = ""
        if let data = try? Data(contentsOf: uploadCredentialsPath),
           let object = try? JSONSerialization.jsonObject(with: data) as? [String: Any] {
            fileAPI = cleanURL(object["api_url"] as? String ?? "")
            fileToken = (object["token"] as? String ?? "").trimmingCharacters(in: .whitespacesAndNewlines)
            fileOrg = (object["org_id"] as? String ?? "").trimmingCharacters(in: .whitespacesAndNewlines)
        }

        return UploadCredentials(
            apiURL: envAPI.isEmpty ? (fileAPI.isEmpty ? defaultAPIURL : fileAPI) : envAPI,
            token: envToken.isEmpty ? fileToken : envToken,
            orgID: envOrg.isEmpty ? fileOrg : envOrg
        )
    }

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
                return "CLI installed at \(installPath.path)"
            }
            return "Bundled gx CLI not found"
        }

        do {
            try FileManager.default.createDirectory(at: installPath.deletingLastPathComponent(), withIntermediateDirectories: true)
            if source.standardizedFileURL.path != installPath.standardizedFileURL.path {
                if FileManager.default.fileExists(atPath: installPath.path) {
                    try FileManager.default.removeItem(at: installPath)
                }
                try FileManager.default.copyItem(at: source, to: installPath)
            }
            try FileManager.default.setAttributes([.posixPermissions: 0o755], ofItemAtPath: installPath.path)
            return "CLI installed at \(installPath.path)"
        } catch {
            return "CLI install failed: \(error.localizedDescription)"
        }
    }
}

private struct DoctorEnvelope: Decodable {
    let doctor: DoctorStatus?
}

private struct DoctorStatus: Decodable {
    let ok: Bool?
    let capture: CaptureStatus?
    let ledger: [LedgerRow]?
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

private struct AuthStatus: Decodable {
    let loggedIn: Bool?
    let login: String?
    let authKind: String?
    let cloudURL: String?
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

    static func authStatus() -> AuthStatus? {
        let result = CommandRunner.run(CLIInstaller.gxExecutable(), ["auth", "status", "--json"], timeout: 30)
        guard result.ok else { return nil }
        return try? JSONDecoder().decode(AuthStatus.self, from: Data(result.stdout.utf8))
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
        let backlog = (capture.pendingExtracts ?? 0) + (capture.pendingSessions ?? 0)
        if backlog > 10 || doctor.ok == false {
            return .yellow
        }
        return .green
    }
}

private struct ActivityEvent {
    let title: String
}

private struct ActivityResult {
    let ok: Bool
    let events: [ActivityEvent]
    let source: String
    let error: String?
}

private enum ActivityClient {
    static func fetch(limit: Int) -> ActivityResult {
        if ProcessInfo.processInfo.environment["GX_ACTIVITY_MOCK"]?.trimmingCharacters(in: .whitespacesAndNewlines) == "1" {
            return ActivityResult(ok: true, events: mockActivity(limit), source: "mock", error: nil)
        }

        let creds = GXConfig.readUploadCredentials()
        guard !creds.token.isEmpty else {
            return ActivityResult(ok: false, events: [], source: "local", error: "Not logged in for activity")
        }
        guard let url = URL(string: "\(creds.apiURL)/v1/activity?limit=\(limit)") else {
            return ActivityResult(ok: false, events: [], source: "local", error: "Invalid activity URL")
        }

        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("application/json", forHTTPHeaderField: "Accept")
        request.setValue("Bearer \(creds.token)", forHTTPHeaderField: "Authorization")
        request.setValue("GX-Menubar/0.1.0", forHTTPHeaderField: "User-Agent")

        let semaphore = DispatchSemaphore(value: 0)
        var loadedData: Data?
        var loadedError: Error?
        let task = URLSession.shared.dataTask(with: request) { data, response, error in
            if let http = response as? HTTPURLResponse, !(200..<300).contains(http.statusCode) {
                loadedError = NSError(domain: "GXActivity", code: http.statusCode, userInfo: [NSLocalizedDescriptionKey: "HTTP \(http.statusCode)"])
            } else {
                loadedData = data
                loadedError = error
            }
            semaphore.signal()
        }
        task.resume()
        if semaphore.wait(timeout: .now() + 15) == .timedOut {
            task.cancel()
            return ActivityResult(ok: false, events: [], source: "api", error: "Activity request timed out")
        }
        if let loadedError {
            return ActivityResult(ok: false, events: mockActivity(min(limit, 3)), source: "fallback", error: loadedError.localizedDescription)
        }
        let events = normalizeActivity(data: loadedData).prefix(limit).map { ActivityEvent(title: activityTitle($0)) }
        return ActivityResult(ok: true, events: Array(events), source: "api", error: nil)
    }

    private static func normalizeActivity(data: Data?) -> [[String: Any]] {
        guard let data,
              let object = try? JSONSerialization.jsonObject(with: data) else {
            return []
        }
        if let array = object as? [[String: Any]] {
            return array
        }
        if let dict = object as? [String: Any] {
            if let events = dict["events"] as? [[String: Any]] {
                return events
            }
            if let items = dict["items"] as? [[String: Any]] {
                return items
            }
        }
        return []
    }

    private static func activityTitle(_ event: [String: Any]) -> String {
        for key in ["title", "summary", "type"] {
            let value = (event[key] as? String ?? "").trimmingCharacters(in: .whitespacesAndNewlines)
            if !value.isEmpty {
                if key == "type" {
                    return value.replacingOccurrences(of: "_", with: " ")
                }
                return value.count > 72 ? String(value.prefix(69)) + "..." : value
            }
        }
        return "Activity event"
    }

    private static func mockActivity(_ limit: Int) -> [ActivityEvent] {
        [
            "PR Summary posted for gx/main",
            "Cursor session uploaded",
            "Capture sync completed"
        ].prefix(limit).map { ActivityEvent(title: $0) }
    }
}

private enum PauseState {
    static var pausePath: URL {
        GXConfig.gxHome.appendingPathComponent("pause-capture")
    }

    static var legacyPausePath: URL {
        GXConfig.gxHome.appendingPathComponent("capture-paused")
    }

    static func isPaused() -> Bool {
        let value = ProcessInfo.processInfo.environment["GX_CAPTURE_PAUSED"]?.lowercased().trimmingCharacters(in: .whitespacesAndNewlines) ?? ""
        if ["1", "true", "yes"].contains(value) {
            return true
        }
        return FileManager.default.fileExists(atPath: pausePath.path) || FileManager.default.fileExists(atPath: legacyPausePath.path)
    }

    static func setPaused(_ paused: Bool) throws {
        try FileManager.default.createDirectory(at: GXConfig.gxHome, withIntermediateDirectories: true)
        if paused {
            try "paused\n".write(to: pausePath, atomically: true, encoding: .utf8)
            try FileManager.default.setAttributes([.posixPermissions: 0o600], ofItemAtPath: pausePath.path)
            return
        }
        for path in [pausePath, legacyPausePath] where FileManager.default.fileExists(atPath: path.path) {
            try FileManager.default.removeItem(at: path)
        }
    }
}

private enum MCPInstructions {
    private static var localMCPPath: String {
        let cwd = URL(fileURLWithPath: FileManager.default.currentDirectoryPath)
        let direct = cwd.appendingPathComponent("mcp/dist/gx-mcp").path
        let parent = cwd.deletingLastPathComponent().appendingPathComponent("mcp/dist/gx-mcp").path
        return firstExecutable([direct, parent]) ?? direct
    }

    private static let applicationMCPPath = "/Applications/GX.app/Contents/Resources/bin/gx-mcp"

    static func mcpExecutablePath() -> String {
        let candidates = [
            bundledMCPPath(),
            localMCPPath,
            applicationMCPPath
        ].compactMap { $0 }
        return firstExecutable(candidates) ?? candidates.first ?? applicationMCPPath
    }

    static func serverReady() -> Bool {
        FileManager.default.isExecutableFile(atPath: mcpExecutablePath())
    }

    static func serverLocationLabel() -> String {
        if mcpExecutablePath() == localMCPPath {
            return "Local MCP binary"
        }
        return "Bundled MCP binary"
    }

    static func cursorCommand() -> String {
        "cursor mcp add gx -- env GX_BINARY=\(shellQuote(CLIInstaller.gxExecutable())) \(shellQuote(mcpExecutablePath()))"
    }

    static func codexCommand() -> String {
        "codex mcp add gx --env \(shellQuote("GX_BINARY=\(CLIInstaller.gxExecutable())")) -- \(shellQuote(mcpExecutablePath()))"
    }

    static func claudeJSON() -> String {
        let args = [
            "GX_BINARY=\(CLIInstaller.gxExecutable())",
            mcpExecutablePath()
        ]
        let value: [String: Any] = [
            "mcpServers": [
                "gx": [
                    "command": "env",
                    "args": args
                ]
            ]
        ]
        let data = try? JSONSerialization.data(withJSONObject: value, options: [.prettyPrinted, .sortedKeys])
        return data.flatMap { String(data: $0, encoding: .utf8) } ?? "{}"
    }

    static func fullText() -> String {
        """
        MCP binary:
        \(mcpExecutablePath())

        Cursor:
        \(cursorCommand())

        Codex:
        \(codexCommand())

        Claude Desktop JSON:
        \(claudeJSON())

        The app installs the CLI at:
        \(CLIInstaller.installPath.path)

        MCP runs over stdio from a standalone gx-mcp binary. Local menu-bar runs use the repo-local mcp/dist/gx-mcp build when present; installed app runs use bundled resources. Cloud context uses gx auth credentials from disk.
        """
    }

    private static func bundledMCPPath() -> String? {
        guard let resourceURL = Bundle.main.resourceURL else { return nil }
        return resourceURL.appendingPathComponent("bin", isDirectory: true)
            .appendingPathComponent("gx-mcp").path
    }

    private static func firstExecutable(_ paths: [String]) -> String? {
        paths.first { FileManager.default.isExecutableFile(atPath: $0) }
    }

    private static func shellQuote(_ value: String) -> String {
        "'" + value.replacingOccurrences(of: "'", with: "'\\''") + "'"
    }
}

private final class GXMenuBarApp: NSObject, NSApplicationDelegate {
    private let statusItem = NSStatusBar.system.statusItem(withLength: NSStatusItem.variableLength)
    private let worker = DispatchQueue(label: "dev.gx.menubar.worker", qos: .utility)
    private var timer: Timer?
    private var refreshing = false

    private var cliStatus = "CLI not installed yet"
    private var doctor: DoctorStatus?
    private var doctorRawJSON = ""
    private var doctorError: String?
    private var authStatus: AuthStatus?
    private var activity = ActivityResult(ok: false, events: [], source: "local", error: "Not loaded yet")
    private var capturePaused = PauseState.isPaused()
    private var lastUpdated: Date?

    func applicationDidFinishLaunching(_ notification: Notification) {
        NSApp.setActivationPolicy(.accessory)
        configureStatusItem()
        rebuildMenu()
        refreshAll(installCLI: true)
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

    private func refreshAll(installCLI: Bool) {
        guard !refreshing else { return }
        refreshing = true
        rebuildMenu()
        worker.async { [weak self] in
            let installMessage = installCLI ? CLIInstaller.installBundledCLI() : self?.cliStatus
            let doctorLoad = DoctorClient.fetch()
            let auth = DoctorClient.authStatus()
            let activity = ActivityClient.fetch(limit: 10)
            let paused = PauseState.isPaused()

            DispatchQueue.main.async {
                guard let self else { return }
                if let installMessage {
                    self.cliStatus = installMessage
                }
                self.doctor = doctorLoad.envelope?.doctor
                self.doctorRawJSON = doctorLoad.rawJSON
                self.doctorError = doctorLoad.error
                self.authStatus = auth
                self.activity = activity
                self.capturePaused = paused
                self.lastUpdated = Date()
                self.refreshing = false
                self.rebuildMenu()
            }
        }
    }

    private func rebuildMenu() {
        let menu = NSMenu()
        menu.addItem(disabled("GX \(appVersion())"))
        menu.addItem(disabled(cliStatus))
        if let updated = lastUpdated {
            menu.addItem(disabled("Updated \(Self.timeFormatter.string(from: updated))"))
        } else if refreshing {
            menu.addItem(disabled("Updating..."))
        }
        menu.addItem(.separator())

        menu.addItem(submenuItem(title: doctorLabel(), submenu: doctorMenu()))
        menu.addItem(submenuItem(title: "Stats", submenu: statsMenu()))
        menu.addItem(submenuItem(title: "Activity", submenu: activityMenu()))
        menu.addItem(.separator())

        menu.addItem(actionItem("Run Doctor Now", #selector(runDoctorNow)))
        menu.addItem(actionItem(capturePaused ? "Resume Capture" : "Pause Capture", #selector(toggleCapturePause), state: capturePaused ? .on : .off))
        menu.addItem(actionItem("Install or Update CLI", #selector(installCLI)))
        menu.addItem(.separator())

        menu.addItem(actionItem("Open GX Console", #selector(openConsole)))
        menu.addItem(submenuItem(title: "MCP Setup", submenu: mcpMenu()))
        menu.addItem(.separator())
        menu.addItem(actionItem("Reveal GX in Finder", #selector(revealInFinder)))
        menu.addItem(actionItem("Quit", #selector(quit), keyEquivalent: "q"))
        statusItem.menu = menu
    }

    private func doctorMenu() -> NSMenu {
        let menu = NSMenu()
        if let error = doctorError, doctor == nil {
            menu.addItem(disabled(error))
        }
        guard let capture = doctor?.capture else {
            menu.addItem(disabled("Run gx init and gx login to finish setup"))
            return menu
        }
        if capture.hookApplicable == false {
            menu.addItem(disabled("Pre-push hook: not checked outside repo"))
        } else {
            menu.addItem(disabled(capture.hookInstalled == true ? "Pre-push hook: ok" : "Pre-push hook: missing"))
        }
        let repoHookTotal = capture.repoHooksTotal ?? capture.repoHooks?.count ?? 0
        if repoHookTotal == 0 {
            menu.addItem(disabled("Registered repo hooks: none"))
        } else if capture.repoHooksOK == true {
            menu.addItem(disabled("Registered repo hooks: \(repoHookTotal) ok"))
        } else {
            menu.addItem(disabled("Registered repo hooks: \(capture.repoHooksMissing ?? 0) missing, \(capture.repoHooksUnreachable ?? 0) unreachable"))
            if let repoHooks = capture.repoHooks {
                for repoHook in repoHooks.filter({ $0.hookInstalled != true || $0.gitReachable == false }).prefix(4) {
                    menu.addItem(disabled("- \(repoHook.repoRoot ?? "repo")"))
                }
            }
        }
        if capture.uploadAuthed == true {
            menu.addItem(disabled("Upload auth: ok"))
        } else if let error = capture.uploadAuthError, !error.isEmpty {
            menu.addItem(disabled("Upload auth: \(error)"))
        } else {
            menu.addItem(disabled("Upload auth: missing"))
        }
        menu.addItem(disabled(capture.cursorReachable == true ? "Cursor vscdb: ok" : "Cursor vscdb: missing"))
        let backlog = (capture.pendingExtracts ?? 0) + (capture.pendingSessions ?? 0)
        menu.addItem(disabled("Staging backlog: \(backlog) pending"))
        if let diskFree = capture.diskFreeGB {
            menu.addItem(disabled(capture.diskWarn == true ? "Disk free: \(diskFree) GB low" : "Disk free: \(diskFree) GB"))
        }
        if let api = capture.uploadAPI, !api.isEmpty {
            menu.addItem(disabled("Upload API: \(api)"))
        }
        menu.addItem(.separator())
        menu.addItem(actionItem("Copy gx doctor JSON", #selector(copyDoctorJSON)))
        return menu
    }

    private func statsMenu() -> NSMenu {
        let menu = NSMenu()
        if let summary = doctor?.diagnose?.summary {
            menu.addItem(disabled(summary))
            menu.addItem(.separator())
        }
        if let capture = doctor?.capture {
            menu.addItem(disabled("Pending extracts: \(capture.pendingExtracts ?? 0)"))
            menu.addItem(disabled("Pending sessions: \(capture.pendingSessions ?? 0)"))
            if let diskFree = capture.diskFreeGB {
                menu.addItem(disabled("Disk free: \(diskFree) GB"))
            }
        } else {
            menu.addItem(disabled("Stats unavailable"))
        }
        if let ledger = doctor?.ledger, !ledger.isEmpty {
            menu.addItem(.separator())
            for row in ledger.prefix(5) {
                let agent = row.agent ?? "agent"
                let status = row.status ?? "unknown"
                let calls = row.calls ?? "0"
                menu.addItem(disabled("\(agent): \(status), \(calls) calls"))
            }
        }
        return menu
    }

    private func activityMenu() -> NSMenu {
        let menu = NSMenu()
        if activity.events.isEmpty {
            menu.addItem(disabled(activity.error ?? "No recent activity"))
        } else {
            for event in activity.events.prefix(10) {
                menu.addItem(disabled(event.title))
            }
        }
        menu.addItem(.separator())
        menu.addItem(actionItem("Refresh Activity", #selector(runDoctorNow)))
        return menu
    }

    private func mcpMenu() -> NSMenu {
        let menu = NSMenu()
        let serverReady = MCPInstructions.serverReady()
        menu.addItem(disabled(MCPInstructions.serverLocationLabel() + (serverReady ? ": ready" : ": missing")))
        menu.addItem(.separator())
        menu.addItem(actionItem("Show Instructions", #selector(showMCPInstructions)))
        menu.addItem(actionItem("Copy Cursor Command", #selector(copyCursorMCPCommand)))
        menu.addItem(actionItem("Copy Codex Command", #selector(copyCodexMCPCommand)))
        menu.addItem(actionItem("Copy Claude JSON", #selector(copyClaudeMCPJSON)))
        return menu
    }

    private func doctorLabel() -> String {
        switch DoctorClient.state(for: doctor) {
        case .green:
            return "Doctor: OK"
        case .yellow:
            return "Doctor: warnings"
        case .red:
            return "Doctor: needs attention"
        case .unknown:
            return doctor == nil ? "Doctor: unavailable" : "Doctor: unknown"
        }
    }

    private func appVersion() -> String {
        Bundle.main.object(forInfoDictionaryKey: "CFBundleShortVersionString") as? String ?? "0.1.0"
    }

    private func consoleURL() -> URL {
        let envConsole = GXConfig.cleanURL(ProcessInfo.processInfo.environment["GX_CONSOLE_URL"] ?? "")
        if !envConsole.isEmpty, let url = URL(string: envConsole) {
            return url
        }
        let rawCloudURL = GXConfig.cleanURL(authStatus?.cloudURL ?? ProcessInfo.processInfo.environment["GX_CLOUD_URL"] ?? "")
        if !rawCloudURL.isEmpty {
            let base = rawCloudURL.replacingOccurrences(of: #"/gx/pr$"#, with: "", options: .regularExpression)
            if let url = URL(string: base) {
                return url
            }
        }
        return URL(string: defaultConsoleURL)!
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
        refreshAll(installCLI: true)
    }

    @objc private func toggleCapturePause() {
        do {
            try PauseState.setPaused(!capturePaused)
            capturePaused.toggle()
        } catch {
            doctorError = "Pause toggle failed: \(error.localizedDescription)"
        }
        rebuildMenu()
    }

    @objc private func openConsole() {
        NSWorkspace.shared.open(consoleURL())
    }

    @objc private func revealInFinder() {
        NSWorkspace.shared.activateFileViewerSelecting([Bundle.main.bundleURL])
    }

    @objc private func copyDoctorJSON() {
        copyToPasteboard(doctorRawJSON.isEmpty ? "{}" : doctorRawJSON)
    }

    @objc private func showMCPInstructions() {
        NSApp.activate(ignoringOtherApps: true)
        let alert = NSAlert()
        alert.messageText = "GX MCP Setup"
        alert.informativeText = MCPInstructions.fullText()
        alert.addButton(withTitle: "Copy Cursor Command")
        alert.addButton(withTitle: "Copy Codex Command")
        alert.addButton(withTitle: "Copy Claude JSON")
        alert.addButton(withTitle: "OK")
        let response = alert.runModal()
        if response == .alertFirstButtonReturn {
            copyToPasteboard(MCPInstructions.cursorCommand())
        } else if response == .alertSecondButtonReturn {
            copyToPasteboard(MCPInstructions.codexCommand())
        } else if response == .alertThirdButtonReturn {
            copyToPasteboard(MCPInstructions.claudeJSON())
        }
    }

    @objc private func copyCursorMCPCommand() {
        copyToPasteboard(MCPInstructions.cursorCommand())
    }

    @objc private func copyCodexMCPCommand() {
        copyToPasteboard(MCPInstructions.codexCommand())
    }

    @objc private func copyClaudeMCPJSON() {
        copyToPasteboard(MCPInstructions.claudeJSON())
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
