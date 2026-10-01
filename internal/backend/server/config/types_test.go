package config

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func testModelAdapter(displayName string, sortValue int) ModelAdapterConfig {
	return ModelAdapterConfig{
		Sort:            sortValue,
		DisplayName:     displayName,
		Type:            "openai",
		BaseURL:         "https://api.example.com/v1",
		APIKey:          "test-key",
		TooltipData:     displayName,
		ModelID:         displayName,
		ReasoningEffort: "medium",
		OpenAIEndpoint:  "/v1/responses",
	}
}

func TestNormalizeModelAdapterConfigsPreservesLegacyArrayOrder(t *testing.T) {
	adapters, err := NormalizeModelAdapterConfigs([]ModelAdapterConfig{
		testModelAdapter("first", 0),
		testModelAdapter("second", 0),
		testModelAdapter("third", 0),
	})
	if err != nil {
		t.Fatalf("NormalizeModelAdapterConfigs returned error: %v", err)
	}

	for index, expectedName := range []string{"first", "second", "third"} {
		if adapters[index].DisplayName != expectedName {
			t.Fatalf("adapter %d = %q, want %q", index, adapters[index].DisplayName, expectedName)
		}
		if adapters[index].Sort != index+1 {
			t.Fatalf("adapter %d sort = %d, want %d", index, adapters[index].Sort, index+1)
		}
	}
}

func TestNormalizeModelAdapterConfigsUsesStableExplicitSort(t *testing.T) {
	adapters, err := NormalizeModelAdapterConfigs([]ModelAdapterConfig{
		testModelAdapter("legacy", 0),
		testModelAdapter("third", 30),
		testModelAdapter("first", 10),
		testModelAdapter("second-a", 20),
		testModelAdapter("second-b", 20),
	})
	if err != nil {
		t.Fatalf("NormalizeModelAdapterConfigs returned error: %v", err)
	}

	expectedNames := []string{"first", "second-a", "second-b", "third", "legacy"}
	for index, expectedName := range expectedNames {
		if adapters[index].DisplayName != expectedName {
			t.Fatalf("adapter %d = %q, want %q", index, adapters[index].DisplayName, expectedName)
		}
		if adapters[index].Sort != index+1 {
			t.Fatalf("adapter %d sort = %d, want %d", index, adapters[index].Sort, index+1)
		}
	}
}

func TestNormalizeModelAdapterConfigsAllowsBlankReasoningEffort(t *testing.T) {
	adapter := testModelAdapter("non-reasoning-model", 1)
	adapter.ReasoningEffort = ""

	adapters, err := NormalizeModelAdapterConfigs([]ModelAdapterConfig{adapter})
	if err != nil {
		t.Fatalf("NormalizeModelAdapterConfigs returned error: %v", err)
	}
	if got := adapters[0].ReasoningEffort; got != "" {
		t.Fatalf("ReasoningEffort = %q, want blank", got)
	}
}

func TestNormalizeModelAdapterConfigsRejectsUnknownReasoningEffort(t *testing.T) {
	adapter := testModelAdapter("invalid-reasoning-effort", 1)
	adapter.ReasoningEffort = "unsupported"

	if _, err := NormalizeModelAdapterConfigs([]ModelAdapterConfig{adapter}); err == nil {
		t.Fatal("NormalizeModelAdapterConfigs should reject an unknown reasoning effort")
	}
}

func TestNormalizeConfigPreservesCommitSettings(t *testing.T) {
	normalized, err := NormalizeConfig(Config{
		CommitModelHash: "  abc123  ",
		CommitPrompt:    "  自定义提交提示词  ",
	})
	if err != nil {
		t.Fatalf("NormalizeConfig returned error: %v", err)
	}
	if normalized.CommitModelHash != "abc123" {
		t.Fatalf("CommitModelHash = %q, want %q", normalized.CommitModelHash, "abc123")
	}
	if normalized.CommitPrompt != "自定义提交提示词" {
		t.Fatalf("CommitPrompt = %q, want %q", normalized.CommitPrompt, "自定义提交提示词")
	}
}

func TestNormalizeConfigKeepsCommitSettingsEmptyByDefault(t *testing.T) {
	normalized, err := NormalizeConfig(Config{})
	if err != nil {
		t.Fatalf("NormalizeConfig returned error: %v", err)
	}
	if normalized.CommitModelHash != "" || normalized.CommitPrompt != "" {
		t.Fatalf("commit settings should default to empty, got hash=%q prompt=%q", normalized.CommitModelHash, normalized.CommitPrompt)
	}
}

func TestConfigYAMLRoundTripKeepsCommitSettings(t *testing.T) {
	marshaled, err := yaml.Marshal(Config{CommitModelHash: "hash-1", CommitPrompt: "提示词"})
	if err != nil {
		t.Fatalf("yaml.Marshal returned error: %v", err)
	}
	var loaded Config
	if err := yaml.Unmarshal(marshaled, &loaded); err != nil {
		t.Fatalf("yaml.Unmarshal returned error: %v", err)
	}
	if loaded.CommitModelHash != "hash-1" || loaded.CommitPrompt != "提示词" {
		t.Fatalf("round trip lost commit settings, got hash=%q prompt=%q", loaded.CommitModelHash, loaded.CommitPrompt)
	}
}
