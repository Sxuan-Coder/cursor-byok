<script setup>
import ProviderMark from "@/components/brand/ProviderMark.vue";
import { useActiveModelIds } from "@/composables/useActiveModelIds";
import { appState } from "@/state/appState";
import { providerTypeBrand, providerTypeLabel } from "@/utils/providerLabel";
import { computed } from "vue";
import { useRouter } from "vue-router";

const MAX_ROWS = 6;

const router = useRouter();
const activeIds = useActiveModelIds();

const rows = computed(() => appState.modelAdapters.slice(0, MAX_ROWS));

function isActive(adapter) {
  return activeIds.value.has(adapter.id);
}

function openModelConfig(openEditor = false) {
  router.push(
    openEditor
      ? { path: "/model-config", query: { new: "1" } }
      : { path: "/model-config" },
  );
}
</script>

<template>
  <section class="rounded-[var(--radius-lg)] border border-[var(--border)] bg-[var(--bg-card)]">
    <header
      class="flex h-[42px] items-center justify-between gap-3 border-b border-[var(--border-subtle)] px-3.5"
    >
      <h2 class="text-[13px] font-medium text-[var(--text-primary)]">常用配置</h2>
      <button
        type="button"
        class="flex h-[22px] w-[22px] items-center justify-center rounded-[6px] text-[var(--text-secondary)] transition-colors hover:bg-[var(--bg-hover)] hover:text-[var(--brand)]"
        aria-label="新增模型配置"
        title="新增模型配置"
        @click="openModelConfig(true)"
      >
        <span class="icon-[mdi--plus] text-[16px]"></span>
      </button>
    </header>

    <div
      v-if="rows.length === 0"
      class="flex h-[120px] items-center justify-center px-4 text-center text-[12px] text-[var(--text-faint)]"
    >
      尚未配置任何模型渠道
    </div>

    <ul v-else class="p-2">
      <li v-for="adapter in rows" :key="adapter.id">
        <button
          type="button"
          class="flex w-full items-center gap-2.5 rounded-[9px] px-2 py-[7px] text-left transition-colors hover:bg-[var(--bg-hover)]"
          :title="adapter.displayName"
          @click="openModelConfig()"
        >
          <span
            class="flex h-[26px] w-[26px] shrink-0 items-center justify-center rounded-[7px] border border-[var(--border)] bg-[var(--bg-card-soft)]"
          >
            <ProviderMark :brand="providerTypeBrand(adapter.type)" :size="14" />
          </span>
          <span class="flex min-w-0 flex-1 flex-col gap-[3px]">
            <span class="truncate text-[12.5px] text-[var(--text-primary)]">
              {{ adapter.displayName }}
            </span>
            <span class="truncate text-[10.5px] text-[var(--text-faint)]">
              {{ providerTypeLabel(adapter.type) }}
            </span>
          </span>
          <span
            v-if="isActive(adapter)"
            class="flex shrink-0 items-center gap-1 text-[10.5px] text-[var(--brand)]"
          >
            <span class="h-[5px] w-[5px] rounded-full bg-[var(--brand)]"></span>
            <span>使用中</span>
          </span>
        </button>
      </li>
    </ul>
  </section>
</template>