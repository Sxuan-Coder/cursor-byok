<script setup>
import AppLogo from "@/components/brand/AppLogo.vue";
import AccountMenu from "@/components/shell/AccountMenu.vue";
import { appState, checkForAppUpdates } from "@/state/appState";
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";

const route = useRoute();
const router = useRouter();

const NAV_ITEMS = [
  { key: "overview", label: "概览", icon: "icon-[mdi--view-dashboard-outline]", path: "/" },
  { key: "records", label: "记录", icon: "icon-[mdi--text-box-outline]", path: "/records" },
  { key: "models", label: "模型", icon: "icon-[mdi--view-grid-outline]", path: "/model-config" },
  { key: "usage", label: "用量", icon: "icon-[mdi--chart-donut]", path: "/usage" },
  { key: "settings", label: "设置", icon: "icon-[mdi--cog-outline]", path: "/settings" },
];

const versionLabel = computed(() => `v${appState.appVersion || "..."}`);

function isActive(path) {
  return route.path === path;
}
</script>

<template>
  <nav
    class="flex h-full w-[var(--sidebar-w)] shrink-0 flex-col border-r border-[var(--border-subtle)] bg-[var(--bg-sidebar)] px-3 pb-3 pt-3"
  >
    <div class="flex items-center gap-2 px-1 pb-3">
      <span class="flex h-[30px] w-[30px] shrink-0 items-center justify-center rounded-[9px] border border-[var(--brand-border)] bg-[var(--brand-soft)]">
        <AppLogo :size="17" />
      </span>
      <span class="flex min-w-0 flex-col leading-tight">
        <span class="truncate text-[13px] font-semibold text-[var(--text-primary)]">
          Cursor助手
        </span>
        <span class="truncate text-[9.5px] text-[var(--text-faint)]">
          BYOK 多模型代理客户端
        </span>
      </span>
    </div>

    <ul class="flex flex-col gap-1">
      <li v-for="item in NAV_ITEMS" :key="item.key">
        <button
          type="button"
          class="relative flex w-full items-center gap-2.5 rounded-[9px] px-3 py-[9px] text-left text-[13px] transition-colors"
          :class="
            isActive(item.path)
              ? 'bg-[var(--brand-soft-strong)] font-medium text-[var(--brand)]'
              : 'text-[var(--text-secondary)] hover:bg-[var(--bg-hover)] hover:text-[var(--text-primary)]'
          "
          :aria-current="isActive(item.path) ? 'page' : undefined"
          @click="router.push(item.path)"
        >
          <span
            v-if="isActive(item.path)"
            class="absolute left-0 top-1/2 h-[16px] w-[3px] -translate-y-1/2 rounded-full bg-[var(--brand)]"
          ></span>
          <span :class="[item.icon, 'shrink-0 text-[17px]']"></span>
          <span class="truncate">{{ item.label }}</span>
        </button>
      </li>
    </ul>

    <div class="mt-auto flex flex-col gap-2 pt-3">
      <AccountMenu placement="top-start" variant="row" />

      <button
        type="button"
        class="flex items-center gap-1.5 rounded-[7px] px-1 py-0.5 text-[11px] text-[var(--text-faint)] transition-colors hover:text-[var(--text-secondary)] disabled:opacity-60"
        :disabled="appState.updateState === 'checking'"
        @click="checkForAppUpdates()"
      >
        <span>{{ versionLabel }}</span>
        <span>检查更新</span>
      </button>
    </div>
  </nav>
</template>
