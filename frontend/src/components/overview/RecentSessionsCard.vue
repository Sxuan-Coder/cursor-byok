<script setup>
import { useRecentActivity } from "@/composables/useRecentActivity";
import { useRouter } from "vue-router";

const router = useRouter();
const { loading, error, rows, empty } = useRecentActivity();

function viewAll() {
  router.push("/records");
}
</script>

<template>
  <section
    class="rounded-[var(--radius-lg)] border border-[var(--border)] bg-[var(--bg-card)]"
  >
    <header
      class="flex h-[42px] items-center justify-between gap-3 border-b border-[var(--border-subtle)] px-3.5"
    >
      <h2 class="text-[13px] font-medium text-[var(--text-primary)]">最近会话</h2>
      <button
        type="button"
        class="flex items-center gap-0.5 text-[11.5px] text-[var(--text-secondary)] transition-colors hover:text-[var(--brand)]"
        @click="viewAll"
      >
        <span>查看全部</span>
        <span class="icon-[mdi--chevron-right] text-[14px]"></span>
      </button>
    </header>

    <div v-if="loading" class="flex h-[120px] items-center justify-center text-[12px] text-[var(--text-faint)]">
      加载中...
    </div>

    <div
      v-else-if="empty"
      class="flex h-[120px] flex-col items-center justify-center gap-1 px-4 text-center"
    >
      <span class="text-[12px] text-[var(--text-faint)]">暂无可用的会话记录</span>
      <span v-if="error" class="text-[11px] text-[var(--text-faint)]">已回退为模型排行</span>
    </div>

    <ul v-else class="p-2">
      <li v-for="row in rows" :key="row.id">
        <button
          type="button"
          class="flex w-full items-center gap-2.5 rounded-[9px] px-2 py-[7px] text-left transition-colors hover:bg-[var(--bg-hover)]"
          :title="`${row.title} · ${row.tag} · ${row.detail}`"
          @click="viewAll"
        >
          <span
            class="flex h-[28px] w-[28px] shrink-0 items-center justify-center rounded-[8px] border border-[var(--border)] bg-[var(--bg-card-soft)]"
          >
            <span class="icon-[mdi--message-text-outline] text-[14px] text-[var(--text-muted)]"></span>
          </span>
          <span class="flex min-w-0 flex-1 flex-col gap-[4px]">
            <span class="truncate text-[12.5px] text-[var(--text-primary)]">{{ row.title }}</span>
            <span
              class="flex w-fit shrink-0 items-center gap-1 rounded-[5px] border border-[var(--border)] px-1.5 py-[1px] text-[10px] text-[var(--text-muted)]"
            >
              {{ row.tag }}
            </span>
          </span>
          <span class="shrink-0 self-start pt-[2px] text-[11px] text-[var(--text-faint)]">
            {{ row.time }}
          </span>
        </button>
      </li>
    </ul>
  </section>
</template>