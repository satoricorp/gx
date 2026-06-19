package clitui

import (
	"fmt"
	"strings"

	"github.com/satoricorp/gx/internal/vcs"
)

// UseInteractive reports whether gx should run bubbletea pickers.
func UseInteractive(in any) bool {
	return useInteractive(in)
}

func renderCommandLine(invocation string, accentCommand bool) string {
	parts := strings.Fields(invocation)
	if len(parts) == 0 {
		return ""
	}
	line := prompt("$") + " "
	if accentCommand {
		line += command(strings.Join(parts, " "))
	} else {
		line += text(invocation)
	}
	return line
}

// RenderStatusSummary renders the default gx status output.
func RenderStatusSummary(snapshot vcs.StatusSnapshot) string {
	if len(snapshot.Bookmarks) == 0 {
		return muted("No GX work units recorded yet.")
	}

	lines := []string{
		renderCommandLine("gx status", true),
		muted(fmt.Sprintf("%d bookmarks · %s", len(snapshot.Bookmarks), snapshot.RepoLabel)),
	}

	for _, bookmark := range snapshot.Bookmarks {
		lines = append(lines, renderBookmarkSummary(bookmark, snapshot.BaseRef)...)
	}
	lines = append(lines, muted(fmt.Sprintf("⎇ base · ⌂ local · ⇡ on origin")))
	return strings.Join(lines, "\n")
}

func renderBookmarkSummary(bookmark vcs.BookmarkSnapshot, baseRef string) []string {
	marker, nameStyle := bookmarkMarker(bookmark)
	title := nameStyle(bookmark.Stack.Name)
	if bookmark.Current {
		title += "  " + muted("← current")
	}
	lines := []string{
		fmt.Sprintf("%s  %s", marker, title),
		"     " + muted(bookmarkSummaryDetail(bookmark, baseRef)),
	}
	return lines
}

func bookmarkMarker(bookmark vcs.BookmarkSnapshot) (string, func(string) string) {
	if bookmark.Current {
		return accent("●"), text
	}
	if bookmark.ChangeCount > 0 && bookmark.ApprovedCount == bookmark.ChangeCount {
		return success("●"), secondary
	}
	return muted("○"), secondary
}

func bookmarkSummaryDetail(bookmark vcs.BookmarkSnapshot, baseRef string) string {
	ref := firstNonEmpty(bookmark.Stack.BaseRef, baseRef, "main")
	parts := []string{"⎇ " + ref}
	if bookmark.FileCount > 0 {
		parts = append(parts, fmt.Sprintf("%d files", bookmark.FileCount))
	} else if bookmark.ChangeCount > 0 {
		parts = append(parts, fmt.Sprintf("%d changes", bookmark.ChangeCount))
	}
	if bookmark.ChangeCount > 0 {
		parts = append(parts, fmt.Sprintf("%d/%d approved", bookmark.ApprovedCount, bookmark.ChangeCount))
	}
	parts = append(parts, vcsBookmarkSyncIcons(bookmark.Stack))
	return strings.Join(parts, " · ")
}

func vcsBookmarkSyncIcons(body vcs.StackInfo) string {
	icons := "⌂"
	if body.RemoteRef != nil && strings.TrimSpace(*body.RemoteRef) != "" || body.Status == "published" {
		icons += "⇡"
	}
	return icons
}

func renderBookmarkHeader(snapshot vcs.StatusSnapshot, bookmark vcs.BookmarkSnapshot, keys string) string {
	currentName := ""
	if bookmark.Current {
		currentName = bookmark.Stack.Name
	}
	for _, item := range snapshot.Bookmarks {
		if item.Current {
			currentName = item.Stack.Name
			break
		}
	}
	return strings.Join([]string{
		secondary(fmt.Sprintf("repo %s", snapshot.RepoLabel)),
		secondary(fmt.Sprintf("bookmarks %d", len(snapshot.Bookmarks))),
		muted("current " + currentName),
		hint(keys),
	}, "  ")
}

func renderBookmarkMeta(bookmark vcs.BookmarkSnapshot) string {
	meta := vcsBookmarkMeta(bookmark)
	if bookmark.Current {
		meta += " · current"
	}
	return muted(meta)
}

func vcsBookmarkMeta(bookmark vcs.BookmarkSnapshot) string {
	body := bookmark.Stack
	parts := []string{
		"alias " + body.Alias,
		"base " + firstNonEmpty(body.BaseRef, "main"),
		fmtChanges(bookmark.ChangeCount),
		fmtApproved(bookmark.ApprovedCount, bookmark.ChangeCount),
		"sync " + syncTargets(body),
	}
	return strings.Join(parts, " · ")
}

func syncTargets(body vcs.StackInfo) string {
	parts := []string{"local"}
	if body.RemoteRef != nil && strings.TrimSpace(*body.RemoteRef) != "" || body.Status == "published" {
		parts = append(parts, "origin")
	}
	return strings.Join(parts, ",")
}

func renderRevisionLine(rev vcs.RevisionSnapshot, selected bool) string {
	prefix := "     "
	if selected {
		prefix = accent("›") + " "
	}
	left := muted(fmt.Sprintf("r%d", rev.Index))
	if rev.ShortID != "" {
		left = muted(revisionSyncArrow(rev) + " " + rev.ShortID)
	}
	parts := []string{prefix + left, value(rev.Description)}
	if rev.Working {
		parts = append(parts, muted("working change"))
	} else if rev.SyncNote != "" {
		parts = append(parts, muted(rev.SyncNote))
	}
	return strings.Join(parts, "  ")
}

func renderRevisionPickerLine(rev vcs.RevisionSnapshot, selected bool) string {
	id := rev.ShortID
	if id == "" {
		id = fmt.Sprintf("r%d", rev.Index)
	}
	line := pickerPrefix(selected) + muted(revisionSyncArrow(rev)+" "+id) + "  " + value(rev.Description)
	if rev.Working {
		line += "  " + secondary("@")
	} else if rev.Published && !hasRemoteSync(rev) {
		line += "  " + secondary("origin")
	} else if rev.SyncNote != "" {
		line += "  " + secondary(compactSyncNote(rev.SyncNote))
	}
	return line
}

func pickerPrefix(selected bool) string {
	if selected {
		return accent("›") + "  "
	}
	return "   "
}

func compactSyncNote(note string) string {
	note = strings.TrimSpace(note)
	switch note {
	case "local only":
		return "local"
	case "working change", "working target":
		return "@"
	case "origin only":
		return "origin"
	default:
		return strings.ReplaceAll(note, ",", "+")
	}
}

func hasRemoteSync(rev vcs.RevisionSnapshot) bool {
	return strings.Contains(rev.SyncNote, "origin")
}

func revisionSyncArrow(rev vcs.RevisionSnapshot) string {
	if rev.Published || hasRemoteSync(rev) {
		return "↓"
	}
	return "↑"
}

// RenderStatusInteractive renders the gx status explorer frame.
func RenderStatusInteractive(snapshot vcs.StatusSnapshot, cursor int) string {
	if len(snapshot.Bookmarks) == 0 {
		return RenderStatusSummary(snapshot)
	}
	if cursor < 0 {
		cursor = 0
	}
	if cursor >= len(snapshot.Bookmarks) {
		cursor = len(snapshot.Bookmarks) - 1
	}

	lines := []string{
		renderCommandLine("gx status", true),
		"    " + renderBookmarkHeader(snapshot, snapshot.Bookmarks[cursor], "j/k move bookmark"),
	}

	for index, bookmark := range snapshot.Bookmarks {
		if index == cursor {
			lines = append(lines, "    "+accent("›")+" "+value(bookmark.Stack.Name)+"  "+renderBookmarkMeta(bookmark))
			for _, rev := range bookmark.Revisions {
				lines = append(lines, "    "+renderRevisionLine(rev, false))
			}
			continue
		}
		lines = append(lines, "         "+value(bookmark.Stack.Name)+"  "+renderBookmarkMeta(bookmark))
		if len(bookmark.Revisions) > 0 {
			hidden := len(bookmark.Revisions)
			lines = append(lines, "             "+hint(fmt.Sprintf("%d revisions hidden in summary", hidden)))
		}
	}
	if cursor < len(snapshot.Bookmarks)-1 {
		lines = append(lines, "    "+divider())
	}
	return strings.Join(lines, "\n")
}

// RenderModifyInteractive renders the gx edit revision picker frame.
func RenderModifyInteractive(snapshot vcs.StatusSnapshot, revisions []vcs.RevisionSnapshot, cursor int) string {
	if len(revisions) == 0 {
		return muted("No revisions available to edit.")
	}
	if cursor < 0 {
		cursor = 0
	}
	if cursor >= len(revisions) {
		cursor = len(revisions) - 1
	}

	bookmark := currentBookmark(snapshot)
	header := secondary(fmt.Sprintf("%s · %d revs", bookmark.Stack.Name, len(revisions))) + "  " +
		hint("j/k move · enter edit · q quit")

	lines := []string{
		renderCommandLine("gx edit", false),
		"    " + header,
	}
	for index, rev := range revisions {
		lines = append(lines, "    "+renderRevisionPickerLine(rev, index == cursor))
	}
	return strings.Join(lines, "\n")
}

// RenderModifyResult renders output after a successful gx edit.
func RenderModifyResult(snapshot vcs.StatusSnapshot, result vcs.ModifyResult, revLabel string) string {
	bookmark := currentBookmark(snapshot)
	lines := []string{
		renderCommandLine("gx edit "+revLabel, false),
	}
	kv := []struct{ label, val, meta string }{
		{"revision", revLabel, shortID(result.CurrentChange.ChangeID, 12)},
		{"message", result.CurrentChange.Description, ""},
		{"bookmark", bookmark.Stack.Name, mutedMeta(firstNonEmpty(bookmark.Stack.BaseRef, snapshot.BaseRef), bookmark.Stack)},
	}
	for _, row := range kv {
		lines = append(lines, renderKV(row.label, row.val, row.meta))
	}
	lines = append(lines, "    "+divider())
	lines = append(lines, renderKV("next", fmt.Sprintf(`gx add -m %q`, result.CurrentChange.Description), ""))
	return strings.Join(lines, "\n")
}

func renderKV(label, val, meta string) string {
	line := "    " + muted(label) + "  " + value(val)
	if meta != "" {
		line += "  " + meta
	}
	return line
}

func mutedMeta(baseRef string, stack vcs.StackInfo) string {
	return muted("base " + firstNonEmpty(baseRef, "main") + " · sync " + syncTargets(stack))
}

// RenderStatusAgent renders machine-readable gx status output.
func RenderStatusAgent(snapshot vcs.StatusSnapshot, filter string) string {
	lines := []string{
		renderCommandLine("gx status --agent", false),
		secondary(fmt.Sprintf(`repo=%s bookmarks=%d current=%q`, snapshot.RepoLabel, len(snapshot.Bookmarks), currentBookmarkName(snapshot))),
	}
	filter = strings.TrimSpace(strings.ToLower(filter))
	for _, bookmark := range snapshot.Bookmarks {
		if filter != "" && !bookmarkMatchesFilter(bookmark, filter) {
			continue
		}
		lines = append(lines, renderBookmarkAgent(bookmark)...)
		if filter != "" {
			break
		}
		if len(bookmark.Revisions) > 3 && filter == "" {
			lines = append(lines, muted(fmt.Sprintf(`revision omitted count=%d reason=summary`, len(bookmark.Revisions)-3)))
		}
	}
	if filter != "" {
		lines = append(lines, divider(), renderCommandLine("gx status --agent "+filter, false))
		bookmark := findBookmark(snapshot, filter)
		if bookmark != nil {
			lines = append(lines,
				muted(fmt.Sprintf(`bookmark=%q repo=%s`, bookmark.Stack.Name, snapshot.RepoLabel)),
				muted(fmt.Sprintf(`base=%s changes=%d approved=%d/%d sync=%s current=%t`,
					firstNonEmpty(bookmark.Stack.BaseRef, snapshot.BaseRef),
					bookmark.ChangeCount,
					bookmark.ApprovedCount,
					bookmark.ChangeCount,
					syncTargets(bookmark.Stack),
					bookmark.Current,
				)),
			)
			for _, rev := range bookmark.Revisions {
				lines = append(lines, renderRevisionAgent(rev))
			}
		}
	}
	return strings.Join(lines, "\n")
}

func renderBookmarkAgent(bookmark vcs.BookmarkSnapshot) []string {
	lines := []string{
		muted(fmt.Sprintf(`bookmark name=%q alias=%s base=%s changes=%d approved=%d/%d sync=%s current=%t`,
			bookmark.Stack.Name,
			bookmark.Stack.Alias,
			firstNonEmpty(bookmark.Stack.BaseRef, "main"),
			bookmark.ChangeCount,
			bookmark.ApprovedCount,
			bookmark.ChangeCount,
			syncTargets(bookmark.Stack),
			bookmark.Current,
		)),
	}
	limit := len(bookmark.Revisions)
	if limit > 5 {
		limit = 5
	}
	for _, rev := range bookmark.Revisions[:limit] {
		lines = append(lines, renderRevisionAgent(rev))
	}
	return lines
}

func renderRevisionAgent(rev vcs.RevisionSnapshot) string {
	extra := ""
	if rev.Working {
		extra = " working=true"
	}
	return muted(fmt.Sprintf(`revision id=%s message=%q sync=%s%s`, rev.ShortID, rev.Description, rev.SyncNote, extra))
}

func divider() string {
	return muted(strings.Repeat("─", 40))
}

func currentBookmark(snapshot vcs.StatusSnapshot) vcs.BookmarkSnapshot {
	for _, bookmark := range snapshot.Bookmarks {
		if bookmark.Current {
			return bookmark
		}
	}
	if len(snapshot.Bookmarks) > 0 {
		return snapshot.Bookmarks[0]
	}
	return vcs.BookmarkSnapshot{}
}

func currentBookmarkName(snapshot vcs.StatusSnapshot) string {
	return currentBookmark(snapshot).Stack.Name
}

func findBookmark(snapshot vcs.StatusSnapshot, filter string) *vcs.BookmarkSnapshot {
	for index, bookmark := range snapshot.Bookmarks {
		if bookmarkMatchesFilter(bookmark, filter) {
			copy := bookmark
			if copy.Stack.Alias == "" {
				copy.Stack.Alias = fmt.Sprintf("b%d", index+1)
			}
			return &copy
		}
	}
	return nil
}

func bookmarkMatchesFilter(bookmark vcs.BookmarkSnapshot, filter string) bool {
	filter = strings.TrimSpace(strings.ToLower(filter))
	body := bookmark.Stack
	return strings.EqualFold(body.Name, filter) ||
		strings.EqualFold(body.Alias, filter) ||
		strings.EqualFold(body.BookmarkName, filter) ||
		strings.Contains(strings.ToLower(body.Name), filter)
}

func countChanges(snapshot vcs.StatusSnapshot, bookmarkName string) int {
	for _, bookmark := range snapshot.Bookmarks {
		if bookmark.Stack.BookmarkName == bookmarkName {
			return bookmark.ChangeCount
		}
	}
	return 0
}

func countApproved(snapshot vcs.StatusSnapshot, bookmarkName string) int {
	for _, bookmark := range snapshot.Bookmarks {
		if bookmark.Stack.BookmarkName == bookmarkName {
			return bookmark.ApprovedCount
		}
	}
	return 0
}

func fmtChanges(n int) string {
	if n == 1 {
		return "1 change"
	}
	return fmt.Sprintf("%d changes", n)
}

func fmtApproved(approved, total int) string {
	if total == 0 {
		return "0/0 approved"
	}
	return fmt.Sprintf("%d/%d approved", approved, total)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func shortID(value string, max int) string {
	if len(value) <= max {
		return value
	}
	return value[:max]
}
