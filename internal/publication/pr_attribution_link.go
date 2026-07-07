package publication

import (
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/satoricorp/gx/internal/codereview"
	"github.com/satoricorp/gx/internal/reviewbundle"
)

var (
	prNumberRefRE    = regexp.MustCompile(`(?i)\bPR\s*#(\d+)\b`)
	prHashNumberRE   = regexp.MustCompile(`#(\d+)\b`)
	prOwnerRepoNumRE = regexp.MustCompile(`(?i)([A-Za-z0-9_.-]+)/([A-Za-z0-9_.-]+)#(\d+)\b`)
	fileLineRefRE    = regexp.MustCompile(`^(.+?):(\d+)$`)
)

type attributionLinkContext struct {
	prURL    string
	owner    string
	repo     string
	sha      string
	catalog  prBodyCatalog
	snippets []codereview.ContextSnippet
}

func newAttributionLinkContext(artifact reviewbundle.Artifact, catalog prBodyCatalog, snippets []codereview.ContextSnippet) attributionLinkContext {
	prURL := pullRequestURL(artifact)
	owner, repo, _ := parseGitHubOwnerRepo(prURL)
	return attributionLinkContext{
		prURL:    prURL,
		owner:    owner,
		repo:     repo,
		sha:      headCommitSHA(artifact, catalog),
		catalog:  catalog,
		snippets: snippets,
	}
}

func parseFileLineRef(ref string) (file string, line int) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return "", 0
	}
	if match := fileLineRefRE.FindStringSubmatch(ref); len(match) == 3 {
		line, err := strconv.Atoi(match[2])
		if err != nil || line <= 0 {
			return match[1], 0
		}
		return match[1], line
	}
	return ref, 0
}

func looksLikeFilePath(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}
	if strings.Contains(value, "/") {
		return true
	}
	switch strings.ToLower(path.Ext(value)) {
	case ".go", ".md", ".json", ".yaml", ".yml", ".ts", ".tsx", ".js", ".jsx", ".py", ".rs", ".toml", ".sql", ".sh":
		return true
	default:
		return false
	}
}

func (ctx attributionLinkContext) pullRequestLink(owner, repo string, number int) string {
	if number <= 0 {
		return ""
	}
	if owner == "" {
		owner = ctx.owner
	}
	if repo == "" {
		repo = ctx.repo
	}
	if owner == "" || repo == "" {
		return ""
	}
	return fmt.Sprintf("https://github.com/%s/%s/pull/%d", owner, repo, number)
}

func (ctx attributionLinkContext) resolvePRRefURL(ref string) string {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return ""
	}
	if match := prOwnerRepoNumRE.FindStringSubmatch(ref); len(match) == 4 {
		number, err := strconv.Atoi(match[3])
		if err == nil {
			return ctx.pullRequestLink(match[1], match[2], number)
		}
	}
	if match := prNumberRefRE.FindStringSubmatch(ref); len(match) == 2 {
		number, err := strconv.Atoi(match[1])
		if err == nil {
			return ctx.pullRequestLink("", "", number)
		}
	}
	if strings.HasPrefix(ref, "#") {
		if match := prHashNumberRE.FindStringSubmatch(ref); len(match) == 2 {
			number, err := strconv.Atoi(match[1])
			if err == nil {
				return ctx.pullRequestLink("", "", number)
			}
		}
	}
	return ""
}

func (ctx attributionLinkContext) resolveFileRefURL(file string, line int) string {
	file = strings.TrimSpace(file)
	if file == "" {
		return ""
	}
	if fileInPRDiff(ctx.catalog, file) {
		if line > 0 && validatedHunkAnchor(ctx.catalog.Hunks, codereview.Finding{File: file, Line: line}) {
			return githubHunkLineLink(ctx.prURL, file, line)
		}
		for _, hunk := range ctx.catalog.Hunks {
			if hunk.File != file {
				continue
			}
			if hunk.Link != "" {
				return hunk.Link
			}
			return githubHunkLineLink(ctx.prURL, file, hunk.ChangedLine)
		}
	}
	return blobPermalink(ctx.prURL, ctx.sha, file, line)
}

func (ctx attributionLinkContext) snippetURLForLabel(label string) string {
	label = strings.TrimSpace(label)
	if label == "" {
		return ""
	}
	for _, snippet := range ctx.snippets {
		url := strings.TrimSpace(snippet.URL)
		if url == "" {
			continue
		}
		if snippet.Ref == label || snippet.SourceLabel == label || snippet.Publisher == label || snippet.Title == label {
			return url
		}
	}
	return ""
}

func (ctx attributionLinkContext) resolveRefURL(ref string) string {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return ""
	}
	if url := ctx.resolvePRRefURL(ref); url != "" {
		return url
	}
	if url := ctx.snippetURLForLabel(ref); url != "" {
		return url
	}
	file, line := parseFileLineRef(ref)
	if !looksLikeFilePath(file) {
		return ""
	}
	return ctx.resolveFileRefURL(file, line)
}

func (ctx attributionLinkContext) linkifyAttributionText(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}
	if url := ctx.resolveRefURL(text); url != "" {
		return fmt.Sprintf("[%s](%s)", text, url)
	}
	replacers := []struct {
		re *regexp.Regexp
		fn func(string) string
	}{
		{prOwnerRepoNumRE, func(match string) string {
			submatch := prOwnerRepoNumRE.FindStringSubmatch(match)
			if len(submatch) != 4 {
				return match
			}
			number, err := strconv.Atoi(submatch[3])
			if err != nil {
				return match
			}
			if url := ctx.pullRequestLink(submatch[1], submatch[2], number); url != "" {
				return fmt.Sprintf("[%s](%s)", match, url)
			}
			return match
		}},
		{prNumberRefRE, func(match string) string {
			submatch := prNumberRefRE.FindStringSubmatch(match)
			if len(submatch) != 2 {
				return match
			}
			number, err := strconv.Atoi(submatch[1])
			if err != nil {
				return match
			}
			if url := ctx.pullRequestLink("", "", number); url != "" {
				return fmt.Sprintf("[%s](%s)", match, url)
			}
			return match
		}},
		{prHashNumberRE, func(match string) string {
			submatch := prHashNumberRE.FindStringSubmatch(match)
			if len(submatch) != 2 {
				return match
			}
			number, err := strconv.Atoi(submatch[1])
			if err != nil {
				return match
			}
			if url := ctx.pullRequestLink("", "", number); url != "" {
				return fmt.Sprintf("[%s](%s)", match, url)
			}
			return match
		}},
	}
	for _, replacer := range replacers {
		text = replacer.re.ReplaceAllStringFunc(text, replacer.fn)
	}
	return text
}

func renderAttributionPart(attr prAttribution, ctx attributionLinkContext) string {
	kind := firstNonEmpty(attr.Kind, "source")
	label := firstNonEmpty(attr.Label, attr.Kind)
	ref := strings.TrimSpace(attr.Ref)
	url := strings.TrimSpace(attr.URL)

	if attr.Opaque {
		if url == "" {
			url = ctx.snippetURLForLabel(label)
		}
		if url != "" {
			return fmt.Sprintf("%s [%s](%s)", kind, label, url)
		}
		return ctx.linkifyAttributionText(label)
	}

	if url == "" && ref != "" {
		url = ctx.resolveRefURL(ref)
	}
	if url == "" && label != "" && label != ref {
		url = ctx.resolveRefURL(label)
	}

	display := label
	if ref != "" {
		display = ref
	}
	if url != "" && display != "" {
		return fmt.Sprintf("%s [%s](%s)", kind, display, url)
	}
	if ref != "" {
		if linked := ctx.linkifyAttributionText(ref); linked != ref {
			return fmt.Sprintf("%s %s", kind, linked)
		}
		return fmt.Sprintf("%s `%s`", kind, ref)
	}
	if linked := ctx.linkifyAttributionText(label); linked != label {
		return linked
	}
	return label
}

func linkProvenanceSource(source string, snippets []codereview.ContextSnippet) string {
	source = strings.TrimSpace(source)
	if source == "" {
		return ""
	}
	ctx := attributionLinkContext{snippets: snippets}
	if url := ctx.snippetURLForLabel(source); url != "" {
		return fmt.Sprintf("[%s](%s)", source, url)
	}
	return source
}
