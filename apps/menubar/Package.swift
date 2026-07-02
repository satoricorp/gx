// swift-tools-version: 5.9

import PackageDescription

let localSparklePath = Context.environment["GX_LOCAL_SPARKLE_XCFRAMEWORK"] ?? ""
let useLocalSparkle = !localSparklePath.isEmpty
let sparkleDependency: Target.Dependency = useLocalSparkle
    ? "Sparkle"
    : .product(name: "Sparkle", package: "Sparkle")
let packageDependencies: [Package.Dependency] = useLocalSparkle
    ? []
    : [.package(url: "https://github.com/sparkle-project/Sparkle", from: "2.7.0")]
let packageTargets: [Target] = (useLocalSparkle
    ? [.binaryTarget(name: "Sparkle", path: localSparklePath)]
    : []) + [
        .executableTarget(
            name: "GXMenuBar",
            dependencies: [
                sparkleDependency
            ],
            linkerSettings: [
                .unsafeFlags(["-Xlinker", "-rpath", "-Xlinker", "@executable_path/../Frameworks"])
            ]
        )
    ]

let package = Package(
    name: "GXMenuBar",
    platforms: [
        .macOS(.v13)
    ],
    products: [
        .executable(name: "GXMenuBar", targets: ["GXMenuBar"])
    ],
    dependencies: packageDependencies,
    targets: packageTargets
)
