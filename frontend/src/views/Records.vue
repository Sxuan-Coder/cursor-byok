<script setup>
import Button from "@/components/ui/Button.vue";
import { useMessage } from "@/composables/useMessage";
import { getUsageReport } from "@/services/clientApi";
import { formatCompactInteger } from "@/utils/numberFormat";
import { formatDateTime, formatDuration, formatRate, formatTokens, statusLabel } from "@/utils/usageFormat";
import { computed, onMounted, ref } from "vue";

const MESSAGE_DURATION = 2400;

const message = useMessage();

const report = ref(null);
const loading = ref(false);
const loadError = ref("");
const modelFilter = ref("all");
const kindFilter = ref("all");

const events = computed(() => {
  const list = report.value?.events;
  return Array.isArray(list) ? list : [];
});

const modelOptions = computed(() => {
  const seen = new Set();
  for (const item of events.value) {
    seen.add(String(item.model || "unknown"));
  }
  return Array.from(seen).sort();
});

const kindOptions = computed(() => {
  const seen = new Set();
  for (const item of events.value) {
    const kind = String(item.kind || "").trim();
    if (kind) {
      seen.add(kind);
    }
  }
  return Array.from(seen).sort();
});

const filteredEvents = computed(() =>
  events.value.filter((item) => {
    if (modelFilter.value !== "all" && String(item.model || "unknown") !== modelFilter.value) {
      return false;
    }
    if (kindFilter.value !== "all" && String(item.kind || "") !== kindFilter.value) {
      return false;
    }
    return true;
  }),
);

const summaryLine = computed(
  () => `共 ${formatTokens(filteredEvents.value.length)} 条记录`,
);

async function loadReport({ notify = false } = {}) {
  loading.value = true;
  loadError.value = "";
  try {
    report.value = await getUsageReport();
    if (notify) {
      message("记录已刷新", { duration: MESSAGE_DURATION });
    }
  } catch (error) {
    loadError.value = String(error || "加载记录失败");
  } finally {
    loading.value = false;
  }
}

const SELECT_CLASS =
  "rounded-[8px] border border-[var(--border)] bg-[var(--bg-card-soft)] px-2 py-1.5 text-[12px] text-[var(--text-secondary)] outline-none";

onMounted(() => {
  void loadReport();
});
</script>

<template>
  <div class="flex h-full min-h-0 flex-col gap-4 overflow-y-auto px-4 pb-4 pt-4">
    <section class="rounded-[var(--radius-lg)] border border-[var(--border)] bg-[var(--bg-card)]">
      <header
        class="flex h-[46px] items-center justify-between gap-3 border-b border-[var(--border-subtle)] px-4"
      >
        <div class="flex min-w-0 items-center gap-2">
          <span class="icon-[mdi--text-box-outline] text-[16px] text-[var(--brand)]"></span>
          <h2 class="text-[13px] font-medium text-[var(--text-primary)]">最近请求记录</h2>
          <span class="truncate text-[11.5px] text-[var(--text-faint)]">{{ summaryLine }}</span>
        </div>
        <div class="flex shrink-0 items-center gap-2">
          <select v-model="modelFilter" :class="SELECT_CLASS" aria-label="按模型筛选">
            <option value="all">全部模型</option>
            <option v-for="name in modelOptions" :key="name" :value="name">{{ name }}</option>
          </select>
          <select v-if="kindOptions.length > 1" v-model="kindFilter" :class="SELECT_CLASS" aria-label="按类型筛选">
            <option value="all">全部类型</option>
            <option v-for="kind in kindOptions" :key="kind" :value="kind">{{ kind }}</option>
          </select>
          <Button variant="default" :disabled="loading" @click="loadReport({ notify: true })">
            {{ loading ? "加载中..." : "刷新" }}
          </Button>
        </div>
      </header>

      <div v-if="loadError" class="px-4 py-3 text-[12px] text-[var(--danger)]">{{ loadError }}</div>

      <div class="max-h-[calc(100vh-220px)] overflow-auto px-4 pb-4">
        <table class="w-full text-[12px]">
          <thead class="sticky top-0 z-10">
            <tr class="bg-[var(--bg-card)] text-left text-[var(--text-faint)]">
              <th class="whitespace-nowrap py-2.5 pr-3 font-normal">时间</th>
              <th class="whitespace-nowrap py-2.5 pr-3 font-normal">模型</th>
              <th class="whitespace-nowrap py-2.5 pr-3 font-normal">Provider</th>
              <th class="whitespace-nowrap py-2.5 pr-3 font-normal">状态</th>
              <th class="whitespace-nowrap py-2.5 pr-3 text-right font-normal">耗时</th>
              <th class="whitespace-nowrap py-2.5 pr-3 text-right font-normal">输入</th>
              <th class="whitespace-nowrap py-2.5 pr-3 text-right font-normal">输出</th>
              <th class="whitespace-nowrap py-2.5 pr-3 text-right font-normal">缓存读</th>
              <th class="whitespace-nowrap py-2.5 pr-3 text-right font-normal">缓存写</th>
              <th class="whitespace-nowrap py-2.5 text-right font-normal">命中率</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="filteredEvents.length === 0">
              <td colspan="10" class="py-10 text-center text-[var(--text-faint)]">
                暂无记录
              </td>
            </tr>
            <tr
              v-for="item in filteredEvents"
              :key="item.eventId"
              class="border-t border-[var(--border-subtle)] text-[var(--text-secondary)]"
              :title="item.errorText || ''"
            >
              <td class="whitespace-nowrap py-2.5 pr-3 font-num">{{ formatDateTime(item.at) }}</td>
              <td class="py-2.5 pr-3 text-[var(--text-primary)]">{{ item.model || "unknown" }}</td>
              <td class="py-2.5 pr-3">{{ item.provider || "--" }}</td>
              <td class="py-2.5 pr-3">
                <span
                  class="rounded-[5px] px-1.5 py-[2px] text-[11px]"
                  :class="
                    item.status === 'completed'
                      ? 'bg-[var(--brand-soft)] text-[var(--brand)]'
                      : 'bg-[var(--danger-soft)] text-[var(--danger)]'
                  "
                >
                  {{ statusLabel(item.status) }}
                </span>
              </td>
              <td class="py-2.5 pr-3 text-right font-num">{{ formatDuration(item.durationMs) }}</td>
              <td class="py-2.5 pr-3 text-right font-num">{{ formatTokens(item.inputTokens) }}</td>
              <td class="py-2.5 pr-3 text-right font-num">{{ formatTokens(item.outputTokens) }}</td>
              <td class="py-2.5 pr-3 text-right font-num">{{ formatTokens(item.cacheReadTokens) }}</td>
              <td class="py-2.5 pr-3 text-right font-num">{{ formatTokens(item.cacheWriteTokens) }}</td>
              <td class="py-2.5 text-right font-num">{{ formatRate(item.cacheHitRate) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <div class="flex items-center justify-between px-1 text-[11.5px] text-[var(--text-faint)]">
      <span>提示：记录来源于本地历史用量汇总，最多保留最近若干条明细。</span>
      <span v-if="events.length">
        过滤后 {{ formatCompactInteger(filteredEvents.length) }} / 共
        {{ formatCompactInteger(events.length) }}
      </span>
    </div>
  </div>
</template>