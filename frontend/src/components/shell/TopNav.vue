<script setup>
import LocaleSelect from "@/components/LocaleSelect.vue";
import AccountMenu from "@/components/shell/AccountMenu.vue";
import { useCursorAccount } from "@/composables/useCursorAccount";
import { useMessage } from "@/composables/useMessage";
import { useLocale } from "@/i18n/runtime";
import {
  appState,
  appViewState,
  syncHomeMetrics,
  syncServiceState,
} from "@/state/appState";
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";

const route = useRoute();
const router = useRouter();
const message = useMessage();
const { signedIn } = useCursorAccount();
const { locale } = useLocale();

const TABS = [
  { key: "overview", label: "概览", icon: "icon-[mdi--view-dashboard-outline]", path: "/" },
  { key: "records", label: "记录", icon: "icon-[mdi--text-box-outline]", path: "/records" },
  { key: "models", label: "模型", icon: "icon-[mdi--view-grid-outline]", path: "/model-config" },
  { key: "usage", label: "用量", icon: "icon-[mdi--chart-donut]", path: "/usage" },
  { key: "settings", label: "设置", icon: "icon-[mdi--cog-outline]", path: "/settings" },
];

const statusText = computed(() => appViewState.serviceStatusText);
const serviceRunning = computed(() => appState.serviceRunning);
const refreshing = computed(() => appState.homeMetricsLoading);

// 顶栏导航与侧边栏重复，窗口较窄或非中文文案更长时只留图标，
// 避免把右侧操作区挤出去（侧边栏仍提供带文字的导航）。
// tab 文字可截断，直接比较 scrollWidth 会测不出溢出，所以先让 tab 按自然宽度
// 展开测量，再把结果写回状态；两次 DOM 更新都在同一微任务批次内，不会出现闪烁。
const navRef = ref(null);
const compactTabs = ref(false);
const probingTabs = ref(false);
let navResizeObserver = null;

async function fitTabs() {
  compactTabs.value = false;
  probingTabs.value = true;
  await nextTick();

  const nav = navRef.value;
  if (nav && nav.clientWidth > 0) {
    compactTabs.value = nav.scrollWidth > nav.clientWidth + 1;
  }
  probingTabs.value = false;
  await nextTick();
}

onMounted(() => {
  void fitTabs();
  navResizeObserver = new ResizeObserver(() => {
    void fitTabs();
  });
  if (navRef.value) {
    navResizeObserver.observe(navRef.value);
  }
});

onBeforeUnmount(() => {
  navResizeObserver?.disconnect();
  navResizeObserver = null;
});

watch(locale, () => {
  void fitTabs();
});

function isActive(path) {
  return route.path === path;
}

async function handleRefresh() {
  const [service, metrics] = await Promise.allSettled([syncServiceState(), syncHomeMetrics()]);
  const metricsFailed = metrics.status === "fulfilled" && metrics.value?.ok === false;
  if (service.status === "rejected" || metricsFailed) {
    message("刷新失败：服务状态或用量数据获取异常");
  }
}
</script>

<template>
  <header
    class="flex h-[var(--topbar-h)] shrink-0 items-center justify-between gap-3 border-b border-[var(--border-subtle)] bg-[var(--bg-content)] px-4"
  >
    <nav ref="navRef" class="flex min-w-0 flex-1 items-center gap-1.5 overflow-hidden">
      <button
        v-for="tab in TABS"
        :key="tab.key"
        type="button"
        class="flex items-center gap-1.5 rounded-[9px] border px-3 py-[6px] text-[12.5px] transition-colors"
        :title="tab.label"
        :aria-label="tab.label"
        :class="[
          probingTabs ? 'shrink-0' : 'min-w-0',
          isActive(tab.path)
            ? 'border-[var(--brand-border)] bg-[var(--brand-soft-strong)] font-medium text-[var(--brand)]'
            : 'border-transparent text-[var(--text-secondary)] hover:bg-[var(--bg-hover)] hover:text-[var(--text-primary)]',
        ]"
        :aria-current="isActive(tab.path) ? 'page' : undefined"
        @click="router.push(tab.path)"
      >
        <span :class="[tab.icon, 'text-[15px] shrink-0']"></span>
        <span v-if="!compactTabs" class="truncate">{{ tab.label }}</span>
      </button>
    </nav>

    <div class="flex shrink-0 items-center gap-1.5">
      <AccountMenu
        v-if="!signedIn"
        placement="bottom-end"
        variant="pill"
        title="登录后可获得 Cursor 插件 Skill 以及 MCP 商店的官方权限"
      />

      <span
        class="flex items-center gap-1.5 rounded-[8px] border px-2.5 py-[5px] text-[11.5px]"
        :class="
          serviceRunning
            ? 'border-[var(--brand-border)] bg-[var(--brand-soft)] text-[var(--brand)]'
            : 'border-[var(--border-strong)] bg-[var(--bg-elevated)] text-[var(--text-secondary)]'
        "
      >
        <span
          class="h-[6px] w-[6px] rounded-full"
          :class="serviceRunning ? 'bg-[var(--brand)]' : 'bg-[var(--text-faint)]'"
        ></span>
        <span>{{ statusText }}</span>
      </span>

      <button
        type="button"
        class="flex h-[30px] w-[30px] items-center justify-center rounded-[8px] text-[var(--text-secondary)] transition-colors hover:bg-[var(--bg-hover)] hover:text-[var(--text-primary)] disabled:opacity-60"
        :disabled="refreshing"
        aria-label="刷新"
        title="刷新"
        @click="handleRefresh"
      >
        <span
          class="icon-[mdi--refresh] text-[16px]"
          :class="{ 'animate-spin': refreshing }"
        ></span>
      </button>

      <LocaleSelect
        density="compact"
        wrapper-class="w-[96px] shrink-0"
        border
        aria-label="界面语言"
        placeholder="语言"
      />
    </div>
  </header>
</template>
