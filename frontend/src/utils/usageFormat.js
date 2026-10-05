import { formatInteger } from "@/utils/numberFormat";

// 用量/记录页共用的展示格式化，避免同一口径在多处重复实现。
export function normalizeAmount(value) {
  const number = Number(value);
  return Number.isFinite(number) ? Math.max(0, number) : 0;
}

export function formatTokens(value) {
  return formatInteger(normalizeAmount(value));
}

export function formatRate(rate) {
  const number = Number(rate);
  if (!Number.isFinite(number)) {
    return "暂无数据";
  }
  return `${(Math.max(0, Math.min(1, number)) * 100).toFixed(2)}%`;
}

export function formatUSD(value) {
  const amount = Number(value);
  if (!Number.isFinite(amount) || amount <= 0) {
    return "$0.00";
  }
  if (amount < 0.01) {
    return "<$0.01";
  }
  return `$${amount.toFixed(2)}`;
}

export function formatDuration(ms) {
  const value = normalizeAmount(ms);
  if (value <= 0) {
    return "--";
  }
  if (value < 1000) {
    return `${value}ms`;
  }
  return `${(value / 1000).toFixed(1)}s`;
}

export function formatDateTime(at) {
  const date = new Date(at);
  if (Number.isNaN(date.getTime())) {
    return "--";
  }
  const pad = (value) => String(value).padStart(2, "0");
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`;
}

// statusLabel 把后端事件状态映射为界面文案。
export function statusLabel(status) {
  switch (String(status || "").trim()) {
    case "completed":
      return "成功";
    case "failed":
    case "error":
      return "失败";
    default:
      return String(status || "").trim() || "--";
  }
}