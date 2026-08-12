package codereview

import (
	"context"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/satoricorp/gx/internal/semantic"
)

func liveReviewResourceRetriever(t *testing.T) ReviewResourceRetriever {
	t.Helper()
	namespace := firstNonEmpty(strings.TrimSpace(os.Getenv("GX_REVIEW_KNOWLEDGE_NAMESPACE")), defaultReviewKnowledgeNamespace)
	return ReviewResourceRetriever{
		Embedder: semantic.NewOpenAIEmbedder(semantic.Config{
			OpenAIAPIKey:         strings.TrimSpace(os.Getenv("OPENAI_API_KEY")),
			OpenAIBaseURL:        normalizeReviewOpenAIBaseURL(firstNonEmpty(os.Getenv("GX_OPENAI_BASE_URL"), os.Getenv("OPENAI_BASE_URL"), "https://api.openai.com")),
			OpenAIEmbeddingModel: firstNonEmpty(os.Getenv("GX_OPENAI_EMBEDDING_MODEL"), "text-embedding-3-small"),
			EmbeddingDimensions:  reviewEnvInt("GX_EMBEDDING_DIMENSIONS", 512),
		}),
		Store: turboPufferReviewResourceStore{
			apiKey:     strings.TrimSpace(os.Getenv("TURBOPUFFER_API_KEY")),
			baseURL:    strings.TrimRight(firstNonEmpty(os.Getenv("GX_TPUF_BASE_URL"), defaultReviewResourceBaseURL), "/"),
			namespace:  namespace,
			httpClient: &http.Client{Timeout: 30 * time.Second},
		},
		Namespace: namespace,
		Limit:     defaultReviewResourceTopK,
	}
}

// TestLiveReviewKnowledgeAnswersEveryChangeShape runs the real corpus query over
// the change shapes a review actually sees, because the question that made this
// worth testing was not "does the query work" but "does it work often enough to
// be worth growing the corpus". A corpus consulted on a third of reviews cannot
// be A/B tested: the arm that queries it less wins or loses on retrieval
// variance rather than on what the corpus contains.
func TestLiveReviewKnowledgeAnswersEveryChangeShape(t *testing.T) {
	requireLiveRetrievalCredentials(t)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	retriever := liveReviewResourceRetriever(t)

	cases := []struct {
		name  string
		opts  Options
		files []string
	}{
		{name: "go security", opts: Options{Scope: "security"}, files: []string{"internal/auth/session.go", "internal/api/handler.go"}},
		{name: "ruby web", opts: Options{}, files: []string{"app/controllers/users_controller.rb", "app/models/user.rb"}},
		{name: "java service", opts: Options{}, files: []string{"src/main/java/org/keycloak/authentication/Authenticator.java"}},
		{name: "python service", opts: Options{}, files: []string{"src/sentry/api/endpoints/organization_details.py"}},
		{name: "typescript frontend", opts: Options{}, files: []string{"web/src/components/Login.tsx"}},
		{name: "sql migration", opts: Options{}, files: []string{"db/migrations/2026_add_users.sql"}},
		{name: "ci workflow", opts: Options{}, files: []string{".github/workflows/release.yml"}},
		{name: "dependency bump", opts: Options{Scope: "dependencies"}, files: []string{"package.json", "package-lock.json"}},
		{name: "deep whole repo", opts: Options{Deep: true, WholeRepo: true}, files: []string{"internal/codereview/engine.go"}},
		{name: "unmapped language", opts: Options{}, files: []string{"src/main.erl", "config/sys.config"}},
		{name: "single doc file", opts: Options{}, files: []string{"README.md"}},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			log := &EvidenceLog{}
			snippets, err := retriever.Retrieve(ctx, RetrieveInput{
				RepoRoot:     "/repo",
				Options:      normalizeOptions(tt.opts),
				Facts:        RepoFacts{Files: tt.files},
				ChangedFiles: tt.files,
				Evidence:     log,
			})
			if err != nil {
				t.Fatalf("Retrieve() error = %v", err)
			}
			statuses := log.Statuses()
			if len(statuses) != 1 {
				t.Fatalf("statuses = %#v, want exactly one review-knowledge line", statuses)
			}
			if statuses[0].State != EvidenceOK {
				t.Fatalf("state = %q (%s), want ok — the corpus answered nothing for this change shape",
					statuses[0].State, statuses[0].Detail)
			}
			if len(snippets) == 0 {
				t.Fatalf("no snippets for %v", tt.files)
			}
			t.Logf("%-20s %d snippet(s)", tt.name, len(snippets))
		})
	}
}
