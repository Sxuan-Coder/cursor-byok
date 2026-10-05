<script setup>
import HeroBanner from "@/components/overview/HeroBanner.vue";
import ModelConfigCard from "@/components/overview/ModelConfigCard.vue";
import QuickActionsCard from "@/components/overview/QuickActionsCard.vue";
import UsageOverviewCard from "@/components/overview/UsageOverviewCard.vue";
import { useActiveModelIds } from "@/composables/useActiveModelIds";
import { useMessage } from "@/composables/useMessage";
import { appState, toggleService } from "@/state/appState";
import { useRouter } from "vue-router";

const router = useRouter();
const message = useMessage();
const activeIds = useActiveModelIds();

async function handleToggleService() {
  const result = await toggleService();
  if (!result.ok) {
    const detail = String(result.error || "服务错误").trim() || "服务错误";
    message(`服务操作失败：${detail}`);
  }
}

function openModelConfig() {
  router.push("/model-config");
}
</script>

<template>
  <div class="flex h-full min-h-0 flex-col gap-3 overflow-y-auto px-3.5 pb-3 pt-2.5">
    <HeroBanner />

    <UsageOverviewCard
      :metrics="appState.homeMetrics"
      :loading="appState.homeMetricsLoading"
      @toggle-service="handleToggleService"
    />

    <ModelConfigCard
      :adapters="appState.modelAdapters"
      :active-ids="activeIds"
      @manage="openModelConfig"
    />

    <QuickActionsCard />
  </div>
</template>