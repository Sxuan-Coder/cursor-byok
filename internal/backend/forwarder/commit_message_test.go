package forwarder

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	modeladapter "cursor/internal/backend/agent/model"
	"cursor/gen/aiserverv1"
	legacyruntime "cursor/internal/runtime"
	promptassets "cursor/prompt"
)

type fakeCommitResolver struct {
	channels     map[string]string
	commitHash   string
	commitPrompt string
	lastHash     string
}

func (fake *fakeCommitResolver) SelectChannelForModel(_ context.Context, modelID string) (*legacyruntime.ResolvedChannel, error) {
	channelID, ok := fake.channels[modelID]
	if !ok {
		return nil, errors.New("channel not available")
	}
	return &legacyruntime.ResolvedChannel{ID: channelID}, nil
}

func (fake *fakeCommitResolver) ProviderStreamIdleTimeout(context.Context) time.Duration {
	return time.Minute
}

func (fake *fakeCommitResolver) CommitModelHash() string {
	return fake.commitHash
}

func (fake *fakeCommitResolver) CommitPrompt() string {
	return fake.commitPrompt
}

func (fake *fakeCommitResolver) LastAgentModelHash() string {
	return fake.lastHash
}

func (fake *fakeCommitResolver) SaveLastAgentModelHash(context.Context, string) error {
	return nil
}

func newCommitTestService(fake *fakeCommitResolver) *Service {
	return &Service{
		resolver:     fake,
		commitConfig: fake,
		modelMemory:  fake,
	}
}

func TestBuildCommitMessagePromptUsesOverride(t *testing.T) {
	request := &aiserverv1.WriteGitCommitMessageRequest{}
	messages, err := buildCommitMessagePrompt(request, []string{"diff --git a/x"}, "  自定义提交提示词  ")
	if err != nil {
		t.Fatalf("buildCommitMessagePrompt returned error: %v", err)
	}
	if len(messages) == 0 {
		t.Fatal("buildCommitMessagePrompt returned no messages")
	}
	if messages[0].Content != "自定义提交提示词" {
		t.Fatalf("system prompt = %q, want %q", messages[0].Content, "自定义提交提示词")
	}
}

func TestBuildCommitMessagePromptFallsBackToEmbeddedAsset(t *testing.T) {
	embedded, err := promptassets.ReadCommitPrompt()
	if err != nil {
		t.Fatalf("ReadCommitPrompt returned error: %v", err)
	}
	request := &aiserverv1.WriteGitCommitMessageRequest{}
	messages, err := buildCommitMessagePrompt(request, []string{"diff --git a/x"}, "  ")
	if err != nil {
		t.Fatalf("buildCommitMessagePrompt returned error: %v", err)
	}
	if messages[0].Content != strings.TrimSpace(embedded) {
		t.Fatalf("system prompt should fall back to embedded asset, got %q", messages[0].Content)
	}
}

func TestResolveCommitMessageModelIDPrefersConfiguredHash(t *testing.T) {
	fake := &fakeCommitResolver{
		channels:   map[string]string{"hash-commit": "hash-commit", "hash-last": "hash-last"},
		commitHash: "hash-commit",
		lastHash:   "hash-last",
	}
	service := newCommitTestService(fake)
	modelID, source, hash, err := service.resolveCommitMessageModelID(context.Background())
	if err != nil {
		t.Fatalf("resolveCommitMessageModelID returned error: %v", err)
	}
	if modelID != "hash-commit" || source != "commit_model_hash" || hash != "hash-commit" {
		t.Fatalf("got modelID=%q source=%q hash=%q, want configured commit hash", modelID, source, hash)
	}
}

func TestResolveCommitMessageModelIDRejectsUnconfiguredModel(t *testing.T) {
	fake := &fakeCommitResolver{
		channels:   map[string]string{},
		commitHash: "missing-hash",
	}
	service := newCommitTestService(fake)
	_, _, _, err := service.resolveCommitMessageModelID(context.Background())
	if err == nil {
		t.Fatal("resolveCommitMessageModelID should reject a commit model that is not configured")
	}
}

func TestResolveCommitMessageModelIDFallsBackToLastAgentModel(t *testing.T) {
	fake := &fakeCommitResolver{
		channels: map[string]string{"hash-last": "hash-last"},
		lastHash: "hash-last",
	}
	service := newCommitTestService(fake)
	modelID, source, hash, err := service.resolveCommitMessageModelID(context.Background())
	if err != nil {
		t.Fatalf("resolveCommitMessageModelID returned error: %v", err)
	}
	if modelID != "hash-last" || source != "last_agent_model_hash" || hash != "hash-last" {
		t.Fatalf("got modelID=%q source=%q hash=%q, want last agent model hash", modelID, source, hash)
	}
}

func TestResolveCommitMessageModelIDDefaultFallback(t *testing.T) {
	fake := &fakeCommitResolver{channels: map[string]string{}}
	service := newCommitTestService(fake)
	modelID, source, hash, err := service.resolveCommitMessageModelID(context.Background())
	if err != nil {
		t.Fatalf("resolveCommitMessageModelID returned error: %v", err)
	}
	if modelID != "" || source != "default_fallback" || hash != "" {
		t.Fatalf("got modelID=%q source=%q hash=%q, want default fallback", modelID, source, hash)
	}
}

func TestCommitPromptOverrideReadsConfig(t *testing.T) {
	fake := &fakeCommitResolver{commitPrompt: "  覆盖提示词  "}
	service := newCommitTestService(fake)
	if got := service.commitPromptOverride(); got != "覆盖提示词" {
		t.Fatalf("commitPromptOverride = %q, want %q", got, "覆盖提示词")
	}
}

var _ modeladapter.ChannelResolver = (*fakeCommitResolver)(nil)
var _ commitConfigProvider = (*fakeCommitResolver)(nil)
var _ agentModelMemory = (*fakeCommitResolver)(nil)