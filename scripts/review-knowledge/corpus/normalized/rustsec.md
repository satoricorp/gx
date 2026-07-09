`cargo-audit`

 Audit
 `Cargo.lock`
 files for crates with security vulnerabilities.

 Get started

 ```
> cargo audit
 Scanning Cargo.lock for vulnerabilities (4 crate dependencies)
Crate: lz4-sys
Version: 1.9.3
Title: Memory corruption in liblz4
Date: 2022-08-25
ID: RUSTSEC-2022-0051
URL: https://rustsec.org/advisories/RUSTSEC-2022-0051
Solution: Upgrade to >=1.9.4
Dependency tree:
lz4-sys 1.9.3
└── crate 0.1.0
error: 1 vulnerability found!

```

 `cargo-deny`

 Audit
 `Cargo.lock` files for crates with security
 vulnerabilities, limit the usage of particular dependencies, their licenses, sources to download
 from, detect multiple versions of same packages in the dependency tree and more.

 Get started
