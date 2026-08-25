// 提供商分组纯函数模块：不依赖 Vue 与 Wails 绑定，便于 node:test 直接测试。
// 分组键为「归一化 baseURL + 访问密钥」，与后端渠道身份归一化规则保持一致。

function asText(value) {
  return typeof value === "string" ? value : String(value ?? "");
}

function normalizeProviderBaseURL(value) {
  const text = asText(value).trim();
  if (!text) {
    return "";
  }
  try {
    const parsed = new URL(text);
    if (parsed.protocol !== "http:" && parsed.protocol !== "https:") {
      return text;
    }
    parsed.protocol = parsed.protocol.toLowerCase();
    parsed.hostname = parsed.hostname.toLowerCase();
    const normalized = parsed.toString().replace(/\/+$/, "");
    return normalized || parsed.toString();
  } catch (_error) {
    return text;
  }
}

// buildProviderGroupKey 返回模型适配器所属提供商的分组键。
export function buildProviderGroupKey(adapter) {
  const source = adapter && typeof adapter === "object" ? adapter : {};
  return [
    normalizeProviderBaseURL(source.baseURL),
    asText(source.apiKey).trim(),
  ].join("\n");
}

// buildProviderGroupKeyFromCredentials 用裸凭证构造分组键，供批量更新/删除匹配。
export function buildProviderGroupKeyFromCredentials(baseURL, apiKey) {
  return buildProviderGroupKey({ baseURL, apiKey });
}

// buildProviderGroups 按提供商把模型适配器聚合成分组列表，保持全局首次出现顺序。
// 返回结构：{ key, baseURL, apiKey, types, adapters }。
export function buildProviderGroups(adapters) {
  const groups = [];
  const groupByKey = new Map();
  for (const adapter of Array.isArray(adapters) ? adapters : []) {
    const key = buildProviderGroupKey(adapter);
    let group = groupByKey.get(key);
    if (!group) {
      group = {
        key,
        baseURL: asText(adapter.baseURL).trim(),
        apiKey: asText(adapter.apiKey).trim(),
        adapters: [],
      };
      groupByKey.set(key, group);
      groups.push(group);
    }
    group.adapters.push(adapter);
  }
  for (const group of groups) {
    group.types = Array.from(new Set(group.adapters.map((adapter) => asText(adapter.type))));
  }
  return groups;
}
