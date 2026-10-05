// 模型渠道类型在前端的统一展示名，避免各组件重复定义。
export function providerTypeLabel(type) {
  return type === "anthropic" ? "Anthropic" : "OpenAI";
}

// providerTypeBrand 归一化为 ProviderMark 支持的品牌标识。
export function providerTypeBrand(type) {
  return type === "anthropic" ? "anthropic" : "openai";
}