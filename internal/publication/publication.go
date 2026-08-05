package publication

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/satoricorp/lgtm/internal/cloud"
	"github.com/satoricorp/lgtm/internal/reviewbundle"
	"github.com/satoricorp/lgtm/internal/semantic"
	"github.com/satoricorp/lgtm/internal/storage"
	"github.com/satoricorp/lgtm/internal/vcs"
)

type Uploader interface {
	UploadReviewArtifact(ctx context.Context, artifact reviewbundle.Artifact) (reviewbundle.Artifact, error)
}

type SemanticIndexer interface {
	IndexBundle(ctx context.Context, bundle reviewbundle.Bundle) (semantic.IndexResult, error)
}

type Result struct {
	Uploaded              bool
	Queued                bool
	QueueID               string
	QueuePath             string
	ReviewID              string
	ReviewURL             string
	IndexStatus           string
	SemanticIndexed       bool
	SemanticChunks        int
	SemanticSessionChunks int
	SemanticCodeChunks    int
	SemanticIndexError    string
	ArtifactPath          string
	Artifact              reviewbundle.Artifact
}

type StackPublisher interface {
	PreparePublish(ctx context.Context, args []string, opts vcs.PushOptions) (vcs.PushResult, error)
	RecordPublish(ctx context.Context, result vcs.PushResult) error
}

type AllStackPublisher interface {
	PrepareAllPublishes(ctx context.Context, args []string, opts vcs.PushOptions) ([]vcs.PushResult, error)
	RecordPublish(ctx context.Context, result vcs.PushResult) error
}

type StackResult struct {
	Push   vcs.PushResult
	Review Result
}

type AllStacksResult struct {
	Stacks []StackResult
}

type Publisher struct {
	uploader Uploader
	indexer  SemanticIndexer
}

type ReviewPublication struct {
	publisher *Publisher
}

func NewPublisher(uploader Uploader) *Publisher {
	indexer, _, err := semantic.NewLocalIndexerFromEnv()
	publisher := &Publisher{uploader: uploader, indexer: indexer}
	if err != nil {
		publisher.indexer = nil
	}
	return publisher
}

func NewPublisherWithIndexer(uploader Uploader, indexer SemanticIndexer) *Publisher {
	return &Publisher{uploader: uploader, indexer: indexer}
}

func NewReviewPublication(publisher *Publisher) *ReviewPublication {
	return &ReviewPublication{publisher: publisher}
}

func PublishStack(ctx context.Context, stack StackPublisher, args []string, opts vcs.PushOptions, publisher *Publisher) (StackResult, error) {
	return NewReviewPublication(publisher).PublishStack(ctx, stack, args, opts)
}

func (p *ReviewPublication) PublishStack(ctx context.Context, stack StackPublisher, args []string, opts vcs.PushOptions) (StackResult, error) {
	if stack == nil {
		return StackResult{}, fmt.Errorf("publish stack engine is required")
	}
	if err := p.configure(); err != nil {
		return StackResult{}, err
	}
	push, err := stack.PreparePublish(ctx, args, opts)
	if err != nil {
		return StackResult{}, err
	}
	return p.PublishPrepared(ctx, push, stack.RecordPublish)
}

func PublishAllStacks(ctx context.Context, stacks AllStackPublisher, args []string, opts vcs.PushOptions, publisher *Publisher) (AllStacksResult, error) {
	return NewReviewPublication(publisher).PublishAllStacks(ctx, stacks, args, opts)
}

func (p *ReviewPublication) PublishAllStacks(ctx context.Context, stacks AllStackPublisher, args []string, opts vcs.PushOptions) (AllStacksResult, error) {
	if stacks == nil {
		return AllStacksResult{}, fmt.Errorf("publish stack engine is required")
	}
	if err := p.configure(); err != nil {
		return AllStacksResult{}, err
	}

	prepared, prepareErr := stacks.PrepareAllPublishes(ctx, args, opts)
	if len(prepared) == 0 && prepareErr != nil {
		return AllStacksResult{}, prepareErr
	}
	results := AllStacksResult{Stacks: make([]StackResult, 0, len(prepared))}
	for _, push := range prepared {
		result, err := p.PublishPrepared(ctx, push, stacks.RecordPublish)
		results.Stacks = append(results.Stacks, result)
		if err != nil {
			return results, err
		}
	}
	if prepareErr != nil {
		return results, prepareErr
	}
	return results, nil
}

func (p *ReviewPublication) PublishPrepared(ctx context.Context, push vcs.PushResult, record func(context.Context, vcs.PushResult) error) (StackResult, error) {
	if err := p.configure(); err != nil {
		return StackResult{Push: push}, err
	}
	review, err := p.publisher.PublishPush(ctx, push)
	if err != nil {
		return StackResult{Push: push}, err
	}
	result := StackResult{Push: push, Review: review}
	if record != nil {
		if err := record(ctx, push); err != nil {
			return result, err
		}
	}
	return result, nil
}

func (p *ReviewPublication) configure() error {
	publisher, err := configuredPublisher(p.publisher)
	if err != nil {
		return err
	}
	p.publisher = publisher
	return nil
}

func configuredPublisher(publisher *Publisher) (*Publisher, error) {
	if publisher != nil {
		return publisher, nil
	}
	client := cloud.NewClient()
	if client == nil {
		return nil, fmt.Errorf("lgtm cloud is not configured; set LGTM_CLOUD_URL or rebuild with cloud endpoints")
	}
	return NewPublisher(client), nil
}

func (p *Publisher) PublishPush(ctx context.Context, push vcs.PushResult) (Result, error) {
	if p == nil || p.uploader == nil {
		return Result{}, nil
	}
	bundle, err := reviewbundle.BuildPush(ctx, push)
	if err != nil {
		return Result{}, err
	}
	return p.PublishArtifact(ctx, reviewbundle.NewArtifact(bundle))
}

func (p *Publisher) PublishBundle(ctx context.Context, bundle reviewbundle.Bundle) (Result, error) {
	return p.PublishArtifact(ctx, reviewbundle.NewArtifact(bundle))
}

func (p *Publisher) PublishArtifact(ctx context.Context, artifact reviewbundle.Artifact) (Result, error) {
	result := Result{}
	if p == nil || p.uploader == nil {
		return result, nil
	}
	if _, err := UpdateGitHubPullRequestBody(ctx, artifact); err != nil {
		return Result{}, err
	}
	if p.indexer != nil {
		indexResult, err := p.indexer.IndexBundle(ctx, artifact.Bundle)
		result.SemanticIndexed = indexResult.Indexed
		result.SemanticChunks = indexResult.Chunks
		result.SemanticSessionChunks = indexResult.SessionChunks
		result.SemanticCodeChunks = indexResult.CodeChunks
		if err != nil {
			result.SemanticIndexError = err.Error()
		}
	} else if _, _, err := semantic.NewLocalIndexerFromEnv(); err != nil {
		result.SemanticIndexError = err.Error()
	}
	artifact.IndexStatus = firstNonEmpty(artifact.IndexStatus, "pending")
	upload, err := p.uploader.UploadReviewArtifact(ctx, artifact)
	if err != nil {
		return Result{}, err
	}
	result.Uploaded = true
	result.ReviewID = upload.ReviewID
	result.ReviewURL = upload.ReviewURL
	result.IndexStatus = upload.IndexStatus
	result.Artifact = upload
	artifactPath, err := WriteLocalArtifact(upload)
	if err != nil {
		return result, err
	}
	result.ArtifactPath = artifactPath
	return result, nil
}

func WriteLocalArtifact(artifact reviewbundle.Artifact) (string, error) {
	dir, err := storage.DefaultDir()
	if err != nil {
		return "", err
	}
	reviewsDir := filepath.Join(dir, "reviews")
	if err := os.MkdirAll(reviewsDir, 0o755); err != nil {
		return "", fmt.Errorf("create local review artifact dir: %w", err)
	}
	filename := localArtifactFilename(artifact)
	data, err := json.MarshalIndent(artifact, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal local review artifact: %w", err)
	}
	path := filepath.Join(reviewsDir, filename)
	if err := os.WriteFile(path, append(data, '\n'), 0o600); err != nil {
		return "", fmt.Errorf("write local review artifact: %w", err)
	}
	return path, nil
}

func localArtifactFilename(artifact reviewbundle.Artifact) string {
	name := firstNonEmpty(artifact.ReviewID, artifact.Push.HeadCommitID, fmt.Sprintf("%d", artifact.CreatedAt), "review")
	return sanitizeArtifactName(name) + ".json"
}

var artifactNamePattern = regexp.MustCompile(`[^A-Za-z0-9._-]+`)

func sanitizeArtifactName(value string) string {
	value = artifactNamePattern.ReplaceAllString(strings.TrimSpace(value), "-")
	value = strings.Trim(value, ".-")
	if value == "" {
		return "review"
	}
	if len(value) > 120 {
		return value[:120]
	}
	return value
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
