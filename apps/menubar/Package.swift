// swift-tools-version: 5.9

import PackageDescription

let package = Package(
    name: "GXMenuBar",
    platforms: [
        .macOS(.v13)
    ],
    products: [
        .executable(name: "GXMenuBar", targets: ["GXMenuBar"])
    ],
    targets: [
        .executableTarget(name: "GXMenuBar")
    ]
)
