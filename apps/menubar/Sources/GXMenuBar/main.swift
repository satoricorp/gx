import AppKit
import Foundation

private let pollInterval: TimeInterval = 60
private let defaultConsoleURL = "https://gx.run"
private let documentationURL = "https://gx.run/documentation"

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

    static func cursorCommand() -> String {
        "cursor mcp add gx -- env GX_BINARY=\(shellQuote(CLIInstaller.gxExecutable())) \(shellQuote(mcpExecutablePath()))"
    }

    static func codexCommand() -> String {
        "codex mcp add gx --env \(shellQuote("GX_BINARY=\(CLIInstaller.gxExecutable())")) -- \(shellQuote(mcpExecutablePath()))"
    }

    static func claudeCodeCommand() -> String {
        "claude mcp add gx -- env GX_BINARY=\(shellQuote(CLIInstaller.gxExecutable())) \(shellQuote(mcpExecutablePath()))"
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

    private var doctor: DoctorStatus?
    private var doctorRawJSON = ""
    private var doctorError: String?
    private var authStatus: AuthStatus?
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
            if installCLI {
                _ = CLIInstaller.installBundledCLI()
            }
            let doctorLoad = DoctorClient.fetch()
            let auth = DoctorClient.authStatus()

            DispatchQueue.main.async {
                guard let self else { return }
                self.doctor = doctorLoad.envelope?.doctor
                self.doctorRawJSON = doctorLoad.rawJSON
                self.doctorError = doctorLoad.error
                self.authStatus = auth
                self.lastUpdated = Date()
                self.refreshing = false
                self.rebuildMenu()
            }
        }
    }

    private func rebuildMenu() {
        let menu = NSMenu()
        menu.addItem(disabled("GX \(appVersion())"))
        if let updated = lastUpdated {
            menu.addItem(disabled("Updated \(Self.timeFormatter.string(from: updated))"))
        } else if refreshing {
            menu.addItem(disabled("Updating..."))
        }
        menu.addItem(.separator())

        menu.addItem(doctorSubmenuItem())
        menu.addItem(submenuItem(title: "Stats", submenu: statsMenu()))
        menu.addItem(.separator())

        menu.addItem(actionItem("Update CLI", #selector(installCLI)))
        menu.addItem(.separator())

        menu.addItem(actionItem("open gx.run", #selector(openConsole)))
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
        } else {
            menu.addItem(disabled("Run gx init and gx login to finish setup"))
        }
        menu.addItem(.separator())
        menu.addItem(actionItem("Run Doctor", #selector(runDoctorNow)))
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

    private func mcpMenu() -> NSMenu {
        let menu = NSMenu()
        menu.addItem(actionItem("Show Instructions", #selector(openMCPDocumentation)))
        menu.addItem(.separator())
        menu.addItem(actionItem("Copy Cursor Install Command", #selector(copyCursorMCPCommand)))
        menu.addItem(actionItem("Copy Codex Install Command", #selector(copyCodexMCPCommand)))
        menu.addItem(actionItem("Copy Claude Code Install Command", #selector(copyClaudeCodeMCPCommand)))
        return menu
    }

    private func doctorSubmenuItem() -> NSMenuItem {
        let item = submenuItem(title: "Doctor", submenu: doctorMenu())
        item.attributedTitle = doctorAttributedTitle()
        return item
    }

    private func doctorAttributedTitle() -> NSAttributedString {
        let title = NSMutableAttributedString(
            string: "Doctor",
            attributes: [.font: NSFont.menuFont(ofSize: 0)]
        )
        title.append(NSAttributedString(
            string: "  •",
            attributes: [
                .font: NSFont.systemFont(ofSize: 8, weight: .bold),
                .foregroundColor: doctorDotColor(),
                .baselineOffset: 1
            ]
        ))
        return title
    }

    private func doctorDotColor() -> NSColor {
        switch doctorMenuState() {
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

    private func doctorMenuState() -> HealthState {
        if doctor == nil, doctorError != nil {
            return .red
        }
        return DoctorClient.state(for: doctor)
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

    @objc private func openConsole() {
        NSWorkspace.shared.open(consoleURL())
    }

    @objc private func copyDoctorJSON() {
        copyToPasteboard(doctorRawJSON.isEmpty ? "{}" : doctorRawJSON)
    }

    @objc private func openMCPDocumentation() {
        NSWorkspace.shared.open(URL(string: documentationURL)!)
    }

    @objc private func copyCursorMCPCommand() {
        copyToPasteboard(MCPInstructions.cursorCommand())
    }

    @objc private func copyCodexMCPCommand() {
        copyToPasteboard(MCPInstructions.codexCommand())
    }

    @objc private func copyClaudeCodeMCPCommand() {
        copyToPasteboard(MCPInstructions.claudeCodeCommand())
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
