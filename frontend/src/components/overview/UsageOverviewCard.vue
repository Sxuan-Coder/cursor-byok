<script setup>
import CacheGauge from "@/components/overview/CacheGauge.vue";
import Button from "@/components/ui/Button.vue";
import Tooltip from "@/components/ui/Tooltip.vue";
import { appState } from "@/state/appState";
import {
  deriveHomeMetrics,
  formatRateLabel,
  formatUSD,
} from "@/utils/homeMetrics";
import { formatCompactInteger, formatInteger } from "@/utils/numberFormat";
import dayjs from "dayjs";
import { computed } from "vue";

const props = defineProps({
  metrics: { type: Object, required: true },
  loading: { type: Boolean, default: false },
});

const emit = defineEmits(["toggle-service"]);

const today = computed(() => dayjs().format("YYYY-MM-DD"));

const derived = computed(() =>
  deriveHomeMetrics(props.metrics, {
    includeCacheWriteInHitRate: appState.includeCacheWriteInHitRate,
  }),
);

const serviceRunning = computed(() => appState.serviceRunning);
const serviceBusy = computed(() => appState.serviceBusy);
const serviceButtonText = computed(() => {
  if (serviceBusy.value) {
    return serviceRunning.value ? "关闭中..." : "启动中...";
  }
  return serviceRunning.value ? "关闭服务" : "启动服务";
});

function pairLine(label, value) {
  return `${label}：${formatInteger(value)}`;
}

const cacheTooltip = computed(() => {
  const detail = derived.value;
  const formula = appState.includeCacheWriteInHitRate
    ? "缓存读取 /（缓存读取 + 缓存创建 + 非缓存输入）"
    : "缓存读取 /（缓存读取 + 非缓存输入）";
  return [
    `当前：${formatRateLabel(detail.selectedHitRate)}`,
    `公式：${formula}`,
    `默认 ${formatRateLabel(detail.defaultHitRate)} / 计入创建 ${formatRateLabel(detail.reuseHitRate)}`,
  ].join("\n");
});

const turnsTooltip = computed(() => {
  const detail = derived.value;
  return [
    "按历史记录里扫描到的回合 summary 汇总。",
    "",
    pairLine("总轮次", detail.turnsTotal),
    pairLine("有效轮次", detail.validTurnsTotal),
    pairLine("异常轮次", detail.invalidTurnsTotal),
    `有效占比：${formatRateLabel(detail.validTurnsRate)}`,
  ].join("\n");
});

const tokensTooltip = computed(() => {
  const detail = derived.value;
  return [
    "总请求 Token 包含 Prompt 和模型输出。",
    "",
    pairLine("总请求", detail.requestTokens),
    pairLine("Prompt", detail.promptTokens),
    pairLine("输出推算", detail.completionTokens),
    pairLine("非缓存输入", detail.inputTokens),
    pairLine("缓存读取", detail.cacheReadTokens),
    pairLine("缓存写入", detail.cacheWriteTokens),
  ].join("\n");
});

const costTooltip = computed(() => {
  const detail = derived.value;
  return [
    "按 Claude Opus 价格估算。",
    "",
    `普通输入：${formatInteger(detail.inputTokens)} × $5/1M = ${formatUSD(detail.cost.input)}`,
    `模型输出：${formatInteger(detail.completionTokens)} × $25/1M = ${formatUSD(detail.cost.output)}`,
    `缓存读取：${formatInteger(detail.cacheReadTokens)} × $0.5/1M = ${formatUSD(detail.cost.cacheRead)}`,
    `缓存写入：${formatInteger(detail.cacheWriteTokens)} × $6.25/1M = ${formatUSD(detail.cost.cacheWrite)}`,
    "",
    `合计：${formatUSD(detail.cost.total)}`,
  ].join("\n");
});
</script>

<template>
  <section
    class="rounded-[var(--radius-lg)] border border-[var(--border)] bg-[var(--bg-card)] transition-opacity"
    :class="{ 'opacity-70': loading }"
  >
    <header
      class="flex h-[38px] items-center justify-between gap-3 border-b border-[var(--border-subtle)] px-3.5"
    >
      <div class="flex items-center gap-2">
        <span class="icon-[mdi--chart-timeline-variant] text-[15px] text-[var(--brand)]"></span>
        <h2 class="text-[13px] font-medium text-[var(--text-primary)]">
          今日使用概览
        </h2>
      </div>
      <div class="flex items-center gap-1.5 text-[11.5px] text-[var(--text-faint)]">
        <span class="icon-[mdi--calendar-blank-outline] text-[14px]"></span>
        <span class="font-num">{{ today }}</span>
      </div>
    </header>

    <div class="grid grid-cols-4">
      <div
        class="flex min-w-0 flex-col justify-between border-r border-[var(--border-subtle)] px-3.5 pb-2.5 pt-2.5"
      >
        <div class="flex items-center gap-1 text-[12px] text-[var(--text-secondary)]">
          <span>缓存命中率</span>
          <Tooltip :content="cacheTooltip" />
        </div>
        <div class="flex justify-center pt-1.5">
          <CacheGauge :rate="derived.selectedHitRate" :size="92" />
        </div>
      </div>

      <div
        class="flex min-w-0 flex-col justify-between border-r border-[var(--border-subtle)] px-3.5 pb-2.5 pt-2.5"
      >
        <div class="flex items-center gap-1 text-[12px] text-[var(--text-secondary)]">
          <span>对话轮次</span>
          <Tooltip :content="turnsTooltip" />
        </div>
        <div class="min-w-0">
          <div
            class="truncate font-num text-[23px] leading-none tracking-tight text-[var(--text-primary)]"
            :title="formatInteger(metrics.turnsTotal)"
          >
            {{ formatCompactInteger(metrics.turnsTotal) }}
          </div>
          <div class="mt-2 truncate text-[11.5px] text-[var(--text-faint)]">
            有效 {{ formatCompactInteger(metrics.validTurnsTotal) }} / 异常
            {{ formatCompactInteger(metrics.invalidTurnsTotal) }}
          </div>
        </div>
      </div>

      <div
        class="flex min-w-0 flex-col justify-between border-r border-[var(--border-subtle)] px-3.5 pb-2.5 pt-2.5"
      >
        <div class="flex items-center gap-1 text-[12px] text-[var(--text-secondary)]">
          <span>Token 消耗</span>
          <Tooltip :content="tokensTooltip" />
        </div>
        <div class="min-w-0">
          <div
            class="truncate font-num text-[23px] leading-none tracking-tight text-[var(--text-primary)]"
            :title="formatInteger(metrics.requestTokensTotal)"
          >
            {{ formatCompactInteger(metrics.requestTokensTotal) }}
          </div>
          <div class="mt-2 truncate text-[11.5px] text-[var(--text-faint)]">
            Prompt {{ formatCompactInteger(metrics.promptTokensTotal) }}
          </div>
        </div>
      </div>

      <div
        class="flex min-w-0 flex-col justify-between px-3.5 pb-2.5 pt-2.5"
      >
        <div class="flex items-center gap-1 text-[12px] text-[var(--text-secondary)]">
          <span>价值估算</span>
          <Tooltip :content="costTooltip" />
        </div>
        <div class="min-w-0">
          <div
            class="truncate font-num text-[23px] leading-none tracking-tight text-[var(--text-primary)]"
            :title="formatUSD(derived.cost.total)"
          >
            {{ formatUSD(derived.cost.total) }}
          </div>
          <div class="mt-2 truncate text-[11.5px] text-[var(--text-faint)]">
            缓存读写
            {{ formatUSD(derived.cost.cacheRead + derived.cost.cacheWrite) }}
          </div>
        </div>
      </div>
    </div>

    <footer
      class="flex items-center justify-between gap-3 border-t border-[var(--border-subtle)] px-3.5 py-[6px]"
    >
      <span
        class="flex items-center gap-1.5 text-[12px]"
        :class="serviceRunning ? 'text-[var(--brand)]' : 'text-[var(--text-secondary)]'"
      >
        <span
          class="h-[6px] w-[6px] rounded-full"
          :class="serviceRunning ? 'bg-[var(--brand)]' : 'bg-[var(--text-faint)]'"
        ></span>
        <span>{{ serviceRunning ? "服务运行中" : "服务未启动" }}</span>
      </span>
      <Button
        variant="primary"
        :disabled="serviceBusy"
        @click="emit('toggle-service')"
      >
        <span
          :class="serviceRunning ? 'icon-[mdi--pause]' : 'icon-[mdi--play]'"
          class="text-[15px]"
        ></span>
        <span>{{ serviceButtonText }}</span>
      </Button>
    </footer>
  </section>
</template>