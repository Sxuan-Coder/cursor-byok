// 首页统计派生逻辑：把后端原始 summary 换算为概览卡展示所需的口径与成本估算。
export const TOKEN_PRICE_PER_MILLION = {
  input: 5,
  output: 25,
  cacheRead: 0.5,
  cacheWrite: 6.25,
};

export function normalizeCount(value) {
  const number = Number(value);
  if (!Number.isFinite(number)) {
    return 0;
  }
  return Math.round(number);
}

function ratio(numerator, denominator) {
  const top = normalizeCount(numerator);
  const bottom = normalizeCount(denominator);
  if (bottom <= 0) {
    return null;
  }
  return top / bottom;
}

function priceTokens(tokens, pricePerMillion) {
  return (normalizeCount(tokens) / 1_000_000) * pricePerMillion;
}

export function formatRateLabel(value) {
  const rate = Number(value);
  if (!Number.isFinite(rate)) {
    return "暂无数据";
  }
  return `${(Math.max(0, Math.min(1, rate)) * 100).toFixed(2)}%`;
}

export function formatUSD(value) {
  const amount = Number(value);
  if (!Number.isFinite(amount)) {
    return "$0.00";
  }
  if (amount > 0 && amount < 0.01) {
    return "<$0.01";
  }
  return `$${amount.toFixed(2)}`;
}

// deriveHomeMetrics 依据缓存统计口径计算命中率与估算成本。
export function deriveHomeMetrics(metrics, { includeCacheWriteInHitRate = false } = {}) {
  const source = metrics && typeof metrics === "object" ? metrics : {};
  const cacheReadTokens = normalizeCount(source.cacheReadTokens);
  const cacheWriteTokens = normalizeCount(source.cacheWriteTokens);
  const promptTokens = normalizeCount(source.promptTokensTotal);
  const requestTokens = normalizeCount(source.requestTokensTotal);

  const inputTokens = Math.max(0, promptTokens - cacheReadTokens - cacheWriteTokens);
  const completionTokens = Math.max(0, requestTokens - promptTokens);

  const defaultHitRate = ratio(cacheReadTokens, cacheReadTokens + inputTokens);
  const reuseHitRate = ratio(
    cacheReadTokens,
    cacheReadTokens + cacheWriteTokens + inputTokens,
  );
  const selectedHitRate = includeCacheWriteInHitRate ? reuseHitRate : defaultHitRate;

  const cost = {
    input: priceTokens(inputTokens, TOKEN_PRICE_PER_MILLION.input),
    output: priceTokens(completionTokens, TOKEN_PRICE_PER_MILLION.output),
    cacheRead: priceTokens(cacheReadTokens, TOKEN_PRICE_PER_MILLION.cacheRead),
    cacheWrite: priceTokens(cacheWriteTokens, TOKEN_PRICE_PER_MILLION.cacheWrite),
  };
  cost.total = cost.input + cost.output + cost.cacheRead + cost.cacheWrite;

  return {
    turnsTotal: normalizeCount(source.turnsTotal),
    validTurnsTotal: normalizeCount(source.validTurnsTotal),
    invalidTurnsTotal: normalizeCount(source.invalidTurnsTotal),
    requestTokens,
    promptTokens,
    completionTokens,
    inputTokens,
    cacheReadTokens,
    cacheWriteTokens,
    defaultHitRate,
    reuseHitRate,
    selectedHitRate,
    validTurnsRate: ratio(source.validTurnsTotal, source.turnsTotal),
    cost,
  };
}