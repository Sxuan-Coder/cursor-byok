<script setup>
import ProviderMark from "@/components/brand/ProviderMark.vue";
import { providerTypeLabel } from "@/utils/providerLabel";
import { computed } from "vue";

const props = defineProps({
  adapters: { type: Array, required: true },
  activeIds: { type: Object, default: () => new Set() },
});

defineEmits(["manage"]);

function isActive(adapter) {
  return props.activeIds instanceof Set ? props.activeIds.has(adapter.id) : false;
}

const rows = computed(() => props.adapters);
</script>

<template>
  <section class="rounded-[var(--radius-lg)] border border-[var(--border)] bg-[var(--bg-card)]">
    <header
      class="flex h-[38px] items-center justify-between gap-3 border-b border-[var(--border-subtle)] px-3.5"
    >
      <div class="flex items-center gap-2">
        <span class="icon-[mdi--cube-outline] text-[15px] text-[var(--brand)]"></span>
        <h2 class="text-[13px] font-medium text-[var(--text-primary)]">当前模型配置</h2>
      </div>
      <button
        type="button"
        class="text-[12px] text-[var(--text-secondary)] transition-colors hover:text-[var(--brand)]"
        @click="$emit('manage')"
      >
        管理配置
      </button>
    </header>

    <div class="p-2">
      <div
        v-if="rows.length === 0"
        class="flex h-[120px] items-center justify-center rounded-[10px] border border-dashed border-[var(--border)] text-[12.5px] text-[var(--text-faint)]"
      >
        尚未配置任何模型渠道
      </div>

      <div v-else class="grid grid-cols-3 gap-1.5">
        <article
          v-for="adapter in rows"
          :key="adapter.id"
          class="rounded-[10px] border bg-[var(--bg-card-soft)] px-2.5 py-2 transition-colors"
          :class="
            isActive(adapter)
              ? 'border-[var(--brand-border)] bg-[var(--brand-soft)]'
              : 'border-[var(--border)] hover:border-[var(--border-strong)]'
          "
        >
          <div class="flex items-center justify-between gap-1.5">
            <div class="flex min-w-0 items-center gap-1.5">
              <ProviderMark :brand="adapter.type" :size="13" />
              <span
                class="truncate text-[11.5px] font-medium tracking-[-0.01em] text-[var(--text-primary)]"
                :title="adapter.displayName"
              >
                {{ adapter.displayName }}
              </span>
            </div>
            <span
              class="flex shrink-0 items-center gap-[2px] rounded-[999px] border px-[3px] py-[2px] text-[9px]"
              :class="
                adapter.type === 'anthropic'
                  ? 'border-[#4a3a33] text-[#d97757]'
                  : 'border-[var(--border-strong)] text-[#c9c9c9]'
              "
            >
              <ProviderMark :brand="adapter.type" :size="8" />
              <span>{{ providerTypeLabel(adapter.type) }}</span>
            </span>
          </div>
          <div
            class="mt-0.5 truncate text-[11px] text-[var(--text-faint)]"
            :title="adapter.modelID"
          >
            {{ adapter.modelID }}
          </div>
          <div
            v-if="isActive(adapter)"
            class="mt-1.5 flex items-center gap-1 text-[10.5px] text-[var(--brand)]"
          >
            <span class="h-[5px] w-[5px] rounded-full bg-[var(--brand)]"></span>
            <span>使用中</span>
          </div>
        </article>
      </div>
    </div>
  </section>
</template>