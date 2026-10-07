package gxtest

import "os"

// EnclosingGitEnv are the variables through which a git process that launches
// a test binary pins that binary's own git commands to its repository and to
// the commit it is making.
//
// Git exports them to every hook, and to `rebase --exec` and `bisect run`. In
// a linked worktree — every checkout under .claude/worktrees is one — that
// includes an absolute GIT_DIR, and a child git honours GIT_DIR over its
// working directory. So inside a hook, `git -C <t.TempDir()> config user.name
// "Test User"` does not configure the temp repository at all: it writes the
// enclosing repository's shared .git/config.
//
// The first group is git's own answer to which variables tie a command to one
// repository: the list `git rev-parse --local-env-vars` prints, which git
// clears itself before it runs a command in a different repository, such as a
// submodule. TestEnclosingGitEnvCoversGitsOwnList keeps it current. The second
// group is the identity and clock of the commit being made. Git exports them
// to commit hooks, they outrank the user.name a test configures, and a frozen
// GIT_AUTHOR_DATE gives every commit a test makes the same timestamp.
var EnclosingGitEnv = []string{
	"GIT_ALTERNATE_OBJECT_DIRECTORIES",
	"GIT_CONFIG",
	"GIT_CONFIG_PARAMETERS",
	"GIT_CONFIG_COUNT",
	"GIT_OBJECT_DIRECTORY",
	"GIT_DIR",
	"GIT_WORK_TREE",
	"GIT_IMPLICIT_WORK_TREE",
	"GIT_GRAFT_FILE",
	"GIT_INDEX_FILE",
	"GIT_NO_REPLACE_OBJECTS",
	"GIT_REPLACE_REF_BASE",
	"GIT_PREFIX",
	"GIT_SHALLOW_FILE",
	"GIT_COMMON_DIR",

	"GIT_AUTHOR_NAME",
	"GIT_AUTHOR_EMAIL",
	"GIT_AUTHOR_DATE",
	"GIT_COMMITTER_NAME",
	"GIT_COMMITTER_EMAIL",
	"GIT_COMMITTER_DATE",
}

// DetachFromEnclosingGit makes a test binary's git commands land in the
// repositories its tests name, even when the binary runs inside a git hook.
// Call it first in TestMain, before anything in the process runs git.
//
// Before this existed, the CLI suite once ran as a prepare-commit-msg hook in
// a linked worktree, with git's GIT_DIR in its environment, and wrote
// `core.bare = true` and `user.name = Test User` into gx's own .git/config
// (2026-08-15). The bare flag broke the checkout at once and was reverted by
// hand. The identity went unnoticed until 2026-10-06, so for seven weeks every
// commit made from that clone carried Test User <test@example.com>.
// Rerunning the suite inside a hook against a decoy repository showed the
// rest: test commits on the decoy's branches, gx's lifecycle hooks installed
// into it, and a fake origin in its config.
//
// It has to be the process and not the git helpers, for two reasons. There
// are more than twenty package-local runGit variants, and a fix that each of
// them has to carry is a fix the next one will forget. More decisively, the code under
// test runs git too, and it inherits os.Environ() like everything else: the
// hooks that the run above installed came from gx's own init path, which no
// test helper sits in front of.
//
// TestEveryPackageThatRunsGitDetachesFromTheEnclosingGit fails for any package
// whose tests run git without calling this from TestMain.
func DetachFromEnclosingGit() {
	for _, key := range EnclosingGitEnv {
		_ = os.Unsetenv(key)
	}
}
