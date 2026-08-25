package historymetrics

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

// ModelPricing 某个模型的 token 单价（美元 / 百万 tokens）。
// 单价 <= 0 表示该项不计费。
type ModelPricing struct {
	Input      float64 `json:"input"`
	Output     float64 `json:"output"`
	CacheRead  float64 `json:"cacheRead"`
	CacheWrite float64 `json:"cacheWrite"`
}

// PricingTable 按 modelKey（provider::model 或 model）配置的价格表。
type PricingTable map[string]ModelPricing

// DefaultPricing 未命中价格表时使用的兜底单价（与首页估算口径一致）。
var DefaultPricing = ModelPricing{
	Input:      5,
	Output:     25,
	CacheRead:  0.5,
	CacheWrite: 6.25,
}

type DailyPoint struct {
	Date             string  `json:"date"`
	ProviderCalls    int64   `json:"providerCalls"`
	TurnsTotal       int64   `json:"turnsTotal"`
	InputTokens      int64   `json:"inputTokens"`
	OutputTokens     int64   `json:"outputTokens"`
	CacheReadTokens  int64   `json:"cacheReadTokens"`
	CacheWriteTokens int64   `json:"cacheWriteTokens"`
	TotalTokens      int64   `json:"totalTokens"`
	EstimatedCostUSD float64 `json:"estimatedCostUSD"`
}

type ModelUsage struct {
	Key               string    `json:"key"`
	Provider          string    `json:"provider"`
	Model             string    `json:"model"`
	ProviderCalls     int64     `json:"providerCalls"`
	InputTokens       int64     `json:"inputTokens"`
	OutputTokens      int64     `json:"outputTokens"`
	CacheReadTokens   int64     `json:"cacheReadTokens"`
	CacheWriteTokens  int64     `json:"cacheWriteTokens"`
	TotalTokens       int64     `json:"totalTokens"`
	CacheHitRate      *float64  `json:"cacheHitRate"`
	EstimatedCostUSD  float64   `json:"estimatedCostUSD"`
	PricingConfigured bool      `json:"pricingConfigured"`
	LastSeenAt        time.Time `json:"lastSeenAt,omitempty"`
}

type EventDetail struct {
	EventID          string    `json:"eventId"`
	Kind             string    `json:"kind"`
	Status           string    `json:"status"`
	At               time.Time `json:"at"`
	RequestID        string    `json:"requestId"`
	Provider         string    `json:"provider"`
	Model            string    `json:"model"`
	DurationMS       int64     `json:"durationMs"`
	ErrorText        string    `json:"errorText"`
	InputTokens      int64     `json:"inputTokens"`
	OutputTokens     int64     `json:"outputTokens"`
	CacheReadTokens  int64     `json:"cacheReadTokens"`
	CacheWriteTokens int64     `json:"cacheWriteTokens"`
	TotalTokens      int64     `json:"totalTokens"`
	CacheHitRate     *float64  `json:"cacheHitRate"`
	UsagePresent     bool      `json:"usagePresent"`
}

// Report 用量报表：摘要 + 按天序列 + 按模型聚合 + 最近请求明细。
type Report struct {
	Summary          Summary       `json:"summary"`
	Daily            []DailyPoint  `json:"daily"`
	Models           []ModelUsage  `json:"models"`
	Events           []EventDetail `json:"events"`
	EstimatedCostUSD float64       `json:"estimatedCostUSD"`
}

type usageReportDocument struct {
	Totals struct {
		ProviderCalls     int64 `json:"provider_calls"`
		TurnsTotal        int64 `json:"turns_total"`
		ValidTurnsTotal   int64 `json:"valid_turns_total"`
		InvalidTurnsTotal int64 `json:"invalid_turns_total"`
		InputTokens       int64 `json:"input_tokens"`
		OutputTokens      int64 `json:"output_tokens"`
		CacheReadTokens   int64 `json:"cache_read_tokens"`
		CacheWriteTokens  int64 `json:"cache_write_tokens"`
		TotalTokens       int64 `json:"total_tokens"`
	} `json:"totals"`
	Daily []struct {
		Date             string `json:"date"`
		ProviderCalls    int64  `json:"provider_calls"`
		TurnsTotal       int64  `json:"turns_total"`
		InputTokens      int64  `json:"input_tokens"`
		OutputTokens     int64  `json:"output_tokens"`
		CacheReadTokens  int64  `json:"cache_read_tokens"`
		CacheWriteTokens int64  `json:"cache_write_tokens"`
		TotalTokens      int64  `json:"total_tokens"`
	} `json:"daily"`
	Models map[string]struct {
		Provider         string    `json:"provider"`
		Model            string    `json:"model"`
		ProviderCalls    int64     `json:"provider_calls"`
		InputTokens      int64     `json:"input_tokens"`
		OutputTokens     int64     `json:"output_tokens"`
		CacheReadTokens  int64     `json:"cache_read_tokens"`
		CacheWriteTokens int64     `json:"cache_write_tokens"`
		TotalTokens      int64     `json:"total_tokens"`
		LastSeenAt       time.Time `json:"last_seen_at"`
	} `json:"models"`
	RecentEvents []struct {
		EventID          string    `json:"event_id"`
		Kind             string    `json:"kind"`
		Status           string    `json:"status"`
		At               time.Time `json:"at"`
		RequestID        string    `json:"request_id"`
		Provider         string    `json:"provider"`
		Model            string    `json:"model"`
		DurationMS       int64     `json:"duration_ms"`
		ErrorText        string    `json:"error_text"`
		InputTokens      int64     `json:"input_tokens"`
		OutputTokens     int64     `json:"output_tokens"`
		CacheReadTokens  int64     `json:"cache_read_tokens"`
		CacheWriteTokens int64     `json:"cache_write_tokens"`
		TotalTokens      int64     `json:"total_tokens"`
		UsagePresent     bool      `json:"usage_present"`
	} `json:"recent_events"`
}

func EstimateCostUSD(input, output, cacheRead, cacheWrite int64, pricing ModelPricing) float64 {
	cost := float64(input)/1_000_000*maxZero(pricing.Input) +
		float64(output)/1_000_000*maxZero(pricing.Output) +
		float64(cacheRead)/1_000_000*maxZero(pricing.CacheRead) +
		float64(cacheWrite)/1_000_000*maxZero(pricing.CacheWrite)
	return cost
}

func maxZero(value float64) float64 {
	if value < 0 {
		return 0
	}
	return value
}

func eventCacheHitRate(input, cacheRead int64) *float64 {
	total := input + cacheRead
	if total <= 0 {
		return nil
	}
	value := float64(cacheRead) / float64(total)
	return &value
}

// LoadUsageReport 读取 usage.json 并生成完整报表。pricing 为 nil 时全部使用默认单价。
func LoadUsageReport(path string, pricing PricingTable) (Report, error) {
	report := Report{
		Daily:  make([]DailyPoint, 0),
		Models: make([]ModelUsage, 0),
		Events: make([]EventDetail, 0),
	}
	body, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return report, nil
		}
		return report, fmt.Errorf("read usage file: %w", err)
	}
	var doc usageReportDocument
	if err := json.Unmarshal(body, &doc); err != nil {
		return report, fmt.Errorf("decode usage file: %w", err)
	}

	totals := Totals{
		InputTokens:        doc.Totals.InputTokens,
		OutputTokens:       doc.Totals.OutputTokens,
		CacheReadTokens:    doc.Totals.CacheReadTokens,
		CacheWriteTokens:   doc.Totals.CacheWriteTokens,
		PromptTokensTotal:  doc.Totals.InputTokens + doc.Totals.CacheReadTokens + doc.Totals.CacheWriteTokens,
		RequestTokensTotal: doc.Totals.TotalTokens,
	}
	report.Summary = Summary{
		ProviderCallsTotal: int(doc.Totals.ProviderCalls),
		TurnsTotal:         int(doc.Totals.TurnsTotal),
		ValidTurnsTotal:    int(doc.Totals.ValidTurnsTotal),
		InvalidTurnsTotal:  int(doc.Totals.InvalidTurnsTotal),
		RequestTokensTotal: totals.RequestTokensTotal,
		PromptTokensTotal:  totals.PromptTokensTotal,
		OutputTokensTotal:  doc.Totals.OutputTokens,
		CacheReadTokens:    totals.CacheReadTokens,
		CacheWriteTokens:   totals.CacheWriteTokens,
		CacheHitRate:       cacheHitRateFromTotals(totals),
	}
	report.EstimatedCostUSD = EstimateCostUSD(
		doc.Totals.InputTokens, doc.Totals.OutputTokens,
		doc.Totals.CacheReadTokens, doc.Totals.CacheWriteTokens,
		resolvePricing(pricing, "", ""),
	)

	for _, item := range doc.Daily {
		point := DailyPoint{
			Date:             item.Date,
			ProviderCalls:    item.ProviderCalls,
			TurnsTotal:       item.TurnsTotal,
			InputTokens:      item.InputTokens,
			OutputTokens:     item.OutputTokens,
			CacheReadTokens:  item.CacheReadTokens,
			CacheWriteTokens: item.CacheWriteTokens,
			TotalTokens:      item.TotalTokens,
		}
		point.EstimatedCostUSD = EstimateCostUSD(
			item.InputTokens, item.OutputTokens,
			item.CacheReadTokens, item.CacheWriteTokens,
			resolvePricing(pricing, "", ""),
		)
		report.Daily = append(report.Daily, point)
	}
	sort.Slice(report.Daily, func(i, j int) bool { return report.Daily[i].Date < report.Daily[j].Date })

	for key, item := range doc.Models {
		entry := ModelUsage{
			Key:              key,
			Provider:         item.Provider,
			Model:            item.Model,
			ProviderCalls:    item.ProviderCalls,
			InputTokens:      item.InputTokens,
			OutputTokens:     item.OutputTokens,
			CacheReadTokens:  item.CacheReadTokens,
			CacheWriteTokens: item.CacheWriteTokens,
			TotalTokens:      item.TotalTokens,
			CacheHitRate:     eventCacheHitRate(item.InputTokens, item.CacheReadTokens),
			LastSeenAt:       item.LastSeenAt,
		}
		pricingForModel, configured := lookupPricing(pricing, item.Provider, item.Model)
		entry.PricingConfigured = configured
		entry.EstimatedCostUSD = EstimateCostUSD(
			item.InputTokens, item.OutputTokens,
			item.CacheReadTokens, item.CacheWriteTokens, pricingForModel,
		)
		report.Models = append(report.Models, entry)
	}
	sort.Slice(report.Models, func(i, j int) bool {
		if report.Models[i].TotalTokens != report.Models[j].TotalTokens {
			return report.Models[i].TotalTokens > report.Models[j].TotalTokens
		}
		return report.Models[i].Key < report.Models[j].Key
	})

	for _, item := range doc.RecentEvents {
		// 明细只展示 provider_call；turn_finalized 是同请求的回合级汇总，
		// 会与 provider_call 重复且没有模型信息。
		if item.Kind != "provider_call" {
			continue
		}
		report.Events = append(report.Events, EventDetail{
			EventID:          item.EventID,
			Kind:             item.Kind,
			Status:           item.Status,
			At:               item.At,
			RequestID:        item.RequestID,
			Provider:         item.Provider,
			Model:            item.Model,
			DurationMS:       item.DurationMS,
			ErrorText:        item.ErrorText,
			InputTokens:      item.InputTokens,
			OutputTokens:     item.OutputTokens,
			CacheReadTokens:  item.CacheReadTokens,
			CacheWriteTokens: item.CacheWriteTokens,
			TotalTokens:      item.TotalTokens,
			CacheHitRate:     eventCacheHitRate(item.InputTokens, item.CacheReadTokens),
			UsagePresent:     item.UsagePresent,
		})
	}
	// recent_events 本身按新事件在前；保持稳定排序即可。
	return report, nil
}

func lookupPricing(pricing PricingTable, provider string, model string) (ModelPricing, bool) {
	if pricing == nil {
		return DefaultPricing, false
	}
	model = strings.TrimSpace(model)
	if model == "" {
		return DefaultPricing, false
	}
	if provider != "" {
		if value, ok := pricing[provider+"::"+model]; ok {
			return value, true
		}
	}
	if value, ok := pricing[model]; ok {
		return value, true
	}
	return DefaultPricing, false
}

func resolvePricing(pricing PricingTable, provider string, model string) ModelPricing {
	value, _ := lookupPricing(pricing, provider, model)
	return value
}
