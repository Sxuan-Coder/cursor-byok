package forwarder

import (
	"testing"
	"time"
)

func TestUsageStoreModelAggregation(t *testing.T) {
	store := NewUsageFileStore(t.TempDir())
	at := time.Date(2026, 8, 25, 10, 0, 0, 0, time.UTC)

	if err := store.UpsertEvent(usageFileEvent{
		EventID:      "req-1",
		Kind:         usageEventKindProvider,
		At:           at,
		RequestID:    "req-1",
		Provider:     "openai",
		Model:        "gpt-x",
		DurationMS:   900,
		InputTokens:  100,
		OutputTokens: 40,
	}); err != nil {
		t.Fatal(err)
	}
	if err := store.UpsertEvent(usageFileEvent{
		EventID:      "req-2",
		Kind:         usageEventKindProvider,
		At:           at.Add(time.Minute),
		RequestID:    "req-2",
		Provider:     "openai",
		Model:        "gpt-y",
		InputTokens:  10,
		OutputTokens: 5,
	}); err != nil {
		t.Fatal(err)
	}

	doc, err := readUsageFileDocument(store.path)
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Models) != 2 {
		t.Fatalf("models = %d, want 2: %+v", len(doc.Models), doc.Models)
	}
	bucket := doc.Models["openai::gpt-x"]
	if bucket.Provider != "openai" || bucket.Model != "gpt-x" {
		t.Fatalf("bucket identity = %+v", bucket)
	}
	if bucket.InputTokens != 100 || bucket.OutputTokens != 40 || bucket.ProviderCalls != 1 {
		t.Fatalf("bucket totals = %+v", bucket)
	}

	// 覆盖同一事件：模型聚合需要先减旧值再加新值。
	if err := store.UpsertEvent(usageFileEvent{
		EventID:      "req-1",
		Kind:         usageEventKindProvider,
		At:           at,
		RequestID:    "req-1",
		Provider:     "openai",
		Model:        "gpt-x",
		InputTokens:  60,
		OutputTokens: 20,
	}); err != nil {
		t.Fatal(err)
	}
	doc, err = readUsageFileDocument(store.path)
	if err != nil {
		t.Fatal(err)
	}
	bucket = doc.Models["openai::gpt-x"]
	if bucket.InputTokens != 60 || bucket.OutputTokens != 20 {
		t.Fatalf("after overwrite bucket = %+v", bucket)
	}
	if doc.Totals.InputTokens != 70 || doc.Totals.OutputTokens != 25 {
		t.Fatalf("totals after overwrite = %+v", doc.Totals)
	}

	// turn 事件不进入模型聚合。
	if err := store.UpsertEvent(usageFileEvent{
		EventID:     "turn::c1::1",
		Kind:        usageEventKindTurn,
		Status:      usageTurnStatusDone,
		At:          at,
		InputTokens: 999,
	}); err != nil {
		t.Fatal(err)
	}
	doc, err = readUsageFileDocument(store.path)
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Models) != 2 {
		t.Fatalf("turn event should not create model bucket, models = %+v", doc.Models)
	}
	if _, ok := doc.Models[""]; ok {
		t.Fatal("empty model key should not be stored")
	}
}
