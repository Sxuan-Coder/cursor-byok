<script setup>
import SideNav from "@/components/shell/SideNav.vue";
import TitleBar from "@/components/shell/TitleBar.vue";
import TopNav from "@/components/shell/TopNav.vue";
import { syncServiceState } from "@/state/appState";
import { isWindows } from "@/utils/isWindows";
import { computed, onMounted, onUnmounted } from "vue";
import { useRoute } from "vue-router";

const route = useRoute();
const title = computed(() => route.meta.title ?? "Cursor助手");
// 所有路由都支持主窗口内导航，关闭一律收起而不是退出，避免误退出应用（退出走托盘）。
const directlyClose = computed(() => route.meta.directlyClose === true);

const SERVICE_POLL_MS = 10000;
let serviceTimer = null;

onMounted(() => {
  serviceTimer = window.setInterval(() => {
    void syncServiceState().catch(() => {});
  }, SERVICE_POLL_MS);
});

onUnmounted(() => {
  if (serviceTimer) {
    window.clearInterval(serviceTimer);
    serviceTimer = null;
  }
});
</script>

<template>
  <div class="flex h-screen w-screen flex-col overflow-hidden bg-[var(--bg-app)] text-[var(--text-primary)]">
    <TitleBar
      :title="title"
      :show-window-controls="isWindows"
      :close-directly="directlyClose"
    />
    <div class="flex min-h-0 flex-1">
      <SideNav />
      <div class="flex min-w-0 flex-1 flex-col">
        <TopNav />
        <main class="flex min-h-0 min-w-0 flex-1 flex-col bg-[var(--bg-content)]">
          <router-view />
        </main>
      </div>
    </div>
  </div>
</template>