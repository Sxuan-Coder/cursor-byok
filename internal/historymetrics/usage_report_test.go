package historymetrics

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadUsageReportMissingFile(t *testing.T) {
	report, err := LoadUsageReport(filepath.Join(t.TempDir(), "usage.json"), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Daily) != 0 || len(report.Models) != 0 || len(report.Events) != 0 {
		t.Fatalf("expected empty report, got %+v", report)
	}
}

func TestLoadUsageReportWithModelsAndEvents(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "usage.json")
	body := `{
  "schema_version": 3,
  "totals": {"provider_calls": 2, "turns_total": 1, "valid_turns_total": 1, "input_tokens": 100, "output_tokens": 200, "cache_read_tokens": 50, "cache_write_tokens": 10, "total_tokens": 360},
  "daily": [{"date": "2026-08-25", "provider_calls": 2, "input_tokens": 100, "output_tokens": 200, "cache_read_tokens": 50, "cache_write_tokens": 10, "total_tokens": 360}],
  "models": {
    "openai::gpt-x": {"provider": "openai", "model": "gpt-x", "provider_calls": 2, "input_tokens": 100, "output_tokens": 200, "cache_read_tokens": 50, "cache_write_tokens": 10, "total_tokens": 360}
  },
  "recent_events": [
    {"event_id": "req-1", "kind": "provider_call", "status": "completed", "request_id": "req-1", "provider": "openai", "model": "gpt-x", "duration_ms": 1200, "input_tokens": 100, "output_tokens": 200, "cache_read_tokens": 50, "cache_write_tokens": 10, "total_tokens": 360, "usage_present": true}
  ]
}`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	report, err := LoadUsageReport(path, PricingTable{
		"openai::gpt-x": {Input: 3, Output: 15, CacheRead: 0.3, CacheWrite: 3.75},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.Summary.ProviderCallsTotal != 2 {
		t.Fatalf("provider calls = %d, want 2", report.Summary.ProviderCallsTotal)
	}
	if len(report.Models) != 1 || report.Models[0].Model != "gpt-x" {
		t.Fatalf("models = %+v", report.Models)
	}
	if !report.Models[0].PricingConfigured {
		t.Fatal("expected pricing configured")
	}
	// (100*3 + 200*15 + 50*0.3 + 10*3.75) / 1e6
	want := (100*3 + 200*15 + 50*0.3 + 10*3.75) / 1_000_000
	if diff := report.Models[0].EstimatedCostUSD - want; diff < -1e-12 || diff > 1e-12 {
		t.Fatalf("cost = %v, want %v", report.Models[0].EstimatedCostUSD, want)
	}
	if len(report.Events) != 1 || report.Events[0].DurationMS != 1200 || report.Events[0].Model != "gpt-x" {
		t.Fatalf("events = %+v", report.Events)
	}
	if report.Events[0].CacheHitRate == nil || *report.Events[0].CacheHitRate != 50.0/150.0 {
		t.Fatalf("cache hit rate = %v", report.Events[0].CacheHitRate)
	}
	if len(report.Daily) != 1 || report.Daily[0].Date != "2026-08-25" {
		t.Fatalf("daily = %+v", report.Daily)
	}
}
