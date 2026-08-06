package codereview

import "testing"

func TestIsVendoredOrGeneratedPath(t *testing.T) {
	vendored := []string{
		"dotcom/node_modules/zod/index.js",
		"node_modules/.bin/tsc",
		"dotcom/.next/dev/server/manifest.json",
		"web/dist/main.js",
		"api/target/debug/app",
		"py/.venv/lib/site-packages/x.py",
		"svc/__pycache__/mod.cpython-312.pyc",
		"ios/Pods/Alamofire/Source/A.swift",
		"infra/.terraform/providers/aws",
		"app/coverage/lcov-report/index.html",
	}
	for _, path := range vendored {
		if !isVendoredOrGeneratedPath(path) {
			t.Errorf("isVendoredOrGeneratedPath(%q) = false, want true", path)
		}
	}

	source := []string{
		"server/src/index.ts",
		"internal/codereview/review.go",
		"docs/cli.mdx",
		// Substring matches must not count: only whole path segments.
		"src/build-config.ts",
		"src/distributed/queue.go",
		"app/outbound/mailer.ts",
		"lib/vendored-notes.md",
		"src/components/BuildBanner.tsx",
	}
	for _, path := range source {
		if isVendoredOrGeneratedPath(path) {
			t.Errorf("isVendoredOrGeneratedPath(%q) = true, want false", path)
		}
	}

	if isVendoredOrGeneratedPath("") || isVendoredOrGeneratedPath("   ") {
		t.Error("empty path must not be treated as vendored")
	}
}
