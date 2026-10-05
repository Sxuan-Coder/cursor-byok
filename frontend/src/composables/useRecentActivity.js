import { getUsageReport } from "@/services/clientApi";
import { computed, onMounted, ref } from "vue";

// 概览页「最近会话」的数据来源。
// 后端目前没有会话/标题数据（internal/historymetrics 与 internal/cursor 均无会话维度），
// 因此按约定回退为模型排行：以用量报表的按模型聚合为准，并附带最近一次调用时间。
const MAX_ROWS = 6;
const MIN_VALID_YEAR = 2000;

function normalizeCount(value) {
  const number = Number(value);
  return Number.isFinite(number) ? Math.max(0, Math.round(number)) : 0;
}

function normalizeProvider(provider) {
  const text = String(provider ?? "").trim();
  return text.toLowerCase();
}

function providerLabel(provider) {
  const key = normalizeProvider(provider);
  if (key === "anthropic") {
    return "Anthropic";
  }
  if (key === "openai") {
    return "OpenAI";
  }
  return String(provider ?? "").trim() || "未知";
}

function providerBrand(provider) {
  return normalizeProvider(provider) === "anthropic" ? "anthropic" : "openai";
}

function parseTime(value) {
  const time = new Date(value).getTime();
  if (!Number.isFinite(time)) {
    return null;
  }
  return new Date(time).getFullYear() < MIN_VALID_YEAR ? null : time;
}

// formatRelativeTime 输出「刚刚 / N 分钟前 / N 小时前 / 昨天 / N 天前 / M月D日」。
export function formatRelativeTime(value) {
  const time = parseTime(value);
  if (time === null) {
    return "";
  }
  const diff = Date.now() - time;
  const minute = 60 * 1000;
  const hour = 60 * minute;
  const day = 24 * hour;

  if (diff < minute) {
    return "刚刚";
  }
  if (diff < hour) {
    return `${Math.floor(diff / minute)} 分钟前`;
  }
  if (diff < day) {
    return `${Math.floor(diff / hour)} 小时前`;
  }
  if (diff < 2 * day) {
    return "昨天";
  }
  if (diff < 30 * day) {
    return `${Math.floor(diff / day)} 天前`;
  }
  const date = new Date(time);
  return `${date.getMonth() + 1}月${date.getDate()}日`;
}

function buildLeaderboardRows(report) {
  const models = Array.isArray(report?.models) ? report.models : [];
  return models
    .slice()
    .sort((left, right) => normalizeCount(right.providerCalls) - normalizeCount(left.providerCalls))
    .slice(0, MAX_ROWS)
    .map((item) => ({
      id: String(item.key || `${item.provider}::${item.model}`),
      title: String(item.model || "未知模型"),
      tag: providerLabel(item.provider),
      brand: providerBrand(item.provider),
      detail: `${normalizeCount(item.providerCalls)} 次调用`,
      time: formatRelativeTime(item.lastSeenAt),
    }));
}

export function useRecentActivity() {
  const loading = ref(false);
  const error = ref("");
  const rows = ref([]);

  async function refresh() {
    loading.value = true;
    error.value = "";
    try {
      const report = await getUsageReport();
      rows.value = buildLeaderboardRows(report);
    } catch (cause) {
      rows.value = [];
      error.value = String(cause || "加载失败");
    } finally {
      loading.value = false;
    }
  }

  onMounted(() => {
    void refresh();
  });

  return {
    loading,
    error,
    rows: computed(() => rows.value),
    empty: computed(() => !loading.value && rows.value.length === 0),
    refresh,
  };
}