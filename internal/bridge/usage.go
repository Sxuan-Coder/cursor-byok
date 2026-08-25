package bridge

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"cursor/internal/appdata"
	"cursor/internal/historymetrics"
)

// UsageService 定义用量报表相关的 Wails service。
type UsageService struct{}

// NewUsageService 创建用量报表 service。
func NewUsageService() *UsageService {
	return &UsageService{}
}

// GetUsageReport 返回完整用量报表（应用当前价格配置的成本估算）。
func (service *UsageService) GetUsageReport() (historymetrics.Report, error) {
	if err := appdata.EnsureAssistantHome(); err != nil {
		return historymetrics.Report{}, err
	}
	pricing, err := service.loadPricing()
	if err != nil {
		return historymetrics.Report{}, err
	}
	return historymetrics.LoadUsageReport(appdata.UsageFilePath(), pricing)
}

// GetPricing 返回当前价格配置；文件不存在时返回空表（前端用默认价展示）。
func (service *UsageService) GetPricing() (historymetrics.PricingTable, error) {
	if err := appdata.EnsureAssistantHome(); err != nil {
		return nil, err
	}
	return service.loadPricing()
}

// SavePricing 保存价格配置并返回应用后的报表。
func (service *UsageService) SavePricing(pricing historymetrics.PricingTable) (historymetrics.Report, error) {
	if err := appdata.EnsureAssistantHome(); err != nil {
		return historymetrics.Report{}, err
	}
	if pricing == nil {
		pricing = historymetrics.PricingTable{}
	}
	path := appdata.PricingFilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return historymetrics.Report{}, fmt.Errorf("create pricing directory: %w", err)
	}
	body, err := json.MarshalIndent(pricing, "", "  ")
	if err != nil {
		return historymetrics.Report{}, fmt.Errorf("encode pricing: %w", err)
	}
	if err := os.WriteFile(path, body, 0o644); err != nil {
		return historymetrics.Report{}, fmt.Errorf("write pricing: %w", err)
	}
	return historymetrics.LoadUsageReport(appdata.UsageFilePath(), pricing)
}

func (service *UsageService) loadPricing() (historymetrics.PricingTable, error) {
	body, err := os.ReadFile(appdata.PricingFilePath())
	if err != nil {
		if os.IsNotExist(err) {
			return historymetrics.PricingTable{}, nil
		}
		return nil, fmt.Errorf("read pricing file: %w", err)
	}
	var pricing historymetrics.PricingTable
	if err := json.Unmarshal(body, &pricing); err != nil {
		return nil, fmt.Errorf("decode pricing file: %w", err)
	}
	if pricing == nil {
		pricing = historymetrics.PricingTable{}
	}
	return pricing, nil
}
