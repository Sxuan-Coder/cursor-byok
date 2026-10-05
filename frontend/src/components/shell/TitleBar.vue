<script setup>
import AppLogo from "@/components/brand/AppLogo.vue";
import { Window } from "@wailsio/runtime";

const props = defineProps({
  title: { type: String, default: "Cursor助手" },
  showWindowControls: { type: Boolean, default: false },
  closeDirectly: { type: Boolean, default: false },
});

async function minimizeWindow() {
  await Window.Minimise();
}

async function toggleMaximizeWindow() {
  await Window.ToggleMaximise();
}

async function closeWindow() {
  if (props.closeDirectly) {
    await Window.Close();
    return;
  }
  // 收起而不是退出，保留托盘常驻能力。
  await new Promise((resolve) => setTimeout(resolve, 200));
  await Window.Hide();
}
</script>

<template>
  <header
    class="drag-region relative flex h-[var(--titlebar-h)] shrink-0 items-center gap-2 border-b border-[var(--border-subtle)] bg-[var(--bg-titlebar)] px-3"
  >
    <span class="flex min-w-0 items-center gap-2">
      <AppLogo :size="15" />
      <span class="text-[12.5px] font-medium text-[var(--text-primary)]">
        {{ title }}
      </span>
      <span class="text-[12px] text-[var(--border-strong)]">|</span>
      <span class="text-[12px] font-medium text-[var(--brand)]">BYOK</span>
      <span class="text-[12px] text-[var(--brand)]">永久免费</span>
    </span>

    <div
      v-if="showWindowControls"
      class="no-drag-region absolute right-2 top-1/2 flex -translate-y-1/2 items-center gap-[2px]"
    >
      <button
        type="button"
        class="flex h-[24px] w-[30px] items-center justify-center rounded-[5px] text-[var(--text-muted)] transition-colors hover:bg-[var(--bg-hover)] hover:text-[var(--text-primary)]"
        aria-label="最小化"
        title="最小化"
        @click="minimizeWindow"
      >
        <span class="icon-[mdi--window-minimize] text-[15px]"></span>
      </button>
      <button
        type="button"
        class="flex h-[24px] w-[30px] items-center justify-center rounded-[5px] text-[var(--text-muted)] transition-colors hover:bg-[var(--bg-hover)] hover:text-[var(--text-primary)]"
        aria-label="最大化"
        title="最大化 / 还原"
        @click="toggleMaximizeWindow"
      >
        <span class="icon-[mdi--window-maximize] text-[14px]"></span>
      </button>
      <button
        type="button"
        class="flex h-[24px] w-[30px] items-center justify-center rounded-[5px] text-[var(--text-muted)] transition-colors hover:bg-[#c0392b] hover:text-white"
        aria-label="关闭"
        title="关闭"
        @click="closeWindow"
      >
        <span class="icon-[mdi--close] text-[16px]"></span>
      </button>
    </div>
  </header>
</template>