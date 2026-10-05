<script setup>
import Button from "@/components/ui/Button.vue";
import Card from "@/components/ui/Card.vue";
import Input from "@/components/ui/Input.vue";
import ContentModal from "@/components/ui/ContentModal.vue";
import { useMessage } from "@/composables/useMessage";
import {
  getUsagePricing,
  getUsageReport,
  saveUsagePricing,
} from "@/services/clientApi";
import { formatCompactInteger, formatInteger } from "@/utils/numberFormat";
import {
  BarController,
  BarElement,
  CategoryScale,
  Chart as ChartJS,
  Legend,
  LinearScale,
  LineController,
  LineElement,
  PointElement,
  Tooltip,
} from "chart.js";
import { computed, onMounted, ref } from "vue";
import { Chart } from "vue-chartjs";

ChartJS.register(
  BarController,
  BarElement,
  CategoryScale,
  Legend,
  LinearScale,
  LineController,
  LineElement,
  PointElement,
  Tooltip,
);

const message = useMessage();

const report = ref(null);
const loading = ref(false);
const loadError = ref("");
const rangeDays = ref(30);
const modelFilter = ref("all");

const pricingOpen = ref(false);
const pricingSaving = ref(false);
const pricingDraft = ref({});
const pricingError = ref("");

const DEFAULT_PRICING = {
  input: 5,
  output: 25,
  cacheRead: 0.5,
  cacheWrite: 6.25,
};

function normalizeNumber(value) {
  const number = Number(value);
  return Number.isFinite(number) ? Math.max(0, number) : 0;
}

function formatTokens(value) {
  return formatInteger(normalizeNumber(value));
}

function formatRate(rate) {
  const number = Number(rate);
  if (!Number.isFinite(number)) {
    return "暂无数据";
  }
  return `${(Math.max(0, Math.min(1, number)) * 100).toFixed(2)}%`;
}

function formatUSD(value) {
  const amount = Number(value);
  if (!Number.isFinite(amount) || amount <= 0) {
    return "$0.00";
  }
  if (amount < 0.01) {
    return "<$0.01";
  }
  return `$${amount.toFixed(2)}`;
}

function formatDuration(ms) {
  const value = normalizeNumber(ms);
  if (value <= 0) {
    return "--";
  }
  if (value < 1000) {
    return `${value}ms`;
  }
  return `${(value / 1000).toFixed(1)}s`;
}

function formatTime(at) {
  const date = new Date(at);
  if (Number.isNaN(date.getTime())) {
    return "--";
  }
  const pad = (n) => String(n).padStart(2, "0");
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`;
}

const summary = computed(() => report.value?.summary ?? {});
const models = computed(() => report.value?.models ?? []);
const events = computed(() => report.value?.events ?? []);
const estimatedTotalCost = computed(() => report.value?.estimatedCostUSD ?? 0);

const rangeOptions = [
  { label: "7天", value: 7 },
  { label: "30天", value: 30 },
  { label: "90天", value: 90 },
  { label: "全部", value: 0 },
];

const daily = computed(() => {
  const points = report.value?.daily ?? [];
  if (!Array.isArray(points) || points.length === 0) {
    return [];
  }
  if (rangeDays.value === 0) {
    return points;
  }
  return points.slice(-rangeDays.value);
});

const modelOptions = computed(() => {
  const seen = new Map();
  for (const item of events.value) {
    const label = item.model || "unknown";
    if (!seen.has(label)) {
      seen.set(label, label);
    }
  }
  return Array.from(seen.values()).sort();
});

const filteredEvents = computed(() => {
  if (modelFilter.value === "all") {
    return events.value;
  }
  return events.value.filter((item) => (item.model || "unknown") === modelFilter.value);
});

const dailyChart = computed(() => ({
  labels: daily.value.map((item) => item.date),
  datasets: [
    {
      type: "bar",
      label: "输入",
      data: daily.value.map((item) => item.inputTokens),
      backgroundColor: "rgba(96, 165, 250, 0.75)",
      stack: "tokens",
      yAxisID: "y",
    },
    {
      type: "bar",
      label: "输出",
      data: daily.value.map((item) => item.outputTokens),
      backgroundColor: "rgba(52, 211, 153, 0.75)",
      stack: "tokens",
      yAxisID: "y",
    },
    {
      type: "bar",
      label: "缓存读",
      data: daily.value.map((item) => item.cacheReadTokens),
      backgroundColor: "rgba(251, 191, 36, 0.75)",
      stack: "tokens",
      yAxisID: "y",
    },
    {
      type: "bar",
      label: "缓存写",
      data: daily.value.map((item) => item.cacheWriteTokens),
      backgroundColor: "rgba(244, 114, 182, 0.75)",
      stack: "tokens",
      yAxisID: "y",
    },
    {
      type: "line",
      label: "估算成本($)",
      data: daily.value.map((item) => Number((item.estimatedCostUSD ?? 0).toFixed(4))),
      borderColor: "rgba(248, 113, 113, 0.9)",
      backgroundColor: "rgba(248, 113, 113, 0.9)",
      tension: 0.3,
      pointRadius: 2,
      yAxisID: "y1",
    },
  ],
}));

const dailyChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  interaction: { mode: "index", intersect: false },
  plugins: {
    legend: { labels: { color: "#a3a3a3", boxWidth: 12 } },
    tooltip: {
      callbacks: {
        label(context) {
          const label = context.dataset.label || "";
          const value = context.parsed.y;
          if (context.dataset.type === "line") {
            return `${label}: $${Number(value).toFixed(4)}`;
          }
          return `${label}: ${formatCompactInteger(value)}`;
        },
      },
    },
  },
  scales: {
    x: {
      stacked: true,
      ticks: { color: "#a3a3a3", maxTicksLimit: 12 },
      grid: { color: "rgba(255,255,255,0.05)" },
    },
    y: {
      stacked: true,
      ticks: {
        color: "#a3a3a3",
        callback(value) {
          return formatCompactInteger(value);
        },
      },
      grid: { color: "rgba(255,255,255,0.05)" },
    },
    y1: {
      position: "right",
      ticks: { color: "rgba(248, 113, 113, 0.9)" },
      grid: { drawOnChartArea: false },
    },
  },
};

const summaryCards = computed(() => [
  { label: "请求数", value: formatTokens(summary.value.providerCallsTotal ?? 0) },
  { label: "输入 Tokens", value: formatTokens(summary.value.promptTokensTotal ?? 0) },
  { label: "输出 Tokens", value: formatTokens(summary.value.outputTokensTotal ?? 0) },
  { label: "缓存命中率", value: formatRate(summary.value.cacheHitRate) },
  { label: "估算总成本", value: formatUSD(estimatedTotalCost.value) },
]);

async function loadReport() {
  loading.value = true;
  loadError.value = "";
  try {
    report.value = await getUsageReport();
  } catch (error) {
    loadError.value = String(error || "加载用量报表失败");
  } finally {
    loading.value = false;
  }
}

function pricingKeyFor(provider, model) {
  const providerText = String(provider || "").trim();
  const modelText = String(model || "").trim() || "unknown";
  return providerText ? `${providerText}::${modelText}` : modelText;
}

async function openPricingEditor() {
  pricingError.value = "";
  const draft = {};
  try {
    const saved = (await getUsagePricing()) || {};
    for (const item of models.value) {
      const key = pricingKeyFor(item.provider, item.model);
      const existing = saved[key];
      draft[key] = {
        label: item.model || key,
        input: existing?.input ?? DEFAULT_PRICING.input,
        output: existing?.output ?? DEFAULT_PRICING.output,
        cacheRead: existing?.cacheRead ?? DEFAULT_PRICING.cacheRead,
        cacheWrite: existing?.cacheWrite ?? DEFAULT_PRICING.cacheWrite,
      };
    }
    for (const [key, value] of Object.entries(saved)) {
      if (draft[key]) {
        continue;
      }
      draft[key] = {
        label: key,
        input: value?.input ?? DEFAULT_PRICING.input,
        output: value?.output ?? DEFAULT_PRICING.output,
        cacheRead: value?.cacheRead ?? DEFAULT_PRICING.cacheRead,
        cacheWrite: value?.cacheWrite ?? DEFAULT_PRICING.cacheWrite,
      };
    }
    pricingDraft.value = draft;
    pricingOpen.value = true;
  } catch (error) {
    message(`读取价格配置失败：${String(error || "")}`);
  }
}

async function handleSavePricing() {
  if (pricingSaving.value) {
    return;
  }
  pricingSaving.value = true;
  pricingError.value = "";
  try {
    const pricing = {};
    for (const [key, value] of Object.entries(pricingDraft.value)) {
      pricing[key] = {
        input: normalizeNumber(value.input),
        output: normalizeNumber(value.output),
        cacheRead: normalizeNumber(value.cacheRead),
        cacheWrite: normalizeNumber(value.cacheWrite),
      };
    }
    report.value = await saveUsagePricing(pricing);
    pricingOpen.value = false;
  } catch (error) {
    pricingError.value = String(error || "保存价格配置失败");
  } finally {
    pricingSaving.value = false;
  }
}

const pricingEntries = computed(() => Object.entries(pricingDraft.value));

onMounted(() => {
  void loadReport();
});
</script>

<template>
  <div class="flex h-full min-h-0 flex-col gap-4 overflow-y-auto px-4 pb-4 pt-4 text-[var(--text-primary)]">
    <Card class="!p-4">
      <div class="flex items-center justify-between gap-3">
        <div>
          <div class="text-[13.5px] font-medium">用量成本报表</div>
          <div class="mt-1 text-[12px] text-[var(--text-secondary)]">按请求、模型与天维度统计 Token 用量与估算成本</div>
        </div>
        <div class="flex items-center gap-2">
          <Button variant="default" :disabled="loading" @click="openPricingEditor">单价设置</Button>
          <Button variant="primary" :disabled="loading" @click="loadReport">
            {{ loading ? "加载中..." : "刷新" }}
          </Button>
        </div>
      </div>
      <div v-if="loadError" class="mt-3 text-sm text-[#f87171]">{{ loadError }}</div>
      <div class="mt-4 grid grid-cols-2 gap-3 md:grid-cols-5">
        <div
          v-for="card in summaryCards"
          :key="card.label"
          class="rounded-lg border border-white/10 bg-white/[0.03] p-3"
        >
          <div class="text-xs text-[#a3a3a3]">{{ card.label }}</div>
          <div class="mt-1 text-lg font-medium">{{ card.value }}</div>
        </div>
      </div>
    </Card>

    <Card class="!p-4">
      <div class="flex items-center justify-between">
        <div class="text-base font-medium">每日趋势</div>
        <div class="flex items-center gap-1">
          <button
            v-for="option in rangeOptions"
            :key="option.value"
            class="rounded-md px-2 py-1 text-sm cursor-pointer"
            :class="rangeDays === option.value ? 'bg-[#333] text-[#e5e5e5]' : 'text-[#a3a3a3] hover:bg-[#2a2a2a]'"
            @click="rangeDays = option.value"
          >
            {{ option.label }}
          </button>
        </div>
      </div>
      <div class="mt-3 h-[260px]">
        <Chart v-if="daily.length > 0" type="bar" :data="dailyChart" :options="dailyChartOptions" />
        <div v-else class="flex h-full items-center justify-center text-sm text-[#a3a3a3]">暂无数据</div>
      </div>
    </Card>

    <Card class="!p-4">
      <div class="text-base font-medium">按模型汇总</div>
      <div class="mt-3 overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="text-left text-[#a3a3a3]">
              <th class="pb-2 pr-3 font-normal">模型</th>
              <th class="pb-2 pr-3 font-normal">Provider</th>
              <th class="pb-2 pr-3 font-normal text-right">请求数</th>
              <th class="pb-2 pr-3 font-normal text-right">输入</th>
              <th class="pb-2 pr-3 font-normal text-right">输出</th>
              <th class="pb-2 pr-3 font-normal text-right">缓存读</th>
              <th class="pb-2 pr-3 font-normal text-right">缓存写</th>
              <th class="pb-2 pr-3 font-normal text-right">命中率</th>
              <th class="pb-2 pr-3 font-normal text-right">估算成本</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="models.length === 0">
              <td colspan="9" class="py-6 text-center text-[#a3a3a3]">暂无数据</td>
            </tr>
            <tr v-for="item in models" :key="item.key" class="border-t border-white/5">
              <td class="py-2 pr-3">{{ item.model || "unknown" }}</td>
              <td class="py-2 pr-3 text-[#a3a3a3]">{{ item.provider || "--" }}</td>
              <td class="py-2 pr-3 text-right">{{ formatTokens(item.providerCalls) }}</td>
              <td class="py-2 pr-3 text-right">{{ formatTokens(item.inputTokens) }}</td>
              <td class="py-2 pr-3 text-right">{{ formatTokens(item.outputTokens) }}</td>
              <td class="py-2 pr-3 text-right">{{ formatTokens(item.cacheReadTokens) }}</td>
              <td class="py-2 pr-3 text-right">{{ formatTokens(item.cacheWriteTokens) }}</td>
              <td class="py-2 pr-3 text-right">{{ formatRate(item.cacheHitRate) }}</td>
              <td class="py-2 pr-3 text-right">{{ formatUSD(item.estimatedCostUSD) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </Card>

    <Card class="!p-4">
      <div class="flex items-center justify-between">
        <div class="text-base font-medium">最近请求</div>
        <select
          v-model="modelFilter"
          class="rounded-md border border-white/10 bg-[#262626] px-2 py-1 text-sm text-[#e5e5e5]"
        >
          <option value="all">全部模型</option>
          <option v-for="name in modelOptions" :key="name" :value="name">{{ name }}</option>
        </select>
      </div>
      <div class="mt-3 max-h-[360px] overflow-auto">
        <table class="w-full text-sm">
          <thead class="sticky top-0 bg-[#1e1e1e]">
            <tr class="text-left text-[#a3a3a3]">
              <th class="pb-2 pr-3 font-normal">时间</th>
              <th class="pb-2 pr-3 font-normal">模型</th>
              <th class="pb-2 pr-3 font-normal">状态</th>
              <th class="pb-2 pr-3 font-normal text-right">耗时</th>
              <th class="pb-2 pr-3 font-normal text-right">输入</th>
              <th class="pb-2 pr-3 font-normal text-right">输出</th>
              <th class="pb-2 pr-3 font-normal text-right">缓存读</th>
              <th class="pb-2 pr-3 font-normal text-right">缓存写</th>
              <th class="pb-2 pr-3 font-normal text-right">命中率</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="filteredEvents.length === 0">
              <td colspan="9" class="py-6 text-center text-[#a3a3a3]">暂无数据</td>
            </tr>
            <tr
              v-for="item in filteredEvents"
              :key="item.eventId"
              class="border-t border-white/5"
              :title="item.errorText || ''"
            >
              <td class="py-2 pr-3 whitespace-nowrap text-[#a3a3a3]">{{ formatTime(item.at) }}</td>
              <td class="py-2 pr-3">{{ item.model || "unknown" }}</td>
              <td class="py-2 pr-3">
                <span
                  class="rounded px-1.5 py-0.5 text-xs"
                  :class="item.status === 'completed' ? 'bg-emerald-500/15 text-emerald-300' : 'bg-red-500/15 text-red-300'"
                >
                  {{ item.status || "--" }}
                </span>
              </td>
              <td class="py-2 pr-3 text-right">{{ formatDuration(item.durationMs) }}</td>
              <td class="py-2 pr-3 text-right">{{ formatTokens(item.inputTokens) }}</td>
              <td class="py-2 pr-3 text-right">{{ formatTokens(item.outputTokens) }}</td>
              <td class="py-2 pr-3 text-right">{{ formatTokens(item.cacheReadTokens) }}</td>
              <td class="py-2 pr-3 text-right">{{ formatTokens(item.cacheWriteTokens) }}</td>
              <td class="py-2 pr-3 text-right">{{ formatRate(item.cacheHitRate) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </Card>

    <ContentModal :open="pricingOpen" title="单价设置（美元 / 百万 Tokens）" size="lg" @close="pricingOpen = false">
      <div class="flex h-full min-h-0 flex-col p-4">
      <div class="mb-3 grid grid-cols-[1fr_repeat(4,90px)] gap-2 text-xs text-[#a3a3a3]">
        <div>模型</div>
        <div class="text-center">输入</div>
        <div class="text-center">输出</div>
        <div class="text-center">缓存读</div>
        <div class="text-center">缓存写</div>
      </div>
      <div class="max-h-[420px] space-y-3 overflow-y-auto">
        <div v-if="pricingEntries.length === 0" class="text-sm text-[#a3a3a3]">暂无可配置模型</div>
        <div
          v-for="[key, value] in pricingEntries"
          :key="key"
          class="grid grid-cols-[1fr_repeat(4,90px)] items-center gap-2"
        >
          <div class="truncate text-sm" :title="key">{{ value.label }}</div>
          <Input v-model="value.input" type="number" min="0" step="0.01" placeholder="输入" />
          <Input v-model="value.output" type="number" min="0" step="0.01" placeholder="输出" />
          <Input v-model="value.cacheRead" type="number" min="0" step="0.01" placeholder="缓存读" />
          <Input v-model="value.cacheWrite" type="number" min="0" step="0.01" placeholder="缓存写" />
        </div>
      </div>
      <div v-if="pricingError" class="mt-3 text-sm text-[#f87171]">{{ pricingError }}</div>
      <div class="mt-4 flex justify-end gap-2">
        <Button variant="default" @click="pricingOpen = false">取消</Button>
        <Button variant="primary" :disabled="pricingSaving" @click="handleSavePricing">
          {{ pricingSaving ? "保存中..." : "保存" }}
        </Button>
      </div>
      </div>
    </ContentModal>
  </div>
</template>
